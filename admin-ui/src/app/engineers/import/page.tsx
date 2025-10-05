'use client';

import { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert } from '@/components/ui/alert';
import { ArrowLeft, Upload, FileText, CheckCircle2, X } from 'lucide-react';

export default function EngineersImportPage() {
  const router = useRouter();
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadComplete, setUploadComplete] = useState(false);
  const [importStats, setImportStats] = useState({
    total: 0,
    success: 0,
    failed: 0,
  });
  const [dragActive, setDragActive] = useState(false);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
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

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!file) return;

    setUploading(true);

    // Simulate upload and processing
    await new Promise((resolve) => setTimeout(resolve, 2000));

    // Simulate import results
    const stats = {
      total: Math.floor(Math.random() * 20) + 10,
      success: 0,
      failed: 0,
    };
    stats.success = stats.total - Math.floor(Math.random() * 3);
    stats.failed = stats.total - stats.success;

    setImportStats(stats);
    setUploadComplete(true);
    setUploading(false);

    // Store in localStorage (in production, would call API)
    const existingEngineers = localStorage.getItem('engineers');
    let engineers = existingEngineers ? JSON.parse(existingEngineers) : [];
    
    // Add new engineers from import
    for (let i = 0; i < stats.success; i++) {
      engineers.push({
        id: `ENG-${Date.now()}-${i}`,
        name: `Engineer ${i + 1}`,
        phone: `+91-98765432${String(i).padStart(2, '0')}`,
        email: `engineer${i + 1}@company.com`,
        location: ['Mumbai', 'Delhi', 'Bangalore', 'Pune'][i % 4],
        specializations: ['MRI Scanner', 'CT Scanner', 'X-Ray', 'Ultrasound'][i % 4],
      });
    }

    localStorage.setItem('engineers', JSON.stringify(engineers));

    // Auto redirect after 2 seconds
    setTimeout(() => {
      router.push('/dashboard');
    }, 2000);
  };

  const handleRemoveFile = () => {
    setFile(null);
    setUploadComplete(false);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-4xl mx-auto">
        <Button
          variant="ghost"
          onClick={() => router.push('/dashboard')}
          className="mb-6"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Dashboard
        </Button>

        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Import Engineers</h1>
          <p className="text-gray-600 mt-2">
            Upload a CSV file to bulk import service engineers
          </p>
        </div>

        {!uploadComplete ? (
          <>
            <Card className="mb-6">
              <CardHeader>
                <CardTitle>CSV Format Requirements</CardTitle>
                <CardDescription>
                  Your CSV file should include the following columns
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="bg-gray-50 p-4 rounded-lg font-mono text-sm">
                  <div className="text-gray-700 mb-2">
                    <strong>Required columns:</strong>
                  </div>
                  <div className="text-gray-600">
                    name, phone, email, location, specializations
                  </div>
                  <div className="mt-4 text-gray-700">
                    <strong>Example:</strong>
                  </div>
                  <div className="text-gray-600 text-xs mt-2">
                    name,phone,email,location,specializations<br />
                    Raj Kumar,+91-9876543210,raj@company.com,Mumbai,MRI Scanner | CT Scanner<br />
                    Priya Shah,+91-9876543211,priya@company.com,Delhi,Ultrasound | ECG
                  </div>
                </div>
                <div className="mt-4">
                  <a
                    href="#"
                    className="text-blue-600 hover:underline text-sm"
                    onClick={(e) => {
                      e.preventDefault();
                      alert('Sample CSV download coming soon!');
                    }}
                  >
                    Download sample template
                  </a>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Upload CSV File</CardTitle>
                <CardDescription>
                  Drag and drop your file or click to browse
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div
                  className={`border-2 border-dashed rounded-lg p-12 text-center transition-colors ${
                    dragActive
                      ? 'border-blue-500 bg-blue-50'
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                  onDragEnter={handleDrag}
                  onDragLeave={handleDrag}
                  onDragOver={handleDrag}
                  onDrop={handleDrop}
                >
                  {file ? (
                    <div className="space-y-4">
                      <div className="flex items-center justify-center gap-3 text-green-600">
                        <FileText className="h-8 w-8" />
                        <div className="text-left">
                          <p className="font-medium">{file.name}</p>
                          <p className="text-sm text-gray-500">
                            {(file.size / 1024).toFixed(2)} KB
                          </p>
                        </div>
                        <button
                          onClick={handleRemoveFile}
                          className="ml-4 p-1 hover:bg-gray-100 rounded"
                        >
                          <X className="h-5 w-5 text-gray-500" />
                        </button>
                      </div>
                      <Button
                        onClick={handleUpload}
                        disabled={uploading}
                        className="w-full"
                      >
                        {uploading ? (
                          <>
                            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                            Importing Engineers...
                          </>
                        ) : (
                          <>
                            <Upload className="mr-2 h-4 w-4" />
                            Import Engineers
                          </>
                        )}
                      </Button>
                    </div>
                  ) : (
                    <>
                      <Upload className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                      <p className="text-lg font-medium text-gray-900 mb-2">
                        Drop your CSV file here
                      </p>
                      <p className="text-sm text-gray-500 mb-4">
                        or click to browse from your computer
                      </p>
                      <input
                        type="file"
                        accept=".csv"
                        onChange={handleFileChange}
                        className="hidden"
                        id="file-upload"
                      />
                      <label htmlFor="file-upload">
                        <Button variant="outline" asChild>
                          <span>Choose File</span>
                        </Button>
                      </label>
                    </>
                  )}
                </div>
              </CardContent>
            </Card>
          </>
        ) : (
          <Card>
            <CardContent className="pt-6">
              <div className="text-center py-8">
                <div className="mb-4">
                  <CheckCircle2 className="h-16 w-16 text-green-500 mx-auto" />
                </div>
                <h3 className="text-2xl font-bold text-gray-900 mb-2">
                  Import Successful!
                </h3>
                <p className="text-gray-600 mb-6">
                  Your engineers have been imported successfully.
                </p>

                <div className="grid grid-cols-3 gap-4 max-w-md mx-auto mb-8">
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="text-3xl font-bold text-gray-900">
                      {importStats.total}
                    </div>
                    <div className="text-sm text-gray-600">Total</div>
                  </div>
                  <div className="bg-green-50 p-4 rounded-lg">
                    <div className="text-3xl font-bold text-green-600">
                      {importStats.success}
                    </div>
                    <div className="text-sm text-gray-600">Success</div>
                  </div>
                  <div className="bg-red-50 p-4 rounded-lg">
                    <div className="text-3xl font-bold text-red-600">
                      {importStats.failed}
                    </div>
                    <div className="text-sm text-gray-600">Failed</div>
                  </div>
                </div>

                <Alert className="max-w-md mx-auto mb-6">
                  <p className="text-sm">
                    Redirecting to dashboard in 2 seconds...
                  </p>
                </Alert>

                <Button onClick={() => router.push('/dashboard')}>
                  Go to Dashboard
                </Button>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
