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
type RevealHighestBidReq struct {
	DeployerPassword string `json:"deployer_password"`
}

const BIDDING_CONTRACT_WASM = "../../../artifacts/bidding_contract.wasm"

func main() {
	// Create host function registry and register GetSomeBid and SaveSomeBid
	// host functions
	hostRegistry := wasmbridge.NewHostFunctionRegistry()

	hostRegistry.Register(bidContractHost.NewDecryptBid())
	hostRegistry.Register(bidContractHost.NewSaveBidInfo())
	hostRegistry.Register(bidContractHost.NewGetBidInfo())
	
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
	encryptedBidStr := hex.EncodeToString(encryptedBid)

	placeBidReq := PlaceBidReq{
		BidderDid:          "b1234did",
		EncryptedBidAmount: encryptedBidStr,
	}

	placeBidReqBytes, err := json.Marshal(placeBidReq)
	if err != nil {
		fmt.Println("failed to marshal placeBidReq")
	}

	placeBidReqStr := string(placeBidReqBytes)
	placeBidReqContractMsg := fmt.Sprintf(`{"place_bid": %v}`, placeBidReqStr)

	fmt.Println(placeBidReqContractMsg)

	contractResult, err := wasmModule.CallFunction(placeBidReqContractMsg)
	if err != nil {
		log.Fatalf("function call failed: %v", err)
		return
	}
	fmt.Println("place_bid status: ", contractResult)

	revealHighestBidReq := RevealHighestBidReq{
		DeployerPassword: "mypassword",
	}
	revealHighestBidReqBytes, err := json.Marshal(revealHighestBidReq)
	if err != nil {
		fmt.Println("failed to marshal revealHighestBidReq")
	}

	revealHighestBidReqStr := string(revealHighestBidReqBytes)
	revealHighestBidContractMsg := fmt.Sprintf(`{"reveal_highest_bid": %v}`, revealHighestBidReqStr)

	fmt.Println(revealHighestBidContractMsg)

	finalContractResult, err := wasmModule.CallFunction(revealHighestBidContractMsg)
	if err != nil {
		log.Fatalf("reveal_highest_bid function call failed: %v", err)
		return
	}
	fmt.Println("Result: ", finalContractResult)

}

