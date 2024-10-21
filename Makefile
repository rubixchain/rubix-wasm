#!/usr/bin/make -f

# Directory containing all contract crates
CONTRACTS_DIR := contracts

# Directory to store artifacts
ARTIFACTS_DIR := artifacts

# Target architecture for WASM
TARGET := wasm32-unknown-unknown

# Cargo build flags for release build targeting WASM
BUILD_FLAGS := --release --target $(TARGET)

# Find all Cargo.toml files within immediate subdirectories of CONTRACTS_DIR
CONTRACTS := $(shell find $(CONTRACTS_DIR) -mindepth 2 -maxdepth 2 -name Cargo.toml | sed 's|/Cargo.toml||')

.PHONY: all build clean

# Default target: build all contracts and copy artifacts
all: build

# Build all contracts and copy their WASM binaries to artifacts directory
build:
	@echo "Starting build of all contracts in '$(CONTRACTS_DIR)'..."
	@mkdir -p $(ARTIFACTS_DIR)
	@for contract in $(CONTRACTS); do \
		contract_name=$$(basename $$contract); \
		echo "Building $$contract_name..."; \
		(cd $$contract && cargo build $(BUILD_FLAGS)) || { echo "Build failed for $$contract_name"; exit 1; }; \
		wasm_file="target/$(TARGET)/release/$$contract_name.wasm"; \
		src_path="$$contract/$$wasm_file"; \
		if [ -f $$src_path ]; then \
			cp $$src_path $(ARTIFACTS_DIR)/; \
			echo "Copied $$src_path to $(ARTIFACTS_DIR)/"; \
		else \
			echo "Warning: WASM file $$src_path not found."; \
		fi; \
	done
	@echo "All contracts built and WASM binaries copied to '$(ARTIFACTS_DIR)' successfully."

build-wasi:
	@echo "Starting build of all contracts in '$(CONTRACTS_DIR)'..."
	@mkdir -p $(ARTIFACTS_DIR)
	@for contract in $(CONTRACTS); do \
		contract_name=$$(basename $$contract); \
		echo "Building $$contract_name..."; \
		(cd $$contract && cargo build $(BUILD_FLAGS)) || { echo "Build failed for $$contract_name"; exit 1; }; \
		wasm_file="target/$(TARGET)/release/$$contract_name.wasm"; \
		src_path="$$contract/$$wasm_file"; \
		if [ -f $$src_path ]; then \
			cp $$src_path $(ARTIFACTS_DIR)/; \
			echo "Copied $$src_path to $(ARTIFACTS_DIR)/"; \
		else \
			echo "Warning: WASM file $$src_path not found."; \
		fi; \
	done
	@echo "All contracts built and WASM binaries copied to '$(ARTIFACTS_DIR)' successfully."

