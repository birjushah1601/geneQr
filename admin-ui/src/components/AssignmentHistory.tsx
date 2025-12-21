'use client';

import { useState, useEffect } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { 
  History, 
  User, 
  Calendar, 
  Clock,
  ArrowRight,
  Loader2
} from 'lucide-react';
import apiClient from '@/lib/api/client';

interface Assignment {
  id: string;
  ticket_id: string;
  engineer_id: string;
  engineer_name: string;
  organization_name: string;
  assignment_tier: string;
  assigned_at: string;
  assigned_by?: string;
  reason?: string;
  status: string;
}

interface AssignmentHistoryProps {
  ticketId: string;
}

export default function AssignmentHistory({ ticketId }: AssignmentHistoryProps) {
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (ticketId) {
      fetchAssignmentHistory();
    }
  }, [ticketId]);

  const fetchAssignmentHistory = async () => {
    setIsLoading(true);
    try {
      const response = await apiClient.get(`/v1/tickets/${ticketId}/assignments/history`);
      setAssignments(response.data.assignments || []);
    } catch (err) {
      console.error('Failed to fetch assignment history:', err);
      setAssignments([]);
    } finally {
      setIsLoading(false);
    }
  };

  const getTierLabel = (tier: string) => {
    switch (tier) {
      case 'tier_1':
        return 'Tier 1: OEM';
      case 'tier_2':
        return 'Tier 2: Partner';
      case 'tier_3':
        return 'Tier 3: Multi-brand';
      case 'tier_4':
        return 'Tier 4: Hospital';
      default:
        return tier;
    }
  };

  const getTierColor = (tier: string) => {
    switch (tier) {
      case 'tier_1':
        return 'bg-purple-100 text-purple-800 border-purple-300';
      case 'tier_2':
        return 'bg-blue-100 text-blue-800 border-blue-300';
      case 'tier_3':
        return 'bg-green-100 text-green-800 border-green-300';
      case 'tier_4':
        return 'bg-gray-100 text-gray-800 border-gray-300';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'completed':
        return 'bg-blue-100 text-blue-800';
      case 'reassigned':
        return 'bg-yellow-100 text-yellow-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <History className="h-5 w-5" />
            Assignment History
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (assignments.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <History className="h-5 w-5" />
            Assignment History
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8 text-gray-500">
            <History className="h-12 w-12 mx-auto mb-3 text-gray-400" />
            <p>No assignment history yet</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <History className="h-5 w-5" />
          Assignment History
          <Badge variant="outline" className="ml-2">
            {assignments.length} {assignments.length === 1 ? 'assignment' : 'assignments'}
          </Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {assignments.map((assignment, index) => (
            <div
              key={assignment.id}
              className="border rounded-lg p-4 relative"
            >
              {/* Connection line to next item */}
              {index < assignments.length - 1 && (
                <div className="absolute left-8 top-full w-0.5 h-4 bg-gray-300" />
              )}

              <div className="flex items-start gap-4">
                {/* Timeline dot */}
                <div className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 ${
                  assignment.status === 'active' 
                    ? 'bg-green-100 text-green-600' 
                    : 'bg-gray-100 text-gray-600'
                }`}>
                  <User className="h-4 w-4" />
                </div>

                {/* Content */}
                <div className="flex-1">
                  <div className="flex items-start justify-between mb-2">
                    <div>
                      <h4 className="font-semibold text-gray-900 flex items-center gap-2">
                        {assignment.engineer_name}
                        {assignment.status === 'active' && (
                          <Badge className="bg-green-100 text-green-800">
                            Current
                          </Badge>
                        )}
                      </h4>
                      <p className="text-sm text-gray-600">{assignment.organization_name}</p>
                    </div>
                    <Badge className={getStatusColor(assignment.status)}>
                      {assignment.status}
                    </Badge>
                  </div>

                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div className="flex items-center gap-2 text-gray-600">
                      <Calendar className="h-4 w-4" />
                      {formatDate(assignment.assigned_at)}
                    </div>
                    <div>
                      <Badge className={getTierColor(assignment.assignment_tier)}>
                        {getTierLabel(assignment.assignment_tier)}
                      </Badge>
                    </div>
                  </div>

                  {assignment.reason && (
                    <div className="mt-3 p-3 bg-yellow-50 border border-yellow-200 rounded text-sm">
                      <strong className="text-yellow-800">Reason for reassignment:</strong>
                      <p className="text-yellow-700 mt-1">{assignment.reason}</p>
                    </div>
                  )}

                  {assignment.assigned_by && (
                    <p className="mt-2 text-xs text-gray-500">
                      Assigned by: {assignment.assigned_by}
                    </p>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Summary */}
        <div className="mt-6 pt-4 border-t">
          <div className="flex items-center justify-between text-sm text-gray-600">
            <span>Total Assignments: {assignments.length}</span>
            {assignments.length > 1 && (
              <span className="text-yellow-600">
                {assignments.length - 1} reassignment{assignments.length > 2 ? 's' : ''}
              </span>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
