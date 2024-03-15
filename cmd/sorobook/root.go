package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/decentrio/soro-book/config"
	"github.com/decentrio/soro-book/manager"
	"github.com/spf13/cobra"
)

var DefaultCometDir = ".soro-book"

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
			// TODO: We need to read config
			cfg := &config.ManagerConfig{}
			m := manager.DefaultNewManager(cfg)

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
