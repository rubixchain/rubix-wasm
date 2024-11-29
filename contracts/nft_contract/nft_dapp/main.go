package main

import (
	"fmt"
	"log"

	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)

const NFT_CONTRACT_WASM = "../../../artifacts/nft_contract.wasm"

func executeAndGetContractResult(wasmModule *wasmbridge.WasmModule, contractInput string) (string, error) {
	// Call the function
	contractResult, err := wasmModule.CallFunction(contractInput)
	if err != nil {
		return "", fmt.Errorf("function call failed: %v", err)
	}

	return contractResult, nil
}

func mintNFTFunc(wasmModule *wasmbridge.WasmModule) {
	//contractInput := `{"create_sample_nft":{"did":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","metadata":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/metadata.json","artifact":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/testimage24.png","port":"20024","quorumtype":2}}`
	contractInput := `{
		"mint_sample_nft": {
		  "name": "rubix1",
		  "nft_info": {
			"did": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu",
			"metadata": "C:\\Users\\allen\\Downloads\\metadata.json",
			"artifact": "C:\\Users\\allen\\Downloads\\5861.pdf"
		  }
		}
	  }`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("unable to execute `mint_sample_nft` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("mint_sample_nft Result: ", result)
}

func transferNFTFunc(wasmModule *wasmbridge.WasmModule) {
	//contractInput := `{"create_sample_nft":{"did":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","metadata":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/metadata.json","artifact":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/testimage24.png","port":"20024","quorumtype":2}}`
	contractInput := `{"transfer_sample_nft":{"name": "rubix1", "nft_info": {"comment":"testing transfer","nft":"QmZ9jQTZJKq3LjKEZb8rdWuvgnc8DHzfvZULY9s5zch4Xw","nft_data":"","nft_value": 1,"owner": "bafybmihxaehnreq4ygnq3re3soob5znuj7hxoku6aeitdukif75umdv2nu", "receiver": "bafybmienjpoihwu2y6grilbvbrrqhleoifb3irz3gu2savjmjivzqw7424"}}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("unable to execute `transfer_sample_nft` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("transfer_sample_nft Result: ", result)
}

func main() {
	// Rubix Node Configs
	nodeAddress := "http://localhost:20006"
	quorumType := 2
	// Create Import function registry
	hostFnRegistry := wasmbridge.NewHostFunctionRegistry()

	// Initialize the WASM module
	wasmModule, err := wasmbridge.NewWasmModule(
		NFT_CONTRACT_WASM,
		hostFnRegistry,
		wasmbridge.WithRubixNodeAddress(nodeAddress),
		wasmbridge.WithQuorumType(quorumType),
	)
	if err != nil {
		log.Fatalf("Failed to initialize WASM module: %v", err)
	}

	//	mintNFTFunc(wasmModule)
	transferNFTFunc(wasmModule)
}
