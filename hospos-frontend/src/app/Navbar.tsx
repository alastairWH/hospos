"use client";
import Link from "next/link";
import { useRouter, usePathname } from "next/navigation";
import { getAuth, clearAuth } from "./auth";
import { useEffect, useState } from "react";

export default function Navbar() {
  const [role, setRole] = useState<string | null>(null);
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    const { role } = getAuth();
    setRole(role);
  }, [pathname]);

  function handleLogout() {
    clearAuth();
    router.replace("/login");
  }

  return (
    <nav className="bg-white dark:bg-gray-900 shadow mb-6 rounded-b-xl px-4 py-3 flex items-center justify-between">
      <div className="flex items-center gap-6">
        <Link href="/dashboard" className="font-bold text-lg text-blue-700 dark:text-blue-300">HOSPOS</Link>
        <Link href="/products" className="hover:underline text-gray-700 dark:text-gray-200">Products</Link>
        {role === "admin" && (
          <Link href="/admin" className="hover:underline text-gray-700 dark:text-gray-200">Admin</Link>
        )}
      </div>
      <div className="flex items-center gap-4">
        {role && (
          <span className="text-gray-500 dark:text-gray-300 text-sm">Role: <span className="font-semibold">{role}</span></span>
        )}
        <button
          onClick={handleLogout}
          className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg shadow transition"
        >
          Logout
        </button>
      </div>
    </nav>
  );
}
