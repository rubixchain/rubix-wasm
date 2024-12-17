package host

import (
	"encoding/json"
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"

	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge/host"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
)

type DecryptBid struct {
	allocFunc *wasmtime.Func
	memory    *wasmtime.Memory
}
type BidderData struct {
	Did string `json:"Did"`
	Bid []byte `json:"bid_amount"`
}

func NewDecryptBid() *DecryptBid {
	return &DecryptBid{}
}

func (h *DecryptBid) Name() string {
	return "ecies_decryption"
}

func (h *DecryptBid) FuncType() *wasmtime.FuncType {
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

func (h *DecryptBid) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int) {
	h.allocFunc = allocFunc
	h.memory = memory
}

func (h *DecryptBid) Callback() wasmbridge.HostFunctionCallBack {
	return h.callback
}
func (h *DecryptBid) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {

	inputArgs, outputArgs := utils.HostFunctionParamExtraction(args, true, true)

	// Extract input bytes and convert to string
	inputBytes, memory, err := utils.ExtractDataFromWASM(caller, inputArgs)
	if err != nil {
		fmt.Println("Failed to extract data from WASM", err)
		return utils.HandleError(err.Error())
	}
	h.memory = memory

	type DecryptionInputData struct {
		Privatekey_path string `json:"Privatekey_path"`
		Data            []byte `json:"data"`
	}

	var contractInputMap DecryptionInputData
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(inputBytes, &contractInputMap)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in callback function:", err3)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error unmarshaling response in callback function: %v", err3))
	}

	encryptedBid := contractInputMap.Data

	fmt.Println("Message: ", contractInputMap)

	decryptedBid, err := EciesDecryption("/home/rubix/Sai-Rubix/rubix-wasm/contracts/bidding_contract/bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a/pvtKey.pem", encryptedBid)
	if err != nil {
		fmt.Println("err")
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("unable to get decrypted string: %v", err))
	}

	if len(decryptedBid) == 0 {
		fmt.Println("Unable to get the decrypted Bid")
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap("Unable to get the decrypted Bid")
	}

	responseStr := decryptedBid

	err = utils.UpdateDataToWASM(caller, h.allocFunc, responseStr, outputArgs)
	if err != nil {
		fmt.Println("Failed to update data to WASM", err)
		return utils.HandleError(err.Error())
	}

	return utils.HandleOk() // Success

}
