'use client';

import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import CSVUploader from '../CSVUploader';
import { ChevronRight, ChevronLeft, Download } from 'lucide-react';

interface EquipmentUploadStepProps {
  data?: any;
  onNext: (data: any) => void;
  onBack?: () => void;
  isFirstStep?: boolean;
  isLastStep?: boolean;
}

const INDUSTRY_TEMPLATES = [
  {
    id: 'radiology',
    name: 'Radiology',
    icon: '🔬',
    description: 'MRI, CT, X-Ray, Ultrasound, PACS',
    template: '/templates/equipment-catalog-radiology-template.csv',
    count: 8
  },
  {
    id: 'cardiology',
    name: 'Cardiology',
    icon: '❤️',
    description: 'Cath Lab, Echo, ECG, Holter, Stress Test',
    template: '/templates/equipment-catalog-cardiology-template.csv',
    count: 8
  },
  {
    id: 'surgical',
    name: 'Surgical',
    icon: '🏥',
    description: 'OR Table, Anesthesia, Laparoscopy, Surgical Robot',
    template: '/templates/equipment-catalog-surgical-template.csv',
    count: 8
  },
  {
    id: 'icu',
    name: 'ICU',
    icon: '🛏️',
    description: 'Ventilator, Patient Monitor, Infusion Pumps',
    template: '/templates/equipment-catalog-icu-template.csv',
    count: 8
  },
  {
    id: 'laboratory',
    name: 'Laboratory',
    icon: '🧪',
    description: 'Hematology, Chemistry, Microbiology, PCR',
    template: '/templates/equipment-catalog-laboratory-template.csv',
    count: 8
  },
];

export default function EquipmentUploadStep({
  data,
  onNext,
  onBack,
  isFirstStep,
  isLastStep
}: EquipmentUploadStepProps) {
  const [selectedIndustry, setSelectedIndustry] = useState<string | null>(
    data?.selectedIndustry || null
  );
  const [uploadResult, setUploadResult] = useState<any>(data?.uploadResult || null);

  const handleUpload = async (file: File, dryRun: boolean) => {
    const formData = new FormData();
    formData.append('csv_file', file);
    formData.append('created_by', 'onboarding-wizard');
    formData.append('dry_run', dryRun.toString());
    formData.append('update_mode', 'false');

    const response = await fetch('http://localhost:8081/api/v1/equipment/catalog/import', {
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
      onNext({ selectedIndustry, uploadResult });
    }
  };

  const handleSkip = () => {
    onNext({ skipped: true });
  };

  const selectedTemplate = INDUSTRY_TEMPLATES.find(t => t.id === selectedIndustry);

  return (
    <div className="space-y-6">
      {/* Industry Selection */}
      {!uploadResult && (
        <>
          <div className="bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-lg p-6">
            <h4 className="font-semibold text-blue-900 mb-2">Choose Your Industry</h4>
            <p className="text-sm text-blue-700 mb-4">
              Select an industry template to get started quickly with pre-configured equipment, or upload your own custom CSV.
            </p>

            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
              {INDUSTRY_TEMPLATES.map((industry) => (
                <button
                  key={industry.id}
                  onClick={() => setSelectedIndustry(industry.id)}
                  className={`
                    p-4 rounded-lg border-2 transition-all text-center
                    ${selectedIndustry === industry.id
                      ? 'border-blue-500 bg-blue-50 shadow-md'
                      : 'border-gray-200 bg-white hover:border-blue-300 hover:shadow'
                    }
                  `}
                >
                  <div className="text-3xl mb-2">{industry.icon}</div>
                  <div className="font-semibold text-sm text-gray-900">{industry.name}</div>
                  <div className="text-xs text-gray-500 mt-1">{industry.count} items</div>
                </button>
              ))}
            </div>
          </div>

          {selectedTemplate && (
            <Card className="p-4 bg-green-50 border-green-200">
              <div className="flex items-center justify-between">
                <div>
                  <h5 className="font-semibold text-green-900 mb-1">
                    {selectedTemplate.icon} {selectedTemplate.name} Template Selected
                  </h5>
                  <p className="text-sm text-green-700">{selectedTemplate.description}</p>
                </div>
                <a href={selectedTemplate.template} download>
                  <Button variant="outline" size="sm" className="flex items-center">
                    <Download className="w-4 h-4 mr-2" />
                    Download
                  </Button>
                </a>
              </div>
            </Card>
          )}
        </>
      )}

      {/* Upload Section */}
      <div className="bg-white border rounded-lg p-6">
        <h4 className="font-semibold text-gray-900 mb-2">Upload Equipment Catalog</h4>
        <p className="text-sm text-gray-600 mb-4">
          Upload a CSV file with your equipment catalog. You can use one of our industry templates or create your own.
        </p>

        <CSVUploader
          onUpload={handleUpload}
          templateUrl={selectedTemplate?.template}
          title="Upload Equipment Catalog CSV"
          description="CSV with product details (code, name, manufacturer, model, category, price, etc.)"
        />
      </div>

      {/* Benefits */}
      {!uploadResult && (
        <Card className="p-6 bg-gradient-to-br from-purple-50 to-pink-50 border-purple-200">
          <h4 className="font-semibold text-purple-900 mb-3">Why Import Equipment Catalog?</h4>
          <div className="grid grid-cols-2 gap-4">
            <div className="flex items-start">
              <div className="w-2 h-2 rounded-full bg-purple-600 mt-2 mr-3"></div>
              <div>
                <p className="font-medium text-purple-900 text-sm">Quick Setup</p>
                <p className="text-xs text-purple-700">Pre-populate your equipment database</p>
              </div>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 rounded-full bg-purple-600 mt-2 mr-3"></div>
              <div>
                <p className="font-medium text-purple-900 text-sm">Standardized Data</p>
                <p className="text-xs text-purple-700">Consistent product information</p>
              </div>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 rounded-full bg-purple-600 mt-2 mr-3"></div>
              <div>
                <p className="font-medium text-purple-900 text-sm">Easy Maintenance</p>
                <p className="text-xs text-purple-700">Track service intervals and lifecycle</p>
              </div>
            </div>
            <div className="flex items-start">
              <div className="w-2 h-2 rounded-full bg-purple-600 mt-2 mr-3"></div>
              <div>
                <p className="font-medium text-purple-900 text-sm">Parts Management</p>
                <p className="text-xs text-purple-700">Link spare parts to equipment</p>
              </div>
            </div>
          </div>
        </Card>
      )}

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
