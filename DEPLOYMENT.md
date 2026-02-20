# Dokumentasi Deploy VPS - Pemdes Payroll

Panduan lengkap deploy aplikasi Pemdes Payroll ke VPS menggunakan Docker.

---

## Persyaratan

### VPS Minimum Specifications
- **OS**: Ubuntu 20.04 / 22.04 atau Debian 11+
- **RAM**: 2 GB minimum (4 GB recommended)
- **Storage**: 20 GB minimum
- **CPU**: 2 core minimum

### Software yang perlu diinstall di VPS
- Docker 20.10+
- Docker Compose 2.0+

---

## 1. Persiapan VPS

### 1.1 Update System
```bash
sudo apt update && sudo apt upgrade -y
```

### 1.2 Install Docker
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Enable dan start Docker
sudo systemctl enable docker
sudo systemctl start docker

# Verifikasi instalasi
docker --version
```

### 1.3 Install Docker Compose
```bash
# Download Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# Beri permission execute
sudo chmod +x /usr/local/bin/docker-compose

# Verifikasi instalasi
docker-compose --version
```

### 1.4 Setup Firewall (Optional tapi Recommended)
```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP & HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable
```

---

## 2. Deploy Aplikasi

### 2.1 Clone Repository
```bash
# Clone dari Git
cd /var
sudo mkdir -p www
cd www
sudo git clone <your-repo-url> pemdes-payroll
cd pemdes-payroll

# Atau upload manual menggunakan SCP/rsync
# scp -r ./pemdes-payroll user@your-vps-ip:/var/www/
```

### 2.2 Konfigurasi Environment
```bash
# Copy file environment example
cp .env.example .env

# Edit file .env sesuai kebutuhan
nano .env
```

**PENTING: Ubah nilai berikut di `.env`:**
```env
MYSQL_ROOT_PASSWORD=ganti-dengan-password-kuat
MYSQL_PASSWORD=ganti-dengan-password-kuat
JWT_SECRET=ganti-dengan-secret-key-panjang-dan-acak
```

### 2.3 Pastikan Backend Entry Point Ada
```bash
# Cek apakah main.go ada di root directory
ls main.go

# File main.go harus ada di root project pemdes-payroll
# Jika tidak ada, pastikan Anda clone repository dengan lengkap
```

### 2.4 Build dan Start Container
```bash
# Build images
docker-compose build

# Start semua container (background)
docker-compose up -d

# Cek status container
docker-compose ps

# Cek logs
docker-compose logs -f
```

---

## 3. Verifikasi Deploy

### 3.1 Cek Container Status
```bash
docker-compose ps
```

Expected output:
```
NAME                  STATUS              PORTS
payroll-backend       Up                  0.0.0.0:3000->3000/tcp
payroll-frontend      Up                  0.0.0.0:80->80/tcp
payroll-mysql         Up (healthy)        0.0.0.0:3306->3306/tcp
```

### 3.2 Cek Logs
```bash
# Backend logs
docker-compose logs backend

# Frontend logs
docker-compose logs frontend

# MySQL logs
docker-compose logs mysql
```

### 3.3 Test Aplikasi
Buka browser dan akses:
- Frontend: `http://your-vps-ip`
- API: `http://your-vps-ip/api/health` (jika ada)

---

## 4. Konfigurasi Domain & SSL (Optional)

### 4.1 Setup Domain
Tambahkan A record di DNS provider:
```
Type: A
Name: @ (atau www)
Value: your-vps-ip
TTL: 3600
```

### 4.2 Install SSL dengan Certbot
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Request SSL certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal sudah otomatis dikonfigurasi
```

Update `nginx.conf` untuk HTTPS:
```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    # ... konfigurasi lainnya
}
```

---

## 5. Manajemen Container

### Perintah Dasar
```bash
# Stop semua container
docker-compose stop

# Start semua container
docker-compose start

# Restart semua container
docker-compose restart

# Stop dan hapus container
docker-compose down

# Stop dan hapus termasuk volumes
docker-compose down -v

# Rebuild setelah code changes
docker-compose up -d --build
```

### Update Aplikasi
```bash
# 1. Pull latest code
cd /var/www/pemdes-payroll
git pull origin main

# 2. Rebuild dan restart
docker-compose up -d --build

# 3. Hapus images lama (optional)
docker image prune -a
```

### Backup Database
```bash
# Backup
docker-compose exec mysql mysqldump -u root -p pemdes_payroll > backup_$(date +%Y%m%d).sql

# Restore
docker-compose exec -T mysql mysql -u root -p pemdes_payroll < backup_20240101.sql
```

---

## 6. Troubleshooting

### Container tidak start
```bash
# Cek logs
docker-compose logs -f [service-name]

# Cek resource usage
docker stats

# Restart spesifik service
docker-compose restart [service-name]
```

### Database connection error
```bash
# Cek MySQL container status
docker-compose logs mysql

# Masuk ke MySQL container
docker-compose exec mysql mysql -u root -p

# Di dalam MySQL
SHOW DATABASES;
SELECT User, Host FROM mysql.user;
```

### Port sudah digunakan
```bash
# Cek port yang digunakan
sudo netstat -tulpn | grep :3000

# Kill proses yang menggunakan port
sudo kill -9 [PID]
```

### Out of memory
```bash
# Cek memory usage
free -h

# Add swap space (jika RAM kurang)
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

---

## 7. Security Hardening

### 7.1 SSH Configuration
```bash
# Edit SSH config
sudo nano /etc/ssh/sshd_config

# Ubah setting berikut:
Port 2222                    # Ganti default port
PermitRootLogin no           # Disable root login
PasswordAuthentication no    # Hanya key-based auth

# Restart SSH
sudo systemctl restart sshd
```

### 7.2 Firewall Setup
```bash
# Hapus aturan lama dan setup baru
sudo ufw reset

# Default deny
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH (custom port)
sudo ufw allow 2222/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable
sudo ufw enable
```

### 7.3 Fail2Ban untuk Brute Force Protection
```bash
# Install fail2ban
sudo apt install fail2ban -y

# Setup jail untuk SSH
sudo nano /etc/fail2ban/jail.local
```

```
[sshd]
enabled = true
port = 2222
maxretry = 3
bantime = 3600
findtime = 600
```

```bash
# Restart fail2ban
sudo systemctl restart fail2ban
```

---

## 8. Monitoring & Logs

### Cek Container Logs
```bash
# Real-time logs
docker-compose logs -f

# Logs spesifik service
docker-compose logs -f backend

# 100 baris terakhir
docker-compose logs --tail=100 backend
```

### Monitoring dengan Docker Stats
```bash
# Resource usage real-time
docker stats

# Spesifik container
docker stats payroll-backend
```

---

## 9. File Structure Hasil Deploy

```
/var/www/pemdes-payroll/
├── main.go                   # Entry point backend
├── go.mod
├── go.sum
├── backend/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repositories/
│   ├── routes/
│   ├── services/
│   └── middleware/
├── frontend/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── nginx.conf
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── ...
├── docker-compose.yml
├── nginx.conf                # Reverse proxy config
├── .env                      # Environment variables
└── .env.example
```

---

## 10. Kontak & Support

Jika mengalami masalah:
1. Cek logs: `docker-compose logs -f`
2. Verifikasi config di `.env`
3. Pastikan semua ports tersedia
4. Cek resource VPS (RAM/CPU)
