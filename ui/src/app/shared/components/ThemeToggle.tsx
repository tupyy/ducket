import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { MoonIcon, SunIcon } from '@patternfly/react-icons';
import { useTheme } from '@app/shared/contexts/ThemeContext';

export interface ThemeToggleProps {
  variant?: 'primary' | 'secondary' | 'tertiary' | 'danger' | 'warning' | 'link' | 'plain';
  size?: 'sm' | 'lg' | 'default';
  isInline?: boolean;
}

export const ThemeToggle: React.FC<ThemeToggleProps> = ({ variant = 'plain', size = 'default', isInline = false }) => {
  const { theme, toggleTheme } = useTheme();

  return (
    <Button
      variant={variant}
      size={size}
      isInline={isInline}
      onClick={toggleTheme}
      aria-label={`Switch to ${theme === 'light' ? 'dark' : 'light'} theme`}
      icon={theme === 'light' ? <MoonIcon /> : <SunIcon />}
    >
      {theme === 'light' ? 'Dark' : 'Light'}
    </Button>
  );
};
