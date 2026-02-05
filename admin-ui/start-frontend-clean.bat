@echo off
cd /d "C:\Users\birju\aby-med\admin-ui"
set NODE_ENV=
echo.
echo ============================================================
echo   FRONTEND STARTING - CLEAN BUILD
echo ============================================================
echo.
echo Compiling Next.js application...
echo.
npm run dev
pause
