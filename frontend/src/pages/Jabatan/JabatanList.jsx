import React, { useEffect, useState } from 'react';
import { jabatanAPI } from '../../services/api';
import Button from '../../components/ui/Button';
import Modal from '../../components/ui/Modal';
import Input from '../../components/ui/Input';

const JabatanList = () => {
  const [jabatan, setJabatan] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingJabatan, setEditingJabatan] = useState(null);
  const [formData, setFormData] = useState({
    nama_jabatan: '',
    gaji_pokok: '',
    tunjangan_jabatan: '',
    tarif_lembur_per_jam: '',
  });

  useEffect(() => {
    fetchJabatan();
  }, []);

  const fetchJabatan = async () => {
    try {
      const response = await jabatanAPI.getAll();
      setJabatan(response.data);
    } catch (error) {
      console.error('Error fetching jabatan:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleOpenModal = (jabatan = null) => {
    if (jabatan) {
      setEditingJabatan(jabatan);
      setFormData({
        nama_jabatan: jabatan.nama_jabatan,
        gaji_pokok: jabatan.gaji_pokok.toString(),
        tunjangan_jabatan: jabatan.tunjangan_jabatan.toString(),
        tarif_lembur_per_jam: (jabatan.tarif_lembur_per_jam || 0).toString(),
      });
    } else {
      setEditingJabatan(null);
      setFormData({
        nama_jabatan: '',
        gaji_pokok: '',
        tunjangan_jabatan: '',
        tarif_lembur_per_jam: '',
      });
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingJabatan(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const data = {
        nama_jabatan: formData.nama_jabatan,
        gaji_pokok: parseFloat(formData.gaji_pokok),
        tunjangan_jabatan: parseFloat(formData.tunjangan_jabatan) || 0,
        tarif_lembur_per_jam: parseFloat(formData.tarif_lembur_per_jam) || 0,
      };

      if (editingJabatan) {
        await jabatanAPI.update(editingJabatan.id, data);
      } else {
        await jabatanAPI.create(data);
      }

      await fetchJabatan();
      handleCloseModal();
    } catch (error) {
      console.error('Error saving jabatan:', error);
      alert('Gagal menyimpan jabatan: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Apakah Anda yakin ingin menghapus jabatan ini?')) return;

    try {
      await jabatanAPI.delete(id);
      await fetchJabatan();
    } catch (error) {
      console.error('Error deleting jabatan:', error);
      alert('Gagal menghapus jabatan: ' + (error.response?.data?.error || error.message));
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(amount);
  };

  const columns = [
    { header: 'No', accessor: 'id', render: (row, index) => index + 1 },
    { header: 'Nama Jabatan', accessor: 'nama_jabatan' },
    { header: 'Gaji Pokok', accessor: 'gaji_pokok', render: (row) => formatCurrency(row.gaji_pokok) },
    { header: 'Tunjangan Jabatan', accessor: 'tunjangan_jabatan', render: (row) => formatCurrency(row.tunjangan_jabatan) },
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
        <h1 className="text-2xl font-bold text-gray-800">Daftar Jabatan</h1>
        <Button onClick={() => handleOpenModal()}>+ Tambah Jabatan</Button>
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
            {jabatan.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-6 py-4 text-center text-gray-500">
                  Tidak ada data jabatan
                </td>
              </tr>
            ) : (
              jabatan.map((row, rowIndex) => (
                <tr key={row.id} className="hover:bg-gray-50">
                  {columns.map((col, colIndex) => (
                    <td key={colIndex} className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {col.render ? col.render(row, rowIndex) : row[col.accessor]}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingJabatan ? 'Edit Jabatan' : 'Tambah Jabatan'}>
        <form onSubmit={handleSubmit}>
          <Input
            label="Nama Jabatan"
            value={formData.nama_jabatan}
            onChange={(e) => setFormData({ ...formData, nama_jabatan: e.target.value })}
            placeholder="Masukkan nama jabatan"
            required
          />
          <Input
            label="Gaji Pokok"
            type="number"
            value={formData.gaji_pokok}
            onChange={(e) => setFormData({ ...formData, gaji_pokok: e.target.value })}
            placeholder="Masukkan gaji pokok"
            required
            min="0"
            step="0.01"
          />
          <Input
            label="Tunjangan Jabatan"
            type="number"
            value={formData.tunjangan_jabatan}
            onChange={(e) => setFormData({ ...formData, tunjangan_jabatan: e.target.value })}
            placeholder="Masukkan tunjangan jabatan"
            min="0"
            step="0.01"
          />
          <Input
            label="Tarif Lembur Per Jam"
            type="number"
            value={formData.tarif_lembur_per_jam}
            onChange={(e) => setFormData({ ...formData, tarif_lembur_per_jam: e.target.value })}
            placeholder="Masukkan tarif lembur per jam"
            min="0"
            step="0.01"
          />
          <div className="flex justify-end gap-3 mt-6">
            <Button type="button" onClick={handleCloseModal} variant="secondary">
              Batal
            </Button>
            <Button type="submit">
              {editingJabatan ? 'Update' : 'Simpan'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default JabatanList;
