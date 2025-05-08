use rubixwasm_std::errors::WasmError;
use rubixwasm_std::{ contract_fn, call_lock_rbt_api };
use rubixwasm_std::helpers:: {LockRBT, generate_unique_id};
use serde::{Deserialize, Serialize};
use lazy_static::lazy_static;
use std::sync::Mutex;

lazy_static! {
    static ref NETWORK_COIN_REGISTRY: Mutex<NetworkCoinRegistry> = Mutex::new(NetworkCoinRegistry::default());
}

// Request structure for network coin creation
#[derive(Serialize, Deserialize)]
pub struct NetworkCoinMintRequest {
    pub did: String,           // Organization DID
    pub token_name: String,    // Network Coin Name
    pub symbol: String,        // Symbol
    pub total_supply: u64,     // Total Supply
    pub rbt_to_lock: u64,      // Amount of RBTs to lock
}

// Response structure for network coin creation
#[derive(Serialize, Deserialize)]
pub struct NetworkCoinMintResponse {
    pub token_address: String,  // Contract address of the new token
    pub token_name: String,     // Network Coin Name
    pub symbol: String,         // Symbol
    pub total_supply: u64,      // Total Supply
    pub locked_rbt: u64,        // Amount of RBTs locked
}

// Structure for token transfer operations
#[derive(Serialize, Deserialize)]
pub struct NetworkCoinTransferRequest {
    pub token_address: String,  // Address of the token to transfer
    pub sender: String,         // Sender's DID
    pub receiver: String,       // Receiver's DID
    pub amount: u64,            // Amount to transfer
    pub comment: String,        // Optional transfer comment
}

// Structure for token balance queries
#[derive(Serialize, Deserialize)]
pub struct NetworkCoinBalanceRequest {
    pub token_address: String,  // Address of the token
    pub owner_did: String,      // DID of the token owner
}

// Structure for token balance responses
#[derive(Serialize, Deserialize)]
pub struct NetworkCoinBalanceResponse {
    pub token_address: String,  // Address of the token
    pub token_name: String,     // Name of the token
    pub symbol: String,         // Symbol of the token
    pub balance: u64,           // Token balance
}

// Internal structure for tracking network coins
#[derive(Serialize, Deserialize, Clone)]
pub struct NetworkCoinInfo {
    pub token_address: String,  // Contract address of the token
    pub creator_did: String,    // DID of the creator organization
    pub token_name: String,     // Network Coin Name
    pub symbol: String,         // Symbol
    pub total_supply: u64,      // Total Supply
    pub locked_rbt: u64,        // Amount of RBTs locked
    pub creation_timestamp: u64, // Creation timestamp
}

// Token registry for keeping track of all network coins
#[derive(Serialize, Deserialize, Default)]
pub struct NetworkCoinRegistry {
    pub coins: Vec<NetworkCoinInfo>,  // List of all network coins
}

// RBT locking request
#[derive(Serialize, Deserialize)]
pub struct LockRBTRequest {
    pub did: String,            // DID of the organization
    pub amount: u64,            // Amount of RBTs to lock
    pub token_address: String,  // Address of the network coin
}

// RBT locking response
#[derive(Serialize, Deserialize)]
pub struct LockRBTResponse {
    pub transaction_id: String, // ID of the locking transaction
    pub status: bool,           // Success status
    pub message: String,        // Status message
}

// Lock record stored on-chain
pub struct RBTLockRecord {
    pub lock_id: String,         // Unique ID for the lock
    pub owner_did: String,       // Organization DID
    pub network_coin: String,    // Address of associated network coin
    pub amount: u64,             // Amount of RBTs locked
    pub timestamp: u64,          // Lock timestamp
    pub status: LockStatus,      // Current status of the lock
}

// Status of a lock
pub enum LockStatus {
    Active,                     // Lock is currently active
    Released,                   // Lock has been released
    Partial(u64),               // Partially released with amount remaining
}

// Structure for transaction history requests
#[derive(Serialize, Deserialize)]
pub struct TransactionHistoryRequest {
    pub token_address: String,    // Address of the token
    pub owner_did: Option<String>, // Optional: filter by owner DID
    pub limit: Option<u32>,       // Optional: limit the number of results
    pub offset: Option<u32>,      // Optional: pagination offset
}

// Structure for transaction records
#[derive(Serialize, Deserialize, Clone)]
pub struct TransactionRecord {
    pub transaction_id: String,    // ID of the transaction
    pub token_address: String,     // Address of the token
    pub transaction_type: String,  // Type of transaction (mint, transfer, etc.)
    pub from_did: String,          // Sender DID
    pub to_did: String,            // Receiver DID
    pub amount: u64,               // Transaction amount
    pub timestamp: u64,            // Transaction timestamp
    pub comment: Option<String>,   // Optional comment
}

