import * as React from 'react';
import {
  PageSection,
  Title,
  Card,
  CardTitle,
  CardBody,
  Grid,
  GridItem,
  Spinner,
  Bullseye,
  Content,
} from '@patternfly/react-core';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';

const Dashboard: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const { transactions, loading } = useAppSelector((state) => state.transactions);

  React.useEffect(() => {
    dispatch(getTransactions({ limit: 1000 }));
  }, [dispatch]);

  const stats = React.useMemo(() => {
    const totalDebit = transactions
      .filter((t) => t.kind === 'debit')
      .reduce((sum, t) => sum + t.amount, 0);
    const totalCredit = transactions
      .filter((t) => t.kind === 'credit')
      .reduce((sum, t) => sum + t.amount, 0);
    const uniqueAccounts = new Set(transactions.map((t) => t.account)).size;
    const uniqueTags = new Set(transactions.flatMap((t) => t.tags)).size;

    return {
      totalTransactions: transactions.length,
      totalDebit,
      totalCredit,
      balance: totalCredit - totalDebit,
      uniqueAccounts,
      uniqueTags,
    };
  }, [transactions]);

  if (loading) {
    return (
      <PageSection>
        <Bullseye>
          <Spinner size="xl" />
        </Bullseye>
      </PageSection>
    );
  }

  const statCards = [
    { title: 'Transactions', value: stats.totalTransactions.toString() },
    { title: 'Total Credits', value: `€${stats.totalCredit.toFixed(2)}` },
    { title: 'Total Debits', value: `€${stats.totalDebit.toFixed(2)}` },
    { title: 'Net Balance', value: `€${stats.balance.toFixed(2)}` },
    { title: 'Accounts', value: stats.uniqueAccounts.toString() },
    { title: 'Tags', value: stats.uniqueTags.toString() },
  ];

  return (
    <PageSection>
      <Title headingLevel="h1" size="lg" style={{ marginBottom: '1.5rem' }}>
        Dashboard
      </Title>

      <Grid hasGutter>
        {statCards.map((card) => (
          <GridItem key={card.title} sm={12} md={6} lg={4} xl={2}>
            <Card isFullHeight>
              <CardTitle>{card.title}</CardTitle>
              <CardBody>
                <Content component="p" style={{ fontSize: '1.8rem', fontWeight: 'bold' }}>
                  {card.value}
                </Content>
              </CardBody>
            </Card>
          </GridItem>
        ))}
      </Grid>
    </PageSection>
  );
};

export { Dashboard };
