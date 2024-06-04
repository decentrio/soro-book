package manager

import (
	"github.com/decentrio/soro-book/aggregation"
	"github.com/decentrio/soro-book/config"
)

func DefaultNewManager(cfg *config.ManagerConfig) *Manager {
	as := aggregation.NewAggregation(cfg.AggregationCfg)
	return NewManager(cfg, as)
}
