'use client';

import React, { useState } from 'react';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { CheckCircle2, Circle, ChevronRight, ChevronLeft } from 'lucide-react';

export interface WizardStep {
  id: string;
  title: string;
  description: string;
  component: React.ComponentType<any>;
}

interface OnboardingWizardProps {
  steps: WizardStep[];
  onComplete?: (data: any) => void;
  onCancel?: () => void;
}

export default function OnboardingWizard({ 
  steps, 
  onComplete, 
  onCancel 
}: OnboardingWizardProps) {
  const [currentStepIndex, setCurrentStepIndex] = useState(0);
  const [completedSteps, setCompletedSteps] = useState<Set<number>>(new Set());
  const [stepData, setStepData] = useState<Record<string, any>>({});
  const [error, setError] = useState('');

  const currentStep = steps[currentStepIndex];
  const isLastStep = currentStepIndex === steps.length - 1;
  const isFirstStep = currentStepIndex === 0;

  const handleNext = (data?: any) => {
    // Save current step data
    if (data) {
      setStepData(prev => ({
        ...prev,
        [currentStep.id]: data
      }));
    }

    // Mark current step as completed
    setCompletedSteps(prev => new Set([...prev, currentStepIndex]));

    if (isLastStep) {
      // Complete wizard
      const allData = {
        ...stepData,
        [currentStep.id]: data
      };
      onComplete?.(allData);
    } else {
      // Move to next step
      setCurrentStepIndex(prev => prev + 1);
      setError('');
    }
  };

  const handleBack = () => {
    if (!isFirstStep) {
      setCurrentStepIndex(prev => prev - 1);
      setError('');
    }
  };

  const handleStepClick = (index: number) => {
    // Allow navigation to completed steps or current step
    if (index <= currentStepIndex) {
      setCurrentStepIndex(index);
      setError('');
    }
  };

  return (
    <div className="w-full max-w-6xl mx-auto">
      {/* Step Indicator */}
      <Card className="p-6 mb-6">
        <div className="flex items-center justify-between">
          {steps.map((step, index) => {
            const isCompleted = completedSteps.has(index);
            const isCurrent = index === currentStepIndex;
            const isAccessible = index <= currentStepIndex;

            return (
              <React.Fragment key={step.id}>
                <button
                  onClick={() => handleStepClick(index)}
                  disabled={!isAccessible}
                  className={`flex items-center ${isAccessible ? 'cursor-pointer' : 'cursor-not-allowed'}`}
                >
                  <div className="flex items-center">
                    {isCompleted ? (
                      <div className="w-10 h-10 rounded-full bg-green-600 text-white flex items-center justify-center">
                        <CheckCircle2 className="w-6 h-6" />
                      </div>
                    ) : isCurrent ? (
                      <div className="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
                        {index + 1}
                      </div>
                    ) : (
                      <div className={`w-10 h-10 rounded-full ${isAccessible ? 'bg-gray-300' : 'bg-gray-200'} text-gray-600 flex items-center justify-center font-bold`}>
                        {index + 1}
                      </div>
                    )}
                    
                    <div className="ml-3 text-left">
                      <p className={`font-semibold ${isCurrent ? 'text-blue-600' : isCompleted ? 'text-green-600' : 'text-gray-600'}`}>
                        {step.title}
                      </p>
                      <p className="text-sm text-gray-500">{step.description}</p>
                    </div>
                  </div>
                </button>

                {index < steps.length - 1 && (
                  <div className={`flex-1 h-1 mx-4 ${isCompleted ? 'bg-green-600' : 'bg-gray-200'}`}></div>
                )}
              </React.Fragment>
            );
          })}
        </div>
      </Card>

      {/* Error Alert */}
      {error && (
        <Alert variant="destructive" className="mb-6">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Step Content */}
      <Card className="p-8">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">
            {currentStep.title}
          </h2>
          <p className="text-gray-600">{currentStep.description}</p>
        </div>

        {/* Render current step component */}
        <div className="mb-8">
          <currentStep.component
            data={stepData[currentStep.id]}
            onNext={handleNext}
            onBack={handleBack}
            isFirstStep={isFirstStep}
            isLastStep={isLastStep}
          />
        </div>

        {/* Navigation Buttons */}
        <div className="flex justify-between border-t pt-6">
          <div>
            {!isFirstStep && (
              <Button
                variant="outline"
                onClick={handleBack}
                className="flex items-center"
              >
                <ChevronLeft className="w-4 h-4 mr-2" />
                Back
              </Button>
            )}
          </div>

          <div className="flex gap-3">
            {onCancel && (
              <Button variant="ghost" onClick={onCancel}>
                Cancel
              </Button>
            )}
          </div>
        </div>
      </Card>

      {/* Progress Bar */}
      <div className="mt-6">
        <div className="flex items-center justify-between text-sm text-gray-600 mb-2">
          <span>Progress</span>
          <span>{Math.round(((currentStepIndex + 1) / steps.length) * 100)}%</span>
        </div>
        <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
          <div
            className="h-full bg-blue-600 transition-all duration-300"
            style={{ width: `${((currentStepIndex + 1) / steps.length) * 100}%` }}
          ></div>
        </div>
      </div>
    </div>
  );
}
