package manager

import (
	"github.com/decentrio/soro-book/lib/log"

	"github.com/decentrio/soro-book/aggregation"
	"github.com/decentrio/soro-book/config"
)

func DefaultNewManager(cfg *config.ManagerConfig, logger log.Logger) *Manager {
	asConfig := &config.AggregationConfig{}
	as := aggregation.NewAggregation(asConfig, logger)
	return NewManager(cfg, as, logger)
}
