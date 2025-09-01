import * as React from 'react';
import {
  EuiPanel,
  EuiTitle,
  EuiBasicTable,
  EuiBasicTableColumn,
  EuiButtonEmpty,
  EuiText,
} from '@elastic/eui';
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

  // Prepare data for EuiBasicTable
  const tableData = [
    ...transactionSummary.data.map((row, index) => ({
      id: index,
      label: row.label,
      count: row.count,
      debitAmount: row.debitAmount,
      creditAmount: row.creditAmount,
      isTotal: false,
    })),
    {
      id: 'total',
      label: 'Total',
      count: transactionSummary.totals.count,
      debitAmount: transactionSummary.totals.debitAmount,
      creditAmount: transactionSummary.totals.creditAmount,
      isTotal: true,
    },
  ];

  const columns: Array<EuiBasicTableColumn<typeof tableData[0]>> = [
    {
      field: 'label',
      name: transactionSummary.type === 'filtered' ? 'Type' : 'Label',
      render: (label: string, item) => {
        if (item.isTotal) {
          return <EuiText size="s" style={{ fontWeight: 'bold' }}>{label}</EuiText>;
        }
        
        const wildcardFilter = `${label}=*`;
        const isFilterActive = selectedLabels.includes(wildcardFilter);
        
        return (
          <EuiButtonEmpty
            size="xs"
            onClick={() => handleLabelClick(label)}
            style={{
              fontWeight: isFilterActive ? 'bold' : 'normal',
              color: isFilterActive ? '#006BB4' : '#007B94',
            }}
            aria-label={`Filter by ${label} labels`}
          >
            {label} {isFilterActive && '(filtered)'}
          </EuiButtonEmpty>
        );
      },
    },
    {
      field: 'count',
      name: 'Count',
      render: (count: number, item) => (
        <EuiText size="s" style={{ fontWeight: item.isTotal ? 'bold' : 'normal' }}>
          {count}
        </EuiText>
      ),
      width: '80px',
    },
    {
      field: 'debitAmount',
      name: 'Debit',
      render: (amount: number, item) => (
        <EuiText 
          size="s" 
          style={{ 
            fontWeight: item.isTotal ? 'bold' : 'normal',
            color: amount > 0 ? '#BD271E' : '#69707D',
          }}
        >
          {amount > 0
            ? amount.toLocaleString('de-DE', {
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
              })
            : '-'}
        </EuiText>
      ),
      width: '120px',
    },
    {
      field: 'creditAmount',
      name: 'Credit',
      render: (amount: number, item) => (
        <EuiText 
          size="s" 
          style={{ 
            fontWeight: item.isTotal ? 'bold' : 'normal',
            color: amount > 0 ? '#017D73' : '#69707D',
          }}
        >
          {amount > 0
            ? amount.toLocaleString('de-DE', {
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
              })
            : '-'}
        </EuiText>
      ),
      width: '120px',
    },
  ];


  return (
    <EuiPanel paddingSize="m" style={{ marginBottom: '1rem' }}>
      <div style={{ marginBottom: '0.5rem' }}>
        <EuiTitle size="xs">
          <h3>{transactionSummary.type === 'filtered' ? 'Filtered Transaction Summary' : 'Summary'}</h3>
        </EuiTitle>
      </div>
      
      <EuiBasicTable
        items={tableData}
        columns={columns}
        rowProps={(item) => ({
          style: item.isTotal ? {
            borderTop: '1px solid #D3DAE6',
            backgroundColor: '#F7F9FC',
          } : undefined,
        })}
      />
    </EuiPanel>
  );
};
