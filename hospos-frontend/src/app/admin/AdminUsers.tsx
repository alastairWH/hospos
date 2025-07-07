"use client";
import React, { useEffect, useState } from "react";

export type User = {
  id: string;
  name: string;
  role: string;
  pin?: string;
};

export default function AdminUsers({ roles }: { roles: string[] }) {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showAdd, setShowAdd] = useState(false);
  const [showReset, setShowReset] = useState<string | null>(null);
  const [newUser, setNewUser] = useState({ name: "", pin: "", role: roles[0] || "" });

  function fetchUsers() {
    setLoading(true);
    fetch("http://localhost:8080/api/users")
      .then((res) => res.json())
      .then((data) => { setUsers(data); setLoading(false); })
      .catch(() => { setError("Failed to load users"); setLoading(false); });
  }

  useEffect(() => { fetchUsers(); }, []);

  function handleAddUser(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    if (!newUser.name || !newUser.pin || !newUser.role) {
      setError("All fields required");
      return;
    }
    if (newUser.pin.length < 3 || newUser.pin.length > 6 || !/^[0-9]+$/.test(newUser.pin)) {
      setError("PIN must be 3-6 digits");
      return;
    }
    fetch("http://localhost:8080/api/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(newUser),
    })
      .then(async (res) => {
        if (!res.ok) {
          const errText = await res.text();
          console.error("Add user error:", errText);
          throw new Error(errText);
        }
        setShowAdd(false);
        setNewUser({ name: "", pin: "", role: roles[0] || "" });
        fetchUsers();
      })
      .catch((err) => {
        setError("Failed to add user");
        if (err instanceof Error) console.error("Add user exception:", err.message);
      });
  }

function handleDeleteUser(id: string, name: string) {
  if (!window.confirm(`Are you sure you want to delete user '${name}'? This action cannot be undone.`)) {
    return;
  }
  fetch(`http://localhost:8080/api/users/${id}`, { method: "DELETE" })
    .then(async (res) => {
      if (!res.ok) {
        const errText = await res.text();
        setError("Failed to delete user");
        console.error("Delete user error:", errText);
        return;
      }
      fetchUsers();
    })
    .catch((err) => {
      setError("Failed to delete user");
      if (err instanceof Error) console.error("Delete user exception:", err.message);
    });
}

  function handleResetPin(id: string, pin: string) {
    if (pin.length < 3 || pin.length > 6 || !/^[0-9]+$/.test(pin)) {
      setError("PIN must be 3-6 digits");
      return;
    }
    fetch(`http://localhost:8080/api/users/${id}/pin`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ pin }),
    })
      .then(() => { setShowReset(null); fetchUsers(); })
      .catch(() => setError("Failed to reset PIN"));
  }

  function handleUpdateRole(id: string, role: string) {
    fetch(`http://localhost:8080/api/users/${id}/role`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ role }),
    })
      .then(() => fetchUsers())
      .catch(() => setError("Failed to update role"));
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-2">
        <h2 className="text-xl font-semibold">Users</h2>
        <button
          className="bg-blue-600 text-white px-3 py-1 rounded disabled:opacity-50"
          onClick={() => setShowAdd(true)}
          disabled={roles.length === 0}
          title={roles.length === 0 ? "Add at least one role first" : undefined}
        >
          Add User
        </button>
      </div>
      {error && <div className="text-red-600 text-sm mb-2">{error}</div>}
      {loading ? <p>Loading...</p> : (
        <table className="min-w-full bg-white dark:bg-gray-800 rounded shadow text-sm">
          <thead>
            <tr>
              <th className="py-2 px-4 text-left">Name</th>
              <th className="py-2 px-4 text-left">Role</th>
              <th className="py-2 px-4 text-left">Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((u) => (
              <tr key={u.id}>
                <td className="py-2 px-4">{u.name}</td>
                <td className="py-2 px-4">
                  <select
                    value={u.role}
                    onChange={e => handleUpdateRole(u.id, e.target.value)}
                    className="border rounded px-2 py-1"
                    aria-label="User role"
                  >
                    {roles.map(r => <option key={r} value={r}>{r}</option>)}
                  </select>
                </td>
                <td className="py-2 px-4 flex gap-2">
                  <button className="text-blue-600 underline" onClick={() => setShowReset(u.id)}>Reset PIN</button>
                  <button className="text-red-600 underline" onClick={() => handleDeleteUser(u.id, u.name)}>Remove</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {showAdd && roles.length > 0 && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <form className="bg-white dark:bg-gray-800 p-6 rounded shadow flex flex-col gap-3 min-w-[300px]" onSubmit={handleAddUser}>
            <h3 className="text-lg font-bold mb-2">Add User</h3>
            <input placeholder="Name" className="border rounded px-2 py-1" value={newUser.name} onChange={e => setNewUser({ ...newUser, name: e.target.value })} />
            <input placeholder="PIN (3-6 digits)" className="border rounded px-2 py-1" value={newUser.pin} onChange={e => setNewUser({ ...newUser, pin: e.target.value })} maxLength={6} />
            <select
              className="border rounded px-2 py-1"
              value={newUser.role}
              onChange={e => setNewUser({ ...newUser, role: e.target.value })}
              aria-label="New user role"
            >
              {roles.map(r => <option key={r} value={r}>{r}</option>)}
            </select>
            <div className="flex gap-2 mt-2">
              <button type="submit" className="bg-blue-600 text-white px-3 py-1 rounded">Add</button>
              <button type="button" className="bg-gray-300 px-3 py-1 rounded" onClick={() => setShowAdd(false)}>Cancel</button>
            </div>
          </form>
        </div>
      )}
      {showReset && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <form className="bg-white dark:bg-gray-800 p-6 rounded shadow flex flex-col gap-3 min-w-[300px]" onSubmit={e => { e.preventDefault(); const pin = (e.target as any).pin.value; handleResetPin(showReset, pin); }}>
            <h3 className="text-lg font-bold mb-2">Reset PIN</h3>
            <input name="pin" placeholder="New PIN (3-6 digits)" className="border rounded px-2 py-1" maxLength={6} />
            <div className="flex gap-2 mt-2">
              <button type="submit" className="bg-blue-600 text-white px-3 py-1 rounded">Reset</button>
              <button type="button" className="bg-gray-300 px-3 py-1 rounded" onClick={() => setShowReset(null)}>Cancel</button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}
