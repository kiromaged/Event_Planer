@echo off
setlocal enabledelayedexpansion

echo ========================================
echo Event Planner - Complete Fix and Run
echo ========================================
echo.

echo Step 1: Checking Docker...
docker version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker is not running. Starting Docker Desktop...
    start "" "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    echo Waiting 90 seconds for Docker to start...
    timeout /t 90 /nobreak >nul
)

echo Step 2: Verifying Docker is ready...
docker ps >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Docker failed to start
    pause
    exit /b 1
)

echo Step 3: Cleaning up old containers...
docker compose down --remove-orphans 2>nul

echo Step 4: Building frontend with fixed Dockerfile...
docker compose build --no-cache frontend
if %errorlevel% neq 0 (
    echo ERROR: Frontend build failed
    pause
    exit /b 1
)

echo Step 5: Building backend...
docker compose build --no-cache backend
if %errorlevel% neq 0 (
    echo ERROR: Backend build failed
    pause
    exit /b 1
)

echo Step 6: Building MySQL...
docker compose build --no-cache mysql
if %errorlevel% neq 0 (
    echo ERROR: MySQL build failed
    pause
    exit /b 1
)

echo Step 7: Starting all services...
docker compose up -d

echo Step 8: Waiting for services to be ready...
timeout /t 15 /nobreak >nul

echo Step 9: Checking service status...
docker compose ps

echo.
echo ========================================
echo SUCCESS! Event Planner is running!
echo ========================================
echo.
echo Open your browser and go to:
echo http://localhost
echo.
echo Other URLs:
echo - Backend API: http://localhost:8080
echo - Backend Health: http://localhost:8080/api/ping
echo - Database: localhost:3307
echo.
pause