module addition-contract-dapp

go 1.21

replace github.com/rubixchain/rubix-wasm/go-wasm-bridge/wasmbridge => ../../../go-wasm-bridge

require (
	github.com/rubixchain/rubix-wasm/go-wasm-bridge v0.1.2
	gorm.io/gorm v1.25.7-0.20240204074919-46816ad31dde
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
)

require (
	github.com/bytecodealliance/wasmtime-go v1.0.0
	gorm.io/driver/sqlite v1.5.7
)
