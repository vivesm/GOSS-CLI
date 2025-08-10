# GOSS Agentic CLI Makefile

.PHONY: all build clean test fmt vet deps install uninstall run help

# Variables
BINARY_NAME=gossai
BUILD_DIR=bin
MAIN_PATH=cmd/goss/main.go
VERSION=$(shell date +%Y%m%d)

# Default target
all: build

# Build the binary
build: deps
	@echo "🔧 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "🔍 Vetting code..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean

# Install binary to system
install: build
	@echo "📦 Installing $(BINARY_NAME) to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	# Build directly to final location to avoid macOS security issues
	go build -ldflags "-X main.version=$(VERSION)" -o ~/.local/bin/$(BINARY_NAME) $(MAIN_PATH)
	chmod +x ~/.local/bin/$(BINARY_NAME)
	# Remove any extended attributes that could cause issues
	xattr -cr ~/.local/bin/$(BINARY_NAME) 2>/dev/null || true
	@echo "✅ Installed successfully as '$(BINARY_NAME)'"
	@echo "💡 Use '$(BINARY_NAME)' to start the CLI"
	@echo "💡 Optional: Add 'alias goss=$(BINARY_NAME)' to your shell config"

# Uninstall binary from system
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	rm -f ~/.local/bin/$(BINARY_NAME)
	@echo "✅ Uninstalled successfully"

# Run the application
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@echo "Prerequisites:"
	@echo "  1. Start LM Studio with Local Server enabled"
	@echo "  2. Load a function-calling capable model"
	@echo "  3. Ensure server is running on http://localhost:1234"
	@echo ""
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with custom base URL
run-custom:
	@echo "🚀 Running $(BINARY_NAME) with custom settings..."
	./$(BUILD_DIR)/$(BINARY_NAME) --base-url $(BASE_URL) --model $(MODEL)

# Development mode - build and run
dev: build run

# Check code quality
check: fmt vet test

# Build release version
release: clean check build
	@echo "🎉 Release build complete!"

# Show help
help:
	@echo "GOSS Agentic CLI - Available targets:"
	@echo ""
	@echo "  build      - Build the binary"
	@echo "  deps       - Install Go dependencies"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format Go code"
	@echo "  vet        - Vet Go code"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install binary to system"
	@echo "  uninstall  - Uninstall binary from system"
	@echo "  run        - Build and run the application"
	@echo "  run-custom - Run with custom BASE_URL and MODEL"
	@echo "  dev        - Development mode (build + run)"
	@echo "  check      - Run code quality checks (fmt + vet + test)"
	@echo "  release    - Build release version"
	@echo "  help       - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run"
	@echo "  make run-custom BASE_URL=http://localhost:8080/v1 MODEL=mistral-7b"
	@echo "  make install    # Installs as 'gossai'"