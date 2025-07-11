"use client";
import { useEffect, useState } from "react";

interface Sale {
  id: string;
  total: number;
  vat: number;
}

export default function TaxReportPage() {
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

  const totalVAT = sales.reduce((sum, s) => sum + (s.vat || 0), 0);
  const totalSales = sales.reduce((sum, s) => sum + (s.total || 0), 0);

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Tax Report</h1>
      {loading && <div>Loading...</div>}
      {error && <div className="text-red-500">{error}</div>}
      {!loading && !error && (
        <div className="bg-white dark:bg-gray-900 rounded shadow p-6">
          <div className="mb-4 text-lg">Total Sales: <span className="font-semibold">£{totalSales.toFixed(2)}</span></div>
          <div className="mb-4 text-lg">Total VAT: <span className="font-semibold">£{totalVAT.toFixed(2)}</span></div>
          <table className="w-full mt-6 border rounded">
            <thead>
              <tr className="bg-gray-100 dark:bg-gray-800">
                <th className="p-2 text-left">Sale ID</th>
                <th className="p-2 text-left">Total (£)</th>
                <th className="p-2 text-left">VAT (£)</th>
              </tr>
            </thead>
            <tbody>
              {sales.map((s) => (
                <tr key={s.id} className="border-t">
                  <td className="p-2">{s.id}</td>
                  <td className="p-2">{s.total.toFixed(2)}</td>
                  <td className="p-2">{s.vat.toFixed(2)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
