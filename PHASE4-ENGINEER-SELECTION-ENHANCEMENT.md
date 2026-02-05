# Phase 4: Engineer Selection Widget Enhancement

**Status:** Ready for Implementation  
**Estimated Time:** 2-3 hours  
**Progress:** 75% → 100%

---

## Overview

Enhance the Engineer Selection Modal to display engineers from the entire service network (Manufacturer + Channel Partners + Sub-Dealers) with proper categorization, as shown in the UI mockup.

---

## Current Implementation

**File:** `admin-ui/src/components/EngineerSelectionModal.tsx`

**Current Features:**
- AI-powered engineer suggestions
- Match scoring
- Level-based filtering (Sr. Engineers, Skills Match, Less Workload)
- Engineer cards with details
- Assignment functionality

---

## Required Enhancements

### 1. Add "Manufacturer Certified" Category/Tab

Add a new filter category that shows engineers from partner network who are manufacturer-certified.

```typescript
const categories = [
  { id: 'best-match', label: 'Best Overall Match', icon: Star },
  { id: 'senior', label: 'Sr. Engineers Only', icon: Award },
  { id: 'workload', label: 'Less Workload', icon: Users },
  { id: 'skills', label: 'Skills Match', icon: CheckCircle },
  { id: 'certified', label: 'Manufacturer Certified', icon: Award }, // NEW
];
```

### 2. Integrate Network Engineers API

Replace or augment the current engineer suggestions with network engineers:

```typescript
// Current API call
const response = await apiClient.get(`/v1/engineers/suggestions?ticket_id=${ticketId}`);

// NEW: Fetch from network engineers API
import { partnersApi } from '@/lib/api/partners';

const fetchNetworkEngineers = async (manufacturerId: string, equipmentId?: string) => {
  const networkData = await partnersApi.getNetworkEngineers(manufacturerId, equipmentId);
  
  // networkData.grouped will have:
  // {
  //   'Manufacturer': [...engineers],
  //   'Channel Partner - ABC': [...engineers],
  //   'Sub-Dealer - XYZ': [...engineers]
  // }
  
  return networkData.engineers; // Returns flat list with category info
};

// Combine with existing AI scoring if needed
const mergeWithAIScoring = (networkEngineers, aiSuggestions) => {
  // Match engineers by ID and merge match_score
  return networkEngineers.map(engineer => {
    const aiMatch = aiSuggestions.find(ai => ai.engineer_id === engineer.id);
    return {
      ...engineer,
      match_score: aiMatch?.match_score || null,
      ai_recommendation: aiMatch?.recommendation || null,
    };
  });
};
```

### 3. Add Organization Type Badges

Update the engineer card to show organization type badge:

```typescript
const getOrgTypeBadge = (category: string) => {
  if (category === 'Manufacturer') {
    return (
      <Badge className="bg-blue-100 text-blue-800 border-blue-300">
        <Building2 className="h-3 w-3 mr-1" />
        Manufacturer
      </Badge>
    );
  }
  
  if (category.startsWith('Channel Partner')) {
    return (
      <Badge className="bg-orange-100 text-orange-800 border-orange-300">
        <Truck className="h-3 w-3 mr-1" />
        {category}
      </Badge>
    );
  }
  
  if (category.startsWith('Sub-Dealer')) {
    return (
      <Badge className="bg-purple-100 text-purple-800 border-purple-300">
        <Users className="h-3 w-3 mr-1" />
        {category}
      </Badge>
    );
  }
};

// Add to engineer card header
<div className="flex items-center gap-2 mb-2">
  <h3 className="text-lg font-semibold">{engineer.name}</h3>
  {getOrgTypeBadge(engineer.category)}
  {engineer.match_score && (
    <Badge variant="outline" className={getMatchScoreColor(engineer.match_score)}>
      {engineer.match_score}% Match
    </Badge>
  )}
</div>
```

### 4. Update Engineer Card Layout

Match the horizontal scrollable layout shown in the image:

