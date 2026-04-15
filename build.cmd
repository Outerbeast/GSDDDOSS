@echo off

setlocal
:: Check for Go
where go >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH.
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

cd src
go build -ldflags="-s -w" -o GSDDDOSS.exe .
move /Y GSDDDOSS.exe ..\GSDDDOSS.exe
