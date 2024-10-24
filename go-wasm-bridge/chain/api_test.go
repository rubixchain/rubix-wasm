package chain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock Handler for GetSmartContractTokenChain API
func mockSmartContractChainHandler(w http.ResponseWriter, r *http.Request) {
	registeredSmartContracts := []string{
		"Qm123",
		"Qm456",
	}

	// Read the request body to check for the `latest` flag
	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	latest, ok := requestData["latest"].(bool)
	if !ok {
		latest = false // Default to false if "latest" is not present or is invalid
	}
	token, ok := requestData["token"].(string)
	if !ok {
		panic("invalid `token` value")
	}

	var response smartContractTokenChainDataApiResponse

	// Read API request body params
	if !contains(registeredSmartContracts, token) {
		response = smartContractTokenChainDataApiResponse{
			BasicResponse: BasicResponse{
				Status:  false,
				Message: "Smart Contract not registered",
			},
			SCTDataReply: []*SmartContractBlockInfo{},
		}
	} else {
		response = smartContractTokenChainDataApiResponse{
			BasicResponse: BasicResponse{
				Status:  true,
				Message: "Success",
			},
			SCTDataReply: []*SmartContractBlockInfo{
				{
					BlockNo:           1,
					BlockId:           "block-id-1",
					SmartContractData: "data-1",
				},
				{
					BlockNo:           2,
					BlockId:           "block-id-2",
					SmartContractData: "data-2",
				},
				{
					BlockNo:           3,
					BlockId:           "block-id-3",
					SmartContractData: "data-3",
				},
			},
		}
	}

	if latest {
		response.SCTDataReply = response.SCTDataReply[len(response.SCTDataReply)-1:]
	}

	// Marshal the response to JSON and write it to the ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func TestGetSmartContractTokenChain_RegisteredContract_CompleteChain(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(mockSmartContractChainHandler))
	defer mockServer.Close()

	contractHash := "Qm123"
	latestBlockOnly := false
	result, err := GetSmartContractTokenChain(mockServer.URL, contractHash, latestBlockOnly)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the response contains the correct number of blocks
	expectedNumBlocks := 3
	if len(result) != expectedNumBlocks {
		t.Errorf("Expected %d blocks, got %d", expectedNumBlocks, len(result))
	}

	// Check if the content of the first block matches the mock data
	expectedBlockID := "block-id-1"
	if result[0].BlockId != expectedBlockID {
		t.Errorf("Expected BlockId %s, got %s", expectedBlockID, result[0].BlockId)
	}

	// Check if the SmartContractData matches the mock data
	expectedSmartContractData := "data-1"
	if result[0].SmartContractData != expectedSmartContractData {
		t.Errorf("Expected SmartContractData %s, got %s", expectedSmartContractData, result[0].SmartContractData)
	}
}

func TestGetSmartContractTokenChain_RegisteredContract_LatestBlock(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(mockSmartContractChainHandler))
	defer mockServer.Close() // Ensure the server is closed when the test finishes

	// Call the GetSmartContractTokenChain function with the mock server URL
	contractHash := "Qm123"
	latestBlockOnly := true
	result, err := GetSmartContractTokenChain(mockServer.URL, contractHash, latestBlockOnly)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the response contains the correct number of blocks
	expectedNumBlocks := 1
	if len(result) != expectedNumBlocks {
		t.Errorf("Expected %d blocks, got %d", expectedNumBlocks, len(result))
	}

	// Check if the content of the first block matches the mock data
	expectedBlockID := "block-id-3"
	if result[0].BlockId != expectedBlockID {
		t.Errorf("Expected BlockId %s, got %s", expectedBlockID, result[0].BlockId)
	}

	// Check if the SmartContractData matches the mock data
	expectedSmartContractData := "data-3"
	if result[0].SmartContractData != expectedSmartContractData {
		t.Errorf("Expected SmartContractData %s, got %s", expectedSmartContractData, result[0].SmartContractData)
	}
}

func TestGetSmartContractTokenChain_NonRegisteredContract(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(mockSmartContractChainHandler))
	defer mockServer.Close()

	contractHash := "Qm125"
	latestBlockOnly := false
	result, err := GetSmartContractTokenChain(mockServer.URL, contractHash, latestBlockOnly)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the response contains the correct number of blocks
	expectedNumBlocks := 0
	if len(result) != expectedNumBlocks {
		t.Errorf("Expected %d blocks, got %d", expectedNumBlocks, len(result))
	}
}