'use client';

import React from 'react';
import OnboardingWizard, { WizardStep } from '@/components/onboarding/OnboardingWizard';
import CompanyProfileStep from '@/components/onboarding/steps/CompanyProfileStep';
import OrganizationsUploadStep from '@/components/onboarding/steps/OrganizationsUploadStep';
import EquipmentUploadStep from '@/components/onboarding/steps/EquipmentUploadStep';
import CompletionStep from '@/components/onboarding/steps/CompletionStep';
import { useRouter } from 'next/navigation';

export default function OnboardingWizardPage() {
  const router = useRouter();

  const steps: WizardStep[] = [
    {
      id: 'company-profile',
      title: 'Company Profile',
      description: 'Basic information about your organization',
      component: CompanyProfileStep
    },
    {
      id: 'organizations-upload',
      title: 'Import Organizations',
      description: 'Bulk import manufacturers, suppliers, partners',
      component: OrganizationsUploadStep
    },
    {
      id: 'equipment-upload',
      title: 'Equipment Catalog',
      description: 'Import equipment by industry',
      component: EquipmentUploadStep
    },
    {
      id: 'completion',
      title: 'Complete',
      description: 'You\'re all set!',
      component: CompletionStep
    },
  ];

  const handleComplete = (data: any) => {
    console.log('Onboarding completed with data:', data);
    // Save to backend or localStorage
    localStorage.setItem('onboarding_data', JSON.stringify(data));
    localStorage.setItem('onboarding_completed', 'true');
  };

  const handleCancel = () => {
    if (confirm('Are you sure you want to cancel onboarding?')) {
      router.push('/dashboard');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            Welcome to ServQR
          </h1>
          <p className="text-lg text-gray-600">
            Let's get your organization set up in just a few steps
          </p>
        </div>

        {/* Wizard */}
        <OnboardingWizard
          steps={steps}
          onComplete={handleComplete}
          onCancel={handleCancel}
        />
      </div>
    </div>
  );
}
