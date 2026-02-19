import React, { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';

const Layout = () => {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const toggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  const getTitle = () => {
    const path = window.location.pathname;
    if (path === '/') return 'Dashboard';
    if (path === '/jabatan') return 'Manajemen Jabatan';
    if (path === '/karyawan') return 'Manajemen Karyawan';
    if (path === '/absensi') return 'Data Absensi';
    if (path === '/lembur') return 'Data Lembur';
    if (path === '/gaji') return 'Manajemen Gaji';
    if (path === '/laporan') return 'Laporan Gaji';
    return 'Sistem Payroll';
  };

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
