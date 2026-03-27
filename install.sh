#!/bin/bash
set -e

echo "Installing QuickShip..."

# 检测操作系统
OS=$(uname -s)
ARCH=$(uname -m)

# 设置二进制文件名
BINARY="qship"

# 检查二进制文件是否存在
if [ ! -f "$BINARY" ]; then
    echo "Error: $BINARY not found in current directory"
    echo "Please run this script in the same directory as the qship binary"
    exit 1
fi

# 添加执行权限
chmod +x "$BINARY"

# 根据操作系统选择安装路径
case "$OS" in
    Linux*)
        INSTALL_PATH="/usr/local/bin"
        ;;
    Darwin*)
        INSTALL_PATH="/usr/local/bin"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# 检查是否需要 sudo
if [ -w "$INSTALL_PATH" ]; then
    mv "$BINARY" "$INSTALL_PATH/"
else
    echo "Installing to $INSTALL_PATH (requires sudo)..."
    sudo mv "$BINARY" "$INSTALL_PATH/"
fi

echo "✓ QuickShip installed successfully to $INSTALL_PATH/$BINARY"
echo ""
echo "Usage:"
echo "  qship init       # Initialize configuration"
echo "  qship deploy dev # Deploy to environment"
echo ""
echo "Run 'qship' to see all commands"
