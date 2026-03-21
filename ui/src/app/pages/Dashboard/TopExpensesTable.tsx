import * as React from 'react';
import { Card, CardTitle, CardBody, Label, Content } from '@patternfly/react-core';
import { Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import { ITransaction } from '@app/shared/models/transaction';

interface TopExpensesTableProps {
  transactions: ITransaction[];
}

const TopExpensesTable: React.FunctionComponent<TopExpensesTableProps> = ({ transactions }) => {
  const formatDate = (dateStr: string) => {
    try {
      const d = new Date(dateStr);
      const day = String(d.getDate()).padStart(2, '0');
      const month = String(d.getMonth() + 1).padStart(2, '0');
      const year = d.getFullYear();
      return `${day}-${month}-${year}`;
    } catch {
      return dateStr;
    }
  };

  if (transactions.length === 0) {
    return (
      <Card>
        <CardTitle>Top 10 Expenses</CardTitle>
        <CardBody>
          <Content component="p">No debit transactions found.</Content>
        </CardBody>
      </Card>
    );
  }

  return (
    <Card>
      <CardTitle>Top 10 Expenses</CardTitle>
      <CardBody>
        <Table aria-label="Top expenses" variant="compact">
          <Thead>
            <Tr>
              <Th>Date</Th>
              <Th>Content</Th>
              <Th>Tags</Th>
              <Th>Amount</Th>
            </Tr>
          </Thead>
          <Tbody>
            {transactions.map((txn) => (
              <Tr key={txn.id}>
                <Td>{formatDate(txn.date)}</Td>
                <Td>
                  <Content component="p" style={{ maxWidth: '350px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {txn.content}
                  </Content>
                </Td>
                <Td>
                  {txn.tags.map((tag, i) => (
                    <Label key={i} color="teal" isCompact style={{ marginRight: 4 }}>
                      {tag}
                    </Label>
                  ))}
                </Td>
                <Td style={{ fontFamily: 'monospace', fontWeight: 'bold' }}>
                  €{txn.amount.toFixed(2)}
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </CardBody>
    </Card>
  );
};

export { TopExpensesTable };
