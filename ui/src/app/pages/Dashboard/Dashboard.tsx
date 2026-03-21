import * as React from 'react';
import {
  PageSection,
  Title,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  Spinner,
  Bullseye,
  Grid,
  GridItem,
  DatePicker,
  Flex,
  FlexItem,
  Content,
} from '@patternfly/react-core';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import {
  fetchOverview,
  fetchTagSummary,
  fetchBalanceTrend,
  fetchTopExpenses,
  setDateRange,
} from '@app/shared/reducers/dashboard.reducer';
import { StatCards } from './StatCards';
import { TagDonutChart } from './TagDonutChart';
import { BalanceTrendChart } from './BalanceTrendChart';
import { TopExpensesTable } from './TopExpensesTable';

const Dashboard: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const { overview, tagSummary, balanceTrend, topExpenses, loading, dateFrom, dateTo } =
    useAppSelector((state) => state.dashboard);

  const fetchAll = React.useCallback(() => {
    const params = { dateFrom, dateTo };
    dispatch(fetchOverview(params));
    dispatch(fetchTagSummary(params));
    dispatch(fetchBalanceTrend(params));
    dispatch(fetchTopExpenses(params));
  }, [dispatch, dateFrom, dateTo]);

  React.useEffect(() => {
    fetchAll();
  }, [fetchAll]);

  const handleDateFromChange = (_evt: React.FormEvent, value: string) => {
    dispatch(setDateRange({ dateFrom: value, dateTo }));
  };

  const handleDateToChange = (_evt: React.FormEvent, value: string) => {
    dispatch(setDateRange({ dateFrom, dateTo: value }));
  };

  return (
    <PageSection>
      <Flex justifyContent={{ default: 'justifyContentSpaceBetween' }} alignItems={{ default: 'alignItemsCenter' }} style={{ marginBottom: '1.5rem' }}>
        <FlexItem>
          <Title headingLevel="h1" size="lg">Dashboard</Title>
        </FlexItem>
        <FlexItem>
          <Flex gap={{ default: 'gapSm' }} alignItems={{ default: 'alignItemsCenter' }}>
            <FlexItem>
              <Content component="small">From</Content>
            </FlexItem>
            <FlexItem>
              <DatePicker
                value={dateFrom}
                onChange={handleDateFromChange}
                placeholder="YYYY-MM-DD"
                aria-label="Start date"
              />
            </FlexItem>
            <FlexItem>
              <Content component="small">To</Content>
            </FlexItem>
            <FlexItem>
              <DatePicker
                value={dateTo}
                onChange={handleDateToChange}
                placeholder="YYYY-MM-DD"
                aria-label="End date"
                rangeStart={dateFrom ? new Date(`${dateFrom}T00:00:00`) : undefined}
              />
            </FlexItem>
          </Flex>
        </FlexItem>
      </Flex>

      {loading && !overview ? (
        <Bullseye style={{ minHeight: '300px' }}>
          <Spinner size="xl" />
        </Bullseye>
      ) : (
        <>
          {overview && <StatCards overview={overview} />}

          <Grid hasGutter style={{ marginTop: '1.5rem' }}>
            <GridItem sm={12} lg={6}>
              <TagDonutChart data={tagSummary} />
            </GridItem>
            <GridItem sm={12} lg={6}>
              <BalanceTrendChart data={balanceTrend} />
            </GridItem>
          </Grid>

          <div style={{ marginTop: '1.5rem' }}>
            <TopExpensesTable transactions={topExpenses} />
          </div>
        </>
      )}
    </PageSection>
  );
};

export { Dashboard };
