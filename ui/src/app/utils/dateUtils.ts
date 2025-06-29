export const calculateDateRange = (timeRange: string) => {
  const now = new Date();
  const endDateValue = now.toISOString().split('T')[0]; // Format: YYYY-MM-DD
  let startDateValue = '';

  switch (timeRange) {
    case 'last 24 hours':
      startDateValue = new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case 'last 2 days':
      startDateValue = new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case 'last 7 days':
      startDateValue = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case 'last 30 days':
      startDateValue = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case 'last 90 days':
      startDateValue = new Date(now.getTime() - 90 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case 'last 6 months':
      const sixMonthsAgo = new Date(now);
      sixMonthsAgo.setMonth(now.getMonth() - 6);
      startDateValue = sixMonthsAgo.toISOString().split('T')[0];
      break;
    case 'last 1 year':
      const oneYearAgo = new Date(now);
      oneYearAgo.setFullYear(now.getFullYear() - 1);
      startDateValue = oneYearAgo.toISOString().split('T')[0];
      break;
    case 'last 2 years':
      const twoYearsAgo = new Date(now);
      twoYearsAgo.setFullYear(now.getFullYear() - 2);
      startDateValue = twoYearsAgo.toISOString().split('T')[0];
      break;
    default:
      startDateValue = endDateValue;
  }

  return { startDateValue, endDateValue };
};

export const getRelativeTimeRange = (startDate: string, endDate: string): string => {
  const timeRanges = [
    'last 24 hours',
    'last 2 days', 
    'last 7 days',
    'last 30 days',
    'last 90 days',
    'last 6 months',
    'last 1 year',
    'last 2 years'
  ];

  // Check if the provided dates match any of the predefined ranges
  for (const range of timeRanges) {
    const { startDateValue, endDateValue } = calculateDateRange(range);
    if (startDate === startDateValue && endDate === endDateValue) {
      return range;
    }
  }

  // If no match found, return custom range format
  const start = new Date(startDate);
  const end = new Date(endDate);
  
  return `${start.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })} - ${end.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric', 
    year: 'numeric',
  })}`;
};
