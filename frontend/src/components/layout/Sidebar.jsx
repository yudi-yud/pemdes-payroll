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
    { path: '/users', icon: '👤', label: 'Users' },
    { path: '/slip-gaji', icon: '📄', label: 'Slip Gaji' },
  ];

  return (
    <aside
      className={`bg-slate-800 h-screen fixed left-0 top-0 flex flex-col transition-all duration-300 z-50 ${
        collapsed ? 'w-20' : 'w-64'
      }`}
      style={{
        scrollbarWidth: 'thin',
        scrollbarColor: '#475569 #1e293b',
      }}
    >
      {/* Header - Fixed */}
      <div className={`flex-shrink-0 border-b border-slate-700 ${collapsed ? 'p-4 text-center' : 'p-6'}`}>
        {!collapsed ? (
          <>
            <h1 className="text-xl font-bold text-white">Sistem Payroll</h1>
            <p className="text-slate-400 text-sm mt-1">Pemerintah Desa</p>
          </>
        ) : (
          <span className="text-2xl">💼</span>
        )}
      </div>

      {/* Nav - Scrollable */}
      <nav
        className="flex-1 overflow-y-auto overflow-x-hidden mt-4"
        style={{
          scrollbarWidth: 'thin',
          scrollbarColor: '#475569 #1e293b',
        }}
      >
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
            <span className={`text-xl flex-shrink-0 ${collapsed ? '' : 'mr-3'}`}>{item.icon}</span>
            {!collapsed && <span className="font-medium whitespace-nowrap">{item.label}</span>}
          </NavLink>
        ))}
      </nav>

      {/* Footer - Fixed */}
      <div className={`flex-shrink-0 border-t border-slate-700 ${collapsed ? 'p-4 text-center' : 'p-6'}`}>
        {collapsed ? (
          <span className="text-xs text-slate-400">©24</span>
        ) : (
          <p className="text-slate-400 text-sm whitespace-nowrap">© 2024 Pemdes Payroll</p>
        )}
      </div>

      {/* Custom Scrollbar for Webkit browsers */}
      <style>{`
        aside nav::-webkit-scrollbar {
          width: 6px;
        }
        aside nav::-webkit-scrollbar-track {
          background: #1e293b;
        }
        aside nav::-webkit-scrollbar-thumb {
          background: #475569;
          border-radius: 3px;
        }
        aside nav::-webkit-scrollbar-thumb:hover {
          background: #64748b;
        }
      `}</style>
    </aside>
  );
};

export default Sidebar;
