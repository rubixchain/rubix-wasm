package nft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/host"
	"github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
	wasmContext "github.com/rubixchain/rubix-wasm/go-wasm-bridge/context"
)

type DoMintNFTApiCall struct {
	allocFunc   *wasmtime.Func
	memory      *wasmtime.Memory
	nodeAddress string
	quorumType  int
}

type MintNFTData struct {
	Did      string `json:"did"`
	Metadata string `json:"metadata"`
	Artifact string `json:"artifact"`
}

type deployNFTReq struct {
	Nft        string `json:"nft"`
	Did        string `json:"did"`
	QuorumType int32  `json:"quorum_type"`
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

func (h *DoMintNFTApiCall) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int, wasmCtx *wasmContext.WasmContext) {
	h.allocFunc = allocFunc
	h.memory = memory
	h.nodeAddress = nodeAddress
	h.quorumType = quorumType
}

func (h *DoMintNFTApiCall) Callback() host.HostFunctionCallBack {
	return h.callback
}

func callCreateNFTAPI(nodeAddress string, mintNFTdata MintNFTData) []byte {
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Add form fields (simple text fields)
	writer.WriteField("did", mintNFTdata.Did)
	// writer.WriteField("UserId", mintNFTdata.Userid)

	// Add the NFTFile to the form
	fmt.Println("Artifact name is:", mintNFTdata.Artifact)
	nftArtifact, err := os.Open(mintNFTdata.Artifact)
	if err != nil {
		fmt.Println("Error opening Artifact file:", err)
		return nil
	}
	defer nftArtifact.Close()

	// Add the NFTFile part to the form
	nftArtifactFile, err := writer.CreateFormFile("artifact", mintNFTdata.Artifact)
	if err != nil {
		fmt.Println("Error creating NFT Artifact file:", err)
		return nil
	}

	_, err = io.Copy(nftArtifactFile, nftArtifact)
	if err != nil {
		fmt.Println("Error copying NFT file content:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error copying NFT file content: %v\n", err))
		return nil
	}

	// Add the NFTFileInfo to the form
	fmt.Println("Metadata file name is:", mintNFTdata.Metadata)
	metadataFileInfo, err := os.Open(mintNFTdata.Metadata)
	if err != nil {
		fmt.Println("Error opening Metadata file:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error opening NFTFileInfo file: %v\n", err))
		return nil
	}
	defer metadataFileInfo.Close()

	// Add the NFTFileInfo part to the form
	metadataFile, err := writer.CreateFormFile("metadata", mintNFTdata.Metadata)
	if err != nil {
		fmt.Println("Error creating NFTFileInfo form file:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error creating NFTFileInfo form file: %v\n", err))
		return nil
	}

	_, err = io.Copy(metadataFile, metadataFileInfo)
	if err != nil {
		fmt.Println("Error copying NFTFileInfo content:", err)
		return nil
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart writer:", err)
		return nil
	}

	// Create the request URL
	url, err := url.JoinPath(nodeAddress, "/api/create-nft")
	if err != nil {
		fmt.Println("Error forming url path for Create NFT API, err: ", err)
		return nil
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Failed to create HTTP request: %v\n", err))
		return nil
	}

	// Set the Content-Type header to multipart/form-data with the correct boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request in generateToken fun:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Error sending http request: %v\n", err))
		return nil
	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Read and print the response body
	createNFTAPIResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		// return []wasmtime.Val{wasmtime.ValI32(1)}, wasmtime.NewTrap(fmt.Sprintf("Failed to read response body: %v\n", err))
		return nil
	}

	fmt.Println("Response Body:", string(createNFTAPIResponse))
	defer resp.Body.Close()

	return createNFTAPIResponse

}

func callDeployNFTAPI(nodeAddress string, quorumType int, mintNFTData MintNFTData, nftId string) error {
	var deployReq deployNFTReq

	deployReq.Did = mintNFTData.Did
	deployReq.Nft = nftId
	deployReq.QuorumType = int32(quorumType)

	bodyJSON, err := json.Marshal(deployReq)
	if err != nil {
		fmt.Println("error in marshaling JSON:", err)
		return err
	}

	deployNFTUrl, err := url.JoinPath(nodeAddress, "/api/deploy-nft")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", deployNFTUrl, bytes.NewBuffer(bodyJSON))
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

	defer resp.Body.Close()

	_, err = utils.SignatureResponse(id, nodeAddress)
	return err
}

func (h *DoMintNFTApiCall) callback(
	caller *wasmtime.Caller,
	args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
	// Validate the number of arguments
	inputArgs, outputArgs := utils.HostFunctionParamExtraction(args, true, true)

	// Extract input bytes
	inputBytes, memory, err := utils.ExtractDataFromWASM(caller, inputArgs)
	if err != nil {
		fmt.Println("Failed to extract data from WASM", err)
		return utils.HandleError(err.Error())
	}
	h.memory = memory // Assign memory to Host struct for future use

	var mintNFTData MintNFTData
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(inputBytes, &mintNFTData)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in callback function:", err3)
		errMsg := "Error unmarshaling response in callback function:" + err3.Error()
		return utils.HandleError(errMsg)
	}

	callCreateNFTAPIResp := callCreateNFTAPI(h.nodeAddress, mintNFTData)
	var unmarshaledResponse map[string]interface{}
	err = json.Unmarshal(callCreateNFTAPIResp, &unmarshaledResponse)
	if err != nil {
		fmt.Println("Error in unmarshaling callCreateNFTAPIResp:", err)
		errMsg := "Error in unmarshalling reponse from creat nft api" + err.Error()
		return utils.HandleError(errMsg)
	}
	nftID := unmarshaledResponse["result"].(string)
	fmt.Println("Create NFT API result:", nftID)

	errDeploy := callDeployNFTAPI(h.nodeAddress, h.quorumType, mintNFTData, nftID)
	if errDeploy != nil {
		errMsg := "Deploy NFT API failed" + errDeploy.Error()
		return utils.HandleError(errMsg)
	}
	responseStr := string(callCreateNFTAPIResp)
	err = utils.UpdateDataToWASM(caller, h.allocFunc, responseStr, outputArgs)
	if err != nil {
		fmt.Println("Failed to update data to WASM", err)
		return utils.HandleError(err.Error())
	}

	return utils.HandleOk() // Success
}
