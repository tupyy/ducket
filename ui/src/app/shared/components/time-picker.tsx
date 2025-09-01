import * as React from 'react';
import {
  EuiSuperDatePicker,
  EuiFlexGroup,
  EuiFlexItem,
} from '@elastic/eui';
import moment from 'moment';

interface TimePickerProps {
  onDateChange?: (startDate: string, endDate: string) => void;
  initialStartDate?: string;
  initialEndDate?: string;
  initialTimeRange?: string;
}

const TimePicker: React.FC<TimePickerProps> = ({ 
  onDateChange, 
  initialStartDate, 
  initialEndDate, 
  initialTimeRange 
}) => {
  // Convert initial dates to EuiSuperDatePicker format
  const getInitialStart = () => {
    if (initialStartDate) {
      return moment(initialStartDate).toISOString();
    }
    return 'now-30d'; // Default to last 30 days
  };

  const getInitialEnd = () => {
    if (initialEndDate) {
      return moment(initialEndDate).toISOString();
    }
    return 'now'; // Default to now
  };

  const [start, setStart] = React.useState(getInitialStart());
  const [end, setEnd] = React.useState(getInitialEnd());
  const [isLoading, setIsLoading] = React.useState(false);

  const parseRelativeTime = (timeExpression: string): moment.Moment => {
    const now = moment();
    
    // Handle 'now' expressions
    if (timeExpression === 'now') {
      return now;
    }
    
    // Handle expressions like 'now/d' (start of day)
    if (timeExpression.includes('/d')) {
      const base = timeExpression.replace('/d', '');
      if (base === 'now') {
        return now.clone().startOf('day');
      }
      // Handle 'now-1d/d' (start of yesterday)
      if (base.includes('now-') && base.includes('d')) {
        const match = base.match(/now-(\d+)d/);
        if (match) {
          const days = parseInt(match[1]);
          return now.clone().subtract(days, 'days').startOf('day');
        }
      }
    }
    
    // Handle expressions like 'now-30d', 'now-1y', etc.
    if (timeExpression.includes('now-')) {
      const match = timeExpression.match(/now-(\d+)([dMy])/);
      if (match) {
        const value = parseInt(match[1]);
        const unit = match[2];
        
        switch (unit) {
          case 'd':
            return now.clone().subtract(value, 'days');
          case 'M':
            return now.clone().subtract(value, 'months');
          case 'y':
            return now.clone().subtract(value, 'years');
        }
      }
    }
    
    // If it's already a proper date/ISO string, parse it directly
    const parsed = moment(timeExpression);
    if (parsed.isValid()) {
      return parsed;
    }
    
    // Fallback to now if we can't parse
    return now;
  };

  const onTimeChange = ({ start: newStart, end: newEnd }: { start: string; end: string }) => {
    setStart(newStart);
    setEnd(newEnd);
    
    // Convert the date picker values to actual dates, handling relative expressions
    const startMoment = parseRelativeTime(newStart);
    const endMoment = parseRelativeTime(newEnd);
    
    const startDateString = startMoment.format('YYYY-MM-DD');
    const endDateString = endMoment.format('YYYY-MM-DD');
    
    onDateChange?.(startDateString, endDateString);
  };

  const onRefresh = ({ start: newStart, end: newEnd }: { start: string; end: string }) => {
    setIsLoading(true);
    onTimeChange({ start: newStart, end: newEnd });
    
    // Simulate loading for refresh
    setTimeout(() => {
      setIsLoading(false);
    }, 1000);
  };

  return (
    <EuiFlexGroup gutterSize="s" alignItems="center">
      <EuiFlexItem grow={false}>
        <EuiSuperDatePicker
          start={start}
          end={end}
          onTimeChange={onTimeChange}
          onRefresh={onRefresh}
          isLoading={isLoading}
          showUpdateButton="iconOnly"
          dateFormat="YYYY-MM-DD"
          canRoundRelativeUnits={false}
          recentlyUsedRanges={[
            { start: 'now-7d/d', end: 'now/d', label: 'Last 7 days' },
            { start: 'now-30d/d', end: 'now/d', label: 'Last 30 days' },
            { start: 'now-90d/d', end: 'now/d', label: 'Last 90 days' },
            { start: 'now-1y/d', end: 'now/d', label: 'Last year' },
          ]}
          commonlyUsedRanges={[
            { start: 'now/d', end: 'now/d', label: 'Today' },
            { start: 'now-1d/d', end: 'now-1d/d', label: 'Yesterday' },
            { start: 'now-7d/d', end: 'now/d', label: 'Last 7 days' },
            { start: 'now-30d/d', end: 'now/d', label: 'Last 30 days' },
            { start: 'now-90d/d', end: 'now/d', label: 'Last 90 days' },
            { start: 'now-1y/d', end: 'now/d', label: 'Last year' },
          ]}
        />
      </EuiFlexItem>
    </EuiFlexGroup>
  );
};

export { TimePicker };