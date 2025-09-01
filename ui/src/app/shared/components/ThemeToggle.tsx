import * as React from 'react';
import { EuiButton, EuiButtonIcon } from '@elastic/eui';
import { useTheme } from '@app/shared/contexts/ThemeContext';

export interface ThemeToggleProps {
  variant?: 'primary' | 'success' | 'warning' | 'danger' | 'text';
  size?: 's' | 'm';
  iconOnly?: boolean;
}

export const ThemeToggle: React.FC<ThemeToggleProps> = ({ 
  variant = 'text', 
  size = 'm', 
  iconOnly = false 
}) => {
  const { theme, toggleTheme } = useTheme();

  const iconType = theme === 'light' ? 'moon' : 'sun';
  const label = `Switch to ${theme === 'light' ? 'dark' : 'light'} theme`;

  if (iconOnly) {
    return (
      <EuiButtonIcon
        iconType={iconType}
        onClick={toggleTheme}
        aria-label={label}
        size={size}
        color={variant}
      />
    );
  }

  return (
    <EuiButton
      iconType={iconType}
      onClick={toggleTheme}
      aria-label={label}
      size={size}
      color={variant}
      fill={false}
    >
      {theme === 'light' ? 'Dark' : 'Light'}
    </EuiButton>
  );
};