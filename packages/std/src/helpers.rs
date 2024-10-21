use super::imports::do_api_call;
use std::slice;
use std::str;
use super::errors::WasmError;


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