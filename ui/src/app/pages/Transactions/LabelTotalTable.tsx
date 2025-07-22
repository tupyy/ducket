import * as React from 'react';
import {
  Title,
  Card,
  CardBody,
} from '@patternfly/react-core';
import { Table, Thead, Tbody, Tr, Th, Td } from '@patternfly/react-table';
import { ITransaction } from '@app/shared/models/transaction';

interface LabelTotalTableProps {
  transactions: ITransaction[];
  startDate?: string;
  endDate?: string;
}

const LabelTotalTable: React.FunctionComponent<LabelTotalTableProps> = ({ transactions}) => {
  // Compute totals by label
  const labelTotals = React.useMemo(() => {
    const labelMap = new Map<string, { totalAmount: number; transactionCount: number }>();

    transactions.forEach((transaction) => {
      transaction.labels.forEach((label) => {
        const labelKey = `${label.key}=${label.value}`;
        const existing = labelMap.get(labelKey) || { totalAmount: 0, transactionCount: 0 };
        labelMap.set(labelKey, {
          totalAmount: existing.totalAmount + transaction.amount,
          transactionCount: existing.transactionCount + 1,
        });
      });
    });

    // Convert to array and sort by total amount (descending)
    return Array.from(labelMap.entries())
      .map(([label, data]) => ({
        label,
        totalAmount: data.totalAmount,
        transactionCount: data.transactionCount,
      }))
      .sort((a, b) => Math.abs(b.totalAmount) - Math.abs(a.totalAmount));
  }, [transactions]);

  // Don't render if no transactions
  if (transactions.length === 0) {
    return null;
  }

  return (
    <Card style={{ height: '100%' }}>
      <CardBody>
        <Title headingLevel="h3" size="lg" style={{ marginBottom: '1rem' }}>
          Total Amounts by Label
        </Title>

        <div style={{ height: '320px', overflowY: 'auto' }}>
          <Table aria-label="Label totals table" variant="compact">
            <Thead>
              <Tr>
                <Th width={40}>Label</Th>
                <Th width={30}>Total Amount</Th>
                <Th width={30}>Transaction Count</Th>
              </Tr>
            </Thead>
            <Tbody>
              {labelTotals.map((labelTotal) => (
                <Tr key={labelTotal.label}>
                  <Td>{labelTotal.label}</Td>
                  <Td>
                    <span
                      style={{
                        color: labelTotal.totalAmount >= 0 ? '#3e8635' : '#c9190b',
                        fontWeight: 'bold',
                      }}
                    >
                      {labelTotal.totalAmount.toLocaleString('fr-FR', {
                        style: 'currency',
                        currency: 'EUR',
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                      })}
                    </span>
                  </Td>
                  <Td>{labelTotal.transactionCount}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </div>
      </CardBody>
    </Card>
  );
};

export { LabelTotalTable };
