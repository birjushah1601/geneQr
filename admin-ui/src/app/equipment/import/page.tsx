'use client';

import { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Upload, FileText, CheckCircle, ArrowLeft } from 'lucide-react';

export default function EquipmentImportPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [dragActive, setDragActive] = useState(false);
  const [importResult, setImportResult] = useState<any>(null);

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
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      const result = {
        total: 400,
        success: 398,
        failed: 2,
      };
      
      setImportResult(result);
      localStorage.setItem('equipment_imported', 'true');
    } catch (error) {
      alert('Failed to import equipment');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4">
          <Button 
            variant="ghost" 
            onClick={() => router.push('/dashboard')}
            className="mb-2"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Dashboard
          </Button>
          <h1 className="text-2xl font-bold">Import Equipment</h1>
          <p className="text-gray-600">Upload a CSV file with your equipment installations</p>
        </div>
      </div>

      {/* Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-3xl mx-auto">
          <Card className="p-8">
            {!importResult ? (
              <>
                <Alert className="mb-6">
                  <AlertDescription>
                    <strong>CSV Format:</strong> equipment_name, serial_number, customer_name, 
                    installation_location, model_number, category, installation_date
                  </AlertDescription>
                </Alert>

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
                        <p className="text-sm text-gray-500">{(file.size / 1024).toFixed(2)} KB</p>
                      </div>
                      <div className="flex gap-4 justify-center">
                        <Button onClick={handleImport} disabled={loading}>
                          {loading ? 'Importing...' : 'Import Equipment'}
                        </Button>
                        <Button variant="outline" onClick={() => setFile(null)} disabled={loading}>
                          Remove
                        </Button>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-4">
                      <Upload className="w-16 h-16 mx-auto text-gray-400" />
                      <div>
                        <p className="text-lg font-semibold mb-2">Drag and drop your CSV file here</p>
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
                    </div>
                  )}
                </div>
              </>
            ) : (
              <div className="text-center space-y-6 py-8">
                <CheckCircle className="w-24 h-24 mx-auto text-green-600" />
                <div>
                  <h3 className="text-2xl font-bold text-green-600 mb-2">Import Successful!</h3>
                  <p className="text-gray-600 mb-6">Your equipment has been imported</p>
                </div>
                
                <div className="bg-gray-50 rounded-lg p-6 max-w-md mx-auto">
                  <div className="grid grid-cols-3 gap-4 text-center">
                    <div>
                      <p className="text-3xl font-bold text-blue-600">{importResult.total}</p>
                      <p className="text-sm text-gray-500">Total</p>
                    </div>
                    <div>
                      <p className="text-3xl font-bold text-green-600">{importResult.success}</p>
                      <p className="text-sm text-gray-500">Success</p>
                    </div>
                    <div>
                      <p className="text-3xl font-bold text-red-600">{importResult.failed}</p>
                      <p className="text-sm text-gray-500">Failed</p>
                    </div>
                  </div>
                </div>

                <Button onClick={() => router.push('/dashboard')}>
                  Go to Dashboard
                </Button>
              </div>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
}
