package aggregation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	backends "github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/xdr"
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

// parseBalanceChangeEvent is a generalization of a subset of the Stellar Asset
// Contract events. Transfer, mint, clawback, and burn events all have two
// addresses and an amount involved. The addresses represent different things in
// different event types (e.g. "from" or "admin"), but the parsing is identical.
// This helper extracts all three parts or returns a generic error if it can't.
func parseBalanceChangeEvent(topics xdr.ScVec, value xdr.ScVal) (
	first string,
	second string,
	amount Int128Parts,
	err error,
) {
	err = ErrNotBalanceChangeEvent
	if len(topics) != 4 {
		return
	}

	firstSc, ok := topics[1].GetAddress()
	if !ok {
		return
	}
	first, err = firstSc.String()
	if err != nil {
		err = errors.Wrap(err, ErrNotBalanceChangeEvent.Error())
		return
	}

	secondSc, ok := topics[2].GetAddress()
	if !ok {
		return
	}
	second, err = secondSc.String()
	if err != nil {
		err = errors.Wrap(err, ErrNotBalanceChangeEvent.Error())
		return
	}

	xdrAmount, ok := value.GetI128()
	if !ok {
		return
	}

	amount = XdrInt128PartsConvert(xdrAmount)

	return first, second, amount, nil
}

func XdrInt128PartsConvert(in xdr.Int128Parts) Int128Parts {
	out := Int128Parts{
		Hi: int64(in.Hi),
		Lo: uint64(in.Lo),
	}

	return out
}
