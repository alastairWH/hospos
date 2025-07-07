import React from "react";
import ProtectedRoute from "../protected-route";

export default function AdminPage() {
  return (
    <ProtectedRoute allowedRoles={["admin"]}>
      <div className="min-h-screen p-8 bg-gray-50 dark:bg-gray-900">
        <h1 className="text-3xl font-bold mb-6">Admin Panel</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white dark:bg-gray-800 p-6 rounded shadow">
            <h2 className="text-xl font-semibold mb-2">Users</h2>
            <p className="text-gray-500">(Connect to /api/users for data)</p>
          </div>
          <div className="bg-white dark:bg-gray-800 p-6 rounded shadow">
            <h2 className="text-xl font-semibold mb-2">Roles</h2>
            <p className="text-gray-500">(Connect to /api/roles for data)</p>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
