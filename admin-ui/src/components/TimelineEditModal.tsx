"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { X, Check, Clock, AlertTriangle, Edit2, Save } from "lucide-react";
import type { PublicMilestone, PublicTimeline } from "@/types";

interface TimelineEditModalProps {
  timeline: PublicTimeline;
  ticketId: string;
  onClose: () => void;
  onSave: (updatedTimeline: any) => Promise<void>;
}

export function TimelineEditModal({ timeline, ticketId, onClose, onSave }: TimelineEditModalProps) {
  const [editedTimeline, setEditedTimeline] = useState(timeline);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const handleMilestoneStatusChange = (index: number, newStatus: string) => {
    const updated = { ...editedTimeline };
    updated.milestones[index] = {
      ...updated.milestones[index],
      status: newStatus,
    };
    setEditedTimeline(updated);
  };

  const handleMilestoneEtaChange = (index: number, newEta: string) => {
    const updated = { ...editedTimeline };
    updated.milestones[index] = {
      ...updated.milestones[index],
      eta: newEta,
    };
    setEditedTimeline(updated);
  };

  const handleResolutionDateChange = (newDate: string) => {
    setEditedTimeline({
      ...editedTimeline,
      estimated_resolution: newDate,
    });
  };

  const handlePartsStatusChange = (newStatus: string) => {
    setEditedTimeline({
      ...editedTimeline,
      parts_status: newStatus as any,
    });
  };

  const handlePartsEtaChange = (newDate: string) => {
    setEditedTimeline({
      ...editedTimeline,
      parts_eta: newDate,
    });
  };

  const handleSave = async () => {
    try {
      setSaving(true);
      setError("");
      await onSave(editedTimeline);
      onClose();
    } catch (err: any) {
      setError(err.message || "Failed to save timeline");
    } finally {
      setSaving(false);
    }
  };

  const formatDateForInput = (dateStr?: string) => {
    if (!dateStr) return "";
    const date = new Date(dateStr);
    return date.toISOString().slice(0, 16); // Format: YYYY-MM-DDTHH:mm
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 overflow-y-auto">
      <Card className="w-full max-w-4xl max-h-[90vh] overflow-y-auto">
        <CardHeader className="border-b sticky top-0 bg-white z-10">
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <Edit2 className="h-5 w-5 text-blue-600" />
              Edit Service Timeline
            </CardTitle>
            <button
              onClick={onClose}
              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <X className="h-5 w-5" />
            </button>
          </div>
          <p className="text-sm text-gray-600 mt-2">
            Adjust milestones, ETAs, and parts status to reflect the actual service progress.
          </p>
        </CardHeader>

        <CardContent className="space-y-6 p-6">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg flex items-center gap-2">
              <AlertTriangle className="h-5 w-5" />
              {error}
            </div>
          )}

          {/* Overall Resolution Date */}
          <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
            <h3 className="font-semibold text-blue-900 mb-3 flex items-center gap-2">
              <Clock className="h-4 w-4" />
              Overall Expected Resolution
            </h3>
            <div className="flex items-center gap-4">
              <label className="text-sm font-medium text-gray-700 min-w-fit">
                Target Completion:
              </label>
              <Input
                type="datetime-local"
                value={formatDateForInput(editedTimeline.estimated_resolution)}
                onChange={(e) => handleResolutionDateChange(e.target.value)}
                className="flex-1"
              />
            </div>
          </div>

          {/* Parts Information */}
          {editedTimeline.requires_parts && (
            <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
              <h3 className="font-semibold text-purple-900 mb-3">Parts Information</h3>
              <div className="space-y-3">
                <div className="flex items-center gap-4">
                  <label className="text-sm font-medium text-gray-700 min-w-fit">
                    Parts Status:
                  </label>
                  <select
                    value={editedTimeline.parts_status || "ordering"}
                    onChange={(e) => handlePartsStatusChange(e.target.value)}
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500"
                  >
                    <option value="ordering">Ordering</option>
                    <option value="in_transit">In Transit</option>
                    <option value="received">Received</option>
                  </select>
                </div>
                <div className="flex items-center gap-4">
                  <label className="text-sm font-medium text-gray-700 min-w-fit">
                    Expected Arrival:
                  </label>
                  <Input
                    type="datetime-local"
                    value={formatDateForInput(editedTimeline.parts_eta)}
                    onChange={(e) => handlePartsEtaChange(e.target.value)}
                    className="flex-1"
                  />
                </div>
              </div>
            </div>
          )}

          {/* Milestone List */}
          <div>
            <h3 className="font-semibold text-gray-900 mb-4">Milestones</h3>
            <div className="space-y-4">
              {editedTimeline.milestones.map((milestone, index) => (
                <Card key={index} className="border-l-4 border-l-blue-500">
                  <CardContent className="p-4">
                    <div className="flex items-start gap-4">
                      <div className="flex-1">
                        <h4 className="font-semibold text-gray-900">{milestone.title}</h4>
                        <p className="text-sm text-gray-600 mt-1">{milestone.description}</p>
                      </div>
                      <div className="flex flex-col gap-2 min-w-[250px]">
                        <div>
                          <label className="text-xs text-gray-600 block mb-1">Status</label>
                          <select
                            value={milestone.status}
                            onChange={(e) => handleMilestoneStatusChange(index, e.target.value)}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                          >
                            <option value="pending">Pending</option>
                            <option value="in_progress">In Progress</option>
                            <option value="completed">Completed</option>
                            <option value="delayed">Delayed</option>
                            <option value="blocked">Blocked</option>
                            <option value="skipped">Skipped</option>
                          </select>
                        </div>
                        {milestone.status !== 'completed' && milestone.status !== 'skipped' && (
                          <div>
                            <label className="text-xs text-gray-600 block mb-1">ETA</label>
                            <Input
                              type="datetime-local"
                              value={formatDateForInput(milestone.eta)}
                              onChange={(e) => handleMilestoneEtaChange(index, e.target.value)}
                              className="text-sm"
                            />
                          </div>
                        )}
                      </div>
                    </div>

                    {/* Quick Actions */}
                    <div className="flex gap-2 mt-3">
                      {milestone.status !== 'completed' && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleMilestoneStatusChange(index, 'completed')}
                          className="text-xs"
                        >
                          <Check className="h-3 w-3 mr-1" />
                          Mark Complete
                        </Button>
                      )}
                      {milestone.status !== 'skipped' && milestone.status !== 'completed' && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleMilestoneStatusChange(index, 'skipped')}
                          className="text-xs"
                        >
                          Skip
                        </Button>
                      )}
                      {milestone.status !== 'blocked' && milestone.status !== 'completed' && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleMilestoneStatusChange(index, 'blocked')}
                          className="text-xs text-orange-600 border-orange-300"
                        >
                          <AlertTriangle className="h-3 w-3 mr-1" />
                          Mark Blocked
                        </Button>
                      )}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>

          {/* Save Actions */}
          <div className="flex gap-3 justify-end pt-4 border-t sticky bottom-0 bg-white">
            <Button
              variant="outline"
              onClick={onClose}
              disabled={saving}
            >
              Cancel
            </Button>
            <Button
              onClick={handleSave}
              disabled={saving}
              className="bg-blue-600 hover:bg-blue-700"
            >
              {saving ? (
                <>
                  <Clock className="h-4 w-4 mr-2 animate-spin" />
                  Saving...
                </>
              ) : (
                <>
                  <Save className="h-4 w-4 mr-2" />
                  Save Timeline
                </>
              )}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
