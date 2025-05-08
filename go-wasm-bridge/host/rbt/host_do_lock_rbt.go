package rbt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bytecodealliance/wasmtime-go"
	wasmContext "github.com/rubixchain/rubix-wasm/go-wasm-bridge/context"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/host"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
)

type LockRBTData struct {
	Did          string `json:"did"`
	Amount       uint64 `json:"amount"`
	TokenAddress string `json:"token_address"`
	QuorumType   int32  `json:"quorum_type"`
}

type LockRBTResponse struct {
	TransactionId string `json:"transaction_id"`
	Status        bool   `json:"status"`
	Message       string `json:"message"`
}

type DoLockRBTApiCall struct {
	allocFunc   *wasmtime.Func
	memory      *wasmtime.Memory
	nodeAddress string
	quorumType  int
}

func NewDoLockRBTApiCall() *DoLockRBTApiCall {
	return &DoLockRBTApiCall{}
}

func (h *DoLockRBTApiCall) Name() string {
	return "do_lock_rbt"
}

func (h *DoLockRBTApiCall) FuncType() *wasmtime.FuncType {
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

func (h *DoLockRBTApiCall) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int, wasmCtx *wasmContext.WasmContext) {
	h.allocFunc = allocFunc
	h.memory = memory
	h.nodeAddress = nodeAddress
	h.quorumType = quorumType
}

func (h *DoLockRBTApiCall) Callback() host.HostFunctionCallBack {
	return h.callback
}

func callLockRBTAPI(nodeAddress string, quorumType int, lockRBTData LockRBTData) (string, error) {
	lockRBTData.QuorumType = int32(quorumType)

	bodyJSON, err := json.Marshal(lockRBTData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}

	// This endpoint would need to be implemented in the Rubix node
	lockRBTUrl, err := url.JoinPath(nodeAddress, "/api/lock-rbt")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", lockRBTUrl, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return "", err
	}

	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return "", err
	}

	// Handle signature if needed, following the pattern of other API calls
	if result, ok := response["result"].(map[string]interface{}); ok {
		if id, ok := result["id"].(string); ok {
			return utils.SignatureResponse(id, nodeAddress)
		}
	}

	return string(respBody), nil
}

func (h *DoLockRBTApiCall) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
	inputArgs, outputArgs := utils.HostFunctionParamExtraction(args, true, true)

	inputBytes, memory, err := utils.ExtractDataFromWASM(caller, inputArgs)
	if err != nil {
		fmt.Println("Failed to extract data from WASM", err)
		return utils.HandleError(err.Error())
	}
	h.memory = memory

	var lockRBTData LockRBTData
	err = json.Unmarshal(inputBytes, &lockRBTData)
	if err != nil {
		fmt.Println("Error unmarshaling input data:", err)
		return utils.HandleError(err.Error())
	}

	lockRBTResponse, err := callLockRBTAPI(h.nodeAddress, h.quorumType, lockRBTData)
	if err != nil {
		fmt.Println("Failed to lock RBT:", err)
		return utils.HandleError(err.Error())
	}

	err = utils.UpdateDataToWASM(caller, h.allocFunc, lockRBTResponse, outputArgs)
	if err != nil {
		fmt.Println("Failed to update data to WASM", err)
		return utils.HandleError(err.Error())
	}

	return utils.HandleOk()
}
