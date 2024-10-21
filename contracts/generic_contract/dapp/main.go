// main.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/wasmbridge"
)

const GENERIC_CONTRACT_WASM = "../../../artifacts/generic_contract.wasm"


func executeAndGetContractResult(wasmModule *wasmbridge.WasmModule, contractInput string) (string, error) {
    // Call the function
    contractResult, err := wasmModule.CallFunction(contractInput)
    if err != nil {
        return "", fmt.Errorf("function call failed: %v", err)
    }


    return contractResult, nil
}

func callSummationOfThreeSumFunc(wasmModule *wasmbridge.WasmModule) {
    contractInput := `{"add_three_nums": {"a": 1, "b": 2, "c": 9883}}`

    result, err := executeAndGetContractResult(wasmModule, contractInput)
    if err != nil {
        fmt.Printf("unable to execute `add_three_sum` Contract function, error: %v\n", err)
        os.Exit(1)        
    }

    fmt.Println("add_three_sum Result: ", result)
}

func callTestVecFunc(wasmModule *wasmbridge.WasmModule) {
    contractInput := `{"test_vec": {"name_list": ["Arnab", "Ghose"]}}`

    result, err := executeAndGetContractResult(wasmModule, contractInput)
    if err != nil {
        fmt.Printf("unable to execute `add_three_sum` Contract function, error: %v\n", err)
        os.Exit(1)        
    }

    fmt.Println("test_vec Result: ", result)
}

func greetingsFunc(wasmModule *wasmbridge.WasmModule) {
    contractInput := `{"greetings": {"name": "Rubix"}}`

    result, err := executeAndGetContractResult(wasmModule, contractInput)
    if err != nil {
        fmt.Printf("unable to execute `add_three_sum` Contract function, error: %v\n", err)
        os.Exit(1)        
    }

    fmt.Println("greetings Result: ", result)
}

func makeSomeApiCallFunc(wasmModule *wasmbridge.WasmModule) {
    contractInput := `{"make_some_api_call": {"url": "https://httpbin.org/range/2?duration=5&chunk_size=10"}}`

    result, err := executeAndGetContractResult(wasmModule, contractInput)
    if err != nil {
        fmt.Printf("unable to execute `make_some_api_call` Contract function, error: %v\n", err)
        os.Exit(1)        
    }

    fmt.Println("make_some_api_call Result: ", result)
}

func main() {
    // Create Import function registry
    hostFnRegistry := wasmbridge.NewHostFunctionRegistry()
    
    // Initialize the WASM module
    wasmModule, err := wasmbridge.NewWasmModule(GENERIC_CONTRACT_WASM, hostFnRegistry)
    if err != nil {
        log.Fatalf("Failed to initialize WASM module: %v", err)
    }

    greetingsFunc(wasmModule)
    callSummationOfThreeSumFunc(wasmModule)
    makeSomeApiCallFunc(wasmModule)
    callTestVecFunc(wasmModule)
}