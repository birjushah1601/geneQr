'use client';

import { useState } from 'react';
import { DiagnosisCard, DiagnosisButton } from '@/components/diagnosis';
import { DiagnosisDecisionFeedback } from '@/lib/api/diagnosis';
import { Brain, CheckCircle, AlertCircle } from 'lucide-react';

// Mock equipment data for demo
const mockEquipment = {
  id: 'eq-001',
  equipment_name: 'MRI Scanner Pro X1',
  manufacturer_name: 'Siemens Healthcare',
  model_number: 'MAGNETOM Vida 3T',
  serial_number: 'SNK-2019-001-MRI',
  qr_code: 'QR001',
  customer_name: 'City General Hospital',
  status: 'operational' as const,
  installation_location: 'Radiology Department - Room 3',
  created_at: '2019-03-15T00:00:00Z',
  updated_at: '2024-11-18T00:00:00Z',
  created_by: 'admin',
  service_count: 12,
  qr_code_url: ''
};

export default function AIDiagnosisDemoPage() {
  const [diagnosis, setDiagnosis] = useState<any>(null);
  const [demoScenario, setDemoScenario] = useState('power_failure');
  const [feedback, setFeedback] = useState<string>('');

  const scenarios = {
    power_failure: {
      title: "Power Supply Failure",
      description: "Equipment won't start, no power LED, display is blank, cooling fan not running, no response to power button press"
    },
    imaging_quality: {
      title: "Image Quality Issues", 
      description: "MRI images are showing artifacts, noise in scan results, poor contrast resolution, patients complaining about longer scan times"
    },
    cooling_system: {
      title: "Cooling System Problem",
      description: "Equipment overheating alarms, reduced performance, helium pressure warnings, compressor making unusual noises"
    },
    software_error: {
      title: "Software Malfunction",
      description: "System crashes during scans, error messages on screen, unable to save images, interface freezing intermittently"
    }
  };

  const handleDiagnosisComplete = (diagnosisResult: any) => {
    setDiagnosis(diagnosisResult);
  };

  const handleDiagnosisAccept = async (diagnosisId: string) => {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500));
    
    setDiagnosis((prev: any) => ({
      ...prev,
      decision_status: 'accepted',
      decided_at: new Date().toISOString(),
      feedback_text: 'Diagnosis accepted by technician'
    }));
    
    setFeedback('Ã¢Å“â€¦ Diagnosis accepted! AI learning improved.');
  };

  const handleDiagnosisReject = async (diagnosisId: string, feedbackText?: string) => {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500));
    
    setDiagnosis((prev: any) => ({
      ...prev,
      decision_status: 'rejected',
      decided_at: new Date().toISOString(),
      feedback_text: feedbackText || 'Diagnosis rejected by technician'
    }));
    
    setFeedback('Ã¢ÂÅ’ Diagnosis rejected. Thank you for the feedback to improve our AI!');
  };

  const resetDemo = () => {
    setDiagnosis(null);
    setFeedback('');
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg shadow-lg p-6 mb-8">
          <div className="flex items-center gap-4 mb-4">
            <div className="p-3 bg-white/20 rounded-lg">
              <Brain className="h-8 w-8" />
            </div>
            <div>
              <h1 className="text-3xl font-bold">Ã°Å¸Â¤â€“ AI-Assisted Diagnosis Demo</h1>
              <p className="text-purple-100">Experience intelligent medical equipment diagnostics with confidence scoring</p>
            </div>
          </div>
          
          <div className="grid grid-cols-3 gap-4 text-center">
            <div className="bg-white/10 rounded-lg p-3">
              <div className="text-2xl font-bold">77.5%</div>
              <div className="text-sm text-purple-100">Avg Confidence</div>
            </div>
            <div className="bg-white/10 rounded-lg p-3">
              <div className="text-2xl font-bold">4</div>
              <div className="text-sm text-purple-100">Analysis Factors</div>
            </div>
            <div className="bg-white/10 rounded-lg p-3">
              <div className="text-2xl font-bold">2-4hr</div>
              <div className="text-sm text-purple-100">Est Resolution</div>
            </div>
          </div>
        </div>

        {/* Equipment Info */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Ã°Å¸ÂÂ¥ Demo Equipment</h2>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="font-medium text-gray-600">Equipment:</span>
              <span className="ml-2">{mockEquipment.equipment_name}</span>
            </div>
            <div>
              <span className="font-medium text-gray-600">Manufacturer:</span>
              <span className="ml-2">{mockEquipment.manufacturer_name}</span>
            </div>
            <div>
              <span className="font-medium text-gray-600">Model:</span>
              <span className="ml-2">{mockEquipment.model_number}</span>
            </div>
            <div>
              <span className="font-medium text-gray-600">Location:</span>
              <span className="ml-2">{mockEquipment.installation_location}</span>
            </div>
          </div>
        </div>

        {/* Scenario Selector */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Ã¢Å¡Â¡ Select Issue Scenario</h2>
          <div className="grid grid-cols-2 gap-4">
            {Object.entries(scenarios).map(([key, scenario]) => (
              <button
                key={key}
                onClick={() => {
                  setDemoScenario(key);
                  resetDemo();
                }}
                className={`p-4 border rounded-lg text-left transition-colors ${
                  demoScenario === key
                    ? 'border-blue-500 bg-blue-50 text-blue-900'
                    : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
                }`}
              >
                <div className="font-medium mb-1">{scenario.title}</div>
                <div className="text-sm text-gray-600 line-clamp-2">{scenario.description}</div>
              </button>
            ))}
          </div>
        </div>

        {/* AI Diagnosis Section */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold">Ã°Å¸â€Â¬ AI Diagnosis</h2>
            {diagnosis && (
              <button
                onClick={resetDemo}
                className="text-sm text-gray-500 hover:text-gray-700 border border-gray-300 px-3 py-1 rounded-md"
              >
                Reset Demo
              </button>
            )}
          </div>

          {/* Issue Description */}
          <div className="mb-6 p-4 bg-gray-50 rounded-lg">
            <h3 className="font-medium text-gray-900 mb-2">Issue Description:</h3>
            <p className="text-gray-700">{scenarios[demoScenario as keyof typeof scenarios].description}</p>
          </div>

          {/* AI Diagnosis Button */}
          <div className="flex justify-center mb-6">
            <DiagnosisButton 
              equipment={mockEquipment}
              description={scenarios[demoScenario as keyof typeof scenarios].description}
              priority="high"
              onDiagnosisComplete={handleDiagnosisComplete}
              className="px-8 py-3 text-base"
            />
          </div>

          {/* Feedback Message */}
          {feedback && (
            <div className={`p-3 rounded-lg mb-6 text-center ${
              feedback.includes('Ã¢Å“â€¦') ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'
            }`}>
              {feedback}
            </div>
          )}
        </div>

        {/* AI Diagnosis Results */}
        {diagnosis && (
          <div className="mb-8">
            <DiagnosisCard 
              diagnosis={diagnosis}
              onAccept={handleDiagnosisAccept}
              onReject={handleDiagnosisReject}
            />
          </div>
        )}

        {/* Demo Instructions */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
          <h3 className="font-semibold text-blue-900 mb-3">Ã°Å¸Å½Â¯ How to Test the Demo</h3>
          <ol className="list-decimal list-inside space-y-2 text-blue-800 text-sm">
            <li>Choose an issue scenario from the cards above</li>
            <li>Click the "Get AI Diagnosis" button to run analysis</li>
            <li>Review the AI diagnosis with confidence scoring</li>
            <li>Accept or reject the diagnosis with feedback</li>
            <li>See how the AI learns from your decisions</li>
          </ol>
        </div>
      </div>
    </div>
  );
}