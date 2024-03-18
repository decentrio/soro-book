package manager

import (
	"github.com/decentrio/soro-book/aggregation"
	"github.com/decentrio/soro-book/config"
)

func DefaultNewManager(cfg *config.ManagerConfig) *Manager {
	asConfig := &config.AggregationConfig{}
	as := aggregation.NewAggregation(asConfig)
	return NewManager(cfg, as)
}
