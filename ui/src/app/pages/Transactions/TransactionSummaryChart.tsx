import * as React from 'react';
import {
  EuiPanel,
  EuiTitle,
  htmlIdGenerator,
} from '@elastic/eui';
import { Chart, Partition, Settings, PartitionLayout, PartialTheme } from '@elastic/charts';
import { useTheme } from '@app/shared/contexts/ThemeContext';
import { ITransactionSummary } from './reducers/transactionSummary.reducer';

interface TransactionSummaryChartProps {
  transactionSummary: ITransactionSummary;
}

export const TransactionSummaryChart: React.FC<TransactionSummaryChartProps> = ({ transactionSummary }) => {
  const { theme } = useTheme();
  const htmlId = htmlIdGenerator();
  const chartId = htmlId();

  // Prepare data for pie chart (excluding total row)
  let pieChartData = transactionSummary.data.map((row) => ({
    label: row.label,
    value: Math.abs(row.debitAmount) + Math.abs(row.creditAmount),
    count: row.count,
  }));

  // Filter out zero values
  pieChartData = pieChartData.filter(item => item.value > 0);

  // Debug logging
  console.log('TransactionSummary data:', transactionSummary.data);
  console.log('Pie chart data:', pieChartData);

  const themeOverrides: PartialTheme = {
    partition: { emptySizeRatio: 0.4 },
  };

  return (
    <EuiPanel paddingSize="m" style={{ marginBottom: '1rem' }}>
      <EuiTitle className="eui-textCenter" size="xs">
        <h3 id={chartId}>Summary Distribution</h3>
      </EuiTitle>
      
      <div style={{ height: '300px' }}>
        {pieChartData.length === 0 ? (
          <div style={{ textAlign: 'center', paddingTop: '100px' }}>
            No data to display
          </div>
        ) : (
          <Chart size={{ height: 300 }}>
            <Settings
              theme={themeOverrides}
              ariaLabelledBy={chartId}
            />
            <Partition
              id="transactionSummaryPie"
              data={pieChartData}
              layout={PartitionLayout.sunburst}
              valueAccessor={(d) => d.value}
              valueFormatter={(value) => `€${value.toFixed(2)}`}
              layers={[
                {
                  groupByRollup: (d) => d.label,
                  shape: {
                    fillColor: (_, sortIndex) => {
                      const colors = [
                        '#006BB4', '#017D73', '#BD271E', '#DD0A73', 
                        '#9170B8', '#CA8EAE', '#D36086', '#E7664C',
                        '#F5A700', '#54B399', '#79AAD9', '#B9A888'
                      ];
                      return colors[sortIndex % colors.length];
                    },
                  },
                },
              ]}
              clockwiseSectors={false}
            />
          </Chart>
        )}
      </div>
    </EuiPanel>
  );
};