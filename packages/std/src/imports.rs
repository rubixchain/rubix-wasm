extern "C" {
    // do_api_call makes a request for an input url
    pub fn do_api_call(
        url_ptr: *const u8,
        url_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
    // do_mint_nft mints an NFT
    pub fn do_mint_nft(
        inputdata_ptr: *const u8,
        inputdata_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
    pub fn do_transfer_nft(
        inputdata_ptr: *const u8,
        inputdata_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
    //do_mint_ft mints fungible tokens
    pub fn do_mint_ft(
        inputdata_ptr: *const u8,
        inputdata_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
    pub fn do_transfer_ft(
        inputdata_ptr: *const u8,
        inputdata_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
    //lock RBT function
    pub fn do_lock_rbt(
        inputdata_ptr: *const u8,
        inputdata_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
}
