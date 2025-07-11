"use client";
import { useEffect, useState } from "react";

interface BookingDetailModalProps {
  open: boolean;
  bookingId: string | null;
  onClose: () => void;
}

interface Booking {
  id: string;
  customerId: string;
  tableNumber: string;
  products: { name: string; qty: number; price: number }[];
  billTotal: number;
  status: string;
  createdAt: string;
  bookingTime: string;
  closedAt?: string;
  notes?: string;
}

export default function BookingDetailModal({ open, bookingId, onClose }: BookingDetailModalProps) {
  const [booking, setBooking] = useState<Booking | null>(null);
  const [loading, setLoading] = useState(false);
  const [editStatus, setEditStatus] = useState<string | null>(null);
  const [editTime, setEditTime] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!open || !bookingId) return;
    setLoading(true);
    fetch(`http://localhost:8080/api/bookings/${bookingId}`)
      .then(res => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then(data => {
        setBooking(data);
        setEditStatus(data.status);
        setEditTime(data.bookingTime ? data.bookingTime.slice(0, 16) : "");
        setLoading(false);
      })
      .catch(() => {
        setBooking(null);
        setLoading(false);
      });
  }, [open, bookingId]);

  const handleSave = async () => {
    if (!bookingId) return;
    setSaving(true);
    // Always send bookingTime in the correct format (YYYY-MM-DDTHH:mm)
    const patchBody: any = { status: editStatus };
    if (editTime) patchBody.bookingTime = editTime;
    await fetch(`http://localhost:8080/api/bookings/${bookingId}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(patchBody)
    });
    setSaving(false);
    onClose();
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-gray-900 rounded shadow-lg p-6 w-full max-w-lg relative">
        <button className="absolute top-2 right-2 text-gray-500 hover:text-gray-800" onClick={onClose}>&times;</button>
        {loading ? (
          <div>Loading...</div>
        ) : !booking ? (
          <div className="text-red-500">Booking not found or failed to load.</div>
        ) : (
          <div>
            <h2 className="text-xl font-bold mb-2">Booking for Table {booking.tableNumber}</h2>
            <div className="mb-2 text-gray-700 dark:text-gray-300">Status: {" "}
              <select value={editStatus ?? booking.status} onChange={e => setEditStatus(e.target.value)} className="border rounded px-2 py-1 ml-2">
                <option value="open">Open</option>
                <option value="closed">Closed</option>
                <option value="cancelled">Cancelled</option>
              </select>
            </div>
            <div className="mb-2 text-gray-700 dark:text-gray-300">Booking Time: <span className="font-semibold">{booking.bookingTime ? new Date(booking.bookingTime).toLocaleString() : "N/A"}</span></div>
            <div className="mb-2 text-gray-700 dark:text-gray-300 flex items-center">Edit Booking Time: {" "}
              <input type="datetime-local" value={editTime ?? booking.bookingTime?.slice(0,16)} onChange={e => setEditTime(e.target.value)} className="border rounded px-2 py-1 ml-2" />
            </div>
            <div className="mb-2 text-gray-700 dark:text-gray-300">Created: {new Date(booking.createdAt).toLocaleString()}</div>
            {booking.closedAt && <div className="mb-2 text-gray-700 dark:text-gray-300">Closed: {new Date(booking.closedAt).toLocaleString()}</div>}
            <div className="mb-2 text-gray-700 dark:text-gray-300">Customer ID: {booking.customerId}</div>
            <div className="mb-2 text-gray-700 dark:text-gray-300">Notes: {booking.notes || <span className="italic">None</span>}</div>
            <div className="mb-4">
              <h3 className="font-semibold mb-1">Products</h3>
              {booking.products.length === 0 ? (
                <div className="text-gray-500 italic">No products added.</div>
              ) : (
                <ul className="divide-y divide-gray-200 dark:divide-gray-700">
                  {booking.products.map((p, i) => (
                    <li key={i} className="py-1 flex justify-between">
                      <span>{p.name} x{p.qty}</span>
                      <span>£{(p.price * p.qty).toFixed(2)}</span>
                    </li>
                  ))}
                </ul>
              )}
            </div>
            <div className="text-lg font-bold text-right mb-4">Bill Total: £{booking.billTotal.toFixed(2)}</div>
            <div className="flex justify-end gap-2">
              <button className="px-4 py-2 rounded bg-gray-200 dark:bg-gray-700" onClick={onClose} disabled={saving}>Cancel</button>
              <button className="px-4 py-2 rounded bg-blue-600 text-white" onClick={handleSave} disabled={saving}>{saving ? "Saving..." : "Save Changes"}</button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
