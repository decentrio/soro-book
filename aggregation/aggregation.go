package aggregation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/decentrio/soro-book/lib/log"
	converter "github.com/decentrio/xdr-converter/converter"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/ingest"
	backends "github.com/stellar/go/ingest/ledgerbackend"
	stellar_log "github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"

	"github.com/decentrio/soro-book/config"
	db "github.com/decentrio/soro-book/database/handlers"
	"github.com/decentrio/soro-book/database/models"
	"github.com/decentrio/soro-book/lib/service"
)

const (
	txQueueSize = 1000
	prepareStep = 64
)

type LedgerWrapper struct {
	ledger models.Ledger
	txs    []TransactionWrapper
}

type Aggregation struct {
	ctx context.Context

	log *stellar_log.Entry

	config backends.CaptiveCoreConfig

	service.BaseService

	cfg *config.AggregationConfig

	backend *backends.CaptiveStellarCore

	// txQueue channel for trigger new tx
	ledgerQueue chan LedgerWrapper

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
		cfg:         cfg,
		ledgerQueue: make(chan LedgerWrapper, txQueueSize),
		isReSync:    false,
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

	as.sequence = uint32(200_000)

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
		case ledger := <-as.ledgerQueue:
			as.handleReceiveNewLedger(ledger)
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		}
	}
}

// handleReceiveTx
func (as *Aggregation) handleReceiveNewLedger(lw LedgerWrapper) {
	// Create Ledger
	// _, err := as.db.CreateLedger(&lw.ledger)
	// if err != nil {
	// 	as.Logger.Error(fmt.Sprintf("Error create ledger %d: %s", lw.ledger.Sequence, err.Error()))
	// }

	// Create Tx and Soroban events
	for _, tw := range lw.txs {
		// tx := tw.GetModelsTransaction()
		// _, err := as.db.CreateTransaction(tx)
		// if err != nil {
		// 	as.Logger.Error(fmt.Sprintf("Error create ledger %d tx %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), err.Error()))
		// }

		// Contract entry
		entries := tw.GetModelsContractDataEntry()
		for _, entry := range entries {
			var keyXdr xdr.ScVal
			var valXdr xdr.ScVal

			_ = keyXdr.UnmarshalBinary(entry.KeyXdr)
			_ = valXdr.UnmarshalBinary(entry.ValueXdr)

			key, _ := converter.ConvertScVal(keyXdr)
			keyJSON, _ := json.Marshal(key)
			val, _ := converter.ConvertScVal(valXdr)
			valJSON, _ := json.Marshal(val)

			fmt.Println("")
			fmt.Println("Ledger ", lw.ledger.Sequence)
			fmt.Printf("Tx: %s\n", tw.Tx.Result.TransactionHash.HexString())
			fmt.Printf("Key: %s\n", keyJSON)
			fmt.Printf("Val: %s\n", valJSON)
			fmt.Println("")
			// _, err := as.db.CreateContractEntry(&entry)
			// if err != nil {
			// 	as.Logger.Error(fmt.Sprintf("Error create contract data entry ledger %d tx %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), err.Error()))
			// 	continue
			// }
		}

		// Check if tx metadata is v3
		// txMetaV3, ok := tw.Tx.UnsafeMeta.GetV3()
		// if !ok {
		// 	continue
		// }

		// if txMetaV3.SorobanMeta == nil {
		// 	continue
		// }

		// // Create Event
		// for _, op := range tw.Ops {
		// 	events := op.GetContractEvents()
		// 	for _, event := range events {
		// 		_, err := as.db.CreateEvent(&event)
		// 		if err != nil {
		// 			as.Logger.Error(fmt.Sprintf("Error create event ledger %d tx %s event %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), event.ContractId, err.Error()))
		// 			continue
		// 		}
		// 	}
		// }
	}
}

func (as *Aggregation) aggregation() {
	for {
		select {
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		default:
			as.getNewLedger()
		}
	}
}

func (as *Aggregation) getNewLedger() {
	from := as.sequence
	to := as.sequence + prepareStep
	ledgerRange := backends.BoundedRange(from, to)
	err := as.backend.PrepareRange(as.ctx, ledgerRange)
	if err != nil {
		//"is greater than max available in history archives"
		time.Sleep(time.Second)
		return
	}
	for seq := from; seq < to; seq++ {
		// get ledger
		ledgerCloseMeta, err := as.backend.GetLedger(as.ctx, seq)
		if err != nil {
			as.Logger.Error(fmt.Sprintf("Error GetLedger %s", err.Error()))
			continue
		}

		ledger := getLedgerFromCloseMeta(ledgerCloseMeta)

		var txWrappers []TransactionWrapper
		var transactions = uint32(0)
		var operations = uint32(0)
		// get tx
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
				as.Logger.Error(fmt.Sprintf("Error txReader %s", err.Error()))
			}

			txWrapper := NewTransactionWrapper(tx, seq)
			txWrappers = append(txWrappers, txWrapper)

			operations += uint32(len(tx.Envelope.Operations()))
			transactions++
		}

		ledger.Transactions = transactions
		ledger.Operations = operations

		lw := LedgerWrapper{
			ledger: ledger,
			txs:    txWrappers,
		}

		go func(lwi LedgerWrapper) {
			as.ledgerQueue <- lwi
		}(lw)
	}
	as.sequence = to
}

// Method allow trigger for resync
func (as *Aggregation) ReSync(block uint64) {
	as.isReSync = true
}
