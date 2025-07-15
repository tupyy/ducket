import * as React from 'react';
import { Card, CardBody, Title, Spinner } from '@patternfly/react-core';
import { ChartPie, ChartThemeColor, ChartLegend } from '@patternfly/react-charts/victory';
import { ITransactionTypeReport } from '@app/shared/models/label';
import { useTheme } from '@app/shared/contexts/ThemeContext';

interface TransactionTypeChartProps {
  data: ITransactionTypeReport[];
  loading?: boolean;
  title?: string;
}

const TransactionTypeChart: React.FC<TransactionTypeChartProps> = ({
  data,
  loading = false,
  title = 'Debit vs Credit Transactions',
}) => {
  const { theme } = useTheme();
  // Convert TransactionTypeReport data to chart format
  const chartData = React.useMemo(() => {
    return data.map((item) => ({
      x: item.type.charAt(0).toUpperCase() + item.type.slice(1), // Capitalize first letter
      y: item.amount,
    }));
  }, [data]);

  // Create legend data with transaction type info
  const legendData = React.useMemo(() => {
    return data.map((item) => ({
      name: `${item.type.charAt(0).toUpperCase() + item.type.slice(1)}: ${item.amount.toLocaleString('fr-FR', {
        style: 'currency',
        currency: 'EUR',
                  minimumFractionDigits: 2,
                  maximumFractionDigits: 2,
      })}`,
    }));
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
        ) : chartData.length > 0 ? (
          <div
            style={{ height: '300px', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}
          >
            <ChartPie
              ariaTitle="Debit vs Credit transaction amounts"
              ariaDesc="Pie chart showing total amounts for debit and credit transactions"
              data={chartData}
              height={280}
              legendComponent={
                <ChartLegend
                  data={legendData}
                  orientation="vertical"
                  style={{
                    labels: {
                      fill: theme === 'dark' ? '#ffffff' : undefined,
                    },
                  }}
                />
              }
              legendOrientation="vertical"
              legendPosition="right"
              padding={{
                bottom: 10,
                left: 10,
                right: 140,
                top: 10,
              }}
              themeColor={ChartThemeColor.multiUnordered}
              width={550}
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

export { TransactionTypeChart };
