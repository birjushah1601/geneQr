'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { CheckCircle, User, Mail, Phone, MapPin, Wrench, Plus, Trash2, SkipForward } from 'lucide-react';

interface Engineer {
  id: string;
  name: string;
  phone: string;
  email: string;
  location: string;
  specializations: string;
}

export default function EngineersSetup() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [engineers, setEngineers] = useState<Engineer[]>([
    { id: '1', name: '', phone: '', email: '', location: '', specializations: '' }
  ]);

  // Get manufacturer from localStorage
  const manufacturer = typeof window !== 'undefined' 
    ? JSON.parse(localStorage.getItem('current_manufacturer') || '{}')
    : {};

  const handleAddEngineer = () => {
    setEngineers([
      ...engineers,
      { id: Date.now().toString(), name: '', phone: '', email: '', location: '', specializations: '' }
    ]);
  };

  const handleRemoveEngineer = (id: string) => {
    if (engineers.length > 1) {
      setEngineers(engineers.filter(e => e.id !== id));
    }
  };

  const handleChange = (id: string, field: keyof Engineer, value: string) => {
    setEngineers(engineers.map(e => 
      e.id === id ? { ...e, [field]: value } : e
    ));
  };

  const handleSubmit = async () => {
    // Filter out empty engineers
    const validEngineers = engineers.filter(e => e.name && e.phone && e.email);
    
    if (validEngineers.length === 0) {
      alert('Please add at least one engineer with name, phone, and email');
      return;
    }

    setLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1500));
      
      // Store engineers
      localStorage.setItem('engineers', JSON.stringify(validEngineers));
      localStorage.setItem('onboarding_complete', 'true');
      
      // Redirect to dashboard
      router.push('/dashboard');
    } catch (error) {
      alert('Failed to save engineers');
    } finally {
      setLoading(false);
    }
  };

  const handleSkip = () => {
    localStorage.setItem('onboarding_complete', 'true');
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
          
          <div className="flex-1 h-1 mx-4 bg-green-600"></div>
          
          <div className="flex items-center opacity-50">
            <div className="w-10 h-10 rounded-full bg-green-600 text-white flex items-center justify-center">
              <CheckCircle className="w-6 h-6" />
            </div>
            <div className="ml-3">
              <p className="font-semibold">Equipment</p>
              <p className="text-sm text-gray-500">Completed</p>
            </div>
          </div>
          
          <div className="flex-1 h-1 mx-4 bg-blue-600"></div>
          
          <div className="flex items-center">
            <div className="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
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
        <h2 className="text-2xl font-bold mb-2">Add Service Engineers</h2>
        <p className="text-gray-600">
          Add your field engineers for <strong>{manufacturer.name || 'your company'}</strong>
        </p>
        <p className="text-sm text-gray-500 mt-2">
          You can also skip this step and add engineers later from the dashboard
        </p>
      </div>

      <Alert className="mb-6">
        <AlertDescription>
          Add engineers manually below, or skip to upload a CSV file from the dashboard later.
        </AlertDescription>
      </Alert>

      {/* Engineers Form */}
      <div className="space-y-6 max-h-96 overflow-y-auto mb-6">
        {engineers.map((engineer, index) => (
          <div key={engineer.id} className="border border-gray-200 rounded-lg p-4 relative">
            <div className="absolute top-4 right-4">
              {engineers.length > 1 && (
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleRemoveEngineer(engineer.id)}
                  className="text-red-600 hover:text-red-700 hover:bg-red-50"
                >
                  <Trash2 className="w-4 h-4" />
                </Button>
              )}
            </div>

            <h3 className="font-semibold mb-4 flex items-center gap-2">
              <User className="w-5 h-5" />
              Engineer #{index + 1}
            </h3>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label htmlFor={`name-${engineer.id}`}>Full Name *</Label>
                <Input
                  id={`name-${engineer.id}`}
                  placeholder="John Doe"
                  value={engineer.name}
                  onChange={(e) => handleChange(engineer.id, 'name', e.target.value)}
                  className="mt-2"
                />
              </div>

              <div>
                <Label htmlFor={`phone-${engineer.id}`}>Phone Number *</Label>
                <Input
                  id={`phone-${engineer.id}`}
                  type="tel"
                  placeholder="+91-9876543210"
                  value={engineer.phone}
                  onChange={(e) => handleChange(engineer.id, 'phone', e.target.value)}
                  className="mt-2"
                />
              </div>

              <div>
                <Label htmlFor={`email-${engineer.id}`}>Email *</Label>
                <Input
                  id={`email-${engineer.id}`}
                  type="email"
                  placeholder="engineer@company.com"
                  value={engineer.email}
                  onChange={(e) => handleChange(engineer.id, 'email', e.target.value)}
                  className="mt-2"
                />
              </div>

              <div>
                <Label htmlFor={`location-${engineer.id}`}>Location</Label>
                <Input
                  id={`location-${engineer.id}`}
                  placeholder="Mumbai, Maharashtra"
                  value={engineer.location}
                  onChange={(e) => handleChange(engineer.id, 'location', e.target.value)}
                  className="mt-2"
                />
              </div>

              <div className="md:col-span-2">
                <Label htmlFor={`specializations-${engineer.id}`}>
                  Specializations (comma-separated)
                </Label>
                <Input
                  id={`specializations-${engineer.id}`}
                  placeholder="MRI Scanner, CT Scanner, X-Ray"
                  value={engineer.specializations}
                  onChange={(e) => handleChange(engineer.id, 'specializations', e.target.value)}
                  className="mt-2"
                />
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Add Engineer Button */}
      <Button
        variant="outline"
        onClick={handleAddEngineer}
        className="w-full mb-6"
      >
        <Plus className="w-4 h-4 mr-2" />
        Add Another Engineer
      </Button>

      {/* Action Buttons */}
      <div className="flex justify-between gap-4 pt-6 border-t">
        <Button 
          variant="outline" 
          onClick={handleSkip}
          disabled={loading}
        >
          <SkipForward className="w-4 h-4 mr-2" />
          Skip for Now
        </Button>
        
        <Button 
          onClick={handleSubmit} 
          disabled={loading}
          className="px-8"
        >
          {loading ? 'Saving...' : 'Complete Setup'}
        </Button>
      </div>
    </Card>
  );
}
