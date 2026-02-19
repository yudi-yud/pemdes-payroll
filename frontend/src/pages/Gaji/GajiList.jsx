import React, { useEffect, useState } from 'react';
import { gajiAPI, karyawanAPI } from '../../services/api';
import Button from '../../components/ui/Button';
import Modal from '../../components/ui/Modal';
import Input from '../../components/ui/Input';

const GajiList = () => {
  const [gaji, setGaji] = useState([]);
  const [karyawan, setKaryawan] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isGenerateModalOpen, setIsGenerateModalOpen] = useState(false);
  const [filterBulan, setFilterBulan] = useState(new Date().getMonth() + 1);
  const [filterTahun, setFilterTahun] = useState(new Date().getFullYear());
  const [editingGaji, setEditingGaji] = useState(null);
  const [formData, setFormData] = useState({
    karyawan_id: '',
    periode_bulan: new Date().getMonth() + 1,
    periode_tahun: new Date().getFullYear(),
    gaji_pokok: '',
    tunjangan_jabatan: '',
    tunjangan_transport: '',
    tunjangan_makan: '',
    lembur: '',
    potongan: '',
  });
  const [generateData, setGenerateData] = useState({
    periode_bulan: new Date().getMonth() + 1,
    periode_tahun: new Date().getFullYear(),
    tunjangan_transport: '',
    tunjangan_makan: '',
  });

  useEffect(() => {
    fetchData();
  }, [filterBulan, filterTahun]);

  const fetchData = async () => {
    try {
      const [gajiRes, karyawanRes] = await Promise.all([
        gajiAPI.getByPeriod(filterBulan, filterTahun),
        karyawanAPI.getByStatus('aktif'),
      ]);
      setGaji(gajiRes.data);
      setKaryawan(karyawanRes.data);
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleOpenModal = async (gaji = null) => {
    if (gaji) {
      setEditingGaji(gaji);
      setFormData({
        karyawan_id: gaji.karyawan_id.toString(),
        periode_bulan: gaji.periode_bulan,
        periode_tahun: gaji.periode_tahun,
        gaji_pokok: gaji.gaji_pokok.toString(),
        tunjangan_jabatan: gaji.tunjangan_jabatan.toString(),
        tunjangan_transport: gaji.tunjangan_transport.toString(),
        tunjangan_makan: gaji.tunjangan_makan.toString(),
        lembur: gaji.lembur.toString(),
        potongan: gaji.potongan.toString(),
      });
    } else {
      setEditingGaji(null);
      setFormData({
        karyawan_id: '',
        periode_bulan: filterBulan,
        periode_tahun: filterTahun,
        gaji_pokok: '',
        tunjangan_jabatan: '',
        tunjangan_transport: '',
        tunjangan_makan: '',
        lembur: '',
        potongan: '',
      });
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingGaji(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const data = {
        ...formData,
        karyawan_id: parseInt(formData.karyawan_id),
        gaji_pokok: parseFloat(formData.gaji_pokok),
        tunjangan_jabatan: parseFloat(formData.tunjangan_jabatan) || 0,
        tunjangan_transport: parseFloat(formData.tunjangan_transport) || 0,
        tunjangan_makan: parseFloat(formData.tunjangan_makan) || 0,
        lembur: parseFloat(formData.lembur) || 0,
        potongan: parseFloat(formData.potongan) || 0,
      };

      if (editingGaji) {
        await gajiAPI.update(editingGaji.id, data);
      } else {
        await gajiAPI.create(data);
      }

      await fetchData();
      handleCloseModal();
    } catch (error) {
      console.error('Error saving gaji:', error);
      alert('Gagal menyimpan gaji: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Apakah Anda yakin ingin menghapus data gaji ini?')) return;

    try {
      await gajiAPI.delete(id);
      await fetchData();
    } catch (error) {
      console.error('Error deleting gaji:', error);
      alert('Gagal menghapus gaji: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleUpdateStatus = async (id, status) => {
    try {
      await gajiAPI.updateStatus(id, status);
      await fetchData();
    } catch (error) {
      console.error('Error updating status:', error);
      alert('Gagal update status: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleGenerateBatch = async (e) => {
    e.preventDefault();
    try {
      const data = {
        periode_bulan: parseInt(generateData.periode_bulan),
        periode_tahun: parseInt(generateData.periode_tahun),
        tunjangan_transport: parseFloat(generateData.tunjangan_transport) || 0,
        tunjangan_makan: parseFloat(generateData.tunjangan_makan) || 0,
      };

      const response = await gajiAPI.generateBatch(data);
      alert(`Berhasil generate ${response.data.created} gaji. ${response.data.skipped?.length || 0} sudah ada.`);
      setIsGenerateModalOpen(false);
      await fetchData();
    } catch (error) {
      console.error('Error generating batch:', error);
      alert('Gagal generate batch: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleKaryawanChange = async (karyawanId) => {
    setFormData({ ...formData, karyawan_id: karyawanId });
    const karyawanData = karyawan.find(k => k.id === parseInt(karyawanId));
    if (karyawanData?.jabatan) {
      setFormData(prev => ({
        ...prev,
        gaji_pokok: karyawanData.jabatan.gaji_pokok.toString(),
        tunjangan_jabatan: karyawanData.jabatan.tunjangan_jabatan.toString(),
      }));
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
    }).format(amount);
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

  const columns = [
    { header: 'Karyawan', render: (row) => row.karyawan?.nama || '-' },
    { header: 'NIK', render: (row) => row.karyawan?.nik || '-' },
    { header: 'Jabatan', render: (row) => row.karyawan?.jabatan?.nama_jabatan || '-' },
    { header: 'Periode', render: (row) => `${bulanOptions[row.periode_bulan - 1]?.label} ${row.periode_tahun}` },
    { header: 'Gaji Pokok', render: (row) => formatCurrency(row.gaji_pokok) },
    { header: 'Tunjangan', render: (row) => formatCurrency(row.tunjangan_jabatan + row.tunjangan_transport + row.tunjangan_makan) },
    { header: 'Lembur', render: (row) => formatCurrency(row.lembur) },
    { header: 'Potongan', render: (row) => formatCurrency(row.potongan) },
    { header: 'Total Gaji', render: (row) => formatCurrency(row.total_gaji) },
    {
      header: 'Status',
      render: (row) => (
        <span className={`px-2 py-1 rounded-full text-xs font-medium ${
          row.status === 'dibayar' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
        }`}>
          {row.status === 'dibayar' ? 'Dibayar' : 'Pending'}
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
          {row.status === 'pending' && (
            <Button onClick={() => handleUpdateStatus(row.id, 'dibayar')} variant="success" className="px-3 py-1 text-sm">
              Bayar
            </Button>
          )}
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
        <h1 className="text-2xl font-bold text-gray-800">Data Gaji</h1>
        <div className="flex gap-3">
          <Button onClick={() => setIsGenerateModalOpen(true)} variant="success">
            Generate Batch
          </Button>
          <Button onClick={() => handleOpenModal()}>+ Tambah Gaji</Button>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-4 mb-6">
        <div className="flex gap-4 items-end">
          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-1">Periode Bulan</label>
            <select
              value={filterBulan}
              onChange={(e) => setFilterBulan(parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {bulanOptions.map(b => (
                <option key={b.value} value={b.value}>{b.label}</option>
              ))}
            </select>
          </div>
          <div className="w-32">
            <label className="block text-sm font-medium text-gray-700 mb-1">Tahun</label>
            <input
              type="number"
              value={filterTahun}
              onChange={(e) => setFilterTahun(parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              min="2000"
              max="2100"
            />
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {columns.map((col, i) => (
                <th key={i} className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap">
                  {col.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {gaji.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-6 py-4 text-center text-gray-500">
                  Tidak ada data gaji untuk periode ini
                </td>
              </tr>
            ) : (
              gaji.map((row) => (
                <tr key={row.id} className="hover:bg-gray-50">
                  {columns.map((col, colIndex) => (
                    <td key={colIndex} className="px-4 py-4 whitespace-nowrap text-sm text-gray-900">
                      {col.render ? col.render(row) : row[col.accessor]}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingGaji ? 'Edit Gaji' : 'Tambah Gaji'} size="lg">
        <form onSubmit={handleSubmit}>
          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">Karyawan</label>
              <select
                value={formData.karyawan_id}
                onChange={(e) => handleKaryawanChange(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
                disabled={!!editingGaji}
              >
                <option value="">-- Pilih Karyawan --</option>
                {karyawan.map((k) => (
                  <option key={k.id} value={k.id}>
                    {k.nama} - {k.jabatan?.nama_jabatan || 'Tanpa Jabatan'}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Periode Bulan</label>
              <select
                value={formData.periode_bulan}
                onChange={(e) => setFormData({ ...formData, periode_bulan: parseInt(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              >
                {bulanOptions.map(b => (
                  <option key={b.value} value={b.value}>{b.label}</option>
                ))}
              </select>
            </div>
            <Input
              label="Periode Tahun"
              type="number"
              value={formData.periode_tahun}
              onChange={(e) => setFormData({ ...formData, periode_tahun: parseInt(e.target.value) })}
              required
              min="2000"
              max="2100"
            />
            <Input
              label="Gaji Pokok"
              type="number"
              value={formData.gaji_pokok}
              onChange={(e) => setFormData({ ...formData, gaji_pokok: e.target.value })}
              required
              min="0"
              step="0.01"
            />
            <Input
              label="Tunjangan Jabatan"
              type="number"
              value={formData.tunjangan_jabatan}
              onChange={(e) => setFormData({ ...formData, tunjangan_jabatan: e.target.value })}
              min="0"
              step="0.01"
            />
            <Input
              label="Tunjangan Transport"
              type="number"
              value={formData.tunjangan_transport}
              onChange={(e) => setFormData({ ...formData, tunjangan_transport: e.target.value })}
              min="0"
              step="0.01"
            />
            <Input
              label="Tunjangan Makan"
              type="number"
              value={formData.tunjangan_makan}
              onChange={(e) => setFormData({ ...formData, tunjangan_makan: e.target.value })}
              min="0"
              step="0.01"
            />
            <Input
              label="Lembur"
              type="number"
              value={formData.lembur}
              onChange={(e) => setFormData({ ...formData, lembur: e.target.value })}
              min="0"
              step="0.01"
            />
            <Input
              label="Potongan"
              type="number"
              value={formData.potongan}
              onChange={(e) => setFormData({ ...formData, potongan: e.target.value })}
              min="0"
              step="0.01"
            />
          </div>
          <div className="flex justify-end gap-3 mt-6">
            <Button type="button" onClick={handleCloseModal} variant="secondary">
              Batal
            </Button>
            <Button type="submit">
              {editingGaji ? 'Update' : 'Simpan'}
            </Button>
          </div>
        </form>
      </Modal>

      <Modal isOpen={isGenerateModalOpen} onClose={() => setIsGenerateModalOpen(false)} title="Generate Gaji Batch">
        <form onSubmit={handleGenerateBatch}>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Periode Bulan</label>
              <select
                value={generateData.periode_bulan}
                onChange={(e) => setGenerateData({ ...generateData, periode_bulan: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              >
                {bulanOptions.map(b => (
                  <option key={b.value} value={b.value}>{b.label}</option>
                ))}
              </select>
            </div>
            <Input
              label="Periode Tahun"
              type="number"
              value={generateData.periode_tahun}
              onChange={(e) => setGenerateData({ ...generateData, periode_tahun: e.target.value })}
              required
              min="2000"
              max="2100"
            />
            <Input
              label="Tunjangan Transport (Semua)"
              type="number"
              value={generateData.tunjangan_transport}
              onChange={(e) => setGenerateData({ ...generateData, tunjangan_transport: e.target.value })}
              placeholder="Opsional"
              min="0"
              step="0.01"
            />
            <Input
              label="Tunjangan Makan (Semua)"
              type="number"
              value={generateData.tunjangan_makan}
              onChange={(e) => setGenerateData({ ...generateData, tunjangan_makan: e.target.value })}
              placeholder="Opsional"
              min="0"
              step="0.01"
            />
          </div>
          <p className="text-sm text-gray-500 mt-4">
            Generate gaji untuk semua karyawan aktif. Gaji pokok dan tunjangan jabatan akan diambil dari data jabatan masing-masing karyawan.
          </p>
          <div className="flex justify-end gap-3 mt-6">
            <Button type="button" onClick={() => setIsGenerateModalOpen(false)} variant="secondary">
              Batal
            </Button>
            <Button type="submit">Generate</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default GajiList;
