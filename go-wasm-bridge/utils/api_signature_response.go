package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type BasicResponse struct {
	Message string `json:"message"`
	Result  string `json:"result"`
	Status  bool   `json:"status"`
}

// SignatureResponse signs the transaction and returns the transaction hash
func SignatureResponse(requestId string, nodeAddress string) (string, error) {
	data := map[string]interface{}{
		"id":       requestId,
		"mode":     0,
		"password": "mypassword",
	}

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Error marshaling JSON:", err)
	}

	url, err := url.JoinPath(nodeAddress, "/api/signature-response")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return "", fmt.Errorf("Error creating HTTP request:", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending HTTP request:", err)
	}

	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %s\n", err)
	}

	var basicResponse *BasicResponse
	if err := json.Unmarshal(data2, &basicResponse); err != nil {
		return "", fmt.Errorf("unable to unmarshal reponse")
	}

	return extractTransactionID(basicResponse.Message), nil
}

func extractTransactionID(txMsg string) string {
	txMsgElems := strings.Split(txMsg, " ")

	if len(txMsgElems) == 0 {
		return ""
	}

	txID := txMsgElems[len(txMsgElems) - 1]
	return txID
}
