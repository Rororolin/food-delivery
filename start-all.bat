@echo off
title Start All Delivery Servers
color 0A

set "PROJECT_DIR=%~dp0"

echo [1/3] Starting Client Server...
start "CLIENT" cmd /k "cd /d "%PROJECT_DIR%food-delivery" && food-delivery.exe"

echo [2/3] Starting Admin Server...
start "ADMIN" cmd /k "cd /d "%PROJECT_DIR%food-delivery-admin" && food-delivery-admin.exe"

echo [3/3] Starting Courier Server...
start "COURIER" cmd /k "cd /d "%PROJECT_DIR%food-delivery-courier" && food-delivery-courier.exe"

echo.
echo All servers started!
pause