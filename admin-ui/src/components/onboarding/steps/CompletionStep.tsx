'use client';

import React from 'react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { 
  CheckCircle, 
  Building2, 
  Users, 
  Package, 
  ArrowRight,
  Download 
} from 'lucide-react';
import { useRouter } from 'next/navigation';

interface CompletionStepProps {
  data?: any;
  onNext: (data: any) => void;
  onBack?: () => void;
  isFirstStep?: boolean;
  isLastStep?: boolean;
}

export default function CompletionStep({
  data,
  onNext,
  onBack,
  isFirstStep,
  isLastStep
}: CompletionStepProps) {
  const router = useRouter();

  const handleGoToDashboard = () => {
    router.push('/dashboard');
  };

  const stats = [
    {
      icon: Building2,
      label: 'Organizations',
      value: data?.organizationsUpload?.uploadResult?.success_count || 0,
      color: 'blue'
    },
    {
      icon: Package,
      label: 'Equipment',
      value: data?.equipmentUpload?.uploadResult?.success_count || 0,
      color: 'green'
    },
    {
      icon: Users,
      label: 'Engineers',
      value: data?.engineersUpload?.uploadResult?.success_count || 0,
      color: 'purple'
    },
  ];

  return (
    <div className="space-y-6">
      {/* Success Message */}
      <div className="text-center py-8">
        <div className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-green-100 mb-4">
          <CheckCircle className="w-12 h-12 text-green-600" />
        </div>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">
          Onboarding Complete!
        </h2>
        <p className="text-lg text-gray-600">
          Your organization is ready to start managing medical equipment.
        </p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-3 gap-4">
        {stats.map((stat, index) => {
          const Icon = stat.icon;
          const colorClasses = {
            blue: 'bg-blue-50 text-blue-600 border-blue-200',
            green: 'bg-green-50 text-green-600 border-green-200',
            purple: 'bg-purple-50 text-purple-600 border-purple-200'
          }[stat.color];

          return (
            <Card key={index} className={`p-6 border-2 ${colorClasses}`}>
              <div className="flex items-center justify-between mb-2">
                <Icon className="w-8 h-8" />
                <span className="text-3xl font-bold">{stat.value}</span>
              </div>
              <p className="font-medium">{stat.label}</p>
            </Card>
          );
        })}
      </div>

      {/* Next Steps */}
      <Card className="p-6 bg-gradient-to-br from-blue-50 to-indigo-50">
        <h3 className="text-lg font-semibold mb-4">What's Next?</h3>
        <div className="space-y-3">
          <div className="flex items-start">
            <div className="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-sm font-bold mr-3 flex-shrink-0 mt-0.5">
              1
            </div>
            <div>
              <p className="font-medium text-gray-900">Explore Your Dashboard</p>
              <p className="text-sm text-gray-600">
                View equipment, track service requests, and manage your organization
              </p>
            </div>
          </div>

          <div className="flex items-start">
            <div className="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-sm font-bold mr-3 flex-shrink-0 mt-0.5">
              2
            </div>
            <div>
              <p className="font-medium text-gray-900">Complete Your Profile</p>
              <p className="text-sm text-gray-600">
                Add more details about your organization and team members
              </p>
            </div>
          </div>

          <div className="flex items-start">
            <div className="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-sm font-bold mr-3 flex-shrink-0 mt-0.5">
              3
            </div>
            <div>
              <p className="font-medium text-gray-900">Import More Data</p>
              <p className="text-sm text-gray-600">
                You can always import more equipment, parts, or engineers later
              </p>
            </div>
          </div>

          <div className="flex items-start">
            <div className="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-sm font-bold mr-3 flex-shrink-0 mt-0.5">
              4
            </div>
            <div>
              <p className="font-medium text-gray-900">Generate QR Codes</p>
              <p className="text-sm text-gray-600">
                Create QR codes for your equipment for easy tracking and maintenance
              </p>
            </div>
          </div>
        </div>
      </Card>

      {/* Resources */}
      <Card className="p-6">
        <h3 className="text-lg font-semibold mb-4">Helpful Resources</h3>
        <div className="grid grid-cols-2 gap-3">
          <Button variant="outline" className="justify-start">
            <Download className="w-4 h-4 mr-2" />
            Download User Guide
          </Button>
          <Button variant="outline" className="justify-start">
            <Package className="w-4 h-4 mr-2" />
            View CSV Templates
          </Button>
        </div>
      </Card>

      {/* Action Buttons */}
      <div className="flex justify-center pt-6">
        <Button
          onClick={handleGoToDashboard}
          size="lg"
          className="px-8 flex items-center"
        >
          Go to Dashboard
          <ArrowRight className="w-5 h-5 ml-2" />
        </Button>
      </div>
    </div>
  );
}
