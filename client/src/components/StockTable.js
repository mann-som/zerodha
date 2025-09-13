import React, { useState } from 'react';
import axios from 'axios';

function StockTable({ stocks, isAdmin, token, fetchStocks, setSelectedStock }) {
  const [error, setError] = useState('');

  const deleteStock = async (id) => {
    try {
      const response = await axios.delete(`http://localhost:8080/api/stocks/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      console.log('Delete stock response:', response.status, response.data);
      setError('');
      fetchStocks();
    } catch (err) {
      console.error('Error deleting stock:', err.response?.status, err.response?.data || err.message);
      setError(err.response?.data?.error || `Failed to delete stock: ${err.message}`);
    }
  };

  return (
    <div className="mb-8">
      <h2 className="text-xl font-bold mb-4">Listed Stocks</h2>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      <table className="w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th className="border p-2">Symbol</th>
            <th className="border p-2">Description</th>
            <th className="border p-2">Current Price</th>
            <th className="border p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {stocks.map((stock) => (
            <tr key={stock.id}>
              <td className="border p-2">{stock.symbol}</td>
              <td className="border p-2">{stock.description}</td>
              <td className="border p-2">{stock.current_price}</td>
              <td className="border p-2">
                <button
                  onClick={() => setSelectedStock({ symbol: stock.symbol, price: stock.current_price, side: 'buy' })}
                  className="bg-green-500 text-white p-1 rounded mr-2 hover:bg-green-600"
                >
                  Buy
                </button>
                <button
                  onClick={() => setSelectedStock({ symbol: stock.symbol, price: stock.current_price, side: 'sell' })}
                  className="bg-red-500 text-white p-1 rounded hover:bg-red-600"
                >
                  Sell
                </button>
                {isAdmin && (
                  <button
                    onClick={() => deleteStock(stock.id)}
                    className="bg-red-500 text-white p-1 rounded hover:bg-red-600 ml-2"
                  >
                    Delete
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default StockTable;