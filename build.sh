#!/bin/bash
set -e

# 从 Git tag 读取版本号，如果没有 tag 则使用 dev
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
OUTPUT_DIR="release"

rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

echo "Building QuickShip v${VERSION}..."

LDFLAGS="-X main.Version=${VERSION}"

GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$OUTPUT_DIR/qship-linux-amd64"
echo "✓ linux-amd64"

GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$OUTPUT_DIR/qship-darwin-amd64"
echo "✓ darwin-amd64"

GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$OUTPUT_DIR/qship-darwin-arm64"
echo "✓ darwin-arm64"

GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$OUTPUT_DIR/qship-windows-amd64.exe"
echo "✓ windows-amd64"

echo ""
echo "Done! Distribution files:"
ls -lh "$OUTPUT_DIR"/
