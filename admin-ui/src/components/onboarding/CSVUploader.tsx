'use client';

import React, { useState, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Card } from '@/components/ui/card';
import { 
  Upload, 
  FileText, 
  CheckCircle, 
  AlertCircle, 
  Download, 
  X,
  Loader2 
} from 'lucide-react';

interface CSVUploaderProps {
  onUpload: (file: File, dryRun: boolean) => Promise<any>;
  templateUrl?: string;
  acceptedFormats?: string;
  maxSizeMB?: number;
  title?: string;
  description?: string;
}

export default function CSVUploader({
  onUpload,
  templateUrl,
  acceptedFormats = '.csv',
  maxSizeMB = 10,
  title = 'Upload CSV File',
  description = 'Drag and drop your CSV file here, or click to browse'
}: CSVUploaderProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [validating, setValidating] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState('');

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const files = Array.from(e.dataTransfer.files);
    const csvFile = files.find(f => f.name.endsWith('.csv'));

    if (csvFile) {
      validateAndSelectFile(csvFile);
    } else {
      setError('Please drop a CSV file');
    }
  }, []);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      validateAndSelectFile(file);
    }
  };

  const validateAndSelectFile = (file: File) => {
    setError('');
    setResult(null);

    // Validate file type
    if (!file.name.endsWith('.csv')) {
      setError('Please select a CSV file');
      return;
    }

    // Validate file size
    const fileSizeMB = file.size / (1024 * 1024);
    if (fileSizeMB > maxSizeMB) {
      setError(`File size must be less than ${maxSizeMB}MB`);
      return;
    }

    setSelectedFile(file);
  };

  const handleValidate = async () => {
    if (!selectedFile) return;

    setValidating(true);
    setError('');
    setResult(null);

    try {
      const response = await onUpload(selectedFile, true); // dry_run = true
      setResult(response);

      if (response.failure_count > 0) {
        setError(`Validation found ${response.failure_count} error(s)`);
      }
    } catch (err: any) {
      setError(err.message || 'Validation failed');
    } finally {
      setValidating(false);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    setUploading(true);
    setError('');

    try {
      const response = await onUpload(selectedFile, false); // dry_run = false
      setResult(response);

      if (response.failure_count === 0) {
        // Success!
      } else {
        setError(`Import completed with ${response.failure_count} error(s)`);
      }
    } catch (err: any) {
      setError(err.message || 'Upload failed');
    } finally {
      setUploading(false);
    }
  };

  const handleClear = () => {
    setSelectedFile(null);
    setResult(null);
    setError('');
  };

  return (
    <div className="space-y-4">
      {/* Template Download */}
      {templateUrl && (
        <div className="flex items-center justify-between p-4 bg-blue-50 rounded-lg border border-blue-200">
          <div className="flex items-center">
            <FileText className="w-5 h-5 text-blue-600 mr-3" />
            <div>
              <p className="font-medium text-blue-900">Need a template?</p>
              <p className="text-sm text-blue-700">Download the CSV template to get started</p>
            </div>
          </div>
          <a href={templateUrl} download>
            <Button variant="outline" size="sm" className="flex items-center">
              <Download className="w-4 h-4 mr-2" />
              Download Template
            </Button>
          </a>
        </div>
      )}

      {/* Upload Area */}
      {!selectedFile ? (
        <div
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          className={`
            border-2 border-dashed rounded-lg p-8 text-center transition-colors
            ${isDragging 
              ? 'border-blue-500 bg-blue-50' 
              : 'border-gray-300 hover:border-gray-400 bg-gray-50'
            }
          `}
        >
          <Upload className={`w-12 h-12 mx-auto mb-4 ${isDragging ? 'text-blue-500' : 'text-gray-400'}`} />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{title}</h3>
          <p className="text-gray-600 mb-4">{description}</p>
          
          <label>
            <input
              type="file"
              accept={acceptedFormats}
              onChange={handleFileSelect}
              className="hidden"
            />
            <Button variant="outline" className="cursor-pointer">
              Select File
            </Button>
          </label>

          <p className="text-sm text-gray-500 mt-4">
            Maximum file size: {maxSizeMB}MB
          </p>
        </div>
      ) : (
        <Card className="p-6">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center">
              <FileText className="w-8 h-8 text-blue-600 mr-3" />
              <div>
                <p className="font-medium text-gray-900">{selectedFile.name}</p>
                <p className="text-sm text-gray-500">
                  {(selectedFile.size / 1024).toFixed(2)} KB
                </p>
              </div>
            </div>
            <Button variant="ghost" size="sm" onClick={handleClear}>
              <X className="w-4 h-4" />
            </Button>
          </div>

          {/* Validation/Upload Buttons */}
          {!result && (
            <div className="flex gap-3">
              <Button
                onClick={handleValidate}
                disabled={validating || uploading}
                variant="outline"
                className="flex-1"
              >
                {validating ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Validating...
                  </>
                ) : (
                  <>
                    <CheckCircle className="w-4 h-4 mr-2" />
                    Validate
                  </>
                )}
              </Button>

              <Button
                onClick={handleUpload}
                disabled={validating || uploading}
                className="flex-1"
              >
                {uploading ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Uploading...
                  </>
                ) : (
                  <>
                    <Upload className="w-4 h-4 mr-2" />
                    Upload
                  </>
                )}
              </Button>
            </div>
          )}
        </Card>
      )}

      {/* Error Alert */}
      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Result */}
      {result && (
        <Card className="p-6">
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h4 className="font-semibold text-lg">
                {result.dry_run ? 'Validation Results' : 'Import Results'}
              </h4>
              {result.failure_count === 0 && (
                <CheckCircle className="w-6 h-6 text-green-600" />
              )}
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div className="p-4 bg-blue-50 rounded-lg">
                <p className="text-sm text-blue-700 font-medium">Total Rows</p>
                <p className="text-2xl font-bold text-blue-900">{result.total_rows}</p>
              </div>
              <div className="p-4 bg-green-50 rounded-lg">
                <p className="text-sm text-green-700 font-medium">Success</p>
                <p className="text-2xl font-bold text-green-900">{result.success_count}</p>
              </div>
              <div className="p-4 bg-red-50 rounded-lg">
                <p className="text-sm text-red-700 font-medium">Errors</p>
                <p className="text-2xl font-bold text-red-900">{result.failure_count}</p>
              </div>
            </div>

            {/* Errors List */}
            {result.errors && result.errors.length > 0 && (
              <div className="mt-4">
                <h5 className="font-medium text-red-900 mb-2">Errors:</h5>
                <div className="space-y-2 max-h-60 overflow-y-auto">
                  {result.errors.map((err: any, idx: number) => (
                    <div key={idx} className="p-3 bg-red-50 border border-red-200 rounded text-sm">
                      <p className="font-medium text-red-900">Row {err.row}:</p>
                      <p className="text-red-700">{err.message}</p>
                      {err.data && <p className="text-red-600 text-xs mt-1">{err.data}</p>}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Actions */}
            {result.dry_run && result.failure_count === 0 && (
              <Button onClick={handleUpload} className="w-full" disabled={uploading}>
                {uploading ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Uploading...
                  </>
                ) : (
                  <>
                    <Upload className="w-4 h-4 mr-2" />
                    Proceed with Import
                  </>
                )}
              </Button>
            )}

            {!result.dry_run && result.failure_count === 0 && (
              <Alert>
                <CheckCircle className="h-4 w-4" />
                <AlertDescription>
                  Successfully imported {result.success_count} record(s)!
                </AlertDescription>
              </Alert>
            )}
          </div>
        </Card>
      )}
    </div>
  );
}
