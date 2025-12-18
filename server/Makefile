.PHONY: help build test run clean lint

# Default target
help:
	@echo "Server targets:"
	@echo "  make build  - Build the server binary"
	@echo "  make test   - Run server tests"
	@echo "  make run    - Run the server locally"
	@echo "  make lint   - Run linter"
	@echo "  make clean  - Remove build artifacts"

# Build the server
build:
	go build -o bin/server ./cmd/server

# Run tests
test:
	go test ./...

# Run the server
run:
	go run ./cmd/server

# Lint
lint:
	golangci-lint run

# Clean
clean:
	rm -rf bin/
