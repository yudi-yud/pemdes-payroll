#!/bin/bash

echo "==================================="
echo "Sistem Payroll Pemdes"
echo "==================================="
echo ""

# Jalankan backend di background
echo "[1/2] Menjalankan Backend (Golang)..."
go run main.go &
BACKEND_PID=$!

# Tunggu backend siap
sleep 3

# Jalankan frontend di background
echo "[2/2] Menjalankan Frontend (React)..."
cd frontend && npm run dev &
FRONTEND_PID=$!

echo ""
echo "==================================="
echo "Server berjalan!"
echo "==================================="
echo "Backend:  http://localhost:3000"
echo "Frontend: http://localhost:5173"
echo ""
echo "Tekan Ctrl+C untuk stop kedua server"
echo ""

# Handle Ctrl+C
trap "kill $BACKEND_PID $FRONTEND_PID; exit" INT

wait
