module network_coin_contract/dapp

go 1.23.2

require github.com/rubixchain/rubix-wasm/go-wasm-bridge v0.1.4

require (
	github.com/bytecodealliance/wasmtime-go v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
)

// Use local version of go-wasm-bridge that includes the RBT package
replace github.com/rubixchain/rubix-wasm/go-wasm-bridge => ../../../go-wasm-bridge
