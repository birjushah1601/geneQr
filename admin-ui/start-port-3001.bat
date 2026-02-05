@echo off
echo Starting on PORT 3001 (fresh port, no cache)...
cd /d C:\Users\birju\aby-med\admin-ui
set PORT=3001
"C:\Program Files\nodejs\node.exe" node_modules\next\dist\bin\next dev -p 3001
pause
