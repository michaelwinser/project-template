.PHONY: help build test test-all test-ci clean up down review review-all
.PHONY: docker-dev docker-shell docker-build docker-clean

# Default target
help:
	@echo "Development targets:"
	@echo "  make build      - Build all components"
	@echo "  make test       - Run adaptive tests (based on changes)"
	@echo "  make test-all   - Run all tests"
	@echo "  make test-ci    - Run CI test suite"
	@echo "  make up         - Start all services"
	@echo "  make down       - Stop all services"
	@echo "  make clean      - Remove build artifacts"
	@echo ""
	@echo "Docker targets:"
	@echo "  make docker-dev   - Build dev container"
	@echo "  make docker-shell - Shell into dev container"
	@echo "  make docker-build - Build production images"
	@echo "  make docker-clean - Remove Docker images/volumes"
	@echo ""
	@echo "Review targets:"
	@echo "  make review     - Run adaptive reviews"
	@echo "  make review-all - Run all reviews"
	@echo ""
	@echo "Component targets:"
	@echo "  make server-*   - Server-specific targets (see server/Makefile)"
	@echo "  make client-*   - Client-specific targets (see client/Makefile)"
	@echo "  make cli-*      - CLI-specific targets (see cli/Makefile)"
	@echo "  make web-*      - Web-specific targets (see web/Makefile)"

# Build all components
build:
	$(MAKE) -C server build
	$(MAKE) -C client build
	$(MAKE) -C cli build
	$(MAKE) -C web build

# Adaptive test - decides what to test based on changes
test:
	@./scripts/adaptive-test.sh

# Run all tests
test-all:
	$(MAKE) -C server test
	$(MAKE) -C client test
	$(MAKE) -C cli test
	$(MAKE) -C web test
	$(MAKE) -C tests/integration test

# CI test suite
test-ci: test-all

# Start services
up:
	docker-compose up -d

# Stop services
down:
	docker-compose down

# Clean build artifacts
clean:
	$(MAKE) -C server clean
	$(MAKE) -C client clean
	$(MAKE) -C cli clean
	$(MAKE) -C web clean
	rm -rf generated/*

# Adaptive review
review:
	@./scripts/adaptive-review.sh

# Run all reviews
review-all:
	@./scripts/review-modularity.sh
	@./scripts/review-security.sh
	@./scripts/review-api-spec.sh

# Component-specific targets (pass-through)
server-%:
	$(MAKE) -C server $*

client-%:
	$(MAKE) -C client $*

cli-%:
	$(MAKE) -C cli $*

web-%:
	$(MAKE) -C web $*

# Docker targets
docker-dev:
	docker-compose build dev

docker-shell:
	docker-compose run --rm dev

docker-build:
	docker-compose build server

docker-clean:
	docker-compose down -v --rmi local
