.PHONY: build test clean run help

BINARY_NAME=subdomain-finder
BUILD_DIR=build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -domain example.com

run-verbose: build
	@echo "Running $(BINARY_NAME) with verbose output..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -domain example.com -verbose

run-custom: build
	@echo "Running $(BINARY_NAME) with custom wordlist..."
	@./$(BUILD_DIR)/$(BINARY_NAME) -domain example.com -wordlist wordlist.txt -verbose

install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

lint: fmt vet

help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Build and run with example.com"
	@echo "  run-verbose   - Build and run with verbose output"
	@echo "  run-custom    - Build and run with custom wordlist"
	@echo "  install-deps  - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run fmt and vet"
	@echo "  help          - Show this help message"
