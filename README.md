# PEMDES PAYROLL

Sistem Penggajian (Payroll) untuk Pemerintah Desa yang dibangun dengan arsitektur Fullstack modern.

## 🏗️ Teknologi

### Backend
- **Go** (Golang) - Bahasa pemrograman utama
- **Fiber** - Web framework berkinerja tinggi
- **GORM** - ORM untuk database MySQL
- **JWT** - Autentikasi token-based

### Frontend
- **React 18** - UI Library
- **Vite** - Build tool & dev server
- **Tailwind CSS** - Utility-first CSS framework
- **React Router** - Client-side routing
- **Axios** - HTTP client
- **date-fns** - Library manipulasi tanggal

### Database & Deployment
- **MySQL** - Database utama
- **Docker** - Containerization
- **Nginx** - Reverse Proxy untuk production

## 📋 Fitur

- ✅ **Manajemen Jabatan** - Kelola data jabatan dan gaji pokok
- ✅ **Manajemen Karyawan** - CRUD data karyawan desa
- ✅ **Manajemen Gaji** - Hitung dan kelola penggajian
- ✅ **Absensi** - Catat kehadiran karyawan
- ✅ **Lembur** - Kelola data lembur dan kompensasi
- ✅ **Laporan** - Generate laporan penggajian
- ✅ **Autentikasi** - Login user dengan JWT

## 🚀 Quick Start

### Prasyarat

- Go 1.21 atau lebih tinggi
- Node.js 18 atau lebih tinggi
- MySQL 8.0 atau lebih tinggi
- Docker (opsional, untuk deployment)

### 1. Clone Repository

```bash
git clone <repository-url>
cd pemdes-payroll
```

### 2. Setup Database

```bash
mysql -u root -p < setup.sql
```

### 3. Setup Environment

Salin file environment example:

```bash
cp .env.example .env
```

Edit `.env` sesuai konfigurasi Anda:

```env
# Database Configuration
MYSQL_ROOT_PASSWORD=susahbanget
MYSQL_DATABASE=pemdes_payroll
MYSQL_USER=payroll_user
MYSQL_PASSWORD=susahbanget
MYSQL_PORT=3306

# Backend Configuration
BACKEND_PORT=3000
JWT_SECRET=RIJNrA&ZNsZN16k.Ep0-RIJNrA&ZNsZN16k.Ep0

# Frontend Configuration
FRONTEND_PORT=80

# Nginx Reverse Proxy (Production)
NGINX_PORT=8080
```

### 4. Install Backend Dependencies

```bash
go mod download
```

### 5. Install Frontend Dependencies

```bash
cd frontend
npm install
cd ..
```

### 6. Jalankan Aplikasi

#### Development Mode (Windows)

```bash
start.bat
```

#### Development Mode (Linux/Mac)

```bash
chmod +x start.sh
./start.sh
```

#### Manual

Terminal 1 - Backend:
```bash
go run main.go
```

Terminal 2 - Frontend:
```bash
cd frontend
npm run dev
```

### 7. Akses Aplikasi

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:3000
- **Default Admin**:
  - Username: `admin`
  - Password: `admin123`

## 📁 Struktur Proyek

```
pemdes-payroll/
├── backend/
│   ├── config/         # Konfigurasi database
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # Custom middleware (CORS, Auth)
│   ├── models/         # Data models
│   ├── repositories/   # Database access layer
│   └── routes/         # Route definitions
├── frontend/
│   ├── src/
│   │   ├── components/ # Reusable components
│   │   │   └── ui/     # UI components (Button, Table, Modal)
│   │   ├── pages/      # Page components
│   │   │   ├── Dashboard/
│   │   │   ├── Karyawan/
│   │   │   ├── Gaji/
│   │   │   └── ...
│   │   ├── main.jsx    # App entry point
│   │   └── index.css   # Global styles
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
├── main.go             # Backend entry point
├── setup.sql           # Database setup script
├── docker-compose.yml  # Docker configuration
├── nginx.conf          # Nginx configuration
├── start.bat           # Windows startup script
└── start.sh            # Linux/Mac startup script
```

## 🔧 API Endpoints

### Health Check
- `GET /health` - Cek status API

### Autentikasi
- `POST /api/auth/login` - Login user
- `POST /api/auth/logout` - Logout user

### Jabatan
- `GET /api/jabatan` - List semua jabatan
- `GET /api/jabatan/:id` - Detail jabatan
- `POST /api/jabatan` - Tambah jabatan
- `PUT /api/jabatan/:id` - Update jabatan
- `DELETE /api/jabatan/:id` - Hapus jabatan

### Karyawan
- `GET /api/karyawan` - List semua karyawan
- `GET /api/karyawan/:id` - Detail karyawan
- `POST /api/karyawan` - Tambah karyawan
- `PUT /api/karyawan/:id` - Update karyawan
- `DELETE /api/karyawan/:id` - Hapus karyawan

### Gaji
- `GET /api/gaji` - List semua gaji
- `GET /api/gaji/:id` - Detail gaji
- `POST /api/gaji` - Hitung & simpan gaji
- `PUT /api/gaji/:id` - Update gaji
- `DELETE /api/gaji/:id` - Hapus gaji

### Absensi
- `GET /api/absensi` - List absensi
- `POST /api/absensi` - Tambah absensi
- `PUT /api/absensi/:id` - Update absensi

### Lembur
- `GET /api/lembur` - List lembur
- `POST /api/lembur` - Tambah lembur

### Laporan
- `GET /api/laporan/gaji` - Laporan gaji

## 🐳 Docker Deployment

```bash
# Build dan start semua container
docker-compose up -d

# Cek logs
docker-compose logs -f

# Stop container
docker-compose down

# Stop dan hapus volume database
docker-compose down -v
```

## 📝 Environment Variables

| Variable | Deskripsi | Default |
|----------|-----------|---------|
| `MYSQL_ROOT_PASSWORD` | Password root MySQL | - |
| `MYSQL_DATABASE` | Nama database | `pemdes_payroll` |
| `MYSQL_USER` | User database | `payroll_user` |
| `MYSQL_PASSWORD` | Password database | - |
| `MYSQL_PORT` | Port MySQL | `3306` |
| `BACKEND_PORT` | Port backend API | `3000` |
| `JWT_SECRET` | Secret key untuk JWT | - |
| `FRONTEND_PORT` | Port frontend dev | `5173` |
| `NGINX_PORT` | Port Nginx production | `8080` |

## 🔐 Default Credentials

**Admin User:**
- Username: `admin`
- Password: `admin123`

> **Penting**: Ganti password default di production!

## 📦 Building Frontend untuk Production

```bash
cd frontend
npm run build
```

File build akan tersimpan di `frontend/dist/`

## 🛠️ Troubleshooting

### Port sudah digunakan
Jika port 3000 atau 5173 sudah digunakan, ubah di file `.env`:

```env
BACKEND_PORT=3001
FRONTEND_PORT=5174
```

### Database connection failed
Pastikan MySQL sudah running dan credentials di `.env` sudah benar.

### CORS Error
Backend sudah dikonfigurasi dengan CORS middleware. Pastikan origin frontend sudah terdaftar.

## 📄 Lisensi

[MIT License](LICENSE)

## 👨‍💻 Development

Dikembangkan untuk keperluan Pemerintah Desa dalam pengelolaan penggajian pegawai.

---

**Version**: 1.0.0
**Last Updated**: 2026
