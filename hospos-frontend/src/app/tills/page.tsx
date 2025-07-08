"use client";
import React, { useEffect, useState } from "react";
import ProtectedRoute from "../protected-route";
import Card from "../ui/Card";
import Button from "../ui/Button";
import Input from "../ui/Input";
import Alert from "../ui/Alert";

export default function TillsPage() {
  const [locations, setLocations] = useState<any[]>([]);
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [linkCode, setLinkCode] = useState<string | null>(null);

  function fetchLocations() {
    setLoading(true);
    fetch("http://localhost:8080/api/locations")
      .then((res) => res.json())
      .then((data) => {
        setLocations(Array.isArray(data) ? data : []);
        setLoading(false);
      });
  }

  useEffect(() => { fetchLocations(); }, []);

  function handleAddLocation(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    if (!name) {
      setError("Location name required");
      return;
    }
    fetch("http://localhost:8080/api/locations", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    })
      .then((res) => res.json())
      .then((data) => {
        setName("");
        setLinkCode(data.linkCode);
        fetchLocations();
      })
      .catch(() => setError("Failed to add location"));
  }



  return (
    <ProtectedRoute allowedRoles={["admin", "manager"]}>
      <div className="min-h-screen p-8 bg-gray-50 dark:bg-gray-900">
        <h1 className="text-3xl font-bold mb-6">Tills / Locations</h1>
        <form className="flex gap-3 mb-6" onSubmit={handleAddLocation}>
          <Input placeholder="Location name" value={name} onChange={e => setName(e.target.value)} />
          <Button type="submit">Add Location</Button>
        </form>
        {error && <Alert type="error">{error}</Alert>}
        {linkCode && (
          <Alert type="success">Linking code: <span className="font-mono">{linkCode}</span></Alert>
        )}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mt-6">
          {loading ? (
            <div className="flex justify-center items-center h-32">
              <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-600"></div>
            </div>
          ) : locations.length === 0 ? (
            <Alert type="info">No locations found.</Alert>
          ) : (
            locations.map((loc) => (
              <Card key={loc.id}>
                <div className="flex flex-col gap-2">
                  <span className="font-semibold text-lg text-gray-800 dark:text-gray-100">{loc.name}</span>
                  <span className="text-xs text-gray-500 dark:text-gray-300">ID: {loc.id}</span>
                  <span className="text-xs text-gray-500 dark:text-gray-300">Linking code: <span className="font-mono">{loc.linkCode}</span></span>
                </div>
              </Card>
            ))
          )}
        </div>
      </div>
    </ProtectedRoute>
  );
}
