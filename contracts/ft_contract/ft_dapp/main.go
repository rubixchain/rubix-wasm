package main

import (
	"fmt"
	"log"

	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)

const FT_CONTRACT_WASM = "../../../artifacts/ft_contract.wasm"

func executeAndGetContractResult(wasmModule *wasmbridge.WasmModule, contractInput string) (string, error) {
	// Call the function
	contractResult, err := wasmModule.CallFunction(contractInput)
	if err != nil {
		return "", fmt.Errorf("function call failed: %v", err)
	}

	return contractResult, nil
}

func mintFTFunc1(wasmModule *wasmbridge.WasmModule) {
	contractInput := `{"mint_sample_ft":{"name": "rubix1", "ft_info": {
  "did": "bafybmihho2yn53uxmso4yl6rv2dik35l4gmalt3gypl4ohpxuw7iy4a3ke",
  "ft_count": 10,
  "ft_name": "test17",
  "ft_num_start_index": 0,
  "token_count": 1
}}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("unable to execute `mint_sample_ft` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("mint_sample_ft Result: ", result)
}
func mintFTFunc2(wasmModule *wasmbridge.WasmModule) {
	contractInput := `{"mint_sample_ft":{"name": "rubix1", "ft_info": {
  "did": "bafybmier5otnwn7m7yodowfvzfbobwha3cagyrjphrtifa45jug3tafzxi",
  "ft_count": 10,
  "ft_name": "test17",
  "ft_num_start_index": 10,
  "token_count": 1
}}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("unable to execute `mint_sample_ft` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("mint_sample_ft Result: ", result)
}

func transferFTFunc(wasmModule *wasmbridge.WasmModule) {
	contractInput := `{"transfer_sample_ft":{"name": "rubix1", "ft_info": {"comment":"testing ft transfer","ft_count":1,"ft_name":"apple","sender": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu","creatorDID": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu", "receiver": "bafybmienjpoihwu2y6grilbvbrrqhleoifb3irz3gu2savjmjivzqw7424"}}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("unable to execute `transfer_sample_ft` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("transfer_sample_ft Result: ", result)
}

func main() {
	// Rubix Node Configs
	nodeAddress := "http://localhost:20006"
	quorumType := 2

	// Create Import function registry
	hostFnRegistry := wasmbridge.NewHostFunctionRegistry()

	// Initialize the WASM module
	wasmModule, err := wasmbridge.NewWasmModule(
		FT_CONTRACT_WASM,
		hostFnRegistry,
		wasmbridge.WithRubixNodeAddress(nodeAddress),
		wasmbridge.WithQuorumType(quorumType),
	)
	if err != nil {
		log.Fatalf("Failed to initialize WASM module: %v", err)
	}
	mintFTFunc1(wasmModule)
	fmt.Println("***First 100 created, trying to create next 100***")
	mintFTFunc2(wasmModule)
	//transferFTFunc(wasmModule)
}
