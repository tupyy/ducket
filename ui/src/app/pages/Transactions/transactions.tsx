import * as React from 'react';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { PageSection, Toolbar, ToolbarContent, ToolbarGroup, ToolbarItem } from '@patternfly/react-core';
import { TransactionList } from './list';
import { CubesIcon } from '@patternfly/react-icons';
import { Content, EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { TimePicker } from '@app/shared/components/time-picker';

const Transactions: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const transactions = useAppSelector((state) => state.transactions);
  const [dateRange, setDateRange] = React.useState<{ startDate: string; endDate: string } | null>(null);

  React.useEffect(() => {
    // Initial load without filters
    dispatch(getTransactions());
  }, [dispatch]);

  React.useEffect(() => {
    // Fetch transactions when date range changes (backend filtering)
    if (dateRange) {
      dispatch(getTransactions(dateRange));
    } else {
      // Reset to all transactions if no date range
      dispatch(getTransactions());
    }
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
    <PageSection hasBodyWrapper={false}>
      {/* Date Picker Toolbar */}
      <Toolbar>
        <ToolbarContent>
          <ToolbarGroup>
            <ToolbarItem>
              <TimePicker onDateChange={handleDateChange} />
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
    </PageSection>
  );
};

export { Transactions };
