package main

import (
	bidContractHost "bidding-contract/host"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)

type PlaceBidReq struct {
	BidderDid          string `json:"bidder_did"`
	EncryptedBidAmount string `json:"encrypted_bid_amount"`
}

const BIDDING_CONTRACT_WASM = "../../../artifacts/bidding_contract.wasm"

func main() {
	// Create host function registry and register GetSomeBid and SaveSomeBid
	// host functions
	hostRegistry := wasmbridge.NewHostFunctionRegistry()

	hostRegistry.Register(bidContractHost.NewGetBid())
	hostRegistry.Register(bidContractHost.NewSaveBid())
	hostRegistry.Register(bidContractHost.NewDecryptBid())

	// Initialize the WASM module
	wasmModule, err := wasmbridge.NewWasmModule(BIDDING_CONTRACT_WASM, hostRegistry)
	if err != nil {
		log.Fatalf("Failed to initialize WASM module: %v\n", err)
		return
	}

	inputBidAmount := 119.34
	inputBid := fmt.Sprintf(`{"bid_amount": %v}`, inputBidAmount)
	inputBidBytes := []byte(inputBid)

	pubkeyPath := "/home/rubix/Sai-Rubix/rubix-wasm/contracts/bidding_contract/bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a/pubKey.pem"
	encryptedBid := bidContractHost.EciesEncryption(pubkeyPath, inputBidBytes)
	// encodedBidInString := base64.StdEncoding.EncodeToString(encryptedBid)
	encryptedBidStr := hex.EncodeToString(encryptedBid)

	placeBidReq := PlaceBidReq{
		BidderDid:          "bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a",
		EncryptedBidAmount: encryptedBidStr,
	}

	placeBidReqBytes, err := json.Marshal(placeBidReq)
	if err != nil {
		fmt.Println("failed to marshal placeBidReq")
	}

	placeBidReqStr := string(placeBidReqBytes)
	contractMsg := fmt.Sprintf(`{"place_bid": %v}`, placeBidReqStr)

	fmt.Println(contractMsg)

	contractResult, err := wasmModule.CallFunction(contractMsg)
	if err != nil {
		log.Fatalf("function call failed: %v", err)
		return
	}

	fmt.Println("Result: ", contractResult)
}
