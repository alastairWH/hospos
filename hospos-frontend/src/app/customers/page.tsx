"use client";
import { useEffect, useState } from "react";

import Link from "next/link";
import CustomerAddModal from "./CustomerAddModal";

interface Customer {
  id: string;
  name: string;
  email: string;
  phone: string;
  notes: string;
  tags: string[];
  createdAt: string;
}

export default function CustomersPage() {
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [error, setError] = useState("");
  const [showAdd, setShowAdd] = useState(false);

  function fetchCustomers() {
    setLoading(true);
    fetch(`http://localhost:8080/api/customers${search ? `?q=${encodeURIComponent(search)}` : ""}`)
      .then((res) => res.json())
      .then((data) => {
        if (!Array.isArray(data)) {
          setCustomers([]);
          setError(""); // No error, just treat as empty
        } else {
          setCustomers(data);
          setError("");
        }
        setLoading(false);
      })
      .catch((e) => {
        setError("Failed to load customers");
        setLoading(false);
      });
  }
  useEffect(() => {
    fetchCustomers();
    // eslint-disable-next-line
  }, [search]);

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Customers</h1>
      <div className="flex gap-2 mb-4">
        <input
          className="border rounded px-2 py-1 flex-1"
          placeholder="Search by name, email, or phone..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <button
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          onClick={() => setShowAdd(true)}
        >
          Add Customer
        </button>
      </div>
      <CustomerAddModal open={showAdd} onClose={() => setShowAdd(false)} onAdded={fetchCustomers} />
      {loading ? (
        <div>Loading...</div>
      ) : error ? (
        <div className="text-red-500">{error}</div>
      ) : customers.length === 0 ? (
        <div>No customers available.</div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full border">
            <thead>
              <tr className="bg-gray-100 dark:bg-gray-800">
                <th className="p-2 text-left">Name</th>
                <th className="p-2 text-left">Email</th>
                <th className="p-2 text-left">Phone</th>
                <th className="p-2 text-left">Tags</th>
                <th className="p-2 text-left">Created</th>
                <th className="p-2"></th>
              </tr>
            </thead>
            <tbody>
              {customers.map((c) => (
                <tr key={c.id} className="border-b hover:bg-gray-50 dark:hover:bg-gray-900">
                  <td className="p-2 font-medium">
                    <Link href={`/customers/${c.id}`} className="text-blue-700 hover:underline">{c.name}</Link>
                  </td>
                  <td className="p-2">{c.email}</td>
                  <td className="p-2">{c.phone}</td>
                  <td className="p-2">{c.tags?.join(", ")}</td>
                  <td className="p-2">{new Date(c.createdAt).toLocaleDateString()}</td>
                  <td className="p-2">
                    <Link href={`/customers/${c.id}`} className="text-blue-600 hover:underline">View</Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
