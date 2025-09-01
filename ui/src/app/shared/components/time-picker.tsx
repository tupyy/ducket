import * as React from 'react';
import {
  EuiDatePicker,
  EuiFlexGroup,
  EuiFlexItem,
  EuiSelect,
  EuiSelectOption,
} from '@elastic/eui';
import { calculateDateRange, getRelativeTimeRange } from '@app/utils/dateUtils';
import moment from 'moment';

interface TimePickerProps {
  onDateChange?: (startDate: string, endDate: string) => void;
  initialStartDate?: string;
  initialEndDate?: string;
  initialTimeRange?: string;
}

const timeList = [
  { value: 'last 24 hours', text: 'Last 24 hours' },
  { value: 'last 2 days', text: 'Last 2 days' },
  { value: 'last 7 days', text: 'Last 7 days' },
  { value: 'last 30 days', text: 'Last 30 days' },
  { value: 'last 90 days', text: 'Last 90 days' },
  { value: 'last 6 months', text: 'Last 6 months' },
  { value: 'last 1 year', text: 'Last 1 year' },
  { value: 'last 2 years', text: 'Last 2 years' },
  { value: 'custom', text: 'Custom range' },
];

const TimePicker: React.FC<TimePickerProps> = ({ 
  onDateChange, 
  initialStartDate, 
  initialEndDate, 
  initialTimeRange 
}) => {
  const getFirstDayOfMonth = () => {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, '0');
    return `${year}-${month}-01`;
  };

  const getLastDayOfMonth = () => {
    const today = new Date();
    const year = today.getFullYear();
    const month = today.getMonth();
    const lastDay = new Date(year, month + 1, 0).getDate();
    const monthStr = String(month + 1).padStart(2, '0');
    return `${year}-${monthStr}-${String(lastDay).padStart(2, '0')}`;
  };

  const [selectedTimeRange, setSelectedTimeRange] = React.useState<string>(
    initialTimeRange || 'last 30 days'
  );
  const [startDate, setStartDate] = React.useState<moment.Moment | null>(
    initialStartDate ? moment(initialStartDate) : moment(getFirstDayOfMonth())
  );
  const [endDate, setEndDate] = React.useState<moment.Moment | null>(
    initialEndDate ? moment(initialEndDate) : moment(getLastDayOfMonth())
  );

  const handleTimeRangeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    setSelectedTimeRange(value);
    
    if (value !== 'custom') {
      const { startDate: calcStartDate, endDate: calcEndDate } = getRelativeTimeRange(value);
      setStartDate(moment(calcStartDate));
      setEndDate(moment(calcEndDate));
      onDateChange?.(calcStartDate, calcEndDate);
    }
  };

  const handleStartDateChange = (date: moment.Moment | null) => {
    setStartDate(date);
    if (date && endDate) {
      onDateChange?.(date.format('YYYY-MM-DD'), endDate.format('YYYY-MM-DD'));
    }
  };

  const handleEndDateChange = (date: moment.Moment | null) => {
    setEndDate(date);
    if (startDate && date) {
      onDateChange?.(startDate.format('YYYY-MM-DD'), date.format('YYYY-MM-DD'));
    }
  };

  const showCustomDates = selectedTimeRange === 'custom';

  return (
    <EuiFlexGroup gutterSize="s" alignItems="center">
      <EuiFlexItem grow={false}>
        <EuiSelect
          value={selectedTimeRange}
          onChange={handleTimeRangeChange}
          options={timeList}
          compressed
        />
      </EuiFlexItem>
      
      {showCustomDates && (
        <>
          <EuiFlexItem grow={false}>
            <EuiDatePicker
              selected={startDate}
              onChange={handleStartDateChange}
              maxDate={endDate || moment()}
              dateFormat="YYYY-MM-DD"
              compressed
            />
          </EuiFlexItem>
          
          <EuiFlexItem grow={false}>
            <EuiDatePicker
              selected={endDate}
              onChange={handleEndDateChange}
              minDate={startDate}
              maxDate={moment()}
              dateFormat="YYYY-MM-DD"
              compressed
            />
          </EuiFlexItem>
        </>
      )}
    </EuiFlexGroup>
  );
};

export { TimePicker };