package wasmbridge

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)



func WriteString(h HostFunction, caller *wasmtime.Caller, str string) (int32, int32, error) {
	allocFunc := h.GetAllocFunc()
	memory := h.GetMemory()
	
	strBytes := []byte(str)
	strLen := int32(len(strBytes))

	// Allocate memory for the string
	result, err := allocFunc.Call(caller, strLen)
	if err != nil {
		return 0, 0, fmt.Errorf("alloc function call failed: %v", err)
	}

	respPtr, ok := result.(int32)
	if !ok {
		return 0, 0, fmt.Errorf("alloc function did not return i32")
	}

	// Ensure we don't write out of bounds
	memSize := memory.DataSize(caller)
	if uint32(respPtr)+uint32(strLen) > uint32(memSize) {
		return 0, 0, fmt.Errorf("response exceeds memory bounds")
	}

	// Write the string bytes to memory
	copy(memory.UnsafeData(caller)[respPtr:], strBytes)

	return respPtr, strLen, nil
}
