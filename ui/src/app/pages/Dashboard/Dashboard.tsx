import * as React from 'react';
import {Card, PageSection, Title, CardBody, Grid, GridItem, Flex, FlexItem } from '@patternfly/react-core';
import { TimePicker } from '@app/shared/components/time-picker';
import { TagAmountsChart } from '@app/shared/components/TagAmountsChart';
import { TransactionTypeChart } from '@app/shared/components/TransactionTypeChart';
import { MonthlyTagChart } from '@app/shared/components/MonthlyTagChart';
import { TagFilter } from '@app/shared/components/tag-filter';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { calculateTagReport, calculateTransactionTypeReport } from '@app/shared/reducers/tag-report.reducer';
import { calculateMonthlyTagReport } from '@app/shared/reducers/monthly-tag-report.reducer';
import { calculateDateRange, getRelativeTimeRange } from '@app/utils/dateUtils';

const Dashboard: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const { transactions, loading } = useAppSelector((state) => state.transactions);
  const { tagReportData, transactionTypeData, loading: reportLoading } = useAppSelector((state) => state.tagReport);
  const { monthlyTagReports, loading: monthlyReportLoading } = useAppSelector((state) => state.monthlyTagReport);

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
      dispatch(calculateTagReport({ transactions, excludeCredits: true }));
      dispatch(calculateTransactionTypeReport(transactions));
      dispatch(calculateMonthlyTagReport({ transactions, excludeCredits: true }));
    }
  }, [dispatch, transactions]);

  // Get available tags from transactions for the filter
  const availableTags = React.useMemo(() => {
    const tagSet = new Set<string>();
    transactions.forEach((transaction) => {
      transaction.tags.forEach((tag) => {
        tagSet.add(tag.value);
      });
    });
    return Array.from(tagSet).sort();
  }, [transactions]);

  const handleTagFilterChange = (tags: string[]) => {
    setSelectedTag(tags);
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

        <GridItem span={6}>
          <TagAmountsChart data={tagReportData} loading={loading || reportLoading} />
        </GridItem>

        <GridItem span={6}>
          <TransactionTypeChart data={transactionTypeData} loading={loading || reportLoading} />
        </GridItem>

        <GridItem span={12}>
          <Card>
            <CardBody>
              <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
                Monthly Tag Trends
              </Title>
              <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsLg' }}>
                <FlexItem>
                  <TagFilter
                    availableTags={availableTags}
                    selectedTags={selectedTag}
                    onTagsChange={handleTagFilterChange}
                    placeholder="Select a tag to view monthly trends..."
                  />
                </FlexItem>
                <FlexItem>
                  <MonthlyTagChart
                    data={monthlyTagReports}
                    loading={loading || reportLoading || monthlyReportLoading}
                    tagNames={selectedTag}
                    startDate={startDate}
                    endDate={endDate}
                    title={
                      selectedTag.length > 0
                        ? `${selectedTag.join(', ')} - Monthly Trends`
                        : 'Select tags to view trends'
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
