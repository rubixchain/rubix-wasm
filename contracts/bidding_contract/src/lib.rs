use imports::{call_get_bid_from_state, call_save_bid_to_state};
use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;
extern crate serde_json;
use crate::imports::call_ecies_decryption;
use crate::imports::DecryptionInputData;
use hex::decode;
use std::collections::HashMap;

#[derive(Deserialize, Serialize)]
pub struct PlaceBidReq {
    pub bidder_did: String,
    pub encrypted_bid_amount: String,
}

#[contract_fn]
pub fn place_bid(place_bid_req: PlaceBidReq) -> Result<String, WasmError> {

    let input_encrypted_bid = place_bid_req.encrypted_bid_amount;
    let private_key_path = "/home/rubix/Sai-Rubix/rubix-wasm/contracts/bidding_contract/bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a/pvtKey.pem";
    
    //decyption_data is the data which is required to pass as an input arguement to the decryption function
    let decryption_data = DecryptionInputData {
        Privatekey_path: private_key_path.to_string(),
        data: decode(input_encrypted_bid).map_err(|_| WasmError::from("Decoding failed"))?,
    };

    // Use `?` to propagate errors if any
    let decrypted_bid = call_ecies_decryption(decryption_data)
    .map_err(|_| WasmError::from("Decryption failed"))?;

    let bid_map: HashMap<String, f64> = serde_json::from_str(&decrypted_bid)
        .map_err(|_| WasmError::from("Failed to deserialize decrypted bid into HashMap"))?;

   
    // Extract the "bid_amount" key from the map
    let input_bid = bid_map.get("bid_amount")
        .ok_or_else(|| WasmError::from("Key 'bid_amount' not found in decrypted bid"))?;


    // Input Bid cannot be less that 3 RBT
    if *input_bid < 3.0_f64 {
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
        match call_save_bid_to_state(*input_bid) {
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
    use serde::{Deserialize,Serialize};
    use std::slice;
    use std::str;
    use std::string;
    
    #[derive(Serialize, Deserialize)]
    pub struct DecryptionInputData {
        pub Privatekey_path:String,
        pub data:Vec<u8>
    }

    extern "C" {
        pub fn get_bid_from_state(
            out_bid_ptr: *mut *const u8,
            out_bid_len: *mut usize
        ) -> i32;

        pub fn save_bid_to_state(
            in_bid_ptr: *const u8,
            in_bid_len: usize,
        ) -> i32;
        pub fn ecies_decryption(
            input_ptr: *const u8,
            input_len: usize,
            output_ptr: *mut *const u8,
            output_len: *mut usize
        )->i32;
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
            // Convert the bid_amount struct to bytes
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
    pub fn call_ecies_decryption(input_to_decryption:DecryptionInputData) -> Result<String, WasmError> {
        unsafe {
        let input_to_decryption_bytes = serde_json::to_string(&input_to_decryption).unwrap().into_bytes();
        let  input_ptr = input_to_decryption_bytes.as_ptr();
        let  input_len = input_to_decryption_bytes.len();
        // Allocate space for the decrypted_data pointer and length
        let mut decrypted_data_ptr: *const u8 = std::ptr::null();
        let mut decrypted_data_len: usize = 0;
    
  
        let result = ecies_decryption(input_ptr, input_len, &mut decrypted_data_ptr,&mut decrypted_data_len,);
           
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if decrypted_data_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(decrypted_data_ptr, decrypted_data_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }


        }


    }



}