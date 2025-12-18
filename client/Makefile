.PHONY: help build test clean lint install

# Default target
help:
	@echo "Client library targets:"
	@echo "  make build   - Build the client library"
	@echo "  make test    - Run client tests"
	@echo "  make lint    - Run linter"
	@echo "  make install - Install dependencies"
	@echo "  make clean   - Remove build artifacts"

# Install dependencies (if node_modules doesn't exist)
node_modules: package.json
	npm install
	@touch node_modules

# Install target for explicit install
install:
	npm install

# Build (depends on node_modules)
build: node_modules
	npm run build

# Run tests
test: node_modules
	npm test

# Lint
lint: node_modules
	npm run lint

# Clean
clean:
	rm -rf dist/ node_modules/
