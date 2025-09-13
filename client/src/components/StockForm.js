import React, { useState } from 'react';
import axios from 'axios';

function StockForm({ token, fetchStocks }) {
  const [form, setForm] = useState({ symbol: '', description: '', initial_price: '' });
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        ...form,
        initial_price: parseFloat(form.initial_price) || 0, // Convert to float64
      };
      const response = await axios.post('http://localhost:8080/api/stocks', payload, {
        headers: { Authorization: `Bearer ${token}` },
      });
      console.log('Create stock response:', response.status, response.data);
      setForm({ symbol: '', description: '', initial_price: '' });
      setError('');
      fetchStocks();
    } catch (err) {
      console.error('Error creating stock:', err.response?.status, err.response?.data || err.message);
      setError(err.response?.data?.error || 'Failed to create stock: ' + err.message);
    }
  };

  return (
    <div className="mb-8">
      <h2 className="text-xl font-bold mb-4">List New Stock</h2>
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
        <input
          type="text"
          placeholder="Description (e.g., Reliance Industries Ltd)"
          value={form.description}
          onChange={(e) => setForm({ ...form, description: e.target.value })}
          className="border p-2 w-full rounded"
          required
        />
        <input
          type="number"
          placeholder="Initial Price"
          value={form.initial_price}
          onChange={(e) => setForm({ ...form, initial_price: e.target.value })}
          className="border p-2 w-full rounded"
          step="0.01"
          min="0"
          required
        />
        <button type="submit" className="bg-blue-500 text-white p-2 w-full rounded hover:bg-blue-600">
          Create Stock
        </button>
      </form>
    </div>
  );
}

export default StockForm;