'use client';

import React, { useState } from 'react';
import { PartsAssignmentModal } from '@/components/PartsAssignmentModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Package, Wrench, CheckCircle } from 'lucide-react';

interface SelectedPart {
  id: string;
  part_name: string;
  part_number: string;
  quantity: number;
  unit_price: number;
  requires_engineer: boolean;
  engineer_level_required?: string;
}

export default function PartsDemo() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [assignedParts, setAssignedParts] = useState<SelectedPart[]>([]);

  const handleAssignParts = (parts: SelectedPart[]) => {
    setAssignedParts(parts);
    console.log('Assigned parts:', parts);
  };

  const totalCost = assignedParts.reduce((sum, part) => sum + (part.unit_price * part.quantity), 0);
  const requiresEngineer = assignedParts.some(p => p.requires_engineer);
  const highestLevel = assignedParts
    .filter(p => p.engineer_level_required)
    .map(p => p.engineer_level_required)
    .sort()
    .reverse()[0];

  return (
    <div className="container mx-auto p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Parts Assignment Demo</h1>
          <p className="text-muted-foreground mt-2">
            Test the Parts Assignment Modal with live data from the API
          </p>
        </div>
        <Button size="lg" onClick={() => setIsModalOpen(true)}>
          <Package className="h-5 w-5 mr-2" />
          Open Parts Browser
        </Button>
      </div>

      {/* Demo Equipment Card */}
      <Card>
        <CardHeader>
          <CardTitle>Sample Equipment: MRI Scanner GE-3T</CardTitle>
          <CardDescription>
            Click the button above to browse and assign spare parts to this equipment
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-2">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Model:</span>
              <span className="font-medium">GE Discovery MR750</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Serial Number:</span>
              <span className="font-medium">MRI-3T-001</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Location:</span>
              <span className="font-medium">Radiology Department - Room 102</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Assigned Parts Summary */}
      {assignedParts.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <CheckCircle className="h-5 w-5 text-green-500" />
              Assigned Parts ({assignedParts.length})
            </CardTitle>
            <CardDescription>
              Parts successfully assigned to the equipment
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Parts List */}
            <div className="space-y-3">
              {assignedParts.map(part => (
                <div
                  key={part.id}
                  className="flex items-center justify-between p-3 border rounded-lg"
                >
                  <div className="flex-1">
                    <p className="font-semibold">{part.part_name}</p>
                    <p className="text-sm text-muted-foreground">{part.part_number}</p>
                    {part.requires_engineer && (
                      <Badge variant="secondary" className="mt-1">
                        <Wrench className="h-3 w-3 mr-1" />
                        {part.engineer_level_required || 'Engineer Required'}
                      </Badge>
                    )}
                  </div>
                  <div className="text-right">
                    <p className="font-bold">
                      ₹{(part.unit_price * part.quantity).toLocaleString()}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {part.quantity} × ₹{part.unit_price}
                    </p>
                  </div>
                </div>
              ))}
            </div>

            {/* Summary */}
            <div className="border-t pt-4 space-y-2">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Total Parts:</span>
                <span className="font-medium">{assignedParts.length} items</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Total Quantity:</span>
                <span className="font-medium">
                  {assignedParts.reduce((sum, p) => sum + p.quantity, 0)}
                </span>
              </div>
              {requiresEngineer && (
                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">Engineer Required:</span>
                  <Badge variant="secondary">
                    <Wrench className="h-3 w-3 mr-1" />
                    {highestLevel || 'Yes'}
                  </Badge>
                </div>
              )}
              <div className="flex justify-between text-lg font-semibold pt-2 border-t">
                <span>Total Cost:</span>
                <span className="text-primary">₹{totalCost.toLocaleString()}</span>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex gap-3 pt-4">
              <Button variant="outline" onClick={() => setIsModalOpen(true)}>
                <Package className="h-4 w-4 mr-2" />
                Modify Parts
              </Button>
              <Button variant="outline" onClick={() => setAssignedParts([])}>
                Clear All
              </Button>
              <Button className="ml-auto">
                Save Assignment
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Empty State */}
      {assignedParts.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Package className="h-16 w-16 text-muted-foreground mb-4" />
            <p className="text-lg font-medium mb-2">No parts assigned yet</p>
            <p className="text-muted-foreground text-center mb-4">
              Click "Open Parts Browser" to browse and assign spare parts
            </p>
            <Button onClick={() => setIsModalOpen(true)}>
              Get Started
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Parts Assignment Modal */}
      <PartsAssignmentModal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onAssign={handleAssignParts}
        equipmentId="eq-001"
        equipmentName="MRI Scanner GE-3T"
      />
    </div>
  );
}
