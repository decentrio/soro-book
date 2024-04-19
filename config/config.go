package config

import (
	"path/filepath"

	"github.com/stellar/go/ingest/ledgerbackend"
)

const (
	DefaultConfigDir = "config"

	DefaultManagerConfigFileName     = "managerConfig.json"
	DefaultAggregationConfigFileName = "aggregationConfig.json"
)

type Config struct {
	RootDir        string
	ManagerCfg     *ManagerConfig
	AggregationCfg *AggregationConfig
}

func DefaultConfig() *Config {
	return &Config{}
}

func (c *Config) SetRoot(root string) {
	c.RootDir = root
}

func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func (c *Config) ManagerConfigFile() string {
	return rootify(DefaultManagerConfigFileName, c.RootDir)
}

func (c *Config) AggregationConfigFile() string {
	return rootify(DefaultAggregationConfigFileName, c.RootDir)
}

type ManagerConfig struct {
}

func DefaultManagerConfig() ManagerConfig {
	return ManagerConfig{}
}

type AggregationConfig struct {
	archiveURL        string
	networkPassphrase string
	toml              *ledgerbackend.CaptiveCoreToml
}

func DefaultAggregationConfig() AggregationConfig {
	return AggregationConfig{}
}
