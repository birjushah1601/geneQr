'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Textarea } from '@/components/ui/textarea';
import { ArrowLeft, Save } from 'lucide-react';
import engineersApi from '@/lib/api/engineers';
import { useToast } from '@/hooks/use-toast';

export default function NewEngineerPage() {
  const router = useRouter();
  const { toast } = useToast();
  const [saving, setSaving] = useState(false);

  const [formData, setFormData] = useState({
    full_name: '',
    first_name: '',
    last_name: '',
    employee_id: '',
    email: '',
    phone: '',
    whatsapp_number: '',
    employment_type: 'full_time',
    status: 'available',
    mobile_engineer: true,
    on_call_24x7: false,
    max_daily_tickets: 5,
    preferred_contact_method: 'phone',
    coverage_cities: '',
    coverage_states: '',
    skills: '',
  });

  const handleInputChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.full_name || !formData.email) {
      toast({
        title: 'Validation Error',
        description: 'Please fill in all required fields',
        variant: 'destructive',
      });
      return;
    }

    try {
      setSaving(true);
      
      const payload = {
        ...formData,
        coverage_cities: formData.coverage_cities ? formData.coverage_cities.split(',').map(c => c.trim()) : [],
        coverage_states: formData.coverage_states ? formData.coverage_states.split(',').map(s => s.trim()) : [],
        skills: formData.skills ? formData.skills.split(',').map(s => s.trim()) : [],
        active_tickets: 0,
        total_tickets_resolved: 0,
        customer_rating: 0,
        first_time_fix_rate: 0,
      };

      const response = await engineersApi.create(payload);
      
      toast({
        title: 'Success',
        description: 'Engineer created successfully',
      });
      
      router.push(`/engineers/${response.id}`);
    } catch (error: any) {
      console.error('Failed to create engineer:', error);
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to create engineer',
        variant: 'destructive',
      });
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" onClick={() => router.push('/engineers')}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
        <div>
          <h1 className="text-3xl font-bold">Add New Engineer</h1>
          <p className="text-muted-foreground">Create a new field service engineer profile</p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Information */}
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
            <CardDescription>Engineer's personal and contact details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="full_name">Full Name *</Label>
                <Input
                  id="full_name"
                  value={formData.full_name}
                  onChange={(e) => handleInputChange('full_name', e.target.value)}
                  placeholder="John Doe"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="employee_id">Employee ID</Label>
                <Input
                  id="employee_id"
                  value={formData.employee_id}
                  onChange={(e) => handleInputChange('employee_id', e.target.value)}
                  placeholder="EMP-001"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="first_name">First Name</Label>
                <Input
                  id="first_name"
                  value={formData.first_name}
                  onChange={(e) => handleInputChange('first_name', e.target.value)}
                  placeholder="John"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="last_name">Last Name</Label>
                <Input
                  id="last_name"
                  value={formData.last_name}
                  onChange={(e) => handleInputChange('last_name', e.target.value)}
                  placeholder="Doe"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="email">Email *</Label>
                <Input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => handleInputChange('email', e.target.value)}
                  placeholder="engineer@example.com"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="phone">Phone</Label>
                <Input
                  id="phone"
                  type="tel"
                  value={formData.phone}
                  onChange={(e) => handleInputChange('phone', e.target.value)}
                  placeholder="+91-9876543210"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="whatsapp_number">WhatsApp</Label>
                <Input
                  id="whatsapp_number"
                  type="tel"
                  value={formData.whatsapp_number}
                  onChange={(e) => handleInputChange('whatsapp_number', e.target.value)}
                  placeholder="+91-9876543210"
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Employment Details */}
        <Card>
          <CardHeader>
            <CardTitle>Employment Details</CardTitle>
            <CardDescription>Work arrangement and capacity</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="employment_type">Employment Type</Label>
                <Select
                  value={formData.employment_type}
                  onValueChange={(value) => handleInputChange('employment_type', value)}
                >
                  <SelectTrigger id="employment_type">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="full_time">Full Time</SelectItem>
                    <SelectItem value="part_time">Part Time</SelectItem>
                    <SelectItem value="contract">Contract</SelectItem>
                    <SelectItem value="freelance">Freelance</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="status">Current Status</Label>
                <Select
                  value={formData.status}
                  onValueChange={(value) => handleInputChange('status', value)}
                >
                  <SelectTrigger id="status">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="available">Available</SelectItem>
                    <SelectItem value="busy">Busy</SelectItem>
                    <SelectItem value="off_duty">Off Duty</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="max_daily_tickets">Max Daily Tickets</Label>
                <Input
                  id="max_daily_tickets"
                  type="number"
                  value={formData.max_daily_tickets}
                  onChange={(e) => handleInputChange('max_daily_tickets', parseInt(e.target.value))}
                  min="1"
                  max="20"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="preferred_contact_method">Preferred Contact Method</Label>
                <Select
                  value={formData.preferred_contact_method}
                  onValueChange={(value) => handleInputChange('preferred_contact_method', value)}
                >
                  <SelectTrigger id="preferred_contact_method">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="phone">Phone</SelectItem>
                    <SelectItem value="whatsapp">WhatsApp</SelectItem>
                    <SelectItem value="email">Email</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-4">
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="mobile_engineer"
                  checked={formData.mobile_engineer}
                  onCheckedChange={(checked) => handleInputChange('mobile_engineer', checked)}
                />
                <Label htmlFor="mobile_engineer" className="cursor-pointer">
                  Mobile Engineer (can travel to customer locations)
                </Label>
              </div>

              <div className="flex items-center space-x-2">
                <Checkbox
                  id="on_call_24x7"
                  checked={formData.on_call_24x7}
                  onCheckedChange={(checked) => handleInputChange('on_call_24x7', checked)}
                />
                <Label htmlFor="on_call_24x7" className="cursor-pointer">
                  Available 24x7 for emergency calls
                </Label>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Skills & Coverage */}
        <Card>
          <CardHeader>
            <CardTitle>Skills & Coverage Area</CardTitle>
            <CardDescription>Technical expertise and service areas</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="skills">Skills (comma-separated)</Label>
              <Textarea
                id="skills"
                value={formData.skills}
                onChange={(e) => handleInputChange('skills', e.target.value)}
                placeholder="CT Scanner, MRI, X-Ray, Ultrasound"
                rows={3}
              />
              <p className="text-xs text-muted-foreground">
                Enter equipment types or specializations, separated by commas
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="coverage_cities">Coverage Cities (comma-separated)</Label>
                <Textarea
                  id="coverage_cities"
                  value={formData.coverage_cities}
                  onChange={(e) => handleInputChange('coverage_cities', e.target.value)}
                  placeholder="Delhi, Noida, Gurgaon, Faridabad"
                  rows={3}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="coverage_states">Coverage States (comma-separated)</Label>
                <Textarea
                  id="coverage_states"
                  value={formData.coverage_states}
                  onChange={(e) => handleInputChange('coverage_states', e.target.value)}
                  placeholder="Delhi, Haryana, Uttar Pradesh"
                  rows={3}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Actions */}
        <div className="flex justify-end gap-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => router.push('/engineers')}
            disabled={saving}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={saving}>
            {saving ? (
              <>Saving...</>
            ) : (
              <>
                <Save className="h-4 w-4 mr-2" />
                Create Engineer
              </>
            )}
          </Button>
        </div>
      </form>
    </div>
  );
}
