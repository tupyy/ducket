import { useTheme } from '@app/shared/contexts/ThemeContext';

export interface ThemeStyles {
  backgroundColor: string;
  textColor: string;
  borderColor: string;
  cardBackground: string;
  headerBackground: string;
  footerBackground: string;
  primaryColor: string;
  secondaryColor: string;
  successColor: string;
  warningColor: string;
  dangerColor: string;
}

export const useThemeStyles = (): ThemeStyles & { theme: string } => {
  const { theme } = useTheme();

  const lightTheme: ThemeStyles = {
    backgroundColor: '#ffffff',
    textColor: '#151515',
    borderColor: '#d2d2d2',
    cardBackground: '#ffffff',
    headerBackground: '#ffffff',
    footerBackground: '#f8f9fa',
    primaryColor: '#0066cc',
    secondaryColor: '#6c757d',
    successColor: '#3e8635',
    warningColor: '#f0ab00',
    dangerColor: '#c9190b',
  };

  const darkTheme: ThemeStyles = {
    backgroundColor: '#151515',
    textColor: '#ffffff',
    borderColor: '#4f4f4f',
    cardBackground: '#212121',
    headerBackground: '#212121',
    footerBackground: '#1a1a1a',
    primaryColor: '#73bcf7',
    secondaryColor: '#a0a0a0',
    successColor: '#92d400',
    warningColor: '#f4c430',
    dangerColor: '#ff6b6b',
  };

  return {
    ...(theme === 'dark' ? darkTheme : lightTheme),
    theme,
  };
};
