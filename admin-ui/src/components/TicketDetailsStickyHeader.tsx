"use client";

import Link from "next/link";
import { ArrowLeft, Clock, Package, User, Mail, Edit2 } from "lucide-react";
import type { ServiceTicket, TicketStatus, TicketPriority } from "@/types";
import { Button } from "@/components/ui/button";

interface TicketDetailsStickyHeaderProps {
  ticket: ServiceTicket;
  onEditTimeline?: () => void;
  onSendNotification?: () => void;
  onReassign?: () => void;
  onAIDiagnosis?: () => void;
  onPriorityChange?: (priority: TicketPriority) => void;
}

function StatusBadge({ status }: { status: TicketStatus }) {
  const config = {
    new: { bg: "bg-gray-100", text: "text-gray-700", label: "New" },
    assigned: { bg: "bg-indigo-100", text: "text-indigo-700", label: "Assigned" },
    in_progress: { bg: "bg-blue-100", text: "text-blue-700", label: "In Progress" },
    on_hold: { bg: "bg-yellow-100", text: "text-yellow-800", label: "On Hold" },
    resolved: { bg: "bg-green-100", text: "text-green-700", label: "Resolved" },
    closed: { bg: "bg-gray-200", text: "text-gray-800", label: "Closed" },
    cancelled: { bg: "bg-red-100", text: "text-red-700", label: "Cancelled" },
  }[status] || { bg: "bg-gray-100", text: "text-gray-700", label: status };

  return (
    <span className={`${config.bg} ${config.text} px-3 py-1 rounded-full text-sm font-semibold`}>
      {config.label}
    </span>
  );
}

function PriorityBadge({ priority, onChange }: { priority: TicketPriority; onChange?: (p: TicketPriority) => void }) {
  const config = {
    low: { bg: "bg-green-100", text: "text-green-700" },
    medium: { bg: "bg-yellow-100", text: "text-yellow-700" },
    high: { bg: "bg-orange-100", text: "text-orange-700" },
    critical: { bg: "bg-red-100", text: "text-red-700" },
  }[priority];

  return (
    <span className={`${config.bg} ${config.text} px-2 py-1 rounded text-xs font-semibold uppercase`}>
      {priority}
    </span>
  );
}

export function TicketDetailsStickyHeader({
  ticket,
  onEditTimeline,
  onSendNotification,
  onReassign,
  onAIDiagnosis,
  onPriorityChange,
}: TicketDetailsStickyHeaderProps) {
  return (
    <div className="sticky top-0 z-50 bg-white border-b shadow-sm">
      <div className="container mx-auto px-4 py-3">
        {/* Top Row: Back button, Ticket Number, Status, Actions */}
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-4">
            <Link
              href="/tickets"
              className="flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 transition"
            >
              <ArrowLeft className="h-4 w-4" />
              <span className="hidden sm:inline">Back</span>
            </Link>
            <h1 className="text-xl font-bold text-gray-900">{ticket.ticket_number}</h1>
            <StatusBadge status={ticket.status} />
            <PriorityBadge priority={ticket.priority} onChange={onPriorityChange} />
          </div>

          {/* Action Buttons */}
          <div className="flex items-center gap-2">
            {onEditTimeline && (
              <Button
                onClick={onEditTimeline}
                variant="outline"
                size="sm"
                className="hidden md:flex items-center gap-2"
              >
                <Clock className="h-4 w-4" />
                <span>Timeline</span>
              </Button>
            )}
            {onSendNotification && (
              <Button
                onClick={onSendNotification}
                variant="outline"
                size="sm"
                className="hidden md:flex items-center gap-2"
              >
                <Mail className="h-4 w-4" />
                <span>Notify</span>
              </Button>
            )}
            {onReassign && (
              <Button
                onClick={onReassign}
                variant="outline"
                size="sm"
                className="hidden sm:flex items-center gap-2"
              >
                <User className="h-4 w-4" />
                <span>Reassign</span>
              </Button>
            )}
            {onAIDiagnosis && (
              <Button
                onClick={onAIDiagnosis}
                variant="outline"
                size="sm"
                className="hidden lg:flex items-center gap-2"
              >
                <Edit2 className="h-4 w-4" />
                <span>AI Diagnosis</span>
              </Button>
            )}
          </div>
        </div>

        {/* Bottom Row: Metadata (compact) */}
        <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-gray-600">
          <div className="flex items-center gap-1.5">
            <Package className="h-3.5 w-3.5 text-gray-400" />
            <span className="font-medium text-gray-900">{ticket.equipment_name || "N/A"}</span>
          </div>
          <span className="text-gray-300">•</span>
          <div className="flex items-center gap-1.5">
            <User className="h-3.5 w-3.5 text-gray-400" />
            <span>{ticket.customer_name || "No customer"}</span>
          </div>
          {ticket.assigned_engineer_name && (
            <>
              <span className="text-gray-300">•</span>
              <div className="flex items-center gap-1.5">
                <span className="text-xs text-gray-500">Engineer:</span>
                <span className="font-medium text-blue-600">{ticket.assigned_engineer_name}</span>
              </div>
            </>
          )}
          <span className="text-gray-300">•</span>
          <div className="flex items-center gap-1.5">
            <Clock className="h-3.5 w-3.5 text-gray-400" />
            <span className="text-xs">{new Date(ticket.created_at).toLocaleDateString()}</span>
          </div>
        </div>
      </div>
    </div>
  );
}
