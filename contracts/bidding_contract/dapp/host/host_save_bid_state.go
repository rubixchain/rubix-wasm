package host

import (
	"bidding-contract/state"
	"strconv"

	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)

type SaveBid struct {
	allocFunc *wasmtime.Func
	memory    *wasmtime.Memory
}

func NewSaveBid() *SaveBid {
	return &SaveBid{}
}

func (h *SaveBid) Name() string {
	return "save_bid_to_state"
}

func (h *SaveBid) FuncType() *wasmtime.FuncType {
	return wasmtime.NewFuncType(
		[]*wasmtime.ValType{
			wasmtime.NewValType(wasmtime.KindI32), // out_bid_ptr
			wasmtime.NewValType(wasmtime.KindI32), // out_bid_len
		},
		[]*wasmtime.ValType{wasmtime.NewValType(wasmtime.KindI32)}, // return i32
	)
}

func (h *SaveBid) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int) {
	h.allocFunc = allocFunc
	h.memory = memory
}

func (h *SaveBid) Callback() wasmbridge.HostFunctionCallBack {
	return h.callback
}

func (h *SaveBid) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
	// Validate the number of arguments
	if len(args) != 2 {
		errMsg := fmt.Sprintf("%v expects 2 arguments, got %d", h.Name(), len(args))
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	// Extract arguments
	inBidPtr := args[0].I32()
	inBidLenPtr := args[1].I32()

	// Access memory from the caller
	memory := caller.GetExport("memory").Memory()
	if memory == nil {
		errMsg := "memory export not found"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}
	h.memory = memory // Assign memory to Host struct for future use

	// Read the input bid string from WASM memory
	data := memory.UnsafeData(caller)
	if data == nil {
		errMsg := "Failed to get memory data"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	bidAmountStart := int(inBidPtr)
	bidAmountEnd := bidAmountStart + int(inBidLenPtr)

	if bidAmountStart < 0 || bidAmountEnd > len(data) {
		errMsg := "Input bid amount exceeds memory bounds"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	inputBidStr := string(data[bidAmountStart:bidAmountEnd])

	inputBidAmt, err := strconv.ParseFloat(inputBidStr, 64)
	if err != nil {
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(err.Error())
	}

	// Save the state
	stateSaveErr := state.SaveIncomingBid(inputBidAmt)
	if stateSaveErr != nil {
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(stateSaveErr.Error())
	}

	return []wasmtime.Val{wasmtime.ValI32(0)}, nil // Success
}
