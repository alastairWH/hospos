"use client";
import React, { useEffect, useState } from "react";

interface BusinessInfo {
  companyName: string;
  companyAddress: string;
  financeEmail: string;
  vatId: string;
  companyRegNumber: string;
  phone: string;
  website: string;
  logoUrl: string;
  salesIdPrefix: string;
  currency: string;
  defaultTaxRate: string;
  bankDetails: string;
  legalFooter: string;
  openingHours: string;
  socialLinks: string;
  customReceiptMsg: string;
  invoiceFormat: string;
  country: string;
  lastSalesNumber: string;
}

const defaultInfo: BusinessInfo = {
  companyName: "",
  companyAddress: "",
  financeEmail: "",
  vatId: "",
  companyRegNumber: "",
  phone: "",
  website: "",
  logoUrl: "",
  salesIdPrefix: "",
  currency: "GBP",
  defaultTaxRate: "20",
  bankDetails: "",
  legalFooter: "",
  openingHours: "",
  socialLinks: "",
  customReceiptMsg: "",
  invoiceFormat: "",
  country: "UK",
  lastSalesNumber: "0",
};

export default function AdminBusinessInfo() {
  const [info, setInfo] = useState<BusinessInfo>(defaultInfo);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/api/business")
      .then((res) => res.json())
      .then((data) => {
        setInfo({ ...defaultInfo, ...data });
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
      });
  }, []);

  function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) {
    setInfo({ ...info, [e.target.name]: e.target.value });
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError("");
    setSuccess("");
    fetch("http://localhost:8080/api/business", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        ...info,
        defaultTaxRate: parseFloat(info.defaultTaxRate),
        lastSalesNumber: parseInt(info.lastSalesNumber, 10),
        socialLinks: info.socialLinks.split(",").map((s) => s.trim()).filter(Boolean),
      }),
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to save");
        return res.json();
      })
      .then(() => {
        setSuccess("Business info saved!");
        setSaving(false);
      })
      .catch(() => {
        setError("Failed to save business info");
        setSaving(false);
      });
  }

  if (loading) return <div>Loading business info...</div>;

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-xl font-bold mb-2">Business Information</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block font-medium">Company Name</label>
          <input name="companyName" value={info.companyName} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Company Address</label>
          <input name="companyAddress" value={info.companyAddress} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Finance Email</label>
          <input name="financeEmail" value={info.financeEmail} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">VAT ID</label>
          <input name="vatId" value={info.vatId} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Company Reg Number</label>
          <input name="companyRegNumber" value={info.companyRegNumber} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Phone</label>
          <input name="phone" value={info.phone} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Website</label>
          <input name="website" value={info.website} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Logo URL</label>
          <input name="logoUrl" value={info.logoUrl} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Sales ID Prefix</label>
          <input name="salesIdPrefix" value={info.salesIdPrefix} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Currency</label>
          <input name="currency" value={info.currency} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Default Tax Rate (%)</label>
          <input name="defaultTaxRate" value={info.defaultTaxRate} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Bank Details</label>
          <input name="bankDetails" value={info.bankDetails} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Legal Footer</label>
          <input name="legalFooter" value={info.legalFooter} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Opening Hours</label>
          <input name="openingHours" value={info.openingHours} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Social Links (comma separated)</label>
          <input name="socialLinks" value={info.socialLinks} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Custom Receipt Message</label>
          <input name="customReceiptMsg" value={info.customReceiptMsg} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Invoice Format</label>
          <input name="invoiceFormat" value={info.invoiceFormat} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Country</label>
          <input name="country" value={info.country} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
        <div>
          <label className="block font-medium">Last Sales Number</label>
          <input name="lastSalesNumber" value={info.lastSalesNumber} onChange={handleChange} className="border rounded px-2 py-1 w-full" />
        </div>
      </div>
      {error && <div className="text-red-500">{error}</div>}
      {success && <div className="text-green-600">{success}</div>}
      <button type="submit" className="px-4 py-2 rounded bg-blue-600 text-white" disabled={saving}>{saving ? "Saving..." : "Save Business Info"}</button>
    </form>
  );
}
