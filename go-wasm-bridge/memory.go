package wasmbridge

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
)

func ExtractDataFromWASM(caller *wasmtime.Caller, inputArg *utils.WasmArgInfo) ([]byte, *wasmtime.Memory, error) {
	// Access memory from the caller
	memory := caller.GetExport("memory").Memory()
	if memory == nil {
		errMsg := "memory export not found"
		return nil, nil, errors.New(errMsg)
	}

	// Read the input string from WASM memory
	wasmMemory := memory.UnsafeData(caller)
	if wasmMemory == nil {
		errMsg := "Failed to get memory data"
		return nil, nil, errors.New(errMsg)
	}

	// Convert pointers to int for slicing
	inputStart := int(inputArg.DataPtr)
	inputEnd := inputStart + int(inputArg.DataPtrSize)

	// Validate memory bounds
	if inputStart < 0 || inputEnd > len(wasmMemory) {
		errMsg := "input exceeds memory bounds"
		return nil, nil, errors.New(errMsg)
	}

	// Extract input bytes and convert to string
	inputBytes := wasmMemory[inputStart:inputEnd]

	return inputBytes, memory, nil
}

func UpdateDataToWASM(caller *wasmtime.Caller, allocFunction *wasmtime.Func, outputValue string, outputArg *utils.WasmArgInfo) error {
	outputValueLen := int32(len(outputValue))

	// Allocating memory for output
	memory := caller.GetExport("memory").Memory()
	if memory == nil {
		errMsg := "memory export not found"
		return errors.New(errMsg)
	}
	result, err := allocFunction.Call(caller, outputValueLen)
	if err != nil {
		return err
	}
	wasmMemory := memory.UnsafeData(caller)
	if wasmMemory == nil {
		errMsg := "Failed to get memory data"
		return fmt.Errorf("%s", errMsg)
	}

	// Type Cast the allocated pointer
	respPtr, ok := result.(int32)
	if !ok {
		errMsg := "Alloc function did not return i32"
		return fmt.Errorf("%s", errMsg)
	}

	// Get memory size to ensure we don't write out of bounds
	memSize := memory.DataSize(caller)
	if uint32(respPtr)+uint32(outputValueLen) > uint32(memSize) {
		errMsg := "Response exceeds memory bounds"
		return fmt.Errorf("%s", errMsg)
	}

	// Write response bytes to allocated memory
	copy(wasmMemory[respPtr:], []byte(outputValue))

	respPtrPtr := outputArg.DataPtr
	respLenPtr := outputArg.DataPtrSize

	// Write the response pointer back to WASM memory using Little Endian encoding
	binary.LittleEndian.PutUint32(wasmMemory[respPtrPtr:], uint32(respPtr))

	// Write the response length back to WASM memory using Little Endian encoding
	binary.LittleEndian.PutUint32(wasmMemory[respLenPtr:], uint32(outputValueLen))

	return nil
}
