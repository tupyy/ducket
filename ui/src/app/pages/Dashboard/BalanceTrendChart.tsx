import * as React from 'react';
import { Card, CardTitle, CardBody, Content } from '@patternfly/react-core';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { BalanceTrendPoint } from '@app/shared/reducers/dashboard.reducer';

interface BalanceTrendChartProps {
  data: BalanceTrendPoint[];
}

const BalanceTrendChart: React.FunctionComponent<BalanceTrendChartProps> = ({ data }) => {
  if (data.length === 0) {
    return (
      <Card isFullHeight>
        <CardTitle>Balance Trend</CardTitle>
        <CardBody>
          <Content component="p">No data available.</Content>
        </CardBody>
      </Card>
    );
  }

  return (
    <Card isFullHeight>
      <CardTitle>Balance Trend</CardTitle>
      <CardBody>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="month" />
            <YAxis />
            <Tooltip formatter={(value: number) => `€${value.toFixed(2)}`} />
            <Legend />
            <Line type="monotone" dataKey="credit" stroke="#4cb140" name="Credits" strokeWidth={2} />
            <Line type="monotone" dataKey="debit" stroke="#c9190b" name="Debits" strokeWidth={2} />
            <Line type="monotone" dataKey="balance" stroke="#06c" name="Balance" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </CardBody>
    </Card>
  );
};

export { BalanceTrendChart };