// Structure for transaction history response
#[derive(Serialize, Deserialize)]
pub struct TransactionHistoryResponse {
    pub transactions: Vec<TransactionRecord>, // List of transactions
    pub total_count: u32,                    // Total number of transactions
    pub limit: u32,                          // Limit used for this query
    pub offset: u32,                         // Offset used for this query
}

// Mock function to simulate getting a balance
fn get_mock_balance(token_address: &str, owner_did: &str) -> u64 {
    // This is just a placeholder
    if let Some(coin_info) = get_coin_from_registry(token_address) {
        if coin_info.creator_did == owner_did {
            return coin_info.total_supply;
        }
    }
    0
}

// Mock function to simulate getting transaction history
fn get_mock_transactions(request: &TransactionHistoryRequest) -> Vec<TransactionRecord> {
    // Just a placeholder
    let mut transactions = Vec::new();

    if let Some(coin_info) = get_coin_from_registry(&request.token_address) {
        // Add a mock mint transaction
        let mint_transaction = TransactionRecord {
            transaction_id: format!("tx_{}", generate_unique_id()),
            token_address: request.token_address.clone(),
            transaction_type: "mint".to_string(),
            from_did: "system".to_string(),
            to_did: coin_info.creator_did.clone(),
            amount: coin_info.total_supply,
            timestamp: coin_info.creation_timestamp,
            comment: None,
        };

        transactions.push(mint_transaction);
    }

    transactions
}

// Function to add a coin to the registry
fn add_coin_to_registry(coin_info: NetworkCoinInfo) -> Result<(), WasmError> {
    let mut registry = NETWORK_COIN_REGISTRY.lock().map_err(|e|
        WasmError::from(format!("Failed to lock registry: {}", e)))?;

    // Check if a coin with the same address already exists
    if registry.coins.iter().any(|c| c.token_address == coin_info.token_address) {
        return Err(WasmError::from("A coin with this address already exists"));
    }

    registry.coins.push(coin_info);
    Ok(())
}

// Function to look up a coin in the registry
fn get_coin_from_registry(token_address: &str) -> Option<NetworkCoinInfo> {
    let registry = NETWORK_COIN_REGISTRY.lock().ok()?;
    registry.coins.iter()
        .find(|c| c.token_address == token_address)
        .cloned()
}

// Function to update the registry after creating a network coin
fn register_network_coin(
    token_address: String,
    creator_did: String,
    token_name: String,
    symbol: String,
    total_supply: u64,
    locked_rbt: u64
) -> Result<(), WasmError> {
    use std::time::{SystemTime, UNIX_EPOCH};
    let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs();

    let coin_info = NetworkCoinInfo {
        token_address,
        creator_did,
        token_name,
        symbol,
        total_supply,
        locked_rbt,
        creation_timestamp: now,
    };

    add_coin_to_registry(coin_info)
}

#[contract_fn]
pub fn lock_rbt(request: LockRBTRequest) -> Result<String, WasmError> {
    // Validate input parameters
    if request.amount == 0 {
        return Err(WasmError::from("RBT amount must be greater than zero"));
    }

    // Create the lock RBT request
    let lock_request = LockRBT {
        did: request.did.clone(),
        amount: request.amount,
        token_address: request.token_address.clone(),
    };

    // Call the API to lock RBTs
    match call_lock_rbt_api(lock_request) {
        Ok(response) => {
            // Parse the response
            match serde_json::from_str::<LockRBTResponse>(&response) {
                Ok(lock_response) => {
                    if lock_response.status {
                        // Successful lock
                        Ok(serde_json::to_string(&lock_response).unwrap())
                    } else {
                        // Lock failed with a specific message
                        Err(WasmError::from(format!("RBT locking failed: {}", lock_response.message)))
                    }
                },
                Err(e) => Err(WasmError::from(format!("Failed to parse lock response: {}", e))),
            }
        },
        Err(e) => Err(e),
    }
}

