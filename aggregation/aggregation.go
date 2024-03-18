package aggregation

import (
	"context"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/ingest"
	backends "github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/support/log"

	"github.com/decentrio/soro-book/config"
	"github.com/decentrio/soro-book/lib/service"
)

const (
	txQueueSize = 1000
	step        = 1
)

// receive data from Ledger sequence
var fromSeq = uint32(485951)

type txInfo struct {
	txHash string
}

type Aggregation struct {
	ctx context.Context

	log *log.Entry

	config backends.CaptiveCoreConfig

	service.BaseService

	cfg *config.AggregationConfig

	// txQueue channel for trigger new tx
	txQueue chan txInfo

	// isReSync is flag represent if services is
	// re-synchronize
	isReSync bool

	// subscribe services
}

// AggregationOption sets an optional parameter on the State.
type AggregationOption func(*Aggregation)

func NewAggregation(
	cfg *config.AggregationConfig,
	options ...AggregationOption,
) *Aggregation {
	as := &Aggregation{
		cfg:      cfg,
		txQueue:  make(chan txInfo, txQueueSize),
		isReSync: false,
	}

	as.BaseService = *service.NewBaseService("Aggregation", as)
	for _, opt := range options {
		opt(as)
	}

	as.ctx = context.Background()
	as.log = log.New()
	as.log.SetLevel(logrus.ErrorLevel)
	Config.Log = as.log

	return as
}

func (as *Aggregation) OnStart() error {
	// Note that when using goroutines, you need to be careful to ensure that no
	// race conditions occur when accessing the txQueue.
	as.aggregation()

	go as.process()
	return nil
}

func (as *Aggregation) OnStop() error {
	return nil
}

// aggregation process
func (as *Aggregation) process() {
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
func (as *Aggregation) handleReceiveTx(tx txInfo) {
	// filter

	// callback
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
	estimateSeqNext := int64(seqHistoryArchives) + step

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

func (as *Aggregation) aggregation() error {
	backend, err := backends.NewCaptive(Config)
	panicIf(err)
	defer backend.Close()

	for {
		ledgerRange := backends.BoundedRange(fromSeq, fromSeq+step)
		err = backend.PrepareRange(as.ctx, ledgerRange)
		if err != nil {
			//"is greater than max available in history archives"
			err = pauseWaitLedger(as.config, err)
			if err != nil {
				return err
			}
			continue
		}
		for seq := fromSeq; seq < fromSeq+step; seq++ {
			txReader, err := ingest.NewLedgerTransactionReader(
				as.ctx, backend, Config.NetworkPassphrase, seq,
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
					return err
				}

				if tx.Result.Successful() {
					newTxInfo := txInfo{
						txHash: tx.Result.TransactionHash.HexString(),
					}

					go func(tx txInfo) {
						// add txInfo chan txQueue <- tx
						as.txQueue <- tx
					}(newTxInfo)
				}
			}
		}
		fromSeq += step
		break
	}
	return nil
}
