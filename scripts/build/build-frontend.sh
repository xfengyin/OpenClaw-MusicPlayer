#!/bin/bash

# Build script for OpenClaw Music Player Frontend

set -e

echo "Building OpenClaw Music Player Frontend..."

# Change to frontend directory
cd "$(dirname "$0")/../../frontend/web"

# Install dependencies
echo "Installing dependencies..."
npm install

# Build web app
echo "Building web app..."
npm run build

# Copy to electron renderer
echo "Copying to Electron..."
mkdir -p ../electron/renderer
cp -r dist/* ../electron/renderer/

echo "Build complete!"
echo "Web app located in: $(pwd)/dist"
echo "Electron renderer located in: $(pwd)/../electron/renderer"
