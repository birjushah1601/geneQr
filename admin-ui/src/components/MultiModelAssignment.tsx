"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ticketsApi, type AssignmentSuggestionsResponse, type EngineerSuggestion } from "@/lib/api/tickets";
import EngineerCard from "./EngineerCard";
import { Loader2, Target, Award, TrendingDown, UserCheck, Users, CheckCircle } from "lucide-react";

interface MultiModelAssignmentProps {
  ticketId: string;
  onAssignmentComplete?: () => void;
  layout?: "vertical" | "horizontal";
}

export default function MultiModelAssignment({ ticketId, onAssignmentComplete, layout = "vertical" }: MultiModelAssignmentProps) {
  const [activeTab, setActiveTab] = useState("best_match");
  const [selectedEngineer, setSelectedEngineer] = useState<EngineerSuggestion | null>(null);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const queryClient = useQueryClient();

  // Fetch assignment suggestions
  const { data, isLoading, error } = useQuery({
    queryKey: ["assignment-suggestions", ticketId],
    queryFn: async () => {
      const suggestions = await ticketsApi.getAssignmentSuggestions(ticketId);
      
      // Add partner engineers as a separate category
      if (suggestions && suggestions.suggestions_by_model) {
        const allEngineers = Object.values(suggestions.suggestions_by_model).flatMap(
          (model: any) => model.engineers || []
        );
        
        // Filter partner engineers (those not from the main organization)
        const mainOrgId = suggestions.equipment?.manufacturer_org_id;
        const partnerEngineers = allEngineers.filter(
          (eng: any) => eng.organization_id !== mainOrgId
        );
        
        // Add partner engineers category if there are any
        if (partnerEngineers.length > 0) {
          suggestions.suggestions_by_model.partner_engineers = {
            name: "Partner Engineers",
            description: "Engineers from partner organizations",
            engineers: partnerEngineers,
            weight: 1.0
          };
        }
      }
      
      return suggestions;
    },
    staleTime: 30_000,
  });

  // Assign engineer mutation
  const assignMutation = useMutation({
    mutationFn: async (engineer: EngineerSuggestion) => {
      // Validate engineer ID
      if (!engineer.id || engineer.id.trim() === "") {
        throw new Error("Invalid engineer: Engineer ID is missing or empty");
      }
      
      console.log("Assigning engineer:", {
        engineerId: engineer.id,
        engineerName: engineer.name,
        ticketId: ticketId
      });
      
      await ticketsApi.assignEngineerToTicket(ticketId, {
        ticket_id: ticketId,
        engineer_id: engineer.id,
        engineer_name: engineer.name, // Include engineer name
        assignment_tier: "1",
        assignment_tier_name: "Direct Assignment",
        assigned_by: "admin", // TODO: Get from auth context
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["ticket", ticketId] });
      queryClient.invalidateQueries({ queryKey: ["tickets"] });
      setShowConfirmModal(false);
      setSelectedEngineer(null);
      if (onAssignmentComplete) {
        onAssignmentComplete();
      }
    },
    onError: (error: any) => {
      console.error("Assignment failed:", error);
      alert(`Failed to assign engineer: ${error.message || "Unknown error"}`);
    },
  });

  const handleAssignClick = (engineer: EngineerSuggestion) => {
    setSelectedEngineer(engineer);
    setShowConfirmModal(true);
  };

  const confirmAssignment = () => {
    if (selectedEngineer) {
      assignMutation.mutate(selectedEngineer);
    }
  };

  const getTabIcon = (key: string) => {
    switch (key) {
      case "best_match": return <Target className="h-4 w-4" />;
      case "manufacturer_certified": return <Award className="h-4 w-4" />;
      case "low_workload": return <TrendingDown className="h-4 w-4" />;
      case "high_seniority": return <UserCheck className="h-4 w-4" />;
      case "skills_match": return <Users className="h-4 w-4" />;
      case "partner_engineers": return <Users className="h-4 w-4" />;
      default: return null;
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
        <span className="ml-2 text-gray-600">Finding best engineers...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
        <p className="font-medium">Failed to load engineer suggestions</p>
        <p className="text-sm mt-1">{error.message}</p>
      </div>
    );
  }

  if (!data) {
    return null;
  }

  const models = data.suggestions_by_model;
  const activeModel = models[activeTab];

  return (
    <div className="bg-white rounded-lg border">
      {/* Equipment Context */}
      <div className="p-4 border-b bg-gray-50">
        <h3 className="font-semibold text-gray-900 mb-2">Assignment Context</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
          <div>
            <span className="text-gray-600">Equipment:</span>
            <p className="font-medium text-gray-900">{data.equipment.name}</p>
          </div>
          <div>
            <span className="text-gray-600">Manufacturer:</span>
            <p className="font-medium text-gray-900">{data.equipment.manufacturer}</p>
          </div>
          <div>
            <span className="text-gray-600">Priority:</span>
            <p className="font-medium text-gray-900 capitalize">{data.ticket.priority}</p>
          </div>
          <div>
            <span className="text-gray-600">Min Level Required:</span>
            <p className="font-medium text-gray-900">Level {data.ticket.min_level_required}</p>
          </div>
        </div>
      </div>

      {/* Two Column Layout: Filters Left, Results Right */}
      <div className="grid grid-cols-12 gap-0 min-h-[500px]">
        {/* Left Sidebar - Model Tabs */}
        <div className="col-span-3 border-r bg-gray-50 p-4">
          <div className="space-y-2">
            {Object.entries(models).map(([key, model]) => (
              <button
                key={key}
                onClick={() => setActiveTab(key)}
                className={`w-full flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-lg transition-all ${
                  activeTab === key
                    ? 'bg-blue-600 text-white shadow-md'
                    : 'bg-white border border-gray-200 text-gray-700 hover:bg-gray-50 hover:border-gray-300'
                }`}
              >
                {getTabIcon(key)}
                <span className="flex-1 text-left">{model.model_name}</span>
                <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${
                  activeTab === key
                    ? 'bg-white/20 text-white'
                    : 'bg-gray-100 text-gray-700'
                }`}>
                  {model.count}
                </span>
              </button>
            ))}
          </div>
        </div>

        {/* Right Content - Engineers */}
        <div className="col-span-9 p-4">
          <div className="mb-4">
            <h3 className="font-semibold text-gray-900">{activeModel.model_name}</h3>
            <p className="text-sm text-gray-600 mt-1">{activeModel.description}</p>
          </div>

          {activeModel.engineers.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <Users className="h-12 w-12 mx-auto mb-3 text-gray-400" />
              <p>No engineers match this criteria</p>
            </div>
          ) : layout === "horizontal" ? (
            <div className="flex gap-4 overflow-x-auto pb-4 snap-x snap-mandatory scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-gray-100">
              {activeModel.engineers.map((engineer, index) => (
                <div key={`${engineer.id}-${index}`} className="flex-none w-80 snap-start">
                  <EngineerCard
                    engineer={engineer}
                    onAssign={() => handleAssignClick(engineer)}
                    isAssigning={assignMutation.isPending && selectedEngineer?.id === engineer.id}
                  />
                </div>
              ))}
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {activeModel.engineers.map((engineer, index) => (
                <EngineerCard
                  key={`${engineer.id}-${index}`}
                  engineer={engineer}
                  onAssign={() => handleAssignClick(engineer)}
                  isAssigning={assignMutation.isPending && selectedEngineer?.id === engineer.id}
                />
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Confirmation Modal */}
      {showConfirmModal && selectedEngineer && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">Confirm Engineer Assignment</h3>
            <p className="text-gray-600 mb-4">
              Assign <span className="font-semibold">{selectedEngineer.name}</span> to this ticket?
            </p>

            <div className="bg-gray-50 rounded-lg p-3 mb-4 text-sm">
              <p className="text-gray-700"><strong>Level:</strong> {selectedEngineer.engineer_level}</p>
              <p className="text-gray-700"><strong>Organization:</strong> {selectedEngineer.organization_name || 'N/A'}</p>
              {selectedEngineer.workload && (
                <p className="text-gray-700"><strong>Active Tickets:</strong> {selectedEngineer.workload.active_tickets}</p>
              )}
            </div>

            {assignMutation.isError && (
              <div className="bg-red-50 border border-red-200 rounded p-2 mb-4 text-sm text-red-700">
                {assignMutation.error.message}
              </div>
            )}

            <div className="flex gap-3">
              <button
                onClick={() => setShowConfirmModal(false)}
                disabled={assignMutation.isPending}
                className="flex-1 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                onClick={confirmAssignment}
                disabled={assignMutation.isPending}
                className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 flex items-center justify-center gap-2"
              >
                {assignMutation.isPending ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin" />
                    Assigning...
                  </>
                ) : (
                  <>
                    <CheckCircle className="h-4 w-4" />
                    Confirm Assignment
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
