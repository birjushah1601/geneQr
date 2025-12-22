'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import DashboardLayout from '@/components/DashboardLayout';
import {
  ArrowLeft,
  Edit,
  Trash2,
  Phone,
  Mail,
  MapPin,
  Building2,
  Users,
  Star,
  Clock,
  Wrench,
  CheckCircle2,
  Loader2,
  MessageSquare,
} from 'lucide-react';
import engineersApi from '@/lib/api/engineers';
import { useToast } from '@/hooks/use-toast';

export default function EngineerDetailPage() {
  const router = useRouter();
  const params = useParams();
  const { toast } = useToast();
  const engineerId = params?.id as string;

  const [engineer, setEngineer] = useState<any | null>(null);
  const [organizations, setOrganizations] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [orgsLoading, setOrgsLoading] = useState(true);

  useEffect(() => {
    if (engineerId) {
      loadEngineer();
      // Organizations are optional for now; will enable once API is ready
      setOrgsLoading(false);
    }
  }, [engineerId]);

  const loadEngineer = async () => {
    try {
      setLoading(true);
      const data = await engineersApi.getById(engineerId);
      setEngineer(data);
    } catch (error) {
      console.error('Failed to load engineer:', error);
      toast({
        title: 'Error',
        description: 'Failed to load engineer details',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  const loadOrganizations = async () => {
    setOrgsLoading(false);
  };

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this engineer?')) return;

    try {
      await engineersApi.delete(engineerId);
      toast({
        title: 'Success',
        description: 'Engineer deleted successfully',
      });
      router.push('/engineers');
    } catch (error) {
      console.error('Failed to delete engineer:', error);
      toast({
        title: 'Error',
        description: 'Failed to delete engineer',
        variant: 'destructive',
      });
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'available':
        return 'bg-green-500';
      case 'busy':
        return 'bg-yellow-500';
      case 'off_duty':
        return 'bg-gray-500';
      default:
        return 'bg-gray-500';
    }
  };

  const getOrgTypeColor = (orgType?: string) => {
    switch (orgType) {
      case 'manufacturer':
        return 'bg-purple-100 text-purple-800';
      case 'dealer':
        return 'bg-blue-100 text-blue-800';
      case 'hospital':
        return 'bg-red-100 text-red-800';
      case 'distributor':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  if (!engineer) {
    return (
      <div className="p-6">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <p className="text-muted-foreground">Engineer not found</p>
            <Button onClick={() => router.push('/engineers')} className="mt-4">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Engineers
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" onClick={() => router.push('/engineers')}>
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back
          </Button>
          <div>
            <h1 className="text-3xl font-bold">{engineer.full_name}</h1>
            <p className="text-muted-foreground">{engineer.employee_id || 'No Employee ID'}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => router.push(`/engineers/${engineerId}/edit`)}>
            <Edit className="h-4 w-4 mr-2" />
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            <Trash2 className="h-4 w-4 mr-2" />
            Delete
          </Button>
        </div>
      </div>

      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="organizations">Organizations ({organizations.length})</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            {/* Main Info */}
            <Card className="lg:col-span-2">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>Engineer Information</CardTitle>
                  <div className={`w-3 h-3 rounded-full ${getStatusColor(engineer.status)}`} />
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-muted-foreground">Full Name</p>
                    <p className="font-medium">{engineer.full_name}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Employee ID</p>
                    <p className="font-medium">{engineer.employee_id || 'N/A'}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Email</p>
                    <div className="flex items-center gap-2">
                      <Mail className="h-4 w-4 text-muted-foreground" />
                      <p className="font-medium">{engineer.email || 'N/A'}</p>
                    </div>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Phone</p>
                    <div className="flex items-center gap-2">
                      <Phone className="h-4 w-4 text-muted-foreground" />
                      <p className="font-medium">{engineer.phone || 'N/A'}</p>
                    </div>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">WhatsApp</p>
                    <div className="flex items-center gap-2">
                      <MessageSquare className="h-4 w-4 text-muted-foreground" />
                      <p className="font-medium">{engineer.whatsapp_number || 'N/A'}</p>
                    </div>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Preferred Contact</p>
                    <Badge variant="outline">{engineer.preferred_contact_method}</Badge>
                  </div>
                </div>

                <div className="pt-4 border-t">
                  <p className="text-sm text-muted-foreground mb-2">Employment</p>
                  <div className="flex items-center gap-2">
                    <Badge className={getOrgTypeColor(engineer.org_type)}>
                      {engineer.org_type || 'Freelance'}
                    </Badge>
                    <Badge variant="secondary">{engineer.employment_type}</Badge>
                    {engineer.mobile_engineer && <Badge variant="outline">Mobile Engineer</Badge>}
                    {engineer.on_call_24x7 && (
                      <Badge variant="outline">
                        <Clock className="h-3 w-3 mr-1" />
                        24x7
                      </Badge>
                    )}
                  </div>
                </div>

                {engineer.org_name && (
                  <div className="pt-4 border-t">
                    <p className="text-sm text-muted-foreground mb-2">Primary Organization</p>
                    <div className="flex items-center gap-2">
                      <Building2 className="h-4 w-4 text-muted-foreground" />
                      <p className="font-medium">{engineer.org_name}</p>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Status & Stats */}
            <Card>
              <CardHeader>
                <CardTitle>Status & Capacity</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <p className="text-sm text-muted-foreground">Current Status</p>
                  <Badge className="mt-1">{engineer.status}</Badge>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Active Tickets</p>
                  <p className="text-2xl font-bold">
                    {engineer.active_tickets} / {engineer.max_daily_tickets}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Customer Rating</p>
                  <div className="flex items-center gap-2 mt-1">
                    <Star className="h-5 w-5 text-yellow-500 fill-yellow-500" />
                    <span className="text-2xl font-bold">{engineer.customer_rating.toFixed(1)}</span>
                  </div>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">First Time Fix Rate</p>
                  <p className="text-2xl font-bold text-green-600">
                    {engineer.first_time_fix_rate.toFixed(1)}%
                  </p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Total Tickets Resolved</p>
                  <div className="flex items-center gap-2 mt-1">
                    <CheckCircle2 className="h-5 w-5 text-green-600" />
                    <span className="text-2xl font-bold">{engineer.total_tickets_resolved}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Skills & Coverage */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Wrench className="h-5 w-5" />
                  Skills & Expertise
                </CardTitle>
              </CardHeader>
              <CardContent>
                {engineer.skills && engineer.skills.length > 0 ? (
                  <div className="flex flex-wrap gap-2">
                    {engineer.skills.map((skill, idx) => (
                      <Badge key={idx} variant="secondary">
                        {skill}
                      </Badge>
                    ))}
                  </div>
                ) : (
                  <p className="text-muted-foreground">No skills specified</p>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  Coverage Area
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                {engineer.coverage_states && engineer.coverage_states.length > 0 && (
                  <div>
                    <p className="text-sm text-muted-foreground mb-1">States</p>
                    <p className="font-medium">{engineer.coverage_states.join(', ')}</p>
                  </div>
                )}
                {engineer.coverage_cities && engineer.coverage_cities.length > 0 && (
                  <div>
                    <p className="text-sm text-muted-foreground mb-1">Cities</p>
                    <div className="flex flex-wrap gap-1">
                      {engineer.coverage_cities.map((city, idx) => (
                        <Badge key={idx} variant="outline">
                          {city}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
                {(!engineer.coverage_states || engineer.coverage_states.length === 0) &&
                  (!engineer.coverage_cities || engineer.coverage_cities.length === 0) && (
                    <p className="text-muted-foreground">No coverage area specified</p>
                  )}
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Organizations Tab */}
        <TabsContent value="organizations">
          <Card>
            <CardHeader>
              <CardTitle>Organization Memberships</CardTitle>
              <CardDescription>
                Multiple organizations this engineer is authorized to work with
              </CardDescription>
            </CardHeader>
            <CardContent>
              {orgsLoading ? (
                <div className="flex items-center justify-center py-12">
                  <Loader2 className="h-8 w-8 animate-spin" />
                </div>
              ) : organizations.length === 0 ? (
                <div className="text-center py-12">
                  <Building2 className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                  <p className="text-muted-foreground">No organization memberships</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {organizations.map((membership) => (
                    <div
                      key={membership.org_id}
                      className="flex items-center justify-between p-4 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <Building2 className="h-5 w-5 text-muted-foreground" />
                        <div>
                          <p className="font-medium">{membership.org_name || 'Unknown Organization'}</p>
                          <div className="flex items-center gap-2 mt-1">
                            {membership.org_type && (
                              <Badge className={getOrgTypeColor(membership.org_type)} variant="outline">
                                {membership.org_type}
                              </Badge>
                            )}
                            {membership.role && (
                              <Badge variant="secondary">{membership.role}</Badge>
                            )}
                          </div>
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => {
                          /* Navigate to organization */
                        }}
                      >
                        View
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Performance Tab */}
        <TabsContent value="performance">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-sm font-medium">Total Resolved</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-bold">{engineer.total_tickets_resolved}</div>
                <p className="text-xs text-muted-foreground">All-time tickets</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle className="text-sm font-medium">Customer Rating</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-2">
                  <Star className="h-6 w-6 text-yellow-500 fill-yellow-500" />
                  <span className="text-3xl font-bold">{engineer.customer_rating.toFixed(1)}</span>
                </div>
                <p className="text-xs text-muted-foreground">Average rating</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle className="text-sm font-medium">First Time Fix</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-bold text-green-600">
                  {engineer.first_time_fix_rate.toFixed(1)}%
                </div>
                <p className="text-xs text-muted-foreground">Success rate</p>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
      </div>
    </DashboardLayout>
  );
}
