# PDF Generation Service - Barcode CLI
# Copyright (c) 2026 Michael Lechner

.PHONY: all build clean test run help

# Binary names
CLI_BINARY=barcode
SERVER_BINARY=mcp-barcode-server
# Output directory
BIN_DIR=bin

all: build

build:
	@echo "Building..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(CLI_BINARY) ./cmd/barcode
	go build -o $(BIN_DIR)/$(SERVER_BINARY) ./cmd/mcp-server

clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -f *.svg *.png

test:
	@echo "Running tests..."
	go test ./...

run: build
	@./$(BIN_DIR)/$(CLI_BINARY)

run-server: build
	@./$(BIN_DIR)/$(SERVER_BINARY)

help:
	@echo "Available commands:"
	@echo "  make build  - Build the binary in $(BIN_DIR)/"
	@echo "  make clean  - Remove the $(BIN_DIR)/ directory and temporary files"
	@echo "  make test   - Run Go tests"
	@echo "  make run    - Build and run the binary"
	@echo "  make help   - Show this help message"
