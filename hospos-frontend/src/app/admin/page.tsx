"use client";
import React, { useEffect, useState } from "react";
import ProtectedRoute from "../protected-route";
import AdminUsers from "./AdminUsers";
import AdminRoles from "./AdminRoles";
import AdminBusinessInfo from "./AdminBusinessInfo";
import Card from "../ui/Card";

export default function AdminPage() {
  const [roles, setRoles] = useState<string[]>([]);

  // Fetch roles from API
  const fetchRoles = () => {
    fetch("http://localhost:8080/api/roles")
      .then((res) => res.json())
      .then((data) => setRoles(data.map((r: any) => r.role)));
  };

  useEffect(() => { fetchRoles(); }, []);

  return (
    <ProtectedRoute allowedRoles={["admin"]}>
      <div className="min-h-screen p-8">
        <h1 className="text-3xl font-bold mb-6">Admin Panel</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card>
            <AdminUsers roles={roles} />
          </Card>
          <Card>
            <AdminRoles onRolesChanged={fetchRoles} />
          </Card>
        </div>
        <div className="mt-8">
          <Card>
            <AdminBusinessInfo />
          </Card>
        </div>
      </div>
    </ProtectedRoute>
  );
}
