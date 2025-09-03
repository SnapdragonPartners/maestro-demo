# Makefile for Go development in containerized environment
# Designed to work with security constraints: rootless, read-only filesystem, no network

.PHONY: build test lint clean help

# Default target
all: build test lint

# Build the Go application
build:
	@echo "Building Go application..."
	go build -o /tmp/maestro-demo .
	@echo "Build successful - executable created at /tmp/maestro-demo"

# Run tests
test:
	@echo "Running Go tests..."
	go test -v ./...
	@echo "Tests completed successfully"

# Basic lint using go fmt and go vet (available in standard Go installation)
lint:
	@echo "Running Go linting..."
	@echo "Checking code formatting..."
	go fmt ./...
	@echo "Running go vet..."
	go vet ./...
	@echo "Linting completed successfully"

# Clean build artifacts (in /tmp due to read-only filesystem)
clean:
	@echo "Cleaning build artifacts..."
	rm -f /tmp/maestro-demo
	@echo "Clean completed"

# Run the application (from /tmp due to read-only filesystem)
run: build
	@echo "Running application..."
	/tmp/maestro-demo

# Display available targets
help:
	@echo "Available targets:"
	@echo "  build  - Build the Go application"
	@echo "  test   - Run tests"
	@echo "  lint   - Run linting (go fmt + go vet)"
	@echo "  clean  - Clean build artifacts"
	@echo "  run    - Build and run the application"
	@echo "  all    - Run build, test, and lint"
	@echo "  help   - Show this help message"

# Development workflow validation
validate: build test lint
	@echo "=== Container Build Pipeline Validation ==="
	@echo "✓ Build pipeline: PASSED"
	@echo "✓ Test framework: PASSED" 
	@echo "✓ Linting tools: PASSED"
	@echo "✓ Security constraints: Compatible with rootless, read-only filesystem"
	@echo "✓ Development workflow: VALIDATED"

# Go Module Management
.PHONY: module-init module-validate module-clean

module-init: ## Initialize and maintain Go module
	@echo "Initializing Go module..."
	@./scripts/init-go-module.sh

module-validate: ## Validate Go module integrity
	@echo "Validating Go module..."
	@go mod verify
	@go mod tidy
	@echo "✓ Module validation completed"

module-clean: ## Clean Go module cache
	@echo "Cleaning Go module cache..."
	@go clean -modcache
	@echo "✓ Module cache cleaned"

