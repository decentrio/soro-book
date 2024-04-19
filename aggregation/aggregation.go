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
	QueueSize          = 10000
	DefaultPrepareStep = 64
)

type State int32

const (
	LEDGER State = iota
	TX
	CONTRACT
)

type Aggregation struct {
	service.BaseService

	ctx context.Context
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
	isReSync    bool
	prepareStep uint32

	state          State
	startLedgerSeq uint32
	currLedgerSeq  uint32

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
		state:                    LEDGER,
		prepareStep:              DefaultPrepareStep,
		isReSync:                 false,
	}

	as.BaseService = *service.NewBaseService("Aggregation", as)
	for _, opt := range options {
		opt(as)
	}

	as.BaseService.SetLogger(logger.With("module", "aggregation"))

	as.db = db.NewDBHandler()

	as.ctx = context.Background()

	Config := CaptiveCoreConfig([]string{as.cfg.ArchiveURL}, as.cfg.NetworkPassphrase, as.cfg.BinaryPath, nil)
	log := stellar_log.New()
	log.SetLevel(logrus.ErrorLevel)
	Config.Log = log

	as.startLedgerSeq = uint32(as.cfg.LedgerHeight)

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
		time.Sleep(time.Second)
	}
}

// Method allow trigger for resync
func (as *Aggregation) ReSync(block uint64) {
	as.isReSync = true
}
