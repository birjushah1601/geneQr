import React, { useState } from 'react';
import { Brain, AlertCircle, CheckCircle, X, MessageCircle, Eye, Wrench, Package } from 'lucide-react';
import { DiagnosisResponse, getConfidenceBadgeColor, getConfidenceDescription } from '@/lib/api/diagnosis';

interface DiagnosisCardProps {
  diagnosis: DiagnosisResponse;
  onAccept: (diagnosisId: string) => void;
  onReject: (diagnosisId: string, feedback?: string) => void;
  loading?: boolean;
}

export function DiagnosisCard({ diagnosis, onAccept, onReject, loading = false }: DiagnosisCardProps) {
  const [showFeedbackModal, setShowFeedbackModal] = useState(false);
  const [feedbackText, setFeedbackText] = useState('');

  const handleAccept = () => {
    onAccept(diagnosis.diagnosis_id);
  };

  const handleRejectWithFeedback = () => {
    onReject(diagnosis.diagnosis_id, feedbackText);
    setShowFeedbackModal(false);
    setFeedbackText('');
  };

  const confidenceColor = getConfidenceBadgeColor(diagnosis.confidence_level);
  const confidenceDesc = getConfidenceDescription(diagnosis.confidence_level);

  return (
    <div className="bg-white rounded-lg shadow-lg p-6 border border-gray-200">
      {/* Header */}
      <div className="flex items-start justify-between mb-6">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-blue-100 rounded-lg">
            <Brain className="h-6 w-6 text-blue-600" />
          </div>
          <div>
            <h2 className="text-xl font-semibold text-gray-900">AI Diagnosis</h2>
            <p className="text-sm text-gray-500">ID: {diagnosis.diagnosis_id}</p>
          </div>
        </div>
        
        {/* Confidence Badge */}
        <div className="flex flex-col items-end gap-2">
          <div className={`px-3 py-1 rounded-full text-xs font-medium border ${confidenceColor}`}>
            {Math.round(diagnosis.confidence * 100)}% {diagnosis.confidence_level}
          </div>
          <p className="text-xs text-gray-500 text-right max-w-xs">
            {confidenceDesc}
          </p>
        </div>
      </div>

      {/* Primary Diagnosis */}
      <div className="mb-6 p-4 bg-gray-50 rounded-lg">
        <h3 className="text-lg font-medium text-gray-900 mb-2">Primary Diagnosis</h3>
        <div className="grid grid-cols-2 gap-4 mb-3">
          <div>
            <p className="text-sm text-gray-600 font-medium">Issue</p>
            <p className="text-gray-900">{diagnosis.primary_diagnosis.description}</p>
          </div>
          <div>
            <p className="text-sm text-gray-600 font-medium">Category</p>
            <p className="text-gray-900">{diagnosis.primary_diagnosis.problem_category}</p>
          </div>
        </div>
        <div className="mb-3">
          <p className="text-sm text-gray-600 font-medium">Root Cause</p>
          <p className="text-gray-900">{diagnosis.primary_diagnosis.root_cause}</p>
        </div>
        <div>
          <p className="text-sm text-gray-600 font-medium">Reasoning</p>
          <p className="text-gray-900 text-sm">{diagnosis.primary_diagnosis.reasoning_explanation}</p>
        </div>
      </div>

      {/* Confidence Factors */}
      <div className="mb-6">
        <h4 className="text-sm font-medium text-gray-700 mb-2 flex items-center gap-2">
          <Eye className="h-4 w-4" />
          Confidence Based On
        </h4>
        <ul className="space-y-1">
          {diagnosis.ai_metadata.confidence_factors.map((factor, index) => (
            <li key={index} className="text-sm text-gray-600 flex items-center gap-2">
              <div className="w-1.5 h-1.5 bg-blue-500 rounded-full"></div>
              {factor}
            </li>
          ))}
        </ul>
      </div>

      {/* Alternative Diagnoses */}
      {diagnosis.alternate_diagnoses.length > 0 && (
        <div className="mb-6">
          <h4 className="text-sm font-medium text-gray-700 mb-3">Alternative Possibilities</h4>
          <div className="space-y-2">
            {diagnosis.alternate_diagnoses.map((alt, index) => (
              <div key={index} className="p-3 bg-yellow-50 rounded-lg border border-yellow-200">
                <div className="flex justify-between items-start mb-1">
                  <p className="font-medium text-gray-900 text-sm">{alt.description}</p>
                  <span className="text-xs text-yellow-800 bg-yellow-100 px-2 py-1 rounded">
                    {Math.round(alt.confidence)}%
                  </span>
                </div>
                <p className="text-xs text-gray-600">{alt.reasoning_explanation}</p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Recommended Actions */}
      {diagnosis.recommended_actions.length > 0 && (
        <div className="mb-6">
          <h4 className="text-sm font-medium text-gray-700 mb-3 flex items-center gap-2">
            <Wrench className="h-4 w-4" />
            Recommended Actions
          </h4>
          <div className="space-y-2">
            {diagnosis.recommended_actions.map((action, index) => (
              <div key={index} className="p-3 border border-gray-200 rounded-lg">
                <div className="flex justify-between items-start mb-1">
                  <p className="font-medium text-gray-900 text-sm">{action.action}</p>
                  <span className={`text-xs px-2 py-1 rounded ${
                    action.priority === 'high' ? 'bg-red-100 text-red-800' :
                    action.priority === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                    'bg-green-100 text-green-800'
                  }`}>
                    {action.priority}
                  </span>
                </div>
                <p className="text-xs text-gray-600 mb-2">{action.description}</p>
                <div className="text-xs text-gray-500">
                  Est. time: {action.estimated_time}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Required Parts */}
      {diagnosis.required_parts.length > 0 && (
        <div className="mb-6">
          <h4 className="text-sm font-medium text-gray-700 mb-3 flex items-center gap-2">
            <Package className="h-4 w-4" />
            Suggested Parts
          </h4>
          <div className="space-y-2">
            {diagnosis.required_parts.map((part, index) => (
              <div key={index} className="p-3 border border-gray-200 rounded-lg">
                <div className="flex justify-between items-start mb-1">
                  <p className="font-medium text-gray-900 text-sm">{part.part_name}</p>
                  <div className="text-right">
                    <div className="text-xs text-gray-500">
                      {Math.round(part.probability)}% likely needed
                    </div>
                    <div className="text-xs text-gray-600">Qty: {part.quantity}</div>
                  </div>
                </div>
                <p className="text-xs text-gray-600">Part #: {part.part_code}</p>
                <p className="text-xs text-gray-500">{part.manufacturer}</p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Decision Status */}
      {diagnosis.decision_status === 'pending' && (
        <div className="border-t border-gray-200 pt-4">
          <p className="text-sm text-gray-600 mb-4 flex items-center gap-2">
            <AlertCircle className="h-4 w-4 text-orange-500" />
            This is an AI suggestion. Please review and decide.
          </p>
          
          <div className="flex gap-3">
            <button
              onClick={handleAccept}
              disabled={loading}
              className="flex-1 bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <CheckCircle className="h-4 w-4" />
              Accept Diagnosis
            </button>
            
            <button
              onClick={() => setShowFeedbackModal(true)}
              disabled={loading}
              className="flex-1 bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <X className="h-4 w-4" />
              Reject & Provide Feedback
            </button>
          </div>
        </div>
      )}

      {/* Decision Made */}
      {diagnosis.decision_status !== 'pending' && (
        <div className="border-t border-gray-200 pt-4">
          <div className={`p-3 rounded-lg flex items-center gap-2 ${
            diagnosis.decision_status === 'accepted' 
              ? 'bg-green-50 text-green-800' 
              : 'bg-red-50 text-red-800'
          }`}>
            {diagnosis.decision_status === 'accepted' ? (
              <CheckCircle className="h-4 w-4" />
            ) : (
              <X className="h-4 w-4" />
            )}
            <span className="font-medium">
              {diagnosis.decision_status === 'accepted' ? 'Accepted' : 'Rejected'}
            </span>
            {diagnosis.decided_at && (
              <span className="text-sm">
                on {new Date(diagnosis.decided_at).toLocaleDateString()}
              </span>
            )}
          </div>
          {diagnosis.feedback_text && (
            <div className="mt-2 p-2 bg-gray-50 rounded text-sm text-gray-700">
              <strong>Feedback:</strong> {diagnosis.feedback_text}
            </div>
          )}
        </div>
      )}

      {/* Feedback Modal */}
      {showFeedbackModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full">
            <div className="flex items-center</p></div><div className="gap-2 mb-4">
              <MessageCircle className="h-5 w-5 text-blue-600" />
              <h3 className="text-lg font-semibold">Why are you rejecting this diagnosis?</h3>
            </div>
            
            <textarea
              value={feedbackText}
              onChange={(e) => setFeedbackText(e.target.value)}
              placeholder="Please explain what the AI got wrong and what the correct diagnosis should be..."
              className="w-full h-32 p-3 border border-gray-300 rounded-md resize-none focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            
            <div className="flex gap-3 mt-4">
              <button
                onClick={() => setShowFeedbackModal(false)}
                className="flex-1 bg-gray-300 text-gray-700 py-2 px-4 rounded-md hover:bg-gray-400 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleRejectWithFeedback}
                disabled={!feedbackText.trim()}
                className="flex-1 bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700 transition-colors disabled:opacity-50"
              >
                Submit Feedback
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}