import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import UserTable from './UserTable';
import OrderTable from './OrderTable';
import OrderForm from './OrderForm';
import UserForm from './UserForm';

function Dashboard({ token, user, setToken, setUser }) {
  const [users, setUsers] = useState([]);
  const [orders, setOrders] = useState([]);
  const [error, setError] = useState('');

  const fetchUsers = useCallback(async () => {
    try {
      const response = await axios.get('http://localhost:8080/api/users', {
        headers: { Authorization: `Bearer ${token}` }
      });
      console.log('fetchUsers response:', response.status, response.data);
      setUsers(response.data || []);
      setError('');
    } catch (error) {
      console.error('Error fetching users:', error.response?.status, error.response?.data || error.message);
      setError(error.response?.data?.error || 'Failed to fetch users: ' + error.message);
      if (error.response?.status === 401) {
        setToken('');
        setUser(null);
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }
    }
  }, [token, setToken, setUser]);

  const fetchOrders = useCallback(async () => {
    try {
      const response = await axios.get('http://localhost:8080/api/orders', {
        headers: { Authorization: `Bearer ${token}` }
      });
      console.log('fetchOrders response:', response.status, response.data);
      setOrders(response.data || []);
      setError('');
    } catch (error) {
      console.error('Error fetching orders:', error.response?.status, error.response?.data || error.message);
      setError(error.response?.data?.error || 'Failed to fetch orders: ' + error.message);
      if (error.response?.status === 401) {
        setToken('');
        setUser(null);
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }
    }
  }, [token, setToken, setUser]);

  useEffect(() => {
    if (token) {
      fetchUsers();
      fetchOrders();
    }
  }, [token, fetchUsers, fetchOrders]);

  const handleLogout = () => {
    setToken('');
    setUser(null);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  };

  return (
    <div className="container mx-auto p-4">
      <header className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Zerodha Clone - Dashboard</h1>
        <div>
          <span className="mr-4">Logged in as: {user.email} ({user.role})</span>
          <button onClick={handleLogout} className="bg-red-500 text-white p-2 rounded hover:bg-red-600">
            Logout
          </button>
        </div>
      </header>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      {user.role === 'admin' && <UserForm token={token} fetchUsers={fetchUsers} />}
      <OrderForm token={token} userId={user.user_id} fetchOrders={fetchOrders} />
      <OrderTable orders={orders} isAdmin={user.role === 'admin'} fetchOrders={fetchOrders} token={token} />
      {user.role === 'admin' && <UserTable users={users} fetchUsers={fetchUsers} token={token} />}
    </div>
  );
}

export default Dashboard;