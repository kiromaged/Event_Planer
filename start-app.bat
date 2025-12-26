@echo off
echo Starting Event Planner Application...

echo Checking Docker status...
docker version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker is not running. Starting Docker Desktop...
    start "" "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    echo Waiting for Docker to start (60 seconds)...
    timeout /t 60 /nobreak >nul
)

echo Building frontend with fixed routing...
docker compose build --no-cache frontend

echo Starting all services...
docker compose up -d

echo Waiting for services to be ready...
timeout /t 10 /nobreak >nul

echo Checking service status...
docker compose ps

echo.
echo ========================================
echo Event Planner is ready!
echo ========================================
echo.
echo Open your browser and go to:
echo http://localhost
echo.
echo Backend API: http://localhost:8080
echo Database: localhost:3307
echo.
pause