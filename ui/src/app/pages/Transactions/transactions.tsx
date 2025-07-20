import * as React from 'react';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
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
import { TransactionList } from './list';
import { CubesIcon } from '@patternfly/react-icons';
import { EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { TimePicker } from '@app/shared/components/time-picker';
import { calculateDateRange } from '@app/utils/dateUtils';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';

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

  // Calculate transaction summary based on current data
  const transactionSummary = React.useMemo(() => {
    // Always show totals by label key when no local filters are applied in the list component
    // We'll check if there are active filters by comparing source vs filtered transactions
    const hasActiveFilters = filteredTransactions.length !== transactions.transactions.length;
    
    if (hasActiveFilters && filteredTransactions.length > 0) {
      // Show totals of filtered transactions
      const totalAmount = filteredTransactions.reduce((sum, transaction) => sum + transaction.amount, 0);
      const transactionCount = filteredTransactions.length;
      const debitTotal = filteredTransactions
        .filter(t => t.kind === 'debit')
        .reduce((sum, t) => sum + Math.abs(t.amount), 0);
      const creditTotal = filteredTransactions
        .filter(t => t.kind === 'credit')
        .reduce((sum, t) => sum + Math.abs(t.amount), 0);
      
      return {
        type: 'filtered' as const,
        data: [
          { label: 'Total Transactions', count: transactionCount, amount: totalAmount },
          { label: 'Total Debits', count: filteredTransactions.filter(t => t.kind === 'debit').length, amount: -debitTotal },
          { label: 'Total Credits', count: filteredTransactions.filter(t => t.kind === 'credit').length, amount: creditTotal },
        ]
      };
    } else {
      // Show totals by label key only
      const labelKeyTotals: { [key: string]: { count: number; amount: number } } = {};
      
      transactions.transactions.forEach((transaction) => {
        transaction.labels.forEach((label) => {
          if (!labelKeyTotals[label.key]) {
            labelKeyTotals[label.key] = { count: 0, amount: 0 };
          }
          labelKeyTotals[label.key].count += 1;
          labelKeyTotals[label.key].amount += transaction.amount;
        });
      });
      
      const data = Object.entries(labelKeyTotals)
        .map(([key, totals]) => ({
          label: key,
          count: totals.count,
          amount: totals.amount
        }))
        .sort((a, b) => Math.abs(b.amount) - Math.abs(a.amount)); // Sort by absolute amount descending
      
      return {
        type: 'byLabelKey' as const,
        data
      };
    }
  }, [filteredTransactions, transactions.transactions]);

  React.useEffect(() => {
    // Fetch transactions when date range changes (backend filtering)
    // This includes initial load with default date range (last 30 days)
    dispatch(getTransactions({
      startDate: dateRange.startDate,
      endDate: dateRange.endDate
    }));
  }, [dateRange, dispatch]);

  const handleDateChange = (startDate: string, endDate: string) => {
    console.log('Date range changed:', { startDate, endDate });
    console.log('Previous dateRange:', dateRange);
    setDateRange({ startDate, endDate });
    console.log('New dateRange will be:', { startDate, endDate });
  };

  /**
   * Render the transaction summary table
   */
  const renderTransactionSummary = () => (
    <Card style={{ marginBottom: '1rem' }}>
      <CardBody>
        <Content>
          <strong>
            {transactionSummary.type === 'filtered' ? 'Filtered Transaction Summary' : 'Summary by Label Key'}
          </strong>
        </Content>
        <Table aria-label="transaction-summary" variant="compact" style={{ marginTop: '0.5rem' }}>
          <Thead>
            <Tr>
              <Th>{transactionSummary.type === 'filtered' ? 'Type' : 'Label Key'}</Th>
              <Th>Count</Th>
              <Th>Total Amount</Th>
            </Tr>
          </Thead>
          <Tbody>
            {transactionSummary.data.map((row, index) => (
              <Tr key={index}>
                <Td>
                  <Content>
                    <strong>{row.label}</strong>
                  </Content>
                </Td>
                <Td>
                  <Content>{row.count}</Content>
                </Td>
                <Td>
                  <Content 
                    style={{ 
                      color: row.amount >= 0 ? 'var(--pf-v6-global--success-color--100)' : 'var(--pf-v6-global--danger-color--100)',
                      fontWeight: 'bold'
                    }}
                  >
                    {Math.abs(row.amount).toLocaleString('de-DE', { 
                      minimumFractionDigits: 2, 
                      maximumFractionDigits: 2 
                    })}
                  </Content>
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </CardBody>
    </Card>
  );

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
        {transactions.transactions.length > 0 && renderTransactionSummary()}
        
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
                      initialTimeRange="last 30 days"
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
