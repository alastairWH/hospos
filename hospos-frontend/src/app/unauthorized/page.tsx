"use client";

import { useRouter } from "next/navigation";
import { clearAuth } from "../auth";

export default function UnauthorizedPage() {
  const router = useRouter();
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[#101624] text-white">
      <div className="bg-[#182033] p-8 rounded-2xl shadow-xl flex flex-col items-center gap-6">
        <h1 className="text-3xl font-bold">Unauthorized</h1>
        <p className="text-lg text-gray-300">You do not have permission to access this page.</p>
        <button
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-semibold transition"
          onClick={() => {
            clearAuth();
            router.replace("/login");
          }}
        >
          Back to Login
        </button>
      </div>
    </div>
  );
}
