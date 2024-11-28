package wasmbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
)

type TransferFTData struct {
	FTCount    int32  `json:"ft_count"`
	FTName     string `json:"ft_name"`
	CreatorDID string `json:"creatorDID"`
	QuorumType int32  `json:"quorum_type"`
	Comment    string `json:"comment"`
	Receiver   string `json:"receiver"`
	Sender     string `json:"sender"`
}

type DoTransferFTApiCall struct {
	allocFunc   *wasmtime.Func
	memory      *wasmtime.Memory
	nodeAddress string
	quorumType  int
}

func NewDoTransferFTApiCall() *DoTransferFTApiCall {
	return &DoTransferFTApiCall{}
}
func (h *DoTransferFTApiCall) Name() string {
	return "do_transfer_ft"
}
func (h *DoTransferFTApiCall) FuncType() *wasmtime.FuncType {
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

func (h *DoTransferFTApiCall) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int) {
	h.allocFunc = allocFunc
	h.memory = memory
	h.nodeAddress = nodeAddress
	h.quorumType = quorumType
}

func (h *DoTransferFTApiCall) Callback() HostFunctionCallBack {
	return h.callback
}
func callTransferFTAPI(nodeAddress string, quorumType int, transferFTdata TransferFTData) error {
	transferFTdata.QuorumType = int32(quorumType)
	bodyJSON, err := json.Marshal(transferFTdata)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	transferFTUrl, err := url.JoinPath(nodeAddress, "/api/initiate-ft-transfer")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", transferFTUrl, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status in callTransferFTAPI:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return err
	}
	// Process the data as needed
	fmt.Println("Response Body in callTransferFTAPI :", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
		return err3
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)

	_, err = SignatureResponse(id, nodeAddress)
	return err
}

func (h *DoTransferFTApiCall) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
	// Validate the number of arguments
	inputArgs, outputArgs := utils.HostFunctionParamExtraction(args, true, true)

	// Extract input bytes and convert to string
	inputBytes, memory, err := ExtractDataFromWASM(caller, inputArgs)
	if err != nil {
		fmt.Println("Failed to extract data from WASM", err)
		return utils.HandleError(err.Error())
	}
	h.memory = memory // Assign memory to Host struct for future use
	var transferFTData TransferFTData

	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(inputBytes, &transferFTData)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in callback function:", err3)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error unmarshaling response in callback function:", err3))
	}
	callTransferFTAPIRespErr := callTransferFTAPI(h.nodeAddress, h.quorumType, transferFTData)

	if callTransferFTAPIRespErr != nil {
		fmt.Println("failed to transfer NFT", callTransferFTAPIRespErr)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap("failed to transfer NFT")
	}

	responseStr := "success"
	wasmInput := WasmInput{
		Caller:        caller,
		AllocFunction: h.allocFunc,
		Memory:        memory,
		OutputValue:   responseStr,
	}
	err = UpdateDataToWASM(&wasmInput, outputArgs)
	if err != nil {
		fmt.Println("Failed to update data to WASM", err)
		return utils.HandleError(err.Error())
	}

	return utils.HandleOk() // Success

}
