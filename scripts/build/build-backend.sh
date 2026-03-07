#!/bin/bash

# Build script for OpenClaw Music Player Backend

set -e

echo "Building OpenClaw Music Player Backend..."

# Change to backend directory
cd "$(dirname "$0")/../../backend"

# Clean previous builds
rm -rf dist
mkdir -p dist

# Build for multiple platforms
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/server-linux-amd64 cmd/server/main.go

echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/server-darwin-amd64 cmd/server/main.go
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/server-darwin-arm64 cmd/server/main.go

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/server-windows-amd64.exe cmd/server/main.go

echo "Build complete!"
echo "Binaries located in: $(pwd)/dist"
ls -lh dist/
