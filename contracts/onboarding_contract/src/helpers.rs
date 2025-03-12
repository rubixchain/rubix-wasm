use std::ffi::CString;
use rubixwasm_std::errors::WasmError;
use super::imports::do_DB_write;
use super::imports::do_verify_platform_signature;
use std::slice;
use std::str;
use serde::{Serialize,Deserialize};
use serde_json;

pub fn call_DB_write() -> Result<String, WasmError> {
    unsafe {
        // Allocate space for the response pointer and length
        let mut db_query_ptr: *const u8 = std::ptr::null();
        let mut db_query_len: usize = 0;

        // Call the imported host function
        let result = do_DB_write(
            &mut db_query_ptr,
            &mut db_query_len,
        );

        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if db_query_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(db_query_ptr, db_query_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }
    }
}

pub fn call_do_verify_platform_signature(
    pub provider_did: String,
    pub receiver_did: String,
    pub signature: String,) -> Result<(), WasmError> {
    unsafe {
        // Convert the URL string to bytes
        let provider_did_ptr = provider_did_bytes.as_ptr();
        let provider_did_len = provider_did_bytes.len();

        let receiver_did_ptr = receiver_did_bytes.as_ptr();
        let receiver_did_len = receiver_did_bytes.len();

        let signature_ptr = signature_bytes.as_ptr();
        let signature_len = signature_bytes.len();

        // Call the imported host function
        let result = do_verify_platform_signature(
            provider_did_ptr,
            provider_did_len,
            receiver_did_ptr,
            receiver_did_len,
            signature_ptr,
            signature_len,
        );
        
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        Ok(())
    }
}

