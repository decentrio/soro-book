package aggregation

import (
	"context"
	"fmt"
	"time"

	backends "github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"

	"github.com/decentrio/soro-book/config"
	db "github.com/decentrio/soro-book/database/handlers"
	"github.com/decentrio/soro-book/database/models"
	"github.com/decentrio/soro-book/lib/service"
)

const (
	QueueSize          = 1000
	DefaultPrepareStep = 64
)

type Aggregation struct {
	service.BaseService

	ACfg *config.AggregationConfig

	ctx     context.Context
	Cfg     backends.CaptiveCoreConfig
	backend backends.LedgerBackend

	// txQueue channel for trigger new tx
	ledgerQueue              chan xdr.LedgerCloseMeta
	txQueue                  chan TransactionWrapper
	assetContractEventsQueue chan models.StellarAssetContractEvent
	wasmContractEventsQueue  chan models.WasmContractEvent
	contractDataEntrysQueue  chan models.ContractsData

	// isSync is flag represent if services is
	// re-synchronize
	isSync      bool
	prepareStep uint32

	StartLedgerSeq uint32
	CurrLedgerSeq  uint32

	db *db.DBHandler
}

// AggregationOption sets an optional parameter on the State.
type AggregationOption func(*Aggregation)

func NewAggregation(
	cfg *config.AggregationConfig,
	options ...AggregationOption,
) *Aggregation {
	as := &Aggregation{
		ledgerQueue:              make(chan xdr.LedgerCloseMeta, QueueSize),
		txQueue:                  make(chan TransactionWrapper, QueueSize),
		assetContractEventsQueue: make(chan models.StellarAssetContractEvent, QueueSize),
		wasmContractEventsQueue:  make(chan models.WasmContractEvent, QueueSize),
		contractDataEntrysQueue:  make(chan models.ContractsData, QueueSize),
		prepareStep:              DefaultPrepareStep,
		isSync:                   false,
		ACfg:                     cfg,
	}

	as.BaseService = *service.NewBaseService("Aggregation", as)
	for _, opt := range options {
		opt(as)
	}

	logger := log.New().WithField("module", "aggregation")
	logger.SetLevel(log.DebugLevel)
	as.BaseService.SetLogger(logger)

	as.StartLedgerSeq = as.ACfg.StartLedgerHeight
	as.CurrLedgerSeq = as.ACfg.CurrLedgerHeight

	fmt.Println(as.StartLedgerSeq, as.CurrLedgerSeq)

	as.db = db.NewDBHandler()

	as.ctx = context.Background()
	as.backend, as.Cfg = newLedgerBackend(as.ctx, *as.ACfg, as.Logger)
	return as
}

func (as *Aggregation) OnStart() error {
	as.Logger.Info("Start")
	go as.ledgerProcessing()
	go as.transactionProcessing()
	go as.contractDataEntryProcessing()
	go as.contractEventsProcessing()
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

func (as *Aggregation) aggregation() {
	for {
		select {
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		default:
			as.getNewLedger()
		}
		time.Sleep(time.Millisecond)
	}
}
