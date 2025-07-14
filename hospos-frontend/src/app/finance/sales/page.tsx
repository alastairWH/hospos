"use client";
import { useEffect, useState } from "react";

interface SaleProduct {
  product_id: string;
  name: string;
  qty: number;
  price: number;
}

interface SalePayment {
  amount: number;
  method: string;
}

interface Sale {
  id: string;
  products: SaleProduct[];
  total: number;
  vat: number;
  discount: number;
  paid: number;
  payments: SalePayment[];
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
              <th className="p-2 text-left">Products</th>
              <th className="p-2 text-left">Total (£)</th>
              <th className="p-2 text-left">VAT (£)</th>
              <th className="p-2 text-left">Discount (£)</th>
              <th className="p-2 text-left">Paid (£)</th>
              <th className="p-2 text-left">Payments</th>
            </tr>
          </thead>
          <tbody>
            {sales.map((sale) => (
              <tr key={sale.id} className="border-t align-top">
                <td className="p-2">{sale.id}</td>
                <td className="p-2">
                  {Array.isArray(sale.products) && sale.products.length > 0 ? (
                    <ul className="list-disc pl-4">
                      {sale.products.map((p, idx) => (
                        <li key={idx}>
                          {p.name} x{p.qty} (£{p.price.toFixed(2)})
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <span className="text-gray-400">No products</span>
                  )}
                </td>
                <td className="p-2">{sale.total?.toFixed(2) ?? "0.00"}</td>
                <td className="p-2">{sale.vat?.toFixed(2) ?? "0.00"}</td>
                <td className="p-2">{sale.discount?.toFixed(2) ?? "0.00"}</td>
                <td className="p-2">{sale.paid?.toFixed(2) ?? "0.00"}</td>
                <td className="p-2">
                  {Array.isArray(sale.payments) && sale.payments.length > 0 ? (
                    <ul className="list-disc pl-4">
                      {sale.payments.map((pay, idx) => (
                        <li key={idx}>
                          {pay.method}: £{pay.amount.toFixed(2)}
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <span className="text-gray-400">No payments</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
