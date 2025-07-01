import * as React from 'react';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import {
  PageSection,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  Card,
  CardBody,
  Grid,
  GridItem,
} from '@patternfly/react-core';
import { TransactionList } from './list';
import { TagTotalTable } from './TagTotalTable';
import { TagByMonthChart } from './TagByMonthChart';
import { CubesIcon } from '@patternfly/react-icons';
import { Content, EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { TimePicker } from '@app/shared/components/time-picker';
import { calculateDateRange } from '@app/utils/dateUtils';

const Transactions: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const transactions = useAppSelector((state) => state.transactions);
  const { filteredTransactions } = useAppSelector((state) => state.transactionFilter);

  // Initialize with last 30 days default date range
  const defaultDateRange = React.useMemo(() => calculateDateRange('last 30 days'), []);
  const [dateRange, setDateRange] = React.useState<{ startDate: string; endDate: string }>({
    startDate: defaultDateRange.startDateValue,
    endDate: defaultDateRange.endDateValue,
  });

  React.useEffect(() => {
    // Fetch transactions when date range changes (backend filtering)
    // This includes initial load with default date range (last 30 days)
    dispatch(getTransactions(dateRange));
  }, [dateRange, dispatch]);

  const handleDateChange = (startDate: string, endDate: string) => {
    console.log('Date range changed:', { startDate, endDate });
    setDateRange({ startDate, endDate });
  };

  const emptyState = (
    <EmptyState variant={EmptyStateVariant.full} titleText="No transactions" icon={CubesIcon}>
      <EmptyStateBody>
        <Content>
          <Content component="p">Please add some transactions</Content>
        </Content>
      </EmptyStateBody>
    </EmptyState>
  );

  return (
    <PageSection hasBodyWrapper={false} style={{ backgroundColor: 'transparent' }}>
      {/* Tag Analytics Section */}
      <div style={{ padding: '1rem' }}>
        <Grid hasGutter>
          <GridItem span={6}>
            <TagByMonthChart
              transactions={filteredTransactions}
              startDate={dateRange.startDate}
              endDate={dateRange.endDate}
              title="Tag Totals by Month"
            />
          </GridItem>
          <GridItem span={6}>
            <TagTotalTable
              transactions={filteredTransactions}
              startDate={dateRange.startDate}
              endDate={dateRange.endDate}
            />
          </GridItem>
        </Grid>
      </div>

      {/* Main Content Card */}
      <div style={{ padding: '1rem' }}>
        <Card>
          <CardBody>
            {/* Date Picker Toolbar */}
            <Toolbar>
              <ToolbarContent>
                <ToolbarGroup>
                  <ToolbarItem>
                    <TimePicker
                      onDateChange={handleDateChange}
                      initialStartDate={dateRange.startDate}
                      initialEndDate={dateRange.endDate}
                      initialTimeRange="last 30 days"
                    />
                  </ToolbarItem>
                </ToolbarGroup>
              </ToolbarContent>
            </Toolbar>

            {/* Transaction List with integrated tag filter */}
            {transactions.transactions.length == 0 ? (
              emptyState
            ) : (
              <TransactionList transactions={transactions.transactions} />
            )}
          </CardBody>
        </Card>
      </div>
    </PageSection>
  );
};

export { Transactions };
