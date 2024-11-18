package host

import (
	"bidding-contract/state"
	"encoding/binary"

	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
)

type GetBid struct {
	allocFunc *wasmtime.Func
	memory    *wasmtime.Memory
}

func NewGetBid() *GetBid {
	return &GetBid{}
}

func (h *GetBid) Name() string {
	return "get_bid_from_state"
}

func (h *GetBid) FuncType() *wasmtime.FuncType {
	return wasmtime.NewFuncType(
		[]*wasmtime.ValType{
			wasmtime.NewValType(wasmtime.KindI32), // out_bid_ptr
			wasmtime.NewValType(wasmtime.KindI32), // out_bid_len
		},
		[]*wasmtime.ValType{wasmtime.NewValType(wasmtime.KindI32)}, // return i32
	)
}

func (h *GetBid) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int) {
	h.allocFunc = allocFunc
	h.memory = memory
}

func (h *GetBid) Callback() wasmbridge.HostFunctionCallBack {
	return h.callback
}

func (h *GetBid) callback(
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
	bidPtr := args[0].I32()
	bidLenPtr := args[1].I32()

	// Access memory from the caller
	memory := caller.GetExport("memory").Memory()
	if memory == nil {
		errMsg := "memory export not found"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}
	h.memory = memory // Assign memory to Host struct for future use

	// Read the URL string from WASM memory
	data := memory.UnsafeData(caller)
	if data == nil {
		errMsg := "Failed to get memory data"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	// Get the state
	currentBid, err := state.GetCurrentBid()
	if err != nil {
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(err.Error())
	}
	currentBidStr := fmt.Sprintf("%.3f", currentBid)

	// Allocate memory in WASM for the response string
	respLen := int32(len(currentBidStr))
	result, err := h.allocFunc.Call(caller, respLen)
	if err != nil {
		fmt.Printf("Alloc call failed: %v\n", err)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Alloc call failed: %v\n", err))
	}

	// Type assertion to int32 as allocFunc is expected to return i32
	respPtr, ok := result.(int32)
	if !ok {
		errMsg := "Alloc function did not return i32"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	// Get memory size to ensure we don't write out of bounds
	memSize := memory.DataSize(caller)
	if uint32(respPtr)+uint32(respLen) > uint32(memSize) {
		errMsg := "Response exceeds memory bounds"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	// Write response bytes to allocated memory
	copy(data[respPtr:], []byte(currentBidStr))

	// Write the response pointer back to WASM memory using Little Endian encoding
	binary.LittleEndian.PutUint32(data[bidPtr:], uint32(respPtr))

	// Write the response length back to WASM memory using Little Endian encoding
	binary.LittleEndian.PutUint32(data[bidLenPtr:], uint32(respLen))

	return []wasmtime.Val{wasmtime.ValI32(0)}, nil // Success
}
