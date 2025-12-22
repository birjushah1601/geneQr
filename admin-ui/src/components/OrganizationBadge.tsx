'use client';

import { useAuth } from '@/contexts/AuthContext';
import { 
  Factory, 
  Hospital, 
  Truck, 
  ShoppingBag,
  Building2,
  Beaker,
  Camera,
  Shield
} from 'lucide-react';
import { cn } from '@/lib/utils';

interface OrganizationBadgeProps {
  variant?: 'default' | 'compact' | 'large';
  showTooltip?: boolean;
  className?: string;
}

// Organization type configuration
const orgTypeConfig = {
  manufacturer: {
    label: 'Manufacturer',
    icon: Factory,
    color: 'bg-indigo-100 text-indigo-700 border-indigo-200',
    iconColor: 'text-indigo-600',
    description: 'Equipment manufacturer and OEM',
  },
  hospital: {
    label: 'Hospital',
    icon: Hospital,
    color: 'bg-red-100 text-red-700 border-red-200',
    iconColor: 'text-red-600',
    description: 'Medical facility and healthcare provider',
  },
  imaging_center: {
    label: 'Imaging Center',
    icon: Camera,
    color: 'bg-pink-100 text-pink-700 border-pink-200',
    iconColor: 'text-pink-600',
    description: 'Diagnostic imaging facility',
  },
  distributor: {
    label: 'Distributor',
    icon: Truck,
    color: 'bg-purple-100 text-purple-700 border-purple-200',
    iconColor: 'text-purple-600',
    description: 'Equipment distributor and service provider',
  },
  dealer: {
    label: 'Dealer',
    icon: ShoppingBag,
    color: 'bg-green-100 text-green-700 border-green-200',
    iconColor: 'text-green-600',
    description: 'Authorized equipment dealer',
  },
  supplier: {
    label: 'Supplier',
    icon: Beaker,
    color: 'bg-teal-100 text-teal-700 border-teal-200',
    iconColor: 'text-teal-600',
    description: 'Parts and supplies provider',
  },
  system_admin: {
    label: 'System Admin',
    icon: Shield,
    color: 'bg-gray-100 text-gray-700 border-gray-200',
    iconColor: 'text-gray-600',
    description: 'Platform administrator',
  },
};

export default function OrganizationBadge({ 
  variant = 'default', 
  showTooltip = true,
  className 
}: OrganizationBadgeProps) {
  const { organizationContext } = useAuth();

  if (!organizationContext) {
    return null;
  }

  const orgType = organizationContext.organization_type;
  const role = organizationContext.role;

  // Get config for org type, fallback to generic
  const config = orgTypeConfig[orgType as keyof typeof orgTypeConfig] || {
    label: 'Organization',
    icon: Building2,
    color: 'bg-blue-100 text-blue-700 border-blue-200',
    iconColor: 'text-blue-600',
    description: 'Organization member',
  };

  const Icon = config.icon;

  // Variant styles
  const variantStyles = {
    compact: 'px-2 py-1 text-xs',
    default: 'px-3 py-1.5 text-sm',
    large: 'px-4 py-2 text-base',
  };

  const iconSizes = {
    compact: 'h-3 w-3',
    default: 'h-4 w-4',
    large: 'h-5 w-5',
  };

  const badge = (
    <div
      className={cn(
        'inline-flex items-center gap-2 rounded-full border font-medium transition-all',
        config.color,
        variantStyles[variant],
        showTooltip && 'cursor-help',
        className
      )}
      title={showTooltip ? config.description : undefined}
    >
      <Icon className={cn(iconSizes[variant], config.iconColor)} />
      <span>{config.label}</span>
      {role && role !== 'user' && (
        <span className="ml-1 opacity-60 text-xs">
          ({role})
        </span>
      )}
    </div>
  );

  return badge;
}

// Export a mini version for compact spaces
export function OrganizationIcon() {
  const { organizationContext } = useAuth();

  if (!organizationContext) {
    return null;
  }

  const orgType = organizationContext.organization_type;
  const config = orgTypeConfig[orgType as keyof typeof orgTypeConfig] || {
    icon: Building2,
    iconColor: 'text-blue-600',
    color: 'bg-blue-100',
  };

  const Icon = config.icon;

  return (
    <div className={cn('w-8 h-8 rounded-lg flex items-center justify-center', config.color)}>
      <Icon className={cn('h-4 w-4', config.iconColor)} />
    </div>
  );
}
