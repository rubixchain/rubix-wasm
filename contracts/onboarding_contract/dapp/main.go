package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/bytecodealliance/wasmtime-go"
    "github.com/rubixchain/rubix-wasm/go-wasm-bridge/host"
    "github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils"
)

// Database instance
var db *gorm.DB

type OnboardingProvider struct {
    ProviderDID           string `json:"provider_did" gorm:"not null"`
    ReceiverDID           string `json:"receiver_did" gorm:"not null"`
    ProviderName          string `json:"provider_name" gorm:"not null"`
    InfrastructureDetails string `json:"infrastructure_details"`
    Signature             string `json:"signature" gorm:"not null"`
}

// Initialize Database
func initDB() {
    var err error
    db, err = gorm.Open(sqlite.Open("onboarding_contract.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err) 

    db.AutoMigrate(&OnboardingProvider{})
    }
}

// Onboarding Provider Function
func onboardingProvider(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var request OnboardingProvider
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    if request.ProviderDID == "" || request.ProviderName == "" {
        http.Error(w, "Provider DID and name cannot be empty", http.StatusBadRequest)
        return
    }

    if request.ReceiverDID == "" {
        http.Error(w, "Receiver DID cannot be empty", http.StatusBadRequest)
        return
    }

    // Save to database
    if err := db.Create(&request).Error; err != nil {
        http.Error(w, "Failed to onboard provider", http.StatusInternalServerError)
        return
    }

    response := fmt.Sprintf("Provider onboarded: %s", request.ProviderDID)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))
}

type DoDbWrite struct {
    allocFunc *wasmtime.Func
    memory    *wasmtime.Memory
}

func NewDoDbWrite() *DoDbWrite {
    return &DoDbWrite{}
}

func (h *DoDbWrite) Name() string {
    return "do_DB_write"
}

func (h *DoDbWrite) FuncType() *wasmtime.FuncType {
    return wasmtime.NewFuncType(
        []*wasmtime.ValType{
            wasmtime.NewValType(wasmtime.KindI32), // resp_ptr_ptr
            wasmtime.NewValType(wasmtime.KindI32), // resp_ptr_len
        },
        []*wasmtime.ValType{wasmtime.NewValType(wasmtime.KindI32)}, // return i32
    )
}

func (h *DoDbWrite) Initialize(allocFunc, deallocFunc *wasmtime.Func, memory *wasmtime.Memory, nodeAddress string, quorumType int) {
    h.allocFunc = allocFunc
    h.memory = memory
}

func (h *DoDbWrite) Callback() host.HostFunctionCallBack {
    return h.callback
}

func (h *DoDbWrite) callback(
    caller *wasmtime.Caller,
    args []wasmtime.Val,
) ([]wasmtime.Val, *wasmtime.Trap) {
    // Validate number of arguments
    if len(args) < 2 {
        return nil, wasmtime.NewTrap("Invalid arguments: Expected 2")
    }

    respPtrPtr := args[0].I32()
    //respLenPtr := args[1].I32()

    responseStr := "Database Write Successful"

    responseArg := &utils.WasmArgInfo{
        DataPtr:     respPtrPtr,
        DataPtrSize: int32(len(responseStr)),
    }

    err := utils.UpdateDataToWASM(caller, h.allocFunc, responseStr, responseArg)
    if err != nil {
        fmt.Println("Failed to update data to WASM", err)
        return utils.HandleError(err.Error())
    }

    return utils.HandleOk()
}

// Main Function
func main() {
    initDB()

    http.HandleFunc("/api/onboard-provider", onboardingProvider)

    fmt.Println("Server running on port 20000...")
    log.Fatal(http.ListenAndServe(":20000", nil))
}