import React, { useEffect, useState } from 'react';
import { userAPI } from '../../services/api';
import Button from '../../components/ui/Button';
import Modal from '../../components/ui/Modal';
import Input from '../../components/ui/Input';

const UserManagement = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    name: '',
    email: '',
    role: 'karyawan',
  });

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const response = await userAPI.getAll();
      setUsers(response.data);
    } catch (error) {
      console.error('Error fetching users:', error);
      alert('Gagal mengambil data users');
    } finally {
      setLoading(false);
    }
  };

  const handleOpenModal = (user = null) => {
    if (user) {
      setEditingUser(user);
      setFormData({
        username: user.username,
        password: '',
        name: user.name,
        email: user.email || '',
        role: user.role,
      });
    } else {
      setEditingUser(null);
      setFormData({
        username: '',
        password: '',
        name: '',
        email: '',
        role: 'karyawan',
      });
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingUser(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const data = {
        username: formData.username,
        name: formData.name,
        email: formData.email,
        role: formData.role,
      };

      // Only include password for new user or when changing it
      if (formData.password) {
        data.password = formData.password;
      }

      if (editingUser) {
        await userAPI.update(editingUser.id, data);
      } else {
        await userAPI.create(data);
      }

      await fetchUsers();
      handleCloseModal();
    } catch (error) {
      console.error('Error saving user:', error);
      alert('Gagal menyimpan user: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Apakah Anda yakin ingin menghapus user ini?')) return;

    try {
      await userAPI.delete(id);
      await fetchUsers();
    } catch (error) {
      console.error('Error deleting user:', error);
      alert('Gagal menghapus user: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleToggleActive = async (id) => {
    try {
      await userAPI.toggleActive(id);
      await fetchUsers();
    } catch (error) {
      console.error('Error toggling user:', error);
      alert('Gagal mengubah status user');
    }
  };

  const getRoleBadge = (role) => {
    const badges = {
      admin: 'bg-purple-100 text-purple-800',
      hr: 'bg-blue-100 text-blue-800',
      finance: 'bg-green-100 text-green-800',
      karyawan: 'bg-gray-100 text-gray-800',
    };
    const labels = {
      admin: 'Admin',
      hr: 'HR',
      finance: 'Finance',
      karyawan: 'Karyawan',
    };
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${badges[role] || badges.karyawan}`}>
        {labels[role] || role}
      </span>
    );
  };

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
        <h1 className="text-2xl font-bold text-gray-800">Manajemen User</h1>
        <Button onClick={() => handleOpenModal()}>+ Tambah User</Button>
      </div>

      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Username</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {users.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-4 text-center text-gray-500">
                  Tidak ada data user
                </td>
              </tr>
            ) : (
              users.map((row) => (
                <tr key={row.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{row.username}</td>
                  <td className="px-6 py-4 text-sm text-gray-900">{row.name}</td>
                  <td className="px-6 py-4 text-sm text-gray-500">{row.email || '-'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">{getRoleBadge(row.role)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <button
                      onClick={() => handleToggleActive(row.id)}
                      className={`px-3 py-1 rounded text-xs font-medium ${
                        row.is_active
                          ? 'bg-green-100 text-green-800 hover:bg-green-200'
                          : 'bg-red-100 text-red-800 hover:bg-red-200'
                      }`}
                    >
                      {row.is_active ? 'Aktif' : 'Non-Aktif'}
                    </button>
                  </td>
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

      <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingUser ? 'Edit User' : 'Tambah User'}>
        <form onSubmit={handleSubmit}>
          <Input
            label="Username"
            value={formData.username}
            onChange={(e) => setFormData({ ...formData, username: e.target.value })}
            placeholder="Masukkan username"
            required
            disabled={!!editingUser}
          />
          <Input
            label="Password"
            type="password"
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            placeholder={editingUser ? 'Kosongkan jika tidak diubah' : 'Masukkan password'}
            required={!editingUser}
          />
          <Input
            label="Nama Lengkap"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="Masukkan nama lengkap"
            required
          />
          <Input
            label="Email"
            type="email"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            placeholder="Masukkan email"
          />
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Role</label>
            <select
              value={formData.role}
              onChange={(e) => setFormData({ ...formData, role: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="admin">Admin</option>
              <option value="hr">HR</option>
              <option value="finance">Finance</option>
              <option value="karyawan">Karyawan</option>
            </select>
          </div>
          <div className="flex justify-end gap-3 mt-6">
            <Button type="button" onClick={handleCloseModal} variant="secondary">Batal</Button>
            <Button type="submit">{editingUser ? 'Update' : 'Simpan'}</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default UserManagement;