```typescript
// Wrap engineer cards in horizontal scroll container
<div className="flex gap-4 overflow-x-auto pb-4">
  {filteredEngineers.map(engineer => (
    <div 
      key={engineer.id}
      className="min-w-[320px] max-w-[320px] border rounded-lg p-4 flex-shrink-0"
    >
      {/* Card content */}
      <div className="text-center mb-4">
        <div className="w-16 h-16 rounded-full bg-blue-500 mx-auto mb-3 flex items-center justify-center text-white text-2xl font-bold">
          {engineer.name.charAt(0)}
        </div>
        <h3 className="font-semibold text-lg">{engineer.name}</h3>
        <p className="text-sm text-gray-600">Level {engineer.level} • {engineer.skills}</p>
      </div>
      
      {/* Match percentage - prominent */}
      {engineer.match_score && (
        <div className="text-center mb-4">
          <div className={`text-4xl font-bold ${getMatchScoreColor(engineer.match_score)}`}>
            {engineer.match_score}%
          </div>
          <div className="text-sm text-gray-500">Match</div>
        </div>
      )}
      
      {/* Organization badge */}
      <div className="mb-4">
        {getOrgTypeBadge(engineer.category)}
      </div>
      
      {/* Details */}
      <div className="space-y-2 text-sm mb-4">
        <div className="flex items-center gap-2">
          <MapPin className="h-4 w-4" />
          <span>{engineer.location}</span>
        </div>
        {/* Add more details */}
      </div>
      
      {/* Assign button */}
      <Button 
        onClick={() => handleAssign(engineer.id)}
        className="w-full"
      >
        Assign Engineer
      </Button>
    </div>
  ))}
</div>
```

### 5. Filter Logic Update

Add filtering by organization category:

```typescript
const [activeCategory, setActiveCategory] = useState('best-match');
const [showOnlyManufacturer, setShowOnlyManufacturer] = useState(false);

const filteredEngineers = engineers.filter(engineer => {
  // Category filters
  if (activeCategory === 'certified') {
    return engineer.category !== 'Manufacturer'; // Show only partners
  }
  
  if (activeCategory === 'senior') {
    return engineer.level === 'L3' || engineer.level === 'senior';
  }
  
  if (activeCategory === 'skills') {
    return engineer.match_score && engineer.match_score >= 80;
  }
  
  // ... other filters
  
  return true;
}).sort((a, b) => {
  // Sort by match score desc
  return (b.match_score || 0) - (a.match_score || 0);
});
```

---

## Implementation Steps

### Step 1: Update API Integration
1. Import `partnersApi` from `@/lib/api/partners`
2. Get `manufacturerId` and `equipmentId` from ticket context
3. Call `partnersApi.getNetworkEngineers(manufacturerId, equipmentId)`
4. Merge with existing AI suggestions

### Step 2: Add Category Tab
1. Add "Manufacturer Certified" to category list
2. Update filter logic to show only partner engineers
3. Add icon (Award or Certificate)

### Step 3: Update Engineer Cards
1. Add organization type badge to each card
2. Make badge color-coded (blue/orange/purple)
3. Ensure match percentage is prominent
4. Add organization name to card

### Step 4: Update Layout
1. Optional: Change to horizontal scroll layout
2. Or keep vertical list with organization grouping
3. Ensure responsive design

### Step 5: Test
1. Create a ticket with equipment
2. Open engineer selection modal
3. Verify all categories work
4. Verify partner engineers show correctly
5. Verify assignment works for partner engineers

---

## Testing Checklist

- [ ] Engineers from manufacturer show with blue badge
- [ ] Engineers from channel partners show with orange badge
- [ ] Engineers from sub-dealers show with purple badge
- [ ] Match scores display correctly
- [ ] "Manufacturer Certified" filter works
- [ ] Assignment works for all engineer types
- [ ] Equipment-specific filtering works
- [ ] UI is responsive and looks good
- [ ] Horizontal scroll works (if implemented)
- [ ] No engineers from unassociated partners show up

---

**Ready for Implementation!**

This enhancement will complete the Partner Association feature, giving manufacturers full visibility and control over their service network.

**Estimated Lines of Code:** ~150-200 lines (modifications)  
**Files to Modify:** 1 (EngineerSelectionModal.tsx)  
**New Files:** 0  
**API Integration:** Already done (partnersApi exists)