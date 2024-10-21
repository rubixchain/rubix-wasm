# Rubix Wasm 

A set of libraries enabling developers to build Contracts on Rubix with ease!

## Project structure

- `packages/` - Rust packages that provide abstractions to write Smart Contracts
    - `derive` - It contains `contract_fn` which needs to implemented on Smart Contract methods.
    - `std` - It provides all the imports, error handling and memory management features.

- `go-wasm-bridge` - Golang bindings which helps interacting with the WASM binary

