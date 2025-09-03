# Variables
BINARY_NAME=app
BINARY_PATH=/tmp/$(BINARY_NAME)
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
.PHONY: all
all: clean build test

# Build target
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_PATH) ./...
	@echo "Build complete: $(BINARY_PATH)"

# Test target
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests complete"

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_PATH)
	@echo "Clean complete"

# Lint target - updated to use golangci-lint with PATH
.PHONY: lint
lint:
	@echo "Running linter..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && golangci-lint run
	@echo "Linting complete"

# Development targets
.PHONY: dev
dev: clean build test lint

# Run target
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BINARY_PATH)

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all     - Clean, build, and test"
	@echo "  build   - Build the application"
	@echo "  test    - Run tests"
	@echo "  clean   - Clean build artifacts"
	@echo "  lint    - Run golangci-lint"
	@echo "  dev     - Run clean, build, test, and lint"
	@echo "  run     - Build and run the application"
	@echo "  help    - Show this help message"
