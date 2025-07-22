import * as React from 'react';
import { Card, CardBody, Content } from '@patternfly/react-core';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { ITransactionSummary } from './reducers/transactionSummary.reducer';

interface TransactionSummaryProps {
  transactionSummary: ITransactionSummary;
}

export const TransactionSummary: React.FC<TransactionSummaryProps> = ({ transactionSummary }) => {
  return (
    <Card style={{ marginBottom: '1rem' }}>
      <CardBody>
        <Content>
          <strong>{transactionSummary.type === 'filtered' ? 'Filtered Transaction Summary' : 'Summary'}</strong>
        </Content>
        <Table
          aria-label="transaction-summary"
          variant="compact"
          borders={false}
          style={{ marginTop: '0.5rem' }}
        >
          <Thead>
            <Tr>
              <Th>{transactionSummary.type === 'filtered' ? 'Type' : 'Label'}</Th>
              <Th>Count</Th>
              <Th>Debit</Th>
              <Th>Credit</Th>
            </Tr>
          </Thead>
          <Tbody>
            {transactionSummary.data.map((row, index) => (
              <Tr key={index}>
                <Td>
                  <Content>{row.label}</Content>
                </Td>
                <Td>
                  <Content>{row.count}</Content>
                </Td>
                <Td>
                  <Content
                    style={{
                      color:
                        row.debitAmount > 0
                          ? 'var(--pf-v6-global--danger-color--100)'
                          : 'var(--pf-v6-global--palette--black-600)',
                    }}
                  >
                    {row.debitAmount > 0
                      ? row.debitAmount.toLocaleString('de-DE', {
                          minimumFractionDigits: 2,
                          maximumFractionDigits: 2,
                        })
                      : '-'}
                  </Content>
                </Td>
                <Td>
                  <Content
                    style={{
                      color:
                        row.creditAmount > 0
                          ? 'var(--pf-v6-global--success-color--100)'
                          : 'var(--pf-v6-global--palette--black-600)',
                    }}
                  >
                    {row.creditAmount > 0
                      ? row.creditAmount.toLocaleString('de-DE', {
                          minimumFractionDigits: 2,
                          maximumFractionDigits: 2,
                        })
                      : '-'}
                  </Content>
                </Td>
              </Tr>
            ))}
            {/* Total row */}
            <Tr
              style={{
                borderTopWidth: '1px',
                borderTopStyle: 'solid',
                borderTopColor: 'var(--pf-v6-global--BorderColor--100)',
                backgroundColor: 'var(--pf-v6-global--palette--blue-50)',
              }}
            >
              <Td>
                <Content>
                  <strong>Total</strong>
                </Content>
              </Td>
              <Td>
                <Content style={{ fontWeight: 'bold' }}>{transactionSummary.totals.count}</Content>
              </Td>
              <Td>
                <Content
                  style={{
                    color: 'var(--pf-v6-global--danger-color--100)',
                    fontWeight: 'bold',
                  }}
                >
                  {transactionSummary.totals.debitAmount > 0
                    ? transactionSummary.totals.debitAmount.toLocaleString('de-DE', {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                      })
                    : '-'}
                </Content>
              </Td>
              <Td>
                <Content
                  style={{
                    fontWeight: 'bold',
                  }}
                >
                  {transactionSummary.totals.creditAmount > 0
                    ? transactionSummary.totals.creditAmount.toLocaleString('de-DE', {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                      })
                    : '-'}
                </Content>
              </Td>
            </Tr>
          </Tbody>
        </Table>
      </CardBody>
    </Card>
  );
};
