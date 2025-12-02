'use client';

import { useEffect, useState, Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import { equipmentApi } from '@/lib/api/equipment';
import { diagnosisApi, DiagnosisDecisionFeedback } from '@/lib/api/diagnosis';
import { Equipment } from '@/types';
import { Loader2, AlertCircle, CheckCircle, Package } from 'lucide-react';
import { DiagnosisCard, DiagnosisButton } from '@/components/diagnosis';
import { PartsAssignmentModal } from '@/components/PartsAssignmentModal';

function ServiceRequestPageInner() {
  const searchParams = useSearchParams();
  const qrCode = searchParams?.get('qr');
  
  const [equipment, setEquipment] = useState<Equipment | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [success, setSuccess] = useState(false);
  
  // AI Diagnosis state
  const [diagnosis, setDiagnosis] = useState<any>(null);
  const [diagnosisLoading, setDiagnosisLoading] = useState(false);
  
  // Parts Assignment state
  const [isPartsModalOpen, setIsPartsModalOpen] = useState(false);
  const [assignedParts, setAssignedParts] = useState<any[]>([]);
  
  const [formData, setFormData] = useState({
    description: '',
    priority: 'medium',
    requestedBy: '',
  });

  useEffect(() => {
    if (!qrCode) {
      setError('No QR code provided. Please scan a QR code to create a service request.');
      setLoading(false);
      return;
    }

    // Fetch equipment by QR code
    const fetchEquipment = async () => {
      try {
        setLoading(true);
        const data = await equipmentApi.getByQRCode(qrCode);
        setEquipment(data);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch equipment:', err);
        setError(`Equipment not found for QR code: ${qrCode}`);
      } finally {
        setLoading(false);
      }
    };

    fetchEquipment();
  }, [qrCode]);

  // Handle AI diagnosis completion
  const handleDiagnosisComplete = (diagnosisResult: any) => {
    setDiagnosis(diagnosisResult);
  };

  // Handle parts assignment
  const handlePartsAssign = (parts: any[]) => {
    setAssignedParts(parts);
    console.log('Parts assigned:', parts);
  };

  // Handle diagnosis accept/reject
  const handleDiagnosisAccept = async (diagnosisId: string) => {
    try {
      const feedback: DiagnosisDecisionFeedback = {
        decision: 'accepted',
        user_id: 1, // Would be real user ID
        user_role: 'technician',
        feedback_text: 'Diagnosis accepted by user'
      };
      
      await diagnosisApi.submitFeedback(diagnosisId, feedback);
      
      // Update local diagnosis state
      setDiagnosis((prev: any) => ({
        ...prev,
        decision_status: 'accepted',
        decided_at: new Date().toISOString(),
        feedback_text: feedback.feedback_text
      }));
      
    } catch (err) {
      alert('Failed to submit feedback');
    }
  };

  const handleDiagnosisReject = async (diagnosisId: string, feedbackText?: string) => {
    try {
      const feedback: DiagnosisDecisionFeedback = {
        decision: 'rejected',
        user_id: 1, // Would be real user ID
        user_role: 'technician',
        feedback_text: feedbackText || 'Diagnosis rejected by user'
      };
      
      await diagnosisApi.submitFeedback(diagnosisId, feedback);
      
      // Update local diagnosis state
      setDiagnosis((prev: any) => ({
        ...prev,
        decision_status: 'rejected',
        decided_at: new Date().toISOString(),
        feedback_text: feedback.feedback_text
      }));
      
    } catch (err) {
      alert('Failed to submit feedback');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!equipment) return;
    
    try {
      setSubmitting(true);
      
      // Create real service ticket
      const payload = {
        equipment_id: (equipment as any).id,
        qr_code: (equipment as any).qr_code || (equipment as any).qrCode,
        serial_number: (equipment as any).serial_number || (equipment as any).serialNumber,
        equipment_name: (equipment as any).equipment_name || (equipment as any).name,
        customer_id: (equipment as any).customer_id,
        customer_name: (equipment as any).customer_name || (equipment as any).customerName,
        customer_phone: '9999999999',
        issue_category: 'breakdown',
        issue_description: formData.description,
        priority: formData.priority as any,
        source: 'web',
        created_by: formData.requestedBy || 'web-user',
        notes: diagnosis?.summary ? `AI suggestion: ${diagnosis.summary}` : undefined,
        parts_requested: assignedParts.map(part => ({
          part_number: part.part_number,
          description: part.part_name,
          quantity: part.quantity,
          unit_price: part.unit_price,
          total_price: part.unit_price * part.quantity
        })),
      };
      const created = await (await import('@/lib/api/tickets')).ticketsApi.create(payload as any);
      console.log('Ticket created', created);
      
      setSuccess(true);
      setFormData({ description: '', priority: 'medium', requestedBy: '' });
      
    } catch (err) {
      alert(`Failed to create service request: ${err instanceof Error ? err.message : 'Unknown error'}`);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading equipment details...</p>
        </div>
      </div>
    );
  }

  if (error || !equipment) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
        <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-6">
          <div className="flex items-center gap-3 text-red-600 mb-4">
            <AlertCircle className="h-6 w-6" />
            <h2 className="text-lg font-semibold">Error</h2>
          </div>
          <p className="text-gray-700 mb-4">{error}</p>
          <div className="bg-gray-50 p-4 rounded-md">
            <p className="text-sm text-gray-600">
              <strong>QR Code:</strong> {qrCode || 'Not provided'}
            </p>
          </div>
        </div>
      </div>
    );
  }

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
        <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-6">
          <div className="flex items-center gap-3 text-green-600 mb-4">
            <CheckCircle className="h-6 w-6" />
            <h2 className="text-lg font-semibold">Service Request Created!</h2>
          </div>
          <p className="text-gray-700 mb-6">
            Your service request has been submitted successfully. Our team will contact you soon.
          </p>
          <div className="bg-gray-50 p-4 rounded-md mb-4">
            <p className="text-sm text-gray-600 mb-1">
              <strong>Equipment:</strong> {(equipment as any).equipment_name || (equipment as any).name || 'N/A'}
            </p>
            <p className="text-sm text-gray-600">
              <strong>Serial Number:</strong> {(equipment as any).serial_number || (equipment as any).serialNumber || 'N/A'}
            </p>
          </div>
          <button
            onClick={() => setSuccess(false)}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors"
          >
            Create Another Request
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-2xl mx-auto">
        <div className="bg-white rounded-lg shadow-lg p-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">
            Create Service Request
          </h1>

          {/* Equipment Details */}
          <div className="bg-blue-50 border border-blue-200 rounded-md p-4 mb-6">
            <h2 className="text-lg font-semibold text-blue-900 mb-3">Equipment Details</h2>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-blue-700 font-medium">Equipment Name</p>
                <p className="text-blue-900">{(equipment as any).equipment_name || (equipment as any).name || 'N/A'}</p>
              </div>
              <div>
                <p className="text-sm text-blue-700 font-medium">Serial Number</p>
                <p className="text-blue-900">{(equipment as any).serial_number || (equipment as any).serialNumber || 'N/A'}</p>
              </div>
              <div>
                <p className="text-sm text-blue-700 font-medium">Manufacturer</p>
                <p className="text-blue-900">{(equipment as any).manufacturer_name || (equipment as any).manufacturerName || 'N/A'}</p>
              </div>
              <div>
                <p className="text-sm text-blue-700 font-medium">Model</p>
                <p className="text-blue-900">{(equipment as any).model_number || (equipment as any).modelNumber || 'N/A'}</p>
              </div>
              {((equipment as any).customer_name || (equipment as any).customerName) && (
                <div className="col-span-2">
                  <p className="text-sm text-blue-700 font-medium">Hospital/Location</p>
                  <p className="text-blue-900">{(equipment as any).customer_name || (equipment as any).customerName}</p>
                </div>
              )}
              <div className="col-span-2">
                <p className="text-sm text-blue-700 font-medium">QR Code</p>
                <p className="text-blue-900 font-mono text-sm">{(equipment as any).qr_code || (equipment as any).qrCode || 'N/A'}</p>
              </div>
            </div>
          </div>

          {/* Service Request Form */}
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="requestedBy" className="block text-sm font-medium text-gray-700 mb-2">
                Your Name *
              </label>
              <input
                type="text"
                id="requestedBy"
                required
                value={formData.requestedBy}
                onChange={(e) => setFormData({ ...formData, requestedBy: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Enter your name"
              />
            </div>

            <div>
              <label htmlFor="priority" className="block text-sm font-medium text-gray-700 mb-2">
                Priority *
              </label>
              <select
                id="priority"
                required
                value={formData.priority}
                onChange={(e) => setFormData({ ...formData, priority: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="low">Low - Routine maintenance</option>
                <option value="medium">Medium - Issue affecting performance</option>
                <option value="high">High - Critical issue, equipment down</option>
              </select>
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
                Issue Description *
              </label>
              <textarea
                id="description"
                required
                rows={5}
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Describe the issue or service needed..."
              />
            </div>

            {/* AI Diagnosis Section */}
            {formData.description && (
              <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-sm font-medium text-purple-900">🤖 AI Assistant</h3>
                  <DiagnosisButton 
                    equipment={equipment}
                    description={formData.description}
                    priority={formData.priority}
                    onDiagnosisComplete={handleDiagnosisComplete}
                  />
                </div>
                <p className="text-xs text-purple-700">
                  Get AI-powered diagnosis suggestions based on your issue description
                </p>
              </div>
            )}

            {/* Parts Assignment Section */}
            <div className="bg-green-50 border border-green-200 rounded-lg p-4">
              <div className="flex items-center justify-between mb-3">
                <div>
                  <h3 className="text-sm font-medium text-green-900 flex items-center gap-2">
                    <Package className="h-4 w-4" />
                    Spare Parts Needed
                  </h3>
                  {assignedParts.length > 0 && (
                    <p className="text-xs text-green-700 mt-1">
                      {assignedParts.length} part{assignedParts.length > 1 ? 's' : ''} assigned • ₹{assignedParts.reduce((sum, p) => sum + (p.unit_price * p.quantity), 0).toLocaleString()}
                    </p>
                  )}
                </div>
                <button
                  type="button"
                  onClick={() => setIsPartsModalOpen(true)}
                  className="px-4 py-2 bg-green-600 text-white text-sm rounded-md hover:bg-green-700 transition-colors"
                >
                  {assignedParts.length > 0 ? 'Modify Parts' : 'Add Parts'}
                </button>
              </div>
              {assignedParts.length === 0 ? (
                <p className="text-xs text-green-700">
                  Select spare parts needed for this service request
                </p>
              ) : (
                <div className="mt-2 space-y-2">
                  {assignedParts.slice(0, 3).map((part) => (
                    <div key={part.id} className="flex justify-between text-xs bg-white p-2 rounded border border-green-100">
                      <span className="font-medium">{part.part_name}</span>
                      <span className="text-gray-600">
                        {part.quantity}x • ₹{part.unit_price * part.quantity}
                      </span>
                    </div>
                  ))}
                  {assignedParts.length > 3 && (
                    <p className="text-xs text-green-600 text-center">
                      +{assignedParts.length - 3} more part{assignedParts.length - 3 > 1 ? 's' : ''}
                    </p>
                  )}
                </div>
              )}
            </div>

            <div className="flex gap-4">
              <button
                type="submit"
                disabled={submitting}
                className="flex-1 bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                {submitting && <Loader2 className="h-5 w-5 animate-spin" />}
                {submitting ? 'Submitting...' : 'Submit Service Request'}
              </button>
            </div>
          </form>

          <div className="mt-6 pt-6 border-t border-gray-200">
            <p className="text-xs text-gray-500 text-center">
              By submitting this request, you agree to our service terms and conditions.
              Our team will contact you within 24 hours.
            </p>
          </div>
        </div>

        {/* AI Diagnosis Results */}
        {diagnosis && (
          <div className="mt-8">
            <DiagnosisCard 
              diagnosis={diagnosis}
              onAccept={handleDiagnosisAccept}
              onReject={handleDiagnosisReject}
              loading={diagnosisLoading}
            />
          </div>
        )}
      </div>

      {/* Parts Assignment Modal */}
      <PartsAssignmentModal
        open={isPartsModalOpen}
        onClose={() => setIsPartsModalOpen(false)}
        onAssign={handlePartsAssign}
        equipmentId={(equipment as any)?.id || 'unknown'}
        equipmentName={(equipment as any)?.equipment_name || (equipment as any)?.name || 'Equipment'}
      />
    </div>
  );
}

export default function ServiceRequestPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center bg-gray-50"><div className="text-center"><Loader2 className="h-10 w-10 animate-spin text-blue-600 mx-auto mb-4" /><p className="text-gray-600">Loading...</p></div></div>}>
      <ServiceRequestPageInner />
    </Suspense>
  );
}
