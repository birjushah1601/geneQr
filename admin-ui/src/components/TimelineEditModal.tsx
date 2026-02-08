"use client";

import React, { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { X, Check, Clock, AlertTriangle, Edit2, Save, ChevronUp, ChevronDown, GripVertical } from "lucide-react";
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
  const [blockerComments, setBlockerComments] = useState<{[key: number]: string}>({});
  const [customMilestones, setCustomMilestones] = useState<PublicMilestone[]>([]);
  const [draggedIndex, setDraggedIndex] = useState<number | null>(null);

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

  const handleAddMilestone = () => {
    const newMilestone: PublicMilestone = {
      stage: `custom_${Date.now()}`,
      title: "New Custom Stage",
      description: "Add description...",
      status: "pending",
      is_active: false,
    };
    
    const updated = { ...editedTimeline };
    // Insert before the last milestone (Resolution)
    updated.milestones.splice(updated.milestones.length - 1, 0, newMilestone);
    setEditedTimeline(updated);
  };

  const handleRemoveMilestone = (index: number) => {
    const updated = { ...editedTimeline };
    updated.milestones.splice(index, 1);
    setEditedTimeline(updated);
  };

  const handleMilestoneTitleChange = (index: number, newTitle: string) => {
    const updated = { ...editedTimeline };
    updated.milestones[index] = {
      ...updated.milestones[index],
      title: newTitle,
    };
    setEditedTimeline(updated);
  };

  const handleMilestoneDescriptionChange = (index: number, newDesc: string) => {
    const updated = { ...editedTimeline };
    updated.milestones[index] = {
      ...updated.milestones[index],
      description: newDesc,
    };
    setEditedTimeline(updated);
  };

  const handleBlockerCommentChange = (index: number, comment: string) => {
    setBlockerComments({
      ...blockerComments,
      [index]: comment,
    });
  };

  const handleMoveMilestoneUp = (index: number) => {
    if (index === 0) return; // Can't move first item up
    
    const updated = { ...editedTimeline };
    const milestones = [...updated.milestones];
    // Swap with previous
    [milestones[index - 1], milestones[index]] = [milestones[index], milestones[index - 1]];
    updated.milestones = milestones;
    setEditedTimeline(updated);
  };

  const handleMoveMilestoneDown = (index: number) => {
    if (index >= editedTimeline.milestones.length - 1) return; // Can't move last item down
    
    const updated = { ...editedTimeline };
    const milestones = [...updated.milestones];
    // Swap with next
    [milestones[index], milestones[index + 1]] = [milestones[index + 1], milestones[index]];
    updated.milestones = milestones;
    setEditedTimeline(updated);
  };

  const handleDragStart = (e: React.DragEvent<HTMLDivElement>, index: number) => {
    setDraggedIndex(index);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>, index: number) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>, targetIndex: number) => {
    e.preventDefault();
    
    if (draggedIndex === null || draggedIndex === targetIndex) return;
    
    const updated = { ...editedTimeline };
    const milestones = [...updated.milestones];
    const [removed] = milestones.splice(draggedIndex, 1);
    milestones.splice(targetIndex, 0, removed);
    updated.milestones = milestones;
    setEditedTimeline(updated);
    setDraggedIndex(null);
  };

  const handleDragEnd = () => {
    setDraggedIndex(null);
  };

  const handleSave = async () => {
    try {
      setSaving(true);
      setError("");
      
      // Save timeline
      await onSave({
        ...editedTimeline,
        blocker_comments: blockerComments,
      });
      
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
    
    // Format in local timezone for datetime-local input
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="w-full max-w-4xl max-h-[90vh] overflow-y-auto overflow-x-visible" style={{ isolation: 'isolate' }}>
        <Card className="w-full">
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
              <div className="flex-1 relative" style={{
                backgroundColor: '#F3F4F6',
                border: '2px solid #9CA3AF',
                borderRadius: '8px',
                padding: '8px',
                boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
              }}>
                <Input
                  type="datetime-local"
                  value={formatDateForInput(editedTimeline.estimated_resolution)}
                  onChange={(e) => handleResolutionDateChange(e.target.value)}
                  className="w-full border-0 outline-none"
                  style={{ 
                    colorScheme: 'light',
                    backgroundColor: '#F3F4F6',
                    fontSize: '14px',
                    fontWeight: '500'
                  }}
                />
              </div>
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
                  <div className="flex-1 relative" style={{
                    backgroundColor: '#F3F4F6',
                    border: '2px solid #9CA3AF',
                    borderRadius: '8px',
                    padding: '8px',
                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                  }}>
                    <Input
                      type="datetime-local"
                      value={formatDateForInput(editedTimeline.parts_eta)}
                      onChange={(e) => handlePartsEtaChange(e.target.value)}
                      className="w-full border-0 outline-none"
                      style={{ 
                        colorScheme: 'light',
                        backgroundColor: '#F3F4F6',
                        fontSize: '14px',
                        fontWeight: '500'
                      }}
                    />
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Milestone List */}
          <div>
            <div className="flex items-center justify-between mb-4">
              <h3 className="font-semibold text-gray-900">Milestones</h3>
              <Button
                size="sm"
                variant="outline"
                onClick={handleAddMilestone}
                className="text-xs"
              >
                + Add Custom Milestone
              </Button>
            </div>
            <div className="space-y-4">
              {editedTimeline.milestones.map((milestone, index) => (
                <Card 
                  key={index} 
                  className={`border-l-4 border-l-blue-500 transition-all ${
                    draggedIndex === index ? 'opacity-50 scale-95' : ''
                  }`}
                  draggable
                  onDragStart={(e) => handleDragStart(e, index)}
                  onDragOver={(e) => handleDragOver(e, index)}
                  onDrop={(e) => handleDrop(e, index)}
                  onDragEnd={handleDragEnd}
                >
                  <CardContent className="p-4">
                    <div className="flex items-start gap-2">
                      {/* Drag Handle & Move Buttons */}
                      <div className="flex flex-col items-center gap-1 pt-1">
                        <GripVertical className="h-5 w-5 text-gray-400 cursor-move" title="Drag to reorder" />
                        <div className="flex flex-col gap-0.5">
                          <button
                            onClick={() => handleMoveMilestoneUp(index)}
                            disabled={index === 0}
                            className="p-1 hover:bg-gray-100 rounded disabled:opacity-30 disabled:cursor-not-allowed"
                            title="Move up"
                          >
                            <ChevronUp className="h-4 w-4 text-gray-600" />
                          </button>
                          <button
                            onClick={() => handleMoveMilestoneDown(index)}
                            disabled={index === editedTimeline.milestones.length - 1}
                            className="p-1 hover:bg-gray-100 rounded disabled:opacity-30 disabled:cursor-not-allowed"
                            title="Move down"
                          >
                            <ChevronDown className="h-4 w-4 text-gray-600" />
                          </button>
                        </div>
                      </div>
                      
                      <div className="flex items-start gap-4 flex-1">
                      <div className="flex-1">
                        {milestone.stage.startsWith('custom_') ? (
                          <Input
                            value={milestone.title}
                            onChange={(e) => handleMilestoneTitleChange(index, e.target.value)}
                            placeholder="Milestone title"
                            className="font-semibold mb-2"
                          />
                        ) : (
                          <h4 className="font-semibold text-gray-900">{milestone.title}</h4>
                        )}
                        {milestone.stage.startsWith('custom_') ? (
                          <Input
                            value={milestone.description}
                            onChange={(e) => handleMilestoneDescriptionChange(index, e.target.value)}
                            placeholder="Milestone description"
                            className="text-sm"
                          />
                        ) : (
                          <p className="text-sm text-gray-600 mt-1">{milestone.description}</p>
                        )}
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
                            <div className="relative" style={{
                              backgroundColor: '#F3F4F6',
                              border: '2px solid #9CA3AF',
                              borderRadius: '8px',
                              padding: '6px',
                              boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                            }}>
                              <Input
                                type="datetime-local"
                                value={formatDateForInput(milestone.eta)}
                                onChange={(e) => handleMilestoneEtaChange(index, e.target.value)}
                                className="w-full border-0 outline-none text-sm"
                                style={{ 
                                  colorScheme: 'light',
                                  backgroundColor: '#F3F4F6',
                                  fontSize: '13px',
                                  fontWeight: '500'
                                }}
                              />
                            </div>
                          </div>
                        )}
                      </div>
                      </div>
                    </div>

                    {/* Blocker Comment (when marked blocked or delayed) */}
                    {(milestone.status === 'blocked' || milestone.status === 'delayed') && (
                      <div className="mt-3 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                        <label className="text-xs font-semibold text-yellow-900 block mb-2">
                          📝 Explain to Customer (will be added as comment):
                        </label>
                        <textarea
                          value={blockerComments[index] || ''}
                          onChange={(e) => handleBlockerCommentChange(index, e.target.value)}
                          placeholder="e.g., 'Waiting for specialized part delivery from Germany. Expected arrival: Feb 20'"
                          className="w-full px-3 py-2 text-sm border border-yellow-300 rounded-lg focus:ring-2 focus:ring-yellow-500 resize-none"
                          rows={3}
                        />
                        <p className="text-xs text-gray-600 mt-1">
                          This will be posted as a comment on the ticket for customer transparency.
                        </p>
                      </div>
                    )}

                    {/* Quick Actions */}
                    <div className="flex gap-2 mt-3 flex-wrap">
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
                      {milestone.stage.startsWith('custom_') && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleRemoveMilestone(index)}
                          className="text-xs text-red-600 border-red-300"
                        >
                          <X className="h-3 w-3 mr-1" />
                          Remove
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
    </div>
  );
}
