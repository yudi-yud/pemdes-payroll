import React, { useEffect, useState } from 'react';
import { portalAPI } from '../../services/api';

const SlipGaji = () => {
  const [slips, setSlips] = useState([]);
  const [selectedSlip, setSelectedSlip] = useState(null);
  const [loading, setLoading] = useState(false);
  const [bulan, setBulan] = useState(new Date().getMonth() + 1);
  const [tahun, setTahun] = useState(new Date().getFullYear());

  useEffect(() => {
    fetchSlips();
  }, [bulan, tahun]);

  const fetchSlips = async () => {
    setLoading(true);
    try {
      const response = await portalAPI.getMySlips(bulan, tahun);
      // Handle null response from backend
      setSlips(response.data || []);
    } catch (error) {
      console.error('Error fetching slips:', error);
      alert('Gagal mengambil data slip gaji');
    } finally {
      setLoading(false);
    }
  };

  const viewSlip = async (id) => {
    try {
      const response = await portalAPI.getSlip(id);
      setSelectedSlip(response.data);
    } catch (error) {
      console.error('Error fetching slip:', error);
      alert('Gagal mengambil detail slip gaji');
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

  const getStatusBadge = (status) => {
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${
        status === 'dibayar' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
      }`}>
        {status === 'dibayar' ? 'Dibayar' : 'Pending'}
      </span>
    );
  };

  const printSlip = () => {
    if (!selectedSlip) return;
    const printContent = document.getElementById('slip-content');
    const originalContent = document.body.innerHTML;

    document.body.innerHTML = printContent.innerHTML;
    window.print();
    document.body.innerHTML = originalContent;
    window.location.reload();
  };

  return (
    <div className="max-w-6xl mx-auto">
      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <h1 className="text-2xl font-bold text-gray-800 mb-6">Slip Gaji Saya</h1>
        <div className="flex gap-4 items-end">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Bulan</label>
            <select
              value={bulan}
              onChange={(e) => setBulan(parseInt(e.target.value))}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {bulanOptions.map((b) => (
                <option key={b.value} value={b.value}>{b.label}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Tahun</label>
            <input
              type="number"
              value={tahun}
              onChange={(e) => setTahun(parseInt(e.target.value))}
              className="w-24 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              min="2020"
              max="2100"
            />
          </div>
        </div>
      </div>

      {selectedSlip ? (
        <div className="bg-white rounded-lg shadow-md p-8">
          <div className="flex justify-between items-start mb-6">
            <button
              onClick={() => setSelectedSlip(null)}
              className="text-blue-600 hover:text-blue-800 mb-4"
            >
              ← Kembali
            </button>
            <button
              onClick={printSlip}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              🖨 Cetak Slip
            </button>
          </div>

          <div id="slip-content" className="border-2 border-gray-300 p-8">
            <div className="text-center mb-8 pb-4 border-b-2 border-gray-300">
              <h1 className="text-2xl font-bold text-gray-800">SLIP GAJI</h1>
              <p className="text-gray-600">Pemerintah Desa</p>
            </div>

            <div className="grid grid-cols-2 gap-4 mb-6">
              <div>
                <p className="text-sm text-gray-500">Nama Karyawan</p>
                <p className="font-medium text-gray-800">{selectedSlip.karyawan?.nama || '-'}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">NIK</p>
                <p className="font-medium text-gray-800">{selectedSlip.karyawan?.nik || '-'}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Jabatan</p>
                <p className="font-medium text-gray-800">{selectedSlip.karyawan?.jabatan?.nama_jabatan || '-'}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Periode</p>
                <p className="font-medium text-gray-800">
                  {bulanOptions.find(b => b.value === selectedSlip.periode_bulan)?.label} {selectedSlip.periode_tahun}
                </p>
              </div>
            </div>

            <table className="w-full mb-6">
              <thead>
                <tr className="bg-gray-100">
                  <th className="border border-gray-300 px-4 py-2 text-left">Keterangan</th>
                  <th className="border border-gray-300 px-4 py-2 text-right">Jumlah</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td className="border border-gray-300 px-4 py-2">Gaji Pokok</td>
                  <td className="border border-gray-300 px-4 py-2 text-right">{formatCurrency(selectedSlip.gaji_pokok)}</td>
                </tr>
                <tr>
                  <td className="border border-gray-300 px-4 py-2">Tunjangan Jabatan</td>
                  <td className="border border-gray-300 px-4 py-2 text-right">{formatCurrency(selectedSlip.tunjangan_jabatan)}</td>
                </tr>
                <tr>
                  <td className="border border-gray-300 px-4 py-2">Tunjangan Transport</td>
                  <td className="border border-gray-300 px-4 py-2 text-right">{formatCurrency(selectedSlip.tunjangan_transport)}</td>
                </tr>
                <tr>
                  <td className="border border-gray-300 px-4 py-2">Tunjangan Makan</td>
                  <td className="border border-gray-300 px-4 py-2 text-right">{formatCurrency(selectedSlip.tunjangan_makan)}</td>
                </tr>
                <tr>
                  <td className="border border-gray-300 px-4 py-2">Lembur</td>
                  <td className="border border-gray-300 px-4 py-2 text-right">{formatCurrency(selectedSlip.lembur)}</td>
                </tr>
                <tr className="bg-red-50">
                  <td className="border border-gray-300 px-4 py-2 font-medium">Potongan</td>
                  <td className="border border-gray-300 px-4 py-2 text-right text-red-600">{formatCurrency(selectedSlip.potongan)}</td>
                </tr>
                <tr className="bg-blue-50 font-bold">
                  <td className="border border-gray-300 px-4 py-2">TOTAL GAJI BERSIH</td>
                  <td className="border border-gray-300 px-4 py-2 text-right text-blue-600">{formatCurrency(selectedSlip.total_gaji)}</td>
                </tr>
              </tbody>
            </table>

            <div className="grid grid-cols-2 gap-8 mt-8 pt-4 border-t-2 border-gray-300">
              <div className="text-sm text-gray-500">
                <p>Tanggal Cetak: {new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</p>
              </div>
              <div className="text-right text-sm text-gray-500">
                <p>Mengetahui,</p>
                <p className="font-medium">HRD Pemerintah Desa</p>
              </div>
            </div>
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-md overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Periode</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Gaji Pokok</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tunjangan</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Lembur</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Potongan</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total Gaji</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Aksi</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {slips.length === 0 ? (
                <tr>
                  <td colSpan={8} className="px-6 py-4 text-center text-gray-500">
                    {loading ? 'Loading...' : 'Tidak ada data slip gaji untuk periode ini'}
                  </td>
                </tr>
              ) : (
                slips.map((slip) => (
                  <tr key={slip.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm text-gray-900">
                      {bulanOptions.find(b => b.value === slip.periode_bulan)?.label} {slip.periode_tahun}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">{formatCurrency(slip.gaji_pokok)}</td>
                    <td className="px-6 py-4 text-sm text-gray-900">
                      {formatCurrency(slip.tunjangan_jabatan + slip.tunjangan_transport + slip.tunjangan_makan)}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">{formatCurrency(slip.lembur)}</td>
                    <td className="px-6 py-4 text-sm text-gray-900">{formatCurrency(slip.potongan)}</td>
                    <td className="px-6 py-4 text-sm font-medium text-gray-900">{formatCurrency(slip.total_gaji)}</td>
                    <td className="px-6 py-4 text-sm">{getStatusBadge(slip.status)}</td>
                    <td className="px-6 py-4 text-sm">
                      <button
                        onClick={() => viewSlip(slip.id)}
                        className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
                      >
                        Lihat Slip
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default SlipGaji;
