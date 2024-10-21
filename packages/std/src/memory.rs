use std::mem;

/// Allocates a block of memory of the specified size and returns a pointer to it.
/// Uses `Vec::with_capacity` and `mem::forget` to manage memory without Rust's automatic deallocation.
#[no_mangle]
pub extern "C" fn alloc(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    mem::forget(buf); // Prevent Rust from freeing the memory
    ptr
}

/// Deallocates previously allocated memory by reconstructing the `Vec` from the raw pointer,
/// allowing Rust to manage its deallocation.
#[no_mangle]
pub extern "C" fn dealloc(ptr: *mut u8, size: usize) {
    unsafe {
        let _ = Vec::from_raw_parts(ptr, size, size);
    }
}
