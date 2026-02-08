"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Clock, Package, CheckCircle, Loader2, AlertTriangle, AlertCircle, User } from "lucide-react";

interface PublicMilestone {
  stage: string;
  title: string;
  description: string;
  status: string;
  eta?: string;
  completed_at?: string;
  is_active: boolean;
}

interface PublicTimeline {
  overall_status: string;
  status_message: string;
  current_stage: string;
  current_stage_desc: string;
  next_stage?: string;
  next_stage_desc?: string;
  estimated_resolution?: string;
  time_remaining: string;
  requires_parts: boolean;
  parts_status?: string;
  parts_eta?: string;
  assigned_engineer?: string;
  priority: string;
  is_urgent: boolean;
  milestones: PublicMilestone[];
  progress_percentage: number;
}

interface TicketTimelineProps {
  timeline: PublicTimeline;
  isPublic?: boolean;
}

export function TicketTimeline({ timeline, isPublic = false }: TicketTimelineProps) {
  const formatDate = (dateStr?: string) => {
    if (!dateStr) return "TBD";
    const date = new Date(dateStr);
    return date.toLocaleString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "numeric",
      minute: "2-digit",
      hour12: true,
    });
  };

  const formatTime = (dateStr?: string) => {
    if (!dateStr) return "";
    const date = new Date(dateStr);
    return date.toLocaleString("en-US", {
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "2-digit",
      hour12: true,
    });
  };

  const getStatusBadge = (status: string) => {
    const styles = {
      on_track: "bg-green-100 text-green-700",
      at_risk: "bg-yellow-100 text-yellow-700",
      delayed: "bg-red-100 text-red-700",
      blocked: "bg-orange-100 text-orange-700",
    };
    const labels = {
      on_track: "On Track",
      at_risk: "At Risk",
      delayed: "Delayed",
      blocked: "Waiting",
    };
    return (
      <span className={`px-3 py-1 rounded-full text-sm font-semibold ${styles[status as keyof typeof styles] || styles.on_track}`}>
        {labels[status as keyof typeof labels] || status}
      </span>
    );
  };

  return (
    <div className={isPublic ? "p-4" : "space-y-6"}>
      {/* Compact Overall Status - Public View */}
      {isPublic ? (
        <div className="space-y-4">
          {/* Status Summary */}
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                {getStatusBadge(timeline.overall_status)}
                {timeline.is_urgent && (
                  <span className="px-2 py-0.5 bg-orange-100 text-orange-700 text-xs font-semibold rounded-full">
                    High Priority
                  </span>
                )}
              </div>
              <p className="text-sm text-gray-600">{timeline.status_message}</p>
            </div>
            <div className="text-right">
              <p className="text-xs text-gray-500">Target Completion</p>
              <p className="text-sm font-semibold text-gray-900">{formatDate(timeline.estimated_resolution)}</p>
              <p className="text-xs text-blue-600 font-medium">{timeline.time_remaining}</p>
            </div>
          </div>

          {/* Compact Progress Bar */}
          <div>
            <div className="flex items-center justify-between mb-1.5">
              <span className="text-xs text-gray-600">Progress</span>
              <span className="text-xs font-semibold text-blue-600">{timeline.progress_percentage}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div 
                className="bg-gradient-to-r from-blue-500 to-indigo-600 h-2 rounded-full transition-all"
                style={{ width: `${timeline.progress_percentage}%` }}
              />
            </div>
          </div>
        </div>
      ) : (
        /* Full Status Card - Admin View */
        <Card className="bg-gradient-to-r from-blue-50 to-indigo-50 border-blue-200">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-5 w-5 text-blue-600" />
                Expected Timeline
              </CardTitle>
              {getStatusBadge(timeline.overall_status)}
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-700 mb-4">
              {timeline.status_message}
            </p>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-gray-600">Target Resolution</p>
                <p className="text-xl font-semibold text-gray-900">
                  {formatDate(timeline.estimated_resolution)}
                </p>
                <p className="text-sm text-blue-600 mt-1 font-medium">
                  {timeline.time_remaining}
                </p>
              </div>
              
              {timeline.is_urgent && (
                <div className="flex items-center gap-2 text-orange-600">
                  <AlertCircle className="h-5 w-5" />
                  <div>
                    <p className="font-semibold">High Priority</p>
                    <p className="text-sm text-gray-600">Fast-tracked service</p>
                  </div>
                </div>
              )}
            </div>

            {/* Progress Bar */}
            <div className="mt-6">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-600">Overall Progress</span>
                <span className="text-sm font-semibold text-blue-600">{timeline.progress_percentage}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2.5">
                <div 
                  className="bg-gradient-to-r from-blue-500 to-indigo-600 h-2.5 rounded-full transition-all duration-500"
                  style={{ width: `${timeline.progress_percentage}%` }}
                />
              </div>
              <p className="text-xs text-gray-500 mt-2">
                Currently: {timeline.current_stage_desc}
              </p>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Parts Status - Compact for Public */}
      {timeline.requires_parts && isPublic && (
        <div className="flex items-center gap-3 p-3 bg-purple-50 border border-purple-200 rounded-lg">
          <Package className="h-5 w-5 text-purple-600 flex-shrink-0" />
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-purple-900">Parts Required</p>
            <p className="text-xs text-purple-700">
              {timeline.parts_status?.replace('_', ' ') || 'Pending'}
              {timeline.parts_eta && ` • ETA: ${formatDate(timeline.parts_eta)}`}
            </p>
          </div>
        </div>
      )}

      {/* Parts Status - Full for Admin */}
      {timeline.requires_parts && !isPublic && (
        <Card className="border-purple-200">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-purple-700">
              <Package className="h-5 w-5" />
              Parts Required
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <p className="text-sm text-gray-600">Status</p>
                <p className="text-lg font-semibold capitalize">
                  {timeline.parts_status?.replace('_', ' ') || 'Pending'}
                </p>
              </div>
              {timeline.parts_eta && (
                <div className="flex-1 text-right">
                  <p className="text-sm text-gray-600">Expected Arrival</p>
                  <p className="text-lg font-semibold text-purple-700">
                    {formatDate(timeline.parts_eta)}
                  </p>
                </div>
              )}
            </div>
            <div className="mt-4 p-3 bg-purple-50 rounded-lg">
              <p className="text-sm text-purple-800">
                📦 We've identified the required parts and are working to complete the repair once they arrive.
              </p>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Milestone Timeline - Compact for Public */}
      {isPublic ? (
        <div className="space-y-3">
          <p className="text-xs font-semibold text-gray-700 uppercase tracking-wide">Service Steps</p>
          <div className="relative">
            {timeline.milestones.map((milestone, index) => (
              <div 
                key={milestone.stage}
                className="flex gap-3 pb-3 last:pb-0"
              >
                {/* Timeline line */}
                {index < timeline.milestones.length - 1 && (
                  <div 
                    className={`absolute left-3 top-6 w-0.5 ${
                      milestone.status === 'completed' ? 'bg-green-300' : 'bg-gray-200'
                    }`}
                    style={{ height: 'calc(100% - 1.5rem)' }}
                  />
                )}
                
                {/* Milestone icon */}
                <div className={`
                  relative z-10 flex-shrink-0 w-6 h-6 rounded-full flex items-center justify-center transition-all
                  ${milestone.status === 'completed' ? 'bg-green-500 text-white' : ''}
                  ${milestone.is_active ? 'bg-blue-500 text-white' : ''}
                  ${milestone.status === 'pending' ? 'bg-gray-200 text-gray-400' : ''}
                  ${milestone.status === 'blocked' ? 'bg-yellow-500 text-white' : ''}
                  ${milestone.status === 'delayed' ? 'bg-red-500 text-white' : ''}
                `}>
                  {milestone.status === 'completed' && <CheckCircle className="h-4 w-4" />}
                  {milestone.is_active && <Loader2 className="h-4 w-4" />}
                  {milestone.status === 'pending' && <Clock className="h-3 w-3" />}
                  {milestone.status === 'blocked' && <AlertTriangle className="h-4 w-4" />}
                  {milestone.status === 'delayed' && <AlertCircle className="h-4 w-4" />}
                </div>

                {/* Milestone content */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-baseline justify-between gap-2">
                    <h3 className={`text-sm font-medium ${
                      milestone.is_active ? 'text-blue-600' : 'text-gray-900'
                    }`}>
                      {milestone.title}
                    </h3>
                    {milestone.eta && (
                      <span className="text-xs text-gray-500 whitespace-nowrap">
                        {milestone.status === 'completed' ? '✓' : formatDate(milestone.eta)}
                      </span>
                    )}
                  </div>
                  {milestone.is_active && (
                    <p className="text-xs text-gray-600 mt-0.5">{milestone.description}</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      ) : (
        /* Full Milestone Timeline - Admin View */
        <Card>
          <CardHeader>
            <CardTitle>Service Journey</CardTitle>
            <p className="text-sm text-gray-500 mt-1">Step-by-step progress of your service request</p>
          </CardHeader>
          <CardContent>
          <div className="relative">
            {timeline.milestones.map((milestone, index) => (
              <div 
                key={milestone.stage}
                className={`flex gap-4 ${
                  index === timeline.milestones.length - 1 ? 'pb-0' : 'pb-8'
                }`}
              >
                {/* Timeline line */}
                {index < timeline.milestones.length - 1 && (
                  <div 
                    className={`absolute left-4 top-10 w-0.5 h-full ${
                      milestone.status === 'completed' ? 'bg-green-300' : 'bg-gray-200'
                    }`}
                    style={{ height: 'calc(100% - 2.5rem)' }}
                  />
                )}
                
                {/* Milestone icon */}
                <div className={`
                  relative z-10 flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center transition-all
                  ${milestone.status === 'completed' ? 'bg-green-500 text-white shadow-lg' : ''}
                  ${milestone.is_active ? 'bg-blue-500 text-white animate-pulse shadow-lg' : ''}
                  ${milestone.status === 'pending' ? 'bg-gray-200 text-gray-500' : ''}
                  ${milestone.status === 'blocked' ? 'bg-yellow-500 text-white' : ''}
                  ${milestone.status === 'delayed' ? 'bg-red-500 text-white' : ''}
                `}>
                  {milestone.status === 'completed' && <CheckCircle className="h-5 w-5" />}
                  {milestone.is_active && <Loader2 className="h-5 w-5 animate-spin" />}
                  {milestone.status === 'pending' && <Clock className="h-4 w-4" />}
                  {milestone.status === 'blocked' && <AlertTriangle className="h-5 w-5" />}
                  {milestone.status === 'delayed' && <AlertCircle className="h-5 w-5" />}
                </div>

                {/* Milestone content */}
                <div className="flex-1 pb-8">
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <h3 className={`font-semibold ${
                        milestone.is_active ? 'text-blue-600 text-lg' : 'text-gray-900'
                      }`}>
                        {milestone.title}
                      </h3>
                      <p className="text-sm text-gray-600 mt-1">
                        {milestone.description}
                      </p>
                      {milestone.is_active && (
                        <p className="text-sm text-blue-600 mt-2 font-medium flex items-center gap-1">
                          <span className="inline-block w-2 h-2 bg-blue-600 rounded-full animate-pulse"></span>
                          Currently in progress
                        </p>
                      )}
                    </div>
                    <div className="text-right">
                      {milestone.completed_at && (
                        <span className="text-sm text-green-600 font-medium">
                          ✓ {formatTime(milestone.completed_at)}
                        </span>
                      )}
                      {!milestone.completed_at && milestone.eta && milestone.status !== 'completed' && (
                        <span className="text-sm text-gray-500">
                          ETA: {formatTime(milestone.eta)}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
      )}

      {/* Assigned Engineer */}
      {timeline.assigned_engineer && (
        <Card className="border-indigo-200">
          <CardContent className="py-4">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-semibold text-lg shadow-lg">
                {timeline.assigned_engineer.split(' ').map(n => n[0]).join('').substring(0, 2)}
              </div>
              <div>
                <p className="text-sm text-gray-600">Your Assigned Engineer</p>
                <p className="font-semibold text-lg text-gray-900">{timeline.assigned_engineer}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
