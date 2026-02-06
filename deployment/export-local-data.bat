@echo off
REM Export Local Database Data to Production
REM This script exports data from local PostgreSQL

echo.
echo === Local Database Export ===
echo.

set LOCAL_HOST=localhost
set LOCAL_PORT=5432
set LOCAL_USER=postgres
set LOCAL_DB=medical_equipment
set OUTPUT_FILE=local_data_export.sql

echo Exporting data from local database...
echo   Host: %LOCAL_HOST%
echo   Database: %LOCAL_DB%
echo   Output: %OUTPUT_FILE%
echo.

REM Check if pg_dump exists
where pg_dump >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: pg_dump not found. Please install PostgreSQL client tools.
    echo Download from: https://www.postgresql.org/download/windows/
    pause
    exit /b 1
)

REM Export data
echo Enter PostgreSQL password when prompted...
echo.

pg_dump -h %LOCAL_HOST% -p %LOCAL_PORT% -U %LOCAL_USER% -d %LOCAL_DB% --data-only --inserts --disable-triggers --no-owner --no-privileges -f %OUTPUT_FILE%

if %ERRORLEVEL% EQU 0 (
    echo.
    echo [SUCCESS] Data exported successfully!
    echo.
    echo File: %OUTPUT_FILE%
    echo.
    echo Next steps:
    echo 1. Upload this file to your production server
    echo 2. Run the import script on the server
    echo.
    echo Upload command:
    echo   scp %OUTPUT_FILE% root@158.69.118.34:/tmp/
    echo.
) else (
    echo.
    echo [ERROR] Export failed!
    echo.
)

pause
