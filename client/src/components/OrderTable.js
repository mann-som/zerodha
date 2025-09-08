import React from 'react';
import axios from 'axios';

function OrderTable({ orders, isAdmin, fetchOrders, token }) {
  const deleteOrder = async (id) => {
    try {
      await axios.delete(`http://localhost:8080/api/orders/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      fetchOrders();
    } catch (err) {
      alert('Error deleting order: ' + (err.response?.data?.error || 'Unknown error'));
    }
  };

  return (
    <div className="mb-8">
      <h2 className="text-xl font-bold mb-4">Order Book</h2>
      <table className="w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th className="border p-2">ID</th>
            <th className="border p-2">User ID</th>
            <th className="border p-2">Symbol</th>
            <th className="border p-2">Side</th>
            <th className="border p-2">Quantity</th>
            <th className="border p-2">Price</th>
            <th className="border p-2">Status</th>
            {isAdmin && <th className="border p-2">Actions</th>}
          </tr>
        </thead>
        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td className="border p-2">{order.id}</td>
              <td className="border p-2">{order.user_id}</td>
              <td className="border p-2">{order.symbol}</td>
              <td className="border p-2">{order.side}</td>
              <td className="border p-2">{order.quantity}</td>
              <td className="border p-2">{order.price}</td>
              <td className="border p-2">{order.status}</td>
              {isAdmin && (
                <td className="border p-2">
                  <button
                    onClick={() => deleteOrder(order.id)}
                    className="bg-red-500 text-white p-1 rounded hover:bg-red-600"
                  >
                    Delete
                  </button>
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default OrderTable;