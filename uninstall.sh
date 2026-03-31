#!/bin/bash
set -e

BINARY="qship"
INSTALL_PATH="/usr/local/bin/$BINARY"
CONFIG_DIR="$HOME/.quickship"
CONFIG_FILE="$HOME/.quickship.yaml"

echo "Uninstalling QuickShip..."

# 删除二进制文件
if [ -f "$INSTALL_PATH" ]; then
    if [ -w "/usr/local/bin" ]; then
        rm "$INSTALL_PATH"
    else
        sudo rm "$INSTALL_PATH"
    fi
    echo "✓ Removed binary: $INSTALL_PATH"
else
    echo "Binary not found at $INSTALL_PATH"
fi

# 删除配置目录
if [ -d "$CONFIG_DIR" ]; then
    rm -rf "$CONFIG_DIR"
    echo "✓ Removed config directory: $CONFIG_DIR"
fi

# 删除配置文件
if [ -f "$CONFIG_FILE" ]; then
    rm "$CONFIG_FILE"
    echo "✓ Removed config file: $CONFIG_FILE"
fi

echo ""
echo "✓ QuickShip uninstalled successfully!"
