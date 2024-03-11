package aggregation

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/ingest"
	backends "github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/support/log"

	aggConfig "github.com/decentrio/soro-book/aggregation/config"
	"github.com/decentrio/soro-book/aggregation/rpc"
	"github.com/decentrio/soro-book/config"
	"github.com/decentrio/soro-book/lib/service"
)

const (
	txQueueSize = 1000
	step        = 64
)

// receive data from Ledger sequence
var fromSeq = uint32(485951)

type txInfo struct {
	txHash string
}

type Aggregation struct {
	service.BaseService

	cfg config.AggregationConfig

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
	cfg config.AggregationConfig,
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

	return as
}

func (as *Aggregation) OnStart() error {
	ctx := context.Background()
	// Only log errors from the backend to keep output cleaner.
	lg := log.New()
	lg.SetLevel(logrus.ErrorLevel)
	aggConfig.Config.Log = lg

	backend, err := backends.NewCaptive(aggConfig.Config)
	panicIf(err)
	defer backend.Close()

	for {
		ledgerRange := backends.BoundedRange(fromSeq, fromSeq+step)
		err = backend.PrepareRange(ctx, ledgerRange)
		if err != nil {
			//"is greater than max available in history archives"
			err = pauseWaitLedger(err)
			if err != nil {
				return err
			}
			continue
		}

		for seq := fromSeq; seq <= fromSeq+step; seq++ {
			txReader, err := ingest.NewLedgerTransactionReader(
				ctx, backend, aggConfig.Config.NetworkPassphrase, seq,
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
			go as.process()
		}
		fromSeq += step + 1
	}
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

func panicIf(err error) {
	if err != nil {
		panic(fmt.Errorf("an error occurred, panicking: %s", err))
	}
}

func pauseWaitLedger(err error) error {
	if !strings.Contains(err.Error(), "is greater than max available in history archives") {
		// if not err by LatestLedger: is greater than max available in history archives
		return err
	}

	latestLedger, err := rpc.GetLatestLedger()
	if err != nil {
		return err
	}
	// Ledger closing time is 4s/ledger
	ledgerClosingTime := 4 * time.Second

	numLedgerWait := 64 - int64(latestLedger-fromSeq)
	timeWait := numLedgerWait * ledgerClosingTime.Nanoseconds()
	time.Sleep(time.Duration(timeWait))
	return nil
}
