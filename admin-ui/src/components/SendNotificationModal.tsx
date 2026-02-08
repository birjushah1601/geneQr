"use client";

import { useState, useEffect } from "react";
import { X, Mail, Loader2 } from "lucide-react";
import { apiClient } from "@/lib/api/client";
import { ticketsApi } from "@/lib/api/tickets";

interface SendNotificationModalProps {
  ticketId: string;
  ticketNumber: string;
  customerEmail?: string;
  customerPhone?: string;
  ticket?: any;
  onClose: () => void;
  onSuccess?: () => void;
}

export function SendNotificationModal({ 
  ticketId, 
  ticketNumber, 
  customerEmail = '',
  customerPhone = '',
  ticket,
  onClose, 
  onSuccess 
}: SendNotificationModalProps) {
  const [email, setEmail] = useState(customerEmail || "");
  const [comment, setComment] = useState("");
  const [loadingComments, setLoadingComments] = useState(true);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  // Fetch comments and generate summary on mount
  useEffect(() => {
    const generateSummary = async () => {
      try {
        // Fetch comments and timeline for the ticket
        const response = await ticketsApi.getComments(ticketId);
        const comments = response.comments || [];
        
        // Fetch timeline
        let timeline = null;
        try {
          const timelineResponse = await apiClient.get(`/v1/tickets/${ticketId}/timeline`);
          timeline = timelineResponse.data;
        } catch (err) {
          console.error('Failed to fetch timeline:', err);
        }
        
        // Generate ticket summary with history
        let summary = `Ticket Update: ${ticketNumber}\n\n`;
        summary += `Equipment: ${ticket?.equipment_name || 'N/A'}\n`;
        summary += `Status: ${ticket?.status || 'N/A'}\n`;
        summary += `Priority: ${ticket?.priority || 'N/A'}\n`;
        summary += `Created: ${ticket?.created_at ? new Date(ticket.created_at).toLocaleDateString() : 'N/A'}\n\n`;
        
        if (ticket?.assigned_engineer_name) {
          summary += `Assigned Engineer: ${ticket.assigned_engineer_name}\n\n`;
        }
        
        summary += `Issue: ${ticket?.issue_description || 'N/A'}\n\n`;
        
        // Add timeline information
        if (timeline) {
          summary += `📅 Service Timeline:\n`;
          summary += `${'='.repeat(50)}\n\n`;
          
          summary += `Expected Resolution: ${timeline.estimated_resolution ? new Date(timeline.estimated_resolution).toLocaleString() : 'TBD'}\n`;
          summary += `Time Remaining: ${timeline.time_remaining}\n`;
          summary += `Progress: ${timeline.progress_percentage}%\n`;
          summary += `Status: ${timeline.status_message}\n\n`;
          
          if (timeline.requires_parts) {
            summary += `📦 Parts Required: ${timeline.parts_status || 'Pending'}\n`;
            if (timeline.parts_eta) {
              summary += `Parts Expected: ${new Date(timeline.parts_eta).toLocaleString()}\n`;
            }
            summary += `\n`;
          }
          
          summary += `Service Journey:\n`;
          timeline.milestones.forEach((m: any) => {
            const icon = m.status === 'completed' ? '✓' : 
                         m.is_active ? '→' : '○';
            summary += `${icon} ${m.title}`;
            if (m.completed_at) {
              summary += ` - Completed: ${new Date(m.completed_at).toLocaleDateString()}`;
            } else if (m.eta && !m.completed_at) {
              summary += ` - ETA: ${new Date(m.eta).toLocaleDateString()}`;
            }
            summary += `\n`;
          });
          summary += `\n`;
        }
        
        if (comments && comments.length > 0) {
          summary += `Activity History:\n`;
          summary += `${'='.repeat(50)}\n\n`;
          
          comments.forEach((comment: any) => {
            const date = new Date(comment.created_at).toLocaleString();
            summary += `[${date}] ${comment.author_name || 'System'}:\n`;
            summary += `${comment.comment}\n\n`;
          });
        } else {
          summary += `No comments yet.\n\n`;
        }
        
        if (ticket?.tracking_url) {
          summary += `\nYou can track this ticket anytime at: ${ticket.tracking_url}\n`;
        }
        
        setComment(summary);
      } catch (err) {
        console.error('Failed to fetch comments:', err);
        // Fallback to basic summary
        let summary = `Ticket Update: ${ticketNumber}\n\n`;
        summary += `Equipment: ${ticket?.equipment_name || 'N/A'}\n`;
        summary += `Status: ${ticket?.status || 'N/A'}\n`;
        summary += `\nPlease contact support for more details.`;
        setComment(summary);
      } finally {
        setLoadingComments(false);
      }
    };

    generateSummary();
  }, [ticketId, ticketNumber, ticket]);

  const handleSend = async () => {
    if (!email || !comment) {
      setError("Please provide both email and comment");
      return;
    }

    // Basic email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setError("Please provide a valid email address");
      return;
    }

    setSending(true);
    setError("");

    try {
      await apiClient.post(`/v1/tickets/${ticketId}/notify`, {
        email,
        comment,
      });

      setSuccess(true);
      setTimeout(() => {
        onSuccess?.();
        onClose();
      }, 1500);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || "Failed to send notification");
      setSending(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-lg w-full">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-blue-100 rounded-lg">
              <Mail className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <h2 className="text-lg font-semibold">Send Notification</h2>
              <p className="text-sm text-gray-500">Ticket {ticketNumber}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 p-2 hover:bg-gray-100 rounded-full transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Body */}
        <div className="p-6 space-y-4">
          {success ? (
            <div className="bg-green-50 border border-green-200 rounded-lg p-4 text-center">
              <div className="text-green-600 text-lg mb-2">✓</div>
              <p className="text-green-800 font-medium">Notification sent successfully!</p>
            </div>
          ) : (
            <>
              {/* Email Input */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Recipient Email
                  {!customerEmail && customerPhone && (
                    <span className="ml-2 text-xs text-amber-600">
                      (No email on record - Phone: {customerPhone})
                    </span>
                  )}
                </label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder={customerEmail ? customerEmail : "Enter customer email address"}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  disabled={sending}
                />
                {!customerEmail && (
                  <p className="text-xs text-gray-500 mt-1">
                    This ticket was created before email field was added. Please enter customer's email manually.
                  </p>
                )}
              </div>

              {/* Comment Input */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Message Summary
                  {loadingComments && (
                    <span className="ml-2 text-xs text-gray-500">
                      <Loader2 className="inline h-3 w-3 animate-spin" /> Loading ticket history...
                    </span>
                  )}
                </label>
                <textarea
                  value={comment}
                  onChange={(e) => setComment(e.target.value)}
                  placeholder="Loading ticket summary..."
                  rows={15}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-vertical font-mono text-sm"
                  disabled={sending || loadingComments}
                />
                <p className="text-xs text-gray-500 mt-1">
                  Pre-filled with ticket details and comment history. Edit as needed.
                </p>
              </div>

              {/* Error Message */}
              {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                  <p className="text-sm text-red-800">{error}</p>
                </div>
              )}
            </>
          )}
        </div>

        {/* Footer */}
        {!success && (
          <div className="flex items-center justify-end gap-3 p-6 border-t bg-gray-50">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
              disabled={sending}
            >
              Cancel
            </button>
            <button
              onClick={handleSend}
              disabled={sending || !email || !comment}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              {sending ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Sending...
                </>
              ) : (
                <>
                  <Mail className="h-4 w-4" />
                  Send Notification
                </>
              )}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
