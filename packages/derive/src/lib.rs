use proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, ItemFn};

#[proc_macro_attribute]
pub fn contract_fn(_attr: TokenStream, item: TokenStream) -> TokenStream {
    // Parse the input function
    let input = parse_macro_input!(item as ItemFn);

    // Extract the function name and input/output types
    let func_name = &input.sig.ident;
    let input_type = match input.sig.inputs.first() {
        Some(syn::FnArg::Typed(arg)) => &arg.ty,
        _ => panic!("Expected a function with one argument"),
    };
    let output_type = match &input.sig.output {
        syn::ReturnType::Type(_, ty) => ty,
        _ => panic!("Expected a function with a return type"),
    };

    // Generate a new name for the wrapper function
    let wrapper_func_name = syn::Ident::new(&format!("{}_", func_name), func_name.span());

    // Generate the wrapper function
    let expanded = quote! {
        // Original function
        #input

        // Generated wrapper function
        #[no_mangle]
        pub extern "C" fn #wrapper_func_name(
            input_ptr: *mut u8,
            input_len: usize,
            output_ptr_ptr: *mut *mut u8,
            output_len_ptr: *mut usize
        ) -> i32 {
            use std::slice;
            use std::ptr;
            use serde::{Serialize, Deserialize};
            use serde_json;

            // Safety: Ensure the pointers are valid
            unsafe {
                // Deserialize input data
                let input_data = slice::from_raw_parts(input_ptr, input_len);
                let input: #input_type = match serde_json::from_slice(input_data) {
                    Ok(data) => data,
                    Err(_) => return 1, // Error during deserialization
                };

                // Call the original function
                let result: #output_type = #func_name(input);

                match result {
                    Ok(success_str) => {
                        // Serialize the success string
                        let serialized_output = match serde_json::to_vec(&success_str) {
                            Ok(data) => data,
                            Err(_) => return 1, // Error during serialization
                        };

                        // Allocate memory for output data
                        let output_len = serialized_output.len();
                        let output_ptr = ::rubixwasm_std::memory::alloc(output_len);

                        // Write serialized data to output_ptr
                        ptr::copy_nonoverlapping(serialized_output.as_ptr(), output_ptr, output_len);
                        *output_ptr_ptr = output_ptr;
                        *output_len_ptr = output_len;

                        0 // Success
                    },
                    Err(wasm_error) => {
                        // Serialize the error message
                        let serialized_error = match serde_json::to_vec(&wasm_error.msg) {
                            Ok(data) => data,
                            Err(_) => return 1, // Error during serialization
                        };

                        // Allocate memory for error data
                        let error_len = serialized_error.len();
                        let error_ptr = ::rubixwasm_std::memory::alloc(error_len);

                        // Write serialized error to error_ptr
                        ptr::copy_nonoverlapping(serialized_error.as_ptr(), error_ptr, error_len);
                        *output_ptr_ptr = error_ptr;
                        *output_len_ptr = error_len;

                        1 // Error
                    }
                }
            }
        }
    };

    TokenStream::from(expanded)
}
