use super::imports::do_api_call;
use super::imports::do_mint_nft;
use super::imports::do_transfer_nft;
use super::imports::do_mint_ft;
use super::imports::do_transfer_ft;
use std::slice;
use std::str;
use super::errors::WasmError;
use serde::{Serialize,Deserialize};
use serde_json;
use crate::imports::do_lock_rbt;

#[derive(Serialize, Deserialize)]
pub struct CreateNft {
    pub did:        String, 
    pub metadata:    String,
    pub artifact:    String,
}

#[derive(Serialize, Deserialize)]
pub struct TransferNft{
    pub comment:    String, 
    pub nft:        String,
    pub nft_data:   String,
    pub nft_value:  f64,
    pub owner:      String,
    pub receiver:    String,
}

#[derive(Serialize, Deserialize)]
pub struct TransferFt{
    pub comment:    String, 
    pub ft_count:        i32,
    pub ft_name:   String,
    pub creator_did:      String,
    pub sender: String,
    pub receiver:    String,
}

#[derive(Serialize, Deserialize)]
pub struct MintFt {
    pub did:        String, 
    pub ft_count:    i32,
    pub ft_name:    String,
    pub token_count: i32,
}

#[derive(Serialize, Deserialize)]
pub struct LockRBT {
    pub did: String,
    pub amount: u64,
    pub token_address: String,
}


// call_do_api_call is helper function for do_api_call import function 
pub fn call_do_api_call(url: &str) -> Result<String, WasmError> {
    unsafe {
        // Convert the URL string to bytes
        let url_bytes = url.as_bytes();
        let url_ptr = url_bytes.as_ptr();
        let url_len = url_bytes.len();

        // Allocate space for the response pointer and length
        let mut resp_ptr: *const u8 = std::ptr::null();
        let mut resp_len: usize = 0;

        // Call the imported host function
        let result = do_api_call(
            url_ptr,
            url_len,
            &mut resp_ptr,
            &mut resp_len,
        );
        
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if resp_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(resp_ptr, resp_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }
    }
}
// call_mint_nft_api is helper function for do_mint_nft import function 
pub fn call_mint_nft_api(input_data: CreateNft ) -> Result<String, WasmError> {
    unsafe {
        // Convert the input data to bytes
        let input_bytes = serde_json::to_string(&input_data).unwrap().into_bytes();

        // let input_bytes = input_data.as_bytes();
        let input_ptr = input_bytes.as_ptr();
        let input_len = input_bytes.len();

        // Allocate space for the response pointer and length
        let mut resp_ptr: *const u8 = std::ptr::null();
        let mut resp_len: usize = 0;

        // Call the imported host functionrubixwasm_std::
        let result = do_mint_nft(
            input_ptr,
            input_len,
            &mut resp_ptr,
            &mut resp_len,
        );
        
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if resp_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(resp_ptr, resp_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }
    }
}
pub fn call_transfer_nft_api(input_data: TransferNft) -> Result<String, WasmError> {
    unsafe {
        // Convert the input data to bytes
        let input_bytes = serde_json::to_string(&input_data).unwrap().into_bytes();

        // let input_bytes = input_data.as_bytes();
        let input_ptr = input_bytes.as_ptr();
        let input_len = input_bytes.len();

        // Allocate space for the response pointer and length
        let mut resp_ptr: *const u8 = std::ptr::null();
        let mut resp_len: usize = 0;

        // Call the imported host functionrubixwasm_std::
        let result = do_transfer_nft(
            input_ptr,
            input_len,
            &mut resp_ptr,
            &mut resp_len,
        );
        
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if resp_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(resp_ptr, resp_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }
    }

}

// call_mint_ft_api is helper function for do_mint_ft import function 
pub fn call_mint_ft_api(input_data: MintFt ) -> Result<String, WasmError> {
    unsafe {
        // Convert the input data to bytes
        let input_bytes = serde_json::to_string(&input_data).unwrap().into_bytes();

        // let input_bytes = input_data.as_bytes();
        let input_ptr = input_bytes.as_ptr();
        let input_len = input_bytes.len();

        // Allocate space for the response pointer and length
        let mut resp_ptr: *const u8 = std::ptr::null();
        let mut resp_len: usize = 0;

        // Call the imported host functionrubixwasm_std::
        let result = do_mint_ft(
            input_ptr,
            input_len,
            &mut resp_ptr,
            &mut resp_len,
        );
        
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if resp_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(resp_ptr, resp_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }
    }
}

pub fn call_transfer_ft_api(input_data: TransferFt) -> Result<String, WasmError> {
    unsafe {
        // Convert the input data to bytes
        let input_bytes = serde_json::to_string(&input_data).unwrap().into_bytes();

        // let input_bytes = input_data.as_bytes();
        let input_ptr = input_bytes.as_ptr();
        let input_len = input_bytes.len();

        // Allocate space for the response pointer and length
        let mut resp_ptr: *const u8 = std::ptr::null();
        let mut resp_len: usize = 0;

        // Call the imported host functionrubixwasm_std::
        let result = do_transfer_ft(
            input_ptr,
            input_len,
            &mut resp_ptr,
            &mut resp_len,
        );
        
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if resp_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(resp_ptr, resp_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }
    }
}

pub fn call_lock_rbt_api(input_data: LockRBT) -> Result<String, WasmError> {
    unsafe {
        // Convert the input data to bytes
        let input_bytes = serde_json::to_string(&input_data).unwrap().into_bytes();
        let input_ptr = input_bytes.as_ptr();
        let input_len = input_bytes.len();

        // Allocate space for the response pointer and length
        let mut resp_ptr: *const u8 = std::ptr::null();
        let mut resp_len: usize = 0;

        // Call the imported host function
        let result = do_lock_rbt(
            input_ptr,
            input_len,
            &mut resp_ptr,
            &mut resp_len,
        );

        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if resp_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(resp_ptr, resp_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }
    }
}

// Add your own implementation of generate_unique_id
pub fn generate_unique_id() -> String {
    // Counter-based approach for blockchain environment
    use std::sync::atomic::{AtomicU64, Ordering};
    static COUNTER: AtomicU64 = AtomicU64::new(1);
    let id = COUNTER.fetch_add(1, Ordering::SeqCst);
    format!("nct_{:x}", id)
}