"use client";
import React, { useState, useEffect } from "react";
import Modal from "../ui/Modal";
import Input from "../ui/Input";
import Button from "../ui/Button";
import Alert from "../ui/Alert";

export default function ProductAddModal({ open, onClose, onProductAdded }: {
  open: boolean;
  onClose: () => void;
  onProductAdded?: () => void;
}) {
  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [category, setCategory] = useState("");
  const [categories, setCategories] = useState<string[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (open) {
      fetch("http://localhost:8080/api/categories")
        .then((res) => res.json())
        .then((data) => setCategories(data.map((c: any) => c.name)));
    }
  }, [open]);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    if (!name || !price || !category) {
      setError("All fields are required");
      return;
    }
    if (isNaN(Number(price)) || Number(price) < 0) {
      setError("Price must be a positive number");
      return;
    }
    setLoading(true);
    fetch("http://localhost:8080/api/products", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, price: Number(price), category }),
    })
      .then((res) => {
        setLoading(false);
        if (!res.ok) throw new Error();
        setName("");
        setPrice("");
        setCategory("");
        onClose();
        if (onProductAdded) onProductAdded();
      })
      .catch(() => {
        setLoading(false);
        setError("Failed to add product");
      });
  }

  return (
    <Modal open={open} onClose={onClose} title="Add Product">
      <form className="flex flex-col gap-3" onSubmit={handleSubmit}>
        <Input placeholder="Product name" value={name} onChange={e => setName(e.target.value)} />
        <Input placeholder="Price" value={price} onChange={e => setPrice(e.target.value)} type="number" min="0" step="0.01" />
        <label htmlFor="category-select" className="sr-only">
          Category
        </label>
        <select
          id="category-select"
          className="border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
          value={category}
          onChange={e => setCategory(e.target.value)}
        >
          <option value="">Select category</option>
          {categories.map((cat) => (
            <option key={cat} value={cat}>{cat}</option>
          ))}
        </select>
        {error && <Alert type="error">{error}</Alert>}
        <Button type="submit" disabled={loading}>{loading ? "Adding..." : "Add Product"}</Button>
      </form>
    </Modal>
  );
}
