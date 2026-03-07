#!/bin/bash

# Build script for OpenClaw Music Player Desktop App

set -e

echo "Building OpenClaw Music Player Desktop App..."

# Build frontend first
"$(dirname "$0")/build-frontend.sh"

# Build backend
echo "Building backend..."
"$(dirname "$0")/build-backend.sh"

# Change to electron directory
cd "$(dirname "$0")/../../frontend/electron"

# Install dependencies
echo "Installing Electron dependencies..."
npm install

# Copy backend binary to extra resources
echo "Copying backend binary..."
mkdir -p dist
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    cp ../../backend/dist/server-linux-amd64 dist/server
elif [[ "$OSTYPE" == "darwin"* ]]; then
    cp ../../backend/dist/server-darwin-amd64 dist/server
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    cp ../../backend/dist/server-windows-amd64.exe dist/server.exe
fi

# Build Electron app
echo "Building Electron app..."

# Detect platform
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    npx electron-builder --linux
elif [[ "$OSTYPE" == "darwin"* ]]; then
    npx electron-builder --mac
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    npx electron-builder --win
else
    echo "Unknown platform, building for all..."
    npx electron-builder
fi

echo "Build complete!"
echo "Desktop app located in: $(pwd)/dist"
