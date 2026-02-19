import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { laporanAPI } from '../../services/api';
import Button from '../../components/ui/Button';

const LaporanGaji = () => {
  const [laporan, setLaporan] = useState([]);
  const [rekap, setRekap] = useState(null);
  const [loading, setLoading] = useState(false);
  const [bulan, setBulan] = useState(new Date().getMonth() + 1);
  const [tahun, setTahun] = useState(new Date().getFullYear());

  useEffect(() => {
    fetchLaporan();
  }, [bulan, tahun]);

  const fetchLaporan = async () => {
    setLoading(true);
    try {
      const [laporanRes, rekapRes] = await Promise.all([
        laporanAPI.getLaporanGaji(bulan, tahun),
        laporanAPI.getRekapGaji(bulan, tahun),
      ]);
      setLaporan(laporanRes.data);
      setRekap(rekapRes.data);
    } catch (error) {
      console.error('Error fetching laporan:', error);
    } finally {
      setLoading(false);
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

  const handlePrint = () => {
    window.print();
  };

  const handleExportCSV = () => {
    if (laporan.length === 0) return;

    const headers = ['NIK', 'Nama Karyawan', 'Jabatan', 'Gaji Pokok', 'Tunjangan Jabatan', 'Transport', 'Makan', 'Lembur', 'Potongan', 'Total Gaji', 'Status'];
    const csvContent = [
      headers.join(','),
      ...laporan.map(row => [
        row.nik,
        `"${row.nama_karyawan}"`,
        `"${row.jabatan}"`,
        row.gaji_pokok,
        row.tunjangan_jabatan,
        row.tunjangan_transport,
        row.tunjangan_makan,
        row.lembur,
        row.potongan,
        row.total_gaji,
        row.status,
      ].join(','))
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `laporan_gaji_${bulan}_${tahun}.csv`;
    link.click();
  };

  const handleExportExcel = async () => {
    try {
      const response = await axios.get(`http://localhost:3000/api/laporan/export/excel?bulan=${bulan}&tahun=${tahun}`, {
        responseType: 'blob',
      });

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.download = `Laporan_Gaji_${bulanOptions[bulan - 1].label}_${tahun}.xlsx`;
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (error) {
      console.error('Error exporting Excel:', error);
      alert('Gagal export Excel: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleExportPDF = async (karyawanId, nama) => {
    try {
      const response = await axios.get(`http://localhost:3000/api/laporan/export/karyawan/${karyawanId}/pdf`, {
        responseType: 'blob',
      });

      const url = window.URL.createObjectURL(new Blob([response.data], { type: 'application/pdf' }));
      const link = document.createElement('a');
      link.href = url;
      link.download = `Laporan_Gaji_${nama.replace(/\s+/g, '_')}.pdf`;
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (error) {
      console.error('Error exporting PDF:', error);
      alert('Gagal export PDF: ' + (error.response?.data?.error || error.message));
    }
  };

  const columns = [
    { header: 'No', render: (row, index) => index + 1 },
    { header: 'NIK', accessor: 'nik' },
    { header: 'Nama Karyawan', accessor: 'nama_karyawan' },
    { header: 'Jabatan', accessor: 'jabatan' },
    { header: 'Gaji Pokok', render: (row) => formatCurrency(row.gaji_pokok) },
    { header: 'Tunj. Jabatan', render: (row) => formatCurrency(row.tunjangan_jabatan) },
    { header: 'Transport', render: (row) => formatCurrency(row.tunjangan_transport) },
    { header: 'Makan', render: (row) => formatCurrency(row.tunjangan_makan) },
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
        <div className="flex gap-2 print:hidden">
          <button
            onClick={() => handleExportPDF(row.karyawan_id, row.nama_karyawan)}
            className="px-3 py-1 bg-red-600 text-white rounded-lg hover:bg-red-700 text-sm font-medium transition-colors"
            title="Export PDF Per Karyawan"
          >
            📄 PDF
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-6 print:hidden">
        <h1 className="text-2xl font-bold text-gray-800">Laporan Gaji</h1>
        <div className="flex gap-3">
          <Button onClick={handleExportExcel} variant="success">📊 Export Excel</Button>
          <Button onClick={handleExportCSV} variant="secondary">Export CSV</Button>
          <Button onClick={handlePrint}>Cetak</Button>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-4 mb-6 print:hidden">
        <div className="flex gap-4 items-end">
          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-1">Periode Bulan</label>
            <select
              value={bulan}
              onChange={(e) => setBulan(parseInt(e.target.value))}
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
              value={tahun}
              onChange={(e) => setTahun(parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              min="2000"
              max="2100"
            />
          </div>
        </div>
      </div>

      {rekap && (
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Rekapitulasi Gaji - {bulanOptions[bulan - 1]?.label} {tahun}</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="bg-blue-50 rounded-lg p-4">
              <p className="text-sm text-gray-600">Total Karyawan</p>
              <p className="text-2xl font-bold text-blue-600">{rekap.total_karyawan}</p>
            </div>
            <div className="bg-green-50 rounded-lg p-4">
              <p className="text-sm text-gray-600">Total Gaji</p>
              <p className="text-2xl font-bold text-green-600">{formatCurrency(rekap.total_gaji)}</p>
            </div>
            <div className="bg-yellow-50 rounded-lg p-4">
              <p className="text-sm text-gray-600">Status Pending</p>
              <p className="text-2xl font-bold text-yellow-600">{rekap.status_pending}</p>
            </div>
            <div className="bg-purple-50 rounded-lg p-4">
              <p className="text-sm text-gray-600">Status Dibayar</p>
              <p className="text-2xl font-bold text-purple-600">{rekap.status_dibayar}</p>
            </div>
          </div>
          <div className="mt-4 grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
            <div>
              <span className="text-gray-600">Total Gaji Pokok: </span>
              <span className="font-semibold">{formatCurrency(rekap.total_gaji_pokok)}</span>
            </div>
            <div>
              <span className="text-gray-600">Total Tunjangan: </span>
              <span className="font-semibold">{formatCurrency(rekap.total_tunjangan)}</span>
            </div>
            <div>
              <span className="text-gray-600">Total Lembur: </span>
              <span className="font-semibold">{formatCurrency(rekap.total_lembur)}</span>
            </div>
            <div>
              <span className="text-gray-600">Total Potongan: </span>
              <span className="font-semibold text-red-600">{formatCurrency(rekap.total_potongan)}</span>
            </div>
          </div>
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-md overflow-x-auto">
          <div className="p-4 border-b print:px-0 print:pt-0">
            <h2 className="text-lg font-semibold">Daftar Gaji - {bulanOptions[bulan - 1]?.label} {tahun}</h2>
          </div>
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
              {laporan.length === 0 ? (
                <tr>
                  <td colSpan={columns.length} className="px-6 py-4 text-center text-gray-500">
                    Tidak ada data gaji untuk periode ini
                  </td>
                </tr>
              ) : (
                laporan.map((row, index) => (
                  <tr key={row.id} className="hover:bg-gray-50">
                    {columns.map((col, colIndex) => (
                      <td key={colIndex} className="px-4 py-4 whitespace-nowrap text-sm text-gray-900">
                        {col.render ? col.render(row, index) : row[col.accessor]}
                      </td>
                    ))}
                  </tr>
                ))
              )}
            </tbody>
            {laporan.length > 0 && (
              <tfoot className="bg-gray-50 font-semibold">
                <tr>
                  <td colSpan={10} className="px-4 py-3 text-right">GRAND TOTAL:</td>
                  <td className="px-4 py-3">
                    {formatCurrency(laporan.reduce((sum, row) => sum + row.total_gaji, 0))}
                  </td>
                  <td></td>
                  <td></td>
                </tr>
              </tfoot>
            )}
          </table>
        </div>
      )}
    </div>
  );
};

export default LaporanGaji;
