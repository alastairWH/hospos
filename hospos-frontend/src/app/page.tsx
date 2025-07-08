"use client";
import Link from "next/link";

export default function LandingPage() {
  return (
    <main className="min-h-screen bg-base-200 flex items-center justify-center">
      <div className="text-center space-y-6">
        {/* App Name */}
        <h1 className="text-6xl font-bold text-base-content">hospos</h1>
        {/* Slogan */}
        <p className="text-2xl text-base-content opacity-75">
          hospos - Open Source Hospitality POS system
        </p>
        {/* Buttons */}
        <div className="flex justify-center space-x-4">
          <Link href="/admin" className="btn btn-primary">Go to Login</Link>
        </div>
      </div>
    </main>
  );
}

