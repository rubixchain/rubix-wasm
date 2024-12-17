package host

import (
	"bidding-contract/state"
	"encoding/json"
	"fmt"
	"github.com/bytecodealliance/wasmtime-go"
	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge/host"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
)

type SaveBidInfo struct {
	allocFunc *wasmtime.Func
	memory    *wasmtime.Memory
}

func NewSaveBidInfo() *SaveBidInfo {
	return &SaveBidInfo{}
}

func (h *SaveBidInfo) Name() string {
	return "save_bidding_info"
}

func (h *SaveBidInfo) FuncType() *wasmtime.FuncType {
	return wasmtime.NewFuncType(
		[]*wasmtime.ValType{
			wasmtime.NewValType(wasmtime.KindI32), // out_bid_ptr
			wasmtime.NewValType(wasmtime.KindI32), // out_bid_len
		},
		[]*wasmtime.ValType{wasmtime.NewValType(wasmtime.KindI32)}, // return i32
	)
}

func (h *SaveBidInfo) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int) {
	h.allocFunc = allocFunc
	h.memory = memory
}

func (h *SaveBidInfo) Callback() wasmbridge.HostFunctionCallBack {
	return h.callback
}

func (h *SaveBidInfo) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
	inputArgs, _ := utils.HostFunctionParamExtraction(args, true, false)
	
	// Extract input bytes and convert to string
	inputBytes,memory, err := utils.ExtractDataFromWASM(caller, inputArgs)
	if err != nil {
		fmt.Println("Failed to extract data from WASM", err)
		return utils.HandleError(err.Error())
	}
	h.memory = memory 
	// Deserialize the struct (assuming JSON format)
	var placeBidReq struct {
		BidderDID          string `json:"bidder_did"`
		EncryptedBidAmount string `json:"encrypted_bid_amount"`
	}

	if err := json.Unmarshal(inputBytes, &placeBidReq); err != nil {
		errMsg := fmt.Sprintf("Failed to deserialize PlaceBidReq: %v", err)
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	fmt.Printf("Received PlaceBidReq: %+v\n", placeBidReq)

	// Save the struct
	stateSaveErr := state.SavePlaceBidReqState(placeBidReq)
	if stateSaveErr != nil {
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(stateSaveErr.Error())
	}

	return []wasmtime.Val{wasmtime.ValI32(0)}, nil // Success
}
