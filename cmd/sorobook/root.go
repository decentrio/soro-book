package main

import (
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
func ParseConfig(cmd *cobra.Command) (*cfg.ManagerConfig, error) {
	conf := cfg.DefaultConfig()

	home, err := cmd.Flags().GetString(cli.HomeFlag)
	if err != nil {
		return nil, err
	}

	conf.RootDir = home
	conf.SetRoot(conf.RootDir)

	managerConfigFile := conf.ManagerConfigFile()
	if cfg.FileExists(managerConfigFile) {
		conf.LoadManagerConfig(managerConfigFile)
	}

	var aggregationConfig cfg.AggregationConfig
	aggregationConfigFile := conf.AggregationConfigFile()
	if cfg.FileExists(aggregationConfigFile) {
		aggregationConfig = cfg.LoadAggregationConfig(aggregationConfigFile)
	} else {
		aggregationConfig = cfg.DefaultAggregationConfig()
	}
	conf.AggregationCfg = &aggregationConfig

	return conf, nil
}
