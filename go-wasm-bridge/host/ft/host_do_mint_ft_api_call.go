package ft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"net/url"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/host"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
	wasmContext "github.com/rubixchain/rubix-wasm/go-wasm-bridge/context"
)

type DoMintFTApiCall struct {
	allocFunc   *wasmtime.Func
	memory      *wasmtime.Memory
	nodeAddress string
	quorumType  int
}

type MintFTData struct {
	Did             string `json:"did"`
	FtCount         int32  `json:"ft_count"`
	FtName          string `json:"ft_name"`
	FtNumStartIndex int32  `json:"ft_num_start_index"`
	TokenCount      int32  `json:"token_count"`
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

func (h *DoMintFTApiCall) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int, wasmCtx *wasmContext.WasmContext) {
	h.allocFunc = allocFunc
	h.memory = memory
	h.nodeAddress = nodeAddress
	h.quorumType = quorumType
}

func (h *DoMintFTApiCall) Callback() host.HostFunctionCallBack {
	return h.callback
}

func callCreateFTAPI(nodeAddress string, mintFTdata MintFTData) (string, error) {
	fmt.Println("The body in create-ft api :", mintFTdata)
	requestBody, err := json.Marshal(mintFTdata)
	if err != nil {
		fmt.Println("Error marshalling mintFTdata :", err)
		return "", err
	}

	// Create the request URL
	requestURL, err := url.JoinPath(nodeAddress, "/api/create-ft")
	if err != nil {
		return "", err
	}

	// Add ftNumStartIndex as a query parameter
	u, err := url.Parse(requestURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}
	query := u.Query()
	query.Set("ftNumStartIndex", fmt.Sprintf("%d", mintFTdata.FtNumStartIndex))
	u.RawQuery = query.Encode()
	finalURL := u.String()

	req, err := http.NewRequest("POST", finalURL, bytes.NewBuffer(requestBody))
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

	return utils.SignatureResponse(id, nodeAddress)
}

func (h *DoMintFTApiCall) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
	// Validate the number of arguments
	inputArgs, outputArgs := utils.HostFunctionParamExtraction(args, true, true)

	// Extract input bytes and convert to string
	inputBytes, memory, err := utils.ExtractDataFromWASM(caller, inputArgs)
	if err != nil {
		fmt.Println("Failed to extract data from WASM", err)
		return utils.HandleError(err.Error())
	}
	h.memory = memory // Assign memory to Host struct for future use

	var mintFTData MintFTData
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(inputBytes, &mintFTData)
	if err3 != nil {
		fmt.Println("Error unmarshaling mintftdata in callback function:", err3)
		return utils.HandleError(err3.Error())
	}

	callCreateFTAPIResp, err := callCreateFTAPI(h.nodeAddress, mintFTData)
	if err != nil {
		fmt.Println("Error calling CreateFTAPI in callback function:", err)
		return utils.HandleError(err.Error())
	}
	fmt.Println("The api response from create ft api :", callCreateFTAPIResp)

	err = utils.UpdateDataToWASM(caller, h.allocFunc, callCreateFTAPIResp, outputArgs)
	if err != nil {
		fmt.Println("Failed to update data to WASM", err)
		return utils.HandleError(err.Error())
	}

	return utils.HandleOk() // Success
}
