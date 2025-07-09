"use client";


import React, { useEffect, useState } from "react";
import ProtectedRoute from "../protected-route";
import Card from "../ui/Card";
import Button from "../ui/Button";
import DiscountAddModal from "./DiscountAddModal";
import Alert from "../ui/Alert";



function Countdown({ expiresAt, onExpire }: { expiresAt: string, onExpire?: () => void }) {
  const [timeLeft, setTimeLeft] = useState<number>(() => Math.max(0, Math.floor((new Date(expiresAt).getTime() - Date.now()) / 1000)));
  useEffect(() => {
    if (timeLeft <= 0) {
      onExpire && onExpire();
      return;
    }
    const interval = setInterval(() => {
      setTimeLeft(Math.max(0, Math.floor((new Date(expiresAt).getTime() - Date.now()) / 1000)));
    }, 1000);
    return () => clearInterval(interval);
  }, [expiresAt, timeLeft, onExpire]);
  if (timeLeft <= 0) return <span className="text-red-500">Expired</span>;
  const h = Math.floor(timeLeft / 3600);
  const m = Math.floor((timeLeft % 3600) / 60);
  const s = timeLeft % 60;
  return <span>{h > 0 && `${h}h `}{m > 0 && `${m}m `}{s}s left</span>;
}

export default function DiscountsPage() {
  const [discounts, setDiscounts] = useState<any[]>([]);
  const [addOpen, setAddOpen] = useState(false);
  const [error, setError] = useState("");

  const fetchDiscounts = () => {
    fetch("http://localhost:8080/api/discounts")
      .then((res) => res.json())
      .then((data) => setDiscounts(Array.isArray(data) ? data : []));
  };

  useEffect(() => { fetchDiscounts(); }, []);

  const handleDelete = async (id: string) => {
    if (!window.confirm("Delete this discount?")) return;
    setError("");
    const res = await fetch(`http://localhost:8080/api/discounts/${id}`, { method: "DELETE" });
    if (res.ok) fetchDiscounts();
    else setError("Failed to delete");
  };

  const handleRenew = async (id: string) => {
    setError("");
    const res = await fetch(`http://localhost:8080/api/discounts/${id}/renew`, { method: "PATCH" });
    if (res.ok) fetchDiscounts();
    else setError("Failed to renew");
  };

  const now = Date.now();
  // Only use discounts with a valid _id
  // Accept any discount with a truthy _id
  const validDiscounts = discounts.filter(d => d._id);
  const getId = (d: any) => (typeof d._id === "string" ? d._id : d._id?.toString?.() || String(d._id));
  // Show all discounts as active if no expiresAt and no active field (for legacy data)
  const active = validDiscounts.filter(d => (d.active === undefined && !d.expiresAt) || (d.active && (!d.expiresAt || new Date(d.expiresAt).getTime() > now)));
  const expired = validDiscounts.filter(d => (d.active === false) || (d.expiresAt && new Date(d.expiresAt).getTime() <= now));

  return (
    <ProtectedRoute allowedRoles={["admin"]}>
      <div className="min-h-screen p-8">
        <h1 className="text-3xl font-bold mb-6">Discounts</h1>
        {error && <Alert type="error">{error}</Alert>}
        <div className="mb-6 flex gap-2">
          <Button onClick={() => setAddOpen(true)}>Add Discount</Button>
        </div>
        <DiscountAddModal open={addOpen} onClose={() => setAddOpen(false)} onAdded={fetchDiscounts} />
        <h2 className="text-xl font-semibold mt-8 mb-2">Active Discounts</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {active.length === 0 && <div className="text-gray-500">No active discounts.</div>}
          {active.map((d) => (
            <Card key={getId(d)}>
              <div className="flex flex-col gap-1">
                <div className="font-bold text-lg">{d.name}</div>
                <div className="text-blue-700 dark:text-blue-300 font-semibold">{d.percent}% off</div>
                <div className="text-sm text-gray-500">
                  {d.type === "static"
                    ? "Static discount"
                    : (d.type === "code" && d.code)
                      ? `Code: ${d.code}`
                      : d.type === "code"
                        ? "Code-based discount (no code)"
                        : "Unknown type"}
                </div>
                {d.expiresAt && (
                  <div className="text-xs text-orange-600 mt-1">
                    <Countdown expiresAt={d.expiresAt} onExpire={fetchDiscounts} />
                  </div>
                )}
                <div className="flex gap-2 mt-2">
                  <Button onClick={() => handleDelete(getId(d))} className="bg-red-600 hover:bg-red-700">Delete</Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
        <h2 className="text-xl font-semibold mt-10 mb-2">Expired Discounts</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {expired.length === 0 && <div className="text-gray-500">No expired discounts.</div>}
          {expired.map((d) => (
            <Card key={getId(d)}>
              <div className="flex flex-col gap-1">
                <div className="font-bold text-lg">{d.name}</div>
                <div className="text-blue-700 dark:text-blue-300 font-semibold">{d.percent}% off</div>
                <div className="text-sm text-gray-500">
                  {d.type === "static"
                    ? "Static discount"
                    : (d.type === "code" && d.code)
                      ? `Code: ${d.code}`
                      : d.type === "code"
                        ? "Code-based discount (no code)"
                        : "Unknown type"}
                </div>
                {d.expiresAt && (
                  <div className="text-xs text-red-600 mt-1">Expired</div>
                )}
                <div className="flex gap-2 mt-2">
                  <Button onClick={() => handleRenew(getId(d))} className="bg-green-600 hover:bg-green-700">Renew</Button>
                  <Button onClick={() => handleDelete(getId(d))} className="bg-red-600 hover:bg-red-700">Delete</Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      </div>
    </ProtectedRoute>
  );
}