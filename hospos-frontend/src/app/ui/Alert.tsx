"use client";
import React from "react";

export default function Alert({ type = "info", children }: { type?: "info" | "error" | "success"; children: React.ReactNode }) {
  const color =
    type === "error"
      ? "bg-red-50 border-red-200 text-red-700"
      : type === "success"
      ? "bg-green-50 border-green-200 text-green-700"
      : "bg-blue-50 border-blue-200 text-blue-700";
  return (
    <div className={`rounded-lg border px-4 py-2 mb-2 text-sm ${color}`}>{children}</div>
  );
}
