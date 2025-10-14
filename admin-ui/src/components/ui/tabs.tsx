"use client";
import React, { createContext, useContext, useMemo, useState } from 'react';

type TabsContextValue = {
  value: string;
  setValue: (v: string) => void;
};

const TabsContext = createContext<TabsContextValue | null>(null);

export function Tabs({ defaultValue, value, onValueChange, className = '', children }: {
  defaultValue?: string;
  value?: string;
  onValueChange?: (v: string) => void;
  className?: string;
  children?: React.ReactNode;
}) {
  const [internal, setInternal] = useState(defaultValue || '');
  const current = value ?? internal;
  const setValue = (v: string) => {
    if (onValueChange) onValueChange(v);
    else setInternal(v);
  };
  const ctx = useMemo(() => ({ value: current, setValue }), [current]);
  return (
    <TabsContext.Provider value={ctx}>
      <div className={className}>{children}</div>
    </TabsContext.Provider>
  );
}

export function TabsList({ className = '', children }: { className?: string; children?: React.ReactNode }) {
  return <div className={`flex gap-2 ${className}`}>{children}</div>;
}

export function TabsTrigger({ value, className = '', children }: { value: string; className?: string; children?: React.ReactNode }) {
  const ctx = useContext(TabsContext)!;
  const active = ctx.value === value;
  return (
    <button
      type="button"
      onClick={() => ctx.setValue(value)}
      className={`px-3 py-1.5 rounded-md border text-sm ${active ? 'bg-blue-600 text-white border-blue-600' : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'} ${className}`}
      aria-pressed={active}
    >
      {children}
    </button>
  );
}

export function TabsContent({ value, className = '', children }: { value: string; className?: string; children?: React.ReactNode }) {
  const ctx = useContext(TabsContext)!;
  if (ctx.value !== value) return null;
  return <div className={className}>{children}</div>;
}
