"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import AdminRoute from '@/components/auth/AdminRoute';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081';

interface ImportResult {
  total: number;
  successful: number;
  failed: number;
  errors?: Array<{ row: number; error: string }>;
  organizations?: Array<{ id: string; name: string }>;
}

export default function ImportManufacturersPage() {
  const router = useRouter();
  const [file, setFile] = useState<File | null>(null);
  const [isDryRun, setIsDryRun] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [result, setResult] = useState<ImportResult | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
      setError('');
      setResult(null);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!file) {
      setError('Please select a CSV file');
      return;
    }

    setError('');
    setIsLoading(true);

    try {
      const formData = new FormData();
      formData.append('csv_file', file);
      formData.append('dry_run', isDryRun.toString());

      const response = await fetch(`${API_BASE_URL}/api/v1/organizations/import`, {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      if (response.ok) {
        const data = await response.json();
        setResult(data);
        
        // If actual import (not dry run) was successful, redirect after a delay
        if (!isDryRun && data.successful > 0) {
          setTimeout(() => {
            router.push('/manufacturers');
          }, 3000);
        }
      } else {
        const data = await response.json();
        setError(data.error?.message || 'Failed to import organizations');
      }
    } catch (err) {
      setError('Network error. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const downloadTemplate = () => {
    const csvContent = `name,type,email,phone,address,city,state,country,postal_code,website,contact_person,status
Acme Medical Equipment,manufacturer,contact@acme.com,+1-555-0100,"123 Main St",Boston,MA,USA,02101,https://acme.com,John Doe,active
City Hospital,hospital,admin@cityhospital.com,+1-555-0200,"456 Healthcare Ave",New York,NY,USA,10001,https://cityhospital.com,Jane Smith,active
MedSupply Co,supplier,sales@medsupply.com,+1-555-0300,"789 Supply Rd",Chicago,IL,USA,60601,https://medsupply.com,Bob Johnson,active`;

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'organizations_template.csv';
    a.click();
    window.URL.revokeObjectURL(url);
  };

  return (
    <AdminRoute>
      <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => router.back()}
            className="text-blue-600 hover:text-blue-700 font-medium mb-4 flex items-center"
          >
            ← Back
          </button>
          <h1 className="text-3xl font-bold text-gray-900">Import Manufacturers</h1>
          <p className="mt-2 text-gray-600">Bulk import organizations from CSV file</p>
        </div>

        {/* Instructions */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-6 mb-6">
          <h2 className="text-lg font-semibold text-blue-900 mb-3">📋 How to Import</h2>
          <ol className="list-decimal list-inside space-y-2 text-blue-800">
            <li>Download the CSV template below</li>
            <li>Fill in your organization data (keep the header row)</li>
            <li>Upload the completed CSV file</li>
            <li>Run a "Dry Run" first to validate your data</li>
            <li>If validation passes, uncheck "Dry Run" and import</li>
          </ol>
          <button
            onClick={downloadTemplate}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 font-medium"
          >
            📥 Download CSV Template
          </button>
        </div>

        {/* CSV Format Info */}
        <div className="bg-white shadow-md rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">CSV Format</h2>
          <div className="overflow-x-auto">
            <table className="min-w-full text-sm">
              <thead>
                <tr className="bg-gray-50">
                  <th className="px-3 py-2 text-left font-medium text-gray-700">Column</th>
                  <th className="px-3 py-2 text-left font-medium text-gray-700">Required</th>
                  <th className="px-3 py-2 text-left font-medium text-gray-700">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">name</td>
                  <td className="px-3 py-2 text-green-600">Yes</td>
                  <td className="px-3 py-2 text-gray-600">Organization name</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">type</td>
                  <td className="px-3 py-2 text-green-600">Yes</td>
                  <td className="px-3 py-2 text-gray-600">manufacturer, hospital, clinic, supplier, etc.</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">email</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">Contact email</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">phone</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">Contact phone</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">address</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">Street address</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">city</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">City</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">state</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">State/Province</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">country</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">Country</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">postal_code</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">Postal/ZIP code</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">website</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">Website URL</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">contact_person</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">Contact person name</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-mono text-xs">status</td>
                  <td className="px-3 py-2 text-gray-500">No</td>
                  <td className="px-3 py-2 text-gray-600">active or inactive (default: active)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        {/* Upload Form */}
        <form onSubmit={handleSubmit} className="bg-white shadow-md rounded-lg p-6 space-y-6">
          {error && (
            <div className="rounded-md bg-red-50 p-4">
              <div className="text-sm text-red-800">{error}</div>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Upload CSV File
            </label>
            <input
              type="file"
              accept=".csv"
              onChange={handleFileChange}
              className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-medium file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            />
            {file && (
              <p className="mt-2 text-sm text-gray-600">
                Selected: <span className="font-medium">{file.name}</span> ({(file.size / 1024).toFixed(1)} KB)
              </p>
            )}
          </div>

          <div className="flex items-center">
            <input
              type="checkbox"
              id="dryRun"
              checked={isDryRun}
              onChange={(e) => setIsDryRun(e.target.checked)}
              className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label htmlFor="dryRun" className="ml-2 block text-sm text-gray-900">
              Dry Run (validate only, don't import)
            </label>
          </div>

          <div className="flex justify-end space-x-4 pt-4 border-t">
            <button
              type="button"
              onClick={() => router.back()}
              className="px-6 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 font-medium"
              disabled={isLoading}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading || !file}
              className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? 'Processing...' : isDryRun ? 'Validate' : 'Import'}
            </button>
          </div>
        </form>

        {/* Results */}
        {result && (
          <div className="mt-6 bg-white shadow-md rounded-lg p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              {isDryRun ? 'Validation Results' : 'Import Results'}
            </h2>
            
            <div className="grid grid-cols-3 gap-4 mb-6">
              <div className="bg-blue-50 p-4 rounded-lg">
                <div className="text-2xl font-bold text-blue-600">{result.total}</div>
                <div className="text-sm text-blue-800">Total Rows</div>
              </div>
              <div className="bg-green-50 p-4 rounded-lg">
                <div className="text-2xl font-bold text-green-600">{result.successful}</div>
                <div className="text-sm text-green-800">Successful</div>
              </div>
              <div className="bg-red-50 p-4 rounded-lg">
                <div className="text-2xl font-bold text-red-600">{result.failed}</div>
                <div className="text-sm text-red-800">Failed</div>
              </div>
            </div>

            {result.errors && result.errors.length > 0 && (
              <div className="mb-4">
                <h3 className="text-lg font-semibold text-red-900 mb-2">Errors:</h3>
                <div className="bg-red-50 border border-red-200 rounded-md p-4 max-h-64 overflow-y-auto">
                  {result.errors.map((err, idx) => (
                    <div key={idx} className="text-sm text-red-800 mb-2">
                      <span className="font-medium">Row {err.row}:</span> {err.error}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {!isDryRun && result.successful > 0 && (
              <div className="rounded-md bg-green-50 p-4">
                <div className="text-sm text-green-800">
                  ✅ Successfully imported {result.successful} organization(s). Redirecting to manufacturers list...
                </div>
              </div>
            )}

            {isDryRun && result.successful === result.total && (
              <div className="rounded-md bg-green-50 p-4">
                <div className="text-sm text-green-800">
                  ✅ All rows validated successfully! You can now uncheck "Dry Run" and import.
                </div>
              </div>
            )}
          </div>
        )}
      </div>
      </div>
    </AdminRoute>
  );
}
