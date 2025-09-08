import React, { useState } from 'react';
import axios from 'axios';

function OrderForm({ token, userId, fetchOrders }) {
  const [form, setForm] = useState({ symbol: '', side: 'buy', quantity: '', price: '' });
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        ...form,
        quantity: parseInt(form.quantity, 10) || 0, 
        price: parseFloat(form.price) || 0, 
        user_id: userId,
      };
      const response = await axios.post('http://localhost:8080/api/orders', payload, {
        headers: { Authorization: `Bearer ${token}` },
      });
      console.log('Create order response:', response.status, response.data);
      setForm({ symbol: '', side: 'buy', quantity: '', price: '' });
      setError('');
      fetchOrders();
    } catch (err) {
      console.error('Error creating order:', err.response?.status, err.response?.data || err.message);
      if (err.response?.status === 204) {
        setForm({ symbol: '', side: 'buy', quantity: '', price: '' });
        setError('');
        fetchOrders();
      } else {
        setError(err.response?.data?.error || 'Failed to create order: ' + err.message);
      }
    }
  };

  return (
    <div className="mb-8">
      <h2 className="text-xl font-bold mb-4">Create Order</h2>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          type="text"
          placeholder="Symbol (e.g., RELIANCE)"
          value={form.symbol}
          onChange={(e) => setForm({ ...form, symbol: e.target.value })}
          className="border p-2 w-full rounded"
          required
        />
        <select
          value={form.side}
          onChange={(e) => setForm({ ...form, side: e.target.value })}
          className="border p-2 w-full rounded"
          required
        >
          <option value="buy">Buy</option>
          <option value="sell">Sell</option>
        </select>
        <input
          type="number"
          placeholder="Quantity"
          value={form.quantity}
          onChange={(e) => setForm({ ...form, quantity: e.target.value })}
          className="border p-2 w-full rounded"
          required
        />
        <input
          type="number"
          placeholder="Price"
          value={form.price}
          onChange={(e) => setForm({ ...form, price: e.target.value })}
          className="border p-2 w-full rounded"
          required
        />
        <button type="submit" className="bg-blue-500 text-white p-2 w-full rounded hover:bg-blue-600">
          Create Order
        </button>
      </form>
    </div>
  );
}

export default OrderForm;