#[contract_fn]
pub fn create_network_coin(request: NetworkCoinMintRequest) -> Result<String, WasmError> {
    // 1. Validate input parameters
    if request.token_name.is_empty() {
        return Err(WasmError::from("Token name cannot be empty"));
    }

    if request.symbol.is_empty() {
        return Err(WasmError::from("Token symbol cannot be empty"));
    }

    if request.total_supply == 0 {
        return Err(WasmError::from("Total supply must be greater than zero"));
    }

    if request.rbt_to_lock == 0 {
        return Err(WasmError::from("RBT to lock must be greater than zero"));
    }

    // 2. Create a placeholder for the token address
    // In a real implementation, this would be generated by the blockchain
    let token_address = format!("nct_{}", generate_unique_id());

    // 3. Lock the RBTs
    let lock_request = LockRBTRequest {
        did: request.did.clone(),
        amount: request.rbt_to_lock,
        token_address: token_address.clone(),
    };

    match lock_rbt(lock_request) {
        Ok(lock_response_str) => {
            // Parse the lock response
            let lock_response: LockRBTResponse = serde_json::from_str(&lock_response_str)
                .map_err(|e| WasmError::from(format!("Failed to parse lock response: {}", e)))?;

            if !lock_response.status {
                return Err(WasmError::from(format!("RBT locking failed: {}", lock_response.message)));
            }

            // 4. Mint the network coin
            let mint_request = rubixwasm_std::helpers::MintFt {
                did: request.did.clone(),
                ft_count: request.total_supply as i32,
                ft_name: request.token_name.clone(),
                token_count: 1, // Always mint one token type
            };

            match rubixwasm_std::call_mint_ft_api(mint_request) {
                Ok(_mint_response) => {
                    register_network_coin(
                        token_address.clone(),
                        request.did.clone(),
                        request.token_name.clone(),
                        request.symbol.clone(),
                        request.total_supply,
                        request.rbt_to_lock
                    )?;

                    let network_coin_response = NetworkCoinMintResponse {
                        token_address: token_address,
                        token_name: request.token_name,
                        symbol: request.symbol,
                        total_supply: request.total_supply,
                        locked_rbt: request.rbt_to_lock,
                    };

                    Ok(serde_json::to_string(&network_coin_response).unwrap())
                },
                Err(e) => {
                    Err(WasmError::from(format!("Failed to mint network coin: {}", e.msg)))
                }
            }
        },
        Err(e) => Err(e),
    }
}

#[contract_fn]
pub fn get_token_balance(request: NetworkCoinBalanceRequest) -> Result<String, WasmError> {
    // Validate the token exists in our registry
    if get_coin_from_registry(&request.token_address).is_none() {
        return Err(WasmError::from("Network coin not found"));
    }

    // In a real implementation, we would query the blockchain for the actual balance
    // For now, we'll return a mock balance
    let balance = get_mock_balance(&request.token_address, &request.owner_did);

    // Get coin info from registry
    let coin_info = get_coin_from_registry(&request.token_address)
        .ok_or_else(|| WasmError::from("Network coin not found"))?;

    // Create and return the balance response
    let balance_response = NetworkCoinBalanceResponse {
        token_address: request.token_address,
        token_name: coin_info.token_name.clone(),
        symbol: coin_info.symbol.clone(),
        balance,
    };

    Ok(serde_json::to_string(&balance_response).unwrap())
}

#[contract_fn]
pub fn transfer_network_coin(request: NetworkCoinTransferRequest) -> Result<String, WasmError> {
    // Validate the token exists in our registry
    if get_coin_from_registry(&request.token_address).is_none() {
        return Err(WasmError::from("Network coin not found"));
    }

    // Validate the sender has enough balance
    let sender_balance = get_mock_balance(&request.token_address, &request.sender);
    if sender_balance < request.amount {
        return Err(WasmError::from("Insufficient balance"));
    }

    // Create the transfer request
    let transfer_request = rubixwasm_std::helpers::TransferFt {
        comment: request.comment,
        ft_count: request.amount as i32,
        ft_name: get_coin_from_registry(&request.token_address).unwrap().token_name.clone(),
        creator_did: get_coin_from_registry(&request.token_address).unwrap().creator_did.clone(),
        sender: request.sender,
        receiver: request.receiver,
    };

    // Call the transfer API
    match rubixwasm_std::call_transfer_ft_api(transfer_request) {
        Ok(response) => {
            // Return the response
            Ok(response)
        },
        Err(e) => Err(e),
    }
}

#[contract_fn]
pub fn get_transaction_history(request: TransactionHistoryRequest) -> Result<String, WasmError> {
    // Validate the token exists in our registry
    if get_coin_from_registry(&request.token_address).is_none() {
        return Err(WasmError::from("Network coin not found"));
    }

    // In a real implementation, we would query the blockchain for transaction history
    // For now, we'll return mock data
    let transactions = get_mock_transactions(&request);

    // Create the response
    let response = TransactionHistoryResponse {
        transactions: transactions.clone(),
        total_count: transactions.len() as u32,
        limit: request.limit.unwrap_or(10),
        offset: request.offset.unwrap_or(0),
    };

    Ok(serde_json::to_string(&response).unwrap())
}

// Add these at the appropriate place in your lib.rs file
#[derive(Serialize, Deserialize)]
pub struct GetAllNetworkCoinsRequest {
    // Empty struct - no fields needed
}

#[contract_fn]
pub fn get_all_network_coins(_request: GetAllNetworkCoinsRequest) -> Result<String, WasmError> {
    // Get the registry
    let registry = NETWORK_COIN_REGISTRY.lock()
        .map_err(|e| WasmError::from(format!("Failed to lock registry: {}", e)))?;

    // Return serialized list of coins
    Ok(serde_json::to_string(&registry.coins).unwrap())
}