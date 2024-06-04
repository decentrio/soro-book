package manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/decentrio/soro-book/aggregation"
	"github.com/decentrio/soro-book/config"
	"github.com/decentrio/soro-book/lib/service"
	"github.com/stellar/go/support/log"
)

// Manager is the root service that manage all services
type Manager struct {
	service.BaseService

	// config of Manager
	cfg *config.ManagerConfig

	// aggregation services
	as *aggregation.Aggregation
}

const (
	PaddingLedger = 2560
)

// StateOption sets an optional parameter on the State.
type ManagerOption func(*Manager)

// NewBaseService creates a new manager.
func NewManager(
	cfg *config.ManagerConfig,
	as *aggregation.Aggregation,
	options ...ManagerOption,
) *Manager {
	m := &Manager{
		cfg: cfg,
		as:  as,
	}

	m.BaseService = *service.NewBaseService("Manager", m)
	for _, opt := range options {
		opt(m)
	}

	m.BaseService.SetLogger(log.New().WithField("module", "manager"))

	return m
}

func (m *Manager) OnStart() error {
	m.Logger.Info("Start")
	m.as.Start()
	return nil
}

func (m *Manager) OnStop() error {
	m.Logger.Info("Stop")
	m.as.Stop()

	asConfig := *m.as.ACfg
	asConfig.StartLedgerHeight = m.as.StartLedgerSeq - PaddingLedger

	bz, err := json.Marshal(asConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(m.cfg.AggregationConfigFile())
	err = config.WriteState(m.cfg.AggregationConfigFile(), bz, 0777)
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(time.Second)
	return nil
}
