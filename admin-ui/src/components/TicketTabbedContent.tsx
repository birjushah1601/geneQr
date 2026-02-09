"use client";

import { useState } from "react";
import { MessageSquare, Package, Paperclip, History } from "lucide-react";

interface Tab {
  id: string;
  label: string;
  icon: React.ReactNode;
  badge?: number;
}

interface TicketTabbedContentProps {
  comments?: React.ReactNode;
  parts?: React.ReactNode;
  attachments?: React.ReactNode;
  history?: React.ReactNode;
  commentsCount?: number;
  partsCount?: number;
  attachmentsCount?: number;
  historyCount?: number;
}

export function TicketTabbedContent({
  comments,
  parts,
  attachments,
  history,
  commentsCount = 0,
  partsCount = 0,
  attachmentsCount = 0,
  historyCount = 0,
}: TicketTabbedContentProps) {
  const [activeTab, setActiveTab] = useState("comments");

  const tabs: Tab[] = [
    {
      id: "comments",
      label: "Comments",
      icon: <MessageSquare className="h-4 w-4" />,
      badge: commentsCount,
    },
    {
      id: "parts",
      label: "Parts",
      icon: <Package className="h-4 w-4" />,
      badge: partsCount,
    },
    {
      id: "attachments",
      label: "Attachments",
      icon: <Paperclip className="h-4 w-4" />,
      badge: attachmentsCount,
    },
    {
      id: "history",
      label: "History",
      icon: <History className="h-4 w-4" />,
      badge: historyCount,
    },
  ];

  return (
    <div className="bg-white border rounded-lg shadow-sm">
      {/* Tab Navigation */}
      <div className="border-b">
        <div className="flex overflow-x-auto">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors whitespace-nowrap
                ${
                  activeTab === tab.id
                    ? "text-blue-600 border-b-2 border-blue-600 bg-blue-50"
                    : "text-gray-600 hover:text-gray-900 hover:bg-gray-50"
                }
              `}
            >
              {tab.icon}
              <span>{tab.label}</span>
              {tab.badge ? (
                <span
                  className={`
                    ml-1 px-2 py-0.5 rounded-full text-xs font-semibold
                    ${
                      activeTab === tab.id
                        ? "bg-blue-600 text-white"
                        : "bg-gray-200 text-gray-600"
                    }
                  `}
                >
                  {tab.badge}
                </span>
              ) : null}
            </button>
          ))}
        </div>
      </div>

      {/* Tab Content */}
      <div className="p-4">
        {activeTab === "comments" && (comments || <EmptyState message="No comments yet" />)}
        {activeTab === "parts" && (parts || <EmptyState message="No parts assigned" />)}
        {activeTab === "attachments" && (attachments || <EmptyState message="No attachments" />)}
        {activeTab === "history" && (history || <EmptyState message="No history available" />)}
      </div>
    </div>
  );
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="text-center py-8 text-gray-500">
      <p className="text-sm">{message}</p>
    </div>
  );
}
