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

func mintFTFunc(wasmModule *wasmbridge.WasmModule) {
	contractInput := `{"mint_sample_ft":{"name": "rubix1", "ft_info": {
  "did": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu",
  "ftcount": 100,
  "ftname": "apple",
  "tokencount": 1,
  "port":"20006"
}}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("unable to execute `mint_sample_ft` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("mint_sample_ft Result: ", result)
}

func transferFTFunc(wasmModule *wasmbridge.WasmModule) {
	contractInput := `{"transfer_sample_ft":{"name": "rubix1", "ft_info": {"comment":"testing ft transfer","FTCount":1,"FTName":"apple","sender": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu","creatorDID": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu", "receiver": "bafybmienjpoihwu2y6grilbvbrrqhleoifb3irz3gu2savjmjivzqw7424","port":"20006","type":2}}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("unable to execute `transfer_sample_ft` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("transfer_sample_ft Result: ", result)
}

func main() {
	// Create Import function registry
	hostFnRegistry := wasmbridge.NewHostFunctionRegistry()
	fmt.Println(hostFnRegistry)
	// Initialize the WASM module
	wasmModule, err := wasmbridge.NewWasmModule(FT_CONTRACT_WASM, hostFnRegistry)
	if err != nil {
		log.Fatalf("Failed to initialize WASM module: %v", err)
	}

	//mintFTFunc(wasmModule)
	transferFTFunc(wasmModule)
}