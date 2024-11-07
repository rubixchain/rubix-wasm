extern "C" {
    // do_api_call makes a request for an input url
    pub fn do_api_call(
        url_ptr: *const u8,
        url_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
    // do_createnft_api_callmakes a request for an input data
    pub fn do_mint_nft(
        inputdata_ptr: *const u8,
        inputdata_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;

}