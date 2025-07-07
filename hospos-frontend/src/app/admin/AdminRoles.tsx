"use client";
import React, { useEffect, useState } from "react";

export default function AdminRoles({ onRolesChanged }: { onRolesChanged?: () => void }) {
  const [roles, setRoles] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showAdd, setShowAdd] = useState(false);
  const [newRole, setNewRole] = useState("");

  function fetchRoles() {
    setLoading(true);
    fetch("http://localhost:8080/api/roles")
      .then((res) => res.json())
      .then((data) => { setRoles(data.map((r: any) => r.role)); setLoading(false); })
      .catch(() => { setError("Failed to load roles"); setLoading(false); });
  }

  useEffect(() => { fetchRoles(); }, []);

  function handleAddRole(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    if (!newRole) {
      setError("Role required");
      return;
    }
    fetch("http://localhost:8080/api/roles", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ role: newRole }),
    })
      .then((res) => {
        if (!res.ok) throw new Error();
        setShowAdd(false);
        setNewRole("");
        fetchRoles();
        if (onRolesChanged) onRolesChanged();
      })
      .catch(() => setError("Failed to add role"));
  }

  function handleDeleteRole(role: string) {
    fetch(`http://localhost:8080/api/roles/${role}`, { method: "DELETE" })
      .then(() => {
        fetchRoles();
        if (onRolesChanged) onRolesChanged();
      })
      .catch(() => setError("Failed to delete role"));
  }

  return (
    <div className="max-w-xl mx-auto p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold text-gray-800 dark:text-gray-100">Roles</h2>
        <button className="bg-blue-600 hover:bg-blue-700 transition text-white px-4 py-2 rounded-lg shadow" onClick={() => setShowAdd(true)}>Add Role</button>
      </div>
      {error && <div className="text-red-600 text-sm mb-2 rounded bg-red-50 p-2 border border-red-200">{error}</div>}
      {loading ? (
        <div className="flex justify-center items-center h-24">
          <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-600"></div>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4">
          {roles.map((r) => (
            <div key={r} className="bg-white dark:bg-gray-800 rounded-xl shadow p-4 flex items-center justify-between border border-gray-100 dark:border-gray-700">
              <span className="font-semibold text-lg text-gray-800 dark:text-gray-100">{r}</span>
              <button className="text-red-600 hover:underline" onClick={() => handleDeleteRole(r)}>Remove</button>
            </div>
          ))}
        </div>
      )}
      {showAdd && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <form className="bg-white dark:bg-gray-800 p-6 rounded shadow flex flex-col gap-3 min-w-[300px]" onSubmit={handleAddRole}>
            <h3 className="text-lg font-bold mb-2">Add Role</h3>
            <input placeholder="Role name" className="border rounded px-2 py-1" value={newRole} onChange={e => setNewRole(e.target.value)} />
            <div className="flex gap-2 mt-2">
              <button type="submit" className="bg-blue-600 text-white px-3 py-1 rounded">Add</button>
              <button type="button" className="bg-gray-300 px-3 py-1 rounded" onClick={() => setShowAdd(false)}>Cancel</button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}
