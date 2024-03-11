package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLatestLedger() (uint32, error) {
	requestData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      8675309,
		"method":  "getLatestLedger",
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Lỗi khi chuyển đổi dữ liệu JSON:", err)
		return 0, err
	}

	response, err := http.Post("https://soroban-testnet.stellar.org", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Lỗi khi gửi yêu cầu POST:", err)
		return 0, err
	}
	defer response.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		fmt.Println("Lỗi khi đọc phản hồi JSON:", err)
		return 0, err
	}

	result := responseData["result"].(map[string]interface{})
	latestLedger := result["sequence"].(float64)

	return uint32(latestLedger), nil
}
