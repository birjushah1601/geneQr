"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { apiClient } from "@/lib/api/client";
import { Loader2, CheckCircle, Clock, AlertCircle, Package, User, Calendar, MessageSquare, History } from "lucide-react";
import { TicketTimeline } from "@/components/TicketTimeline";
import type { PublicTimeline } from "@/types";

interface PublicComment {
  comment: string;
  author_name: string;
  created_at: string;
  author_role: string;
}

interface PublicStatusEvent {
  from_status?: string;
  to_status: string;
  changed_by: string;
  changed_at: string;
  comment?: string;
}

interface PublicTicketView {
  ticket_number: string;
  status: string;
  priority: string;
  equipment_name: string;
  issue_description: string;
  created_at: string;
  updated_at: string;
  comments: PublicComment[];
  status_history: PublicStatusEvent[];
  assigned_engineer?: string;
}

export default function TrackTicketPage() {
  const params = useParams();
  const token = params.token as string;
  
  const [ticket, setTicket] = useState<PublicTicketView | null>(null);
  const [timeline, setTimeline] = useState<PublicTimeline | null>(null);
  const [loading, setLoading] = useState(true);
  const [timelineLoading, setTimelineLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchTicket = async () => {
      try {
        setLoading(true);
        const response = await apiClient.get(`/v1/track/${token}`);
        setTicket(response.data);
      } catch (err: any) {
        setError(err.response?.data?.error?.message || "Invalid or expired tracking link");
      } finally {
        setLoading(false);
      }
    };

    const fetchTimeline = async () => {
      try {
        setTimelineLoading(true);
        // Get ticket ID first from the token
        const ticketResponse = await apiClient.get(`/v1/track/${token}`);
        if (ticketResponse.data?.ticket_id) {
          const timelineResponse = await apiClient.get(`/v1/tickets/${ticketResponse.data.ticket_id}/timeline`);
          setTimeline(timelineResponse.data);
        }
      } catch (err: any) {
        console.error("Failed to fetch timeline:", err);
        // Don't show error - timeline is optional enhancement
      } finally {
        setTimelineLoading(false);
      }
    };

    if (token) {
      fetchTicket();
      fetchTimeline();
    }
  }, [token]);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "new":
        return "bg-blue-100 text-blue-800 border-blue-200";
      case "assigned":
        return "bg-purple-100 text-purple-800 border-purple-200";
      case "in_progress":
        return "bg-yellow-100 text-yellow-800 border-yellow-200";
      case "resolved":
        return "bg-green-100 text-green-800 border-green-200";
      case "closed":
        return "bg-gray-100 text-gray-800 border-gray-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case "new":
        return <AlertCircle className="h-5 w-5" />;
      case "assigned":
      case "in_progress":
        return <Clock className="h-5 w-5" />;
      case "resolved":
      case "closed":
        return <CheckCircle className="h-5 w-5" />;
      default:
        return <Clock className="h-5 w-5" />;
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case "critical":
        return "bg-red-100 text-red-800 border-red-300";
      case "high":
        return "bg-orange-100 text-orange-800 border-orange-300";
      case "medium":
        return "bg-yellow-100 text-yellow-800 border-yellow-300";
      case "low":
        return "bg-green-100 text-green-800 border-green-300";
      default:
        return "bg-gray-100 text-gray-800 border-gray-300";
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading ticket information...</p>
        </div>
      </div>
    );
  }

  if (error || !ticket) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
        <div className="bg-white rounded-lg shadow-xl p-8 max-w-md w-full text-center">
          <AlertCircle className="h-16 w-16 text-red-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Tracking Link Invalid</h1>
          <p className="text-gray-600 mb-6">{error || "This tracking link is invalid or has expired."}</p>
          <a
            href="/"
            className="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            Go to Homepage
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-3xl font-bold text-gray-900">{ticket.ticket_number}</h1>
            <div className={`px-4 py-2 rounded-full border-2 flex items-center gap-2 ${getStatusColor(ticket.status)}`}>
              {getStatusIcon(ticket.status)}
              <span className="font-semibold capitalize">{ticket.status.replace("_", " ")}</span>
            </div>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div className="flex items-center gap-2">
              <Package className="h-4 w-4 text-gray-400" />
              <span className="text-gray-600">Equipment:</span>
              <span className="font-medium">{ticket.equipment_name}</span>
            </div>
            <div className="flex items-center gap-2">
              <span className={`px-3 py-1 rounded-full text-xs font-semibold border ${getPriorityColor(ticket.priority)}`}>
                {ticket.priority.toUpperCase()} Priority
              </span>
            </div>
            <div className="flex items-center gap-2">
              <Calendar className="h-4 w-4 text-gray-400" />
              <span className="text-gray-600">Created:</span>
              <span className="font-medium">{formatDate(ticket.created_at)}</span>
            </div>
          </div>

          {ticket.assigned_engineer && (
            <div className="mt-4 p-3 bg-blue-50 rounded-lg flex items-center gap-2">
              <User className="h-5 w-5 text-blue-600" />
              <span className="text-sm text-gray-600">Assigned Engineer:</span>
              <span className="text-sm font-semibold text-blue-900">{ticket.assigned_engineer}</span>
            </div>
          )}
        </div>

        {/* Issue Description */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <MessageSquare className="h-5 w-5 text-gray-600" />
            Issue Description
          </h2>
          <p className="text-gray-700 whitespace-pre-line">{ticket.issue_description}</p>
        </div>

        {/* Status History */}
        {ticket.status_history && ticket.status_history.length > 0 && (
          <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <History className="h-5 w-5 text-gray-600" />
              Status History
            </h2>
            <div className="space-y-4">
              {ticket.status_history.map((event, index) => (
                <div key={index} className="flex items-start gap-4 pb-4 border-b last:border-b-0">
                  <div className="flex-shrink-0 mt-1">
                    <div className={`w-3 h-3 rounded-full ${getStatusColor(event.to_status).split(' ')[0]}`} />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      {event.from_status && (
                        <span className="text-sm text-gray-500 line-through">{event.from_status}</span>
                      )}
                      <span className="text-sm text-gray-400">â†’</span>
                      <span className="text-sm font-semibold text-gray-900 capitalize">
                        {event.to_status.replace("_", " ")}
                      </span>
                    </div>
                    <div className="text-xs text-gray-500">
                      {event.changed_by} â€¢ {formatDate(event.changed_at)}
                    </div>
                    {event.comment && (
                      <p className="text-sm text-gray-600 mt-1">{event.comment}</p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Service Timeline & ETA */}
        {timeline && !timelineLoading && (
          <div className="mb-6">
            <TicketTimeline timeline={timeline} />
          </div>
        )}

        {/* Comments */}
        {ticket.comments && ticket.comments.length > 0 && (
          <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <MessageSquare className="h-5 w-5 text-gray-600" />
              Updates & Comments
            </h2>
            <div className="space-y-4">
              {ticket.comments.map((comment, index) => (
                <div key={index} className="p-4 bg-gray-50 rounded-lg">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <div className="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center text-white text-sm font-semibold">
                        {comment.author_name.charAt(0).toUpperCase()}
                      </div>
                      <div>
                        <p className="text-sm font-semibold text-gray-900">{comment.author_name}</p>
                        <p className="text-xs text-gray-500 capitalize">{comment.author_role}</p>
                      </div>
                    </div>
                    <span className="text-xs text-gray-500">{formatDate(comment.created_at)}</span>
                  </div>
                  <p className="text-gray-700 whitespace-pre-line">{comment.comment}</p>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* No Activity Message */}
        {(!ticket.comments || ticket.comments.length === 0) && 
         (!ticket.status_history || ticket.status_history.length === 0) && (
          <div className="bg-white rounded-lg shadow-lg p-8 text-center">
            <Clock className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-600">No updates yet. We'll notify you once work begins!</p>
          </div>
        )}

        {/* Footer */}
        <div className="bg-white rounded-lg shadow-lg p-6 mt-6 text-center">
          <p className="text-sm text-gray-600">
            Thank you for your service request. We're working on resolving your issue as quickly as possible.
          </p>
          <p className="text-xs text-gray-500 mt-2">
            Last updated: {formatDate(ticket.updated_at)}
          </p>
        </div>
      </div>
    </div>
  );
}
