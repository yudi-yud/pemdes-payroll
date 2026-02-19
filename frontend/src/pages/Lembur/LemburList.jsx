import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { lemburAPI, karyawanAPI } from '../../services/api';
import Button from '../../components/ui/Button';
import Modal from '../../components/ui/Modal';
import Input from '../../components/ui/Input';

const LemburList = () => {
  const [lembur, setLembur] = useState([]);
  const [karyawan, setKaryawan] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isApproveModalOpen, setIsApproveModalOpen] = useState(false);
  const [selectedLembur, setSelectedLembur] = useState(null);
  const [filterPeriod, setFilterPeriod] = useState({
    bulan: new Date().getMonth() + 1,
    tahun: new Date().getFullYear(),
  });
  const [editingLembur, setEditingLembur] = useState(null);
  const [formData, setFormData] = useState({
    karyawan_id: '',
    tanggal: '',
    jam_mulai: '',
    jam_selesai: '',
    total_jam: '',
    keterangan: '',
  });

  useEffect(() => {
    fetchKaryawan();
  }, []);

  useEffect(() => {
    fetchLembur();
  }, [filterPeriod]);

  const fetchKaryawan = async () => {
    try {
      const response = await karyawanAPI.getByStatus('aktif');
      setKaryawan(response.data);
    } catch (error) {
      console.error('Error fetching karyawan:', error);
    }
  };

  const fetchLembur = async () => {
    setLoading(true);
    try {
      const response = await lemburAPI.getByPeriod(filterPeriod.bulan, filterPeriod.tahun);
      setLembur(response.data);
    } catch (error) {
      console.error('Error fetching lembur:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleOpenModal = (lembur = null) => {
    if (lembur) {
      setEditingLembur(lembur);
      setFormData({
        karyawan_id: lembur.karyawan_id.toString(),
        tanggal: lembur.tanggal ? lembur.tanggal.split('T')[0] : '',
        jam_mulai: lembur.jam_mulai || '',
        jam_selesai: lembur.jam_selesai || '',
        total_jam: lembur.total_jam.toString(),
        keterangan: lembur.keterangan || '',
      });
    } else {
      setEditingLembur(null);
      setFormData({
        karyawan_id: '',
        tanggal: new Date().toISOString().split('T')[0],
        jam_mulai: '',
        jam_selesai: '',
        total_jam: '',
        keterangan: '',
      });
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingLembur(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const data = {
        ...formData,
        karyawan_id: parseInt(formData.karyawan_id),
        total_jam: parseFloat(formData.total_jam),
      };

      if (editingLembur) {
        await axios.put(`http://localhost:3000/api/lembur/${editingLembur.id}`, data);
      } else {
        // Get employee position to get overtime rate
        const karyawanData = karyawan.find(k => k.id === parseInt(data.karyawan_id));
        if (karyawanData?.jabatan) {
          data.tarif_per_jam = karyawanData.jabatan.tarif_lembur_per_jam || 0;
        }
        await lemburAPI.create(data);
      }

      await fetchLembur();
      handleCloseModal();
    } catch (error) {
      console.error('Error saving lembur:', error);
      alert('Gagal menyimpan lembur: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Apakah Anda yakin ingin menghapus data lembur ini?')) return;

    try {
      await lemburAPI.delete(id);
      await fetchLembur();
    } catch (error) {
      console.error('Error deleting lembur:', error);
      alert('Gagal menghapus lembur: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleOpenApproveModal = (lembur) => {
    setSelectedLembur(lembur);
    setIsApproveModalOpen(true);
  };

  const handleApprove = async (status) => {
    if (!selectedLembur) return;

    try {
      await lemburAPI.approve(selectedLembur.id, status, null);
      await fetchLembur();
      setIsApproveModalOpen(false);
      setSelectedLembur(null);
    } catch (error) {
      console.error('Error approving lembur:', error);
      alert('Gagal update status lembur: ' + (error.response?.data?.error || error.message));
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
    { value: 'pending', label: 'Pending', color: 'bg-yellow-100 text-yellow-800' },
    { value: 'disetujui', label: 'Disetujui', color: 'bg-green-100 text-green-800' },
    { value: 'ditolak', label: 'Ditolak', color: 'bg-red-100 text-red-800' },
  ];

  const getStatusBadge = (status) => {
    const s = statusOptions.find(opt => opt.value === status);
    return s ? (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${s.color}`}>
        {s.label}
      </span>
    ) : status;
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(amount);
  };

  // Calculate total lembur for current period
  const totalLembur = lembur
    .filter(l => l.status === 'disetujui')
    .reduce((sum, l) => sum + l.total_nominal, 0);

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Data Lembur</h1>
        <Button onClick={() => handleOpenModal()}>+ Tambah Lembur</Button>
      </div>

      <div className="bg-white rounded-lg shadow-md p-4 mb-6">
        <div className="flex gap-4 items-end">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Periode Bulan</label>
            <select
              value={filterPeriod.bulan}
              onChange={(e) => setFilterPeriod({ ...filterPeriod, bulan: parseInt(e.target.value) })}
              className="w-40 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
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
              value={filterPeriod.tahun}
              onChange={(e) => setFilterPeriod({ ...filterPeriod, tahun: parseInt(e.target.value) })}
              className="w-24 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              min="2020"
              max="2100"
            />
          </div>
          <div className="flex-1"></div>
          <div className="text-right">
            <p className="text-sm text-gray-600">Total Lembur (Disetujui)</p>
            <p className="text-xl font-bold text-green-600">{formatCurrency(totalLembur)}</p>
          </div>
        </div>
      </div>

      {loading ? (
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
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Jam</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total Jam</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tarif/Jam</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total Nominal</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Keterangan</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {lembur.length === 0 ? (
                <tr>
                  <td colSpan={10} className="px-6 py-4 text-center text-gray-500">
                    Tidak ada data lembur untuk periode ini
                  </td>
                </tr>
              ) : (
                lembur.map((row) => (
                  <tr key={row.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{row.tanggal}</td>
                    <td className="px-6 py-4 text-sm text-gray-900">{row.karyawan?.nama}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {row.jam_mulai} - {row.jam_selesai}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{row.total_jam}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{formatCurrency(row.tarif_per_jam)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{formatCurrency(row.total_nominal)}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{row.keterangan}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">{getStatusBadge(row.status)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <div className="flex gap-2">
                        {row.status === 'pending' && (
                          <>
                            <button onClick={() => handleOpenModal(row)} className="px-2 py-1 bg-blue-600 text-white rounded text-xs hover:bg-blue-700">Edit</button>
                            <button onClick={() => handleOpenApproveModal(row)} className="px-2 py-1 bg-green-600 text-white rounded text-xs hover:bg-green-700">Setuju</button>
                            <button onClick={() => handleApprove('ditolak')} className="px-2 py-1 bg-red-600 text-white rounded text-xs hover:bg-red-700">Tolak</button>
                          </>
                        )}
                        {(row.status === 'disetujui' || row.status === 'ditolak') && (
                          <button onClick={() => handleDelete(row.id)} className="px-2 py-1 bg-gray-600 text-white rounded text-xs hover:bg-gray-700">Hapus</button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}

      <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingLembur ? 'Edit Lembur' : 'Tambah Lembur'}>
        <form onSubmit={handleSubmit}>
          <Input
            label="Karyawan"
            type="select"
            value={formData.karyawan_id}
            onChange={(e) => setFormData({ ...formData, karyawan_id: e.target.value })}
            required
            disabled={!!editingLembur}
            options={karyawan.map(k => ({ value: k.id, label: `${k.nama} - ${k.jabatan?.nama_jabatan || 'Tanpa Jabatan'}` }))}
          />
          <Input
            label="Tanggal"
            type="date"
            value={formData.tanggal}
            onChange={(e) => setFormData({ ...formData, tanggal: e.target.value })}
            required
            disabled={!!editingLembur}
          />
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Jam Mulai"
              value={formData.jam_mulai}
              onChange={(e) => setFormData({ ...formData, jam_mulai: e.target.value })}
              placeholder="17:00"
              disabled={!!editingLembur}
            />
            <Input
              label="Jam Selesai"
              value={formData.jam_selesai}
              onChange={(e) => setFormData({ ...formData, jam_selesai: e.target.value })}
              placeholder="20:00"
              disabled={!!editingLembur}
            />
          </div>
          <Input
            label="Total Jam"
            type="number"
            value={formData.total_jam}
            onChange={(e) => setFormData({ ...formData, total_jam: e.target.value })}
            placeholder="Contoh: 2.5"
            required
            min="0.5"
            step="0.5"
          />
          <Input
            label="Keterangan"
            value={formData.keterangan}
            onChange={(e) => setFormData({ ...formData, keterangan: e.target.value })}
            placeholder="Pekerjaan yang dikerjakan"
          />
          <div className="flex justify-end gap-3 mt-6">
            <Button type="button" onClick={handleCloseModal} variant="secondary">Batal</Button>
            <Button type="submit">{editingLembur ? 'Update' : 'Simpan'}</Button>
          </div>
        </form>
      </Modal>

      <Modal isOpen={isApproveModalOpen} onClose={() => setIsApproveModalOpen(false)} title="Setujui Lembur">
        <div className="text-center">
          <p className="mb-4">Apakah Anda ingin menyetujui lembur ini?</p>
          <div className="flex gap-3 justify-center">
            <Button onClick={() => handleApprove('ditolak')} variant="danger">Tolak</Button>
            <Button onClick={() => handleApprove('disetujui')} variant="success">Setujui</Button>
            <Button onClick={() => setIsApproveModalOpen(false)} variant="secondary">Batal</Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default LemburList;
