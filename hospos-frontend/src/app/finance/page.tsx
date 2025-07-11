"use client";
import Link from "next/link";

export default function FinanceDashboard() {
  return (
    <div className="p-8 max-w-5xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">Finance Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Link href="/finance/sales" className="block bg-blue-100 dark:bg-blue-900 rounded-lg p-6 shadow hover:bg-blue-200 dark:hover:bg-blue-800 transition">
          <h2 className="text-xl font-semibold mb-2">Sales Reports</h2>
          <p>View and analyze sales by day, product, or staff.</p>
        </Link>
        <Link href="/finance/payments" className="block bg-green-100 dark:bg-green-900 rounded-lg p-6 shadow hover:bg-green-200 dark:hover:bg-green-800 transition">
          <h2 className="text-xl font-semibold mb-2">Payments</h2>
          <p>See all payment transactions and methods.</p>
        </Link>
        <Link href="/finance/tax" className="block bg-yellow-100 dark:bg-yellow-900 rounded-lg p-6 shadow hover:bg-yellow-200 dark:hover:bg-yellow-800 transition">
          <h2 className="text-xl font-semibold mb-2">Tax Reports</h2>
          <p>Review VAT and tax breakdowns for accounting.</p>
        </Link>
        <Link href="/finance/receipts" className="block bg-purple-100 dark:bg-purple-900 rounded-lg p-6 shadow hover:bg-purple-200 dark:hover:bg-purple-800 transition">
          <h2 className="text-xl font-semibold mb-2">Receipts</h2>
          <p>Search and view past receipts and invoices.</p>
        </Link>
      </div>
    </div>
  );
}
