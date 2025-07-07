"use client";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { getAuth } from "./auth";

export default function ProtectedRoute({ children, allowedRoles }: { children: React.ReactNode; allowedRoles: string[] }) {
  const router = useRouter();
  useEffect(() => {
    const { token, role } = getAuth();
    if (!token || !role || !allowedRoles.includes(role)) {
      router.replace("/login");
    }
  }, [allowedRoles, router]);
  return <>{children}</>;
}
