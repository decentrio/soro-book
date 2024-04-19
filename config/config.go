package config

import (
	"fmt"
	"path/filepath"
	"runtime"

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
	ArchiveURL        string
	NetworkPassphrase string
	BinaryPath        string
	Core              *ledgerbackend.CaptiveCoreToml
	LedgerHeight      uint64
}

func DefaultAggregationConfig() AggregationConfig {
	var binaryPath string
	os := runtime.GOOS
	switch os {
	case "darwin":
		binaryPath = "../bin/stellar-core-mac"
	case "linux":
		binaryPath = "../bin/stellar-core-linux"
	default:
		fmt.Printf("%s.\n", os)
	}

	return AggregationConfig{
		ArchiveURL:        "https://history.stellar.org/prd/core-testnet/core_testnet_002",
		NetworkPassphrase: "Test SDF Network ; September 2015",
		BinaryPath:        binaryPath,
		LedgerHeight:      2,
	}
}
