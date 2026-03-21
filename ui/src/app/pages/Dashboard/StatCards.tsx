import * as React from 'react';
import { Card, CardTitle, CardBody, Grid, GridItem, Content } from '@patternfly/react-core';
import { SummaryOverview } from '@app/shared/reducers/dashboard.reducer';

interface StatCardsProps {
  overview: SummaryOverview;
}

const StatCards: React.FunctionComponent<StatCardsProps> = ({ overview }) => {
  const cards = [
    { title: 'Transactions', value: overview.total_transactions.toLocaleString() },
    { title: 'Total Credits', value: `€${overview.total_credit.toFixed(2)}` },
    { title: 'Total Debits', value: `€${overview.total_debit.toFixed(2)}` },
    { title: 'Net Balance', value: `€${overview.balance.toFixed(2)}` },
  ];

  return (
    <Grid hasGutter>
      {cards.map((card) => (
        <GridItem key={card.title} sm={12} md={6} lg={3}>
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
  );
};

export { StatCards };
