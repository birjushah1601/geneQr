// ============================================================================
// Optional AI Analysis Button
// ============================================================================

import { useState } from 'react';
import { Sparkles, Loader2, CheckCircle, AlertCircle } from 'lucide-react';
import { runAIDiagnosis } from '@/lib/utils/diagnosisHelpers';

interface AIAnalysisButtonProps {
  equipment: any;
  description: string;
  priority?: string;
  files?: File[];
  disabled?: boolean;
  onAnalysisComplete?: (diagnosis: any) => void;
  onAnalysisError?: (error: Error) => void;
}

export function AIAnalysisButton({
  equipment,
  description,
  priority = 'medium',
  files = [],
  disabled = false,
  onAnalysisComplete,
  onAnalysisError,
}: AIAnalysisButtonProps) {
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [progress, setProgress] = useState('');
  const [analysisComplete, setAnalysisComplete] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAnalyze = async () => {
    if (!description.trim()) {
      setError('Please provide a description first');
      return;
    }

    setIsAnalyzing(true);
    setError(null);
    setAnalysisComplete(false);
    setProgress('Initializing...');

    try {
      const diagnosis = await runAIDiagnosis({
        equipment,
        description,
        priority,
        files,
        onProgress: (msg) => setProgress(msg),
      });

      setAnalysisComplete(true);
      setProgress('Analysis complete!');
      
      // Notify parent component
      onAnalysisComplete?.(diagnosis);

      // Auto-hide success message after 3 seconds
      setTimeout(() => {
        setAnalysisComplete(false);
        setProgress('');
      }, 3000);

    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Analysis failed';
      setError(errorMessage);
      setProgress('');
      onAnalysisError?.(err instanceof Error ? err : new Error(errorMessage));
    } finally {
      setIsAnalyzing(false);
    }
  };

  return (
    <div className="space-y-3">
      <button
        type="button"
        onClick={handleAnalyze}
        disabled={disabled || isAnalyzing || !description.trim()}
        className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-gradient-to-r from-purple-600 to-indigo-600 text-white rounded-lg hover:from-purple-700 hover:to-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-md hover:shadow-lg"
      >
        {isAnalyzing ? (
          <>
            <Loader2 className="h-5 w-5 animate-spin" />
            <span>Analyzing...</span>
          </>
        ) : analysisComplete ? (
          <>
            <CheckCircle className="h-5 w-5" />
            <span>Analysis Complete!</span>
          </>
        ) : (
          <>
            <Sparkles className="h-5 w-5" />
            <span>Get AI Diagnosis{files.length > 0 ?  with ${files.length} Image(s) : ''}</span>
          </>
        )}
      </button>

      {/* Progress Message */}
      {progress && !error && (
        <div className="flex items-center gap-2 text-sm text-purple-700 bg-purple-50 px-4 py-2 rounded-lg border border-purple-200">
          {isAnalyzing && <Loader2 className="h-4 w-4 animate-spin" />}
          {analysisComplete && <CheckCircle className="h-4 w-4" />}
          <span>{progress}</span>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="flex items-center gap-2 text-sm text-red-700 bg-red-50 px-4 py-2 rounded-lg border border-red-200">
          <AlertCircle className="h-4 w-4" />
          <span>{error}</span>
        </div>
      )}

      {/* File Info */}
      {files.length > 0 && !isAnalyzing && !analysisComplete && (
        <div className="text-xs text-gray-600 bg-blue-50 px-3 py-2 rounded border border-blue-200">
          <div className="flex items-center gap-2">
            <svg className="h-4 w-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>
              {files.length} file(s) will be analyzed for error codes, damage, and component identification
            </span>
          </div>
        </div>
      )}
    </div>
  );
}
