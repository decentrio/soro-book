package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cfg "github.com/decentrio/soro-book/config"
	"github.com/decentrio/soro-book/lib/cli"
	"github.com/decentrio/soro-book/lib/log"
	"github.com/decentrio/soro-book/manager"
	"github.com/spf13/cobra"
)

var (
	DefaultCometDir = ".soro-book"
	logger          = log.NewSRLogger(log.NewSyncWriter(os.Stdout))
)

var rootCmd = &cobra.Command{
	Use: "sorobook",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

// NewRunNodeCmd returns the command that allows the CLI to start a node.
// It can be used with a custom PrivValidator and in-process ABCI application.
func NewRunNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"node", "run"},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := ParseConfig(cmd)
			if err != nil {
				return err
			}
			m := manager.DefaultNewManager(config, logger)

			if err := m.Start(); err != nil {
				return fmt.Errorf("failed to start node: %w", err)
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			go func() {
				for _ = range c {
					if m.IsRunning() {
						if err := m.Stop(); err != nil {
							fmt.Printf(err.Error())
						}
					}
					os.Exit(0)
				}
			}()

			// Run forever.
			select {}
		},
	}

	return cmd
}

// ParseConfig retrieves the default environment configuration,
// sets up the CometBFT root and ensures that the root exists
func ParseConfig(cmd *cobra.Command) (*cfg.Config, error) {
	conf := cfg.DefaultConfig()

	home, err := cmd.Flags().GetString(cli.HomeFlag)
	if err != nil {
		return nil, err
	}

	conf.RootDir = home
	conf.SetRoot(conf.RootDir)

	var managerConfig cfg.ManagerConfig
	managerConfigFile := conf.ManagerConfigFile()
	if FileExists(managerConfigFile) {
		managerConfig = LoadManagerConfig(managerConfigFile)
	} else {
		managerConfig = cfg.DefaultManagerConfig()
	}
	conf.ManagerCfg = &managerConfig

	var aggregationConfig cfg.AggregationConfig
	aggregationConfigFile := conf.AggregationConfigFile()
	if FileExists(aggregationConfigFile) {
		aggregationConfig = LoadAggregationConfig(aggregationConfigFile)
	} else {
		aggregationConfig = cfg.DefaultAggregationConfig()
	}
	conf.AggregationCfg = &aggregationConfig

	return conf, nil
}

func LoadManagerConfig(path string) cfg.ManagerConfig {
	bz, err := os.ReadFile(path)
	if err != nil {
		os.Exit(1)
	}

	var config cfg.ManagerConfig
	err = json.Unmarshal(bz, &config)
	if err != nil {
		os.Exit(1)
	}

	return config
}

func LoadAggregationConfig(path string) cfg.AggregationConfig {
	bz, err := os.ReadFile(path)
	if err != nil {
		os.Exit(1)
	}

	var config cfg.AggregationConfig
	err = json.Unmarshal(bz, &config)
	if err != nil {
		os.Exit(1)
	}

	return config
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
