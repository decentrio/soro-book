package main

import (
	"os"
	"path/filepath"

	"github.com/decentrio/soro-book/lib/cli"
)

func main() {
	rootCmd.AddCommand(NewRunNodeCmd())
	cmd := cli.PrepareBaseCmd(rootCmd, "CMT", os.ExpandEnv(filepath.Join("$HOME", DefaultCometDir)))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
