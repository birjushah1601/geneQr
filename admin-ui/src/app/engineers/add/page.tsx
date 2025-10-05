'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Plus, Trash2, Save } from 'lucide-react';

interface Engineer {
  id: string;
  name: string;
  phone: string;
  email: string;
  location: string;
  specializations: string;
}

export default function AddEngineerPage() {
  const router = useRouter();
  const [engineers, setEngineers] = useState<Engineer[]>([
    {
      id: '1',
      name: '',
      phone: '',
      email: '',
      location: '',
      specializations: '',
    },
  ]);

  const addEngineer = () => {
    setEngineers([
      ...engineers,
      {
        id: Date.now().toString(),
        name: '',
        phone: '',
        email: '',
        location: '',
        specializations: '',
      },
    ]);
  };

  const removeEngineer = (id: string) => {
    if (engineers.length > 1) {
      setEngineers(engineers.filter((eng) => eng.id !== id));
    }
  };

  const updateEngineer = (id: string, field: keyof Engineer, value: string) => {
    setEngineers(
      engineers.map((eng) =>
        eng.id === id ? { ...eng, [field]: value } : eng
      )
    );
  };

  const handleSave = () => {
    // Validate required fields
    const invalidEngineers = engineers.filter(
      (eng) => !eng.name || !eng.phone || !eng.email
    );

    if (invalidEngineers.length > 0) {
      alert('Please fill in name, phone, and email for all engineers');
      return;
    }

    // Save to localStorage
    const existingEngineers = localStorage.getItem('engineers');
    let allEngineers = existingEngineers ? JSON.parse(existingEngineers) : [];

    const newEngineers = engineers.map((eng, idx) => ({
      ...eng,
      id: `ENG-${Date.now()}-${idx}`,
    }));

    allEngineers = [...allEngineers, ...newEngineers];
    localStorage.setItem('engineers', JSON.stringify(allEngineers));

    alert(`Successfully added ${engineers.length} engineer(s)!`);
    router.push('/engineers');
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <Button
          variant="ghost"
          onClick={() => router.push('/engineers')}
          className="mb-6"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Engineers
        </Button>

        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Add Engineers</h1>
          <p className="text-gray-600 mt-2">
            Manually add service engineers to your team
          </p>
        </div>

        {/* Engineers Forms */}
        <div className="space-y-6">
          {engineers.map((engineer, index) => (
            <Card key={engineer.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Engineer {index + 1}</CardTitle>
                    <CardDescription>
                      Fill in the engineer's details
                    </CardDescription>
                  </div>
                  {engineers.length > 1 && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => removeEngineer(engineer.id)}
                    >
                      <Trash2 className="h-4 w-4 text-red-600" />
                    </Button>
                  )}
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor={`name-${engineer.id}`}>
                      Name <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id={`name-${engineer.id}`}
                      placeholder="Full name"
                      value={engineer.name}
                      onChange={(e) =>
                        updateEngineer(engineer.id, 'name', e.target.value)
                      }
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor={`phone-${engineer.id}`}>
                      Phone <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id={`phone-${engineer.id}`}
                      type="tel"
                      placeholder="+91-9876543210"
                      value={engineer.phone}
                      onChange={(e) =>
                        updateEngineer(engineer.id, 'phone', e.target.value)
                      }
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor={`email-${engineer.id}`}>
                      Email <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id={`email-${engineer.id}`}
                      type="email"
                      placeholder="engineer@company.com"
                      value={engineer.email}
                      onChange={(e) =>
                        updateEngineer(engineer.id, 'email', e.target.value)
                      }
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor={`location-${engineer.id}`}>
                      Location
                    </Label>
                    <Input
                      id={`location-${engineer.id}`}
                      placeholder="City or region"
                      value={engineer.location}
                      onChange={(e) =>
                        updateEngineer(engineer.id, 'location', e.target.value)
                      }
                    />
                  </div>

                  <div className="space-y-2 md:col-span-2">
                    <Label htmlFor={`specializations-${engineer.id}`}>
                      Specializations
                    </Label>
                    <Input
                      id={`specializations-${engineer.id}`}
                      placeholder="MRI Scanner, CT Scanner, X-Ray"
                      value={engineer.specializations}
                      onChange={(e) =>
                        updateEngineer(
                          engineer.id,
                          'specializations',
                          e.target.value
                        )
                      }
                    />
                    <p className="text-xs text-gray-500">
                      Separate multiple specializations with commas
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Actions */}
        <div className="mt-6 flex gap-4">
          <Button
            variant="outline"
            onClick={addEngineer}
            className="flex-1"
          >
            <Plus className="mr-2 h-4 w-4" />
            Add Another Engineer
          </Button>
          <Button
            onClick={handleSave}
            className="flex-1"
          >
            <Save className="mr-2 h-4 w-4" />
            Save All Engineers
          </Button>
        </div>

        {/* Help Card */}
        <Card className="mt-6 border-blue-200 bg-blue-50">
          <CardContent className="pt-6">
            <div className="flex gap-3">
              <div className="text-blue-600">
                <svg
                  className="w-5 h-5"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fillRule="evenodd"
                    d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                    clipRule="evenodd"
                  />
                </svg>
              </div>
              <div className="text-sm text-blue-900">
                <p className="font-medium mb-1">Tip: Bulk Import Available</p>
                <p>
                  If you have multiple engineers to add, consider using the{' '}
                  <button
                    onClick={() => router.push('/engineers/import')}
                    className="underline font-medium"
                  >
                    CSV import
                  </button>{' '}
                  feature for faster data entry.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
