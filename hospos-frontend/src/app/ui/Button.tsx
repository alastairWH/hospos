"use client";
import React from "react";

export default function Button({
  children,
  className = "",
  ...props
}: React.ButtonHTMLAttributes<HTMLButtonElement> & { children: React.ReactNode }) {
  return (
    <button
      className={`bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg shadow transition font-medium focus:outline-none focus:ring-2 focus:ring-blue-400 disabled:opacity-50 ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}
