package aggregation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	backends "github.com/stellar/go/ingest/ledgerbackend"
)

func GetLatestLedger(config backends.CaptiveCoreConfig) (uint32, error) {
	requestData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      8675309,
		"method":  "getLatestLedger",
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error converting JSON data:", err)
		return 0, err
	}

	if config.NetworkPassphrase != "Test SDF Network ; September 2015" {
		// todo: mainnet
		return rpcGetLatestLedger("", requestBody)

	}

	url := "https://soroban-testnet.stellar.org"
	return rpcGetLatestLedger(url, requestBody)
}

func rpcGetLatestLedger(url string, requestBody []byte) (uint32, error) {
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return 0, err
	}
	defer response.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		fmt.Println("Error reading JSON response:", err)
		return 0, err
	}

	result := responseData["result"].(map[string]interface{})
	latestLedger := result["sequence"].(float64)

	return uint32(latestLedger), nil
}

var ErrNotBalanceChangeEvent = errors.New("event doesn't represent a balance change")
