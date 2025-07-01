import * as React from 'react';
import {
  Title,
  Card,
  CardBody,
} from '@patternfly/react-core';
import { Table, Thead, Tbody, Tr, Th, Td } from '@patternfly/react-table';
import { ITransaction } from '@app/shared/models/transaction';

interface TagTotal {
  tag: string;
  totalAmount: number;
  transactionCount: number;
}

interface TagTotalTableProps {
  transactions: ITransaction[];
  startDate?: string;
  endDate?: string;
}

const TagTotalTable: React.FunctionComponent<TagTotalTableProps> = ({ transactions, startDate, endDate }) => {
  // Compute totals by tag
  const tagTotals = React.useMemo(() => {
    const tagMap = new Map<string, { totalAmount: number; transactionCount: number }>();

    transactions.forEach((transaction) => {
      transaction.tags.forEach((tag) => {
        const existing = tagMap.get(tag.value) || { totalAmount: 0, transactionCount: 0 };
        tagMap.set(tag.value, {
          totalAmount: existing.totalAmount + transaction.amount,
          transactionCount: existing.transactionCount + 1,
        });
      });
    });

    // Convert to array and sort by total amount (descending)
    return Array.from(tagMap.entries())
      .map(([tag, data]) => ({
        tag,
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
    <Card style={{ height: '400px' }}>
      <CardBody>
        <Title headingLevel="h3" size="lg" style={{ marginBottom: '1rem' }}>
          Total Amounts by Tag
        </Title>
        
        <div style={{ height: '320px', overflowY: 'auto' }}>
          <Table aria-label="Tag totals table" variant="compact">
            <Thead>
              <Tr>
                <Th width={40}>Tag</Th>
                <Th width={30}>Total Amount</Th>
                <Th width={30}>Transaction Count</Th>
              </Tr>
            </Thead>
            <Tbody>
              {tagTotals.map((tagTotal) => (
                <Tr key={tagTotal.tag}>
                  <Td>{tagTotal.tag}</Td>
                  <Td>
                    <span
                      style={{
                        color: tagTotal.totalAmount >= 0 ? '#3e8635' : '#c9190b',
                        fontWeight: 'bold',
                      }}
                    >
                      {tagTotal.totalAmount.toLocaleString('fr-FR', {
                        style: 'currency',
                        currency: 'EUR',
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                      })}
                    </span>
                  </Td>
                  <Td>{tagTotal.transactionCount}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </div>
      </CardBody>
    </Card>
  );
};

export { TagTotalTable }; 