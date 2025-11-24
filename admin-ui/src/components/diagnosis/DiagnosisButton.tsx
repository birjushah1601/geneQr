import React, { useState } from 'react';
import { Brain, Loader2, Sparkles } from 'lucide-react';
import { DiagnosisRequest, DiagnosisResponse, diagnosisApi, extractSymptoms } from '@/lib/api/diagnosis';
import { aiDiagnosisEngine } from '@/lib/ai/openai';
import { Equipment } from '@/types';

interface DiagnosisButtonProps {
  equipment: Equipment;
  description: string;
  priority: string;
  onDiagnosisComplete?: (diagnosis: DiagnosisResponse) => void;
  className?: string;
}

export function DiagnosisButton({ 
  equipment, 
  description, 
  priority,
  onDiagnosisComplete,
  className = ''
}: DiagnosisButtonProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const runDiagnosis = async () => {
    if (!description.trim()) {
      setError('Please provide an issue description first');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      // Extract symptoms from description
      const symptoms = extractSymptoms(description);

      // Create diagnosis request
      const request: DiagnosisRequest = {
        ticket_id: Math.floor(Math.random() * 100000), // Would be real ticket ID in production
        equipment_id: equipment.id,
        symptoms,
        description,
        user_context: {
          user_id: 1, // Would be real user ID
          username: 'admin', // Would be real username
          role: 'technician' // Would be real role
        },
        options: {
          include_vision_analysis: true,
          include_historical_context: true,
          include_similar_tickets: true,
          max_similar_tickets: 5,
          min_confidence_threshold: 0.3
        }
      };

      // Use AI diagnosis engine for intelligent analysis
      console.log('🧠 Running AI diagnosis for:', equipment.equipment_name);
      const diagnosis = await aiDiagnosisEngine.analyzeDiagnosis(request, equipment);
      
      if (onDiagnosisComplete) {
        onDiagnosisComplete(diagnosis);
      }

    } catch (err) {
      console.error('Diagnosis failed:', err);
      setError(err instanceof Error ? err.message : 'Failed to run AI diagnosis');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <button
        onClick={runDiagnosis}
        disabled={loading || !description.trim()}
        className={`
          flex items-center justify-center gap-2 px-4 py-2 rounded-md 
          transition-colors border font-medium text-sm
          ${!description.trim() 
            ? 'bg-gray-100 text-gray-400 border-gray-200 cursor-not-allowed'
            : loading
            ? 'bg-blue-50 text-blue-600 border-blue-200'
            : 'bg-gradient-to-r from-purple-50 to-blue-50 text-purple-700 border-purple-200 hover:from-purple-100 hover:to-blue-100 hover:border-purple-300'
          }
          ${className}
        `}
      >
        {loading ? (
          <>
            <Loader2 className="h-4 w-4 animate-spin" />
            Running AI Diagnosis...
          </>
        ) : (
          <>
            <Brain className="h-4 w-4" />
            <Sparkles className="h-3 w-3" />
            Get AI Diagnosis
          </>
        )}
      </button>
      
      {error && (
        <p className="text-red-600 text-xs mt-1">{error}</p>
      )}
      
      {!description.trim() && (
        <p className="text-gray-500 text-xs mt-1">
          Add issue description to enable AI diagnosis
        </p>
      )}
    </div>
  );
}