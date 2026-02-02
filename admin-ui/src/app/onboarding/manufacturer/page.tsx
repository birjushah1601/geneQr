'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Building2, Mail, Phone, User, Globe, MapPin } from 'lucide-react';

export default function ManufacturerOnboarding() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  const [formData, setFormData] = useState({
    name: '',
    contact_person: '',
    email: '',
    phone: '',
    address: '',
    website: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      // Validate
      if (!formData.name || !formData.contact_person || !formData.email || !formData.phone) {
        throw new Error('Please fill in all required fields');
      }

      // Store manufacturer data in localStorage for now
      // In production, this would call the API
      const manufacturerId = `MFR-${Date.now()}`;
      const manufacturer = {
        id: manufacturerId,
        ...formData,
        created_at: new Date().toISOString(),
      };
      
      localStorage.setItem('current_manufacturer', JSON.stringify(manufacturer));
      
      // Redirect to equipment import
      router.push('/onboarding/equipment');
    } catch (err: any) {
      setError(err.message || 'Failed to create manufacturer');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <Card className="p-8">
      {/* Progress Steps */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center">
            <div className="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
              1
            </div>
            <div className="ml-3">
              <p className="font-semibold">Manufacturer Details</p>
              <p className="text-sm text-gray-500">Basic information</p>
            </div>
          </div>
          
          <div className="flex-1 h-1 mx-4 bg-gray-200"></div>
          
          <div className="flex items-center opacity-50">
            <div className="w-10 h-10 rounded-full bg-gray-300 text-gray-600 flex items-center justify-center font-bold">
              2
            </div>
            <div className="ml-3">
              <p className="font-semibold">Equipment Import</p>
              <p className="text-sm text-gray-500">Upload CSV</p>
            </div>
          </div>
          
          <div className="flex-1 h-1 mx-4 bg-gray-200"></div>
          
          <div className="flex items-center opacity-50">
            <div className="w-10 h-10 rounded-full bg-gray-300 text-gray-600 flex items-center justify-center font-bold">
              3
            </div>
            <div className="ml-3">
              <p className="font-semibold">Engineers</p>
              <p className="text-sm text-gray-500">Service team</p>
            </div>
          </div>
        </div>
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="text-center mb-6">
          <h2 className="text-2xl font-bold mb-2">Welcome to ServQR</h2>
          <p className="text-gray-600">
            Let's start by setting up your manufacturer profile
          </p>
        </div>

        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Company Name */}
          <div className="md:col-span-2">
            <Label htmlFor="name" className="flex items-center gap-2">
              <Building2 className="w-4 h-4" />
              Company Name *
            </Label>
            <Input
              id="name"
              name="name"
              placeholder="e.g., Siemens Healthineers"
              value={formData.name}
              onChange={handleChange}
              required
              className="mt-2"
            />
          </div>

          {/* Contact Person */}
          <div>
            <Label htmlFor="contact_person" className="flex items-center gap-2">
              <User className="w-4 h-4" />
              Contact Person *
            </Label>
            <Input
              id="contact_person"
              name="contact_person"
              placeholder="John Doe"
              value={formData.contact_person}
              onChange={handleChange}
              required
              className="mt-2"
            />
          </div>

          {/* Email */}
          <div>
            <Label htmlFor="email" className="flex items-center gap-2">
              <Mail className="w-4 h-4" />
              Email *
            </Label>
            <Input
              id="email"
              name="email"
              type="email"
              placeholder="contact@company.com"
              value={formData.email}
              onChange={handleChange}
              required
              className="mt-2"
            />
          </div>

          {/* Phone */}
          <div>
            <Label htmlFor="phone" className="flex items-center gap-2">
              <Phone className="w-4 h-4" />
              Phone Number *
            </Label>
            <Input
              id="phone"
              name="phone"
              type="tel"
              placeholder="+91-9876543210"
              value={formData.phone}
              onChange={handleChange}
              required
              className="mt-2"
            />
          </div>

          {/* Website */}
          <div>
            <Label htmlFor="website" className="flex items-center gap-2">
              <Globe className="w-4 h-4" />
              Website (Optional)
            </Label>
            <Input
              id="website"
              name="website"
              type="url"
              placeholder="https://www.company.com"
              value={formData.website}
              onChange={handleChange}
              className="mt-2"
            />
          </div>

          {/* Address */}
          <div className="md:col-span-2">
            <Label htmlFor="address" className="flex items-center gap-2">
              <MapPin className="w-4 h-4" />
              Address (Optional)
            </Label>
            <Input
              id="address"
              name="address"
              placeholder="Street, City, State, Country"
              value={formData.address}
              onChange={handleChange}
              className="mt-2"
            />
          </div>
        </div>

        {/* Submit Button */}
        <div className="flex justify-end gap-4 pt-6 border-t">
          <Button
            type="submit"
            disabled={loading}
            className="px-8"
          >
            {loading ? 'Creating...' : 'Next: Import Equipment'}
          </Button>
        </div>
      </form>
    </Card>
  );
}
