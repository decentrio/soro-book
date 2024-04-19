package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultConfigDir = "config"

	DefaultManagerConfigFileName     = "managerConfig.json"
	DefaultAggregationConfigFileName = "aggregationConfig.json"
)

type ManagerConfig struct {
	RootDir        string
	AggregationCfg *AggregationConfig
}

func DefaultConfig() *ManagerConfig {
	return &ManagerConfig{}
}

func (c *ManagerConfig) SetRoot(root string) {
	c.RootDir = root
}

func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func (c *ManagerConfig) ManagerConfigFile() string {
	return rootify(DefaultManagerConfigFileName, c.RootDir)
}

func (c *ManagerConfig) AggregationConfigFile() string {
	return rootify(DefaultAggregationConfigFileName, c.RootDir)
}

func (c *ManagerConfig) LoadManagerConfig(path string) {
	bz, err := os.ReadFile(path)
	if err != nil {
		os.Exit(1)
	}

	err = json.Unmarshal(bz, c)
	if err != nil {
		os.Exit(1)
	}
}

type AggregationConfig struct {
	ArchiveURL        string `json:"url,omitempty"`
	NetworkPassphrase string `json:"network_passphrase,omitempty"`
	BinaryPath        string `json:"binary_path,omitempty"`
	LedgerHeight      uint32 `json:"ledger_height,omitempty"`
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
		LedgerHeight:      1193200,
	}
}

func LoadAggregationConfig(path string) AggregationConfig {
	bz, err := os.ReadFile(path)
	if err != nil {
		os.Exit(1)
	}

	var config AggregationConfig
	err = json.Unmarshal(bz, &config)
	if err != nil {
		os.Exit(1)
	}

	return config
}

func WriteState(path string, content []byte, mode os.FileMode) error {
	if !FileExists(path) {
		if err := os.MkdirAll(filepath.Dir(path), mode); err != nil {
			return err
		}
		os.Create(path)
	}

	return os.WriteFile(path, content, mode)
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
