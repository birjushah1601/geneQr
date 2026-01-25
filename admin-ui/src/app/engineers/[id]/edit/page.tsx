'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Textarea } from '@/components/ui/textarea';
import { ArrowLeft, Save, Loader2 } from 'lucide-react';
import engineersApi from '@/lib/api/engineers';
import { useToast } from '@/hooks/use-toast';

export default function EditEngineerPage() {
  const router = useRouter();
  const params = useParams();
  const { toast } = useToast();
  const engineerId = params.id as string;
  
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const [formData, setFormData] = useState({
    full_name: '',
    first_name: '',
    last_name: '',
    employee_id: '',
    email: '',
    phone: '',
    whatsapp_number: '',
    engineer_level: '1',
    experience_years: 0,
    employment_type: 'full_time',
    status: 'available',
    mobile_engineer: true,
    on_call_24x7: false,
    max_daily_tickets: 5,
    preferred_contact_method: 'phone',
    coverage_cities: '',
    coverage_states: '',
    skills: '',
    equipment_types: [] as string[],
  });

  // Fetch existing engineer data
  useEffect(() => {
    const fetchEngineer = async () => {
      try {
        const engineer = await engineersApi.getById(engineerId);
        
        setFormData({
          full_name: (engineer as any).full_name || (engineer as any).name || '',
          first_name: (engineer as any).first_name || '',
          last_name: (engineer as any).last_name || '',
          employee_id: (engineer as any).employee_id || '',
          email: (engineer as any).email || '',
          phone: (engineer as any).phone || '',
          whatsapp_number: (engineer as any).whatsapp_number || (engineer as any).whatsapp || '',
          engineer_level: String((engineer as any).engineer_level || 1),
          experience_years: (engineer as any).experience_years || 0,
          employment_type: (engineer as any).employment_type || 'full_time',
          status: (engineer as any).status || 'available',
          mobile_engineer: (engineer as any).mobile_engineer !== false,
          on_call_24x7: (engineer as any).on_call_24x7 || false,
          max_daily_tickets: (engineer as any).max_daily_tickets || 5,
          preferred_contact_method: (engineer as any).preferred_contact_method || 'phone',
          coverage_cities: (engineer as any).coverage_cities?.join(', ') || '',
          coverage_states: (engineer as any).coverage_states?.join(', ') || '',
          skills: (engineer as any).skills?.join(', ') || (engineer as any).specializations?.join(', ') || '',
          equipment_types: (engineer as any).specializations || (engineer as any).equipment_types || [],
        });
      } catch (error) {
        console.error('Failed to fetch engineer:', error);
        toast({
          title: 'Error',
          description: 'Failed to load engineer data',
          variant: 'destructive',
        });
      } finally {
        setLoading(false);
      }
    };

    fetchEngineer();
  }, [engineerId]);

  const handleInputChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.full_name || !formData.email) {
      toast({
        title: 'Validation Error',
        description: 'Please fill in all required fields (Name, Email)',
        variant: 'destructive',
      });
      return;
    }

    try {
      setSaving(true);
      
      const payload = {
        ...formData,
        engineer_level: parseInt(formData.engineer_level),
        coverage_cities: formData.coverage_cities ? formData.coverage_cities.split(',').map(c => c.trim()) : [],
        coverage_states: formData.coverage_states ? formData.coverage_states.split(',').map(s => s.trim()) : [],
        skills: formData.skills ? formData.skills.split(',').map(s => s.trim()) : [],
        specializations: formData.equipment_types,
      };

      await engineersApi.update(engineerId, payload);
      
      toast({
        title: 'Success',
        description: 'Engineer updated successfully',
      });
      
      router.push(`/engineers/${engineerId}`);
    } catch (error: any) {
      console.error('Failed to update engineer:', error);
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to update engineer',
        variant: 'destructive',
      });
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" onClick={() => router.push(`/engineers/${engineerId}`)}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
        <div>
          <h1 className="text-3xl font-bold">Edit Engineer</h1>
          <p className="text-muted-foreground">Update engineer profile</p>
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
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="employee_id">Employee ID</Label>
                <Input
                  id="employee_id"
                  value={formData.employee_id}
                  onChange={(e) => handleInputChange('employee_id', e.target.value)}
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
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="whatsapp_number">WhatsApp</Label>
                <Input
                  id="whatsapp_number"
                  type="tel"
                  value={formData.whatsapp_number}
                  onChange={(e) => handleInputChange('whatsapp_number', e.target.value)}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Employment Details */}
        <Card>
          <CardHeader>
            <CardTitle>Employment Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="engineer_level">Engineer Level *</Label>
                <Select
                  value={formData.engineer_level}
                  onValueChange={(value) => handleInputChange('engineer_level', value)}
                >
                  <SelectTrigger id="engineer_level">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="1">Level 1 (Junior)</SelectItem>
                    <SelectItem value="2">Level 2 (Mid-Level)</SelectItem>
                    <SelectItem value="3">Level 3 (Senior)</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="experience_years">Experience (Years)</Label>
                <Input
                  id="experience_years"
                  type="number"
                  min="0"
                  value={formData.experience_years}
                  onChange={(e) => handleInputChange('experience_years', parseInt(e.target.value) || 0)}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="status">Status</Label>
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

            <div className="space-y-2">
              <Label>Equipment Types Expertise</Label>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                {['MRI Scanner', 'CT Scanner', 'X-Ray Machine', 'Ultrasound', 'ECG Machine', 'Dialysis Machine', 'Ventilator', 'Anesthesia Machine'].map((equipmentType) => (
                  <div key={equipmentType} className="flex items-center space-x-2">
                    <Checkbox
                      id={`equipment-${equipmentType}`}
                      checked={formData.equipment_types.includes(equipmentType)}
                      onCheckedChange={(checked) => {
                        if (checked) {
                          handleInputChange('equipment_types', [...formData.equipment_types, equipmentType]);
                        } else {
                          handleInputChange('equipment_types', formData.equipment_types.filter(t => t !== equipmentType));
                        }
                      }}
                    />
                    <Label htmlFor={`equipment-${equipmentType}`} className="cursor-pointer text-sm font-normal">
                      {equipmentType}
                    </Label>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="skills">Skills (comma-separated)</Label>
              <Textarea
                id="skills"
                value={formData.skills}
                onChange={(e) => handleInputChange('skills', e.target.value)}
                rows={3}
              />
            </div>
          </CardContent>
        </Card>

        {/* Coverage */}
        <Card>
          <CardHeader>
            <CardTitle>Coverage Area</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="coverage_cities">Coverage Cities</Label>
                <Input
                  id="coverage_cities"
                  value={formData.coverage_cities}
                  onChange={(e) => handleInputChange('coverage_cities', e.target.value)}
                  placeholder="Mumbai, Delhi, Bangalore"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="coverage_states">Coverage States</Label>
                <Input
                  id="coverage_states"
                  value={formData.coverage_states}
                  onChange={(e) => handleInputChange('coverage_states', e.target.value)}
                  placeholder="Maharashtra, Delhi, Karnataka"
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Submit Button */}
        <div className="flex justify-end gap-4">
          <Button type="button" variant="outline" onClick={() => router.push(`/engineers/${engineerId}`)}>
            Cancel
          </Button>
          <Button type="submit" disabled={saving}>
            {saving ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Saving...
              </>
            ) : (
              <>
                <Save className="h-4 w-4 mr-2" />
                Save Changes
              </>
            )}
          </Button>
        </div>
      </form>
    </div>
  );
}
