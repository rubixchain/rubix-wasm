use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;
// use rubixwasm_std::call_do_api_call;
use rubixwasm_std::{call_mint_nft_api, call_transfer_nft_api};
// use rubixwasm_std::CreateNft;
use rubixwasm_std::helpers::{CreateNft, TransferNft};

pub const WHITELIST: &[&str] = &["rubix1", "rubix2"];

#[derive(Serialize, Deserialize)]
pub struct MintSampleNFTReq {
    pub name: String,
    pub nft_info: CreateNft
}

#[contract_fn]
pub fn mint_sample_nft(mint_sample_nft_req: MintSampleNFTReq)-> Result<String, WasmError>{
    // Minting is allowed only for those whose names are whitelisted
    let input_name = mint_sample_nft_req.name;

    if !WHITELIST.contains(&input_name.as_str()) {
        return Err(WasmError::from(format!("name {} is not allowed to mint the sample NFT", &input_name)));
    }

    let nft_mint_info = mint_sample_nft_req.nft_info;

    match call_mint_nft_api(nft_mint_info){
        Ok(resp) => {
            Ok(resp)
        },
        Err(e) => {
            Err(e)
        }
    }
}


#[derive(Serialize, Deserialize)]
pub struct TransferSampleNFTReq {
    pub name: String,
    pub nft_info: TransferNft
}

#[contract_fn]
pub fn transfer_sample_nft(transfer_sample_nft_req: TransferSampleNFTReq)-> Result<String, WasmError>{
    // Minting is allowed only for those whose names are whitelisted
    let input_name = transfer_sample_nft_req.name;

    if !WHITELIST.contains(&input_name.as_str()) {
        return Err(WasmError::from(format!("name {} is not allowed to transfer sample NFTs", &input_name)));
    }

    let nft_transfer_info = transfer_sample_nft_req.nft_info;

    match call_transfer_nft_api(nft_transfer_info){
        Ok(resp) => {
            Ok(resp)
        },
        Err(e) => {
            Err(e)
        }
    }
}
//pub fn transfer_sample_nft(){}

