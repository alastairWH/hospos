"use client";
import { useEffect, useState } from "react";

interface Sale {
  id: string;
  product_id: string;
  quantity: number;
  total: number;
  vat: number;
}

export default function SalesReportPage() {
  const [sales, setSales] = useState<Sale[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/api/sales")
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setSales(data);
        } else {
          setSales([]);
        }
        setLoading(false);
      })
      .catch(() => {
        setError("Failed to load sales data");
        setLoading(false);
      });
  }, []);

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Sales Reports</h1>
      {loading && <div>Loading...</div>}
      {error && <div className="text-red-500">{error}</div>}
      {!loading && !error && (
        <table className="w-full border rounded shadow bg-white dark:bg-gray-900">
          <thead>
            <tr className="bg-gray-100 dark:bg-gray-800">
              <th className="p-2 text-left">Sale ID</th>
              <th className="p-2 text-left">Product ID</th>
              <th className="p-2 text-left">Quantity</th>
              <th className="p-2 text-left">Total (£)</th>
              <th className="p-2 text-left">VAT (£)</th>
            </tr>
          </thead>
          <tbody>
            {sales.map((sale) => (
              <tr key={sale.id} className="border-t">
                <td className="p-2">{sale.id}</td>
                <td className="p-2">{sale.product_id}</td>
                <td className="p-2">{sale.quantity}</td>
                <td className="p-2">{sale.total.toFixed(2)}</td>
                <td className="p-2">{sale.vat.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
