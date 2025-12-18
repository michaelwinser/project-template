.PHONY: help build test dev clean lint typecheck install

# Default target
help:
	@echo "Web client targets:"
	@echo "  make build     - Build the web client"
	@echo "  make test      - Run web client tests"
	@echo "  make dev       - Watch for changes and rebuild"
	@echo "  make lint      - Run linter"
	@echo "  make typecheck - Run TypeScript type checking"
	@echo "  make install   - Install dependencies"
	@echo "  make clean     - Remove build artifacts"

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

# Watch mode
dev: node_modules
	npm run dev

# Lint
lint: node_modules
	npm run lint

# Type check
typecheck: node_modules
	npm run typecheck

# Clean
clean:
	rm -rf dist/ node_modules/
