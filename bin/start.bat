@echo off
:: ---------------------------------------------------------------
:: Auto-elevate to Administrator
:: ---------------------------------------------------------------
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo Requesting administrator privileges...
    powershell -Command "Start-Process '%~f0' -Verb runAs"
    exit /b
)

:: ---------------------------------------------------------------
:: Start Service Script
:: ---------------------------------------------------------------
set SERVICE_NAME=PrintRawWeb

echo Starting service "%SERVICE_NAME%"...
sc.exe start "%SERVICE_NAME%" >nul 2>&1

if %errorlevel% neq 0 (
    echo Failed to start service "%SERVICE_NAME%".
    timeout /t 3 >nul
    exit /b
)

echo Service "%SERVICE_NAME%" started successfully.
exit /b
