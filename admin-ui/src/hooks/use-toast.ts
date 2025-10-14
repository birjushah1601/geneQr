type ToastOptions = { title?: string; description?: string; variant?: 'default' | 'destructive' };

export function useToast() {
  const toast = ({ title, description }: ToastOptions) => {
    // Minimal no-op toast; replace with UI library if desired
    if (typeof window !== 'undefined') {
      const parts = [title, description].filter(Boolean).join(' - ');
      if (parts) console.log('[toast]', parts);
    }
  };
  return { toast };
}
