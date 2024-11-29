module ft_dapp

go 1.22.4

require github.com/rubixchain/rubix-wasm/go-wasm-bridge v0.0.0-20241118120653-62085b862d17

require github.com/bytecodealliance/wasmtime-go v1.0.0 // indirect

replace github.com/rubixchain/rubix-wasm/go-wasm-bridge/wasmbridge => ../../../go-wasm-bridge

replace github.com/rubixchain/rubix-wasm/go-wasm-bridge/utils => ../../../go-wasm-bridge
