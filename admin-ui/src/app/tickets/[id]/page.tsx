"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ticketsApi } from "@/lib/api/tickets";
import { apiClient } from "@/lib/api/client";
import type { ServiceTicket, TicketPriority, TicketStatus } from "@/types";
import { ArrowLeft, Loader2, Package, User, Calendar, Wrench, Pause, Play, CheckCircle, XCircle, AlertTriangle, FileText, MessageSquare, Paperclip, Upload } from "lucide-react";
import { attachmentsApi } from "@/lib/api/attachments";

function StatusBadge({ status }: { status: TicketStatus }) {
  const color = {
    new: "bg-gray-100 text-gray-700",
    assigned: "bg-indigo-100 text-indigo-700",
    in_progress: "bg-blue-100 text-blue-700",
    on_hold: "bg-yellow-100 text-yellow-800",
    resolved: "bg-green-100 text-green-700",
    closed: "bg-gray-200 text-gray-800",
    cancelled: "bg-red-100 text-red-700",
  }[status];
  return <span className={`px-2 py-0.5 rounded text-xs font-medium ${color}`}>{status.replaceAll("_", " ")}</span>;
}

export default function TicketDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const qc = useQueryClient();

  const { data: ticket, isLoading, refetch } = useQuery<ServiceTicket>({
    queryKey: ["ticket", id],
    queryFn: () => ticketsApi.getById(id),
    enabled: !!id,
  });

  const { data: parts } = useQuery<{ ticket_id: string; count: number; parts: any[] }>({
    queryKey: ["ticket", id, "parts"],
    queryFn: async () => (await apiClient.get(`/v1/tickets/${id}/parts`)).data,
    enabled: !!id,
  });

  const { data: attachmentList, refetch: refetchAttachments, isLoading: loadingAttachments } = useQuery({
    queryKey: ["ticket", id, "attachments"],
    queryFn: () => attachmentsApi.list({ ticket_id: String(id), page_size: 50 }),
  });

  const [uploading, setUploading] = useState(false);
  const onUpload = async (file: File) => {
    try {
      setUploading(true);
      await attachmentsApi.upload({ file, ticketId: String(id), category: "issue_photo", source: "admin_ui" });
      await refetchAttachments();
    } finally {
      setUploading(false);
    }
  };

  const [engineerName, setEngineerName] = useState("");

  const assign = useMutation({
    mutationFn: async () => {
      return (await apiClient.post(`/v1/tickets/${id}/assign`, { engineer_name: engineerName, engineer_id: engineerName, assigned_by: "admin" })).data;
    },
    onSuccess: () => { qc.invalidateQueries({ queryKey: ["ticket", id] }); setEngineerName(""); },
  });

  const post = (path: string, body?: any) => apiClient.post(path, body ?? {}).then(r => r.data);

  const ack = useMutation({ mutationFn: () => post(`/v1/tickets/${id}/acknowledge`, { acknowledged_by: "admin" }), onSuccess: () => refetch() });
  const start = useMutation({ mutationFn: () => post(`/v1/tickets/${id}/start`, { started_by: "admin" }), onSuccess: () => refetch() });
  const hold = useMutation({ mutationFn: () => post(`/v1/tickets/${id}/hold`, { reason: "Awaiting parts", changed_by: "admin" }), onSuccess: () => refetch() });
  const resume = useMutation({ mutationFn: () => post(`/v1/tickets/${id}/resume`, { resumed_by: "admin" }), onSuccess: () => refetch() });
  const resolve = useMutation({ mutationFn: () => post(`/v1/tickets/${id}/resolve`, { resolution_notes: "Resolved by admin" }), onSuccess: () => refetch() });
  const close = useMutation({ mutationFn: () => post(`/v1/tickets/${id}/close`, { closed_by: "admin" }), onSuccess: () => refetch() });
  const cancel = useMutation({ mutationFn: () => post(`/v1/tickets/${id}/cancel`, { reason: "Cancelled by admin", cancelled_by: "admin" }), onSuccess: () => refetch() });

  if (isLoading || !ticket) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button onClick={() => router.back()} className="p-2 rounded hover:bg-gray-100"><ArrowLeft className="h-5 w-5" /></button>
            <h1 className="text-lg font-semibold">Ticket {ticket.ticket_number}</h1>
            <StatusBadge status={ticket.status} />
          </div>
          <Link href="/tickets" className="text-sm text-blue-600 hover:underline">All tickets</Link>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6 grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left: details */}
        <div className="lg:col-span-2 space-y-4">
          <div className="bg-white border rounded p-4">
            <h2 className="text-base font-semibold mb-3">Details</h2>
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div className="flex items-center gap-2"><Package className="h-4 w-4 text-gray-400" /><span className="text-gray-500">Equipment</span><span className="font-medium">{ticket.equipment_name}</span></div>
              <div className="flex items-center gap-2"><User className="h-4 w-4 text-gray-400" /><span className="text-gray-500">Customer</span><span className="font-medium">{ticket.customer_name}</span></div>
              <div className="flex items-center gap-2"><Calendar className="h-4 w-4 text-gray-400" /><span className="text-gray-500">Created</span><span className="font-medium">{new Date(ticket.created_at).toLocaleString()}</span></div>
              <div><span className="text-gray-500">Priority</span><div className="mt-1"><span className="px-2 py-0.5 rounded text-xs font-medium bg-gray-100">{ticket.priority}</span></div></div>
            </div>
            <div className="mt-4">
              <p className="text-xs text-gray-500 mb-1">Issue</p>
              <p className="text-sm whitespace-pre-line">{ticket.issue_description}</p>
            </div>
          </div>

          <div className="bg-white border rounded p-4">
            <h2 className="text-base font-semibold mb-3 flex items-center gap-2"><MessageSquare className="h-4 w-4" /> Comments</h2>
            {/* Simple add comment box */}
            <CommentBox ticketId={id} onAdded={() => qc.invalidateQueries({ queryKey: ["ticket", id, "comments"] })} />
            <CommentsList ticketId={id} />
          </div>

          <div className="bg-white border rounded p-4">
            <div className="flex items-center justify-between mb-3">
              <h2 className="text-base font-semibold flex items-center gap-2"><Paperclip className="h-4 w-4" /> Attachments</h2>
              <label className="inline-flex items-center gap-2 px-3 py-1.5 border rounded text-sm cursor-pointer">
                <Upload className="h-4 w-4" /> {uploading ? "Uploading..." : "Upload"}
                <input type="file" className="hidden" onChange={(e) => { const f = e.target.files?.[0]; if (f) onUpload(f); e.currentTarget.value = ""; }} disabled={uploading} />
              </label>
            </div>
            {loadingAttachments ? (
              <p className="text-sm text-gray-500">Loading attachments...</p>
            ) : attachmentList?.items?.length ? (
              <ul className="divide-y">
                {attachmentList.items.map(a => (
                  <li key={a.id} className="py-2 text-sm flex items-center justify-between">
                    <div>
                      <div className="font-medium">{a.fileName}</div>
                      <div className="text-gray-500">{(a.fileSize/1024).toFixed(1)} KB • {new Date(a.uploadDate).toLocaleString()}</div>
                    </div>
                    <span className="px-2 py-0.5 rounded text-xs bg-gray-100">{a.status}</span>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-sm text-gray-500">No attachments yet.</p>
            )}
          </div>

          <div className="bg-white border rounded p-4">
            <h2 className="text-base font-semibold mb-3 flex items-center gap-2"><Package className="h-4 w-4" /> Parts</h2>
            {parts?.parts?.length ? (
              <ul className="divide-y">
                {parts.parts.map((p) => (
                  <li key={p.spare_part_id} className="py-2 flex items-center justify-between text-sm">
                    <div>
                      <div className="font-medium">{p.part_name}</div>
                      <div className="text-gray-500">{p.part_number}</div>
                    </div>
                    <div className="text-right text-gray-600">
                      {p.quantity_required ? <div>Qty: {p.quantity_required}</div> : null}
                      {p.unit_price ? <div>₹{p.unit_price}</div> : null}
                    </div>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-sm text-gray-500">No parts suggested.</p>
            )}
          </div>
        </div>

        {/* Right: actions */}
        <div className="space-y-4">
          <div className="bg-white border rounded p-4">
            <h3 className="text-sm font-semibold mb-3">Assign engineer</h3>
            <div className="flex gap-2">
              <input value={engineerName} onChange={(e) => setEngineerName(e.target.value)} placeholder="Engineer name or ID" className="flex-1 border rounded px-3 py-2 text-sm" />
              <button onClick={() => assign.mutate()} disabled={!engineerName || assign.isLoading} className="px-3 py-2 bg-indigo-600 text-white rounded text-sm disabled:opacity-50">{assign.isLoading ? "Assigning..." : "Assign"}</button>
            </div>
            {ticket.assigned_engineer_name && (
              <p className="text-xs text-gray-500 mt-2">Assigned to: <span className="font-medium">{ticket.assigned_engineer_name}</span></p>
            )}
          </div>

          <div className="bg-white border rounded p-4">
            <h3 className="text-sm font-semibold mb-3">Actions</h3>
            <div className="grid grid-cols-2 gap-2">
              <button onClick={() => ack.mutate()} className="px-3 py-2 border rounded text-sm flex items-center justify-center gap-2"><CheckCircle className="h-4 w-4" /> Acknowledge</button>
              <button onClick={() => start.mutate()} className="px-3 py-2 border rounded text-sm flex items-center justify-center gap-2"><Wrench className="h-4 w-4" /> Start</button>
              <button onClick={() => hold.mutate()} className="px-3 py-2 border rounded text-sm flex items-center justify-center gap-2"><Pause className="h-4 w-4" /> Hold</button>
              <button onClick={() => resume.mutate()} className="px-3 py-2 border rounded text-sm flex items-center justify-center gap-2"><Play className="h-4 w-4" /> Resume</button>
              <button onClick={() => resolve.mutate()} className="px-3 py-2 border rounded text-sm flex items-center justify-center gap-2"><CheckCircle className="h-4 w-4" /> Resolve</button>
              <button onClick={() => close.mutate()} className="px-3 py-2 border rounded text-sm flex items-center justify-center gap-2"><FileText className="h-4 w-4" /> Close</button>
              <button onClick={() => cancel.mutate()} className="px-3 py-2 border rounded text-sm flex items-center justify-center gap-2 col-span-2 text-red-600 border-red-300"><XCircle className="h-4 w-4" /> Cancel</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function CommentsList({ ticketId }: { ticketId: string }) {
  const { data, isLoading } = useQuery<{ comments: any[] }>({
    queryKey: ["ticket", ticketId, "comments"],
    queryFn: () => ticketsApi.getComments(ticketId),
  });
  if (isLoading) return <p className="text-sm text-gray-500">Loading comments...</p>;
  return (
    <div className="mt-3 border rounded divide-y">
      {data?.comments?.length ? data.comments.map((c) => (
        <div key={c.id} className="p-3 text-sm">
          <div className="flex items-center justify-between">
            <div className="font-medium">{c.author_name || "User"}</div>
            <div className="text-xs text-gray-500">{new Date(c.created_at).toLocaleString()}</div>
          </div>
          <p className="mt-1 whitespace-pre-line">{c.comment}</p>
        </div>
      )) : <div className="p-3 text-sm text-gray-500">No comments yet.</div>}
    </div>
  );
}

function CommentBox({ ticketId, onAdded }: { ticketId: string; onAdded: () => void }) {
  const [text, setText] = useState("");
  const add = useMutation({
    mutationFn: () => ticketsApi.addComment(ticketId, { comment: text }),
    onSuccess: () => { setText(""); onAdded(); },
  });
  return (
    <div className="flex gap-2 mb-3">
      <input value={text} onChange={(e) => setText(e.target.value)} placeholder="Add a comment" className="flex-1 border rounded px-3 py-2 text-sm" />
      <button onClick={() => add.mutate()} disabled={!text} className="px-3 py-2 bg-blue-600 text-white rounded text-sm disabled:opacity-50">Post</button>
    </div>
  );
}
