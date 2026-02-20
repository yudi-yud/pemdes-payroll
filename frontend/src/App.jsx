import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { AuthTimeoutProvider } from './contexts/AuthTimeoutContext';
import ErrorBoundary from './components/ErrorBoundary';
import PrivateRoute from './components/PrivateRoute';
import ProtectedRoute from './components/ProtectedRoute';
import Layout from './components/layout/Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import JabatanList from './pages/Jabatan/JabatanList';
import KaryawanList from './pages/Karyawan/KaryawanList';
import GajiList from './pages/Gaji/GajiList';
import LaporanGaji from './pages/Laporan/LaporanGaji';
import AbsensiList from './pages/Absensi/AbsensiList';
import LemburList from './pages/Lembur/LemburList';
import UserManagement from './pages/Users/UserManagement';
import SlipGaji from './pages/Portal/SlipGaji';

function App() {
  return (
    <ErrorBoundary>
      <AuthProvider>
        <AuthTimeoutProvider>
          <BrowserRouter>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route
                path="/"
                element={
                  <PrivateRoute>
                    <Layout />
                  </PrivateRoute>
                }
              >
                <Route index element={<Dashboard />} />
                <Route path="jabatan" element={<JabatanList />} />
                <Route path="karyawan" element={<KaryawanList />} />
                <Route path="absensi" element={<AbsensiList />} />
                <Route path="lembur" element={<LemburList />} />
                <Route
                  path="gaji"
                  element={
                    <ProtectedRoute allowedRoles={['admin', 'finance']}>
                      <GajiList />
                    </ProtectedRoute>
                  }
                />
                <Route path="laporan" element={<LaporanGaji />} />
                <Route
                  path="users"
                  element={
                    <ProtectedRoute allowedRoles={['admin']}>
                      <UserManagement />
                    </ProtectedRoute>
                  }
                />
                <Route path="slip-gaji" element={<SlipGaji />} />
                <Route path="*" element={<Navigate to="/" replace />} />
              </Route>
            </Routes>
          </BrowserRouter>
        </AuthTimeoutProvider>
      </AuthProvider>
    </ErrorBoundary>
  );
}

export default App;
