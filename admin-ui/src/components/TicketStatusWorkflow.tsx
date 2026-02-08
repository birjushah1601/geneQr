"use client";

import React from "react";
import { CheckCircle, Clock, Wrench, Pause, Play, FileCheck, XCircle, ArrowRight, Info } from "lucide-react";
import type { TicketStatus } from "@/types";

interface TicketStatusWorkflowProps {
  currentStatus: TicketStatus;
  onStatusChange?: (newStatus: TicketStatus) => void;
}

// Define the workflow rules and color scheme
const statusConfig = {
  new: {
    label: "New",
    icon: Clock,
    color: "gray",
    bgColor: "bg-gray-50",
    borderColor: "border-gray-300",
    textColor: "text-gray-700",
    badgeColor: "bg-gray-100",
    nextStates: ["assigned"],
    description: "Ticket just created, awaiting assignment"
  },
  assigned: {
    label: "Assigned",
    icon: CheckCircle,
    color: "indigo",
    bgColor: "bg-indigo-50",
    borderColor: "border-indigo-300",
    textColor: "text-indigo-700",
    badgeColor: "bg-indigo-100",
    nextStates: ["in_progress", "cancelled"],
    description: "Assigned to engineer, ready to start"
  },
  in_progress: {
    label: "In Progress",
    icon: Wrench,
    color: "blue",
    bgColor: "bg-blue-50",
    borderColor: "border-blue-300",
    textColor: "text-blue-700",
    badgeColor: "bg-blue-100",
    nextStates: ["on_hold", "resolved", "cancelled"],
    description: "Engineer actively working on the issue"
  },
  on_hold: {
    label: "On Hold",
    icon: Pause,
    color: "yellow",
    bgColor: "bg-yellow-50",
    borderColor: "border-yellow-300",
    textColor: "text-yellow-800",
    badgeColor: "bg-yellow-100",
    nextStates: ["in_progress", "cancelled"],
    description: "Paused - waiting for parts, approval, or external action"
  },
  resolved: {
    label: "Resolved",
    icon: CheckCircle,
    color: "green",
    bgColor: "bg-green-50",
    borderColor: "border-green-300",
    textColor: "text-green-700",
    badgeColor: "bg-green-100",
    nextStates: ["closed", "in_progress"],
    description: "Issue fixed, awaiting customer confirmation or closure"
  },
  closed: {
    label: "Closed",
    icon: FileCheck,
    color: "gray",
    bgColor: "bg-gray-50",
    borderColor: "border-gray-400",
    textColor: "text-gray-800",
    badgeColor: "bg-gray-200",
    nextStates: [],
    description: "Ticket completed and archived"
  },
  cancelled: {
    label: "Cancelled",
    icon: XCircle,
    color: "red",
    bgColor: "bg-red-50",
    borderColor: "border-red-300",
    textColor: "text-red-700",
    badgeColor: "bg-red-100",
    nextStates: [],
    description: "Ticket cancelled - no further action"
  }
};

// Map action buttons to their target status
const actionToStatus: { [key: string]: TicketStatus } = {
  acknowledge: "assigned",
  start: "in_progress",
  hold: "on_hold",
  resume: "in_progress",
  resolve: "resolved",
  close: "closed",
  cancel: "cancelled"
};

