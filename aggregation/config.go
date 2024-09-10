package aggregation

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/decentrio/soro-book/config"
	"github.com/stellar/go/historyarchive"
	"github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/log"
)

var (
	//go:embed configs/captive-core-pubnet.cfg
	PubnetDefaultConfig []byte

	//go:embed configs/captive-core-testnet.cfg
	TestnetDefaultConfig []byte
)

const (
	Pubnet  = "pubnet"
	Testnet = "testnet"
)

var (
	// PublicNetworkhistoryArchiveURLs is a list of history archive URLs for stellar 'pubnet'
	PublicNetworkhistoryArchiveURLs = []string{
		"https://history.stellar.org/prd/core-live/core_live_001/",
		"https://history.stellar.org/prd/core-live/core_live_002/",
		"https://history.stellar.org/prd/core-live/core_live_003/",
		"https://stellar-history-sg-sin.satoshipay.io/",
		"https://stellar-history-us-iowa.satoshipay.io/",
		"https://archive.v4.stellar.lobstr.co/",
		"https://hongkong.stellar.whalestack.com/history/",
		"https://hercules-history.publicnode.org/",
		"https://lyra-history.publicnode.org/",
		"https://stellar-full-history3.bdnodes.net/",
		"https://archive.v5.stellar.lobstr.co/",
		"https://stellar-history-ins.franklintempleton.com/azinsshf401/",
	}

	// TestNetworkhistoryArchiveURLs is a list of history archive URLs for stellar 'testnet'
	TestNetworkhistoryArchiveURLs = []string{
		"https://history.stellar.org/prd/core-testnet/core_testnet_001/",
		"https://history.stellar.org/prd/core-testnet/core_testnet_002/",
		"https://history.stellar.org/prd/core-testnet/core_testnet_003",
	}
)

const (
	// PublicNetworkPassphrase is the pass phrase used for every transaction intended for the public stellar network
	PublicNetworkPassphrase = "Public Global Stellar Network ; September 2015"
	// TestNetworkPassphrase is the pass phrase used for every transaction intended for the SDF-run test network
	TestNetworkPassphrase = "Test SDF Network ; September 2015"
)

func newLedgerBackend(ctx context.Context, config config.AggregationConfig, log *log.Entry) (ledgerbackend.LedgerBackend, ledgerbackend.CaptiveCoreConfig) {
	// generate CaptiveCoreConfig
	var (
		networkPassphrase  string
		historyArchiveURLs []string
		captiveCoreConfig  []byte
	)
	// Default network config
	switch config.Network {
	case Pubnet:
		networkPassphrase = network.PublicNetworkPassphrase
		historyArchiveURLs = network.PublicNetworkhistoryArchiveURLs
		captiveCoreConfig = PubnetDefaultConfig

	case Testnet:
		networkPassphrase = network.TestNetworkPassphrase
		historyArchiveURLs = network.TestNetworkhistoryArchiveURLs
		captiveCoreConfig = TestnetDefaultConfig

	default:
		log.Fatalf("Invalid network %s", config.Network)
	}

	params := ledgerbackend.CaptiveCoreTomlParams{
		NetworkPassphrase:  networkPassphrase,
		HistoryArchiveURLs: historyArchiveURLs,
		UseDB:              true,
	}
	captiveCoreToml, err := ledgerbackend.NewCaptiveCoreTomlFromData(captiveCoreConfig, params)
	if err != nil {
		log.WithError(err)
	}

	captiveConfig := ledgerbackend.CaptiveCoreConfig{
		BinaryPath:          config.BinaryPath,
		NetworkPassphrase:   params.NetworkPassphrase,
		HistoryArchiveURLs:  params.HistoryArchiveURLs,
		CheckpointFrequency: historyarchive.DefaultCheckpointFrequency,
		Log:                 log.WithField("subservice", "stellar-core"),
		Toml:                captiveCoreToml,
		UserAgent:           "ledger-exporter",
		UseDB:               true,
	}
	// Create a new captive core backend
	backend, err := ledgerbackend.NewCaptive(captiveConfig)
	if err != nil {
		log.WithError(err)
	}

	// var ledgerRange ledgerbackend.Range
	// if config.EndLedgerHeight == 0 {
	// 	ledgerRange = ledgerbackend.UnboundedRange(config.StartLedgerHeight)
	// 	log.Info("running in online mode")
	// } else {
	// 	ledgerRange = ledgerbackend.BoundedRange(config.StartLedgerHeight, config.EndLedgerHeight)
	// 	log.Info("running in offline mode")
	// }

	// err = backend.PrepareRange(ctx, ledgerRange)

	return backend, captiveConfig
}

func panicIf(err error) {
	if err != nil {
		panic(fmt.Errorf("an error occurred, panicking: %s", err))
	}
}
