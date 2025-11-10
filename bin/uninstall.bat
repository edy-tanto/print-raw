@echo off
:: -------------------- Auto-elevate --------------------
net session >nul 2>&1
if %errorlevel% neq 0 (
    powershell -Command "Start-Process '%~f0' -Verb runAs"
    exit /b
)

:: -------------------- Config --------------------------
set SERVICE_NAME=PrintRawWeb

:: Stop if running (ignore errors)
sc.exe stop "%SERVICE_NAME%" >nul 2>&1
timeout /t 2 >nul

:: Delete service
sc.exe delete "%SERVICE_NAME%"
if %errorlevel% neq 0 goto :fail_delete

exit /b 0

:fail_delete
echo Failed to delete service "%SERVICE_NAME%".
timeout /t 5 >nul
exit /b 1
