"use client";
import React, { useEffect, useState } from "react";
import ProtectedRoute from "../protected-route";
import Card from "../ui/Card";
import Button from "../ui/Button";
import Alert from "../ui/Alert";

type Product = {
  id: string;
  name: string;
  price: number;
};

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:8080/api/products")
      .then((res) => res.json())
      .then((data) => {
        setProducts(data);
        setLoading(false);
      });
  }, []);

  return (
    <ProtectedRoute allowedRoles={["admin", "manager"]}>
      <div className="min-h-screen p-8 bg-gray-50 dark:bg-gray-900">
        <h1 className="text-3xl font-bold mb-6">Products</h1>
        {loading ? (
          <div className="flex justify-center items-center h-32">
            <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-600"></div>
          </div>
        ) : products.length === 0 ? (
          <Alert type="info">No products found.</Alert>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {products.map((product) => (
              <Card key={product.id}>
                <div className="flex flex-col gap-2">
                  <div className="flex items-center justify-between">
                    <span className="font-semibold text-lg text-gray-800 dark:text-gray-100">{product.name}</span>
                    <span className="text-blue-700 dark:text-blue-300 font-bold text-lg">£{product.price.toFixed(2)}</span>
                  </div>
                  {/* Example action button, can be extended */}
                  <Button className="mt-2 w-full">View Details</Button>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </ProtectedRoute>
  );
}
