"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { ticketsApi } from "@/lib/api/tickets";
import MultiModelAssignment from "@/components/MultiModelAssignment";
import { ArrowLeft, Loader2 } from "lucide-react";
import Link from "next/link";

export default function AssignEngineerPage() {
  const params = useParams();
  const router = useRouter();
  const ticketId = params.id as string;

  const { data: ticket, isLoading } = useQuery({
    queryKey: ["ticket", ticketId],
    queryFn: () => ticketsApi.getById(ticketId),
    enabled: !!ticketId,
  });

  const handleAssignmentComplete = () => {
    // Redirect back to ticket detail page
    router.push(`/tickets/${ticketId}`);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (!ticket) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600 mb-4">Ticket not found</p>
          <Link href="/tickets" className="text-blue-600 hover:underline">
            Back to tickets
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center gap-3">
            <Link
              href={`/tickets/${ticketId}`}
              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <ArrowLeft className="h-5 w-5" />
            </Link>
            <div>
              <h1 className="text-xl font-semibold text-gray-900">
                Assign Engineer
              </h1>
              <p className="text-sm text-gray-600">
                Ticket #{ticket.ticket_number} - {ticket.equipment_name}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="container mx-auto px-4 py-6">
        <MultiModelAssignment
          ticketId={ticketId}
          onAssignmentComplete={handleAssignmentComplete}
        />
      </div>
    </div>
  );
}
