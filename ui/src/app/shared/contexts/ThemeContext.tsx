import * as React from 'react';
import { EuiProvider } from '@elastic/eui';
import { EuiThemeColorMode } from '@elastic/eui/src/services/theme';

export type ThemeType = 'light' | 'dark';
export type EuiColorMode = EuiThemeColorMode;

export interface ThemeContextType {
  theme: ThemeType;
  setTheme: (theme: ThemeType) => void;
  toggleTheme: () => void;
}

const ThemeContext = React.createContext<ThemeContextType | undefined>(undefined);

export interface ThemeProviderProps {
  children: React.ReactNode;
  defaultTheme?: ThemeType;
}

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children, defaultTheme = 'light' }) => {
  const [theme, setTheme] = React.useState<ThemeType>(() => {
    // Check localStorage for saved theme preference
    if (typeof window !== 'undefined') {
      const savedTheme = localStorage.getItem('finante-theme') as ThemeType;
      return savedTheme || defaultTheme;
    }
    return defaultTheme;
  });

  // Save theme to localStorage whenever it changes
  React.useEffect(() => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('finante-theme', theme);
    }
  }, [theme]);

  const toggleTheme = React.useCallback(() => {
    setTheme((prevTheme) => (prevTheme === 'light' ? 'dark' : 'light'));
  }, []);

  const contextValue = React.useMemo(
    () => ({
      theme,
      setTheme,
      toggleTheme,
    }),
    [theme, toggleTheme],
  );

  return (
    <ThemeContext.Provider value={contextValue}>
      <EuiProvider colorMode={theme}>
        {children}
      </EuiProvider>
    </ThemeContext.Provider>
  );
};

export const useTheme = (): ThemeContextType => {
  const context = React.useContext(ThemeContext);
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
};

export { ThemeContext };
