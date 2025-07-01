import * as React from 'react';
import {
  Title,
  Card,
  CardBody,
} from '@patternfly/react-core';
import { Table, Thead, Tbody, Tr, Th, Td } from '@patternfly/react-table';
import { useAppSelector, useAppDispatch } from '@app/shared/store';
import { setSourceTransactions } from '@app/shared/reducers/transaction-filter.reducer';

interface TagTotal {
  tag: string;
  totalAmount: number;
  transactionCount: number;
}

interface TagTotalTableProps {
  startDate?: string;
  endDate?: string;
}

const TagTotalTable: React.FunctionComponent<TagTotalTableProps> = ({ startDate, endDate }) => {
  const dispatch = useAppDispatch();
  const { filteredTransactions, sourceTransactions } = useAppSelector((state) => state.transactionFilter);
  const { transactions } = useAppSelector((state) => state.transactions);

  // Initialize filter reducer if it's empty but we have transactions
  React.useEffect(() => {
    if (sourceTransactions.length === 0 && transactions.length > 0) {
      dispatch(setSourceTransactions(transactions));
    }
  }, [sourceTransactions.length, transactions.length, transactions, dispatch]);

  // Compute totals by tag
  const tagTotals = React.useMemo(() => {
    const tagMap = new Map<string, { totalAmount: number; transactionCount: number }>();

    filteredTransactions.forEach((transaction) => {
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
  }, [filteredTransactions]);

  // Don't render if no transactions
  if (filteredTransactions.length === 0) {
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