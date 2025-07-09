"use client";
import React, { useState } from "react";
import Modal from "../ui/Modal";
import Input from "../ui/Input";
import Button from "../ui/Button";
import Alert from "../ui/Alert";

export default function DiscountAddModal({ open, onClose, onAdded }: {
  open: boolean;
  onClose: () => void;
  onAdded: () => void;
}) {
  const [name, setName] = useState("");
  const [percent, setPercent] = useState("");
  const [type, setType] = useState<"static" | "code">("static");
  const [timed, setTimed] = useState(false);
  const [duration, setDuration] = useState(60); // minutes
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    if (!name.trim() || !percent.trim() || isNaN(Number(percent)) || Number(percent) <= 0) {
      setError("Please enter a valid name and percentage.");
      setLoading(false);
      return;
    }
    if (type === "code" && !code.trim()) {
      setError("Please enter a code for code-based discounts.");
      setLoading(false);
      return;
    }
    const payload: any = {
      name: name.trim(),
      percent: Number(percent),
      type,
    };
    if (type === "code") payload.code = code.trim();
    if (timed) {
      const expiresAt = new Date(Date.now() + duration * 60 * 1000);
      payload.expiresAt = expiresAt.toISOString();
    }
    try {
      const res = await fetch("http://localhost:8080/api/discounts", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const data = await res.json();
        setError(data.error || "Failed to add discount");
      } else {
        setName("");
        setPercent("");
        setCode("");
        setType("static");
        onAdded();
        onClose();
      }
    } catch (e) {
      setError("Network error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal open={open} onClose={onClose} title="Add Discount">
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        {error && <Alert type="error">{error}</Alert>}
        <Input
          placeholder="Discount name"
          value={name}
          onChange={e => setName(e.target.value)}
          required
        />
        <Input
          placeholder="Percent (e.g. 10 for 10%)"
          value={percent}
          onChange={e => setPercent(e.target.value)}
          type="number"
          min={1}
          max={100}
          required
        />
        <div className="flex gap-4 items-center">
          <label className="font-medium">Type:</label>
          <label className="flex items-center gap-1">
            <input type="radio" name="type" value="static" checked={type === "static"} onChange={()=>setType("static")}/>
            Static
          </label>
          <label className="flex items-center gap-1">
            <input type="radio" name="type" value="code" checked={type === "code"} onChange={()=>setType("code")}/>
            Code
          </label>
        </div>
        <div className="flex gap-4 items-center">
          <label className="font-medium">Timed?</label>
          <label className="flex items-center gap-1">
            <input
              type="checkbox"
              checked={timed}
              onChange={e => setTimed(e.target.checked)}
              title="Enable timed discount"
            />
            Timed
          </label>
          {timed && (
            <>
              <label className="ml-2">Duration (minutes):</label>
              <Input
                type="number"
                min={1}
                max={1440}
                value={duration}
                onChange={e => setDuration(Number(e.target.value))}
                className="w-24"
              />
            </>
          )}
        </div>
        {type === "code" && (
          <Input
            placeholder="Discount code (e.g. mon10)"
            value={code}
            onChange={e => setCode(e.target.value)}
            required
          />
        )}
        <Button type="submit" disabled={loading}>{loading ? "Adding..." : "Add Discount"}</Button>
      </form>
    </Modal>
  );
}
