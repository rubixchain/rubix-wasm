use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;
use rubixwasm_std::call_do_api_call;

#[derive(Serialize, Deserialize)]
pub struct AddThreeNumsReq {
    pub a: u32,
    pub b: u32,
    pub c: u32,
}

#[derive(Serialize, Deserialize)]
pub struct TestVecReq {
    pub name_list: Vec<String>
}

#[derive(Serialize, Deserialize)]
pub struct GetSomeRespReq {
    pub url: String
}

#[derive(Serialize, Deserialize)]
pub struct GreetingsReq {
    pub name: String
}

#[derive(Serialize, Deserialize)]
pub struct GreetingsRes {
    pub result: String
}

#[contract_fn]
pub fn add_three_nums(input: AddThreeNumsReq) -> Result<String, WasmError> {
    if input.b == 0 {
        return Err(WasmError::from("Parameter 'b' cannot be zero"))
    }

    let sum = input.a + input.b + input.c;
    Ok(sum.to_string())
}

#[contract_fn]
pub fn test_vec(input: TestVecReq) -> Result<String, WasmError> {
    if input.name_list.len() == 0 {
        return Err(WasmError::from("name_list cannot be empty"))
    }

    let names = input.name_list;

    Ok(names.join("-"))
}

#[contract_fn]
pub fn make_some_api_call(inp: GetSomeRespReq) -> Result<String, WasmError> {
    let url = inp.url;
    
    match call_do_api_call(&url) {
        Ok(resp) => {
            Ok(resp)
        },
        Err(e) => {
            Err(e)
        }
    }
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
