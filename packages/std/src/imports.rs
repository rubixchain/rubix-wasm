extern "C" {
    // do_api_call makes a request for an input url
    pub fn do_api_call(
        url_ptr: *const u8,
        url_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
}