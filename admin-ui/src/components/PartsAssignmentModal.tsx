'use client';

import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { 
  Search, 
  ShoppingCart, 
  Package, 
  Wrench, 
  DollarSign,
  AlertCircle,
  CheckCircle,
  X,
  Plus,
  Minus
} from 'lucide-react';

// Types
interface SparePart {
  id: string;
  part_number: string;
  part_name: string;
  category: string;
  subcategory?: string;
  description?: string;
  unit_price: number;
  currency: string;
  is_available: boolean;
  stock_status: string;
  requires_engineer: boolean;
  engineer_level_required?: string;
  installation_time_minutes?: number;
  lead_time_days?: number;
  minimum_order_quantity: number;
  image_url?: string;
  photos?: string[];
}

interface SelectedPart extends SparePart {
  quantity: number;
}

interface PartsAssignmentModalProps {
  open: boolean;
  onClose: () => void;
  onAssign: (parts: SelectedPart[]) => void;
  equipmentId?: string;
  equipmentName?: string;
}

const CATEGORIES = [
  'component',
  'consumable',
  'accessory',
  'sensor',
  'filter',
  'battery'
];

export function PartsAssignmentModal({
  open,
  onClose,
  onAssign,
  equipmentId,
  equipmentName
}: PartsAssignmentModalProps) {
  const [parts, setParts] = useState<SparePart[]>([]);
  const [selectedCategory, setSelectedCategory] = useState('All Categories');
  const [filteredParts, setFilteredParts] = useState<SparePart[]>([]);
  const [selectedParts, setSelectedParts] = useState<Map<string, SelectedPart>>(new Map());
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const [requiresEngineerFilter, setRequiresEngineerFilter] = useState<boolean | null>(null);
  const [activeTab, setActiveTab] = useState('browse');

  // Fetch parts from API
  useEffect(() => {
    if (open) {
      fetchParts();
    }
  }, [open]);

  const fetchParts = async () => {
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8081/api/v1/catalog/parts', {
        headers: {
          'X-Tenant-ID': 'default',
        },
      });
      
      if (!response.ok) throw new Error('Failed to fetch parts');
      
      const data = await response.json();
      setParts(data.parts || []);
      setFilteredParts(data.parts || []);
    } catch (error) {
      console.error('Error fetching parts:', error);
      setParts([]);
      setFilteredParts([]);
    } finally {
      setLoading(false);
    }
  };

  // Filter parts based on search and filters
  useEffect(() => {
    let filtered = [...parts];

    // Category filter (multiselect)
    if (selectedCategories.length > 0) {
      filtered = filtered.filter(p => selectedCategories.includes(p.category));
    }

    // Search filter
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(p =>
        p.part_name.toLowerCase().includes(query) ||
        p.part_number.toLowerCase().includes(query) ||
        p.description?.toLowerCase().includes(query)
      );
    }

    // Engineer requirement filter
    if (requiresEngineerFilter !== null) {
      filtered = filtered.filter(p => p.requires_engineer === requiresEngineerFilter);
    }

    setFilteredParts(filtered);
  }, [parts, searchQuery, selectedCategories, requiresEngineerFilter]);

  // Toggle category selection
  const toggleCategory = (category: string) => {
    setSelectedCategories(prev =>
      prev.includes(category)
        ? prev.filter(c => c !== category)
        : [...prev, category]
    );
  };

  // Add part to selection
  const handleSelectPart = (part: SparePart) => {
    const newSelected = new Map(selectedParts);
    if (newSelected.has(part.id)) {
      newSelected.delete(part.id);
    } else {
      newSelected.set(part.id, { ...part, quantity: part.minimum_order_quantity || 1 });
    }
    setSelectedParts(newSelected);
  };

  // Update quantity
  const updateQuantity = (partId: string, delta: number) => {
    const newSelected = new Map(selectedParts);
    const part = newSelected.get(partId);
    if (part) {
      const newQty = Math.max(part.minimum_order_quantity || 1, part.quantity + delta);
      newSelected.set(partId, { ...part, quantity: newQty });
      setSelectedParts(newSelected);
    }
  };

  // Calculate totals
  const calculateTotals = () => {
    let totalCost = 0;
    let requiresEngineer = false;
    let highestEngineerLevel = '';
    let totalInstallTime = 0;

    selectedParts.forEach(part => {
      totalCost += part.unit_price * part.quantity;
      if (part.requires_engineer) {
        requiresEngineer = true;
        if (part.engineer_level_required) {
          // L3 > L2 > L1
          if (highestEngineerLevel === '' || 
              part.engineer_level_required > highestEngineerLevel) {
            highestEngineerLevel = part.engineer_level_required;
          }
        }
      }
      if (part.installation_time_minutes) {
        totalInstallTime += part.installation_time_minutes * part.quantity;
      }
    });

    return { totalCost, requiresEngineer, highestEngineerLevel, totalInstallTime };
  };

  const { totalCost, requiresEngineer, highestEngineerLevel, totalInstallTime } = calculateTotals();

  // Handle assignment
  const handleAssign = () => {
    const partsArray = Array.from(selectedParts.values());
    onAssign(partsArray);
    onClose();
  };

  // Reset modal
  const handleClose = () => {
    setSelectedParts(new Map());
    setSearchQuery('');
    setSelectedCategory('All Categories');
    setRequiresEngineerFilter(null);
    setActiveTab('browse');
    onClose();
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-6xl max-h-[90vh]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Package className="h-5 w-5" />
            Assign Spare Parts
            {equipmentName && (
              <span className="text-sm font-normal text-muted-foreground">
                to {equipmentName}
              </span>
            )}
          </DialogTitle>
          <DialogDescription>
            Browse and select spare parts to assign. Selected parts will be added to your cart.
          </DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="browse">
              <Search className="h-4 w-4 mr-2" />
              Browse Parts ({filteredParts.length})
            </TabsTrigger>
            <TabsTrigger value="cart">
              <ShoppingCart className="h-4 w-4 mr-2" />
              Cart ({selectedParts.size})
            </TabsTrigger>
          </TabsList>

          {/* Browse Tab */}
          <TabsContent value="browse" className="space-y-4 mt-4">
            {/* Search and Filters Row */}
            <div className="flex items-center gap-3">
              {/* Search Bar */}
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search parts by name or part number..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>

              {/* Engineer Filter Buttons - Right Aligned */}
              <Button
                variant={requiresEngineerFilter === true ? 'default' : 'outline'}
                size="sm"
                onClick={() => setRequiresEngineerFilter(
                  requiresEngineerFilter === true ? null : true
                )}
              >
                <Wrench className="h-4 w-4 mr-1" />
                Needs Engineer
              </Button>
              <Button
                variant={requiresEngineerFilter === false ? 'default' : 'outline'}
                size="sm"
                onClick={() => setRequiresEngineerFilter(
                  requiresEngineerFilter === false ? null : false
                )}
              >
                <CheckCircle className="h-4 w-4 mr-1" />
                Self-Service
              </Button>
            </div>

            {/* Category Filter Badges */}
            <div className="flex items-center gap-2 flex-wrap">
              <span className="text-sm text-muted-foreground font-medium">Categories:</span>
              {CATEGORIES.map(cat => (
                <Badge
                  key={cat}
                  variant={selectedCategories.includes(cat) ? 'default' : 'outline'}
                  className="cursor-pointer"
                  onClick={() => toggleCategory(cat)}
                >
                  {cat}
                  {selectedCategories.includes(cat) && (
                    <X className="h-3 w-3 ml-1" />
                  )}
                </Badge>
              ))}
              {selectedCategories.length > 0 && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setSelectedCategories([])}
                  className="h-6 text-xs"
                >
                  Clear All
                </Button>
              )}
            </div>

            {/* Parts List */}
            <ScrollArea className="h-[400px] rounded-md border">
              {loading ? (
                <div className="flex items-center justify-center h-full">
                  <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
                    <p className="text-sm text-muted-foreground">Loading parts...</p>
                  </div>
                </div>
              ) : filteredParts.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full p-8">
                  <Package className="h-12 w-12 text-muted-foreground mb-3" />
                  <p className="text-lg font-medium">No parts found</p>
                  <p className="text-sm text-muted-foreground">Try adjusting your filters</p>
                </div>
              ) : (
                <div className="grid gap-3 p-4">
                  {filteredParts.map(part => {
                    const isSelected = selectedParts.has(part.id);
                    return (
                      <Card 
                        key={part.id} 
                        className={`cursor-pointer transition-all hover:shadow-md ${
                          isSelected ? 'ring-2 ring-primary' : ''
                        }`}
                        onClick={() => handleSelectPart(part)}
                      >
                        <CardContent className="p-4">
                          <div className="flex items-start gap-4">
                            <Checkbox
                              checked={isSelected}
                              onCheckedChange={() => handleSelectPart(part)}
                              onClick={(e) => e.stopPropagation()}
                            />
                            
                            {/* Part Image */}
                            {part.image_url && (
                              <div className="w-20 h-20 flex-shrink-0 rounded-md overflow-hidden bg-gray-100">
                                <img 
                                  src={part.image_url} 
                                  alt={part.part_name}
                                  className="w-full h-full object-cover"
                                  onError={(e) => {
                                    // Hide image if it fails to load
                                    (e.target as HTMLImageElement).style.display = 'none';
                                  }}
                                />
                              </div>
                            )}
                            
                            <div className="flex-1 space-y-2">
                              <div className="flex items-start justify-between">
                                <div>
                                  <h4 className="font-semibold">{part.part_name}</h4>
                                  <p className="text-sm text-muted-foreground">
                                    {part.part_number}
                                  </p>
                                </div>
                                
                                <div className="text-right">
                                  <p className="text-lg font-bold">
                                    ₹{part.unit_price.toLocaleString()}
                                  </p>
                                  <p className="text-xs text-muted-foreground">{part.currency}</p>
                                </div>
                              </div>

                              {part.description && (
                                <p className="text-sm text-muted-foreground line-clamp-2">
                                  {part.description}
                                </p>
                              )}

                              <div className="flex flex-wrap gap-2">
                                <Badge variant="outline">{part.category}</Badge>
                                
                                {part.is_available ? (
                                  <Badge variant="default" className="bg-green-500">
                                    <CheckCircle className="h-3 w-3 mr-1" />
                                    In Stock
                                  </Badge>
                                ) : (
                                  <Badge variant="destructive">
                                    <AlertCircle className="h-3 w-3 mr-1" />
                                    Out of Stock
                                  </Badge>
                                )}

                                {part.requires_engineer && (
                                  <Badge variant="secondary">
                                    <Wrench className="h-3 w-3 mr-1" />
                                    {part.engineer_level_required || 'Engineer Required'}
                                  </Badge>
                                )}

                                {part.installation_time_minutes && (
                                  <Badge variant="outline">
                                    ⏱️ {part.installation_time_minutes}min
                                  </Badge>
                                )}
                              </div>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    );
                  })}
                </div>
              )}
            </ScrollArea>
          </TabsContent>

          {/* Cart Tab */}
          <TabsContent value="cart" className="space-y-4">
            <ScrollArea className="h-[400px] rounded-md border">
              {selectedParts.size === 0 ? (
                <div className="flex flex-col items-center justify-center h-full p-8">
                  <ShoppingCart className="h-12 w-12 text-muted-foreground mb-3" />
                  <p className="text-lg font-medium">Your cart is empty</p>
                  <p className="text-sm text-muted-foreground">Add parts from the Browse tab</p>
                </div>
              ) : (
                <div className="p-4 space-y-3">
                  {Array.from(selectedParts.values()).map(part => (
                    <Card key={part.id}>
                      <CardContent className="p-4">
                        <div className="flex items-center gap-4">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleSelectPart(part)}
                          >
                            <X className="h-4 w-4" />
                          </Button>

                          <div className="flex-1">
                            <h4 className="font-semibold">{part.part_name}</h4>
                            <p className="text-sm text-muted-foreground">{part.part_number}</p>
                          </div>

                          <div className="flex items-center gap-2">
                            <Button
                              variant="outline"
                              size="icon"
                              onClick={() => updateQuantity(part.id, -1)}
                              disabled={part.quantity <= (part.minimum_order_quantity || 1)}
                            >
                              <Minus className="h-4 w-4" />
                            </Button>
                            <span className="w-12 text-center font-medium">{part.quantity}</span>
                            <Button
                              variant="outline"
                              size="icon"
                              onClick={() => updateQuantity(part.id, 1)}
                            >
                              <Plus className="h-4 w-4" />
                            </Button>
                          </div>

                          <div className="text-right min-w-[100px]">
                            <p className="font-bold">
                              ₹{(part.unit_price * part.quantity).toLocaleString()}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              ₹{part.unit_price} × {part.quantity}
                            </p>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </ScrollArea>

            {/* Summary */}
            {selectedParts.size > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle>Order Summary</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Total Parts:</span>
                    <span className="font-medium">{selectedParts.size} items</span>
                  </div>
                  
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Total Quantity:</span>
                    <span className="font-medium">
                      {Array.from(selectedParts.values()).reduce((sum, p) => sum + p.quantity, 0)}
                    </span>
                  </div>

                  {totalInstallTime > 0 && (
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Est. Install Time:</span>
                      <span className="font-medium">
                        {Math.floor(totalInstallTime / 60)}h {totalInstallTime % 60}m
                      </span>
                    </div>
                  )}

                  {requiresEngineer && (
                    <div className="flex justify-between items-center">
                      <span className="text-muted-foreground">Engineer Required:</span>
                      <Badge variant="secondary">
                        <Wrench className="h-3 w-3 mr-1" />
                        {highestEngineerLevel || 'Yes'}
                      </Badge>
                    </div>
                  )}

                  <Separator />

                  <div className="flex justify-between text-lg">
                    <span className="font-semibold">Total Cost:</span>
                    <span className="font-bold text-primary">
                      ₹{totalCost.toLocaleString()}
                    </span>
                  </div>
                </CardContent>
              </Card>
            )}
          </TabsContent>
        </Tabs>

        <DialogFooter>
          <Button variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          <Button 
            onClick={handleAssign} 
            disabled={selectedParts.size === 0}
          >
            <CheckCircle className="h-4 w-4 mr-2" />
            Assign {selectedParts.size} Part{selectedParts.size !== 1 ? 's' : ''}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
