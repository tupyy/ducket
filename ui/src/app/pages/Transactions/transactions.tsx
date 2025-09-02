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
import { IHierarchicalSummary } from './utils/calculateSummary';
import { calculateHierarchicalSummary } from './utils/calculateSummary';
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
  
  // State for expanded summary keys
  const [expandedSummaryKeys, setExpandedSummaryKeys] = React.useState<Set<string>>(new Set());

  // Calculate summary locally instead of using Redux
  const transactionSummary = React.useMemo(() => {
    if (transactions.transactions.length > 0 && filteredTransactions.length > 0) {
      return calculateHierarchicalSummary(filteredTransactions);
    }
    return null;
  }, [transactions.transactions, filteredTransactions]);

  // Create filtered summary for chart based on expanded keys
  const chartSummary = React.useMemo(() => {
    if (!transactionSummary) return null;
    
    const hasAnyExpanded = expandedSummaryKeys.size > 0;
    
    if (hasAnyExpanded) {
      // Show only the children of expanded keys
      const filteredData = transactionSummary.data
        .filter(parentRow => expandedSummaryKeys.has(parentRow.label))
        .flatMap(parentRow => parentRow.children || []);
      
      const filteredTotals = filteredData.reduce(
        (acc, row) => ({
          count: acc.count + row.count,
          debitAmount: acc.debitAmount + row.debitAmount,
          creditAmount: acc.creditAmount + row.creditAmount,
        }),
        { count: 0, debitAmount: 0, creditAmount: 0 }
      );

      return {
        type: 'hierarchical' as const,
        data: filteredData,
        totals: filteredTotals,
      };
    }
    
    return transactionSummary;
  }, [transactionSummary, expandedSummaryKeys]);


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
        {transactions.transactions.length > 0 && transactionSummary && chartSummary && (
          <EuiFlexGroup gutterSize="l">
            <EuiFlexItem grow={3}>
              <TransactionSummary 
                transactionSummary={transactionSummary}
                expandedKeys={expandedSummaryKeys}
                setExpandedKeys={setExpandedSummaryKeys}
              />
            </EuiFlexItem>
            <EuiFlexItem grow={1}>
              <TransactionSummaryChart transactionSummary={chartSummary} />
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
