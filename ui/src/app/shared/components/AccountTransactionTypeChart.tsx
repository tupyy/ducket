import * as React from 'react';
import { Card, CardBody, Title, Spinner } from '@patternfly/react-core';
import { ChartDonutThreshold, ChartThemeColor, ChartLegend } from '@patternfly/react-charts/victory';
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

  // Group data by account for better visualization
  const chartData = React.useMemo(() => {
    const accountGroups: { [key: number]: { debit: number; credit: number } } = {};
    
    data.forEach((item) => {
      if (!accountGroups[item.account]) {
        accountGroups[item.account] = { debit: 0, credit: 0 };
      }
      accountGroups[item.account][item.type] = item.amount;
    });

    // Convert to chart format
    const result: Array<{ x: string; y: number; account: number; type: string }> = [];
    Object.entries(accountGroups).forEach(([account, amounts]) => {
      const accountNumber = parseInt(account, 10);
      if (amounts.debit > 0) {
        result.push({
          x: `Account ${accountNumber} - Debit`,
          y: amounts.debit,
          account: accountNumber,
          type: 'debit'
        });
      }
      if (amounts.credit > 0) {
        result.push({
          x: `Account ${accountNumber} - Credit`,
          y: amounts.credit,
          account: accountNumber,
          type: 'credit'
        });
      }
    });

    return result.sort((a, b) => a.account - b.account || a.type.localeCompare(b.type));
  }, [data]);

  // Create legend data with account and transaction type info
  const legendData = React.useMemo(() => {
    return chartData.map((item) => ({
      name: `${item.x}: ${item.y.toLocaleString('fr-FR', {
        style: 'currency',
        currency: 'EUR',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      })}`,
    }));
  }, [chartData]);

  // Calculate total for the donut threshold chart
  const totalAmount = React.useMemo(() => {
    return chartData.reduce((sum, item) => sum + item.y, 0);
  }, [chartData]);

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
        ) : chartData.length > 0 ? (
          <div
            style={{ height: '350px', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}
          >
            <ChartDonutThreshold
              ariaTitle="Debit vs Credit by account"
              ariaDesc="Donut chart showing debit and credit amounts grouped by account"
              data={chartData}
              height={330}
              legendComponent={
                <ChartLegend
                  data={legendData}
                  orientation="vertical"
                  style={{
                    labels: {
                      fill: theme === 'dark' ? '#ffffff' : undefined,
                      fontSize: 12,
                    },
                  }}
                />
              }
              legendOrientation="vertical"
              legendPosition="right"
              padding={{
                bottom: 20,
                left: 20,
                right: 200,
                top: 20,
              }}
              subTitle="Total Amount"
              title={`€${totalAmount.toLocaleString('fr-FR', {
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
              })}`}
              themeColor={ChartThemeColor.multiUnordered}
              width={650}
            />
          </div>
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