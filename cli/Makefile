.PHONY: help build test run clean lint

# Default target
help:
	@echo "CLI targets:"
	@echo "  make build  - Build the CLI binary"
	@echo "  make test   - Run CLI tests"
	@echo "  make run    - Run the CLI"
	@echo "  make lint   - Run linter"
	@echo "  make clean  - Remove build artifacts"

# Build the CLI
build:
	go build -o bin/cli ./cmd/cli

# Run tests
test:
	go test ./...

# Run
run:
	go run ./cmd/cli

# Lint
lint:
	golangci-lint run

# Clean
clean:
	rm -rf bin/
