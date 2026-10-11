# Go の gate 一式 (版の照合 / go mod tidy の drift / vet / go test / lint / 既知脆弱性)。
# 写しは書き換えない (gate-kit-verify が落とす)。gate を外すなら include の前に
#   GATE_KIT_GO_SKIP := go-vuln
# を書く (理由もコメントで残す)。外したものは `make gate-kit-gates` に skip: として出る。
include third_party/gate-kit/go/gate.mk

.PHONY: help build test lint ci clean install run

# Default target
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  ci         - Run the full gate (gate-kit go gates)"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install the binary"
	@echo "  run        - Run the application"

# Full gate. `make gate-kit-gates` lists what runs and what is skipped.
ci: $(GATE_KIT_GO_GATES)
	@echo "✓ ci passed"

# Build the binary
build:
	@echo "Building..."
	go build -o rspd ./cmd/rspd

# Run tests (人が読む用。ci が回すのは gate-kit の go-test = `go test ./...`)
test:
	@echo "Running tests..."
	go test -v -cover ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

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
