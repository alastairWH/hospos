import React, { useEffect, useState } from "react";
import ProtectedRoute from "../protected-route";

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
          <p>Loading...</p>
        ) : (
          <table className="min-w-full bg-white dark:bg-gray-800 rounded shadow">
            <thead>
              <tr>
                <th className="py-2 px-4 text-left">Name</th>
                <th className="py-2 px-4 text-left">Price</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <tr key={product.id}>
                  <td className="py-2 px-4">{product.name}</td>
                  <td className="py-2 px-4">£{product.price.toFixed(2)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </ProtectedRoute>
  );
}
