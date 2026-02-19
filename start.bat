@echo off
REM Script untuk menjalankan Backend dan Frontend sekaligus

echo ===================================
echo Sistem Payroll Pemdes
echo ===================================
echo.

REM Cek apakah MySQL sudah berjalan
echo [1/4] Menjalankan Backend (Golang)...
start "Backend Server" cmd /k "cd /d %~dp0 && go run main.go"

REM Tunggu sebentar agar backend siap
timeout /t 3 /nobreak >nul

echo [2/4] Menjalankan Frontend (React)...
start "Frontend Server" cmd /k "cd /d %~dp0\frontend && npm run dev"

echo.
echo ===================================
echo Server berjalan!
echo ===================================
echo Backend:  http://localhost:3000
echo Frontend: http://localhost:5173
echo.
echo Tekan sembarang tombol untuk menutup jendela ini...
pause >nul
