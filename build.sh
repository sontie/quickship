#!/bin/bash
set -e

VERSION="1.0.0"
OUTPUT_DIR="dist"

rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

echo "Building QuickShip v${VERSION}..."

# Linux amd64
GOOS=linux GOARCH=amd64 go build -o "$OUTPUT_DIR/qship"
cp install.sh "$OUTPUT_DIR/"
(cd "$OUTPUT_DIR" && tar czf "../quickship-${VERSION}-linux-amd64.tar.gz" qship install.sh)
rm "$OUTPUT_DIR/qship"
echo "✓ linux-amd64"

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o "$OUTPUT_DIR/qship"
cp install.sh "$OUTPUT_DIR/"
(cd "$OUTPUT_DIR" && tar czf "../quickship-${VERSION}-darwin-amd64.tar.gz" qship install.sh)
rm "$OUTPUT_DIR/qship"
echo "✓ darwin-amd64"

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o "$OUTPUT_DIR/qship"
cp install.sh "$OUTPUT_DIR/"
(cd "$OUTPUT_DIR" && tar czf "../quickship-${VERSION}-darwin-arm64.tar.gz" qship install.sh)
rm "$OUTPUT_DIR/qship"
echo "✓ darwin-arm64"

# Windows
GOOS=windows GOARCH=amd64 go build -o "$OUTPUT_DIR/qship.exe"
cp install.bat "$OUTPUT_DIR/"
(cd "$OUTPUT_DIR" && zip -q "../quickship-${VERSION}-windows-amd64.zip" qship.exe install.bat)
rm "$OUTPUT_DIR/qship.exe"
echo "✓ windows-amd64"

rm -rf "$OUTPUT_DIR"

echo ""
echo "Done! Distribution packages:"
ls -lh quickship-${VERSION}-*
