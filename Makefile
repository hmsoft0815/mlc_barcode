# PDF Generation Service - Barcode CLI
# Copyright (c) 2026 Michael Lechner

.PHONY: all build build-linux build-windows build-macos build-all clean test run help

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

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(CLI_BINARY)-linux-amd64 ./cmd/barcode
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(SERVER_BINARY)-linux-amd64 ./cmd/mcp-server

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(CLI_BINARY)-windows-amd64.exe ./cmd/barcode
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(SERVER_BINARY)-windows-amd64.exe ./cmd/mcp-server

build-macos:
	@echo "Building for macOS (Apple Silicon)..."
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(CLI_BINARY)-macos-arm64 ./cmd/barcode
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(SERVER_BINARY)-macos-arm64 ./cmd/mcp-server
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(CLI_BINARY)-macos-amd64 ./cmd/barcode
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(SERVER_BINARY)-macos-amd64 ./cmd/mcp-server

build-all: build-linux build-windows build-macos

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
	@echo "  make build         - Build the binary in $(BIN_DIR)/"
	@echo "  make build-linux   - Build the binary for Linux amd64"
	@echo "  make build-windows - Build the binary for Windows amd64"
	@echo "  make build-macos   - Build the binary for macOS (arm64 & amd64)"
	@echo "  make build-all     - Build all cross-compile targets"
	@echo "  make clean         - Remove the $(BIN_DIR)/ directory and temporary files"
	@echo "  make test          - Run Go tests"
	@echo "  make run           - Build and run the binary"
	@echo "  make help          - Show this help message"
