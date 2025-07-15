import * as React from 'react';
import { Card, CardBody, Title, Spinner } from '@patternfly/react-core';
import {
  Chart,
  ChartAxis,
  ChartBar,
  ChartGroup,
  ChartThemeColor,
  ChartLegend,
  ChartVoronoiContainer,
} from '@patternfly/react-charts/victory';
import { ITransaction } from '@app/shared/models/transaction';
import { useTheme } from '@app/shared/contexts/ThemeContext';

interface LabelByMonthChartProps {
  transactions: ITransaction[];
  startDate?: string;
  endDate?: string;
  title?: string;
}

interface MonthlyLabelData {
  label: string;
  month: string;
  amount: number;
}

const LabelByMonthChart: React.FC<LabelByMonthChartProps> = ({
  transactions,
  startDate,
  endDate,
  title = 'Label Totals by Month',
}) => {
  const { theme } = useTheme();

  // Validate date range spans multiple months
  const dateRangeError = React.useMemo(() => {
    if (!startDate || !endDate) {
      return 'Date range is required to display monthly trends';
    }

    const start = new Date(startDate);
    const end = new Date(endDate);

    // Check if dates are valid
    if (isNaN(start.getTime()) || isNaN(end.getTime())) {
      return 'Invalid date range provided';
    }

    // Check if dates span multiple months
    const startMonthYear = `${start.getFullYear()}-${start.getMonth()}`;
    const endMonthYear = `${end.getFullYear()}-${end.getMonth()}`;

    if (startMonthYear === endMonthYear) {
      return 'Date range must span multiple months to display monthly trends';
    }

    return null;
  }, [startDate, endDate]);

  // Process data to get monthly totals by label
  const monthlyData = React.useMemo(() => {
    if (dateRangeError || transactions.length === 0) {
      return [];
    }

    const start = new Date(startDate!);
    const end = new Date(endDate!);

    // Filter transactions within date range
    const transactionsInRange = transactions.filter((transaction) => {
      const transactionDate = new Date(transaction.date);
      return transactionDate >= start && transactionDate <= end;
    });

    // Group by month and label
    const monthlyLabelTotals: { [key: string]: number } = {};

    transactionsInRange.forEach((transaction) => {
      const transactionDate = new Date(transaction.date);
      const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
      const monthKey = `${monthNames[transactionDate.getMonth()]}-${transactionDate.getFullYear()}`;

      transaction.labels.forEach((label) => {
        const labelKey = `${label.key}=${label.value}`;
        const key = `${monthKey}|${labelKey}`;
        monthlyLabelTotals[key] = (monthlyLabelTotals[key] || 0) + transaction.amount;
      });
    });

    // Convert to array format
    const data: MonthlyLabelData[] = [];
    Object.entries(monthlyLabelTotals).forEach(([key, amount]) => {
      const [month, label] = key.split('|');
      data.push({ month, label, amount });
    });

    // Sort by actual date for proper chronological order
    return data.sort((a, b) => {
      const parseMonth = (monthStr: string) => {
        const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
        const [monthName, year] = monthStr.split('-');
        const monthIndex = monthNames.indexOf(monthName);
        return new Date(parseInt(year), monthIndex).getTime();
      };
      return parseMonth(a.month) - parseMonth(b.month);
    });
  }, [transactions, startDate, endDate, dateRangeError]);

  // Get unique labels and months for chart organization
  const uniqueLabels = React.useMemo(() => {
    const labels = new Set(monthlyData.map((d) => d.label));
    return Array.from(labels).sort();
  }, [monthlyData]);

  const uniqueMonths = React.useMemo(() => {
    const months = new Set(monthlyData.map((d) => d.month));
    const parseMonth = (monthStr: string) => {
      const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
      const [monthName, year] = monthStr.split('-');
      const monthIndex = monthNames.indexOf(monthName);
      return new Date(parseInt(year), monthIndex).getTime();
    };
    return Array.from(months).sort((a, b) => parseMonth(a) - parseMonth(b));
  }, [monthlyData]);

  // Prepare chart data grouped by label
  const chartData = React.useMemo(() => {
    return uniqueLabels.map((label) => {
      const labelData = uniqueMonths.map((month) => {
        const dataPoint = monthlyData.find((d) => d.label === label && d.month === month);
        return {
          x: month,
          y: dataPoint ? dataPoint.amount : 0,
          name: label,
        };
      });
      return { label, data: labelData };
    });
  }, [uniqueLabels, uniqueMonths, monthlyData]);

  // Calculate max value for Y axis
  const maxYValue = React.useMemo(() => {
    const allValues = monthlyData.map((d) => Math.abs(d.amount));
    if (allValues.length === 0) return 1000;
    const max = Math.max(...allValues);
    return Math.ceil(max * 1.1);
  }, [monthlyData]);

  // Create legend data
  const legendData = React.useMemo(() => {
    return uniqueLabels.map((label) => ({ name: label }));
  }, [uniqueLabels]);

  // Show error message if date range is invalid
  if (dateRangeError) {
    return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
            <div style={{ textAlign: 'center', color: 'var(--pf-v6-global--danger-color--100)' }}>
              <strong>Error:</strong> {dateRangeError}
            </div>
          </div>
        </CardBody>
      </Card>
    );
  }

  // Show no data message
  if (monthlyData.length === 0) {
    return (
      <Card>
        <CardBody>
          <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
            {title}
          </Title>
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
            <div style={{ textAlign: 'center', color: 'var(--pf-v6-global--Color--200)' }}>
              No transaction data available for the selected date range
            </div>
          </div>
        </CardBody>
      </Card>
    );
  }

  return (
    <Card style={{ height: '100%', width: '100%' }}>
      <CardBody>
        <Title headingLevel="h3" size="md" style={{ marginBottom: '1rem' }}>
          {title}
        </Title>

        <div style={{ height: '400px', width: '100%' }}>
          <Chart
            ariaDesc="Monthly label totals grouped bar chart"
            ariaTitle={title}
            containerComponent={
              <ChartVoronoiContainer
                labels={({ datum }) =>
                  `${datum.name}: ${datum.y.toLocaleString('fr-FR', {
                    style: 'currency',
                    currency: 'EUR',
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                  })}`
                }
                constrainToVisibleArea
              />
            }
            height={380}
            width={700}
            padding={{
              bottom: 60,
              left: 80,
              right: 150,
              top: 20,
            }}
            themeColor={ChartThemeColor.multiUnordered}
            domainPadding={{ x: 50 }}
            domain={{ y: [0, maxYValue] }}
          >
            <ChartAxis
              style={{
                tickLabels: {
                  fill: theme === 'dark' ? '#ffffff' : undefined,
                  fontSize: 14,
                },
              }}
            />
            <ChartAxis
              dependentAxis
              tickFormat={(value) => `${(value / 1000).toFixed(0)}k€`}
              style={{
                tickLabels: {
                  fill: theme === 'dark' ? '#ffffff' : undefined,
                  fontSize: 14,
                },
              }}
            />
            <ChartGroup offset={10}>
              {chartData.map((series) => (
                <ChartBar key={series.label} data={series.data} name={series.label} />
              ))}
            </ChartGroup>
            <ChartLegend
              data={legendData}
              orientation="vertical"
              x={570}
              y={50}
              style={{
                labels: {
                  fill: theme === 'dark' ? '#ffffff' : undefined,
                  fontSize: 14,
                },
              }}
            />
          </Chart>
        </div>
      </CardBody>
    </Card>
  );
};

export { LabelByMonthChart }; 