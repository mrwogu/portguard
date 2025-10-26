#!/bin/bash
# Test local release build
# This script simulates what the GitHub Actions release workflow does

set -e

echo "🚀 Testing release build process..."
echo ""

# Get version
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0-test")}
echo "📦 Version: $VERSION"
echo ""

# Clean previous builds
echo "🧹 Cleaning previous builds..."
make clean
echo ""

# Run tests
echo "🧪 Running tests..."
make test
echo ""

# Run linter
echo "🔍 Running linter..."
if ! make lint; then
    echo "❌ Linter failed. Fix issues before release."
    exit 1
fi
echo ""

# Build for all platforms
echo "🔨 Building for all platforms..."
VERSION=$VERSION make build-all
echo ""

# Package binaries
echo "📦 Packaging binaries..."
cd dist

# Linux AMD64
echo "  - Linux AMD64..."
tar -czf portguard-linux-amd64.tar.gz portguard-linux-amd64
shasum -a 256 portguard-linux-amd64.tar.gz > portguard-linux-amd64.tar.gz.sha256

# Linux ARM64
echo "  - Linux ARM64..."
tar -czf portguard-linux-arm64.tar.gz portguard-linux-arm64
shasum -a 256 portguard-linux-arm64.tar.gz > portguard-linux-arm64.tar.gz.sha256

# Darwin AMD64
echo "  - Darwin AMD64..."
tar -czf portguard-darwin-amd64.tar.gz portguard-darwin-amd64
shasum -a 256 portguard-darwin-amd64.tar.gz > portguard-darwin-amd64.tar.gz.sha256

# Darwin ARM64
echo "  - Darwin ARM64..."
tar -czf portguard-darwin-arm64.tar.gz portguard-darwin-arm64
shasum -a 256 portguard-darwin-arm64.tar.gz > portguard-darwin-arm64.tar.gz.sha256

# Windows AMD64
echo "  - Windows AMD64..."
zip -q portguard-windows-amd64.zip portguard-windows-amd64.exe
shasum -a 256 portguard-windows-amd64.zip > portguard-windows-amd64.zip.sha256

cd ..
echo ""

# List artifacts
echo "✅ Release artifacts created:"
echo ""
ls -lh dist/*.tar.gz dist/*.zip dist/*.sha256 | awk '{printf "  %s  %s\n", $5, $9}'
echo ""

# Verify checksums
echo "🔐 Verifying checksums..."
cd dist
for file in *.sha256; do
    if shasum -a 256 -c "$file" >/dev/null 2>&1; then
        echo "  ✓ ${file%.sha256}"
    else
        echo "  ✗ ${file%.sha256}"
    fi
done
cd ..
echo ""

# Test one binary
echo "🧪 Testing binary..."
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
fi

BINARY="dist/portguard-${PLATFORM}-${ARCH}"
if [ -f "$BINARY" ]; then
    echo "  Running: $BINARY --version"
    $BINARY --version
    echo ""
else
    echo "  ⚠️  Binary not found for this platform: $BINARY"
    echo ""
fi

echo "✅ Release build test complete!"
echo ""
echo "📋 Next steps:"
echo "  1. Create and push a tag: git tag -a $VERSION -m 'Release $VERSION' && git push origin $VERSION"
echo "  2. GitHub Actions will automatically create a release"
echo "  3. Check: https://github.com/mrwogu/portguard/releases"
