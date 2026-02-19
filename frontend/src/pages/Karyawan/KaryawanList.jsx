import React, { useEffect, useState } from 'react';
import { karyawanAPI, jabatanAPI } from '../../services/api';
import Button from '../../components/ui/Button';
import Modal from '../../components/ui/Modal';
import Input from '../../components/ui/Input';

const KaryawanList = () => {
  const [karyawan, setKaryawan] = useState([]);
  const [jabatan, setJabatan] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [editingKaryawan, setEditingKaryawan] = useState(null);
  const [formData, setFormData] = useState({
    nik: '',
    nama: '',
    email: '',
    telepon: '',
    alamat: '',
    jabatan_id: '',
    tanggal_bergabung: '',
    status: 'aktif',
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [karyawanRes, jabatanRes] = await Promise.all([
        karyawanAPI.getAll(),
        jabatanAPI.getAll(),
      ]);
      setKaryawan(karyawanRes.data);
      setJabatan(jabatanRes.data);
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async () => {
    if (!searchTerm) {
      fetchData();
      return;
    }
    try {
      const response = await karyawanAPI.search(searchTerm);
      setKaryawan(response.data);
    } catch (error) {
      console.error('Error searching:', error);
    }
  };

  const handleOpenModal = (karyawan = null) => {
    if (karyawan) {
      setEditingKaryawan(karyawan);
      setFormData({
        nik: karyawan.nik,
        nama: karyawan.nama,
        email: karyawan.email || '',
        telepon: karyawan.telepon || '',
        alamat: karyawan.alamat || '',
        jabatan_id: karyawan.jabatan_id?.toString() || '',
        tanggal_bergabung: karyawan.tanggal_bergabung?.split('T')[0] || '',
        status: karyawan.status || 'aktif',
      });
    } else {
      setEditingKaryawan(null);
      setFormData({
        nik: '',
        nama: '',
        email: '',
        telepon: '',
        alamat: '',
        jabatan_id: '',
        tanggal_bergabung: '',
        status: 'aktif',
      });
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingKaryawan(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const data = { ...formData };
      data.jabatan_id = data.jabatan_id ? parseInt(data.jabatan_id) : null;

      if (editingKaryawan) {
        await karyawanAPI.update(editingKaryawan.id, data);
      } else {
        await karyawanAPI.create(data);
      }

      await fetchData();
      handleCloseModal();
    } catch (error) {
      console.error('Error saving karyawan:', error);
      alert('Gagal menyimpan karyawan: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Apakah Anda yakin ingin menghapus karyawan ini?')) return;

    try {
      await karyawanAPI.delete(id);
      await fetchData();
    } catch (error) {
      console.error('Error deleting karyawan:', error);
      alert('Gagal menghapus karyawan: ' + (error.response?.data?.error || error.message));
    }
  };

  const columns = [
    { header: 'NIK', accessor: 'nik' },
    { header: 'Nama', accessor: 'nama' },
    { header: 'Jabatan', render: (row) => row.jabatan?.nama_jabatan || '-' },
    { header: 'Email', accessor: 'email', render: (row) => row.email || '-' },
    { header: 'Telepon', accessor: 'telepon', render: (row) => row.telepon || '-' },
    {
      header: 'Status',
      render: (row) => (
        <span className={`px-2 py-1 rounded-full text-xs font-medium ${
          row.status === 'aktif' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
        }`}>
          {row.status === 'aktif' ? 'Aktif' : 'Non-Aktif'}
        </span>
      ),
    },
    {
      header: 'Aksi',
      render: (row) => (
        <div className="flex gap-2">
          <Button onClick={() => handleOpenModal(row)} variant="secondary" className="px-3 py-1 text-sm">
            Edit
          </Button>
          <Button onClick={() => handleDelete(row.id)} variant="danger" className="px-3 py-1 text-sm">
            Hapus
          </Button>
        </div>
      ),
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
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Daftar Karyawan</h1>
        <Button onClick={() => handleOpenModal()}>+ Tambah Karyawan</Button>
      </div>

      <div className="bg-white rounded-lg shadow-md p-4 mb-6">
        <div className="flex gap-4">
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
            placeholder="Cari berdasarkan NIK, Nama, atau Email..."
            className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <Button onClick={handleSearch}>Cari</Button>
          <Button onClick={fetchData} variant="secondary">Reset</Button>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {columns.map((col, i) => (
                <th key={i} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  {col.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {karyawan.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-6 py-4 text-center text-gray-500">
                  Tidak ada data karyawan
                </td>
              </tr>
            ) : (
              karyawan.map((row) => (
                <tr key={row.id} className="hover:bg-gray-50">
                  {columns.map((col, colIndex) => (
                    <td key={colIndex} className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {col.render ? col.render(row) : row[col.accessor]}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingKaryawan ? 'Edit Karyawan' : 'Tambah Karyawan'} size="lg">
        <form onSubmit={handleSubmit}>
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="NIK"
              value={formData.nik}
              onChange={(e) => setFormData({ ...formData, nik: e.target.value })}
              placeholder="Masukkan NIK"
              required
            />
            <Input
              label="Nama"
              value={formData.nama}
              onChange={(e) => setFormData({ ...formData, nama: e.target.value })}
              placeholder="Masukkan nama lengkap"
              required
            />
            <Input
              label="Email"
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              placeholder="email@example.com"
            />
            <Input
              label="Telepon"
              value={formData.telepon}
              onChange={(e) => setFormData({ ...formData, telepon: e.target.value })}
              placeholder="08xxxxxxxxxx"
            />
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-1">Jabatan</label>
              <select
                value={formData.jabatan_id}
                onChange={(e) => setFormData({ ...formData, jabatan_id: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">-- Pilih Jabatan --</option>
                {jabatan.map((j) => (
                  <option key={j.id} value={j.id}>
                    {j.nama_jabatan}
                  </option>
                ))}
              </select>
            </div>
            <Input
              label="Tanggal Bergabung"
              type="date"
              value={formData.tanggal_bergabung}
              onChange={(e) => setFormData({ ...formData, tanggal_bergabung: e.target.value })}
            />
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select
                value={formData.status}
                onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="aktif">Aktif</option>
                <option value="non_aktif">Non-Aktif</option>
              </select>
            </div>
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">Alamat</label>
            <textarea
              value={formData.alamat}
              onChange={(e) => setFormData({ ...formData, alamat: e.target.value })}
              placeholder="Masukkan alamat lengkap"
              rows="3"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div className="flex justify-end gap-3 mt-6">
            <Button type="button" onClick={handleCloseModal} variant="secondary">
              Batal
            </Button>
            <Button type="submit">
              {editingKaryawan ? 'Update' : 'Simpan'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default KaryawanList;
