"use client";
import { useState } from "react";
import Modal from "../ui/Modal";

export default function CustomerAddModal({ open, onClose, onAdded }: {
  open: boolean;
  onClose: () => void;
  onAdded: () => void;
}) {
  const [form, setForm] = useState({ name: "", email: "", phone: "", notes: "", tags: "" });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) {
    setForm({ ...form, [e.target.name]: e.target.value });
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    fetch("http://localhost:8080/api/customers", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: form.name,
        email: form.email,
        phone: form.phone,
        notes: form.notes,
        tags: form.tags.split(",").map(t => t.trim()).filter(Boolean),
      }),
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to add customer");
        return res.json();
      })
      .then(() => {
        setForm({ name: "", email: "", phone: "", notes: "", tags: "" });
        setLoading(false);
        onAdded();
        onClose();
      })
      .catch(() => {
        setError("Failed to add customer");
        setLoading(false);
      });
  }

  return (
    <Modal open={open} onClose={onClose} title="Add Customer">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block mb-1 font-medium">Name</label>
          <input name="name" required className="border rounded px-2 py-1 w-full" value={form.name} onChange={handleChange} placeholder="Full name" title="Customer name" />
        </div>
        <div>
          <label className="block mb-1 font-medium">Email</label>
          <input name="email" type="email" className="border rounded px-2 py-1 w-full" value={form.email} onChange={handleChange} placeholder="Email address" title="Email address" />
        </div>
        <div>
          <label className="block mb-1 font-medium">Phone</label>
          <input name="phone" className="border rounded px-2 py-1 w-full" value={form.phone} onChange={handleChange} placeholder="Phone number" title="Phone number" />
        </div>
        <div>
          <label className="block mb-1 font-medium">Tags <span className="text-gray-400">(comma separated)</span></label>
          <input name="tags" className="border rounded px-2 py-1 w-full" value={form.tags} onChange={handleChange} placeholder="e.g. VIP, allergy" title="Tags" />
        </div>
        <div>
          <label className="block mb-1 font-medium">Notes</label>
          <textarea name="notes" className="border rounded px-2 py-1 w-full" value={form.notes} onChange={handleChange} placeholder="Notes about this customer" title="Notes" />
        </div>
        {error && <div className="text-red-500">{error}</div>}
        <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700" disabled={loading}>
          {loading ? "Adding..." : "Add Customer"}
        </button>
      </form>
    </Modal>
  );
}
