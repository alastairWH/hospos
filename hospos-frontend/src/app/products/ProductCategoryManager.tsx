"use client";
import React, { useEffect, useState } from "react";
import Card from "../ui/Card";
import Button from "../ui/Button";
import Input from "../ui/Input";
import Alert from "../ui/Alert";

export default function ProductCategoryManager({ onCategoryAdded }: { onCategoryAdded?: () => void }) {
  const [categories, setCategories] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showAdd, setShowAdd] = useState(false);
  const [newCategory, setNewCategory] = useState("");

  function fetchCategories() {
    setLoading(true);
    fetch("http://localhost:8080/api/categories")
      .then((res) => res.json())
      .then((data) => { setCategories(data.map((c: any) => c.name)); setLoading(false); })
      .catch(() => { setError("Failed to load categories"); setLoading(false); });
  }

  useEffect(() => { fetchCategories(); }, []);

  function handleAddCategory(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    if (!newCategory) {
      setError("Category name required");
      return;
    }
    fetch("http://localhost:8080/api/categories", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: newCategory }),
    })
      .then((res) => {
        if (!res.ok) throw new Error();
        setShowAdd(false);
        setNewCategory("");
        fetchCategories();
        if (onCategoryAdded) onCategoryAdded();
      })
      .catch(() => setError("Failed to add category"));
  }

  return (
    <Card className="mb-6">
      <div className="flex justify-between items-center mb-2">
        <h2 className="text-xl font-semibold">Product Categories</h2>
        <Button onClick={() => setShowAdd(true)}>Add Category</Button>
      </div>
      {error && <Alert type="error">{error}</Alert>}
      {loading ? (
        <div className="flex justify-center items-center h-16">
          <div className="animate-spin rounded-full h-6 w-6 border-t-2 border-b-2 border-blue-600"></div>
        </div>
      ) : (
        <ul className="flex flex-wrap gap-2">
          {categories.map((cat) => (
            <li key={cat} className="bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded-full px-4 py-1 text-sm font-medium">
              {cat}
            </li>
          ))}
        </ul>
      )}
      {showAdd && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <form className="bg-white dark:bg-gray-800 p-6 rounded-2xl shadow flex flex-col gap-3 min-w-[300px] border border-gray-100 dark:border-gray-700" onSubmit={handleAddCategory}>
            <h3 className="text-lg font-bold mb-2">Add Category</h3>
            <Input placeholder="Category name" value={newCategory} onChange={e => setNewCategory(e.target.value)} />
            <div className="flex gap-2 mt-2">
              <Button type="submit">Add</Button>
              <Button type="button" className="bg-gray-300 text-gray-800 dark:bg-gray-700 dark:text-gray-200" onClick={() => setShowAdd(false)}>Cancel</Button>
            </div>
          </form>
        </div>
      )}
    </Card>
  );
}
