package aggregation

import (
	"context"
	"time"

	"github.com/decentrio/soro-book/lib/log"
	"github.com/sirupsen/logrus"
	backends "github.com/stellar/go/ingest/ledgerbackend"
	stellar_log "github.com/stellar/go/support/log"

	"github.com/decentrio/soro-book/config"
	db "github.com/decentrio/soro-book/database/handlers"
	"github.com/decentrio/soro-book/database/models"
	"github.com/decentrio/soro-book/lib/service"
)

const (
	QueueSize   = 1000
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
	ledgerQueue              chan LedgerWrapper
	txQueue                  chan TransactionWrapper
	assetContractEventsQueue chan models.StellarAssetContractEvent
	wasmContractEventsQueue  chan models.WasmContractEvent
	contractDataEntrysQueue  chan models.Contract

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
		cfg:                      cfg,
		ledgerQueue:              make(chan LedgerWrapper, QueueSize),
		txQueue:                  make(chan TransactionWrapper, QueueSize),
		assetContractEventsQueue: make(chan models.StellarAssetContractEvent, QueueSize),
		wasmContractEventsQueue:  make(chan models.WasmContractEvent, QueueSize),
		contractDataEntrysQueue:  make(chan models.Contract, QueueSize),
		isReSync:                 false,
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

	as.sequence = uint32(10)

	var err error
	as.backend, err = backends.NewCaptive(Config)
	panicIf(err)

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

// Method allow trigger for resync
func (as *Aggregation) ReSync(block uint64) {
	as.isReSync = true
}
