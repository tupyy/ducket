import * as React from 'react';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { PageSection } from '@patternfly/react-core';
import { TransactionList } from './list';
import { CubesIcon } from '@patternfly/react-icons';
import { Content, EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';

const Transactions: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const transactions = useAppSelector((state) => state.transactions);

  React.useEffect(() => {
    dispatch(getTransactions());
  }, []);

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
      {transactions.transactions.length == 0 ? (
        emptyState
      ) : (
        <TransactionList transactions={transactions.transactions} />
      )}
    </PageSection>
  );
};

export { Transactions };
