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
import { PartsAssignmentModal } from "@/components/PartsAssignmentModal";
import { diagnosisApi, extractSymptoms } from "@/lib/api/diagnosis";

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
  const [aiAnalyzing, setAiAnalyzing] = useState(false);
  
  const onUpload = async (file: File) => {
    try {
      setUploading(true);
      const uploadResult = await attachmentsApi.upload({ 
        file, 
        ticketId: String(id), 
        category: "issue_photo", 
        source: "admin_ui" 
      });
      await refetchAttachments();
      
      // Trigger AI analysis for image files
      if (file.type.startsWith('image/')) {
        setAiAnalyzing(true);
        try {
          // Convert file to base64 for AI analysis
          const base64 = await fileToBase64(file);
          
          // Call real AI diagnosis API with vision analysis
          await diagnosisApi.analyze({
            ticket_id: Number(id),
            equipment_id: ticket.equipment_id || "",
            symptoms: extractSymptoms(ticket.issue_description || ""),
            description: ticket.issue_description || "",
            images: [base64],
            options: {
              include_vision_analysis: true,
              include_historical_context: true,
              include_similar_tickets: true,
            }
          });
          
          // Refresh ticket data to show AI diagnosis results
          await refetch();
        } catch (error) {
          console.error("AI analysis failed:", error);
        } finally {
          setAiAnalyzing(false);
        }
      }
    } finally {
      setUploading(false);
    }
  };

  // Helper function to convert file to base64
  const fileToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        const base64 = reader.result as string;
        // Remove data URL prefix (e.g., "data:image/png;base64,")
        const base64String = base64.split(',')[1];
        resolve(base64String);
      };
      reader.onerror = reject;
    });
  };

  const [engineerName, setEngineerName] = useState("");
  const [isPartsModalOpen, setIsPartsModalOpen] = useState(false);

  // Fetch engineers list for dropdown
  const { data: engineersData } = useQuery({
    queryKey: ["engineers"],
    queryFn: () => apiClient.get("/engineers?limit=100"),
    staleTime: 60_000,
  });
  const engineers = (engineersData as any)?.data?.items || [];

  // Handle parts assignment
  const handlePartsAssign = async (assignedParts: any[]) => {
    try {
      // Call real API endpoint to update ticket parts
      await apiClient.patch(`/v1/tickets/${id}/parts`, {
        parts: assignedParts
      });
      
      // Refresh the parts query
      qc.invalidateQueries({ queryKey: ["ticket", id, "parts"] });
      setIsPartsModalOpen(false);
    } catch (error) {
      console.error("Failed to assign parts:", error);
      alert("Failed to assign parts. Please try again.");
    }
  };

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
              <h2 className="text-base font-semibold flex items-center gap-2">
                <Paperclip className="h-4 w-4" /> Attachments
                {aiAnalyzing && (
                  <span className="inline-flex items-center gap-1 px-2 py-1 bg-purple-100 text-purple-700 rounded text-xs">
                    <Loader2 className="h-3 w-3 animate-spin" />
                    AI Analyzing...
                  </span>
                )}
              </h2>
              <label className="inline-flex items-center gap-2 px-3 py-1.5 border rounded text-sm cursor-pointer hover:bg-gray-50 transition-colors">
                <Upload className="h-4 w-4" /> {uploading ? "Uploading..." : "Upload"}
                <input 
                  type="file" 
                  className="hidden" 
                  accept="image/*,video/*,.pdf,.doc,.docx" 
                  onChange={(e) => { const f = e.target.files?.[0]; if (f) onUpload(f); e.currentTarget.value = ""; }} 
                  disabled={uploading} 
                />
              </label>
            </div>
            {loadingAttachments ? (
              <p className="text-sm text-gray-500">Loading attachments...</p>
            ) : attachmentList?.items?.length ? (
              <>
                <ul className="divide-y">
                  {attachmentList.items.map(a => {
                    const isImage = a.fileName.match(/\.(jpg|jpeg|png|gif|webp)$/i);
                    const isVideo = a.fileName.match(/\.(mp4|mov|avi|webm)$/i);
                    return (
                      <li key={a.id} className="py-3 text-sm">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <div className="font-medium">{a.fileName}</div>
                              {(isImage || isVideo) && (
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 bg-purple-50 text-purple-600 rounded text-xs">
                                  <AlertTriangle className="h-3 w-3" />
                                  AI Ready
                                </span>
                              )}
                            </div>
                            <div className="text-gray-500 mt-1">
                              {(a.fileSize/1024).toFixed(1)} KB • {new Date(a.uploadDate).toLocaleString()}
                            </div>
                            {(isImage || isVideo) && (
                              <div className="mt-2 text-xs text-purple-600">
                                💡 This file can be analyzed by AI for automatic diagnosis
                              </div>
                            )}
                          </div>
                          <span className="px-2 py-0.5 rounded text-xs bg-gray-100 ml-2">{a.status}</span>
                        </div>
                      </li>
                    );
                  })}
                </ul>
                <div className="mt-3 pt-3 border-t">
                  <div className="text-sm text-gray-600">
                    <div className="flex items-center justify-between">
                      <span>Total Attachments:</span>
                      <span className="font-medium">{attachmentList.items.length}</span>
                    </div>
                    <div className="flex items-center justify-between mt-1">
                      <span>AI-Analyzable:</span>
                      <span className="font-medium text-purple-600">
                        {attachmentList.items.filter(a => 
                          a.fileName.match(/\.(jpg|jpeg|png|gif|webp|mp4|mov|avi|webm)$/i)
                        ).length}
                      </span>
                    </div>
                  </div>
                </div>
              </>
            ) : (
              <div className="text-center py-6">
                <Paperclip className="h-12 w-12 text-gray-300 mx-auto mb-2" />
                <p className="text-sm text-gray-500 mb-1">No attachments yet.</p>
                <p className="text-xs text-gray-400">Upload images or videos for AI-powered analysis</p>
              </div>
            )}
          </div>

          <div className="bg-white border rounded p-4">
            <div className="flex items-center justify-between mb-3">
              <h2 className="text-base font-semibold flex items-center gap-2"><Package className="h-4 w-4" /> Parts</h2>
              <button 
                onClick={() => setIsPartsModalOpen(true)}
                className="px-3 py-1.5 bg-green-600 text-white rounded text-sm hover:bg-green-700 transition-colors flex items-center gap-2"
              >
                <Package className="h-4 w-4" />
                Assign Parts
              </button>
            </div>
            {parts?.parts?.length ? (
              <>
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
                <div className="mt-3 pt-3 border-t">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Total Parts:</span>
                    <span className="font-medium">{parts.parts.length}</span>
                  </div>
                  <div className="flex justify-between text-sm mt-1">
                    <span className="text-gray-600">Total Cost:</span>
                    <span className="font-medium">₹{parts.parts.reduce((sum, p) => sum + ((p.unit_price || 0) * (p.quantity_required || 1)), 0).toLocaleString()}</span>
                  </div>
                </div>
              </>
            ) : (
              <p className="text-sm text-gray-500">No parts assigned yet. Click "Assign Parts" to add parts.</p>
            )}
          </div>
        </div>

        {/* Right: actions */}
        <div className="space-y-4">
          <div className="bg-white border rounded p-4">
            <h3 className="text-sm font-semibold mb-3">Assign engineer</h3>
            <div className="space-y-2">
              <select 
                value={engineerName} 
                onChange={(e) => setEngineerName(e.target.value)} 
                className="w-full border rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                <option value="">Select an engineer...</option>
                {engineers.map((eng: any) => (
                  <option key={eng.id} value={eng.id}>
                    {eng.name} - {eng.skills?.join(', ')} - {eng.home_region}
                  </option>
                ))}
              </select>
              <button 
                onClick={() => assign.mutate()} 
                disabled={!engineerName || assign.isLoading} 
                className="w-full px-3 py-2 bg-indigo-600 text-white rounded text-sm disabled:opacity-50 hover:bg-indigo-700 transition-colors"
              >
                {assign.isLoading ? "Assigning..." : "Assign Engineer"}
              </button>
            </div>
            {ticket.assigned_engineer_name && (
              <div className="mt-3 pt-3 border-t">
                <p className="text-xs text-gray-500">Currently assigned to:</p>
                <p className="text-sm font-medium text-gray-900 flex items-center gap-2 mt-1">
                  <User className="h-4 w-4 text-gray-400" />
                  {ticket.assigned_engineer_name}
                </p>
              </div>
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

      {/* Parts Assignment Modal */}
      <PartsAssignmentModal
        open={isPartsModalOpen}
        onClose={() => setIsPartsModalOpen(false)}
        onAssign={handlePartsAssign}
        equipmentId={ticket.equipment_id || "unknown"}
        equipmentName={ticket.equipment_name || "Equipment"}
      />
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
