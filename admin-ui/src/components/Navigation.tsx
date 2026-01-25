'use client';

import { useRouter, usePathname } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import OrganizationBadge from './OrganizationBadge';
import { 
  Home,
  Package, 
  Wrench, 
  Users, 
  Building2,
  LayoutDashboard,
  Settings,
  LogOut,
  Factory,
  Hospital,
  Truck,
  Sparkles
} from 'lucide-react';
import { cn } from '@/lib/utils';

interface NavItem {
  label: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  allowedOrgTypes: string[]; // Empty array = all types
}

// Navigation items configuration based on organization type
const navigationConfig: NavItem[] = [
  {
    label: 'Dashboard',
    href: '/dashboard',
    icon: LayoutDashboard,
    allowedOrgTypes: [], // All org types
  },
  {
    label: 'AI Onboarding',
    href: '/onboarding/ai-wizard',
    icon: Sparkles,
    allowedOrgTypes: ['system', 'manufacturer'], // System admins and manufacturer admins
  },
  {
    label: 'Equipment',
    href: '/equipment',
    icon: Package,
    allowedOrgTypes: [], // All org types
  },
  {
    label: 'Service Tickets',
    href: '/tickets',
    icon: Wrench,
    allowedOrgTypes: [], // All org types
  },
  {
    label: 'Engineers',
    href: '/engineers',
    icon: Users,
    allowedOrgTypes: ['distributor', 'dealer', 'manufacturer'], // Service providers
  },
  {
    label: 'Organizations',
    href: '/organizations',
    icon: Building2,
    allowedOrgTypes: ['system_admin'], // Admin only
  },
  {
    label: 'Manufacturers',
    href: '/manufacturers',
    icon: Factory,
    allowedOrgTypes: ['system_admin'], // Admin only
  },
];

export default function Navigation() {
  const router = useRouter();
  const pathname = usePathname();
  const { organizationContext, user, logout } = useAuth();

  // Filter navigation items based on organization type
  const visibleNavItems = navigationConfig.filter(item => {
    // If no restrictions, show to all
    if (item.allowedOrgTypes.length === 0) {
      return true;
    }

    // Special handling for AI Onboarding - only for admins
    if (item.href === '/onboarding/ai-wizard') {
      const orgType = organizationContext?.organization_type;
      const role = organizationContext?.role;
      // System admins or manufacturer admins only
      return orgType === 'system' || (orgType === 'manufacturer' && role === 'admin');
    }

    // Check if user's org type is in allowed list
    if (organizationContext?.organization_type) {
      return item.allowedOrgTypes.includes(organizationContext.organization_type);
    }

    // Check if user role is system_admin
    if (organizationContext?.role === 'system_admin' || organizationContext?.role === 'super_admin') {
      return true;
    }

    return false;
  });

  return (
    <div className="w-64 bg-white border-r min-h-screen flex flex-col fixed left-0 top-0 bottom-0 shadow-sm">
      {/* Logo/Header */}
      <div className="p-4 border-b">
        <div className="flex items-center gap-2 mb-3">
          <div className="w-8 h-8 rounded-lg bg-blue-600 flex items-center justify-center text-white font-bold">
            GM
          </div>
          <span className="font-bold text-lg">GeneQR</span>
        </div>
        {organizationContext && (
          <OrganizationBadge variant="compact" />
        )}
      </div>

      {/* Navigation Items - scrollable if overflow */}
      <nav className="flex-1 p-4 space-y-1 overflow-y-auto">
        {visibleNavItems.map((item) => {
          const Icon = item.icon;
          const isActive = pathname === item.href || pathname?.startsWith(item.href + '/');
          
          return (
            <button
              key={item.href}
              onClick={() => router.push(item.href)}
              className={cn(
                'w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 border-l-4',
                isActive
                  ? 'bg-blue-600 text-white shadow-md font-semibold border-blue-800'
                  : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900 border-transparent'
              )}
            >
              <Icon className={cn("h-4 w-4", isActive && "font-bold")} />
              {item.label}
            </button>
          );
        })}
      </nav>

      {/* User Profile & Logout */}
      <div className="p-4 border-t">
        <div className="mb-3 px-3 py-2 bg-gray-50 rounded-lg">
          <p className="text-sm font-medium text-gray-900 truncate">
            {user?.name || user?.email || 'User'}
          </p>
          <p className="text-xs text-gray-500 truncate">{organizationContext?.role || 'User'}</p>
        </div>
        
        <Button
          onClick={() => logout()}
          variant="outline"
          size="sm"
          className="w-full"
        >
          <LogOut className="h-4 w-4 mr-2" />
          Logout
        </Button>
      </div>
    </div>
  );
}
