@echo off
:: -------------------- Auto-elevate --------------------
net session >nul 2>&1
if %errorlevel% neq 0 (
    powershell -Command "Start-Process '%~f0' -Verb runAs"
    exit /b
)

:: -------------------- Config --------------------------
set SERVICE_NAME=PrintRawWeb

:: If service doesn't exist, exit quietly
sc.exe query "%SERVICE_NAME%" >nul 2>&1
if %errorlevel% neq 0 exit /b 0

:: Stop (ignore error if already stopped)
sc.exe stop "%SERVICE_NAME%" >nul 2>&1

:: Optional: brief wait to allow stop to complete
timeout /t 2 >nul

exit /b 0
