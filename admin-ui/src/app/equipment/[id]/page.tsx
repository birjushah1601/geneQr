"use client";

import { useEffect, useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import { equipmentApi } from "@/lib/api/equipment";
import { ticketsApi } from "@/lib/api/tickets";
import type { Equipment, ServiceTicket } from "@/types";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ArrowLeft, QrCode, Download, RefreshCw } from "lucide-react";

type Tab = "overview" | "specs" | "history" | "qr";

export default function EquipmentDetailPage({ params }: { params: { id: string } }) {
  const router = useRouter();
  const { id } = params;

  const [equipment, setEquipment] = useState<Equipment | null>(null);
  const [tickets, setTickets] = useState<ServiceTicket[]>([]);
  const [tab, setTab] = useState<Tab>("overview");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [regenLoading, setRegenLoading] = useState(false);

  useEffect(() => {
    let active = true;
    const run = async () => {
      try {
        setLoading(true);
        setError(null);
        const eq = await equipmentApi.getById(id);
        if (!active) return;
        setEquipment(eq);
        try {
          const t = await ticketsApi.getByEquipment(id);
          if (Array.isArray((t as any).tickets)) setTickets((t as any).tickets);
        } catch {}
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load equipment");
      } finally {
        if (active) setLoading(false);
      }
    };
    run();
    return () => {
      active = false;
    };
  }, [id]);

  const qrImageUrl = useMemo(() => {
    if (!equipment) return "";
    // Prefer inline base64 image when available; fallback to server image endpoint
    if ((equipment as any).qr_code_image) {
      return `data:image/png;base64,${(equipment as any).qr_code_image}`;
    }
    return `${process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8081/api"}/v1/equipment/qr/image/${equipment.id}`;
  }, [equipment]);

  const handleRegenerateQR = async () => {
    if (!equipment) return;
    try {
      setRegenLoading(true);
      await equipmentApi.generateQRCode(equipment.id);
      // Re-fetch equipment to get fresh base64 image
      const refreshed = await equipmentApi.getById(equipment.id);
      setEquipment(refreshed);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to regenerate QR");
    } finally {
      setRegenLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <p className="text-gray-600">Loading…</p>
      </div>
    );
  }

  if (error || !equipment) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Card className="max-w-md w-full">
          <CardHeader>
            <CardTitle className="text-red-600">Unable to load</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-gray-700 mb-4">{error || "Equipment not found"}</p>
            <Button className="w-full" onClick={() => router.push("/equipment")}>Back to list</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="mb-6 flex items-center justify-between">
          <Button variant="ghost" onClick={() => router.push("/equipment")}> 
            <ArrowLeft className="mr-2 h-4 w-4" /> Back
          </Button>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => router.push(`/equipment/${equipment.id}/edit`)}>Edit</Button>
          </div>
        </div>

        <h1 className="text-2xl font-bold text-gray-900 mb-1">{equipment.equipment_name}</h1>
        <p className="text-gray-600 mb-6">Serial: {equipment.serial_number} • {equipment.manufacturer_name}</p>

        {/* Tabs */}
        <div className="flex gap-2 mb-4">
          {(["overview","specs","history","qr"] as Tab[]).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={`px-4 py-2 rounded-md border text-sm font-medium ${tab===t?"bg-blue-600 text-white border-blue-600":"bg-white text-gray-700 border-gray-300 hover:bg-gray-50"}`}
            >
              {t.toUpperCase()}
            </button>
          ))}
        </div>

        {tab === "overview" && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle>Details</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-gray-700 space-y-2">
                <div><span className="text-gray-500">Model:</span> {equipment.model_number || "N/A"}</div>
                <div><span className="text-gray-500">Category:</span> {equipment.category || "N/A"}</div>
                <div><span className="text-gray-500">Customer:</span> {equipment.customer_name}</div>
                <div><span className="text-gray-500">Location:</span> {equipment.installation_location || "N/A"}</div>
                <div><span className="text-gray-500">Status:</span> {equipment.status}</div>
                <div><span className="text-gray-500">Installed:</span> {equipment.installation_date || "N/A"}</div>
                <div><span className="text-gray-500">Warranty:</span> {equipment.warranty_expiry || "N/A"}</div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Recent Service</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-gray-700 space-y-2">
                <div><span className="text-gray-500">Last Service:</span> {equipment.last_service_date || "N/A"}</div>
                <div><span className="text-gray-500">Count:</span> {equipment.service_count}</div>
                <div className="pt-2">
                  <Button size="sm" variant="outline" onClick={() => router.push(`/tickets?equipment_id=${equipment.id}`)}>View tickets</Button>
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {tab === "specs" && (
          <Card>
            <CardHeader>
              <CardTitle>Specifications</CardTitle>
            </CardHeader>
            <CardContent>
              {equipment.specifications ? (
                <pre className="text-xs bg-gray-50 p-4 rounded-md overflow-x-auto">{JSON.stringify(equipment.specifications, null, 2)}</pre>
              ) : (
                <p className="text-sm text-gray-600">No specifications available.</p>
              )}
            </CardContent>
          </Card>
        )}

        {tab === "history" && (
          <Card>
            <CardHeader>
              <CardTitle>Service History</CardTitle>
            </CardHeader>
            <CardContent>
              {tickets.length === 0 ? (
                <p className="text-sm text-gray-600">No tickets found for this equipment.</p>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="text-left text-gray-500">
                        <th className="py-2 pr-4">Ticket</th>
                        <th className="py-2 pr-4">Status</th>
                        <th className="py-2 pr-4">Priority</th>
                        <th className="py-2 pr-4">Created</th>
                      </tr>
                    </thead>
                    <tbody>
                      {tickets.map((t) => (
                        <tr key={t.id} className="border-t">
                          <td className="py-2 pr-4">{t.ticket_number}</td>
                          <td className="py-2 pr-4">{t.status}</td>
                          <td className="py-2 pr-4">{t.priority}</td>
                          <td className="py-2 pr-4">{new Date(t.created_at).toLocaleString()}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </CardContent>
          </Card>
        )}

        {tab === "qr" && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle>QR Code</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-col items-center gap-4">
                  <div className="border rounded-md p-4 bg-white">
                    <img id="qr-img" src={qrImageUrl} alt="QR code" className="w-56 h-56 object-contain" onError={(e)=>{(e.currentTarget as HTMLImageElement).style.display='none'}} />
                  </div>
                  <div className="flex gap-2">
                    <Button variant="outline" onClick={() => equipmentApi.downloadQRLabel(equipment.id)}>
                      <Download className="mr-2 h-4 w-4" /> Download PDF Label
                    </Button>
                    <Button variant="outline" onClick={handleRegenerateQR} disabled={regenLoading}>
                      <RefreshCw className="mr-2 h-4 w-4" /> {regenLoading ? "Regenerating…" : "Regenerate"}
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Scan Payload</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-sm text-gray-700">
                  <div><span className="text-gray-500">QR:</span> {equipment.qr_code || "N/A"}</div>
                  <div className="mt-2">
                    <Button size="sm" variant="outline" onClick={() => router.push(`/test-qr?qr=${encodeURIComponent(equipment.qr_code || equipment.id)}`)}>
                      <QrCode className="mr-2 h-4 w-4" /> Test scan workflow
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
}
