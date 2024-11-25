package host

import (
	// "bidding-contract/state"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
	ecies "github.com/ecies/go/v2"
	wasmbridge "github.com/rubixchain/rubix-wasm/go-wasm-bridge"
	seal "github.com/rubixchain/rubixgoplatform/crypto"

	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
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

func (h *DecryptBid) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory) {
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
	// Validate the number of arguments
	if len(args) != 4 {
		errMsg := fmt.Sprintf("%v expects 4 arguments, got %d", h.Name(), len(args))
		fmt.Println(errMsg)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(errMsg)
	}

	// Extract arguments
	inputPtr := args[0].I32()
	inputLen := args[1].I32()
	outputptrPtr := args[2].I32()
	outputlenPtr := args[3].I32()

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

	var contractInputMap map[string][]byte
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(inputBytes, &contractInputMap)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in callback function:", err3)
	}

	var bidderData BidderData
	err := json.Unmarshal(contractInputMap["place_bid"], &bidderData)
	if err != nil {
		fmt.Println("Error unmarshaling bidderData in callback function:", err)
	}

	encryptedBid := bidderData.Bid

	decryptedBid := eciesDecryption("/home/rubix/Sai-Rubix/rubix-wasm/contracts/bidding_contract/bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a/pvtKey.pem", encryptedBid)

	responseStr := decryptedBid

	// Allocate memory in WASM for the response string
	respLen := int32(len(responseStr))
	result, err := h.allocFunc.Call(caller, respLen)
	if err != nil {
		fmt.Printf("Alloc call failed: %v\n", err)
		return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Alloc call failed: %v\n", err))
	}

	// Type assertion to int32 as allocFunc is expected to return i32
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
	binary.LittleEndian.PutUint32(data[outputptrPtr:], uint32(respPtr))

	// Write the response length back to WASM memory using Little Endian encoding
	binary.LittleEndian.PutUint32(data[outputlenPtr:], uint32(respLen))

	return []wasmtime.Val{wasmtime.ValI32(0)}, nil // Success

}

// ConvertSecp256k1ToEcies converts a secp256k1 private key to an ECIES private key.
func ConvertSecp256k1privkeyToEcies(privKey *secp256k1.PrivateKey) (*ecies.PrivateKey, error) {
	// Serialize the private key to get the private scalar bytes
	privKeyBytes := privKey.Serialize()

	// Convert the private scalar bytes to a big.Int
	d := new(big.Int).SetBytes(privKeyBytes)
	// Create an ECIES public key from the secp256k1 public key
	pubKey := privKey.PubKey()
	eciesPubKey := &ecies.PublicKey{
		X:     pubKey.X(),
		Y:     pubKey.Y(),
		Curve: secp256k1.S256(),
	}

	// Create an ECIES private key from the D value and the ECIES public key
	eciesPrivKey := &ecies.PrivateKey{
		PublicKey: eciesPubKey,
		D:         d,
	}

	return eciesPrivKey, nil
}
func eciesDecryption(privkey_path string, encrypted_data []byte) (plaintext string) {
	read_encodedprivkey, err := os.ReadFile(privkey_path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("privatekey which is read from given privkey.pem file is ", read_encodedprivkey)
	pemdecoded_privkey, rest := pem.Decode(read_encodedprivkey)
	fmt.Println("pemdecoded privkey is ", pemdecoded_privkey)
	fmt.Println("rest part while pem decoding privkey is ", rest)
	password := "mypassword"
	unsealedprivkey, err := seal.UnSeal(password, (pemdecoded_privkey).Bytes)
	fmt.Println("Decrypted Private key is ", unsealedprivkey)
	parsedprivkey := secp256k1.PrivKeyFromBytes(unsealedprivkey)
	ecies_privkey, err := ConvertSecp256k1privkeyToEcies(parsedprivkey)
	plaintext_bytes, err := ecies.Decrypt(ecies_privkey, encrypted_data)
	plaintext_string := string(plaintext_bytes)
	return plaintext_string
}
