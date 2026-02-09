/**
 * Feature Flags System
 * Controls visibility of features based on URL query parameters
 * For demo purposes - can hide/show features without code changes
 */

export type FeatureFlag = 
  | 'PhoneLogin'
  | 'AIOnboarding'
  | 'AddNew'  // Generic flag for add/import buttons
  | 'AIDiagnosis';

/**
 * Check if a feature is enabled via URL query parameter
 * Usage: ?enable=FeatureName or ?enable=FeatureName,AnotherFeature
 */
export function isFeatureEnabled(feature: FeatureFlag): boolean {
  if (typeof window === 'undefined') return false;
  
  const params = new URLSearchParams(window.location.search);
  const enabledFeatures = params.get('enable');
  
  if (!enabledFeatures) return false;
  
  // Support comma-separated list: ?enable=Feature1,Feature2
  const features = enabledFeatures.split(',').map(f => f.trim());
  return features.includes(feature);
}

/**
 * Check multiple features at once
 */
export function areAnyFeaturesEnabled(...features: FeatureFlag[]): boolean {
  return features.some(feature => isFeatureEnabled(feature));
}

/**
 * Get all enabled features from URL
 */
export function getEnabledFeatures(): FeatureFlag[] {
  if (typeof window === 'undefined') return [];
  
  const params = new URLSearchParams(window.location.search);
  const enabledFeatures = params.get('enable');
  
  if (!enabledFeatures) return [];
  
  return enabledFeatures.split(',').map(f => f.trim()) as FeatureFlag[];
}

/**
 * Add feature flag to current URL
 */
export function addFeatureToUrl(feature: FeatureFlag): string {
  const url = new URL(window.location.href);
  const enabled = getEnabledFeatures();
  
  if (!enabled.includes(feature)) {
    enabled.push(feature);
  }
  
  url.searchParams.set('enable', enabled.join(','));
  return url.toString();
}

/**
 * Demo-specific defaults
 * These features are HIDDEN by default for demo
 */
export const DEMO_HIDDEN_FEATURES: FeatureFlag[] = [
  'PhoneLogin',
  'AIOnboarding',
  'AddNew',
  'AIDiagnosis'
];

/**
 * Check if we're in demo mode (all features hidden by default)
 */
export function isDemoMode(): boolean {
  // Could be configured via environment variable
  return process.env.NEXT_PUBLIC_DEMO_MODE === 'true';
}
