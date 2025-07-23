import * as React from 'react';
import { Card, CardBody, Content, Button } from '@patternfly/react-core';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { ITransactionSummary } from './reducers/transactionSummary.reducer';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { setSelectedLabels, setShowOnlyUnlabeled } from './reducers/transaction-filter.reducer';

interface TransactionSummaryProps {
  transactionSummary: ITransactionSummary;
}

export const TransactionSummary: React.FC<TransactionSummaryProps> = ({ transactionSummary }) => {
  const dispatch = useAppDispatch();
  const { selectedLabels, showOnlyUnlabeled } = useAppSelector((state) => state.transactionFilter);

  /**
   * Handle clicking on a label to add wildcard filter
   * @param labelKey - The label key to filter by (e.g., "category")
   */
  const handleLabelClick = (labelKey: string) => {
    const wildcardFilter = `${labelKey}=*`;
    
    // If showing only unlabeled transactions, clear that filter first
    if (showOnlyUnlabeled) {
      dispatch(setShowOnlyUnlabeled(false));
    }
    
    // Check if the wildcard filter already exists
    if (!selectedLabels.includes(wildcardFilter)) {
      // Add the wildcard filter to existing selected labels
      const newSelectedLabels = [...selectedLabels, wildcardFilter];
      dispatch(setSelectedLabels(newSelectedLabels));
    }
  };

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
            {transactionSummary.data.map((row, index) => {
              const wildcardFilter = `${row.label}=*`;
              const isFilterActive = selectedLabels.includes(wildcardFilter);
              
              return (
                <Tr key={index}>
                  <Td>
                    <Button
                      variant="link"
                      onClick={() => handleLabelClick(row.label)}
                      style={{
                        padding: 0,
                        fontSize: 'inherit',
                        fontWeight: isFilterActive ? 'bold' : 'inherit',
                        textDecoration: 'none',
                        color: isFilterActive 
                          ? 'var(--pf-v6-global--primary-color--100)' 
                          : 'var(--pf-v6-global--link--Color)',
                        cursor: 'pointer',
                      }}
                      aria-label={`Filter by ${row.label} labels`}
                    >
                      {row.label} {isFilterActive && '(filtered)'}
                    </Button>
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
              );
            })}
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
