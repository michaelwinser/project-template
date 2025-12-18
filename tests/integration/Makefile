.PHONY: help test test-go

# Default target
help:
	@echo "Integration test targets:"
	@echo "  make test     - Run all integration tests"
	@echo "  make test-go  - Run Go integration tests"

# Run all integration tests
test: test-go

# Run Go integration tests (requires server to be running)
test-go:
	@echo "Running Go integration tests..."
	@echo "Note: Ensure the server is running at TEST_SERVER_URL (default: http://localhost:8080)"
	go test -v ./...
