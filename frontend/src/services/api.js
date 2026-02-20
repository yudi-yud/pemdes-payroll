import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to include token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Jabatan API
export const jabatanAPI = {
  getAll: () => api.get('/api/jabatan'),
  getById: (id) => api.get(`/api/jabatan/${id}`),
  create: (data) => api.post('/api/jabatan', data),
  update: (id, data) => api.put(`/api/jabatan/${id}`, data),
  delete: (id) => api.delete(`/api/jabatan/${id}`),
};

// Karyawan API
export const karyawanAPI = {
  getAll: () => api.get('/api/karyawan'),
  getById: (id) => api.get(`/api/karyawan/${id}`),
  search: (q) => api.get(`/api/karyawan/search?q=${q}`),
  getByStatus: (status) => api.get(`/api/karyawan/status/${status}`),
  create: (data) => api.post('/api/karyawan', data),
  update: (id, data) => api.put(`/api/karyawan/${id}`, data),
  delete: (id) => api.delete(`/api/karyawan/${id}`),
};

// Gaji API
export const gajiAPI = {
  getAll: () => api.get('/api/gaji'),
  getById: (id) => api.get(`/api/gaji/${id}`),
  getByPeriod: (bulan, tahun) => api.get(`/api/gaji/period?bulan=${bulan}&tahun=${tahun}`),
  getByKaryawanId: (id) => api.get(`/api/gaji/karyawan/${id}`),
  create: (data) => api.post('/api/gaji', data),
  update: (id, data) => api.put(`/api/gaji/${id}`, data),
  delete: (id) => api.delete(`/api/gaji/${id}`),
  updateStatus: (id, status) => api.patch(`/api/gaji/${id}/status`, { status }),
  generateBatch: (data) => api.post('/api/gaji/generate-batch', data),
};

// Laporan API
export const laporanAPI = {
  getLaporanGaji: (bulan, tahun) => api.get(`/api/laporan/gaji?bulan=${bulan}&tahun=${tahun}`),
  getRiwayatGajiKaryawan: (id) => api.get(`/api/laporan/gaji/karyawan/${id}`),
  getRekapGaji: (bulan, tahun) => api.get(`/api/laporan/rekap?bulan=${bulan}&tahun=${tahun}`),
};

// Absensi API
export const absensiAPI = {
  getAll: () => api.get('/api/absensi'),
  getById: (id) => api.get(`/api/absensi/${id}`),
  getByKaryawan: (id, startDate, endDate) => api.get(`/api/absensi/karyawan/${id}?start_date=${startDate}&end_date=${endDate}`),
  getRekap: (karyawanId, bulan, tahun) => api.get(`/api/absensi/rekap/${karyawanId}?bulan=${bulan}&tahun=${tahun}`),
  create: (data) => api.post('/api/absensi', data),
  update: (id, data) => api.put(`/api/absensi/${id}`, data),
  delete: (id) => api.delete(`/api/absensi/${id}`),
  exportExcel: () => {
    window.location.href = `${API_BASE_URL}/api/absensi/export/excel`;
  },
  exportPDF: (karyawanId, bulan, tahun) => {
    window.location.href = `${API_BASE_URL}/api/absensi/export/karyawan/${karyawanId}/pdf?bulan=${bulan}&tahun=${tahun}`;
  },
};

// Lembur API
export const lemburAPI = {
  getAll: () => api.get('/api/lembur'),
  getById: (id) => api.get(`/api/lembur/${id}`),
  getByPeriod: (bulan, tahun) => api.get(`/api/lembur/period?bulan=${bulan}&tahun=${tahun}`),
  getByKaryawan: (id) => api.get(`/api/lembur/karyawan/${id}`),
  create: (data) => api.post('/api/lembur', data),
  update: (id, data) => api.put(`/api/lembur/${id}`, data),
  delete: (id) => api.delete(`/api/lembur/${id}`),
  approve: (id, status, approverId = null) => api.patch(`/api/lembur/${id}/approve`, { status, approver_id: approverId }),
};

// Health check
export const healthAPI = {
  check: () => api.get('/health'),
};

// User Management API
export const userAPI = {
  getAll: () => api.get('/api/auth/users'),
  create: (data) => api.post('/api/auth/users', data),
  update: (id, data) => api.put(`/api/auth/users/${id}`, data),
  delete: (id) => api.delete(`/api/auth/users/${id}`),
  toggleActive: (id) => api.patch(`/api/auth/users/${id}/toggle`),
  changePassword: (data) => api.put('/api/auth/change-password', data),
};

// Employee Portal API
export const portalAPI = {
  getMySlips: (bulan, tahun) => api.get(`/api/gaji/my-slips?bulan=${bulan}&tahun=${tahun}`),
  getSlip: (id) => api.get(`/api/gaji/slip/${id}`),
};

export default api;
