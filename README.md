# Rubix Wasm 

A set of libraries enabling developers to build WASM Contracts on Rubix with ease!

## Project structure

- `packages/` - Rust packages that provide abstractions to write Smart Contracts
    - `derive` - It contains `contract_fn` which needs to implemented on Smart Contract methods.
    - `std` - It provides all the imports, error handling and memory management features.

- `go-wasm-bridge` - Golang bindings which helps interacting with the WASM binary

## Contracts

- [Generic Contract]() - A simple contract with sample functions to demostrate `rubix-wasm` features
- [Bidding Contract](https://github.com/rubixchain/rubix-wasm/tree/main/contracts/bidding_contract) - A simple contract which takes a Bid amount and stores it when the provided Bid amount is larger than the current Bid amount. 

## Usage

The [generic contract](https://github.com/rubixchain/rubix-wasm/tree/main/contracts/generic_contract) is a good starting point to explore the usage of libraries and the paradigm to write contracts.

### 1. Write WASM Contract in Rust

Initially, we will provide support for Rust when it comes to writing contracts, as support for additional languages will be added in the future. Writing contract functions is as simple as writing any function in Rust. However, these functions won't be available for use outside of the WASM environment. To export these functions, you need to import and use the `contract_fn` macro, which is available through [rubixwasm-std](https://github.com/rubixchain/rubix-wasm/tree/main/packages/std).

Exported contract functions must adhere to a specific signature. They can only accept a single struct argument that implements [Serde](https://serde.rs/)'s `Serialize` and `Deserialize` traits and must return a `Result<String, WasmError>`. The expected return value must be serialized to a String. The `WasmError` is a string-like type responsible for handling errors and can be imported from [rubixwasm-std](https://github.com/rubixchain/rubix-wasm/tree/main/packages/std).

Following is an illustration of an exported contract function:

```rs
/// contracts/generic_contract/src/lib.rs

use rubixwasm_std::errors::WasmError;
use rubixwasm_std::contract_fn;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct GreetingsReq {
    pub name: String
}

#[derive(Serialize, Deserialize)]
pub struct GreetingsRes {
    pub result: String
}

#[contract_fn]
pub fn greetings(inp: GreetingsReq) -> Result<String, WasmError> {
    let input_name = inp.name;
    if input_name.clone().len() < 3 {
        return Err(WasmError::from("Your name must be alteast 3 characters long"))
    }

    let greeting_string = format!("Hello, {}", input_name);

    let result = GreetingsRes { result: greeting_string };
    let stringifyed_result = serde_json::to_string(&result).expect("unable to serialize struct for Contract function Greetings");

    Ok(stringifyed_result)
}
```

Build the project with target as `wasm32-unknown-unknown`. The WASM binary will be generated at `target/wasm32-unknown-unknown/debug`.

### 2. Execution of Contract

The [go-wasm-bridge](https://github.com/rubixchain/rubix-wasm/tree/main/go-wasm-bridge) will help us executing our WASM Contract binary. Refer [here](https://github.com/rubixchain/rubix-wasm/blob/main/contracts/generic_contract/dapp/main.go) for its usage.


The Go package provides us a `CallFunction()` which lets us call the exported Contract function. Refer [this](https://github.com/rubixchain/rubix-wasm/blob/1d6dc0b989a7d5278da6e71c8b29739fa6e6cdb3/contracts/generic_contract/dapp/main.go#L51) function which is responsible for calling the `greetings` contract function. [Here](https://github.com/rubixchain/rubix-wasm/blob/1d6dc0b989a7d5278da6e71c8b29739fa6e6cdb3/contracts/generic_contract/dapp/main.go#L52) we are letting the WASM runtime know about the function we are interested to execute. This input stringifyed JSON syntax for calling Contract functions is similar to EVM and CosmWasm contracts.




