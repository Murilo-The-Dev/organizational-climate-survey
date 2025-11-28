"use client";

import * as React from "react";

interface RadioGroupProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  value?: string;
  defaultValue?: string;
  onChange?: (value: string) => void;
  onValueChange?: (value: string) => void; // Adiciona suporte para onValueChange
}

export function RadioGroup({ 
  children, 
  value, 
  defaultValue,
  onChange, 
  onValueChange,
  ...props 
}: RadioGroupProps) {
  const [internalValue, setInternalValue] = React.useState(defaultValue || "");
  const currentValue = value ?? internalValue;

  const handleChange = (newValue: string) => {
    if (!value) {
      setInternalValue(newValue);
    }
    onChange?.(newValue);
    onValueChange?.(newValue); // Suporta ambas as APIs
  };

  return (
    <div {...props}>
      {React.Children.map(children, (child) =>
        React.isValidElement(child) 
          ? React.cloneElement(child as React.ReactElement<any>, { 
              checked: child.props.value === currentValue, 
              onChange: handleChange 
            }) 
          : child
      )}
    </div>
  );
}

interface RadioGroupItemProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  value: string;
  checked?: boolean;
  onChange?: (value: string) => void;
}

export const RadioGroupItem = React.forwardRef<HTMLInputElement, RadioGroupItemProps>(
  ({ value, checked, onChange, className, ...props }, ref) => {
    return (
      <input 
        type="radio" 
        value={value} 
        checked={checked}
        onChange={() => onChange?.(value)}
        className={className}
        {...props} 
        ref={ref} 
      />
    );
  }
);

RadioGroupItem.displayName = "RadioGroupItem";