use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;
// use rubixwasm_std::call_do_api_call;
use rubixwasm_std::call_mint_nft_api;
// use rubixwasm_std::CreateNft;
use rubixwasm_std::helpers::CreateNft;

#[derive(Serialize, Deserialize)]
pub struct AddThreeNumsReq {
    pub a: u32,
    pub b: u32,
    pub c: u32,
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
pub fn create_sample_nft(input_data: CreateNft)-> Result<String, WasmError>{
    Ok("6".to_string())

    // match call_mint_nft_api(input_data){
    //     Ok(resp) => {
    //         Ok(resp)
    //     },
    //     Err(e) => {
    //         Err(e)
    //     }
    // }

} 
//pub fn transfer_sample_nft(){}

