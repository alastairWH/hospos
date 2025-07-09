"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { setAuth } from "../auth";
import Input from "../ui/Input";
import Button from "../ui/Button";
import Alert from "../ui/Alert";

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
      setAuth(data.token, data.role, data.name || name);
      router.push("/dashboard");
    } catch {
      setError("Server error");
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center">
      <form
        className="bg-white dark:bg-gray-800 p-8 rounded-2xl shadow-xl w-full max-w-sm flex flex-col gap-4 border border-gray-100 dark:border-gray-700"
        onSubmit={handleSubmit}
      >
        <h1 className="text-3xl font-bold mb-2 text-center">Login</h1>
        <Input
          type="text"
          placeholder="Name"
          value={name}
          onChange={e => setName(e.target.value)}
          autoFocus
        />
        <Input
          type="password"
          placeholder="PIN (3-6 digits)"
          value={pin}
          onChange={e => setPin(e.target.value)}
          maxLength={6}
        />
        {error && <Alert type="error">{error}</Alert>}
        <Button type="submit">Sign In</Button>
      </form>
    </div>
  );
}
