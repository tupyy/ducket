import * as React from 'react';
import {
  Card,
  CardBody,
  PageSection,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  Content,
} from '@patternfly/react-core';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { calculateTransactionSummary } from './reducers/transactionSummary.reducer';
import { setDateRange } from './reducers/transaction-filter.reducer';
import { TransactionList } from './list';
import { TransactionSummary } from './TransactionSummary';
import { CubesIcon } from '@patternfly/react-icons';
import { EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { TimePicker } from '@app/shared/components/time-picker';
import { calculateDateRange } from '@app/utils/dateUtils';

const Transactions: React.FC = () => {
  const dispatch = useAppDispatch();
  const transactions = useAppSelector((state) => state.transactions);
  const { filteredTransactions, dateRange } = useAppSelector((state) => state.transactionFilter);
  const { summary: transactionSummary, loading: summaryLoading } = useAppSelector((state) => state.transactionSummary);

  React.useEffect(() => {
    // Calculate summary when transactions or filtered transactions change
    if (transactions.transactions.length > 0) {
      dispatch(
        calculateTransactionSummary({
          allTransactions: transactions.transactions,
          filteredTransactions: filteredTransactions,
        })
      );
    }
  }, [dispatch, transactions.transactions, filteredTransactions]);

  React.useEffect(() => {
    // Fetch transactions when date range changes (backend filtering)
    // This includes initial load with the date range from reducer state
    dispatch(
      getTransactions({
        startDate: dateRange.startDate,
        endDate: dateRange.endDate,
      })
    );
  }, [dateRange, dispatch]);

  const handleDateChange = (startDate: string, endDate: string) => {
    console.log('Date range changed:', { startDate, endDate });
    console.log('Previous dateRange:', dateRange);
    dispatch(setDateRange({ startDate, endDate }));
    console.log('New dateRange will be:', { startDate, endDate });
  };

  // Determine initial time range for TimePicker component
  const getInitialTimeRange = () => {
    const defaultDateRange = calculateDateRange('last 30 days');
    if (
      dateRange.startDate === defaultDateRange.startDateValue &&
      dateRange.endDate === defaultDateRange.endDateValue
    ) {
      return 'last 30 days';
    }
    return 'custom';
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
      <div style={{ padding: '1rem' }}>
        {/* Transaction Summary Card */}
        {transactions.transactions.length > 0 && transactionSummary && (
          <TransactionSummary transactionSummary={transactionSummary} />
        )}

        {/* Main Transaction List Card */}
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
                      initialTimeRange={getInitialTimeRange()}
                    />
                  </ToolbarItem>
                </ToolbarGroup>
              </ToolbarContent>
            </Toolbar>

            {/* Transaction List with integrated label filter */}
            {transactions.transactions.length === 0 ? (
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
