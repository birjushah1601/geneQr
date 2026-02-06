@echo off
REM Export data from local Docker PostgreSQL container

echo.
echo === Export Data from Docker Container ===
echo.

set OUTPUT_FILE=local_data_export.sql
set CONTAINER_NAME=med_platform_pg
set DATABASE_NAME=med_platform

echo Checking if Docker container is running...
docker ps --filter "name=%CONTAINER_NAME%" --format "{{.Names}}" | findstr /C:"%CONTAINER_NAME%" >nul 2>&1

if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Docker container '%CONTAINER_NAME%' is not running.
    echo.
    echo Please start your local development environment:
    echo   cd dev
    echo   docker-compose up -d
    echo.
    pause
    exit /b 1
)

echo.
echo Exporting data from Docker container...
echo Container: %CONTAINER_NAME%
echo Database: %DATABASE_NAME%
echo Output: %OUTPUT_FILE%
echo.

REM Export data from Docker container
docker exec %CONTAINER_NAME% pg_dump -U postgres -d %DATABASE_NAME% --data-only --inserts --disable-triggers --no-owner --no-privileges > %OUTPUT_FILE%

if %ERRORLEVEL% EQU 0 (
    echo.
    echo [SUCCESS] Data exported successfully!
    echo.
    
    REM Get file size
    for %%A in (%OUTPUT_FILE%) do set FILE_SIZE=%%~zA
    echo File: %OUTPUT_FILE%
    echo Size: %FILE_SIZE% bytes
    echo.
    echo Next steps:
    echo 1. Upload this file to your production server
    echo.
    echo Upload command ^(use Git Bash or WSL^):
    echo   scp %OUTPUT_FILE% root@158.69.118.34:/tmp/
    echo.
    echo Or use WinSCP / FileZilla to upload the file
    echo.
) else (
    echo.
    echo [ERROR] Export failed!
    echo Make sure the Docker container is running and accessible.
    echo.
)

pause
