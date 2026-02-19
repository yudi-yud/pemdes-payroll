import React, { useEffect, useState } from 'react';
import { jabatanAPI, karyawanAPI, gajiAPI } from '../services/api';

const Dashboard = () => {
  const [stats, setStats] = useState({
    totalKaryawan: 0,
    totalJabatan: 0,
    totalGaji: 0,
    totalPending: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      const [karyawanRes, jabatanRes, gajiRes] = await Promise.all([
        karyawanAPI.getAll(),
        jabatanAPI.getAll(),
        gajiAPI.getAll(),
      ]);

      const currentMonth = new Date().getMonth() + 1;
      const currentYear = new Date().getFullYear();

      const totalGaji = gajiRes.data
        .filter(g => g.periode_bulan === currentMonth && g.periode_tahun === currentYear)
        .reduce((sum, g) => sum + g.total_gaji, 0);

      const totalPending = gajiRes.data.filter(g => g.status === 'pending').length;

      setStats({
        totalKaryawan: karyawanRes.data.length,
        totalJabatan: jabatanRes.data.length,
        totalGaji: totalGaji,
        totalPending: totalPending,
      });
    } catch (error) {
      console.error('Error fetching stats:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(amount);
  };

  const cards = [
    {
      title: 'Total Karyawan',
      value: stats.totalKaryawan,
      icon: '👥',
      color: 'bg-blue-500',
      textColor: 'text-blue-600',
    },
    {
      title: 'Total Jabatan',
      value: stats.totalJabatan,
      icon: '💼',
      color: 'bg-green-500',
      textColor: 'text-green-600',
    },
    {
      title: 'Total Gaji Bulan Ini',
      value: formatCurrency(stats.totalGaji),
      icon: '💰',
      color: 'bg-yellow-500',
      textColor: 'text-yellow-600',
    },
    {
      title: 'Gaji Pending',
      value: stats.totalPending,
      icon: '⏳',
      color: 'bg-red-500',
      textColor: 'text-red-600',
    },
  ];

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {cards.map((card, index) => (
          <div key={index} className="bg-white rounded-lg shadow-md p-6 border-l-4 border-blue-500 hover:shadow-lg transition-shadow">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-500 text-sm font-medium">{card.title}</p>
                <p className="text-2xl font-bold text-gray-800 mt-2">{card.value}</p>
              </div>
              <div className="text-4xl">{card.icon}</div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Selamat Datang</h3>
          <p className="text-gray-600">
            Selamat datang di Sistem Payroll Pemerintah Desa. Sistem ini digunakan
            untuk mengelola data jabatan, karyawan, penggajian, dan laporan.
          </p>
          <div className="mt-4 space-y-2">
            <p className="text-sm text-gray-500">📊 <strong>Dashboard</strong> - Ringkasan statistik</p>
            <p className="text-sm text-gray-500">💼 <strong>Jabatan</strong> - Kelola posisi/jabatan</p>
            <p className="text-sm text-gray-500">👥 <strong>Karyawan</strong> - Kelola data karyawan</p>
            <p className="text-sm text-gray-500">💰 <strong>Gaji</strong> - Kelola penggajian</p>
            <p className="text-sm text-gray-500">📈 <strong>Laporan</strong> - Lihat laporan gaji</p>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Cara Penggunaan</h3>
          <ol className="list-decimal list-inside space-y-2 text-gray-600">
            <li>Buat <strong>Jabatan</strong> terlebih dahulu dengan mengatur gaji pokok</li>
            <li>Tambahkan <strong>Karyawan</strong> dan assign ke jabatan</li>
            <li><strong>Generate Gaji</strong> per periode (bulanan)</li>
            <li>Lihat <strong>Laporan</strong> untuk rekapitulasi</li>
          </ol>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
