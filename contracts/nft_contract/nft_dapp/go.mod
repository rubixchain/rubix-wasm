module nft_dapp

go 1.22.4

require github.com/bytecodealliance/wasmtime-go v1.0.0 // indirect

replace github.com/rubixchain/rubix-wasm/go-wasm-bridge/wasmbridge => ../../../go-wasm-bridge

require github.com/rubixchain/rubix-wasm/go-wasm-bridge/wasmbridge v0.0.0-00010101000000-000000000000
