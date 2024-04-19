package manager

import (
	"github.com/decentrio/soro-book/lib/log"

	"github.com/decentrio/soro-book/aggregation"
	"github.com/decentrio/soro-book/config"
)

func DefaultNewManager(cfg *config.Config, logger log.Logger) *Manager {
	as := aggregation.NewAggregation(cfg.AggregationCfg, logger)
	return NewManager(cfg.ManagerCfg, as, logger)
}
