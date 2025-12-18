.PHONY: help build test dev clean lint

# Default target
help:
	@echo "Web client targets:"
	@echo "  make build  - Build the web client"
	@echo "  make test   - Run web client tests"
	@echo "  make dev    - Run dev server"
	@echo "  make lint   - Run linter"
	@echo "  make clean  - Remove build artifacts"

# Build
build:
	npm run build

# Run tests
test:
	npm test

# Dev server
dev:
	npm run dev

# Lint
lint:
	npm run lint

# Clean
clean:
	rm -rf dist/ node_modules/
