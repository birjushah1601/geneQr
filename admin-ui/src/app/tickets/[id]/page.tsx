"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ticketsApi } from "@/lib/api/tickets";
import { apiClient } from "@/lib/api/client";
import type { ServiceTicket, TicketPriority, TicketStatus } from "@/types";
import { ArrowLeft, Loader2, Package, User, Calendar, Wrench, Pause, Play, CheckCircle, XCircle, AlertTriangle, FileText, MessageSquare, Paperclip, Upload, Brain, Sparkles, TrendingUp, Lightbulb, Shield, Trash, X, Mail } from "lucide-react";
import { AIDiagnosisModal } from "@/components/AIDiagnosisModal";
import { attachmentsApi } from "@/lib/api/attachments";
import { PartsAssignmentModal } from "@/components/PartsAssignmentModal";
import { diagnosisApi, extractSymptoms } from "@/lib/api/diagnosis";
import MultiModelAssignment from "@/components/MultiModelAssignment";
import EngineerSelectionModal from "@/components/EngineerSelectionModal";
import AssignmentHistory from "@/components/AssignmentHistory";
import DashboardLayout from "@/components/DashboardLayout";
import { SendNotificationModal } from "@/components/SendNotificationModal";
import { useAuth } from "@/contexts/AuthContext";

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
  const { session } = useAuth();
  const [showDiagnosisModal, setShowDiagnosisModal] = useState(false);
  const [showReassignMultiModel, setShowReassignMultiModel] = useState(false);
  const [showEngineerSelection, setShowEngineerSelection] = useState(false);
  const [showNotificationModal, setShowNotificationModal] = useState(false);
  const router = useRouter();
  const qc = useQueryClient();

  const { data: ticket, isLoading, refetch } = useQuery<ServiceTicket>({
    queryKey: ["ticket", id],
    queryFn: () => ticketsApi.getById(id),
    enabled: !!id,
  });

  const { data: parts, error: partsError, isLoading: partsLoading } = useQuery<{ ticket_id: string; count: number; parts: any[] }>({
    queryKey: ["ticket", id, "parts"],
    queryFn: async () => {
      console.log(`Fetching parts for ticket: ${id}`);
      const response = await apiClient.get(`/v1/tickets/${id}/parts`);
      console.log('Parts API response:', response.data);
      return response.data;
    },
    enabled: !!id,
  });

  useEffect(() => {
    if (partsError) {
      console.error('Error fetching parts:', partsError);
    }
    if (parts) {
      console.log(`Parts loaded: ${parts.count} part(s)`, parts.parts);
    }
  }, [parts, partsError]);

  const { data: attachmentList, refetch: refetchAttachments, isLoading: loadingAttachments } = useQuery({
    queryKey: ["ticket", id, "attachments"],
    queryFn: () => attachmentsApi.list({ ticket_id: String(id), page_size: 50 }),
  });

  // Fetch AI diagnosis history for this ticket (disabled when AI is not configured)
  const { data: diagnosisHistory, isLoading: loadingDiagnosis, refetch: refetchDiagnosis } = useQuery({
    queryKey: ["ticket", id, "diagnosis"],
    queryFn: () => diagnosisApi.getHistoryByTicket(String(id)),
    enabled: false, // Disabled - AI diagnosis endpoint not implemented
  });

  const [uploading, setUploading] = useState(false);
  const [aiAnalyzing, setAiAnalyzing] = useState(false);
  
  const onDelete = async (attachmentId: string) => {
    if (!confirm('Are you sure you want to delete this attachment?')) return;
    
    try {
      await attachmentsApi.delete(attachmentId);
      await refetchAttachments();
      alert('Attachment deleted successfully!');
    } catch (error) {
      console.error('Failed to delete attachment:', error);
      alert('Failed to delete attachment. Please try again.');
    }
  };
  
  const onDeleteComment = async (commentId: string) => {
    if (!confirm('Are you sure you want to delete this comment?')) return;
    
    try {
      await ticketsApi.deleteComment(id as string, commentId);
      qc.invalidateQueries({ queryKey: ['ticket', id, 'comments'] });
      alert('Comment deleted successfully!');
    } catch (error) {
      console.error('Failed to delete comment:', error);
      alert('Failed to delete comment. Please try again.');
    }
  };

  const triggerAIAnalysis = async () => {
    if (!ticket) return;
    setAiAnalyzing(true);
    try {
      await diagnosisApi.analyze({
        ticket_id: String(id),
        equipment_id: ticket.equipment_id || "",
        symptoms: extractSymptoms(ticket.issue_description || ""),
        description: ticket.issue_description || "",
        options: {
          include_vision_analysis: true,
          include_historical_context: true,
          include_similar_tickets: true,
        }
      });
      await refetchDiagnosis();
      alert('AI analysis completed successfully!');
    } catch (error) {
      console.error('AI analysis failed:', error);
      alert('AI analysis failed. Please try again.');
    } finally {
      setAiAnalyzing(false);
    }
  };

  const handleReassignEngineer = async (engineerId: string, engineerName: string) => {
    await apiClient.post(`/v1/tickets/${id}/assign`, {
      engineer_name: engineerName,
      engineer_id: engineerId,
      assigned_by: "admin"
    });
    qc.invalidateQueries({ queryKey: ["ticket", id] });
  };

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
  const [engineerFilter, setEngineerFilter] = useState<'all' | 'own' | 'partners'>('all');

  // Fetch engineers list for dropdown (including partner engineers)
  const { data: engineersData } = useQuery({
    queryKey: ["engineers", "with-partners"],
    queryFn: () => apiClient.get("/v1/engineers?limit=100&include_partners=true"),
    staleTime: 60_000,
  });
  const allEngineers = (engineersData as any)?.data?.engineers || [];
  
  // Get current user's organization for filtering
  const userOrgId = (session as any)?.user?.organization_id;
  
  // Filter engineers based on selection
  const engineers = allEngineers.filter((eng: any) => {
    if (engineerFilter === 'all') return true;
    if (engineerFilter === 'own') return eng.organization_id === userOrgId;
    if (engineerFilter === 'partners') return eng.organization_id !== userOrgId;
    return true;
  });

  // Handle parts assignment
  const handlePartsAssign = async (assignedParts: any[]) => {
    console.log('Ticket Detail - Parts assigned:', assignedParts);
    console.log('Ticket ID:', id);
    console.log('Parts count:', assignedParts.length);
    
    try {
      // Create ticket_parts entries for each assigned part
      for (const part of assignedParts) {
        await apiClient.post(`/v1/tickets/${id}/parts`, {
          spare_part_id: part.id,
          quantity_required: part.quantity || 1,
          unit_price: part.unit_price,
          total_price: (part.unit_price || 0) * (part.quantity || 1),
          is_critical: part.requires_engineer || false,
          status: 'pending',
          notes: `Added via admin UI for ${part.part_name}`
        });
      }
      
      console.log('Parts successfully added to ticket');
      
      // Refresh the parts query
      qc.invalidateQueries({ queryKey: ["ticket", id, "parts"] });
      setIsPartsModalOpen(false);
      
      // Show success message
      alert(`Successfully assigned ${assignedParts.length} part(s) to ticket!`);
    } catch (error) {
      console.error("Failed to assign parts:", error);
      console.error("Error details:", JSON.stringify(error, null, 2));
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
    <DashboardLayout>
      <div className="bg-white border-b -mx-6 -mt-6 px-6 mb-6">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button onClick={() => router.back()} className="p-2 rounded hover:bg-gray-100"><ArrowLeft className="h-5 w-5" /></button>
            <h1 className="text-lg font-semibold">Ticket {ticket.ticket_number}</h1>
            <StatusBadge status={ticket.status} />
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={() => setShowDiagnosisModal(true)}
              className="flex items-center gap-2 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors text-sm"
            >
              <Brain className="h-4 w-4" />
              AI Diagnosis
            </button>
            <button
              onClick={() => setShowNotificationModal(true)}
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm"
            >
              <Mail className="h-4 w-4" />
              Send Notification
            </button>
            <Link href="/tickets" className="text-sm text-blue-600 hover:underline">All tickets</Link>
          </div>
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
              <div>
                <span className="text-gray-500">Priority</span>
                <div className="mt-1">
                  {(() => {
                    // Get user info from localStorage to check permissions
                    const userStr = localStorage.getItem('user');
                    const user = userStr ? JSON.parse(userStr) : null;
                    const orgType = user?.organization_type || user?.organizationType || '';
                    const canEditPriority = orgType === 'system' || orgType === 'manufacturer';
                    
                    // Debug: log to console
                    console.log('Priority edit check:', { orgType, canEditPriority, user });
                    
                    if (canEditPriority) {
                      return (
                        <select
                          value={ticket.priority}
                          onChange={async (e) => {
                            const newPriority = e.target.value;
                            if (confirm(`Change priority to ${newPriority}?`)) {
                              try {
                                await apiClient.patch(`/v1/tickets/${id}/priority`, { priority: newPriority });
                                refetch();
                              } catch (err) {
                                alert('Failed to update priority. Make sure you are logged in as an admin.');
                              }
                            }
                          }}
                          className="px-2 py-1 rounded text-xs font-medium bg-gray-100 border border-gray-300 hover:bg-gray-50 cursor-pointer"
                        >
                          <option value="low">Low</option>
                          <option value="medium">Medium</option>
                          <option value="high">High</option>
                          <option value="critical">Critical</option>
                        </select>
                      );
                    } else {
                      return (
                        <span className="px-2 py-0.5 rounded text-xs font-medium bg-gray-100">
                          {ticket.priority}
                        </span>
                      );
                    }
                  })()}
                </div>
              </div>
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
                                  <p className="text-gray-500 text-xs mt-1">â±ï¸ Est. Time: {action.estimated_time}</p>
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
            <CommentsList ticketId={id as string} onDeleteComment={onDeleteComment} />
          </div>

          {/* Engineer Assignment Section */}
          {!ticket.assigned_engineer_name && (
            <div className="bg-white border rounded p-4">
              <h2 className="text-base font-semibold mb-3 flex items-center gap-2">
                <User className="h-4 w-4" /> Assign Engineer
              </h2>
              
              {/* Smart Engineer Selection Button */}
              <button
                onClick={() => setShowEngineerSelection(true)}
                className="w-full mb-4 px-4 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 transition-all shadow-md hover:shadow-lg flex items-center justify-center gap-2 font-medium"
              >
                <Sparkles className="h-5 w-5" />
                Smart Engineer Selection
              </button>
              
              {/* Or use multi-model assignment */}
              <div className="relative mb-4">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-200"></div>
                </div>
                <div className="relative flex justify-center text-xs">
                  <span className="px-2 bg-white text-gray-500">or use manual assignment</span>
                </div>
              </div>
              
              <MultiModelAssignment 
                ticketId={id} 
                onAssignmentComplete={() => refetch()}
                layout="horizontal"
              />
            </div>
          )}
          
          {/* Assignment History - Show for all tickets */}
          {ticket.assigned_engineer_name && (
            <AssignmentHistory ticketId={id} />
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

                  <button
                    onClick={() => setShowReassignMultiModel(true)}
                    className="mt-2 text-sm text-blue-600 hover:text-blue-800 underline"
                  >
                    Reassign Engineer
                  </button>
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
            ) : attachmentList?.data?.items?.length ? (
              <>
                <ul className="divide-y">
                  {attachmentList.data.items.map(a => {
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
                                ðŸ’¡ This file can be analyzed by AI for automatic diagnosis
                              </div>
                            )}
                          </div>
                          <div className="flex items-center gap-2">
                            <span className="px-2 py-0.5 rounded text-xs bg-gray-100">{a.status}</span>
                            <button
                              onClick={() => onDelete(a.id)}
                              className="p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded transition-colors"
                              title="Delete attachment"
                            >
                              <Trash className="h-4 w-4" />
                            </button>
                          </div>
                        </div>
                      </li>
                    );
                  })}
                </ul>
                <div className="mt-3 pt-3 border-t">
                  <div className="text-sm text-gray-600">
                    <div className="flex items-center justify-between">
                      <span>Total Attachments:</span>
                      <span className="font-medium">{attachmentList.data.items.length}</span>
                    </div>
                    <div className="flex items-center justify-between mt-1">
                      <span>AI-Analyzable:</span>
                      <span className="font-medium text-purple-600">
                        {attachmentList.data.items.filter(a => 
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
                    <li key={p.id || p.spare_part_id} className="py-2 flex items-center justify-between text-sm gap-3">
                      <div className="flex-1">
                        <div className="font-medium">{p.part_name}</div>
                        <div className="text-gray-500">{p.part_number}</div>
                      </div>
                      <div className="text-right text-gray-600">
                        {p.quantity_required ? <div>Qty: {p.quantity_required}</div> : null}
                        {p.unit_price ? <div>â‚¹{p.unit_price}</div> : null}
                      </div>
                      <button
                        onClick={async () => {
                          if (confirm(`Remove ${p.part_name} from this ticket?`)) {
                            try {
                              await apiClient.delete(`/v1/tickets/${id}/parts/${p.id}`);
                              qc.invalidateQueries({ queryKey: ["ticket", id, "parts"] });
                            } catch (err) {
                              alert('Failed to remove part');
                            }
                          }
                        }}
                        className="p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded transition-colors flex-shrink-0"
                        title="Remove part"
                      >
                        <Trash className="h-4 w-4" />
                      </button>
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
                    <span className="font-medium">â‚¹{parts.parts.reduce((sum, p) => sum + ((p.unit_price || 0) * (p.quantity_required || 1)), 0).toLocaleString()}</span>
                  </div>
                </div>
              </>
            ) : (
              <p className="text-sm text-gray-500">No parts assigned yet. Click "Assign Parts" to add parts.</p>
            )}
          </div>
        </div>
      </div>

      {/* AI Diagnosis Modal */}
      <AIDiagnosisModal
        isOpen={showDiagnosisModal}
        onClose={() => setShowDiagnosisModal(false)}
        diagnosis={diagnosisHistory}
        isLoading={aiAnalyzing}
        onTriggerAnalysis={triggerAIAnalysis}
      />

      {/* Send Notification Modal */}
      {showNotificationModal && (
        <SendNotificationModal
          ticketId={ticket.id}
          ticketNumber={ticket.ticket_number}
          customerEmail={ticket.customer_phone || ''}  // TODO: Add customer_email field to ServiceTicket
          customerPhone={ticket.customer_phone}
          ticket={ticket}
          onClose={() => setShowNotificationModal(false)}
          onSuccess={() => {
            setShowNotificationModal(false);
            refetch();
          }}
        />
      )}

      {/* Engineer Reassignment - Multi-Model Assignment */}
      {showReassignMultiModel && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-7xl w-full max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b px-6 py-4 flex items-center justify-between z-10">
              <div>
                <h2 className="text-xl font-semibold text-gray-900">Reassign Engineer</h2>
                <p className="text-sm text-gray-600 mt-1">
                  Current: {ticket?.assigned_engineer_name || "None"}
                </p>
              </div>
              <button
                onClick={() => setShowReassignMultiModel(false)}
                className="text-gray-400 hover:text-gray-600 p-2 hover:bg-gray-100 rounded-full transition-colors"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            <div className="p-6">
              <MultiModelAssignment 
                ticketId={id} 
                onAssignmentComplete={() => {
                  setShowReassignMultiModel(false);
                  refetch();
                }}
                layout="horizontal"
              />
            </div>
          </div>
        </div>
      )}

      {/* Parts Assignment Modal */}
      <PartsAssignmentModal
        open={isPartsModalOpen}
        onClose={() => setIsPartsModalOpen(false)}
        onAssign={handlePartsAssign}
        equipmentId={ticket.equipment_id || "unknown"}
        equipmentName={ticket.equipment_name || "Equipment"}
        existingParts={parts?.parts || []}
      />
      
      {/* Smart Engineer Selection Modal */}
      <EngineerSelectionModal
        isOpen={showEngineerSelection}
        onClose={() => setShowEngineerSelection(false)}
        ticketId={id}
        equipmentName={ticket.equipment_name || "Equipment"}
        onAssignmentSuccess={() => {
          setShowEngineerSelection(false);
          refetch();
        }}
      />
    </DashboardLayout>
  );
}

function CommentsList({ ticketId, onDeleteComment }: { ticketId: string; onDeleteComment: (id: string) => void }) {
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
            <div className="flex items-center gap-2">
              <div className="text-xs text-gray-500">{new Date(c.created_at).toLocaleString()}</div>
              <button
                onClick={() => onDeleteComment(c.id)}
                className="p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded transition-colors"
                title="Delete comment"
              >
                <Trash className="h-3 w-3" />
              </button>
              </div>
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
