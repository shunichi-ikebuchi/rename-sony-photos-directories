.PHONY: help build test lint clean install run

# Default target
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install the binary"
	@echo "  run        - Run the application"

# Build the binary
build:
	@echo "Building..."
	go build -o rspd ./cmd/rspd

# Run tests
test:
	@echo "Running tests..."
	go test -v -cover ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f rspd
	go clean

# Install the binary
install:
	@echo "Installing..."
	go install ./cmd/rspd

# Run the application
run: build
	./rspd
