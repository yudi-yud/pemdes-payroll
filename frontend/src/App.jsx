import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import ErrorBoundary from './components/ErrorBoundary';
import PrivateRoute from './components/PrivateRoute';
import Layout from './components/layout/Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import JabatanList from './pages/Jabatan/JabatanList';
import KaryawanList from './pages/Karyawan/KaryawanList';
import GajiList from './pages/Gaji/GajiList';
import LaporanGaji from './pages/Laporan/LaporanGaji';
import AbsensiList from './pages/Absensi/AbsensiList';
import LemburList from './pages/Lembur/LemburList';

function App() {
  return (
    <ErrorBoundary>
      <AuthProvider>
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
              <Route path="gaji" element={<GajiList />} />
              <Route path="laporan" element={<LaporanGaji />} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ErrorBoundary>
  );
}

export default App;
