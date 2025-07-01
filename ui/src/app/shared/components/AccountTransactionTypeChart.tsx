import * as React from 'react';
import { Card, CardBody, Title, Spinner, Grid, GridItem } from '@patternfly/react-core';
import { ChartPie, ChartThemeColor, ChartLegend } from '@patternfly/react-charts/victory';
import { IAccountTransactionTypeReport } from '@app/shared/models/tag';
import { useTheme } from '@app/shared/contexts/ThemeContext';

interface AccountTransactionTypeChartProps {
  data: IAccountTransactionTypeReport[];
  loading?: boolean;
  title?: string;
}

const AccountTransactionTypeChart: React.FC<AccountTransactionTypeChartProps> = ({
  data,
  loading = false,
  title = 'Debit vs Credit by Account',
}) => {
  const { theme } = useTheme();

  // Group data by account for separate charts
  const accountChartData = React.useMemo(() => {
    const accountGroups: { [key: number]: { debit: number; credit: number } } = {};
    
    data.forEach((item) => {
      if (!accountGroups[item.account]) {
        accountGroups[item.account] = { debit: 0, credit: 0 };
      }
      accountGroups[item.account][item.type] = item.amount;
    });

    // Convert to array of account chart data
    return Object.entries(accountGroups)
      .map(([account, amounts]) => {
        const accountNumber = parseInt(account, 10);
        const chartData: Array<{ x: string; y: number }> = [];
        
        if (amounts.debit > 0) {
          chartData.push({
            x: 'Debit',
            y: amounts.debit
          });
        }
        
        if (amounts.credit > 0) {
          chartData.push({
            x: 'Credit',
            y: amounts.credit
          });
        }

        return {
          account: accountNumber,
          data: chartData,
          balance: amounts.credit - amounts.debit
        };
      })
      .sort((a, b) => a.account - b.account);
  }, [data]);

  return (
    <Card>
      <CardBody>
        <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
          {title}
        </Title>

        {loading ? (
          <div style={{ textAlign: 'center', padding: '2rem' }}>
            <Spinner size="lg" />
          </div>
        ) : accountChartData.length > 0 ? (
          <Grid hasGutter>
            {accountChartData.map((accountData) => (
              <GridItem 
                key={accountData.account} 
                span={accountChartData.length === 1 ? 12 : accountChartData.length === 2 ? 6 : 4}
              >
                <Card>
                  <CardBody>
                    <Title headingLevel="h4" size="md" style={{ marginBottom: '1rem', textAlign: 'center' }}>
                      Account {accountData.account}
                    </Title>
                    <div
                      style={{ 
                        height: '280px', 
                        width: '100%', 
                        display: 'flex', 
                        justifyContent: 'center', 
                        alignItems: 'center' 
                      }}
                    >
                      <ChartPie
                        ariaTitle={`Debit vs Credit for account ${accountData.account}`}
                        ariaDesc={`Pie chart showing debit and credit amounts for account ${accountData.account}`}
                        data={accountData.data}
                        height={260}
                        width={300}
                        legendComponent={
                          <ChartLegend
                            data={accountData.data.map(item => ({
                              name: `${item.x}: ${item.y.toLocaleString('fr-FR', {
                                style: 'currency',
                                currency: 'EUR',
                                minimumFractionDigits: 2,
                                maximumFractionDigits: 2,
                              })}`
                            }))}
                            orientation="vertical"
                            style={{
                              labels: {
                                fill: theme === 'dark' ? '#ffffff' : undefined,
                                fontSize: 11,
                              },
                            }}
                          />
                        }
                        legendOrientation="vertical"
                        legendPosition="bottom"
                        padding={{
                          bottom: 65,
                          left: 20,
                          right: 20,
                          top: 20,
                        }}
                        themeColor={ChartThemeColor.multiUnordered}
                      />
                    </div>
                    <div style={{ textAlign: 'center', marginTop: '0.5rem' }}>
                      <strong>
                        Balance: {accountData.balance.toLocaleString('fr-FR', {
                          style: 'currency',
                          currency: 'EUR',
                          minimumFractionDigits: 2,
                          maximumFractionDigits: 2,
                        })}
                      </strong>
                    </div>
                  </CardBody>
                </Card>
              </GridItem>
            ))}
          </Grid>
        ) : (
          <div style={{ textAlign: 'center', padding: '2rem', color: 'var(--pf-v6-global--Color--200)' }}>
            No transaction data available for the selected date range
          </div>
        )}
      </CardBody>
    </Card>
  );
};

export { AccountTransactionTypeChart }; 