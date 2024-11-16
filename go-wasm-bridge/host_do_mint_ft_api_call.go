package wasmbridge

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	"github.com/bytecodealliance/wasmtime-go"
)

type DoMintFTApiCall struct {
	allocFunc *wasmtime.Func
	memory    *wasmtime.Memory
}
type MintFTData struct {
	Did        string `json:"did"`
	FtCount    int32  `json:"ftcount"`
	FtName     string `json:"ftname"`
	TokenCount int32  `json:"tokencount"`
	Port       string `json:"port"`
}

func NewDoMintFTApiCall() *DoMintFTApiCall {
	return &DoMintFTApiCall{}
}

func (h *DoMintFTApiCall) Name() string {
	return "do_mint_ft"
}

func (h *DoMintFTApiCall) FuncType() *wasmtime.FuncType {
	return wasmtime.NewFuncType(
		[]*wasmtime.ValType{
			wasmtime.NewValType(wasmtime.KindI32), // input_ptr
			wasmtime.NewValType(wasmtime.KindI32), // input_len
			wasmtime.NewValType(wasmtime.KindI32), // resp_ptr_ptr
			wasmtime.NewValType(wasmtime.KindI32), // resp_len_ptr
		},
		[]*wasmtime.ValType{wasmtime.NewValType(wasmtime.KindI32)}, // return i32
	)
}

func (h *DoMintFTApiCall) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory) {
	h.allocFunc = allocFunc
	h.memory = memory
}

func (h *DoMintFTApiCall) Callback() HostFunctionCallBack {
	return h.callback
}

func callCreateFTAPI(mintFTdata MintFTData) (string, error) {
	fmt.Println("The body in create-ft api :", mintFTdata)
	requestBody, err := json.Marshal(mintFTdata)
	if err != nil {
		fmt.Println("Error marshalling mintFTdata :", err)
		return "", err
	}

	// Create the request URL
	url := fmt.Sprintf("http://localhost:%s/api/create-ft", mintFTdata.Port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request in mintft function:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error sending http request: %v\n", err))
		return "", err
	}

	defer resp.Body.Close()
	fmt.Println("The response after calling the api :", resp)

	fmt.Println("Response Status:", resp.Status)

	createFtResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return "", err
	}
	// Process the data as needed
	fmt.Println("Response Body in callTransferFTAPI :", string(createFtResponse))
	var response map[string]interface{}
	err3 := json.Unmarshal(createFtResponse, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
		return "", err3
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)
	signatureResponse := SignatureResponse(id, mintFTdata.Port)

	return signatureResponse, nil

}

func (h *DoMintFTApiCall) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
	// Validate the number of arguments
	if len(args) != 4 {
		errMsg := fmt.Sprintf("%v expects 4 arguments, got %d", h.Name(), len(args))
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	// Extract arguments
	inputPtr := args[0].I32()
	inputLen := args[1].I32()
	respPtrPtr := args[2].I32()
	respLenPtr := args[3].I32()

	// Access memory from the caller
	memory := caller.GetExport("memory").Memory()
	if memory == nil {
		errMsg := "memory export not found"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}
	h.memory = memory // Assign memory to Host struct for future use

	// Read the input string from WASM memory
	data := memory.UnsafeData(caller)
	if data == nil {
		errMsg := "Failed to get memory data"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	// Convert pointers to int for slicing
	inputStart := int(inputPtr)
	inputEnd := inputStart + int(inputLen)

	// Validate memory bounds
	if inputStart < 0 || inputEnd > len(data) {
		errMsg := "input exceeds memory bounds"
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	// Extract input bytes and convert to string
	inputBytes := data[inputStart:inputEnd]

	var mintFTData MintFTData
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(inputBytes, &mintFTData)
	if err3 != nil {
		fmt.Println("Error unmarshaling mintftdata in callback function:", err3)
	}

	callCreateFTAPIResp, err := callCreateFTAPI(mintFTData)
	if err != nil {
		fmt.Println("Error calling CreateFTAPI in callback function:", err)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap("failed to mint ft")
	}
	fmt.Println("The api response from create ft api :", callCreateFTAPIResp)
	responseStr := "success"
	respLen := int32(len(responseStr))
	result, err := h.allocFunc.Call(caller, respLen)
	if err != nil {
		fmt.Printf("Alloc call failed: %v\n", err)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Alloc call failed: %v\n", err))
	}
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
	copy(data[respPtr:], []byte(responseStr))

	// Write the response pointer back to WASM memory using Little Endian encoding
	binary.LittleEndian.PutUint32(data[respPtrPtr:], uint32(respPtr))

	// Write the response length back to WASM memory using Little Endian encoding
	binary.LittleEndian.PutUint32(data[respLenPtr:], uint32(respLen))

	return []wasmtime.Val{wasmtime.ValI32(0)}, nil // Success
}
