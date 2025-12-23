'use client';

import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import CSVUploader from '../CSVUploader';
import { ChevronRight, ChevronLeft } from 'lucide-react';

interface OrganizationsUploadStepProps {
  data?: any;
  onNext: (data: any) => void;
  onBack?: () => void;
  isFirstStep?: boolean;
  isLastStep?: boolean;
}

export default function OrganizationsUploadStep({
  data,
  onNext,
  onBack,
  isFirstStep,
  isLastStep
}: OrganizationsUploadStepProps) {
  const [uploadResult, setUploadResult] = useState<any>(data?.uploadResult || null);

  const handleUpload = async (file: File, dryRun: boolean) => {
    const formData = new FormData();
    formData.append('csv_file', file);
    formData.append('created_by', 'onboarding-wizard');
    formData.append('dry_run', dryRun.toString());
    formData.append('update_mode', 'false');

    const response = await fetch('http://localhost:8081/api/v1/organizations/import', {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      throw new Error('Upload failed');
    }

    const result = await response.json();
    
    if (!dryRun && result.failure_count === 0) {
      setUploadResult(result);
    }

    return result;
  };

  const handleContinue = () => {
    if (uploadResult) {
      onNext({ uploadResult });
    }
  };

  const handleSkip = () => {
    onNext({ skipped: true });
  };

  return (
    <div className="space-y-6">
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h4 className="font-medium text-blue-900 mb-2">Bulk Import Organizations</h4>
        <p className="text-sm text-blue-700">
          Upload a CSV file to import multiple organizations at once. This is useful for:
        </p>
        <ul className="text-sm text-blue-700 list-disc list-inside mt-2 space-y-1">
          <li>Adding multiple manufacturers</li>
          <li>Importing your supplier network</li>
          <li>Setting up hospital and clinic partners</li>
          <li>Bulk importing distributors and dealers</li>
        </ul>
      </div>

      <CSVUploader
        onUpload={handleUpload}
        templateUrl="/templates/organizations-import-template.csv"
        title="Upload Organizations CSV"
        description="Upload a CSV file with organization details (name, type, contact info, address)"
      />

      {/* Navigation Buttons */}
      <div className="flex justify-between pt-6 border-t">
        <div>
          {!isFirstStep && onBack && (
            <Button variant="outline" onClick={onBack} className="flex items-center">
              <ChevronLeft className="w-4 h-4 mr-2" />
              Back
            </Button>
          )}
        </div>

        <div className="flex gap-3">
          <Button variant="ghost" onClick={handleSkip}>
            Skip for now
          </Button>

          {uploadResult && uploadResult.failure_count === 0 && (
            <Button onClick={handleContinue} className="flex items-center">
              Continue
              <ChevronRight className="w-4 h-4 ml-2" />
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
