import React from "react";
import ProtectedRoute from "../protected-route";
import Card from "../ui/Card";

export default function DashboardPage() {
  return (
    <ProtectedRoute allowedRoles={["admin", "manager"]}>
      <div className="min-h-screen p-8">
        <h1 className="text-3xl font-bold mb-6">Dashboard</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card>
            <h2 className="text-xl font-semibold mb-2">Sales Overview</h2>
            <p className="text-gray-500">(Connect to /api/sales for data)</p>
          </Card>
          <Card>
            <h2 className="text-xl font-semibold mb-2">Bookings</h2>
            <p className="text-gray-500">(Connect to /api/bookings for data)</p>
          </Card>
        </div>
      </div>
    </ProtectedRoute>
  );
}
