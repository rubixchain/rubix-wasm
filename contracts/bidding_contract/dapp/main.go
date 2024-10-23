package main

import (
	"fmt"
	"log"

	bidContractHost "bidding-contract/host"

	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)

const BIDDING_CONTRACT_WASM = "../../../artifacts/bidding_contract.wasm"

func main() {
	// Create host function registry and register GetSomeBid and SaveSomeBid
	// host functions
	hostRegistry := wasmbridge.NewHostFunctionRegistry()

	hostRegistry.Register(bidContractHost.NewGetBid())
	hostRegistry.Register(bidContractHost.NewSaveBid())

	// Initialize the WASM module
	wasmModule, err := wasmbridge.NewWasmModule(BIDDING_CONTRACT_WASM, hostRegistry)
	if err != nil {
		log.Fatalf("Failed to initialize WASM module: %v\n", err)
		return
	}

	contractInput := `{"place_bid": {"bid_amount": 104.51}}`

	contractResult, err := wasmModule.CallFunction(contractInput)
	if err != nil {
		log.Fatalf("function call failed: %v", err)
		return
	}

	fmt.Println("Result: ", contractResult)
}
