"use client";
import { useState, useEffect, useRef } from "react";

interface BookingAddModalProps {
  open: boolean;
  onClose: () => void;
  onAdded: () => void;
}

export default function BookingAddModal({ open, onClose, onAdded }: BookingAddModalProps) {
  const [tableNumber, setTableNumber] = useState("");
  const [customerId, setCustomerId] = useState("");
  const [notes, setNotes] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [customerOptions, setCustomerOptions] = useState<{ id: string; name: string }[]>([]);
  const [customerSearch, setCustomerSearch] = useState("");
  const [bookingTime, setBookingTime] = useState("");
  const [showDropdown, setShowDropdown] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!open) return;
    if (customerSearch.length === 0) {
      setCustomerOptions([]);
      return;
    }
    fetch(`http://localhost:8080/api/customers?q=${encodeURIComponent(customerSearch)}`)
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) {
          setCustomerOptions(data.map((c: any) => ({ id: c.id, name: c.name })));
        }
      });
  }, [customerSearch, open]);

  if (!open) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    const res = await fetch("http://localhost:8080/api/bookings", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ tableNumber, customerId, notes, products: [], billTotal: 0, bookingTime })
    });
    if (res.ok) {
      setLoading(false);
      onAdded();
    } else {
      setError("Failed to add booking");
      setLoading(false);
    }
  };

  const handleCustomerSelect = (id: string, name: string) => {
    setCustomerId(id);
    setCustomerSearch(name);
    setShowDropdown(false);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-gray-900 rounded shadow-lg p-6 w-full max-w-md">
        <h2 className="text-xl font-bold mb-4">Add Booking</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-3">
            <label className="block mb-1 font-medium">Table Number</label>
            <input
              className="border rounded px-2 py-1 w-full"
              value={tableNumber}
              onChange={e => setTableNumber(e.target.value)}
              required
              placeholder="Enter table number"
              title="Table Number"
            />
          </div>
          <div className="mb-3">
            <label className="block mb-1 font-medium">Booking Date & Time</label>
            <input
              type="datetime-local"
              className="border rounded px-2 py-1 w-full"
              value={bookingTime}
              onChange={e => setBookingTime(e.target.value)}
              required
            />
          </div>
          <div className="mb-3 relative">
            <label className="block mb-1 font-medium">Customer</label>
            <input
              ref={inputRef}
              className="border rounded px-2 py-1 w-full mb-1"
              placeholder="Search customers..."
              title="Customer"
              value={customerSearch}
              onChange={e => {
                setCustomerSearch(e.target.value);
                setShowDropdown(true);
                setCustomerId("");
              }}
              autoComplete="off"
              onFocus={() => setShowDropdown(true)}
              onBlur={() => setTimeout(() => setShowDropdown(false), 100)}
            />
            {showDropdown && customerOptions.length > 0 && (
              <ul className="absolute z-10 bg-white dark:bg-gray-800 border rounded w-full max-h-40 overflow-y-auto shadow">
                {customerOptions.map((c) => (
                  <li
                    key={c.id}
                    className="px-3 py-2 cursor-pointer hover:bg-blue-100 dark:hover:bg-blue-900"
                    onMouseDown={() => handleCustomerSelect(c.id, c.name)}
                  >
                    {c.name}
                  </li>
                ))}
              </ul>
            )}
            {showDropdown && customerSearch && customerOptions.length === 0 && (
              <div className="absolute z-10 bg-white dark:bg-gray-800 border rounded w-full px-3 py-2 text-gray-500">No results</div>
            )}
          </div>
          <div className="mb-3">
            <label className="block mb-1 font-medium">Notes</label>
            <textarea
              className="border rounded px-2 py-1 w-full"
              value={notes}
              onChange={e => setNotes(e.target.value)}
              placeholder="Enter notes"
              title="Notes"
            />
          </div>
          {error && <div className="text-red-500 mb-2">{error}</div>}
          <div className="flex justify-end gap-2 mt-4">
            <button type="button" className="px-4 py-2 rounded bg-gray-200 dark:bg-gray-700" onClick={onClose} disabled={loading}>Cancel</button>
            <button type="submit" className="px-4 py-2 rounded bg-blue-600 text-white" disabled={loading}>{loading ? "Adding..." : "Add Booking"}</button>
          </div>
        </form>
      </div>
    </div>
  );
}
