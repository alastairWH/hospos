"use client";
import Link from "next/link";
import { useRouter, usePathname } from "next/navigation";
import { getAuth, clearAuth } from "./auth";
import { useEffect, useState } from "react";

export default function Navbar() {
  const [financeOpen, setFinanceOpen] = useState(false);
  const [role, setRole] = useState<string | null>(null);
  const [username, setUsername] = useState<string | null>(null);
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    const { role, username } = getAuth();
    setRole(role);
    setUsername(username);
  }, [pathname]);

  function handleLogout() {
    clearAuth();
    router.replace("/login");
  }

  const isLoggedIn = !!role || !!username;
  if (pathname === "/unauthorized") {
    return (
      <nav className="bg-white dark:bg-gray-900 shadow rounded-b-xl px-4 py-3 flex items-center">
        <Link href="/dashboard" className="flex items-center">
          <img src="/logo-hospos.svg" alt="HOSPOS Logo" className="h-32 w-auto" />
        </Link>
      </nav>
    );
  }
  return (
    <nav className="bg-white dark:bg-gray-900 shadow rounded-b-xl px-8 py-4 flex items-center justify-between min-h-[88px]">
      <div className="flex items-center gap-10">
        <Link href="/dashboard" className="flex items-center mr-10">
          <img src="/logo-hospos.svg" alt="HOSPOS Logo" className="h-20 w-auto" />
        </Link>
        {isLoggedIn && (
          <div className="flex items-center gap-8">
            <Link href="/products" className="hover:underline text-gray-700 dark:text-gray-200">Products</Link>
            <Link href="/discounts" className="hover:underline text-gray-700 dark:text-gray-200">Discounts</Link>
            <Link href="/tills" className="hover:underline text-gray-700 dark:text-gray-200">Tills</Link>
            <Link href="/customers" className="hover:underline text-gray-700 dark:text-gray-200">Customers</Link>
            <Link href="/bookings" className="hover:underline text-gray-700 dark:text-gray-200">Bookings</Link>
            {/* Finance Dropdown */}
            <div className="relative" onMouseLeave={() => setFinanceOpen(false)}>
              <button
                className="hover:underline text-gray-700 dark:text-gray-200 flex items-center gap-1 px-2 py-1 rounded focus:outline-none"
                onClick={() => setFinanceOpen((v) => !v)}
                type="button"
                aria-haspopup="true"
                aria-expanded={financeOpen}
              >
                Finance
                <svg className={`w-4 h-4 ml-1 transition-transform ${financeOpen ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M19 9l-7 7-7-7" /></svg>
              </button>
              {financeOpen && (
                <div className="absolute right-0 mt-2 w-40 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded shadow-lg z-50">
                  <Link href="/finance/reports" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-200">Reports</Link>
                  <Link href="/finance/payments" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-200">Payments</Link>

                </div>
              )}
            </div>

            {role === "admin" && (
              <Link href="/admin" className="hover:underline text-gray-700 dark:text-gray-200">Admin</Link>
            )}
          </div>
        )}
      </div>
      <div className="flex items-center gap-4">
        {isLoggedIn && (
          <>
            <span className="text-gray-500 dark:text-gray-300 text-sm">
              {username && <span>User: <span className="font-semibold">{username}</span></span>}
              {username && role && <span> &middot; </span>}
              {role && <span>Role: <span className="font-semibold">{role}</span></span>}
            </span>
            <button
              onClick={handleLogout}
              className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg shadow transition"
            >
              Logout
            </button>
          </>
        )}
      </div>
    </nav>
  );
}
