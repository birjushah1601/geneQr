"use client";

import { useState } from "react";
import { Brain, X, Loader2, AlertTriangle, CheckCircle, Lightbulb, Package, TrendingUp } from "lucide-react";

interface AIDiagnosisModalProps {
  isOpen: boolean;
  onClose: () => void;
  diagnosis: any;
  isLoading: boolean;
  onTriggerAnalysis: () => void;
}

export function AIDiagnosisModal({ isOpen, onClose, diagnosis, isLoading, onTriggerAnalysis }: AIDiagnosisModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-purple-100 rounded-lg">
              <Brain className="h-6 w-6 text-purple-600" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-gray-900">AI Diagnosis & Analysis</h2>
              <p className="text-sm text-gray-500">Comprehensive analysis of all ticket data</p>
            </div>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
            <X className="h-5 w-5 text-gray-500" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          {isLoading ? (
            <div className="flex flex-col items-center justify-center py-12">
              <Loader2 className="h-12 w-12 text-purple-600 animate-spin mb-4" />
              <p className="text-gray-600">Analyzing ticket data with AI...</p>
              <p className="text-sm text-gray-500 mt-2">This may take a few moments</p>
            </div>
          ) : diagnosis && diagnosis.length > 0 ? (
            <>
              {/* Latest Diagnosis */}
              {diagnosis[0] && (() => {
                const latest = diagnosis[0];
                const confidenceColor =
                  latest.confidence_level === 'HIGH' ? 'bg-green-100 text-green-800 border-green-300' :
                  latest.confidence_level === 'MEDIUM' ? 'bg-yellow-100 text-yellow-800 border-yellow-300' :
                  'bg-red-100 text-red-800 border-red-300';

                return (
                  <div className="space-y-6">
                    {/* Primary Diagnosis */}
                    <div className={`border-2 rounded-lg p-6 ${confidenceColor}`}>
                      <div className="flex items-start justify-between mb-4">
                        <div>
                          <h3 className="text-lg font-semibold mb-1">{latest.primary_diagnosis.problem_type}</h3>
                          <p className="text-sm opacity-80">{latest.primary_diagnosis.problem_category}</p>
                        </div>
                        <div className="flex gap-2">
                          <span className="px-3 py-1 bg-white bg-opacity-50 rounded-full text-xs font-medium">
                            {latest.confidence_level} Confidence ({Math.round(latest.confidence * 100)}%)
                          </span>
                          <span className="px-3 py-1 bg-white bg-opacity-50 rounded-full text-xs font-medium">
                            {latest.primary_diagnosis.severity} Severity
                          </span>
                        </div>
                      </div>
                      <p className="text-sm leading-relaxed">{latest.primary_diagnosis.description}</p>
                      {latest.primary_diagnosis.root_cause && (
                        <div className="mt-4 pt-4 border-t border-current border-opacity-20">
                          <p className="text-sm font-medium mb-1">Root Cause:</p>
                          <p className="text-sm">{latest.primary_diagnosis.root_cause}</p>
                        </div>
                      )}
                    </div>

                    {/* Symptoms */}
                    {latest.primary_diagnosis.symptoms && latest.primary_diagnosis.symptoms.length > 0 && (
                      <div className="bg-gray-50 rounded-lg p-4">
                        <h4 className="font-medium text-gray-900 mb-3 flex items-center gap-2">
                          <AlertTriangle className="h-4 w-4 text-orange-600" />
                          Identified Symptoms
                        </h4>
                        <ul className="space-y-2">
                          {latest.primary_diagnosis.symptoms.map((symptom: string, idx: number) => (
                            <li key={idx} className="flex items-start gap-2 text-sm text-gray-700">
                              <span className="text-orange-600 mt-0.5">•</span>
                              <span>{symptom}</span>
                            </li>
                          ))}
                        </ul>
                      </div>
                    )}

                    {/* Recommended Actions */}
                    {latest.recommended_actions && latest.recommended_actions.length > 0 && (
                      <div className="bg-blue-50 rounded-lg p-4">
                        <h4 className="font-medium text-gray-900 mb-3 flex items-center gap-2">
                          <Lightbulb className="h-4 w-4 text-blue-600" />
                          Recommended Actions
                        </h4>
                        <div className="space-y-3">
                          {latest.recommended_actions.map((action: any, idx: number) => (
                            <div key={idx} className="bg-white rounded p-3">
                              <div className="flex items-start justify-between mb-2">
                                <h5 className="font-medium text-gray-900 text-sm">{action.action}</h5>
                                <div className="flex gap-2">
                                  <span className="px-2 py-0.5 bg-blue-100 text-blue-700 rounded text-xs">
                                    {action.priority}
                                  </span>
                                  <span className="px-2 py-0.5 bg-gray-100 text-gray-700 rounded text-xs">
                                    {action.estimated_time}
                                  </span>
                                </div>
                              </div>
                              <p className="text-xs text-gray-600">{action.description}</p>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Required Parts */}
                    {latest.required_parts && latest.required_parts.length > 0 && (
                      <div className="bg-green-50 rounded-lg p-4">
                        <h4 className="font-medium text-gray-900 mb-3 flex items-center gap-2">
                          <Package className="h-4 w-4 text-green-600" />
                          Required Parts
                        </h4>
                        <div className="space-y-2">
                          {latest.required_parts.map((part: any, idx: number) => (
                            <div key={idx} className="bg-white rounded p-3 flex items-center justify-between">
                              <div>
                                <p className="font-medium text-sm text-gray-900">{part.part_name}</p>
                                <p className="text-xs text-gray-500">{part.part_number}</p>
                              </div>
                              <div className="text-right">
                                <p className="text-sm font-medium text-gray-900">Qty: {part.quantity}</p>
                                <p className="text-xs text-gray-500">{part.urgency}</p>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Vision Analysis */}
                    {latest.vision_analysis && latest.vision_analysis.findings && latest.vision_analysis.findings.length > 0 && (
                      <div className="bg-purple-50 rounded-lg p-4">
                        <h4 className="font-medium text-gray-900 mb-3 flex items-center gap-2">
                          <TrendingUp className="h-4 w-4 text-purple-600" />
                          Visual Analysis Findings
                        </h4>
                        <p className="text-sm text-gray-700 mb-3">{latest.vision_analysis.overall_assessment}</p>
                        <div className="space-y-2">
                          {latest.vision_analysis.findings.map((finding: any, idx: number) => (
                            <div key={idx} className="bg-white rounded p-3">
                              <p className="text-sm font-medium text-gray-900 mb-1">{finding.observation}</p>
                              <p className="text-xs text-gray-600">{finding.significance}</p>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Metadata */}
                    <div className="flex items-center justify-between text-xs text-gray-500 pt-4 border-t">
                      <span>AI Model: {latest.ai_metadata.model}</span>
                      <span>{new Date(latest.created_at).toLocaleString()}</span>
                      {latest.ai_metadata.alternatives_count > 0 && (
                        <span>+{latest.ai_metadata.alternatives_count} alternative diagnoses</span>
                      )}
                    </div>
                  </div>
                );
              })()}
            </>
          ) : (
            <div className="flex flex-col items-center justify-center py-12">
              <Brain className="h-16 w-16 text-gray-300 mb-4" />
              <p className="text-gray-600 mb-4">No AI diagnosis available yet</p>
              <button
                onClick={onTriggerAnalysis}
                className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors flex items-center gap-2"
              >
                <Brain className="h-4 w-4" />
                Trigger AI Analysis
              </button>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="border-t px-6 py-4 bg-gray-50 flex justify-end gap-3">
          {diagnosis && diagnosis.length > 0 && (
            <button
              onClick={onTriggerAnalysis}
              disabled={isLoading}
              className="px-4 py-2 border border-purple-600 text-purple-600 rounded-lg hover:bg-purple-50 transition-colors disabled:opacity-50"
            >
              Re-analyze
            </button>
          )}
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
