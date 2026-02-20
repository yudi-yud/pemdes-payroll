# File Deploy Pemdes Payroll

Berikut adalah file-file yang telah dibuat untuk deployment ke VPS:

## 📁 Struktur File

```
pemdes-payroll/
├── backend/
│   ├── Dockerfile              # Build image backend (Go)
│   └── .dockerignore           # File yang diignore saat build
├── frontend/
│   ├── Dockerfile              # Build image frontend (React + Nginx)
│   ├── .dockerignore           # File yang diignore saat build
│   └── nginx.conf              # Config nginx untuk serving static files
├── docker-compose.yml          # Orchestration semua services
├── nginx.conf                  # Reverse proxy config (production)
├── .env.example                # Template environment variables
├── deploy.sh                   # Script deploy otomatis
└── DEPLOYMENT.md               # Dokumentasi lengkap
```

## 🚀 Cara Deploy

### Option 1: Quick Deploy (Recommended)

```bash
# 1. Upload ke VPS
scp -r ./pemdes-payroll user@your-vps-ip:/var/www/

# 2. SSH ke VPS
ssh user@your-vps-ip

# 3. Jalankan script deploy
cd /var/www/pemdes-payroll
chmod +x deploy.sh
sudo ./deploy.sh
```

### Option 2: Manual Deploy

Lihat panduan lengkap di `DEPLOYMENT.md`

## ⚙️ Konfigurasi

### Environment Variables (.env)

```env
# Database
MYSQL_ROOT_PASSWORD=strong_password_here
MYSQL_DATABASE=pemdes_payroll
MYSQL_USER=payroll_user
MYSQL_PASSWORD=strong_password_here

# Backend
BACKEND_PORT=3000
JWT_SECRET=your-super-secret-jwt-key

# Frontend
FRONTEND_PORT=80
```

## 📋 Ports yang Digunakan

| Service | Internal Port | External Port |
|---------|---------------|---------------|
| Frontend | 80 | 80 |
| Backend | 3000 | 3000 |
| MySQL | 3306 | 3306 |
| Nginx Proxy | 80 | 8080 |

## 🔧 Perintah Penting

```bash
# Build & Start
docker-compose up -d --build

# Stop
docker-compose stop

# Restart
docker-compose restart

# Logs
docker-compose logs -f

# Backup Database
docker-compose exec mysql mysqldump -u root -p pemdes_payroll > backup.sql

# Update Application
git pull origin main
docker-compose up -d --build
```

## 🔐 Security Notes

1. **SELALU** ganti password di `.env` sebelum deploy
2. Setup firewall untuk allow hanya port yang diperlukan
3. Gunakan HTTPS di production dengan Let's Encrypt
4. Disable root login SSH
5. Setup fail2ban untuk brute-force protection

## 📞 Support

Lihat `DEPLOYMENT.md` untuk troubleshooting dan dokumentasi lengkap.
