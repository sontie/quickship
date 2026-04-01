#!/bin/bash
set -e

REPO="sontie/quickship"
VERSION="v$(curl -fsSL https://raw.githubusercontent.com/${REPO}/main/VERSION)"
BINARY="qship"

echo "Installing QuickShip ${VERSION}..."

# 检测操作系统和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# 构建下载 URL
FILENAME="${BINARY}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

echo "Downloading ${FILENAME}..."
TMP_FILE="/tmp/${BINARY}"

if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"
elif command -v wget >/dev/null 2>&1; then
    wget -q "$DOWNLOAD_URL" -O "$TMP_FILE"
else
    echo "Error: curl or wget is required"
    exit 1
fi

chmod +x "$TMP_FILE"

# 安装到系统路径
INSTALL_PATH="/usr/local/bin"

if [ -w "$INSTALL_PATH" ]; then
    mv "$TMP_FILE" "$INSTALL_PATH/$BINARY"
else
    echo "Installing to $INSTALL_PATH (requires sudo)..."
    sudo mv "$TMP_FILE" "$INSTALL_PATH/$BINARY"
fi

echo "✓ QuickShip installed successfully!"
echo ""
echo "Usage:"
echo "  qship version    # Check version"
echo "  qship init       # Initialize configuration"
echo "  qship deploy dev # Deploy to environment"
