"use client";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { getAuth } from "./auth";

export default function ProtectedRoute({ children, allowedRoles }: { children: React.ReactNode; allowedRoles: string[] }) {
  const router = useRouter();
  useEffect(() => {
    const { token, role } = getAuth();
    if (!token || !role) {
      router.replace("/login");
    } else if (!allowedRoles.includes(role)) {
      router.replace("/unauthorized");
    }
  }, [allowedRoles, router]);
  return <>{children}</>;
}
