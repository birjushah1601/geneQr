"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ticketsApi } from "@/lib/api/tickets";
import { apiClient } from "@/lib/api/client";
import type { ServiceTicket, TicketPriority, TicketStatus } from "@/types";
import { ArrowLeft, Loader2, Package, User, Calendar, Wrench, Pause, Play, CheckCircle, XCircle, AlertTriangle, FileText, MessageSquare, Paperclip, Upload, Brain, Sparkles, TrendingUp, Lightbulb, Shield } from "lucide-react";
import { attachmentsApi } from "@/lib/api/attachments";
import { PartsAssignmentModal } from "@/components/PartsAssignmentModal";
import { diagnosisApi, extractSymptoms } from "@/lib/api/diagnosis";
import MultiModelAssignment from "@/components/MultiModelAssignment";

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

  // Fetch AI diagnosis history for this ticket
  const { data: diagnosisHistory, isLoading: loadingDiagnosis, refetch: refetchDiagnosis } = useQuery({
    queryKey: ["ticket", id, "diagnosis"],
    queryFn: () => diagnosisApi.getHistoryByTicket(Number(id)),
    enabled: !!id,
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
          
          // Refresh ticket data and diagnosis to show AI results
          await refetch();
          await refetchDiagnosis();
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

          {/* AI Diagnosis Section */}
          {diagnosisHistory && diagnosisHistory.length > 0 && (
            <div className="bg-gradient-to-br from-purple-50 to-indigo-50 border border-purple-200 rounded-lg p-5 shadow-sm">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-base font-semibold flex items-center gap-2">
                  <Brain className="h-5 w-5 text-purple-600" />
                  <span className="bg-gradient-to-r from-purple-600 to-indigo-600 bg-clip-text text-transparent">
                    AI Diagnosis
                  </span>
                  <Sparkles className="h-4 w-4 text-yellow-500" />
                </h2>
                <span className="text-xs text-purple-600 font-medium">
                  {diagnosisHistory.length} {diagnosisHistory.length === 1 ? 'analysis' : 'analyses'}
                </span>
              </div>

              {/* Latest Diagnosis */}
              {diagnosisHistory[0] && (() => {
                const latestDiagnosis = diagnosisHistory[0];
                const confidenceColor = 
                  latestDiagnosis.confidence_level === 'HIGH' ? 'bg-green-100 text-green-800 border-green-300' :
                  latestDiagnosis.confidence_level === 'MEDIUM' ? 'bg-yellow-100 text-yellow-800 border-yellow-300' :
                  'bg-red-100 text-red-800 border-red-300';

                return (
                  <div className="space-y-4">
                    {/* Primary Diagnosis */}
                    <div className="bg-white rounded-lg p-4 border border-purple-100">
                      <div className="flex items-start justify-between mb-3">
                        <div>
                          <div className="flex items-center gap-2 mb-1">
                            <Shield className="h-4 w-4 text-purple-600" />
                            <h3 className="font-semibold text-gray-900">
                              {latestDiagnosis.primary_diagnosis.problem_type}
                            </h3>
                          </div>
                          <p className="text-xs text-gray-500">
                            {latestDiagnosis.primary_diagnosis.problem_category}
                          </p>
                        </div>
                        <div className="flex flex-col items-end gap-2">
                          <span className={`px-2 py-1 rounded text-xs font-medium border ${confidenceColor}`}>
                            {latestDiagnosis.confidence_level} Confidence ({Math.round(latestDiagnosis.confidence * 100)}%)
                          </span>
                          <span className="px-2 py-1 rounded text-xs bg-orange-100 text-orange-800 border border-orange-300">
                            {latestDiagnosis.primary_diagnosis.severity} Severity
                          </span>
                        </div>
                      </div>

                      <div className="space-y-3 text-sm">
                        <div>
                          <p className="text-gray-700 leading-relaxed">
                            {latestDiagnosis.primary_diagnosis.description}
                          </p>
                        </div>

                        {latestDiagnosis.primary_diagnosis.root_cause && (
                          <div className="bg-purple-50 rounded p-3 border border-purple-100">
                            <p className="text-xs font-semibold text-purple-900 mb-1 flex items-center gap-1">
                              <Lightbulb className="h-3 w-3" />
                              Root Cause
                            </p>
                            <p className="text-gray-700">{latestDiagnosis.primary_diagnosis.root_cause}</p>
                          </div>
                        )}

                        {latestDiagnosis.primary_diagnosis.symptoms && latestDiagnosis.primary_diagnosis.symptoms.length > 0 && (
                          <div>
                            <p className="text-xs font-semibold text-gray-600 mb-2">Detected Symptoms:</p>
                            <div className="flex flex-wrap gap-2">
                              {latestDiagnosis.primary_diagnosis.symptoms.map((symptom: string, idx: number) => (
                                <span key={idx} className="px-2 py-1 bg-blue-50 text-blue-700 rounded text-xs border border-blue-200">
                                  {symptom}
                                </span>
                              ))}
                            </div>
                          </div>
                        )}
                      </div>
                    </div>

                    {/* Recommended Actions */}
                    {latestDiagnosis.recommended_actions && latestDiagnosis.recommended_actions.length > 0 && (
                      <div className="bg-white rounded-lg p-4 border border-blue-100">
                        <h4 className="font-semibold text-sm mb-3 flex items-center gap-2 text-blue-900">
                          <TrendingUp className="h-4 w-4 text-blue-600" />
                          Recommended Actions
                        </h4>
                        <div className="space-y-2">
                          {latestDiagnosis.recommended_actions.slice(0, 3).map((action: any, idx: number) => (
                            <div key={idx} className="flex items-start gap-2 text-sm">
                              <span className="flex-shrink-0 w-5 h-5 rounded-full bg-blue-100 text-blue-700 text-xs flex items-center justify-center font-medium">
                                {idx + 1}
                              </span>
                              <div className="flex-1">
                                <div className="flex items-center gap-2">
                                  <span className="font-medium text-gray-900">{action.action}</span>
                                  <span className={`px-1.5 py-0.5 rounded text-xs ${
                                    action.priority === 'high' ? 'bg-red-100 text-red-700' :
                                    action.priority === 'medium' ? 'bg-yellow-100 text-yellow-700' :
                                    'bg-gray-100 text-gray-700'
                                  }`}>
                                    {action.priority}
                                  </span>
                                </div>
                                <p className="text-gray-600 text-xs mt-0.5">{action.description}</p>
                                {action.estimated_time && (
                                  <p className="text-gray-500 text-xs mt-1">⏱️ Est. Time: {action.estimated_time}</p>
                                )}
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* AI-Suggested Parts */}
                    {latestDiagnosis.required_parts && latestDiagnosis.required_parts.length > 0 && (
                      <div className="bg-white rounded-lg p-4 border border-green-100">
                        <h4 className="font-semibold text-sm mb-3 flex items-center gap-2 text-green-900">
                          <Package className="h-4 w-4 text-green-600" />
                          AI-Suggested Parts
                        </h4>
                        <div className="space-y-2">
                          {latestDiagnosis.required_parts.slice(0, 5).map((part: any, idx: number) => (
                            <div key={idx} className="flex items-center justify-between text-sm bg-green-50 rounded p-2 border border-green-100">
                              <div className="flex-1">
                                <div className="font-medium text-gray-900">{part.part_name}</div>
                                <div className="text-xs text-gray-500">{part.part_code} • {part.manufacturer}</div>
                              </div>
                              <div className="text-right">
                                <div className="text-xs text-gray-600">Qty: {part.quantity}</div>
                                <div className="text-xs font-medium text-green-700">
                                  {Math.round(part.probability * 100)}% match
                                </div>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Vision Analysis */}
                    {latestDiagnosis.vision_analysis && latestDiagnosis.vision_analysis.findings && latestDiagnosis.vision_analysis.findings.length > 0 && (
                      <div className="bg-white rounded-lg p-4 border border-indigo-100">
                        <h4 className="font-semibold text-sm mb-3 flex items-center gap-2 text-indigo-900">
                          <AlertTriangle className="h-4 w-4 text-indigo-600" />
                          Vision Analysis
                        </h4>
                        <p className="text-sm text-gray-700 mb-3">{latestDiagnosis.vision_analysis.overall_assessment}</p>
                        <div className="space-y-2">
                          {latestDiagnosis.vision_analysis.findings.slice(0, 3).map((finding: any, idx: number) => (
                            <div key={idx} className="flex items-start gap-2 text-sm bg-indigo-50 rounded p-2 border border-indigo-100">
                              <span className="text-indigo-600 font-medium text-xs">{finding.category}:</span>
                              <span className="text-gray-700 flex-1">{finding.finding}</span>
                              <span className="text-xs text-gray-500">{Math.round(finding.confidence * 100)}%</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* AI Metadata */}
                    <div className="flex items-center justify-between text-xs text-gray-500 pt-2 border-t border-purple-100">
                      <div className="flex items-center gap-2">
                        <span>AI Model: {latestDiagnosis.ai_metadata.model}</span>
                        <span>•</span>
                        <span>{new Date(latestDiagnosis.created_at).toLocaleString()}</span>
                      </div>
                      {latestDiagnosis.ai_metadata.alternatives_count > 0 && (
                        <span className="text-purple-600 font-medium">
                          +{latestDiagnosis.ai_metadata.alternatives_count} alternative diagnoses
                        </span>
                      )}
                    </div>
                  </div>
                );
              })()}
            </div>
          )}

          {/* Loading State for Diagnosis */}
          {loadingDiagnosis && (
            <div className="bg-purple-50 border border-purple-200 rounded-lg p-4 flex items-center justify-center gap-3">
              <Loader2 className="h-5 w-5 animate-spin text-purple-600" />
              <span className="text-sm text-purple-700">Loading AI diagnosis...</span>
            </div>
          )}

          <div className="bg-white border rounded p-4">
            <h2 className="text-base font-semibold mb-3 flex items-center gap-2"><MessageSquare className="h-4 w-4" /> Comments</h2>
            {/* Simple add comment box */}
            <CommentBox ticketId={id} onAdded={() => qc.invalidateQueries({ queryKey: ["ticket", id, "comments"] })} />
            <CommentsList ticketId={id} />
          </div>

          {/* Engineer Assignment Section */}
          {!ticket.assigned_engineer_name && (
            <div className="bg-white border rounded p-4">
              <h2 className="text-base font-semibold mb-3 flex items-center gap-2">
                <User className="h-4 w-4" /> Assign Engineer
              </h2>
              <MultiModelAssignment 
                ticketId={id} 
                onAssignmentComplete={() => refetch()}
                layout="horizontal"
              />
            </div>
          )}
        </div>

        {/* Right: actions */}
        <div className="space-y-4">
          {/* Currently Assigned Engineer */}
          {ticket.assigned_engineer_name && (
            <div className="bg-white border rounded p-4">
              <h3 className="text-sm font-semibold mb-3">Currently Assigned</h3>
              <div className="flex items-center gap-3 mb-3">
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-semibold">
                  {ticket.assigned_engineer_name.split(' ').map(n => n[0]).join('').substring(0, 2)}
                </div>
                <div>
                  <p className="font-medium text-gray-900">{ticket.assigned_engineer_name}</p>
                  <p className="text-xs text-gray-500">Assigned Engineer</p>
                </div>
              </div>
            </div>
          )}

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

          {/* Attachments Section */}
          <div className="bg-white border rounded p-4">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-semibold flex items-center gap-2">
                <Paperclip className="h-4 w-4" /> Attachments
                {aiAnalyzing && (
                  <span className="inline-flex items-center gap-1 px-2 py-1 bg-purple-100 text-purple-700 rounded text-xs">
                    <Loader2 className="h-3 w-3 animate-spin" />
                    AI Analyzing...
                  </span>
                )}
              </h3>
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
    mutationFn: () => ticketsApi.addComment(ticketId, { 
      comment: text,
      comment_type: "internal",
      author_name: "Admin User"
    }),
    onSuccess: () => { setText(""); onAdded(); },
  });
  return (
    <div className="flex gap-2 mb-3">
      <input value={text} onChange={(e) => setText(e.target.value)} placeholder="Add a comment" className="flex-1 border rounded px-3 py-2 text-sm" />
      <button onClick={() => add.mutate()} disabled={!text} className="px-3 py-2 bg-blue-600 text-white rounded text-sm disabled:opacity-50">Post</button>
    </div>
  );
}