export function TicketStatusWorkflow({ currentStatus, onStatusChange }: TicketStatusWorkflowProps) {
  const current = statusConfig[currentStatus];
  const CurrentIcon = current.icon;

  const getActionForStatus = (targetStatus: TicketStatus): string | null => {
    if (currentStatus === "new" && targetStatus === "assigned") return "acknowledge";
    if (currentStatus === "assigned" && targetStatus === "in_progress") return "start";
    if (currentStatus === "in_progress" && targetStatus === "on_hold") return "hold";
    if (currentStatus === "on_hold" && targetStatus === "in_progress") return "resume";
    if (currentStatus === "in_progress" && targetStatus === "resolved") return "resolve";
    if (currentStatus === "resolved" && targetStatus === "closed") return "close";
    if (targetStatus === "cancelled") return "cancel";
    if (currentStatus === "resolved" && targetStatus === "in_progress") return "reopen";
    return null;
  };

  return (
    <div className="bg-white border rounded-lg p-4 space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold flex items-center gap-2">
          <Info className="h-4 w-4 text-blue-600" />
          Status Workflow
        </h3>
      </div>

      {/* Current Status */}
      <div className={`${current.bgColor} ${current.borderColor} border-2 rounded-lg p-4`}>
        <div className="flex items-center gap-3">
          <div className={`${current.badgeColor} p-2 rounded-lg`}>
            <CurrentIcon className={`h-6 w-6 ${current.textColor}`} />
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-2">
              <span className="text-xs text-gray-500 uppercase font-semibold">Current Status</span>
            </div>
            <h4 className={`text-lg font-bold ${current.textColor}`}>{current.label}</h4>
            <p className="text-xs text-gray-600 mt-1">{current.description}</p>
          </div>
        </div>
      </div>

      {/* Next Possible States */}
      {current.nextStates.length > 0 && (
        <div>
          <div className="flex items-center gap-2 mb-3">
            <ArrowRight className="h-4 w-4 text-gray-400" />
            <span className="text-xs font-semibold text-gray-600 uppercase">Available Actions</span>
          </div>
          
          <div className="grid grid-cols-1 gap-2">
            {current.nextStates.map((nextStatus) => {
              const nextConfig = statusConfig[nextStatus];
              const NextIcon = nextConfig.icon;
              const action = getActionForStatus(nextStatus);
              
              return (
                <button
                  key={nextStatus}
                  onClick={() => onStatusChange?.(nextStatus)}
                  className={`${nextConfig.bgColor} ${nextConfig.borderColor} border-2 rounded-lg p-3 
                    hover:shadow-md transition-all text-left group relative overflow-hidden`}
                >
                  {/* Hover gradient effect */}
                  <div className={`absolute inset-0 bg-gradient-to-r from-transparent to-${nextConfig.color}-100 
                    opacity-0 group-hover:opacity-50 transition-opacity`}></div>
                  
                  <div className="relative flex items-center gap-3">
                    <div className={`${nextConfig.badgeColor} p-2 rounded-lg group-hover:scale-110 transition-transform`}>
                      <NextIcon className={`h-5 w-5 ${nextConfig.textColor}`} />
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <span className={`font-semibold ${nextConfig.textColor}`}>
                          {action === "acknowledge" && "Acknowledge & Assign"}
                          {action === "start" && "Start Working"}
                          {action === "hold" && "Put On Hold"}
                          {action === "resume" && "Resume Work"}
                          {action === "resolve" && "Mark Resolved"}
                          {action === "close" && "Close Ticket"}
                          {action === "cancel" && "Cancel Ticket"}
                          {action === "reopen" && "Reopen Ticket"}
                        </span>
                        <ArrowRight className={`h-4 w-4 ${nextConfig.textColor} group-hover:translate-x-1 transition-transform`} />
                      </div>
                      <p className="text-xs text-gray-600 mt-1">{nextConfig.description}</p>
                    </div>
                  </div>
                </button>
              );
            })}
          </div>
        </div>
      )}

      {/* Visual Flow Diagram */}
      <div className="pt-3 border-t">
        <div className="flex items-center gap-2 mb-3">
          <span className="text-xs font-semibold text-gray-600 uppercase">Complete Workflow</span>
        </div>
        
        <div className="flex items-center gap-1 overflow-x-auto pb-2">
          {(["new", "assigned", "in_progress", "resolved", "closed"] as TicketStatus[]).map((status, index) => {
            const config = statusConfig[status];
            const Icon = config.icon;
            const isCurrent = status === currentStatus;
            const isPast = ["new", "assigned", "in_progress"].indexOf(currentStatus) > 
                          ["new", "assigned", "in_progress"].indexOf(status);
            
            return (
              <React.Fragment key={status}>
                <div className={`flex flex-col items-center ${isCurrent ? 'scale-110' : 'scale-90'} transition-transform`}>
                  <div 
                    className={`${isCurrent ? config.badgeColor + ' ' + config.borderColor + ' border-2' : 'bg-gray-100'} 
                      p-2 rounded-lg ${isCurrent ? 'shadow-md' : ''}`}
                  >
                    <Icon className={`h-4 w-4 ${isCurrent ? config.textColor : 'text-gray-400'}`} />
                  </div>
                  <span className={`text-[10px] mt-1 ${isCurrent ? 'font-bold ' + config.textColor : 'text-gray-500'}`}>
                    {config.label}
                  </span>
                </div>
                
                {index < 4 && (
                  <ArrowRight className={`h-4 w-4 mx-1 ${isPast || isCurrent ? 'text-gray-400' : 'text-gray-300'}`} />
                )}
              </React.Fragment>
            );
          })}
        </div>

        {/* Special States */}
        <div className="flex items-center gap-2 mt-3 flex-wrap">
          {(["on_hold", "cancelled"] as TicketStatus[]).map((status) => {
            const config = statusConfig[status];
            const Icon = config.icon;
            const isCurrent = status === currentStatus;
            
            return (
              <div 
                key={status}
                className={`${isCurrent ? config.badgeColor + ' ' + config.borderColor + ' border-2' : 'bg-gray-50'} 
                  px-2 py-1 rounded-lg flex items-center gap-1.5`}
              >
                <Icon className={`h-3 w-3 ${isCurrent ? config.textColor : 'text-gray-400'}`} />
                <span className={`text-[10px] ${isCurrent ? 'font-bold ' + config.textColor : 'text-gray-500'}`}>
                  {config.label}
                </span>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
