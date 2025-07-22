import * as React from 'react';
import { PageSection, Title, Grid, GridItem, Card, CardBody, Flex, FlexItem } from '@patternfly/react-core';
import { LabelAmountsChart } from '@app/shared/components/LabelAmountsChart';
import { MonthlyLabelChart } from '@app/shared/components/MonthlyLabelChart';
import { LabelFilter } from '@app/shared/components/label-filter';
import { TimePicker } from '@app/shared/components/time-picker';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { calculateLabelReport } from '@app/shared/reducers/label-report.reducer';
import { calculateMonthlyLabelReport } from '@app/shared/reducers/monthly-label-report.reducer';
import { calculateDateRange, getRelativeTimeRange } from '@app/utils/dateUtils';

const Dashboard: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const { transactions, loading } = useAppSelector((state) => state.transactions);
  const { labelReportData, loading: reportLoading } = useAppSelector((state) => state.labelReport);
  const { monthlyLabelReports, loading: monthlyReportLoading } = useAppSelector((state) => state.monthlyLabelReport);

  // Initialize with 1 year default date range
  const defaultDateRange = React.useMemo(() => calculateDateRange('last 1 year'), []);
  const [startDate, setStartDate] = React.useState<string>(defaultDateRange.startDateValue);
  const [endDate, setEndDate] = React.useState<string>(defaultDateRange.endDateValue);
  const [selectedTag, setSelectedTag] = React.useState<string[]>([]);

  const handleDateChange = (start: string, end: string) => {
    setStartDate(start);
    setEndDate(end);
    // Fetch transactions with date filter
    dispatch(getTransactions({ startDate: start, endDate: end }));
  };

  // Calculate reports when transactions change
  React.useEffect(() => {
    if (transactions.length > 0) {
      dispatch(calculateLabelReport({ transactions, excludeCredits: true }));
      dispatch(calculateMonthlyLabelReport({ transactions, excludeCredits: true }));
    }
  }, [dispatch, transactions]);

  // Get available labels from transactions for the filter
  const availableLabels = React.useMemo(() => {
    const labelKeySet = new Set<string>();
    transactions.forEach((transaction) => {
      transaction.labels.forEach((label) => {
        labelKeySet.add(label.key);
      });
    });
    return Array.from(labelKeySet).sort();
  }, [transactions]);

  const handleLabelFilterChange = (labels: string[]) => {
    setSelectedTag(labels);
  };

  // Fetch initial data on component mount with default date range
  React.useEffect(() => {
    dispatch(getTransactions({ startDate, endDate }));
  }, [dispatch, startDate, endDate]);

  return (
    <PageSection hasBodyWrapper={false}>
      <Title headingLevel="h1" size="lg">
        Dashboard
      </Title>

      <Grid hasGutter style={{ marginTop: '1rem' }}>
        <GridItem span={12}>
          <Flex justifyContent={{ default: 'justifyContentFlexStart' }}>
            <FlexItem>
              <TimePicker 
                onDateChange={handleDateChange}
                initialStartDate={startDate}
                initialEndDate={endDate}
              />
            </FlexItem>
          </Flex>

          {startDate && endDate && (
            <div
              style={{
                marginTop: '1rem',
                padding: '1rem',
                backgroundColor: 'var(--pf-v6-global--BackgroundColor--200)',
                borderRadius: '4px',
              }}
            >
              <strong>Selected Range:</strong>{' '}
              {getRelativeTimeRange(startDate, endDate)}
            </div>
          )}
        </GridItem>

        <GridItem span={12}>
          <LabelAmountsChart data={labelReportData} loading={loading || reportLoading} />
        </GridItem>

        <GridItem span={12}>
          <Card>
            <CardBody>
              <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
                Monthly Label Trends
              </Title>
              <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsLg' }}>
                <FlexItem>
                  <LabelFilter
                    availableLabels={availableLabels}
                    selectedLabels={selectedTag}
                    onLabelsChange={handleLabelFilterChange}
                    placeholder="Select label keys to view monthly trends..."
                  />
                </FlexItem>
                <FlexItem>
                  <MonthlyLabelChart
                    data={monthlyLabelReports}
                    loading={loading || reportLoading || monthlyReportLoading}
                    labelNames={selectedTag}
                    startDate={startDate}
                    endDate={endDate}
                    title={
                      selectedTag.length > 0
                        ? `${selectedTag.join(', ')} - Monthly Trends`
                        : 'Select labels to view trends'
                    }
                  />
                </FlexItem>
              </Flex>
            </CardBody>
          </Card>
        </GridItem>
      </Grid>
    </PageSection>
  );
};

export { Dashboard };
