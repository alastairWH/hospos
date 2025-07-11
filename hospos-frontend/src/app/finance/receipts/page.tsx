"use client";
import { useEffect, useState } from "react";

interface Receipt {
  id: string;
  sale_id: string;
  createdAt?: string;
  // Add more fields as needed
}

export default function ReceiptsPage() {
  const [receipts, setReceipts] = useState<Receipt[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/api/receipts")
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setReceipts(data);
        } else {
          setReceipts([]);
        }
        setLoading(false);
      })
      .catch(() => {
        setError("Failed to load receipts data");
        setLoading(false);
      });
  }, []);

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Receipts</h1>
      {loading && <div>Loading...</div>}
      {error && <div className="text-red-500">{error}</div>}
      {!loading && !error && (
        <table className="w-full border rounded shadow bg-white dark:bg-gray-900">
          <thead>
            <tr className="bg-gray-100 dark:bg-gray-800">
              <th className="p-2 text-left">Receipt ID</th>
              <th className="p-2 text-left">Sale ID</th>
              <th className="p-2 text-left">Created At</th>
            </tr>
          </thead>
          <tbody>
            {receipts.map((r) => (
              <tr key={r.id} className="border-t">
                <td className="p-2">{r.id}</td>
                <td className="p-2">{r.sale_id}</td>
                <td className="p-2">{r.createdAt ? new Date(r.createdAt).toLocaleString() : "-"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
