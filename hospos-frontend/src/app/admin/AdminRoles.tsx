"use client";
import React, { useEffect, useState } from "react";

export default function AdminRoles() {
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
      })
      .catch(() => setError("Failed to add role"));
  }

  function handleDeleteRole(role: string) {
    fetch(`http://localhost:8080/api/roles/${role}`, { method: "DELETE" })
      .then(() => fetchRoles())
      .catch(() => setError("Failed to delete role"));
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-2">
        <h2 className="text-xl font-semibold">Roles</h2>
        <button className="bg-blue-600 text-white px-3 py-1 rounded" onClick={() => setShowAdd(true)}>Add Role</button>
      </div>
      {error && <div className="text-red-600 text-sm mb-2">{error}</div>}
      {loading ? <p>Loading...</p> : (
        <table className="min-w-full bg-white dark:bg-gray-800 rounded shadow text-sm">
          <thead>
            <tr>
              <th className="py-2 px-4 text-left">Role</th>
              <th className="py-2 px-4 text-left">Actions</th>
            </tr>
          </thead>
          <tbody>
            {roles.map((r) => (
              <tr key={r}>
                <td className="py-2 px-4">{r}</td>
                <td className="py-2 px-4">
                  <button className="text-red-600 underline" onClick={() => handleDeleteRole(r)}>Remove</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
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
