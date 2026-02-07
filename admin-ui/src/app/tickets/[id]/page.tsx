'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2, Mail, ArrowLeft, Clock, User, Wrench, Package } from 'lucide-react';
import { toast } from 'sonner';

interface Ticket {
  id: string;
  ticket_number: string;
  status: string;
  priority: string;
  equipment_name: string;
  serial_number?: string;
  issue_description: string;
  customer_name: string;
  customer_email: string;
  customer_whatsapp?: string;
  assigned_engineer_name?: string;
  created_at: string;
  updated_at: string;
}

export default function TicketDetailPage() {
  const params = useParams();
  const router = useRouter();
  const ticketId = params.id as string;

  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [loading, setLoading] = useState(true);
  const [sendingEmail, setSendingEmail] = useState(false);

  useEffect(() => {
    fetchTicket();
  }, [ticketId]);

  const fetchTicket = async () => {
    try {
      const response = await fetch(`http://localhost:8081/api/v1/tickets/${ticketId}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch ticket');
      }

      const data = await response.json();
      setTicket(data);
    } catch (error) {
      console.error('Error fetching ticket:', error);
      toast.error('Failed to load ticket details');
    } finally {
      setLoading(false);
    }
  };

  const handleSendEmail = async () => {
    if (!ticket?.customer_email) {
      toast.error('No customer email available');
      return;
    }

    setSendingEmail(true);
    try {
      const response = await fetch(`http://localhost:8081/api/v1/tickets/${ticketId}/send-notification`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          include_comments: true,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to send email');
      }

      const result = await response.json();
      toast.success(`Email sent to ${result.recipient}`);
    } catch (error) {
      console.error('Error sending email:', error);
      toast.error('Failed to send email notification');
    } finally {
      setSendingEmail(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'open': return 'bg-blue-500';
      case 'in_progress': return 'bg-yellow-500';
      case 'resolved': return 'bg-green-500';
      case 'closed': return 'bg-gray-500';
      default: return 'bg-gray-400';
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority?.toLowerCase()) {
      case 'high': return 'bg-red-500';
      case 'medium': return 'bg-orange-500';
      case 'low': return 'bg-green-500';
      default: return 'bg-gray-400';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  if (!ticket) {
    return (
      <div className="container mx-auto py-8">
        <Card className="p-6">
          <p className="text-center text-gray-500">Ticket not found</p>
          <Button onClick={() => router.back()} className="mt-4">
            Go Back
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8 px-4">
      <div className="mb-6 flex items-center justify-between">
        <Button variant="ghost" onClick={() => router.back()}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Tickets
        </Button>

        <Button
          onClick={handleSendEmail}
          disabled={!ticket.customer_email || sendingEmail}
          className="gap-2"
        >
          {sendingEmail ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin" />
              Sending...
            </>
          ) : (
            <>
              <Mail className="w-4 h-4" />
              Send Email Update
            </>
          )}
        </Button>
      </div>

      <Card className="p-6">
        <div className="flex items-start justify-between mb-6">
          <div>
            <h1 className="text-3xl font-bold mb-2">
              {ticket.ticket_number}
            </h1>
            <div className="flex gap-2">
              <Badge className={`${getStatusColor(ticket.status)} text-white`}>
                {ticket.status}
              </Badge>
              <Badge className={`${getPriorityColor(ticket.priority)} text-white`}>
                {ticket.priority} Priority
              </Badge>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Equipment Details */}
          <div>
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <Package className="w-5 h-5" />
              Equipment Information
            </h2>
            <div className="space-y-3">
              <div>
                <label className="text-sm text-gray-500">Equipment Name</label>
                <p className="font-medium">{ticket.equipment_name}</p>
              </div>
              {ticket.serial_number && (
                <div>
                  <label className="text-sm text-gray-500">Serial Number</label>
                  <p className="font-medium">{ticket.serial_number}</p>
                </div>
              )}
            </div>
          </div>

          {/* Customer Details */}
          <div>
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <User className="w-5 h-5" />
              Customer Information
            </h2>
            <div className="space-y-3">
              <div>
                <label className="text-sm text-gray-500">Name</label>
                <p className="font-medium">{ticket.customer_name}</p>
              </div>
              <div>
                <label className="text-sm text-gray-500">Email</label>
                <p className="font-medium">{ticket.customer_email || 'Not provided'}</p>
              </div>
              {ticket.customer_whatsapp && (
                <div>
                  <label className="text-sm text-gray-500">WhatsApp</label>
                  <p className="font-medium">{ticket.customer_whatsapp}</p>
                </div>
              )}
            </div>
          </div>

          {/* Issue Description */}
          <div className="md:col-span-2">
            <h2 className="text-lg font-semibold mb-4">Issue Description</h2>
            <p className="text-gray-700 bg-gray-50 p-4 rounded-lg">
              {ticket.issue_description}
            </p>
          </div>

          {/* Assignment Details */}
          {ticket.assigned_engineer_name && (
            <div>
              <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
                <Wrench className="w-5 h-5" />
                Assigned Engineer
              </h2>
              <p className="font-medium">{ticket.assigned_engineer_name}</p>
            </div>
          )}

          {/* Timeline */}
          <div>
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <Clock className="w-5 h-5" />
              Timeline
            </h2>
            <div className="space-y-2">
              <div>
                <label className="text-sm text-gray-500">Created</label>
                <p className="font-medium">
                  {new Date(ticket.created_at).toLocaleString()}
                </p>
              </div>
              <div>
                <label className="text-sm text-gray-500">Last Updated</label>
                <p className="font-medium">
                  {new Date(ticket.updated_at).toLocaleString()}
                </p>
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}
