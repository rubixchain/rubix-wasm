package utils

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)

type WasmArgInfo struct {
	DataPtr     int32
	DataPtrSize int32
}

func HostFunctionParamExtraction(args []wasmtime.Val, areInputArgsPresent bool, areOutputArgsPresent bool) (*WasmArgInfo, *WasmArgInfo) {
	var nArgs int = len(args)

	if ((nArgs % 2) != 0) || (nArgs == 0) || (nArgs > 4) {
		fmt.Printf("Invalied number of Arguments")
		return nil, nil
	}

	inputArg := &WasmArgInfo{}
	outputArg := &WasmArgInfo{}

	if areInputArgsPresent && areOutputArgsPresent {
		inputArg.DataPtr = args[0].I32()
		inputArg.DataPtrSize = args[1].I32()

		outputArg.DataPtr = args[2].I32()
		outputArg.DataPtrSize = args[3].I32()
	} else {
		if areInputArgsPresent {
			inputArg.DataPtr = args[0].I32()
			inputArg.DataPtrSize = args[1].I32()
		}
		if areOutputArgsPresent {
			outputArg.DataPtr = args[2].I32()
			outputArg.DataPtrSize = args[3].I32()
		}
	}

	return inputArg, outputArg
}
