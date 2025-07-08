"use client";

import React, { useEffect, useState } from "react";
import ProtectedRoute from "../protected-route";
import Card from "../ui/Card";

export default function DiscountsPage(){
    const [discounts, setDiscounts] = useState<string[]>([])

    const fetchDiscounts = () => {
        fetch("http://localhost:8080/api/discounts")
        .then((res) => res.json())
        .then((data)=> setDiscounts(data.map((r: any)=> r.discounts)));
    
    };

    useEffect(()=> { fetchDiscounts();}, []);

    return (
        <ProtectedRoute allowedRoles={["admin"]}>
          <div className="min-h-screen p-8 bg-gray-50 dark:bg-gray-900">
                  <h1 className="text-3xl font-bold mb-6">Discounts</h1>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    
                  </div>
                </div>
        </ProtectedRoute>
    )
}