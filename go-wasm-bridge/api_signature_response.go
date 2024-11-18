package wasmbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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

	// Process the data as needed
	fmt.Println("Response Body in signature response :", string(data2))
	//json encode string
	defer resp.Body.Close()

	return string(data2), nil
}
