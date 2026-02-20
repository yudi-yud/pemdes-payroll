import React, { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import Sidebar from './Sidebar';
import Header from './Header';

const Layout = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const isKaryawan = user?.role === 'karyawan';

  const toggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  const getTitle = () => {
    const path = location.pathname;
    if (path === '/') return 'Dashboard';
    if (path === '/jabatan') return 'Manajemen Jabatan';
    if (path === '/karyawan') return 'Manajemen Karyawan';
    if (path === '/absensi') return 'Data Absensi';
    if (path === '/lembur') return 'Data Lembur';
    if (path === '/gaji') return 'Manajemen Gaji';
    if (path === '/laporan') return 'Laporan Gaji';
    if (path === '/users') return 'Manajemen User';
    if (path === '/slip-gaji') return 'Slip Gaji Saya';
    return 'Sistem Payroll';
  };

  // Redirect karyawan ke slip-gaji jika bukan di halaman slip-gaji
  useEffect(() => {
    if (isKaryawan && location.pathname !== '/slip-gaji') {
      navigate('/slip-gaji', { replace: true });
    }
  }, [isKaryawan, location.pathname, navigate]);

  // Jika role karyawan, gunakan layout sederhana
  if (isKaryawan) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Header title={getTitle()} toggleSidebar={toggleSidebar} />
        <main className="p-8">
          <Outlet />
        </main>
      </div>
    );
  }

  return (
    <div className="flex">
      <Sidebar collapsed={sidebarCollapsed} toggleSidebar={toggleSidebar} />
      <div className={`flex-1 min-h-screen bg-gray-50 transition-all duration-300 ${sidebarCollapsed ? 'ml-20' : 'ml-64'}`}>
        <Header title={getTitle()} toggleSidebar={toggleSidebar} />
        <main className="p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;
