use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;
// use rubixwasm_std::call_do_api_call;
use rubixwasm_std::{call_mint_ft_api, call_transfer_ft_api};
// use rubixwasm_std::CreateNft;
use rubixwasm_std::helpers::{MintFt, TransferFt};

pub const WHITELIST: &[&str] = &["rubix1", "rubix2"];

#[derive(Serialize, Deserialize)]
pub struct MintSampleFTReq {
    pub name: String,
    pub ft_info: MintFt
}

#[contract_fn]
pub fn mint_sample_ft(mint_sample_ft_req: MintSampleFTReq)-> Result<String, WasmError>{
    // Minting is allowed only for those whose names are whitelisted
    let input_name = mint_sample_ft_req.name;

    if !WHITELIST.contains(&input_name.as_str()) {
        return Err(WasmError::from(format!("name {} is not allowed to mint the sample NFT", &input_name)));
    }

    let ft_mint_info = mint_sample_ft_req.ft_info;

    match call_mint_ft_api(ft_mint_info){
        Ok(resp) => {
            Ok(resp)
        },
        Err(e) => {
            Err(e)
        }
    }
}


#[derive(Serialize, Deserialize)]
pub struct TransferSampleFTReq {
    pub name: String,
    pub ft_info: TransferFt
}

#[contract_fn]
pub fn transfer_sample_ft(transfer_sample_ft_req: TransferSampleFTReq)-> Result<String, WasmError>{
    let input_name = transfer_sample_ft_req.name;

    if !WHITELIST.contains(&input_name.as_str()) {
        return Err(WasmError::from(format!("name {} is not allowed to transfer sample FTs", &input_name)));
    }

    let ft_transfer_info = transfer_sample_ft_req.ft_info;

    match call_transfer_ft_api(ft_transfer_info){
        Ok(resp) => {
            Ok(resp)
        },
        Err(e) => {
            Err(e)
        }
    }
}

