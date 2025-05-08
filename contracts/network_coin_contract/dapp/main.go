package main

import (
	"fmt"
	"log"

	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/host/rbt"
)

const NETWORK_COIN_CONTRACT_WASM = "../../../../artifacts/network_coin_contract.wasm"

func executeAndGetContractResult(wasmModule *wasmbridge.WasmModule, contractInput string) (string, error) {
	// Call the function
	contractResult, err := wasmModule.CallFunction(contractInput)
	if err != nil {
		return "", fmt.Errorf("function call failed: %v", err)
	}

	return contractResult, nil
}

// Test function for creating a network coin
func testCreateNetworkCoin(wasmModule *wasmbridge.WasmModule) {
	fmt.Println("\n=== Testing Create Network Coin ===")
	contractInput := `{"create_network_coin": {"did": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu", "token_name": "RubixCoin", "symbol": "RBX", "total_supply": 1000000, "rbt_to_lock": 500}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("Unable to execute `create_network_coin` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("Create Network Coin Result: ", result)
}

// Test function for getting token balance
func testGetTokenBalance(wasmModule *wasmbridge.WasmModule, tokenAddress string) {
	fmt.Println("\n=== Testing Get Token Balance ===")
	contractInput := fmt.Sprintf(`{"get_token_balance": {"token_address": "%s", "owner_did": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu"}}`, tokenAddress)

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("Unable to execute `get_token_balance` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("Token Balance Result: ", result)
}

// Test function for transferring network coins
func testTransferNetworkCoin(wasmModule *wasmbridge.WasmModule, tokenAddress string) {
	fmt.Println("\n=== Testing Transfer Network Coin ===")
	contractInput := fmt.Sprintf(`{"transfer_network_coin": {"token_address": "%s", "sender": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu", "receiver": "bafybmienjpoihwu2y6grilbvbrrqhleoifb3irz3gu2savjmjivzqw7424", "amount": 100, "comment": "Test transfer"}}`, tokenAddress)

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("Unable to execute `transfer_network_coin` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("Transfer Network Coin Result: ", result)
}

// Test function for getting transaction history
func testGetTransactionHistory(wasmModule *wasmbridge.WasmModule, tokenAddress string) {
	fmt.Println("\n=== Testing Get Transaction History ===")
	contractInput := fmt.Sprintf(`{"get_transaction_history": {"token_address": "%s", "limit": 10, "offset": 0}}`, tokenAddress)

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("Unable to execute `get_transaction_history` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("Transaction History Result: ", result)
}

// Test function for getting all network coins
func testGetAllNetworkCoins(wasmModule *wasmbridge.WasmModule) {
	fmt.Println("\n=== Testing Get All Network Coins ===")
	contractInput := `{"get_all_network_coins": {}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("Unable to execute `get_all_network_coins` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("All Network Coins Result: ", result)
}

func main() {
	// Rubix Node Configs
	nodeAddress := "http://localhost:20006"
	quorumType := 2

	// Create Import function registry
	hostFnRegistry := wasmbridge.NewHostFunctionRegistry()

	// Explicitly register the RBT locking function
	hostFnRegistry.Register(rbt.NewDoLockRBTApiCall())

	// Initialize the WASM module
	wasmModule, err := wasmbridge.NewWasmModule(
		NETWORK_COIN_CONTRACT_WASM,
		hostFnRegistry,
		wasmbridge.WithRubixNodeAddress(nodeAddress),
		wasmbridge.WithQuorumType(quorumType),
	)
	if err != nil {
		log.Fatalf("Failed to initialize WASM module: %v", err)
	}

	// Execute test functions
	testCreateNetworkCoin(wasmModule)

	// For these tests, you'll need to update the token address with the one received from creating a network coin
	// You can either hardcode it for testing or parse it from the create_network_coin response
	tokenAddress := "nct_example_address" // Replace with actual token address from creation

	testGetTokenBalance(wasmModule, tokenAddress)
	testTransferNetworkCoin(wasmModule, tokenAddress)
	testGetTransactionHistory(wasmModule, tokenAddress)
	testGetAllNetworkCoins(wasmModule)

	fmt.Println("\nAll tests completed")
}
