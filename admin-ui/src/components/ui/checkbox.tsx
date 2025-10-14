import React from 'react';

export interface CheckboxProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange' | 'type'> {
  onCheckedChange?: (checked: boolean) => void;
}

export function Checkbox({ onCheckedChange, checked, className = '', ...props }: CheckboxProps) {
  return (
    <input
      type="checkbox"
      className={`h-4 w-4 border rounded ${className}`}
      checked={!!checked}
      onChange={(e) => onCheckedChange?.(e.currentTarget.checked)}
      {...props}
    />
  );
}
