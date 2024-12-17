use imports::call_get_bidding_info_from_state_file;
use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;
extern crate serde_json;
use crate::imports::{call_ecies_decryption,call_save_bidding_info};
use crate::imports::{DecryptionInputData, PlaceBidReq, RevealHighestBidReq};
use hex::decode;
use std::collections::HashMap;
use serde::de::Error as SerdeError;


#[derive(Deserialize, Serialize, Debug)]
struct BiddingInfo {
    bidding_info: HashMap<String, String>,
}

#[contract_fn]
pub fn reveal_highest_bid(reveal_highest_bid_req: RevealHighestBidReq) -> Result<String, WasmError> {
    // Retrieve bid info from state file
    let bid_info = match call_get_bidding_info_from_state_file() {
        Ok(bid_info) => bid_info, // Extract bid_info if successful
        Err(e) => return Err(WasmError::from(format!("Failed to get bidding info: {}", e.msg))), // Propagate the error
    };
    // println!("bid_info value is:{:?}",bid_info);
    // Deserialize bid_info into a HashMap
    let deserialized_bid_info: BiddingInfo = serde_json::from_str(&bid_info)
        .map_err(|_| WasmError::from("Failed to deserialize bid info into HashMap"))?;

       // Initialize variables to track the highest bid and bidder
    let mut highest_bid = 0f64; // Track the highest bid
    let mut highest_bidder = String::new(); // Track the highest bidder

    // Iterate over all encrypted bid values and decrypt them
    for (bidder, encrypted_bid) in deserialized_bid_info.bidding_info.iter() {
        // Decode the encrypted bid from hex
        let encrypted_bid_bytes = decode(encrypted_bid)
            .map_err(|_| WasmError::from("Failed to decode encrypted bid from hex"))?;

        // Prepare input for call_ecies_decryption
        let input_data = DecryptionInputData {
            Privatekey_path: reveal_highest_bid_req.deployer_password.clone(),
            data: encrypted_bid_bytes,
        };

        // Decrypt the bid value using ECIES
        let decrypted_bid = match call_ecies_decryption(input_data) {
            Ok(decrypted_bid) => decrypted_bid, // Successfully decrypted bid
            Err(e) => return Err(WasmError::from(format!("Failed to decrypt bid value: {}", e.msg))), // Propagate the error
        };


        // Convert the decrypted string into a f64 value (assuming bids are stored as integers)
        let bid_value: f64 = serde_json::from_str(&decrypted_bid)
        .and_then(|json: serde_json::Value| {
            json.get("bid_amount")
                .and_then(|v| v.as_f64())
                .ok_or_else(|| serde_json::Error::custom("Missing or invalid 'bid_amount' field"))
        })
        .map_err(|_| WasmError::from("Failed to parse decrypted bid as f64"))?;
        // Check if this is the highest bid
        if bid_value > highest_bid {
            highest_bid = bid_value;
            highest_bidder = bidder.clone(); // Update the highest bidder
        }
        
    }
    // Return the highest bid and the associated bidder
    Ok(format!("Highest bid: {}, Bidder: {}", highest_bid, highest_bidder))

}


#[contract_fn]
pub fn place_bid(place_bid_req: PlaceBidReq) -> Result<String, WasmError> {
    
    // Save whatever is coming directly to the

    match call_save_bidding_info(&place_bid_req) {
        Ok(_) => Ok(format!("Bid for did: {} has been saved", place_bid_req.bidder_did)), 
        Err(e) => return Err(WasmError::from(format!("unable to save bidding info: {}", e.msg))),
    }
    
}

// Implementation of the import functions are present in `../dapp/host` directory
mod imports {
    use rubixwasm_std::errors::WasmError;
    use serde::{Deserialize,Serialize};
    use std::slice;
    use std::str;
    
    #[derive(Serialize, Deserialize)]
    pub struct DecryptionInputData {
        pub Privatekey_path:String,
        pub data:Vec<u8>
    }

    #[derive(Deserialize, Serialize)]
    pub struct PlaceBidReq {
        pub bidder_did: String,
        pub encrypted_bid_amount: String,
    }

    #[derive(Deserialize, Serialize)]
    pub struct RevealHighestBidReq {
        pub deployer_password:String,
    }

    extern "C" {
        pub fn ecies_decryption(
            input_ptr: *const u8,
            input_len: usize,
            output_ptr: *mut *const u8,
            output_len: *mut usize
        )->i32;

        pub fn save_bidding_info(
            in_placebid_req_ptr: *const u8,
            in_placebid_req_bid_len: usize,
        ) -> i32;

        pub fn get_bidding_info_from_state_file(
            out_bid_ptr: *mut *const u8,
            out_bid_len: *mut usize
        ) -> i32;

    }
    pub fn call_get_bidding_info_from_state_file() -> Result<String, WasmError> {
        unsafe {
            // Allocate space for the response pointer and length
            let mut bid_info_ptr: *const u8 = std::ptr::null();
            let mut bid_info_len: usize = 0;

            // Call the imported host function
            let result = get_bidding_info_from_state_file(
                &mut bid_info_ptr,
                &mut bid_info_len,
            );
            if result != 0 {
                return Err(WasmError::from(format!("Host function returned error code {}", result)));
            }
             // Ensure the response pointer is not null
             if bid_info_ptr.is_null() {
                return Err(WasmError::from("Response pointer is null".to_string()));
            }
             // Convert the response back to a Rust String
             let response_slice = slice::from_raw_parts(bid_info_ptr, bid_info_len);
             match str::from_utf8(response_slice) {
                 Ok(s) => Ok(s.to_string()),
                 Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
             }

               
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
    
    pub fn call_save_bidding_info(place_bid_req: &PlaceBidReq) -> Result<(), WasmError> {
        unsafe {
            // Convert the place_bid_req struct to string
            let place_bid_req_str = serde_json::to_string(&place_bid_req).map_err(|_| WasmError::from("Failed to serialize struct place_bid_req"))?;
            let place_bid_req_ptr = place_bid_req_str.as_ptr();
            let place_bid_req_len = place_bid_req_str.len();


            // Call the imported host function
            let result = save_bidding_info(
                place_bid_req_ptr,
                place_bid_req_len,
            );
            
            if result != 0 {
                return Err(WasmError::from(format!("Host function returned error code {}", result)));
            }

            Ok(())
        }

    }
}

