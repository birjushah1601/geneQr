"use client";
import React, { createContext, useContext, useMemo, useRef, useState } from 'react';

type Ctx = {
  value?: string;
  onValueChange?: (v: string) => void;
};

const SelectCtx = createContext<Ctx>({});

export function Select({ value, onValueChange, children }: { value?: string; onValueChange?: (v: string) => void; children?: React.ReactNode }) {
  const ctx = useMemo(() => ({ value, onValueChange }), [value, onValueChange]);
  return <SelectCtx.Provider value={ctx}>{children}</SelectCtx.Provider>;
}

export function SelectTrigger({ id, className = '', children, onClick }: React.HTMLAttributes<HTMLButtonElement> & { id?: string }) {
  return (
    <button id={id} type="button" className={`w-full h-10 px-3 border rounded-md text-left ${className}`} onClick={onClick}>
      {children}
    </button>
  );
}

export function SelectValue() {
  const { value } = useContext(SelectCtx);
  return <span>{value || 'Select…'}</span>;
}

export function SelectContent({ className = '', children }: { className?: string; children?: React.ReactNode }) {
  return <div className={`mt-2 border rounded-md bg-white shadow-sm p-1 ${className}`}>{children}</div>;
}

export function SelectItem({ value, children, className = '' }: { value: string; children?: React.ReactNode; className?: string }) {
  const { onValueChange } = useContext(SelectCtx);
  return (
    <div
      role="option"
      className={`px-3 py-2 rounded cursor-pointer hover:bg-gray-100 ${className}`}
      onClick={() => onValueChange && onValueChange(value)}
    >
      {children}
    </div>
  );
}
