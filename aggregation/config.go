package aggregation

import (
	"fmt"
	"runtime"

	"github.com/stellar/go/ingest/ledgerbackend"
)

var (
	Config = CaptiveCoreConfigDefauTestNet()
)

func CaptiveCoreConfigDefauTestNet() ledgerbackend.CaptiveCoreConfig {
	archiveURLs := []string{
		"https://history.stellar.org/prd/core-testnet/core_testnet_001",
		"https://history.stellar.org/prd/core-testnet/core_testnet_002",
		"https://history.stellar.org/prd/core-testnet/core_testnet_003",
	}
	networkPassphrase := "Test SDF Network ; September 2015"

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

	return CaptiveCoreConfig(archiveURLs, networkPassphrase, binaryPath)
}

func CaptiveCoreConfig(archiveURLs []string, networkPassphrase string, binaryPath string) ledgerbackend.CaptiveCoreConfig {
	captiveCoreToml, err := ledgerbackend.NewCaptiveCoreToml(ledgerbackend.CaptiveCoreTomlParams{
		NetworkPassphrase:  networkPassphrase,
		HistoryArchiveURLs: archiveURLs,
	})
	panicIf(err)

	captiveCoreToml, err = captiveCoreToml.CatchupToml()
	panicIf(err)

	return ledgerbackend.CaptiveCoreConfig{
		// Change these based on your environment:
		BinaryPath:         binaryPath,
		NetworkPassphrase:  networkPassphrase,
		HistoryArchiveURLs: archiveURLs,
		Toml:               captiveCoreToml,
	}
}

func panicIf(err error) {
	if err != nil {
		panic(fmt.Errorf("an error occurred, panicking: %s", err))
	}
}
