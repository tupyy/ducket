import * as React from 'react';
import { EuiPanel, EuiTitle, htmlIdGenerator } from '@elastic/eui';
import { PieChart, Pie, Cell, ResponsiveContainer, Sector } from 'recharts';
import { useTheme } from '@app/shared/contexts/ThemeContext';
import { IHierarchicalSummary } from './utils/calculateSummary';

interface TransactionSummaryChartProps {
  transactionSummary: IHierarchicalSummary;
}

const renderActiveShape = (props: any) => {
  const RADIAN = Math.PI / 180;
  const { cx, cy, midAngle, innerRadius, outerRadius, startAngle, endAngle, fill, payload, percent, value } = props;
  const sin = Math.sin(-RADIAN * midAngle);
  const cos = Math.cos(-RADIAN * midAngle);
  const sx = cx + (outerRadius + 10) * cos;
  const sy = cy + (outerRadius + 10) * sin;
  const mx = cx + (outerRadius + 30) * cos;
  const my = cy + (outerRadius + 30) * sin;
  const ex = mx + (cos >= 0 ? 1 : -1) * 22;
  const ey = my;
  const textAnchor = cos >= 0 ? 'start' : 'end';

  return (
    <g>
      <text x={cx} y={cy} dy={8} textAnchor="middle" fill={fill}>
        {payload.label}
      </text>
      <Sector
        cx={cx}
        cy={cy}
        innerRadius={innerRadius}
        outerRadius={outerRadius}
        startAngle={startAngle}
        endAngle={endAngle}
        fill={fill}
      />
      <Sector
        cx={cx}
        cy={cy}
        startAngle={startAngle}
        endAngle={endAngle}
        innerRadius={outerRadius + 6}
        outerRadius={outerRadius + 10}
        fill={fill}
      />
      <path d={`M${sx},${sy}L${mx},${my}L${ex},${ey}`} stroke={fill} fill="none" />
      <circle cx={ex} cy={ey} r={2} fill={fill} stroke="none" />
      <text x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} textAnchor={textAnchor} fill="#333">{`€${value.toFixed(2)}`}</text>
      <text x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} dy={18} textAnchor={textAnchor} fill="#999">
        {`(${(percent * 100).toFixed(2)}%)`}
      </text>
    </g>
  );
};

export const TransactionSummaryChart: React.FC<TransactionSummaryChartProps> = ({ transactionSummary }) => {
  const { theme } = useTheme();
  const htmlId = htmlIdGenerator();
  const chartId = htmlId();
  const [activeIndex, setActiveIndex] = React.useState(0);

  // Prepare data for pie chart (excluding total row)
  let pieChartData = transactionSummary.data.map((row) => {
    // For expanded view, extract just the value part from "key=value"
    const displayLabel = row.label.includes('=') ? row.label.split('=')[1] : row.label;
    
    return {
      label: displayLabel,
      value: Math.abs(row.debitAmount) + Math.abs(row.creditAmount),
      count: row.count,
    };
  });

  // Filter out zero values
  pieChartData = pieChartData.filter((item) => item.value > 0);

  // Check if we're showing filtered/expanded data and get the key
  const isFilteredView = transactionSummary.data.length > 0 && transactionSummary.data.some(item => item.label.includes('='));
  const expandedKey = isFilteredView && transactionSummary.data.length > 0 ? 
    transactionSummary.data[0].label.split('=')[0] : null;

  // Debug logging
  console.log('TransactionSummary data:', transactionSummary.data);
  console.log('Pie chart data:', pieChartData);
  console.log('Expanded key:', expandedKey);

  const colors = [
    '#006BB4',
    '#017D73',
    '#BD271E',
    '#DD0A73',
    '#9170B8',
    '#CA8EAE',
    '#D36086',
    '#E7664C',
    '#F5A700',
    '#54B399',
    '#79AAD9',
    '#B9A888',
  ];

  const onPieEnter = (_: any, index: number) => {
    setActiveIndex(index);
  };

  return (
    <EuiPanel paddingSize="m" style={{ marginBottom: '1rem' }}>
      <div className="eui-textCenter" style={{ marginBottom: '0.5rem' }}>
        <EuiTitle size="xs">
          <h3 id={chartId}>Summary Distribution</h3>
        </EuiTitle>
        {expandedKey && (
          <div style={{ fontSize: '12px', color: '#69707D', marginTop: '4px' }}>
            {expandedKey}
          </div>
        )}
      </div>

      <div style={{ height: '300px', width: '100%' }}>
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              activeShape={renderActiveShape}
              data={pieChartData}
              cx="50%"
              cy="50%"
              innerRadius={60}
              outerRadius={80}
              fill="#8884d8"
              dataKey="value"
              onMouseEnter={onPieEnter}
            >
              {pieChartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
              ))}
            </Pie>
          </PieChart>
        </ResponsiveContainer>
      </div>
    </EuiPanel>
  );
};
