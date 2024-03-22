package aggregation

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/decentrio/soro-book/lib/log"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/ingest"
	backends "github.com/stellar/go/ingest/ledgerbackend"
	stellar_log "github.com/stellar/go/support/log"

	"github.com/decentrio/soro-book/config"
	db "github.com/decentrio/soro-book/database/handlers"
	"github.com/decentrio/soro-book/database/models"
	"github.com/decentrio/soro-book/lib/service"
)

const (
	txQueueSize = 1000
	prepareStep = 64
)

type Aggregation struct {
	ctx context.Context

	log *stellar_log.Entry

	config backends.CaptiveCoreConfig

	service.BaseService

	cfg *config.AggregationConfig

	backend *backends.CaptiveStellarCore

	// txQueue channel for trigger new tx
	txQueue chan ingest.LedgerTransaction

	// isReSync is flag represent if services is
	// re-synchronize
	isReSync bool

	// subscribe services
	sequence uint32

	db *db.DBHandler
}

// AggregationOption sets an optional parameter on the State.
type AggregationOption func(*Aggregation)

func NewAggregation(
	cfg *config.AggregationConfig,
	logger log.Logger,
	options ...AggregationOption,
) *Aggregation {
	as := &Aggregation{
		cfg:      cfg,
		txQueue:  make(chan ingest.LedgerTransaction, txQueueSize),
		isReSync: false,
	}

	as.BaseService = *service.NewBaseService("Aggregation", as)
	for _, opt := range options {
		opt(as)
	}

	as.BaseService.SetLogger(logger.With("module", "aggregation"))

	as.db = db.NewDBHandler()

	as.ctx = context.Background()
	as.log = stellar_log.New()
	as.log.SetLevel(logrus.ErrorLevel)
	Config.Log = as.log

	as.sequence = uint32(100_000)

	var err error
	as.backend, err = backends.NewCaptive(Config)
	panicIf(err)

	return as
}

func (as *Aggregation) OnStart() error {
	as.Logger.Info("Start")
	go as.dataProcessing()
	// Note that when using goroutines, you need to be careful to ensure that no
	// race conditions occur when accessing the txQueue.
	go as.aggregation()
	return nil
}

func (as *Aggregation) OnStop() error {
	as.Logger.Info("Stop")
	as.backend.Close()
	return nil
}

// aggregation process
func (as *Aggregation) dataProcessing() {
	for {
		// Block until state have sync successful
		if as.isReSync {
			continue
		}

		select {
		// Receive a new tx
		case tx := <-as.txQueue:
			as.handleReceiveTx(tx)
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		}
	}
}

// handleReceiveTx
func (as *Aggregation) handleReceiveTx(tx ingest.LedgerTransaction) {
	// Check if tx metadata is v3
	txMetaV3, ok := tx.UnsafeMeta.GetV3()
	if !ok {
		as.Logger.Error("receive tx not a metadata v3")
		return
	}

	if txMetaV3.SorobanMeta == nil {
		as.Logger.Error("nil soroban meta")
		return
	}

	event := &models.Event{
		// Type: txMetaV3.SorobanMeta.Events,
		Ledger: int32(tx.Index),
	}

	// for _, event := range txMetaV3.SorobanMeta.Events {
	// contractEvent, ok := event.Body.GetV0()
	// if !ok {
	// 	as.Logger.Error("Error Soroban event")
	// 	return
	// }
	// topics := contractEvent.Topics
	// }

	_, err := as.db.CreateEvent(event)
	if err != nil {
		as.Logger.Error(err.Error())
	}
}

func (as *Aggregation) aggregation() {
	for {
		select {
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		default:
			as.getNewTx()
		}
	}
}

func (as *Aggregation) getNewTx() {
	from := as.sequence
	to := as.sequence + prepareStep
	ledgerRange := backends.BoundedRange(from, to)
	err := as.backend.PrepareRange(as.ctx, ledgerRange)
	if err != nil {
		//"is greater than max available in history archives"
		err = pauseWaitLedger(as.config, err)
		if err != nil {
			as.Logger.Error(err.Error())
		}

		return
	}
	for seq := from; seq < to; seq++ {
		txReader, err := ingest.NewLedgerTransactionReader(
			as.ctx, as.backend, Config.NetworkPassphrase, seq,
		)
		panicIf(err)
		defer txReader.Close()

		// Read each transaction within the ledger, extract its operations, and
		// accumulate the statistics we're interested in.
		for {
			tx, err := txReader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				as.Logger.Error(err.Error())
			}

			if tx.Result.Successful() {
				as.Logger.Info(fmt.Sprintf("tx received %s", tx.Result.TransactionHash.HexString()))
				go func(txi ingest.LedgerTransaction) {
					as.txQueue <- txi
				}(tx)
			}
		}
	}
	as.sequence = to
}

// Method allow trigger for resync
func (as *Aggregation) ReSync(block uint64) {
	as.isReSync = true
}

// to limit computational resources
func pauseWaitLedger(config backends.CaptiveCoreConfig, err error) error {
	if !strings.Contains(err.Error(), "is greater than max available in history archives") {
		// if not err by LatestLedger: xxx is greater than max available in history archives yyy
		return err
	}

	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindAllString(err.Error(), -1)
	seqHistoryArchives, err := strconv.Atoi(matches[1])

	if err != nil {
		return err
	}
	estimateSeqNext := int64(seqHistoryArchives) + prepareStep

	latestLedger, err := GetLatestLedger(config)
	if err != nil {
		return err
	}

	numLedgerWait := estimateSeqNext - int64(latestLedger) + 1

	if numLedgerWait < 0 {
		return nil
	}
	// Ledger closing time is ~4s/ledger
	ledgerClosingTime := 4 * time.Second
	estimateTimeWait := numLedgerWait * ledgerClosingTime.Nanoseconds()

	time.Sleep(time.Duration(estimateTimeWait))
	return nil
}
