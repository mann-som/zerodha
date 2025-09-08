import React, { useState } from 'react';
import axios from 'axios';

function UserForm({ token, fetchUsers }) {
  const [form, setForm] = useState({ email: '', name: '', password: '', role: 'user', balance: '' });
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        ...form,
        balance: parseFloat(form.balance) || 0, 
      };
      const response = await axios.post('http://localhost:8080/api/users', payload, {
        headers: { Authorization: `Bearer ${token}` },
      });
      console.log('Create user response:', response.status, response.data);
      setForm({ email: '', name: '', password: '', role: 'user', balance: '' });
      setError('');
      fetchUsers();
    } catch (err) {
      console.error('Error creating user:', err.response?.status, err.response?.data || err.message);
      if (err.response?.status === 204) {
        setForm({ email: '', name: '', password: '', role: 'user', balance: '' });
        setError('');
        fetchUsers();
      } else {
        setError(err.response?.data?.error || 'Failed to create user: ' + err.message);
      }
    }
  };

  return (
    <div className="mb-8">
      <h2 className="text-xl font-bold mb-4">Create User</h2>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          type="email"
          placeholder="Email"
          value={form.email}
          onChange={(e) => setForm({ ...form, email: e.target.value })}
          className="border p-2 w-full rounded"
          required
        />
        <input
          type="text"
          placeholder="Name"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          className="border p-2 w-full rounded"
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
          className="border p-2 w-full rounded"
          required
        />
        <select
          value={form.role}
          onChange={(e) => setForm({ ...form, role: e.target.value })}
          className="border p-2 w-full rounded"
          required
        >
          <option value="user">User</option>
          <option value="admin">Admin</option>
        </select>
        <input
          type="number"
          placeholder="Balance"
          value={form.balance}
          onChange={(e) => setForm({ ...form, balance: e.target.value })}
          className="border p-2 w-full rounded"
          required
        />
        <button type="submit" className="bg-blue-500 text-white p-2 w-full rounded hover:bg-blue-600">
          Create User
        </button>
      </form>
    </div>
  );
}

export default UserForm;