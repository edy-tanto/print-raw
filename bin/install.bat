@echo off
:: -------------------- Auto-elevate --------------------
net session >nul 2>&1
if %errorlevel% neq 0 (
    powershell -Command "Start-Process '%~f0' -Verb runAs"
    exit /b
)

:: -------------------- Config --------------------------
set SERVICE_NAME=PrintRawWeb
set SERVICE_DISPLAY=Print Raw Web

:: Get current folder of this .bat file (independent of where it's launched)
set SCRIPT_DIR=%~dp0
:: Remove trailing backslash if present
if "%SCRIPT_DIR:~-1%"=="\" set SCRIPT_DIR=%SCRIPT_DIR:~0,-1%

:: Build full path to EXE (assuming it's in same folder as this .bat)
set EXE_PATH=%SCRIPT_DIR%\print_web_service.exe

echo Installing from "%EXE_PATH%"

:: -------------------- Install Logic -------------------
:: If service already exists, just start it and exit
sc.exe query "%SERVICE_NAME%" >nul 2>&1
if %errorlevel% equ 0 (
    echo Service "%SERVICE_NAME%" already exists. Starting...
    sc.exe start "%SERVICE_NAME%" >nul 2>&1
    exit /b 0
)

:: Create service
sc.exe create "%SERVICE_NAME%" binPath= "%EXE_PATH%" start= auto DisplayName= "%SERVICE_DISPLAY%"
if %errorlevel% neq 0 goto :fail_create

:: Start service
sc.exe start "%SERVICE_NAME%" >nul 2>&1
if %errorlevel% neq 0 goto :fail_start

echo Service "%SERVICE_NAME%" installed and started successfully.
exit /b 0

:fail_create
echo Failed to create service "%SERVICE_NAME%".
timeout /t 5 >nul
exit /b 1

:fail_start
echo Service "%SERVICE_NAME%" was created but failed to start.
timeout /t 5 >nul
exit /b 1
