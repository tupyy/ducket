// Helper function to get account label color based on account number
export const getAccountColor = (account: number): 'purple' | 'teal' | 'green' | 'orange' | 'yellow' | 'red' => {
  const colors: Array<'purple' | 'teal' | 'green' | 'orange' | 'yellow' | 'red'> = [
    'purple', 'teal', 'green', 'orange', 'yellow', 'red'
  ];
  if (account == 1000) {
    return 'green';
  }
  return colors[account % colors.length];
};

// Helper function to get dark theme color for account labels
export const getAccountDarkColor = (account: number): string => {
  const darkColors: string[] = [
    '#b19cd9', // purple
    '#009596', // teal
    '#3e8635', // green
    '#f4c430', // orange
    '#f1c21b', // yellow
    '#c9190b', // red
  ];
  return darkColors[account % darkColors.length];
};
