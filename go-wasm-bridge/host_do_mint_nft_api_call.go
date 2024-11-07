package wasmbridge

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
)

type DoMintNFTApiCall struct {
	allocFunc *wasmtime.Func
	memory    *wasmtime.Memory
}
type MintNFTdata struct {
	Did        string `json:"did"`
	Metadata   string `json:"metadata"`
	Artifact   string `json:"artifact"`
	Port       string `json:"port"`
	QuorumType int32  `json:"Quorumtype"`
}
type deployTokendata struct {
	Nft        string `json:"nft"`
	Did        string `json:"Did"`
	QuorumType int32  `json:"Quorumtype"`
}

func NewDoMintNFTApiCall() *DoMintNFTApiCall {
	return &DoMintNFTApiCall{}
}

func (h *DoMintNFTApiCall) Name() string {
	return "do_mint_nft"
}

func (h *DoMintNFTApiCall) FuncType() *wasmtime.FuncType {
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

func (h *DoMintNFTApiCall) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory) {
	h.allocFunc = allocFunc
	h.memory = memory
}

func (h *DoMintNFTApiCall) Callback() HostFunctionCallBack {
	return h.callback
}
func callCreateNFTAPI(mintNFTdata MintNFTdata) []byte {
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)
	fmt.Println("Printing the data in generateToken function", mintNFTdata)

	// Add form fields (simple text fields)
	writer.WriteField("did", mintNFTdata.Did)
	// writer.WriteField("UserId", mintNFTdata.Userid)

	// Add the NFTFile to the form
	fmt.Println("NFT file name is:", mintNFTdata.Artifact)
	nftFile, err := os.Open(mintNFTdata.Artifact)
	if err != nil {
		fmt.Println("Error opening NFT file:", err)

	}
	defer nftFile.Close()

	// Add the NFTFile part to the form
	nftFormFile, err := writer.CreateFormFile("NFTFile", mintNFTdata.Artifact)
	if err != nil {
		fmt.Println("Error creating NFT form file:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error copying NFT file content: %v\n", err))

	}

	_, err = io.Copy(nftFormFile, nftFile)
	if err != nil {
		fmt.Println("Error copying NFT file content:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error copying NFT file content: %v\n", err))

	}

	// Add the NFTFileInfo to the form
	fmt.Println("NFTFileInfo file name is:", mintNFTdata.Metadata)
	nftFileInfo, err := os.Open(mintNFTdata.Metadata)
	if err != nil {
		fmt.Println("Error opening NFTFileInfo file:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error opening NFTFileInfo file: %v\n", err))

	}
	defer nftFileInfo.Close()

	// Add the NFTFileInfo part to the form
	nftInfoFormFile, err := writer.CreateFormFile("metadata", mintNFTdata.Metadata)
	if err != nil {
		fmt.Println("Error creating NFTFileInfo form file:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error creating NFTFileInfo form file: %v\n", err))

	}

	_, err = io.Copy(nftInfoFormFile, nftFileInfo)
	if err != nil {
		fmt.Println("Error copying NFTFileInfo content:", err)
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart writer:", err)

	}

	// Create the request URL
	url := fmt.Sprintf("http://localhost:%s/api/create-nft", mintNFTdata.Port)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Failed to create HTTP request: %v\n", err))

	}

	// Set the Content-Type header to multipart/form-data with the correct boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request in generateToken fun:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error sending http request: %v\n", err))

	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Read and print the response body
	createNFTAPIResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Failed to read response body: %v\n", err))

	}

	fmt.Println("Response Body:", string(createNFTAPIResponse))
	defer resp.Body.Close()

	return createNFTAPIResponse

}
func callDeployNFTAPI(mintNFTData MintNFTdata, nftId string) error {
	var input deployTokendata
	input.Did = mintNFTData.Did
	input.Nft = nftId
	input.QuorumType = mintNFTData.QuorumType

	bodyJSON, err := json.Marshal(input)
	if err != nil {
		fmt.Println("error in marshaling JSON:", err)
		return err
	}
	url := fmt.Sprintf("http://localhost:%s/api/deploy-nft", mintNFTData.Port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
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
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return err
	}
	// Process the data as needed
	fmt.Println("Response Body in DeployNft :", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
		return err3
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)
	SignatureResponse(id, mintNFTData.Port)

	defer resp.Body.Close()
	return nil

}
func SignatureResponse(requestId string, port string) string {
	data := map[string]interface{}{
		"id":       requestId,
		"mode":     0,
		"password": "mypassword",
	}

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		//	return
	}
	url := fmt.Sprintf("http://localhost:%s/api/signature-response", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		//return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		//return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		//return
	}
	// Process the data as needed
	fmt.Println("Response Body in signature response :", string(data2))
	//json encode string
	defer resp.Body.Close()
	return string(data2)
}

func (h *DoMintNFTApiCall) callback(
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
	// input := string(inputBytes)
	var inputData MintNFTdata
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(inputBytes, &inputData)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in mintToken:", err3)
	}
	// Create a buffer to hold the multipart form data

	// // Make HTTP GET request to the provided input
	// resp, err := http.Get(input)
	// if err != nil {
	// 	fmt.Printf("HTTP request failed: %v\n", err)
	// 	return []wasmtime.Val{wasmtime.ValI32(1)}, nil
	// }
	// defer resp.Body.Close()

	// // Read response body
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Printf("Failed to read response body: %v\n", err)
	// 	return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Failed to read response body: %v\n", err))
	// }
	data2 := callCreateNFTAPI(inputData)
	var unmarshaledResponse map[string]interface{}
	err := json.Unmarshal(data2, &unmarshaledResponse)
	if err != nil {
		fmt.Println("Error in unmarshaling data2:", err)

	}
	respResult := unmarshaledResponse["result"].(string)
	fmt.Println("response result:", respResult)
	// inputData.Nft = respResult
	callDeployNFTAPI(inputData, respResult)
	responseStr := data2

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
	binary.LittleEndian.PutUint32(data[respPtrPtr:], uint32(respPtr))

	// Write the response length back to WASM memory using Little Endian encoding
	binary.LittleEndian.PutUint32(data[respLenPtr:], uint32(respLen))

	return []wasmtime.Val{wasmtime.ValI32(0)}, nil // Success
}
