#!/bin/bash

# Build script for Gemini Agentic CLI
set -e

echo "🔨 Building Gemini Agentic CLI..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ from https://golang.org"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
MIN_VERSION="1.21"

if [ "$(printf '%s\n' "$MIN_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go $MIN_VERSION or newer."
    exit 1
fi

# Create bin directory
mkdir -p bin

# Update dependencies
echo "📦 Updating Go modules..."
go mod tidy

# Build the application
echo "🔧 Building binary..."
go build -ldflags "-X main.version=$(date +%Y%m%d)" -o bin/goss cmd/goss/main.go

# Make executable
chmod +x bin/goss

# Check if build was successful
if [ -f "bin/goss" ]; then
    echo "✅ Build successful! Binary created at: bin/goss"
    echo ""
    echo "🚀 To run:"
    echo "  ./bin/goss"
    echo ""
    echo "📋 Prerequisites:"
    echo "  1. Start LM Studio with Local Server enabled"
    echo "  2. Load a function-calling capable model"
    echo "  3. Ensure server is running on http://localhost:1234"
    echo ""
    echo "🔧 Options:"
    echo "  ./bin/goss --help                    # Show help"
    echo "  ./bin/goss --base-url <url>          # Custom LM Studio URL"
    echo "  ./bin/goss --model <model-name>      # Specify model"
else
    echo "❌ Build failed!"
    exit 1
fi