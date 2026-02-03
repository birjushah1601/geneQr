@echo off
cls
echo ============================================
echo Starting Next.js Fresh (No Cache)
echo ============================================
echo.
cd /d C:\Users\birju\aby-med\admin-ui
echo Checking file...
echo First line should be 'use client';
echo.
echo Starting Next.js...
"C:\Program Files\nodejs\node.exe" node_modules\next\dist\bin\next dev
pause
