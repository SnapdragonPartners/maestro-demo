# Makefile for Go development in containerized environment
# Designed to work with security constraints: rootless, read-only filesystem, no network

.PHONY: build test lint run clean help

# Default target
all: build test lint

# Build the Go application
build:
	@echo "Building Go application..."
	go mod tidy
	mkdir -p bin
	go build -o bin/hello .
	@echo "Build successful - executable created at bin/hello"

# Run tests
test:
	@echo "Running Go tests..."
	go test ./...
	@echo "Tests completed successfully"

# Lint using golangci-lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run
	@echo "Linting completed successfully"

# Run the application directly
run:
	@echo "Running application..."
	go run ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean completed"

# Display available targets
help:
	@echo "Available targets:"
	@echo "  build  - Build the Go application (runs go mod tidy, builds binary at bin/hello)"
	@echo "  test   - Run tests (go test ./...)"
	@echo "  lint   - Run linting (golangci-lint run)"
	@echo "  run    - Run the application (go run ./...)"
	@echo "  clean  - Clean build artifacts (removes bin/ directory)"
	@echo "  all    - Run build, test, and lint"
	@echo "  help   - Show this help message"

# Development workflow validation
validate: build test lint
	@echo "=== Container Build Pipeline Validation ==="
	@echo "✓ Build pipeline: PASSED"
	@echo "✓ Test framework: PASSED" 
	@echo "✓ Linting tools: PASSED"
	@echo "✓ Development workflow: VALIDATED"
