'use client';

import { useState, useRef, useEffect } from 'react';
import { QrCode, Package, Phone, AlertCircle, CheckCircle2, Loader2, ArrowRight, Upload, Camera, X } from 'lucide-react';
import { Html5Qrcode } from 'html5-qrcode';
import { equipmentApi } from '@/lib/api/equipment';
import { ticketsApi } from '@/lib/api/tickets';
import type { Equipment, ServiceTicket } from '@/types';

type WorkflowStep = 'input' | 'lookup' | 'details' | 'success';
type Priority = 'critical' | 'high' | 'medium' | 'low';
type ScanMode = 'upload' | 'camera';

export default function TestQRWorkflowPage() {
  // State management
  const [step, setStep] = useState<WorkflowStep>('input');
  const [qrCode, setQrCode] = useState('');
  const [issueDescription, setIssueDescription] = useState('');
  const [customerPhone, setCustomerPhone] = useState('');
  const [priority, setPriority] = useState<Priority>('medium');
  const [scanMode, setScanMode] = useState<ScanMode>('upload');
  
  const [equipment, setEquipment] = useState<Equipment | null>(null);
  const [ticket, setTicket] = useState<ServiceTicket | null>(null);
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [scanningCamera, setScanningCamera] = useState(false);
  const [uploadedImage, setUploadedImage] = useState<string | null>(null);

  // Refs
  const fileInputRef = useRef<HTMLInputElement>(null);
  const html5QrCodeRef = useRef<Html5Qrcode | null>(null);

  // Cleanup camera on unmount
  useEffect(() => {
    return () => {
      if (html5QrCodeRef.current && scanningCamera) {
        html5QrCodeRef.current.stop().catch(console.error);
      }
    };
  }, [scanningCamera]);

  // Handle file upload
  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setLoading(true);
    setError(null);

    try {
      // Create object URL for preview
      const imageUrl = URL.createObjectURL(file);
      setUploadedImage(imageUrl);

      // Scan QR code from file
      const html5QrCode = new Html5Qrcode('qr-reader');
      const decodedText = await html5QrCode.scanFile(file, false);
      await html5QrCode.clear();
      
      setQrCode(decodedText);
      await lookupEquipment(decodedText);
    } catch (err) {
      setError('Could not read QR code from image. Please try another image or use camera scan.');
      setUploadedImage(null);
    } finally {
      setLoading(false);
    }
  };

  // Start camera scanning
  const startCameraScanning = async () => {
    setScanningCamera(true);
    setError(null);

    try {
      const html5QrCode = new Html5Qrcode('camera-reader');
      html5QrCodeRef.current = html5QrCode;

      await html5QrCode.start(
        { facingMode: 'environment' },
        {
          fps: 10,
          qrbox: (viewfinderWidth: number, viewfinderHeight: number) => {
            const size = Math.floor(Math.min(viewfinderWidth, viewfinderHeight) * 0.8);
            return { width: size, height: size };
          },
        },
        async (decodedText) => {
          // Successfully scanned
          setQrCode(decodedText);
          await stopCameraScanning();
          await lookupEquipment(decodedText);
        },
        (errorMessage) => {
          // Scanning error (can be ignored - continuous scanning)
        }
      );
    } catch (err) {
      setError('Could not access camera. Please check permissions or use file upload.');
      setScanningCamera(false);
    }
  };

  // Stop camera scanning
  const stopCameraScanning = async () => {
    if (html5QrCodeRef.current) {
      try {
        await html5QrCodeRef.current.stop();
        html5QrCodeRef.current = null;
      } catch (err) {
        console.error('Error stopping camera:', err);
      }
    }
    setScanningCamera(false);
  };

  // Look up equipment by QR code
  const lookupEquipment = async (code: string) => {
    if (!code.trim()) {
      setError('Please scan a QR code');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await equipmentApi.getByQRCode(code.trim());
      setEquipment(result);
      setStep('details');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Equipment not found. Please check the QR code.');
    } finally {
      setLoading(false);
    }
  };

  // Step 2: Create ticket with issue details
  const handleCreateTicket = async () => {
    if (!equipment || !issueDescription.trim() || !customerPhone.trim()) {
      setError('Please fill in all required fields');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const ticketData = {
        equipment_id: equipment.id,
        qr_code: equipment.qr_code,
        serial_number: equipment.serial_number,
        customer_phone: customerPhone,
        customer_whatsapp: customerPhone,
        issue_category: 'breakdown',
        issue_description: issueDescription,
        priority: priority,
        source: 'web' as const,
        created_by: 'qr-test-interface',
      };

      const createdTicket = await ticketsApi.create(ticketData);
      setTicket(createdTicket);
      setStep('success');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create ticket');
    } finally {
      setLoading(false);
    }
  };

  // Auto-detect priority based on keywords (same logic as WhatsApp handler)
  const handleIssueChange = (value: string) => {
    setIssueDescription(value);
    
    const lower = value.toLowerCase();
    
    // Critical keywords
    if (/(urgent|emergency|critical|down|not working|stopped|patient)/i.test(lower)) {
      setPriority('critical');
    }
    // High priority keywords
    else if (/(error|alarm|warning|issue|problem|broken)/i.test(lower)) {
      setPriority('high');
    }
    // Medium priority keywords
    else if (/(maintenance|service|check|noise|slow)/i.test(lower)) {
      setPriority('medium');
    } else {
      setPriority('medium');
    }
  };

  // Reset form
  const handleReset = async () => {
    if (scanningCamera) {
      await stopCameraScanning();
    }
    setStep('input');
    setQrCode('');
    setIssueDescription('');
    setCustomerPhone('');
    setPriority('medium');
    setEquipment(null);
    setTicket(null);
    setError(null);
    setUploadedImage(null);
    setScanMode('upload');
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  // Priority badge colors
  const getPriorityColor = (p: Priority) => {
    switch (p) {
      case 'critical': return 'bg-red-100 text-red-800 border-red-300';
      case 'high': return 'bg-orange-100 text-orange-800 border-orange-300';
      case 'medium': return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'low': return 'bg-green-100 text-green-800 border-green-300';
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 p-4 md:p-8">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-600 rounded-full mb-4">
            <QrCode className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            QR Code Workflow Tester
          </h1>
          <p className="text-gray-600">
            Test the WhatsApp-like QR → Ticket creation flow
          </p>
        </div>

        {/* Progress Steps */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
          <div className="flex items-center justify-between">
            {/* Step 1 */}
            <div className="flex-1">
              <div className={`flex items-center ${step === 'input' ? 'text-blue-600' : step === 'lookup' || step === 'details' || step === 'success' ? 'text-green-600' : 'text-gray-400'}`}>
                <div className={`w-8 h-8 rounded-full flex items-center justify-center border-2 font-semibold ${step === 'input' ? 'border-blue-600 bg-blue-50' : step === 'lookup' || step === 'details' || step === 'success' ? 'border-green-600 bg-green-50' : 'border-gray-300'}`}>
                  {step === 'lookup' || step === 'details' || step === 'success' ? '✓' : '1'}
                </div>
                <span className="ml-2 text-sm font-medium hidden md:inline">Scan QR</span>
              </div>
            </div>

            <ArrowRight className={`w-5 h-5 mx-2 ${step === 'details' || step === 'success' ? 'text-green-600' : 'text-gray-300'}`} />

            {/* Step 2 */}
            <div className="flex-1">
              <div className={`flex items-center ${step === 'details' ? 'text-blue-600' : step === 'success' ? 'text-green-600' : 'text-gray-400'}`}>
                <div className={`w-8 h-8 rounded-full flex items-center justify-center border-2 font-semibold ${step === 'details' ? 'border-blue-600 bg-blue-50' : step === 'success' ? 'border-green-600 bg-green-50' : 'border-gray-300'}`}>
                  {step === 'success' ? '✓' : '2'}
                </div>
                <span className="ml-2 text-sm font-medium hidden md:inline">Issue Details</span>
              </div>
            </div>

            <ArrowRight className={`w-5 h-5 mx-2 ${step === 'success' ? 'text-green-600' : 'text-gray-300'}`} />

            {/* Step 3 */}
            <div className="flex-1">
              <div className={`flex items-center ${step === 'success' ? 'text-green-600' : 'text-gray-400'}`}>
                <div className={`w-8 h-8 rounded-full flex items-center justify-center border-2 font-semibold ${step === 'success' ? 'border-green-600 bg-green-50' : 'border-gray-300'}`}>
                  {step === 'success' ? '✓' : '3'}
                </div>
                <span className="ml-2 text-sm font-medium hidden md:inline">Ticket Created</span>
              </div>
            </div>
          </div>
        </div>

        {/* Error Display */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6 flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <p className="text-sm font-medium text-red-800">Error</p>
              <p className="text-sm text-red-700 mt-1">{error}</p>
            </div>
          </div>
        )}

        {/* Main Content Card */}
        <div className="bg-white rounded-lg shadow-md border border-gray-200 p-6 md:p-8">
          {/* Step 1: QR Code Scanning */}
          {step === 'input' && (
            <div className="space-y-6">
              {/* Scan Mode Selector */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-3">
                  Choose Scan Method
                </label>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    onClick={() => {
                      setScanMode('upload');
                      if (scanningCamera) stopCameraScanning();
                    }}
                    className={`p-4 border-2 rounded-lg flex flex-col items-center gap-2 transition-all ${
                      scanMode === 'upload'
                        ? 'border-blue-600 bg-blue-50 text-blue-700'
                        : 'border-gray-300 hover:border-gray-400'
                    }`}
                  >
                    <Upload className="w-6 h-6" />
                    <span className="text-sm font-medium">Upload Image</span>
                  </button>
                  <button
                    onClick={() => {
                      setScanMode('camera');
                      setUploadedImage(null);
                    }}
                    className={`p-4 border-2 rounded-lg flex flex-col items-center gap-2 transition-all ${
                      scanMode === 'camera'
                        ? 'border-blue-600 bg-blue-50 text-blue-700'
                        : 'border-gray-300 hover:border-gray-400'
                    }`}
                  >
                    <Camera className="w-6 h-6" />
                    <span className="text-sm font-medium">Use Camera</span>
                  </button>
                </div>
              </div>

              {/* Upload Mode */}
              {scanMode === 'upload' && !uploadedImage && (
                <div>
                  <div
                    onClick={() => fileInputRef.current?.click()}
                    className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-blue-500 hover:bg-blue-50 transition-all cursor-pointer"
                  >
                    <Upload className="w-12 h-12 text-gray-400 mx-auto mb-3" />
                    <p className="text-sm font-medium text-gray-700 mb-1">
                      Click to upload QR code image
                    </p>
                    <p className="text-xs text-gray-500">
                      PNG, JPG up to 10MB
                    </p>
                  </div>
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    onChange={handleFileUpload}
                    className="hidden"
                  />
                  <div id="qr-reader" className="w-px h-px absolute -left-[9999px]"></div>
                </div>
              )}

              {/* Uploaded Image Preview */}
              {scanMode === 'upload' && uploadedImage && (
                <div className="relative">
                  <img
                    src={uploadedImage}
                    alt="Uploaded QR code"
                    className="w-full max-h-64 object-contain border border-gray-300 rounded-lg"
                  />
                  <button
                    onClick={() => {
                      setUploadedImage(null);
                      setQrCode('');
                      if (fileInputRef.current) fileInputRef.current.value = '';
                    }}
                    className="absolute top-2 right-2 p-2 bg-red-500 hover:bg-red-600 text-white rounded-full"
                  >
                    <X className="w-4 h-4" />
                  </button>
                  {qrCode && (
                    <div className="mt-3 p-3 bg-green-50 border border-green-200 rounded-lg">
                      <p className="text-sm font-medium text-green-800">
                        ✓ QR Code Detected: {qrCode}
                      </p>
                    </div>
                  )}
                </div>
              )}

              {/* Camera Mode */}
              {scanMode === 'camera' && (
                <div>
                  <div id="camera-reader" className="border border-gray-300 rounded-lg overflow-hidden bg-black h-80"></div>
                  
                  {!scanningCamera ? (
                    <button
                      onClick={startCameraScanning}
                      disabled={loading}
                      className="w-full mt-4 py-3 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-300 text-white font-medium rounded-lg transition-colors flex items-center justify-center gap-2"
                    >
                      <Camera className="w-5 h-5" />
                      Start Camera Scan
                    </button>
                  ) : (
                    <button
                      onClick={stopCameraScanning}
                      className="w-full mt-4 py-3 px-4 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition-colors flex items-center justify-center gap-2"
                    >
                      <X className="w-5 h-5" />
                      Stop Scanning
                    </button>
                  )}
                  
                  {qrCode && (
                    <div className="mt-3 p-3 bg-green-50 border border-green-200 rounded-lg">
                      <p className="text-sm font-medium text-green-800">
                        ✓ QR Code Scanned: {qrCode}
                      </p>
                    </div>
                  )}
                </div>
              )}

              {/* Loading State */}
              {loading && (
                <div className="flex items-center justify-center py-4">
                  <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
                  <span className="ml-3 text-gray-600">Processing QR code...</span>
                </div>
              )}
            </div>
          )}

          {/* Step 2: Issue Details */}
          {step === 'details' && equipment && (
            <div className="space-y-6">
              {/* Equipment Info Card */}
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <h3 className="text-sm font-semibold text-blue-900 mb-3 flex items-center gap-2">
                  <Package className="w-4 h-4" />
                  Equipment Found
                </h3>
                <div className="grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <span className="text-blue-700 font-medium">Name:</span>
                    <p className="text-blue-900">{equipment.equipment_name}</p>
                  </div>
                  <div>
                    <span className="text-blue-700 font-medium">Serial:</span>
                    <p className="text-blue-900">{equipment.serial_number}</p>
                  </div>
                  <div>
                    <span className="text-blue-700 font-medium">Customer:</span>
                    <p className="text-blue-900">{equipment.customer_name}</p>
                  </div>
                  <div>
                    <span className="text-blue-700 font-medium">Location:</span>
                    <p className="text-blue-900">{equipment.installation_location || 'N/A'}</p>
                  </div>
                </div>
              </div>

              {/* Customer Phone */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Customer Phone Number *
                </label>
                <div className="relative">
                  <Phone className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                  <input
                    type="tel"
                    value={customerPhone}
                    onChange={(e) => setCustomerPhone(e.target.value)}
                    placeholder="+91 98765 43210"
                    className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>

              {/* Issue Description */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Issue Description *
                </label>
                <textarea
                  value={issueDescription}
                  onChange={(e) => handleIssueChange(e.target.value)}
                  placeholder="Describe the problem with the equipment..."
                  rows={4}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
                />
                <p className="mt-2 text-xs text-gray-500">
                  Keywords like &quot;urgent&quot;, &quot;critical&quot;, &quot;error&quot; will automatically increase priority
                </p>
              </div>

              {/* Priority Display */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Detected Priority
                </label>
                <div className={`inline-flex items-center px-4 py-2 rounded-full border text-sm font-medium ${getPriorityColor(priority)}`}>
                  <span className="w-2 h-2 rounded-full bg-current mr-2"></span>
                  {priority.toUpperCase()}
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-3">
                <button
                  onClick={handleReset}
                  className="flex-1 py-3 px-4 border border-gray-300 hover:bg-gray-50 text-gray-700 font-medium rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreateTicket}
                  disabled={loading || !issueDescription.trim() || !customerPhone.trim()}
                  className="flex-1 py-3 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors flex items-center justify-center gap-2"
                >
                  {loading ? (
                    <>
                      <Loader2 className="w-5 h-5 animate-spin" />
                      Creating Ticket...
                    </>
                  ) : (
                    <>
                      <CheckCircle2 className="w-5 h-5" />
                      Create Ticket
                    </>
                  )}
                </button>
              </div>
            </div>
          )}

          {/* Step 3: Success */}
          {step === 'success' && ticket && equipment && (
            <div className="text-center space-y-6">
              <div className="inline-flex items-center justify-center w-20 h-20 bg-green-100 rounded-full">
                <CheckCircle2 className="w-10 h-10 text-green-600" />
              </div>

              <div>
                <h2 className="text-2xl font-bold text-gray-900 mb-2">
                  Service Request Confirmed!
                </h2>
                <p className="text-gray-600">
                  Your ticket has been created successfully
                </p>
              </div>

              {/* Ticket Details Card */}
              <div className="bg-gradient-to-br from-green-50 to-blue-50 border border-green-200 rounded-lg p-6 text-left">
                <div className="space-y-4">
                  <div>
                    <span className="text-sm text-gray-600">Ticket Number</span>
                    <p className="text-2xl font-bold text-gray-900">{ticket.ticket_number}</p>
                  </div>
                  
                  <div className="grid grid-cols-2 gap-4 pt-4 border-t border-green-200">
                    <div>
                      <span className="text-sm text-gray-600">Equipment</span>
                      <p className="font-medium text-gray-900">{equipment.equipment_name}</p>
                    </div>
                    <div>
                      <span className="text-sm text-gray-600">Serial Number</span>
                      <p className="font-medium text-gray-900">{equipment.serial_number}</p>
                    </div>
                    <div>
                      <span className="text-sm text-gray-600">Priority</span>
                      <div className={`inline-flex items-center px-3 py-1 rounded-full border text-xs font-medium ${getPriorityColor(ticket.priority)}`}>
                        {ticket.priority.toUpperCase()}
                      </div>
                    </div>
                    <div>
                      <span className="text-sm text-gray-600">Status</span>
                      <p className="font-medium text-gray-900">{ticket.status.toUpperCase()}</p>
                    </div>
                  </div>

                  <div className="pt-4 border-t border-green-200">
                    <span className="text-sm text-gray-600">Issue Description</span>
                    <p className="text-gray-900 mt-1">{ticket.issue_description}</p>
                  </div>
                </div>
              </div>

              {/* WhatsApp-like Confirmation Message */}
              <div className="bg-gray-50 border border-gray-200 rounded-lg p-4 text-left">
                <p className="text-xs text-gray-500 mb-2">This is what customer would see on WhatsApp:</p>
                <div className="bg-white rounded-lg p-4 shadow-sm">
                  <p className="text-sm whitespace-pre-line">
                    ✅ <strong>Service Request Confirmed</strong>
                    {'\n\n'}
                    Ticket Number: <strong>{ticket.ticket_number}</strong>
                    {'\n'}
                    Equipment: {equipment.equipment_name}
                    {'\n'}
                    Serial: {equipment.serial_number}
                    {'\n'}
                    Priority: {ticket.priority.toUpperCase()}
                    {'\n\n'}
                    Our engineer will contact you soon.
                    {'\n'}
                    Thank you!
                  </p>
                </div>
              </div>

              <button
                onClick={handleReset}
                className="w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
              >
                Test Another QR Code
              </button>
            </div>
          )}
        </div>

        {/* Info Footer */}
        <div className="mt-6 text-center text-sm text-gray-500">
          <p>💡 This interface simulates the WhatsApp workflow for testing purposes</p>
          <p className="mt-1">Once WhatsApp API keys are configured, the same flow will work via messaging</p>
        </div>
      </div>
    </div>
  );
}
