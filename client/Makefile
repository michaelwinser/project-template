.PHONY: help build test clean lint

# Default target
help:
	@echo "Client library targets:"
	@echo "  make build  - Build the client library"
	@echo "  make test   - Run client tests"
	@echo "  make lint   - Run linter"
	@echo "  make clean  - Remove build artifacts"

# Build
build:
	npm run build

# Run tests
test:
	npm test

# Lint
lint:
	npm run lint

# Clean
clean:
	rm -rf dist/ node_modules/
