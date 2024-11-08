package main

import (
	"fmt"
	"log"

	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/wasmbridge"
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
	contractInput := `{"mint_sample_nft":{"name": "rubix1", "nft_info": {"did":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","metadata":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/metadata.json","artifact":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/testimage25.png","port":"20024","quorumtype":2}}}`

	result, err := executeAndGetContractResult(wasmModule, contractInput)
	if err != nil {
		fmt.Printf("unable to execute `mint_sample_nft` Contract function, error: %v\n", err)
		return
	}
	fmt.Println("mint_sample_nft Result: ", result)
}
func main() {
	// Create Import function registry
	hostFnRegistry := wasmbridge.NewHostFunctionRegistry()

	// Initialize the WASM module
	wasmModule, err := wasmbridge.NewWasmModule(NFT_CONTRACT_WASM, hostFnRegistry)
	if err != nil {
		log.Fatalf("Failed to initialize WASM module: %v", err)
	}

	mintNFTFunc(wasmModule)

}
