use imports::{call_get_bid_from_state, call_save_bid_to_state};
use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;

#[derive(Deserialize, Serialize)]
pub struct PlaceBidReq {
    pub bid_amount: f64,
}

#[contract_fn]
pub fn place_bid(place_bid_req: PlaceBidReq) -> Result<String, WasmError> {
    // Compare input bid
    let input_bid = place_bid_req.bid_amount;

    // Input Bid cannot be less that 3 RBT
    if input_bid < 3.0_f64 {
        return Err(WasmError::from("input bid must atleast be atleast"))
    }
    
    // Get current state
    let current_bid = match call_get_bid_from_state() {
        Ok(bid) => bid,
        Err(e) => {
            return Err(WasmError::from(format!("failed to fetch state from string, err: {}", e.msg)))
        }
    };
    let current_bid_f64: f64 = current_bid.parse().expect("unable to parse current_bid to f64, check value for current_bid");

    if input_bid.gt(&current_bid_f64) {
        match call_save_bid_to_state(input_bid) {
            Ok(()) => {}, 
            Err(e) => return Err(WasmError::from(format!("unable to save state: {}", e.msg))),
        };
    } else {
        return Ok("Input bid is smaller than current bid, current bid state cannot be updated".to_string())
    }

    Ok("Input bid is greater than current bid, state has been updated sucessfully".to_string())
}

// Implementation of the import functions are present in `../dapp/host` directory
mod imports {
    use rubixwasm_std::errors::WasmError;
    use std::slice;
    use std::str;

    extern "C" {
        pub fn get_bid_from_state(
            out_bid_ptr: *mut *const u8,
            out_bid_len: *mut usize
        ) -> i32;

        pub fn save_bid_to_state(
            in_bid_ptr: *const u8,
            in_bid_len: usize,
        ) -> i32;
    }

    pub fn call_get_bid_from_state() -> Result<String, WasmError> {
        unsafe {
            // Allocate space for the response pointer and length
            let mut highest_bid_ptr: *const u8 = std::ptr::null();
            let mut highest_bid_len: usize = 0;

            // Call the imported host function
            let result = get_bid_from_state(
                &mut highest_bid_ptr,
                &mut highest_bid_len,
            );
            
            if result != 0 {
                return Err(WasmError::from(format!("Host function returned error code {}", result)));
            }

            // Ensure the response pointer is not null
            if highest_bid_ptr.is_null() {
                return Err(WasmError::from("Response pointer is null".to_string()));
            }

            // Convert the response back to a Rust String
            let response_slice = slice::from_raw_parts(highest_bid_ptr, highest_bid_len);
            match str::from_utf8(response_slice) {
                Ok(s) => Ok(s.to_string()),
                Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
            }
        }
    }


    pub fn call_save_bid_to_state(bid_amount: f64) -> Result<(), WasmError> {
        unsafe {
            // Convert the URL string to bytes
            let bid_amount_bytes = bid_amount.to_string();
            let bid_amount_ptr = bid_amount_bytes.as_ptr();
            let bid_amount_len = bid_amount_bytes.len();


            // Call the imported host function
            let result = save_bid_to_state(
                bid_amount_ptr,
                bid_amount_len,
            );
            
            if result != 0 {
                return Err(WasmError::from(format!("Host function returned error code {}", result)));
            }

            Ok(())
        }
    }
}