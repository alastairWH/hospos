"use client";
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";

interface Customer {
  id: string;
  name: string;
  email: string;
  phone: string;
  notes: string;
  tags: string[];
  createdAt: string;
}

export default function CustomerDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = params?.id as string;
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    fetch(`http://localhost:8080/api/customers/${id}`)
      .then((res) => {
        if (!res.ok) throw new Error("Not found");
        return res.json();
      })
      .then((data) => {
        setCustomer(data);
        setLoading(false);
      })
      .catch(() => {
        setError("Customer not found");
        setLoading(false);
      });
  }, [id]);

  return (
    <div className="max-w-2xl mx-auto p-6">
      <Link href="/customers" className="text-blue-600 hover:underline">&larr; Back to Customers</Link>
      {loading ? (
        <div>Loading...</div>
      ) : error ? (
        <div className="text-red-500 mt-4">{error}</div>
      ) : customer ? (
        <div className="mt-4">
          <h1 className="text-2xl font-bold mb-2">{customer.name}</h1>
          <div className="mb-2 text-gray-700 dark:text-gray-300">Email: {customer.email || <span className="italic">N/A</span>}</div>
          <div className="mb-2 text-gray-700 dark:text-gray-300">Phone: {customer.phone || <span className="italic">N/A</span>}</div>
          <div className="mb-2 text-gray-700 dark:text-gray-300">Tags: {customer.tags?.length ? customer.tags.join(", ") : <span className="italic">None</span>}</div>
          <div className="mb-2 text-gray-700 dark:text-gray-300">Created: {new Date(customer.createdAt).toLocaleString()}</div>
          <div className="mb-4 text-gray-700 dark:text-gray-300">Notes: {customer.notes || <span className="italic">None</span>}</div>
          {/* Booking history and actions can go here */}
        </div>
      ) : null}
    </div>
  );
}
