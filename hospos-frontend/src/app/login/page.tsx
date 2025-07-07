"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { setAuth } from "../auth";

export default function LoginPage() {
  const [name, setName] = useState("");
  const [pin, setPin] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    if (!name || !pin) {
      setError("Name and PIN required");
      return;
    }
    if (pin.length < 3 || pin.length > 6 || !/^[0-9]+$/.test(pin)) {
      setError("PIN must be 3-6 digits");
      return;
    }
    try {
      const res = await fetch("http://localhost:8080/api/auth", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, pin }),
      });
      if (!res.ok) {
        setError("Invalid credentials");
        return;
      }
      const data = await res.json();
      setAuth(data.token, data.role);
      router.push("/dashboard");
    } catch {
      setError("Server error");
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 dark:bg-gray-900">
      <form
        className="bg-white dark:bg-gray-800 p-8 rounded shadow-md w-full max-w-sm flex flex-col gap-4"
        onSubmit={handleSubmit}
      >
        <h1 className="text-2xl font-bold mb-2 text-center">Login</h1>
        <input
          type="text"
          placeholder="Name"
          value={name}
          onChange={e => setName(e.target.value)}
          className="border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <input
          type="password"
          placeholder="PIN (3-6 digits)"
          value={pin}
          onChange={e => setPin(e.target.value)}
          className="border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          maxLength={6}
        />
        {error && <div className="text-red-600 text-sm">{error}</div>}
        <button
          type="submit"
          className="bg-blue-600 text-white rounded px-4 py-2 font-semibold hover:bg-blue-700 transition"
        >
          Sign In
        </button>
      </form>
    </div>
  );
}
