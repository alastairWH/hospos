"use client";
import { useEffect, useState } from "react";
import BookingAddModal from "./BookingAddModal";
import BookingDetailModal from "./BookingDetailModal";

interface Booking {
  id: string;
  customerId: string;
  tableNumber: string;
  products: any[];
  billTotal: number;
  status: string;
  createdAt: string;
  bookingTime: string;
  closedAt?: string;
  notes?: string;
}

function BookingCalendar({ bookings }: { bookings: any[] }) {
  // Group by bookingTime (not createdAt)
  const grouped = bookings.reduce((acc: Record<string, any[]>, b) => {
    const date = b.bookingTime ? new Date(b.bookingTime).toLocaleDateString() : "Unknown";
    acc[date] = acc[date] || [];
    acc[date].push(b);
    return acc;
  }, {});
  return (
    <div className="overflow-x-auto whitespace-nowrap py-2 mb-4">
      {Object.entries(grouped).map(([date, items]) => (
        <div key={date} className="inline-block align-top mr-4">
          <div className="font-bold text-sm mb-1">{date}</div>
          <div className="space-y-2">
            {items.map((b) => (
              <div key={b.id} className="bg-blue-100 dark:bg-blue-900 rounded px-3 py-2 text-xs mb-1">
                Table {b.tableNumber} - £{b.billTotal.toFixed(2)} - {b.status}
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

function groupBookings(bookings: Booking[]) {
  const now = new Date();
  const todayStr = now.toISOString().slice(0, 10);
  const previous: Booking[] = [];
  const today: Booking[] = [];
  const future: Booking[] = [];
  bookings.forEach((b) => {
    const dateStr = b.bookingTime ? new Date(b.bookingTime).toISOString().slice(0, 10) : "";
    if (dateStr < todayStr) previous.push(b);
    else if (dateStr === todayStr) today.push(b);
    else future.push(b);
  });
  return { previous, today, future };
}

export default function BookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showAdd, setShowAdd] = useState(false);
  const [selectedBookingId, setSelectedBookingId] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);
    fetch("http://localhost:8080/api/bookings")
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        setBookings(Array.isArray(data) ? data : []);
        setLoading(false);
      })
      .catch((err) => {
        setError("Failed to load bookings. Please try again later.");
        setLoading(false);
      });
  }, []);

  const { previous, today, future } = groupBookings(bookings);

  return (
    <div className="max-w-5xl mx-auto p-6">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Bookings</h1>
        <button className="bg-blue-600 text-white px-4 py-2 rounded" onClick={() => setShowAdd(true)}>
          Add Booking
        </button>
      </div>
      <BookingAddModal open={showAdd} onClose={() => setShowAdd(false)} onAdded={() => {
        setShowAdd(false);
        setLoading(true);
        setError(null);
        fetch("http://localhost:8080/api/bookings")
          .then(res => {
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            return res.json();
          })
          .then(data => {
            setBookings(Array.isArray(data) ? data : []);
            setLoading(false);
          })
          .catch(() => {
            setError("Failed to load bookings. Please try again later.");
            setLoading(false);
          });
      }} />
      <div className="mb-8">
        <BookingCalendar bookings={bookings} />
      </div>
      {loading ? (
        <div>Loading...</div>
      ) : error ? (
        <div className="text-red-500 mb-4">{error}</div>
      ) : bookings.length === 0 ? (
        <div>No bookings found.</div>
      ) : (
        <>
          <Section title="Today's Bookings" bookings={today} onSelect={setSelectedBookingId} />
          <Section title="Future Bookings" bookings={future} onSelect={setSelectedBookingId} />
          <Section title="Previous Bookings" bookings={previous} onSelect={setSelectedBookingId} />
        </>
      )}
      <BookingDetailModal open={!!selectedBookingId} bookingId={selectedBookingId} onClose={() => setSelectedBookingId(null)} />
    </div>
  );
}

function Section({ title, bookings, onSelect }: { title: string; bookings: Booking[]; onSelect: (id: string) => void }) {
  if (bookings.length === 0) return null;
  return (
    <div className="mb-8">
      <h2 className="text-lg font-bold mb-2">{title}</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {bookings.map((b) => (
          <button key={b.id} className="block w-full text-left" onClick={() => onSelect(b.id)}>
            <div className="bg-white dark:bg-gray-900 rounded shadow p-4 hover:shadow-lg transition">
              <div className="flex justify-between items-center mb-2">
                <span className="font-semibold text-lg">Table {b.tableNumber}</span>
                <span className={`px-2 py-1 rounded text-xs ${b.status === "open" ? "bg-green-100 text-green-700" : b.status === "closed" ? "bg-gray-200 text-gray-700" : "bg-red-100 text-red-700"}`}>{b.status.charAt(0).toUpperCase() + b.status.slice(1)}</span>
              </div>
              <div className="text-gray-700 dark:text-gray-300 mb-1">Bill: £{b.billTotal.toFixed(2)}</div>
              <div className="text-gray-500 text-xs mb-1">Booking Time: {b.bookingTime ? new Date(b.bookingTime).toLocaleString() : "N/A"}</div>
              <div className="text-gray-500 text-xs mb-1">Created: {new Date(b.createdAt).toLocaleString()}</div>
              {b.closedAt && <div className="text-gray-500 text-xs">Closed: {new Date(b.closedAt).toLocaleString()}</div>}
              {b.notes && <div className="text-gray-400 text-xs mt-2 italic">{b.notes}</div>}
            </div>
          </button>
        ))}
      </div>
    </div>
  );
}
