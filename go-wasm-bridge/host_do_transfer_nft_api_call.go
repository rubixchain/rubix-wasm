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

type TransferNFTData struct {
	NFT        string  `json:"nft"`
	Owner      string  `json:"owner"`
	Receiver   string  `json:"receiver"`
	Comment    string  `json:"comment"`
	NFTValue   float64 `json:"nft_value"`
	NFTData    string  `json:"nft_data"`
	QuorumType int32   `json:"quorum_type"`
}

type DoTransferNFTApiCall struct {
	allocFunc   *wasmtime.Func
	memory      *wasmtime.Memory
	nodeAddress string
	quorumType  int
}

func NewDoTransferNFTApiCall() *DoTransferNFTApiCall {
	return &DoTransferNFTApiCall{}
}
func (h *DoTransferNFTApiCall) Name() string {
	return "do_transfer_nft"
}
func (h *DoTransferNFTApiCall) FuncType() *wasmtime.FuncType {
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

func (h *DoTransferNFTApiCall) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int) {
	h.allocFunc = allocFunc
	h.memory = memory
	h.nodeAddress = nodeAddress
	h.quorumType = quorumType
}

func (h *DoTransferNFTApiCall) Callback() HostFunctionCallBack {
	return h.callback
}
func callTransferNFTAPI(nodeAddress string, quorumType int, transferNFTdata TransferNFTData) error {
	transferNFTdata.QuorumType = int32(quorumType)
	fmt.Println("printing the data in callTransferNFTAPI function is:", transferNFTdata)
	bodyJSON, err := json.Marshal(transferNFTdata)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	transferNFTUrl, err := url.JoinPath(nodeAddress, "/api/execute-nft")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", transferNFTUrl, bytes.NewBuffer(bodyJSON))
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
	fmt.Println("Response Status in callTransferNFTAPI:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return err
	}
	// Process the data as needed
	fmt.Println("Response Body in callTransferNFTAPI :", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
		return err3
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)

	defer resp.Body.Close()

	_, err = SignatureResponse(id, nodeAddress)
	return err
}

func (h *DoTransferNFTApiCall) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
	inputArgs, outputArgs := utils.HostFunctionParamExtraction(args, true, true)

	// Extract input bytes and convert to string
	inputBytes, memory, err := ExtractDataFromWASM(caller, inputArgs)
	if err != nil {
		fmt.Println("Failed to extract data from WASM", err)
		return utils.HandleError(err.Error())
	}
	h.memory = memory // Assign memory to Host struct for future use
	var transferNFTData TransferNFTData

	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(inputBytes, &transferNFTData)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in callback function:", err3)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error unmarshaling response in callback function:", err3))
	}
	callTransferNFTAPIRespErr := callTransferNFTAPI(h.nodeAddress, h.quorumType, transferNFTData)
	if callTransferNFTAPIRespErr != nil {
		fmt.Println("failed to transfer NFT", callTransferNFTAPIRespErr)
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
