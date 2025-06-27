import * as React from 'react';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { 
  PageSection, 
  Toolbar, 
  ToolbarContent, 
  ToolbarItem, 
  ToolbarGroup 
} from '@patternfly/react-core';
import { TransactionList } from './list';
import { CubesIcon } from '@patternfly/react-icons';
import { Content, EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { TimePicker } from '@app/shared/components/time-picker';

const Transactions: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const transactions = useAppSelector((state) => state.transactions);
  const [dateRange, setDateRange] = React.useState<{ startDate: string; endDate: string } | null>(null);

  React.useEffect(() => {
    // Initial load without date filters
    dispatch(getTransactions());
  }, []);

  React.useEffect(() => {
    // Fetch transactions when date range changes
    if (dateRange) {
      dispatch(getTransactions(dateRange));
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
      <Toolbar>
        <ToolbarContent>
          <ToolbarGroup>
            <ToolbarItem>
              <TimePicker onDateChange={handleDateChange} />
            </ToolbarItem>
          </ToolbarGroup>
        </ToolbarContent>
      </Toolbar>
      {transactions.transactions.length == 0 ? (
        emptyState
      ) : (
        <TransactionList transactions={transactions.transactions} />
      )}
    </PageSection>
  );
};

export { Transactions };
