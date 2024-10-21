use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct WasmError {
    pub msg: String
}

impl From<String> for WasmError {
    fn from(msg: String) -> Self {
        WasmError { msg }
    }
}

impl From<&str> for WasmError {
    fn from(msg: &str) -> Self {
        WasmError {
            msg: msg.to_string()
        }
    }
}