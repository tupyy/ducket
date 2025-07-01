import * as React from 'react';
import {
  DatePicker,
  Grid,
  GridItem,
  Dropdown,
  DropdownList,
  DropdownItem,
  MenuToggle,
  MenuToggleElement,
} from '@patternfly/react-core';
import { calculateDateRange, getRelativeTimeRange } from '@app/utils/dateUtils';

interface TimePickerProps {
  onDateChange?: (startDate: string, endDate: string) => void;
  initialStartDate?: string;
  initialEndDate?: string;
  initialTimeRange?: string;
}

const timeList = [
  'last 24 hours',
  'last 2 days',
  'last 7 days',
  'last 30 days',
  'last 90 days',
  'last 6 months',
  'last 1 year',
  'last 2 years',
];

const css: React.CSSProperties = {
  gridTemplateColumns: 'repeat(3, 1fr)',
};

const TimePicker: React.FC<TimePickerProps> = ({ onDateChange, initialStartDate, initialEndDate, initialTimeRange }) => {
  const getFirstDayOfMonth = () => {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, '0'); // +1 because getMonth() is 0-based
    return `${year}-${month}-01`;
  };

  const getTodayDate = () => {
    return new Date().toISOString().split('T')[0];
  };

  const [startDate, setStartDate] = React.useState<string>(initialStartDate || getFirstDayOfMonth());
  const [endDate, setEndDate] = React.useState<string>(initialEndDate || getTodayDate());
  const [isDropdownOpen, setIsDropdownOpen] = React.useState(false);
  
  // Initialize selectedTimeRange based on initial dates or provided range
  const getInitialTimeRange = () => {
    if (initialTimeRange) {
      return initialTimeRange;
    }
    if (initialStartDate && initialEndDate) {
      const detectedRange = getRelativeTimeRange(initialStartDate, initialEndDate);
      // Check if it matches one of our predefined ranges
      return timeList.includes(detectedRange) ? detectedRange : 'Custom range';
    }
    return 'Select time range';
  };

  // Helper function to update time range display based on current dates
  const updateTimeRangeDisplay = (start: string, end: string) => {
    const detectedRange = getRelativeTimeRange(start, end);
    if (timeList.includes(detectedRange)) {
      setSelectedTimeRange(detectedRange);
    } else {
      setSelectedTimeRange('Custom range');
    }
  };
  
  const [selectedTimeRange, setSelectedTimeRange] = React.useState<string>(getInitialTimeRange());

  React.useEffect(() => {
    // Trigger callback with initial values on mount
    onDateChange?.(startDate, endDate);
  }, []); // Only run on mount to avoid infinite loops

  const handleStartDateChange = (_event: any, value: string) => {
    setStartDate(value);

    // If start date is after end date, update end date to match start date
    if (value && endDate && new Date(value) > new Date(endDate)) {
      setEndDate(value);
      updateTimeRangeDisplay(value, value);
      onDateChange?.(value, value);
    } else {
      updateTimeRangeDisplay(value, endDate);
      onDateChange?.(value, endDate);
    }
  };

  const handleEndDateChange = (_event: any, value: string) => {
    setEndDate(value);
    updateTimeRangeDisplay(startDate, value);
    onDateChange?.(startDate, value);
  };

  const validateEndDate = (date: Date): string => {
    if (startDate && date < new Date(startDate)) {
      return 'End date must be after start date';
    }
    return '';
  };

  const handleTimeRangeClick = (timeRange: string) => {
    const { startDateValue, endDateValue } = calculateDateRange(timeRange);
    setStartDate(startDateValue);
    setEndDate(endDateValue);
    setSelectedTimeRange(timeRange);
    setIsDropdownOpen(false);
    onDateChange?.(startDateValue, endDateValue);
  };

  const onToggleClick = () => {
    setIsDropdownOpen(!isDropdownOpen);
  };

  return (
    <React.Fragment>
      <Grid style={css} hasGutter>
        <GridItem span={1}>
          <DatePicker
            key={`start-${startDate}`}
            value={startDate}
            onChange={handleStartDateChange}
            placeholder="Start date"
          />
        </GridItem>
        <GridItem span={1}>
          <DatePicker
            value={endDate}
            onChange={handleEndDateChange}
            placeholder="End date"
            validators={[validateEndDate]}
            rangeStart={startDate ? new Date(startDate) : undefined}
            isDisabled={!startDate}
          />
        </GridItem>
        <GridItem span={1}>
          <Dropdown
            isOpen={isDropdownOpen}
            onOpenChange={(isOpen: boolean) => setIsDropdownOpen(isOpen)}
            toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
              <MenuToggle ref={toggleRef} onClick={onToggleClick} isExpanded={isDropdownOpen}>
                {selectedTimeRange}
              </MenuToggle>
            )}
            ouiaId="TimeRangeDropdown"
            shouldFocusToggleOnSelect
          >
            <DropdownList>
              {timeList.map((item: string, idx: number) => (
                <DropdownItem key={`${idx}`} onClick={() => handleTimeRangeClick(item)}>
                  {item}
                </DropdownItem>
              ))}
            </DropdownList>
          </Dropdown>
        </GridItem>
      </Grid>
    </React.Fragment>
  );
};

export { TimePicker };
