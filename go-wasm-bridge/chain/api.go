package chain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BasicResponse struct {
    Status  bool        `json:"status"`
    Message string      `json:"message"`
    Result  interface{} `json:"result"`
}

type smartContractTokenChainDataApiResponse struct {
	BasicResponse
	SCTDataReply []*SmartContractBlockInfo
}

type SmartContractBlockInfo struct {
	BlockNo           uint64
	BlockId           string
	SmartContractData string
}

// GetSmartContractTokenChain fetches the Smart Contract chain for an input contractHash
//
// If `latestBlockOnly` is true, only the lasted block of the chain will be returned
func GetSmartContractTokenChain(rubixNodeURL, contractHash string, latestBlockOnly bool) ([]*SmartContractBlockInfo, error) {
	// Make a request to /api/get-smart-contract-token-chain-data to retrive the chain
	requestURL := fmt.Sprintf("%v/api/get-smart-contract-token-chain-data", rubixNodeURL)
	
	requestBody := map[string]interface{}{
		"token": contractHash,
		"latest": latestBlockOnly,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractTokenChain: unable to marshal request struct")
	}
	requestBodyReader := bytes.NewReader(requestBodyBytes)

	apiResponse, err := http.Post(requestURL, "application/json", requestBodyReader)
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractTokenChain: unable to get reponse, err: %v", err)
	}
	defer apiResponse.Body.Close()

	apiResonseBody, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractTokenChain: unable to parse response body, err: %v", err)
	}

	var apiResponseStruct *smartContractTokenChainDataApiResponse
	if err := json.Unmarshal(apiResonseBody, &apiResponseStruct); err != nil {
		return nil, fmt.Errorf("GetSmartContractTokenChain: unable to form the token chain as unmarshaling failed, err: %v", err)
	}

	return apiResponseStruct.SCTDataReply, nil
}