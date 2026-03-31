$ErrorActionPreference = "Stop"

$Binary = "qship.exe"
$InstallPath = "$env:USERPROFILE\bin"
$BinaryPath = "$InstallPath\$Binary"
$ConfigDir = "$env:USERPROFILE\.quickship"
$ConfigFile = "$env:USERPROFILE\.quickship.yaml"

Write-Host "Uninstalling QuickShip..." -ForegroundColor Cyan

# 删除二进制文件
if (Test-Path $BinaryPath) {
    Remove-Item $BinaryPath -Force
    Write-Host "✓ Removed binary: $BinaryPath" -ForegroundColor Green
}

# 删除安装目录
if (Test-Path $InstallPath) {
    Remove-Item $InstallPath -Recurse -Force
    Write-Host "✓ Removed install directory: $InstallPath" -ForegroundColor Green
}

# 删除配置目录
if (Test-Path $ConfigDir) {
    Remove-Item $ConfigDir -Recurse -Force
    Write-Host "✓ Removed config directory: $ConfigDir" -ForegroundColor Green
}

# 删除配置文件
if (Test-Path $ConfigFile) {
    Remove-Item $ConfigFile -Force
    Write-Host "✓ Removed config file: $ConfigFile" -ForegroundColor Green
}

# 从 PATH 中移除
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -like "*$InstallPath*") {
    $NewPath = ($UserPath -split ';' | Where-Object { $_ -ne $InstallPath }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Host "✓ Removed from PATH" -ForegroundColor Green
}

Write-Host ""
Write-Host "✓ QuickShip uninstalled successfully!" -ForegroundColor Green
Write-Host "Please restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
