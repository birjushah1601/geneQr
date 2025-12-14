"use client";

import { useState } from "react";
import { User, X, Loader2, UserCheck } from "lucide-react";

interface EngineerReassignModalProps {
  isOpen: boolean;
  onClose: () => void;
  currentEngineer: string | null;
  engineers: any[];
  onReassign: (engineerId: string, engineerName: string) => Promise<void>;
  isLoading: boolean;
}

export function EngineerReassignModal({
  isOpen,
  onClose,
  currentEngineer,
  engineers,
  onReassign,
  isLoading,
}: EngineerReassignModalProps) {
  const [selectedEngineer, setSelectedEngineer] = useState("");
  const [submitting, setSubmitting] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async () => {
    if (!selectedEngineer) return;

    const engineer = engineers.find((e) => e.id === selectedEngineer);
    if (!engineer) return;

    setSubmitting(true);
    try {
      await onReassign(engineer.id, engineer.name);
      onClose();
    } catch (error) {
      console.error("Failed to reassign:", error);
      alert("Failed to reassign engineer. Please try again.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
        {/* Header */}
        <div className="border-b px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-blue-100 rounded-lg">
              <User className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <h2 className="text-lg font-semibold text-gray-900">
                {currentEngineer ? "Reassign Engineer" : "Assign Engineer"}
              </h2>
              <p className="text-sm text-gray-500">Select a new engineer for this ticket</p>
            </div>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
            <X className="h-5 w-5 text-gray-500" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          {currentEngineer && (
            <div className="mb-4 p-3 bg-gray-50 rounded-lg">
              <p className="text-sm text-gray-600">Currently assigned to:</p>
              <p className="font-medium text-gray-900">{currentEngineer}</p>
            </div>
          )}

          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
            </div>
          ) : engineers && engineers.length > 0 ? (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Select Engineer
              </label>
              <select
                value={selectedEngineer}
                onChange={(e) => setSelectedEngineer(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={submitting}
              >
                <option value="">-- Select an engineer --</option>
                {engineers.map((engineer) => (
                  <option key={engineer.id} value={engineer.id}>
                    {engineer.name} {engineer.specialization ? `(${engineer.specialization})` : ""}
                  </option>
                ))}
              </select>

              {selectedEngineer && (
                <div className="mt-4 p-3 bg-blue-50 rounded-lg">
                  <p className="text-sm text-blue-800 flex items-center gap-2">
                    <UserCheck className="h-4 w-4" />
                    Engineer will be notified upon assignment
                  </p>
                </div>
              )}
            </div>
          ) : (
            <div className="text-center py-8">
              <User className="h-12 w-12 text-gray-300 mx-auto mb-3" />
              <p className="text-gray-600">No engineers available</p>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="border-t px-6 py-4 bg-gray-50 flex justify-end gap-3">
          <button
            onClick={onClose}
            disabled={submitting}
            className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-100 transition-colors disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={!selectedEngineer || submitting}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 flex items-center gap-2"
          >
            {submitting ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Assigning...
              </>
            ) : (
              <>
                <UserCheck className="h-4 w-4" />
                {currentEngineer ? "Reassign" : "Assign"}
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
