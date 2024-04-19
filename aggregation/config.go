package aggregation

import (
	"fmt"

	"github.com/stellar/go/ingest/ledgerbackend"
)

func CaptiveCoreConfig(archiveURLs []string, networkPassphrase string, binaryPath string, core *ledgerbackend.CaptiveCoreToml) ledgerbackend.CaptiveCoreConfig {
	captiveCoreToml, err := ledgerbackend.NewCaptiveCoreToml(ledgerbackend.CaptiveCoreTomlParams{
		NetworkPassphrase:  networkPassphrase,
		HistoryArchiveURLs: archiveURLs,
	})
	panicIf(err)

	if core == nil {
		core, err = captiveCoreToml.CatchupToml()
		panicIf(err)
	}

	return ledgerbackend.CaptiveCoreConfig{
		// Change these based on your environment:
		BinaryPath:         binaryPath,
		NetworkPassphrase:  networkPassphrase,
		HistoryArchiveURLs: archiveURLs,
		Toml:               core,
	}
}

func panicIf(err error) {
	if err != nil {
		panic(fmt.Errorf("an error occurred, panicking: %s", err))
	}
}
