import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { absensiAPI, karyawanAPI } from '../../services/api';
import Button from '../../components/ui/Button';
import Modal from '../../components/ui/Modal';
import Input from '../../components/ui/Input';

const AbsensiList = () => {
  const [absensi, setAbsensi] = useState([]);
  const [karyawan, setKaryawan] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [viewMode, setViewMode] = useState('list'); // 'list' or 'rekap'
  const [selectedKaryawan, setSelectedKaryawan] = useState('');
  const [rekapData, setRekapData] = useState(null);
  const [currentPeriod, setCurrentPeriod] = useState({
    bulan: new Date().getMonth() + 1,
    tahun: new Date().getFullYear(),
  });
  const [editingAbsensi, setEditingAbsensi] = useState(null);
  const [formData, setFormData] = useState({
    karyawan_id: '',
    tanggal: '',
    jam_masuk: '',
    jam_keluar: '',
    status: 'hadir',
    keterangan: '',
  });

  useEffect(() => {
    fetchKaryawan();
  }, []);

  const fetchKaryawan = async () => {
    try {
      const response = await karyawanAPI.getByStatus('aktif');
      setKaryawan(response.data);
    } catch (error) {
      console.error('Error fetching karyawan:', error);
    }
  };

  const fetchAbsensi = async () => {
    setLoading(true);
    try {
      const response = await absensiAPI.getAll();
      setAbsensi(response.data);
    } catch (error) {
      console.error('Error fetching absensi:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchRekap = async () => {
    if (!selectedKaryawan) return;
    setLoading(true);
    try {
      const response = await absensiAPI.getRekap(
        selectedKaryawan,
        currentPeriod.bulan,
        currentPeriod.tahun
      );
      setRekapData(response.data);
    } catch (error) {
      console.error('Error fetching rekap:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (viewMode === 'list') {
      fetchAbsensi();
    } else {
      fetchRekap();
    }
  }, [viewMode, selectedKaryawan, currentPeriod]);

  const handleOpenModal = (absensi = null) => {
    if (absensi) {
      setEditingAbsensi(absensi);
      setFormData({
        karyawan_id: absensi.karyawan_id.toString(),
        tanggal: absensi.tanggal ? absensi.tanggal.split('T')[0] : '',
        jam_masuk: absensi.jam_masuk || '',
        jam_keluar: absensi.jam_keluar || '',
        status: absensi.status || 'hadir',
        keterangan: absensi.keterangan || '',
      });
    } else {
      setEditingAbsensi(null);
      setFormData({
        karyawan_id: '',
        tanggal: new Date().toISOString().split('T')[0],
        jam_masuk: '',
        jam_keluar: '',
        status: 'hadir',
        keterangan: '',
      });
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingAbsensi(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const data = {
        ...formData,
        karyawan_id: parseInt(formData.karyawan_id),
      };

      if (editingAbsensi) {
        await axios.put(`http://localhost:3000/api/absensi/${editingAbsensi.id}`, data);
      } else {
        await absensiAPI.create(data);
      }

      await fetchAbsensi();
      handleCloseModal();
    } catch (error) {
      console.error('Error saving absensi:', error);
      alert('Gagal menyimpan absensi: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Apakah Anda yakin ingin menghapus data absensi ini?')) return;

    try {
      await absensiAPI.delete(id);
      await fetchAbsensi();
    } catch (error) {
      console.error('Error deleting absensi:', error);
      alert('Gagal menghapus absensi: ' + (error.response?.data?.error || error.message));
    }
  };

  const bulanOptions = [
    { value: 1, label: 'Januari' },
    { value: 2, label: 'Februari' },
    { value: 3, label: 'Maret' },
    { value: 4, label: 'April' },
    { value: 5, label: 'Mei' },
    { value: 6, label: 'Juni' },
    { value: 7, label: 'Juli' },
    { value: 8, label: 'Agustus' },
    { value: 9, label: 'September' },
    { value: 10, label: 'Oktober' },
    { value: 11, label: 'November' },
    { value: 12, label: 'Desember' },
  ];

  const statusOptions = [
    { value: 'hadir', label: 'Hadir', color: 'bg-green-100 text-green-800' },
    { value: 'izin', label: 'Izin', color: 'bg-blue-100 text-blue-800' },
    { value: 'sakit', label: 'Sakit', color: 'bg-yellow-100 text-yellow-800' },
    { value: 'alpha', label: 'Alpha', color: 'bg-red-100 text-red-800' },
  ];

  const getStatusBadge = (status) => {
    const s = statusOptions.find(opt => opt.value === status);
    return s ? (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${s.color}`}>
        {s.label}
      </span>
    ) : status;
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Data Absensi</h1>
        <div className="flex gap-3">
          <div className="flex items-center gap-2 bg-white rounded-lg shadow p-2">
            <button
              onClick={() => setViewMode('list')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${viewMode === 'list' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100'}`}
            >
              📋 List Absensi
            </button>
            <button
              onClick={() => setViewMode('rekap')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${viewMode === 'rekap' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100'}`}
            >
              📊 Rekap Bulanan
            </button>
          </div>
          {viewMode === 'list' ? (
            <>
              <Button onClick={() => absensiAPI.exportExcel()} variant="success">📊 Export Excel</Button>
              <Button onClick={() => handleOpenModal()}>+ Tambah Absensi</Button>
            </>
          ) : selectedKaryawan ? (
            <Button onClick={() => absensiAPI.exportPDF(selectedKaryawan, currentPeriod.bulan, currentPeriod.tahun)} variant="success">📄 Export PDF</Button>
          ) : null}
        </div>
      </div>

      {viewMode === 'rekap' && (
        <div className="bg-white rounded-lg shadow-md p-4 mb-6">
          <div className="flex gap-4 items-end">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-700 mb-1">Pilih Karyawan</label>
              <select
                value={selectedKaryawan}
                onChange={(e) => setSelectedKaryawan(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">-- Pilih Karyawan --</option>
                {karyawan.map((k) => (
                  <option key={k.id} value={k.id}>{k.nama}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Bulan</label>
              <select
                value={currentPeriod.bulan}
                onChange={(e) => setCurrentPeriod({ ...currentPeriod, bulan: parseInt(e.target.value) })}
                className="w-32 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {bulanOptions.map(b => (
                  <option key={b.value} value={b.value}>{b.label}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Tahun</label>
              <input
                type="number"
                value={currentPeriod.tahun}
                onChange={(e) => setCurrentPeriod({ ...currentPeriod, tahun: parseInt(e.target.value) })}
                className="w-24 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                min="2020"
                max="2100"
              />
            </div>
          </div>

          {rekapData && (
            <div className="mt-6 grid grid-cols-2 md:grid-cols-5 gap-4">
              <div className="bg-green-50 rounded-lg p-4">
                <p className="text-sm text-gray-600">Hadir</p>
                <p className="text-2xl font-bold text-green-600">{rekapData.hadir || 0}</p>
              </div>
              <div className="bg-blue-50 rounded-lg p-4">
                <p className="text-sm text-gray-600">Izin</p>
                <p className="text-2xl font-bold text-blue-600">{rekapData.izin || 0}</p>
              </div>
              <div className="bg-yellow-50 rounded-lg p-4">
                <p className="text-sm text-gray-600">Sakit</p>
                <p className="text-2xl font-bold text-yellow-600">{rekapData.sakit || 0}</p>
              </div>
              <div className="bg-red-50 rounded-lg p-4">
                <p className="text-sm text-gray-600">Alpha</p>
                <p className="text-2xl font-bold text-red-600">{rekapData.alpha || 0}</p>
              </div>
              <div className="bg-purple-50 rounded-lg p-4">
                <p className="text-sm text-gray-600">Total Hari</p>
                <p className="text-2xl font-bold text-purple-600">
                  {(rekapData.hadir || 0) + (rekapData.izin || 0) + (rekapData.sakit || 0) + (rekapData.alpha || 0)}
                </p>
              </div>
            </div>
          )}
        </div>
      )}

      {viewMode === 'list' ? (
        loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        ) : (
          <div className="bg-white rounded-lg shadow-md overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Karyawan</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Jam Masuk</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Jam Keluar</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Keterangan</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {absensi.length === 0 ? (
                  <tr>
                    <td colSpan={8} className="px-6 py-4 text-center text-gray-500">
                      Tidak ada data absensi
                    </td>
                  </tr>
                ) : (
                  absensi.map((row) => (
                    <tr key={row.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{row.tanggal}</td>
                      <td className="px-6 py-4 text-sm text-gray-900">{row.karyawan?.nama}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{row.jam_masuk}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{row.jam_keluar}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm">{getStatusBadge(row.status)}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">{row.keterangan}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm">
                        <div className="flex gap-2">
                          <button onClick={() => handleOpenModal(row)} className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700">Edit</button>
                          <button onClick={() => handleDelete(row.id)} className="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700">Hapus</button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        )
      ) : (
        <div className="bg-white rounded-lg shadow-md p-6">
          {!selectedKaryawan ? (
            <p className="text-center text-gray-500">Silakan pilih karyawan untuk melihat rekapitulasi absensi</p>
          ) : loading ? (
            <div className="flex items-center justify-center h-64">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            </div>
          ) : !rekapData ? (
            <p className="text-center text-gray-500">Belum ada data absensi untuk periode ini</p>
          ) : null}
        </div>
      )}

      <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingAbsensi ? 'Edit Absensi' : 'Tambah Absensi'}>
        <form onSubmit={handleSubmit}>
          <Input
            label="Karyawan"
            type="select"
            value={formData.karyawan_id}
            onChange={(e) => setFormData({ ...formData, karyawan_id: e.target.value })}
            required
            options={karyawan.map(k => ({ value: k.id, label: k.nama }))}
          />
          <Input
            label="Tanggal"
            type="date"
            value={formData.tanggal}
            onChange={(e) => setFormData({ ...formData, tanggal: e.target.value })}
            required
          />
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Jam Masuk"
              value={formData.jam_masuk}
              onChange={(e) => setFormData({ ...formData, jam_masuk: e.target.value })}
              placeholder="08:00"
            />
            <Input
              label="Jam Keluar"
              value={formData.jam_keluar}
              onChange={(e) => setFormData({ ...formData, jam_keluar: e.target.value })}
              placeholder="17:00"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select
              value={formData.status}
              onChange={(e) => setFormData({ ...formData, status: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              {statusOptions.map(opt => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </select>
          </div>
          <Input
            label="Keterangan"
            value={formData.keterangan}
            onChange={(e) => setFormData({ ...formData, keterangan: e.target.value })}
            placeholder="Keterangan tambahan (opsional)"
          />
          <div className="flex justify-end gap-3 mt-6">
            <Button type="button" onClick={handleCloseModal} variant="secondary">Batal</Button>
            <Button type="submit">{editingAbsensi ? 'Update' : 'Simpan'}</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default AbsensiList;
