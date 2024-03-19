package manager

import (
	"time"

	"github.com/decentrio/soro-book/aggregation"
	"github.com/decentrio/soro-book/config"
	"github.com/decentrio/soro-book/lib/log"
	"github.com/decentrio/soro-book/lib/service"
)

// Manager is the root service that manage all services
type Manager struct {
	service.BaseService

	// config of Manager
	cfg *config.ManagerConfig

	// aggregation services
	as *aggregation.Aggregation
}

// StateOption sets an optional parameter on the State.
type ManagerOption func(*Manager)

// NewBaseService creates a new manager.
func NewManager(
	cfg *config.ManagerConfig,
	as *aggregation.Aggregation,
	log log.Logger,
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

	m.BaseService.SetLogger(log.With("module", "manager"))

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
	time.Sleep(time.Second)
	return nil
}
