'use client';

import { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Upload, FileText, CheckCircle, XCircle, ArrowRight, SkipForward } from 'lucide-react';

export default function EquipmentImport() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [dragActive, setDragActive] = useState(false);
  const [importResult, setImportResult] = useState<any>(null);

  // Get manufacturer from localStorage
  const manufacturer = typeof window !== 'undefined' 
    ? JSON.parse(localStorage.getItem('current_manufacturer') || '{}')
    : {};

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const droppedFile = e.dataTransfer.files[0];
      if (droppedFile.type === 'text/csv' || droppedFile.name.endsWith('.csv')) {
        setFile(droppedFile);
      } else {
        alert('Please upload a CSV file');
      }
    }
  }, []);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  const handleImport = async () => {
    if (!file) return;
    
    setLoading(true);
    try {
      // Simulate import (in production, this would call the API)
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Mock result
      const result = {
        total: 400,
        success: 398,
        failed: 2,
        equipmentIds: Array.from({ length: 398 }, (_, i) => `EQ-${i + 1}`),
      };
      
      setImportResult(result);
      localStorage.setItem('equipment_imported', 'true');
      
      // Auto-redirect after 2 seconds
      setTimeout(() => {
        router.push('/onboarding/engineers');
      }, 2000);
    } catch (error) {
      alert('Failed to import equipment');
    } finally {
      setLoading(false);
    }
  };

  const handleSkip = () => {
    router.push('/onboarding/engineers');
  };

  const handleSkipAll = () => {
    router.push('/dashboard');
  };

  return (
    <Card className="p-8">
      {/* Progress Steps */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center opacity-50">
            <div className="w-10 h-10 rounded-full bg-green-600 text-white flex items-center justify-center">
              <CheckCircle className="w-6 h-6" />
            </div>
            <div className="ml-3">
              <p className="font-semibold">Manufacturer</p>
              <p className="text-sm text-gray-500">Completed</p>
            </div>
          </div>
          
          <div className="flex-1 h-1 mx-4 bg-blue-600"></div>
          
          <div className="flex items-center">
            <div className="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
              2
            </div>
            <div className="ml-3">
              <p className="font-semibold">Equipment Import</p>
              <p className="text-sm text-gray-500">CSV upload</p>
            </div>
          </div>
          
          <div className="flex-1 h-1 mx-4 bg-gray-200"></div>
          
          <div className="flex items-center opacity-50">
            <div className="w-10 h-10 rounded-full bg-gray-300 text-gray-600 flex items-center justify-center font-bold">
              3
            </div>
            <div className="ml-3">
              <p className="font-semibold">Engineers</p>
              <p className="text-sm text-gray-500">Service team</p>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="text-center mb-6">
        <h2 className="text-2xl font-bold mb-2">Import Equipment</h2>
        <p className="text-gray-600">
          Upload a CSV file with your equipment installations for <strong>{manufacturer.name || 'your company'}</strong>
        </p>
      </div>

      {!importResult ? (
        <>
          {/* CSV Format Info */}
          <Alert className="mb-6">
            <AlertDescription>
              <strong>CSV Format Required:</strong> equipment_name, serial_number, customer_name, 
              installation_location, model_number, category, installation_date
              <br />
              <a 
                href="/sample-equipment.csv" 
                download 
                className="text-blue-600 hover:underline mt-2 inline-block"
              >
                Download sample CSV template
              </a>
            </AlertDescription>
          </Alert>

          {/* Upload Area */}
          <div
            className={`border-2 border-dashed rounded-lg p-12 text-center transition-colors ${
              dragActive ? 'border-blue-600 bg-blue-50' : 'border-gray-300 bg-gray-50'
            }`}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
          >
            {file ? (
              <div className="space-y-4">
                <FileText className="w-16 h-16 mx-auto text-blue-600" />
                <div>
                  <p className="font-semibold text-lg">{file.name}</p>
                  <p className="text-sm text-gray-500">
                    {(file.size / 1024).toFixed(2)} KB
                  </p>
                </div>
                <div className="flex gap-4 justify-center">
                  <Button onClick={handleImport} disabled={loading}>
                    {loading ? 'Importing...' : 'Import Equipment'}
                  </Button>
                  <Button 
                    variant="outline" 
                    onClick={() => setFile(null)}
                    disabled={loading}
                  >
                    Remove
                  </Button>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                <Upload className="w-16 h-16 mx-auto text-gray-400" />
                <div>
                  <p className="text-lg font-semibold mb-2">
                    Drag and drop your CSV file here
                  </p>
                  <p className="text-gray-500 mb-4">or</p>
                  <label className="cursor-pointer">
                    <input
                      type="file"
                      className="hidden"
                      accept=".csv"
                      onChange={handleFileSelect}
                    />
                    <span className="inline-flex items-center justify-center px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700">
                      Browse Files
                    </span>
                  </label>
                </div>
                <p className="text-sm text-gray-400">
                  Supported format: CSV (up to 10MB)
                </p>
              </div>
            )}
          </div>

          {/* Action Buttons */}
          <div className="flex justify-between gap-4 pt-6 mt-6 border-t">
            <Button 
              variant="outline" 
              onClick={handleSkip}
              disabled={loading}
            >
              <SkipForward className="w-4 h-4 mr-2" />
              Skip for Now
            </Button>
            
            <div className="flex gap-4">
              <Button 
                variant="ghost" 
                onClick={handleSkipAll}
                disabled={loading}
              >
                Complete Setup Later
              </Button>
              <Button 
                onClick={handleImport} 
                disabled={!file || loading}
              >
                Import & Continue
                <ArrowRight className="w-4 h-4 ml-2" />
              </Button>
            </div>
          </div>
        </>
      ) : (
        // Import Success
        <div className="text-center space-y-6 py-8">
          <CheckCircle className="w-24 h-24 mx-auto text-green-600" />
          <div>
            <h3 className="text-2xl font-bold text-green-600 mb-2">
              Import Successful!
            </h3>
            <p className="text-gray-600 mb-6">
              Your equipment has been imported successfully
            </p>
          </div>
          
          <div className="bg-gray-50 rounded-lg p-6 max-w-md mx-auto">
            <div className="grid grid-cols-3 gap-4 text-center">
              <div>
                <p className="text-3xl font-bold text-blue-600">{importResult.total}</p>
                <p className="text-sm text-gray-500">Total Records</p>
              </div>
              <div>
                <p className="text-3xl font-bold text-green-600">{importResult.success}</p>
                <p className="text-sm text-gray-500">Successful</p>
              </div>
              <div>
                <p className="text-3xl font-bold text-red-600">{importResult.failed}</p>
                <p className="text-sm text-gray-500">Failed</p>
              </div>
            </div>
          </div>

          <p className="text-gray-500">
            Redirecting to engineer setup...
          </p>
        </div>
      )}
    </Card>
  );
}
