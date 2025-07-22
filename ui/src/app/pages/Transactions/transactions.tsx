import * as React from 'react';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { setDateRange } from '@app/shared/reducers/transaction-filter.reducer';
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
  const { filteredTransactions, dateRange } = useAppSelector((state) => state.transactionFilter);

  // Calculate transaction summary based on current data
  const transactionSummary = React.useMemo(() => {
    // Always show totals by label key when no local filters are applied in the list component
    // We'll check if there are active filters by comparing source vs filtered transactions
    const hasActiveFilters = filteredTransactions.length !== transactions.transactions.length;

    if (hasActiveFilters && filteredTransactions.length > 0) {
      // Show totals of filtered transactions
      const transactionCount = filteredTransactions.length;
      const debitTransactions = filteredTransactions.filter(t => t.kind === 'credit');
      const creditTransactions = filteredTransactions.filter(t => t.kind === 'debit');
      const debitTotal = debitTransactions.reduce((sum, t) => sum + Math.abs(t.amount), 0);
      const creditTotal = creditTransactions.reduce((sum, t) => sum + Math.abs(t.amount), 0);

      return {
        type: 'filtered' as const,
        data: [
          {
            label: 'Total Transactions',
            count: transactionCount,
            debitAmount: debitTotal,
            creditAmount: creditTotal
          },
        ]
      };
    } else {
      // Show totals by label key only
      const labelKeyTotals: { [key: string]: { count: number; debitAmount: number; creditAmount: number } } = {};

      transactions.transactions.forEach((transaction) => {
        transaction.labels.forEach((label) => {
          if (!labelKeyTotals[label.key]) {
            labelKeyTotals[label.key] = { count: 0, debitAmount: 0, creditAmount: 0 };
          }
          labelKeyTotals[label.key].count += 1;

          if (transaction.kind === 'credit') {
            labelKeyTotals[label.key].debitAmount += Math.abs(transaction.amount);
          } else if (transaction.kind === 'debit') {
            labelKeyTotals[label.key].creditAmount += Math.abs(transaction.amount);
          }
        });
      });

      const data = Object.entries(labelKeyTotals)
        .map(([key, totals]) => ({
          label: key,
          count: totals.count,
          debitAmount: totals.debitAmount,
          creditAmount: totals.creditAmount
        }))
        .sort((a, b) => (Math.abs(b.debitAmount) + Math.abs(b.creditAmount)) - (Math.abs(a.debitAmount) + Math.abs(a.creditAmount))); // Sort by total absolute amount descending

      return {
        type: 'byLabelKey' as const,
        data
      };
    }
  }, [filteredTransactions, transactions.transactions]);

  React.useEffect(() => {
    // Fetch transactions when date range changes (backend filtering)
    // This includes initial load with the date range from reducer state
    dispatch(getTransactions({
      startDate: dateRange.startDate,
      endDate: dateRange.endDate
    }));
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
    if (dateRange.startDate === defaultDateRange.startDateValue &&
        dateRange.endDate === defaultDateRange.endDateValue) {
      return 'last 30 days';
    }
    return 'custom';
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
              <Th>Debit</Th>
              <Th>Credit</Th>
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
                      color: row.debitAmount > 0 ? 'var(--pf-v6-global--danger-color--100)' : 'var(--pf-v6-global--palette--black-600)',
                      fontWeight: 'bold'
                    }}
                  >
                    {row.debitAmount > 0 ? row.debitAmount.toLocaleString('de-DE', {
                      minimumFractionDigits: 2,
                      maximumFractionDigits: 2
                    }) : '-'}
                  </Content>
                </Td>
                <Td>
                  <Content
                    style={{
                      color: row.creditAmount > 0 ? 'var(--pf-v6-global--success-color--100)' : 'var(--pf-v6-global--palette--black-600)',
                      fontWeight: 'bold'
                    }}
                  >
                    {row.creditAmount > 0 ? row.creditAmount.toLocaleString('de-DE', {
                      minimumFractionDigits: 2,
                      maximumFractionDigits: 2
                    }) : '-'}
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
