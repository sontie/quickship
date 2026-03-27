@echo off
echo Installing QuickShip...

if not exist qship.exe (
    echo Error: qship.exe not found in current directory
    exit /b 1
)

set INSTALL_PATH=%USERPROFILE%\bin

if not exist "%INSTALL_PATH%" (
    mkdir "%INSTALL_PATH%"
)

copy qship.exe "%INSTALL_PATH%\" >nul

echo.
echo QuickShip installed to %INSTALL_PATH%
echo.
echo Please add %INSTALL_PATH% to your PATH environment variable:
echo 1. Open System Properties ^> Environment Variables
echo 2. Edit PATH and add: %INSTALL_PATH%
echo 3. Restart your terminal
echo.
echo Or run directly: %INSTALL_PATH%\qship.exe
pause
