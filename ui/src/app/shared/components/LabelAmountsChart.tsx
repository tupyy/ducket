import * as React from 'react';
import { Card, CardBody, Title, Spinner } from '@patternfly/react-core';
import { ChartPie, ChartThemeColor, ChartLegend } from '@patternfly/react-charts/victory';
import { ITagReport } from '@app/shared/models/label';
import { useTheme } from '@app/shared/contexts/ThemeContext';

interface LabelAmountsChartProps {
  data: ITagReport[];
  loading?: boolean;
  title?: string;
}

const LabelAmountsChart: React.FC<LabelAmountsChartProps> = ({
  data,
  loading = false,
  title = 'Transaction Amounts by Label',
}) => {
  const { theme } = useTheme();
  // Convert LabelReport data to chart format
  const chartData = React.useMemo(() => {
    return data
      .map((item) => ({
        x: item.tag, // Keep as 'tag' for backward compatibility with existing interfaces
        y: item.amount,
      }))
      .sort((a, b) => b.y - a.y)
      .slice(0, 10); // Show top 10 labels
  }, [data]);

  // Create legend data from label amounts
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
              ariaTitle="Transaction amounts by label"
              ariaDesc="Pie chart showing total transaction amounts grouped by label"
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
              width={500}
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

export { LabelAmountsChart }; 