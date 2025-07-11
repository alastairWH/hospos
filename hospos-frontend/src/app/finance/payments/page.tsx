"use client";
import { useEffect, useState } from "react";

interface Payment {
  id: string;
  sale_id: string;
  amount: number;
  method: string;
}

export default function PaymentsPage() {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/api/payments")
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setPayments(data);
        } else {
          setPayments([]);
        }
        setLoading(false);
      })
      .catch(() => {
        setError("Failed to load payments data");
        setLoading(false);
      });
  }, []);

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Payments</h1>
      {loading && <div>Loading...</div>}
      {error && <div className="text-red-500">{error}</div>}
      {!loading && !error && (
        <table className="w-full border rounded shadow bg-white dark:bg-gray-900">
          <thead>
            <tr className="bg-gray-100 dark:bg-gray-800">
              <th className="p-2 text-left">Payment ID</th>
              <th className="p-2 text-left">Sale ID</th>
              <th className="p-2 text-left">Amount (£)</th>
              <th className="p-2 text-left">Method</th>
            </tr>
          </thead>
          <tbody>
            {payments.map((p) => (
              <tr key={p.id} className="border-t">
                <td className="p-2">{p.id}</td>
                <td className="p-2">{p.sale_id}</td>
                <td className="p-2">{p.amount.toFixed(2)}</td>
                <td className="p-2">{p.method}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
