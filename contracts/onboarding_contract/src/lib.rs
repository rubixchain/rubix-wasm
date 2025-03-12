use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;
use rubixwasm_std::call_do_api_call;
pub mod helpers;
pub use helpers::call_DB_write;
pub use helpers::call_do_verify_platform_signature;

#[derive(Serialize, Deserialize)]
pub struct OnboardingProvider {
    pub provider_did: String,
    pub receiver_did: String,
    //pub provider_name: String,
    pub infrastructure_details: String,
    pub signature: String,
}

#[derive(Serialize, Deserialize)]
pub struct ProviderInfo {
    pub provider_did: String,
    pub provider_name: String,
    pub cpu_provided: String,
    pub gpu_provided: String,
}

pub fn onboard_provider(provider: OnboardingProvider) {
    println!("Onboarding provider with DID: {}", provider.provider_did);
    println!("Receiver DID: {}", provider.receiver_did);
    println!("Infrastructure Details: {}", provider.infrastructure_details);
    println!("Signature: {}", provider.signature);
}

pub fn register_provider_info(info: ProviderInfo) {
    println!("Registering provider: {}", info.provider_name);
    println!("CPU Provided: {}", info.cpu_provided);
    println!("GPU Provided: {}", info.gpu_provided);
}

#[contract_fn]
pub fn onboarding_provider(input: OnboardingProvider) -> Result<String, WasmError> {
    if input.provider_did.is_empty() {
        return Err(WasmError::from("Platform DID cannot be empty"));
    }
    if input.receiver_did.is_empty() {
        return Err(WasmError::from("Receiver DID cannot be empty"));
    }
    call_DB_write(input: )
}

#[contract_fn]
pub fn get_provider_info(input: ProviderInfo) -> Result<, WasmError> {
    if input.provider_did.is_empty() || input.provider_name.is_empty() {
        return Err(WasmError::from("Platform DID and provider name cannot be empty"));
    }
    if input.cpu_provided.is_empty() || input.gpu_provided.is_empty() {
        return Err(WasmError::from("CPU and GPU details cannot be empty"));
    }
    call_DB_write(input: )
}


