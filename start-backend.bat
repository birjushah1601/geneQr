@echo off
set DB_HOST=localhost
set DB_PORT=5433
set DB_USER=postgres
set DB_PASSWORD=postgres
set DB_NAME=medplatform
set PORT=8081
set BASE_URL=http://localhost:3000
set ENABLED_MODULES=equipment-registry
set DATABASE_URL=postgres://postgres:postgres@localhost:5433/medplatform?sslmode=disable

start /B medical-platform.exe > backend.log 2>&1
