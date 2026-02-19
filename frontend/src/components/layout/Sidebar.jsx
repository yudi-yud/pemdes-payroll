import React from 'react';
import { NavLink } from 'react-router-dom';

const Sidebar = ({ collapsed = false, toggleSidebar }) => {
  const menuItems = [
    { path: '/', icon: '📊', label: 'Dashboard' },
    { path: '/jabatan', icon: '💼', label: 'Jabatan' },
    { path: '/karyawan', icon: '👥', label: 'Karyawan' },
    { path: '/absensi', icon: '📅', label: 'Absensi' },
    { path: '/lembur', icon: '⏰', label: 'Lembur' },
    { path: '/gaji', icon: '💰', label: 'Gaji' },
    { path: '/laporan', icon: '📈', label: 'Laporan' },
  ];

  return (
    <aside className={`bg-slate-800 min-h-screen fixed left-0 top-0 transition-all duration-300 z-50 ${collapsed ? 'w-20' : 'w-64'}`}>
      <div className={`p-6 ${collapsed ? 'px-4 text-center' : ''}`}>
        {!collapsed ? (
          <>
            <h1 className="text-xl font-bold text-white">Sistem Payroll</h1>
            <p className="text-slate-400 text-sm mt-1">Pemerintah Desa</p>
          </>
        ) : (
          <span className="text-2xl">💼</span>
        )}
      </div>

      <nav className="mt-6">
        {menuItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            end={item.path === '/'}
            className={({ isActive }) => `
              flex items-center ${collapsed ? 'justify-center px-4' : 'px-6'} py-3 text-slate-300 hover:bg-slate-700 hover:text-white transition-colors
              ${isActive ? 'bg-slate-900 text-white border-r-4 border-blue-500' : ''}
            `}
            title={collapsed ? item.label : ''}
          >
            <span className="text-xl ${collapsed ? '' : 'mr-3'}">{item.icon}</span>
            {!collapsed && <span className="font-medium">{item.label}</span>}
          </NavLink>
        ))}
      </nav>

      <div className={`absolute bottom-0 left-0 right-0 border-t border-slate-700 ${collapsed ? 'p-4 text-center' : 'p-6'}`}>
        {collapsed ? (
          <span className="text-xs text-slate-400">©24</span>
        ) : (
          <p className="text-slate-400 text-sm">© 2024 Pemdes Payroll</p>
        )}
      </div>
    </aside>
  );
};

export default Sidebar;
