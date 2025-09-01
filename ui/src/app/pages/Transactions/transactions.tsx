import * as React from 'react';
import {
  EuiPageSection,
  EuiPanel,
  EuiFlexGroup,
  EuiFlexItem,
  EuiEmptyPrompt,
  EuiSpacer,
} from '@elastic/eui';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { calculateTransactionSummary } from './reducers/transactionSummary.reducer';
import { setDateRange } from './reducers/transaction-filter.reducer';
import { TransactionList } from './list';
import { TransactionSummary } from './TransactionSummary';
import { TransactionSummaryChart } from './TransactionSummaryChart';
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
    <EuiEmptyPrompt
      icon="cube"
      title={<h2>No transactions</h2>}
      body="Please add some transactions"
    />
  );

  return (
    <EuiPageSection style={{ backgroundColor: 'transparent' }}>
      <div style={{ padding: '1rem' }}>
        <h1 style={{ position: 'absolute', left: '-10000px' }}>Transactions</h1>
        {/* Transaction Summary and Chart */}
        {transactions.transactions.length > 0 && transactionSummary && (
          <EuiFlexGroup gutterSize="l">
            <EuiFlexItem grow={true}>
              <TransactionSummary transactionSummary={transactionSummary} />
            </EuiFlexItem>
            <EuiFlexItem grow={false} style={{ minWidth: '350px', maxWidth: '350px' }}>
              <TransactionSummaryChart transactionSummary={transactionSummary} />
            </EuiFlexItem>
          </EuiFlexGroup>
        )}

        {/* Main Transaction List Panel */}
        <EuiPanel paddingSize="l">
          {/* Date Picker */}
          <EuiFlexGroup gutterSize="s" alignItems="center">
            <EuiFlexItem grow={false}>
              <TimePicker
                onDateChange={handleDateChange}
                initialStartDate={dateRange.startDate}
                initialEndDate={dateRange.endDate}
                initialTimeRange={getInitialTimeRange()}
              />
            </EuiFlexItem>
          </EuiFlexGroup>
          
          <EuiSpacer size="m" />

          {/* Transaction List with integrated label filter */}
          {transactions.transactions.length === 0 ? (
            emptyState
          ) : (
            <TransactionList transactions={transactions.transactions} />
          )}
        </EuiPanel>
      </div>
    </EuiPageSection>
  );
};

export { Transactions };
