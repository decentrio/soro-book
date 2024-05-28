package config

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	Network           string `json:"network,omitempty"`
	BinaryPath        string `json:"binary_path,omitempty"`
	StartLedgerHeight uint32 `json:"start_ledger_height,omitempty"`
	EndLedgerHeight   uint32 `json:"end_ledger_height,omitempty"`
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
