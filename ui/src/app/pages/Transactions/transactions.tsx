import * as React from 'react';
import {
  Card,
  CardBody,
  PageSection,
  Title,
  Pagination,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  Button,
  Grid,
  GridItem,
  Content,
} from '@patternfly/react-core';
import { ITransaction } from '@app/shared/models/transaction';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { LabelByMonthChart } from './LabelByMonthChart';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { calculateTransactionSummary } from './reducers/transactionSummary.reducer';
import { setDateRange } from './reducers/transaction-filter.reducer';
import { TransactionList } from './list';
import { CubesIcon } from '@patternfly/react-icons';
import { EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { TimePicker } from '@app/shared/components/time-picker';
import { calculateDateRange } from '@app/utils/dateUtils';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';

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

  /**
   * Render the transaction summary table
   */
  const renderTransactionSummary = () => {
    if (!transactionSummary) return null;
    
    return (
      <Card style={{ marginBottom: '1rem' }}>
        <CardBody>
          <Content>
            <strong>{transactionSummary.type === 'filtered' ? 'Filtered Transaction Summary' : 'Summary'}</strong>
          </Content>
          <Table aria-label="transaction-summary" variant="compact" style={{ marginTop: '0.5rem' }}>
            <Thead>
              <Tr>
                <Th>{transactionSummary.type === 'filtered' ? 'Type' : 'Label'}</Th>
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
                      color:
                        row.debitAmount > 0
                          ? 'var(--pf-v6-global--danger-color--100)'
                          : 'var(--pf-v6-global--palette--black-600)',
                      fontWeight: 'bold',
                    }}
                  >
                    {row.debitAmount > 0
                      ? row.debitAmount.toLocaleString('de-DE', {
                          minimumFractionDigits: 2,
                          maximumFractionDigits: 2,
                        })
                      : '-'}
                  </Content>
                </Td>
                <Td>
                  <Content
                    style={{
                      color:
                        row.creditAmount > 0
                          ? 'var(--pf-v6-global--success-color--100)'
                          : 'var(--pf-v6-global--palette--black-600)',
                      fontWeight: 'bold',
                    }}
                  >
                    {row.creditAmount > 0
                      ? row.creditAmount.toLocaleString('de-DE', {
                          minimumFractionDigits: 2,
                          maximumFractionDigits: 2,
                        })
                      : '-'}
                  </Content>
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </CardBody>
    </Card>
    );
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
