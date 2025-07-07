"use client";
import React, { useEffect, useState } from "react";
import ProtectedRoute from "../protected-route";
import AdminUsers from "./AdminUsers";
import AdminRoles from "./AdminRoles";

export default function AdminPage() {
  const [roles, setRoles] = useState<string[]>([]);
  useEffect(() => {
    fetch("http://localhost:8080/api/roles")
      .then((res) => res.json())
      .then((data) => setRoles(data.map((r: any) => r.role)));
  }, []);
  return (
    <ProtectedRoute allowedRoles={["admin"]}>
      <div className="min-h-screen p-8 bg-gray-50 dark:bg-gray-900">
        <h1 className="text-3xl font-bold mb-6">Admin Panel</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white dark:bg-gray-800 p-6 rounded shadow">
            <AdminUsers roles={roles} />
          </div>
          <div className="bg-white dark:bg-gray-800 p-6 rounded shadow">
            <AdminRoles />
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
