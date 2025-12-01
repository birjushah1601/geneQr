"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { attachmentsApi } from "@/lib/api/attachments";

export default function UnassignedAttachmentsPage() {
  const qc = useQueryClient();
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["attachments", "unassigned"],
    queryFn: () => attachmentsApi.list({ page: 1, page_size: 50, unassigned: true, source: "whatsapp" }),
  });

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4">
          <h1 className="text-lg font-semibold">Unassigned Attachments</h1>
          <p className="text-sm text-gray-500">Media received via WhatsApp that is not yet linked to a ticket</p>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6">
        {isLoading ? (
          <p className="text-sm text-gray-500">Loading…</p>
        ) : error ? (
          <p className="text-sm text-red-600">Failed to load: {(error as any)?.message || "Unknown error"}</p>
        ) : (data?.items?.length ? (
          <div className="bg-white border rounded divide-y">
            {data.items.map((att) => (
              <AttachmentRow key={att.id} id={att.id} fileName={att.fileName} fileType={att.fileType} fileSize={att.fileSize} uploadDate={att.uploadDate} onLinked={() => refetch()} />
            ))}
          </div>
        ) : (
          <p className="text-sm text-gray-500">No unassigned attachments found.</p>
        ))}
      </div>
    </div>
  );
}

function AttachmentRow({ id, fileName, fileType, fileSize, uploadDate, onLinked }: { id: string; fileName: string; fileType: string; fileSize: number; uploadDate: string; onLinked: () => void }) {
  const [ticketId, setTicketId] = useState("");
  const link = useMutation({
    mutationFn: async () => attachmentsApi.link(id, ticketId),
    onSuccess: () => { setTicketId(""); onLinked(); },
  });

  return (
    <div className="p-3 flex items-center justify-between gap-4">
      <div className="min-w-0">
        <div className="font-medium truncate">{fileName}</div>
        <div className="text-xs text-gray-500">{fileType} • {(fileSize/1024).toFixed(1)} KB • {new Date(uploadDate).toLocaleString()}</div>
        <div className="text-xs text-gray-400">ID: {id}</div>
      </div>
      <div className="flex items-center gap-2">
        <input value={ticketId} onChange={(e) => setTicketId(e.target.value)} placeholder="Ticket ID" className="border rounded px-2 py-1 text-sm" />
        <button onClick={() => link.mutate()} disabled={!ticketId || link.isLoading} className="px-3 py-1.5 bg-indigo-600 text-white rounded text-sm disabled:opacity-50">
          {link.isLoading ? "Linking…" : "Link"}
        </button>
      </div>
    </div>
  );
}
