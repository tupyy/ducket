import * as React from 'react';
import {
  EuiBasicTable,
  EuiBasicTableColumn,
  EuiButtonEmpty,
  EuiPanel,
  EuiText,
  EuiTitle,
} from '@elastic/eui';
import { IHierarchicalSummary, IHierarchicalSummaryData } from './utils/calculateSummary';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { setSelectedLabels, setShowOnlyUnlabeled } from './reducers/transaction-filter.reducer';
import { useThemeStyles } from '@app/shared/hooks/useThemeStyles';

interface TransactionSummaryProps {
  transactionSummary: IHierarchicalSummary;
  expandedKeys: Set<string>;
  setExpandedKeys: React.Dispatch<React.SetStateAction<Set<string>>>;
}

export const TransactionSummary: React.FC<TransactionSummaryProps> = ({ 
  transactionSummary, 
  expandedKeys, 
  setExpandedKeys 
}) => {
  const dispatch = useAppDispatch();
  const { selectedLabels, showOnlyUnlabeled } = useAppSelector((state) => state.transactionFilter);
  const themeStyles = useThemeStyles();

  /**
   * Handle clicking on a parent label key to expand/collapse
   * @param labelKey - The label key to expand/collapse
   */
  const handleLabelClick = (labelKey: string) => {
    const newExpandedKeys = new Set(expandedKeys);
    if (expandedKeys.has(labelKey)) {
      newExpandedKeys.delete(labelKey);
    } else {
      newExpandedKeys.add(labelKey);
    }
    setExpandedKeys(newExpandedKeys);
  };

  // Prepare hierarchical data for EuiBasicTable
  const tableData: Array<{
    id: string | number;
    label: string;
    count: number;
    debitAmount: number;
    creditAmount: number;
    isTotal: boolean;
    isParent: boolean;
    isChild: boolean;
    parentKey?: string;
  }> = [];

  const hasAnyExpanded = expandedKeys.size > 0;

  if (hasAnyExpanded) {
    // Filtered view: only show expanded labels and their children
    transactionSummary.data.forEach((parentRow, parentIndex) => {
      if (expandedKeys.has(parentRow.label)) {
        // Add parent row but with zero totals (for collapsing back)
        tableData.push({
          id: `parent-${parentIndex}`,
          label: parentRow.label,
          count: 0,
          debitAmount: 0,
          creditAmount: 0,
          isTotal: false,
          isParent: true,
          isChild: false,
        });

        // Add children
        if (parentRow.children) {
          parentRow.children.forEach((childRow, childIndex) => {
            tableData.push({
              id: `child-${parentIndex}-${childIndex}`,
              label: childRow.label,
              count: childRow.count,
              debitAmount: childRow.debitAmount,
              creditAmount: childRow.creditAmount,
              isTotal: false,
              isParent: false,
              isChild: true,
              parentKey: parentRow.label,
            });
          });
        }
      }
    });
  } else {
    // Default view: show all parent rows collapsed
    transactionSummary.data.forEach((parentRow, parentIndex) => {
      tableData.push({
        id: `parent-${parentIndex}`,
        label: parentRow.label,
        count: parentRow.count,
        debitAmount: parentRow.debitAmount,
        creditAmount: parentRow.creditAmount,
        isTotal: false,
        isParent: true,
        isChild: false,
      });
    });
  }

  // Calculate totals for the filtered view
  const filteredTotals = hasAnyExpanded ? 
    tableData
      .filter(item => !item.isTotal && item.isChild) // Only count child rows
      .reduce(
        (acc, row) => ({
          count: acc.count + row.count,
          debitAmount: acc.debitAmount + row.debitAmount,
          creditAmount: acc.creditAmount + row.creditAmount,
        }),
        { count: 0, debitAmount: 0, creditAmount: 0 }
      ) : transactionSummary.totals;

  // Add total row
  tableData.push({
    id: 'total',
    label: 'Total',
    count: filteredTotals.count,
    debitAmount: filteredTotals.debitAmount,
    creditAmount: filteredTotals.creditAmount,
    isTotal: true,
    isParent: false,
    isChild: false,
  });

  const columns: Array<EuiBasicTableColumn<typeof tableData[0]>> = [
    {
      field: 'label',
      name: 'Label',
      render: (label: string, item) => {
        if (item.isTotal) {
          return <EuiText size="s" style={{ fontWeight: 'bold', color: themeStyles.textColor }}>{label}</EuiText>;
        }

        if (item.isChild) {
          return (
            <EuiText size="s" style={{ paddingLeft: '1.5rem', color: themeStyles.textColor }}>
              {label}
            </EuiText>
          );
        }

        // In expanded view, show regular rows (which are the key=value pairs)
        if (hasAnyExpanded && !item.isParent) {
          return <EuiText size="s" style={{ color: themeStyles.textColor }}>{label}</EuiText>;
        }

        if (item.isParent) {
          const isExpanded = expandedKeys.has(label);
          return (
            <EuiButtonEmpty
              size="xs"
              onClick={() => handleLabelClick(label)}
              style={{
                fontWeight: 'normal',
              }}
              aria-label={`${isExpanded ? 'Collapse' : 'Expand'} ${label} details`}
            >
              {isExpanded ? '▼' : '▶'} {label}
            </EuiButtonEmpty>
          );
        }

        return <EuiText size="s">{label}</EuiText>;
      },
    },
    {
      field: 'count',
      name: 'Count',
      render: (count: number, item) => (
        <EuiText size="s" style={{ 
          fontWeight: item.isTotal ? 'bold' : 'normal',
          color: (item.isTotal || item.isChild || (hasAnyExpanded && !item.isParent)) ? themeStyles.textColor : undefined,
        }}>
          {item.isParent && hasAnyExpanded ? '-' : count}
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
            color: (item.isTotal || item.isChild || (hasAnyExpanded && !item.isParent)) ? themeStyles.textColor : (amount > 0 ? '#BD271E' : '#69707D'),
          }}
        >
          {item.isParent && hasAnyExpanded ? '-' : 
           amount > 0
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
            color: (item.isTotal || item.isChild || (hasAnyExpanded && !item.isParent)) ? themeStyles.textColor : (amount > 0 ? '#017D73' : '#69707D'),
          }}
        >
          {item.isParent && hasAnyExpanded ? '-' : 
           amount > 0
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
    <EuiPanel 
      paddingSize="m" 
      style={{ 
        marginBottom: '1rem',
      }}
    >
      <div style={{ marginBottom: '0.5rem' }}>
        <EuiTitle size="xs">
          <h3>Transaction Summary by Label Key</h3>
        </EuiTitle>
      </div>

      <EuiBasicTable
        items={tableData}
        columns={columns}
        rowProps={(item) => ({
          style: item.isTotal ? {
            borderTop: '2px solid #D3DAE6',
            backgroundColor: themeStyles.cardBackground,
            marginTop: '8px',
          } : item.isChild ? {
            backgroundColor: themeStyles.cardBackground,
          } : (hasAnyExpanded && !item.isParent) ? {
            backgroundColor: themeStyles.cardBackground,
          } : undefined,
        })}
      />
    </EuiPanel>
  );
};
