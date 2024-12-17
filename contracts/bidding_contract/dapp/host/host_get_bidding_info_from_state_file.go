package host

import (
	"bidding-contract/state"
	"encoding/json"
	"fmt"
	"github.com/bytecodealliance/wasmtime-go"
	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge/host"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
)

type GetBidInfo struct {
	allocFunc *wasmtime.Func
	memory    *wasmtime.Memory
}

// NewGetBidInfo creates and returns a new instance of GetBidInfo.
func NewGetBidInfo() *GetBidInfo {
	return &GetBidInfo{}
}

// Name returns the name of the host function implemented.
func (h *GetBidInfo) Name() string {
	return "get_bidding_info_from_state_file"
}

func (h *GetBidInfo) FuncType() *wasmtime.FuncType {
	return wasmtime.NewFuncType(
		[]*wasmtime.ValType{
			wasmtime.NewValType(wasmtime.KindI32), // out_bid_ptr
			wasmtime.NewValType(wasmtime.KindI32), // out_bid_len
		},
		[]*wasmtime.ValType{wasmtime.NewValType(wasmtime.KindI32)}, // return i32
	)
}

func (h *GetBidInfo) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int) {
	h.allocFunc = allocFunc
	h.memory = memory
}

func (h *GetBidInfo) Callback() wasmbridge.HostFunctionCallBack {
	return h.callback
}

func (h *GetBidInfo) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {

	inputArgs, outputArgs := utils.HostFunctionParamExtraction(args, false, true)

	_, memory, err := utils.ExtractDataFromWASM(caller, inputArgs)
	if err != nil {
		fmt.Println("Failed to extract data from WASM", err)
		return utils.HandleError(err.Error())
	}

	h.memory = memory

	// Get all the bid information from the statefile
	bidInfoFromStateFile, err := state.GetBidInfo()
	if err != nil {
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(err.Error())
	}

	bidInfoFromStateFileBytes, err := json.Marshal(bidInfoFromStateFile)
	if err != nil {
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("failed to marshall bidinfo from state file, err: %v", err.Error()))
	}
	bidInfoStr := string(bidInfoFromStateFileBytes)

	err = utils.UpdateDataToWASM(caller, h.allocFunc, bidInfoStr, outputArgs)
	if err != nil {
		fmt.Println("Failed to update data to WASM", err)
		return utils.HandleError(err.Error())
	}

	return utils.HandleOk() // Success
}
