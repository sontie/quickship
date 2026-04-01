$ErrorActionPreference = "Stop"

$REPO = "sontie/quickship"
$VERSION = "v" + (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$REPO/main/VERSION" -UseBasicParsing).Content.Trim()
$BINARY = "qship.exe"

Write-Host "Installing QuickShip $VERSION..." -ForegroundColor Green

# 检测架构
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# 构建下载 URL
$FILENAME = "qship-windows-$ARCH.exe"
$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

Write-Host "Downloading $FILENAME..."
$TMP_FILE = "$env:TEMP\$BINARY"

try {
    Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TMP_FILE -UseBasicParsing
} catch {
    Write-Host "Error: Failed to download $DOWNLOAD_URL" -ForegroundColor Red
    exit 1
}

# 安装到用户目录
$INSTALL_DIR = "$env:USERPROFILE\bin"
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

Move-Item -Path $TMP_FILE -Destination "$INSTALL_DIR\$BINARY" -Force

# 检查 PATH
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$INSTALL_DIR*") {
    Write-Host ""
    Write-Host "Adding $INSTALL_DIR to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$INSTALL_DIR", "User")
    Write-Host "Please restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "✓ QuickShip installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Usage:"
Write-Host "  qship version    # Check version"
Write-Host "  qship init       # Initialize configuration"
Write-Host "  qship deploy dev # Deploy to environment"
