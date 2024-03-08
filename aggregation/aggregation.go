package aggregation

import (
	"github.com/decentrio/soro-book/config"
	"github.com/decentrio/soro-book/lib/service"
)

const (
	txQueueSize = 1000
)

type txInfo struct {
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
