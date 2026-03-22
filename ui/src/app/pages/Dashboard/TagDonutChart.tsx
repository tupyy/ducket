import * as React from 'react';
import { Card, CardTitle, CardBody, Content } from '@patternfly/react-core';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';
import { TagSummary } from '@app/shared/reducers/dashboard.reducer';

interface TagDonutChartProps {
  data: TagSummary[];
}

const COLORS = [
  '#06c', '#4cb140', '#f0ab00', '#c9190b', '#a18fff',
  '#009596', '#ec7a08', '#7d1007', '#3e8635', '#002f5d',
];

const TagDonutChart: React.FunctionComponent<TagDonutChartProps> = ({ data }) => {
  const MAX_SLICES = 8;
  const sorted = [...data].sort((a, b) => b.total_debit - a.total_debit);
  const top = sorted.slice(0, MAX_SLICES);
  const rest = sorted.slice(MAX_SLICES);

  const chartData = top.map((d) => ({ name: d.tag, value: d.total_debit }));
  if (rest.length > 0) {
    const otherTotal = rest.reduce((sum, d) => sum + d.total_debit, 0);
    chartData.push({ name: `Other (${rest.length})`, value: otherTotal });
  }

  if (chartData.length === 0) {
    return (
      <Card isFullHeight>
        <CardTitle>Spending by Tag</CardTitle>
        <CardBody>
          <Content component="p">No tagged transactions found.</Content>
        </CardBody>
      </Card>
    );
  }

  return (
    <Card isFullHeight>
      <CardTitle>Spending by Tag</CardTitle>
      <CardBody>
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              innerRadius={60}
              outerRadius={100}
              dataKey="value"
              nameKey="name"
              label={({ name, percent }) => `${name} ${((percent ?? 0) * 100).toFixed(0)}%`}
            >
              {chartData.map((_, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip formatter={(value) => `€${Number(value).toFixed(2)}`} />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      </CardBody>
    </Card>
  );
};

export { TagDonutChart };
