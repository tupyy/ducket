import * as React from 'react';
import {
  EuiInMemoryTable,
  EuiBasicTableColumn,
  EuiTableActionsColumnType,
  EuiButtonEmpty,
  EuiButtonIcon,
  EuiFlexGroup,
  EuiFlexItem,
  EuiBadge,
  EuiText,
  EuiSpacer,
  EuiPanel,
  EuiFormRow,
  EuiFieldText,
  EuiSwitch,
  EuiButton,
  EuiLoadingSpinner,
  EuiCallOut,
  EuiSelect,
  EuiSelectOption,
  EuiToolTip,
} from '@elastic/eui';
import { ILabelTransaction, ITransaction } from '@app/shared/models/transaction';
import { LabelFilter } from '@app/shared/components/label-filter';
import { AddLabelModal } from './AddLabelModal';
import { RemoveLabelModal } from './RemoveLabelModal';
import { CreateRuleModal } from './CreateRuleModal';
import { TransactionDetailsModal } from './TransactionDetailsModal';
import { useTheme } from '@app/shared/contexts/ThemeContext';
import { getAccountColor, getAccountDarkColor } from '@app/utils/colorUtils';
import { safeFormatDateString } from '@app/utils/dateUtils';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import {
  setSourceTransactions,
  setSelectedLabels,
  setSelectedTransactionTypes,
  setSelectedAccounts,
  setDescriptionFilter,
  setShowOnlyUnlabeled,
  clearAllFilters,
  applyFilters,
  setPage,
  setPerPage,
  setSorting,
} from './reducers/transaction-filter.reducer';
import {
  clearAddLabelToTransactionSuccess,
  removeLabelFromTransaction,
  updateTransactionInfo,
} from '@app/shared/reducers/transaction.reducer';

export interface ITransactionListProps {
  transactions: Array<ITransaction> | [];
}

const TransactionList: React.FunctionComponent<ITransactionListProps> = ({ transactions }) => {
  const { theme } = useTheme();
  const dispatch = useAppDispatch();

  // Get all table state from the filter reducer
  const {
    filteredTransactions,
    selectedLabels,
    selectedTransactionTypes,
    selectedAccounts,
    descriptionFilter,
    showOnlyUnlabeled,
    page,
    perPage,
    sortDirection,
    sortIndex,
    filtering,
  } = useAppSelector((state) => state.transactionFilter);

  // Get transaction update state
  const { updatingInfo, errorMessage } = useAppSelector((state) => state.transactions);

  // Modal states
  const [isTransactionDetailsModalOpen, setIsTransactionDetailsModalOpen] = React.useState(false);
  const [selectedTransactionForDetails, setSelectedTransactionForDetails] = React.useState<ITransaction | null>(null);
  const [isBulkLabelModalOpen, setIsBulkLabelModalOpen] = React.useState(false);
  const [isRemoveLabelModalOpen, setIsRemoveLabelModalOpen] = React.useState(false);
  const [labelToRemove, setLabelToRemove] = React.useState<{
    transaction: ITransaction;
    label: ILabelTransaction;
  } | null>(null);
  const [isCreateRuleModalOpen, setIsCreateRuleModalOpen] = React.useState(false);
  const [selectedTransactionForRule, setSelectedTransactionForRule] = React.useState<ITransaction | null>(null);

  // Selection state for bulk operations
  const [selectedItems, setSelectedItems] = React.useState<ITransaction[]>([]);

  // Calculate available values for filters
  const availableLabels = React.useMemo(() => {
    const labelSet = new Set<string>();
    const keySet = new Set<string>();

    transactions.forEach((transaction) => {
      transaction.labels.forEach((label) => {
        labelSet.add(`${label.key}=${label.value}`);
        keySet.add(label.key);
      });
    });

    keySet.forEach((key) => {
      labelSet.add(`${key}=*`);
    });

    const labels = Array.from(labelSet);
    return labels.sort((a, b) => {
      const [keyA, valueA] = a.split('=');
      const [keyB, valueB] = b.split('=');
      if (keyA !== keyB) return keyA.localeCompare(keyB);
      if (valueA === '*' && valueB !== '*') return -1;
      if (valueA !== '*' && valueB === '*') return 1;
      return valueA.localeCompare(valueB);
    });
  }, [transactions]);

  const availableAccounts = React.useMemo(() => {
    const accounts = Array.from(new Set(transactions.map(t => t.account)));
    return accounts.sort();
  }, [transactions]);

  const availableTransactionTypes = React.useMemo(() => {
    const types = Array.from(new Set(transactions.map(t => t.kind)));
    return types.sort();
  }, [transactions]);

  // Color mapping function for labels based on key
  const getLabelColor = (labelKey: string): string => {
    const colors = [
      'primary', 'success', 'warning', 'danger', 'accent', 
      'default', 'subdued', 'hollow'
    ];
    
    // Create a simple hash from the label key
    let hash = 0;
    for (let i = 0; i < labelKey.length; i++) {
      const char = labelKey.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash; // Convert to 32bit integer
    }
    
    // Use absolute value and modulo to get a consistent color index
    const colorIndex = Math.abs(hash) % colors.length;
    return colors[colorIndex];
  };

  // Initialize transactions in the filter store
  React.useEffect(() => {
    dispatch(setSourceTransactions(transactions));
  }, [transactions, dispatch]);

  // Apply filters when filter criteria change
  React.useEffect(() => {
    dispatch(applyFilters({
      selectedLabels,
      selectedTransactionTypes,
      selectedAccounts,
      descriptionFilter,
      showOnlyUnlabeled,
    }));
  }, [selectedLabels, selectedTransactionTypes, selectedAccounts, descriptionFilter, showOnlyUnlabeled, dispatch]);

  const handleRemoveLabel = (transaction: ITransaction, label: ILabelTransaction) => {
    setLabelToRemove({ transaction, label });
    setIsRemoveLabelModalOpen(true);
  };

  const handleConfirmRemoveLabel = () => {
    if (labelToRemove) {
      dispatch(removeLabelFromTransaction({
        transactionHref: labelToRemove.transaction.href,
        key: labelToRemove.label.key,
        value: labelToRemove.label.value,
      }));
    }
    setIsRemoveLabelModalOpen(false);
    setLabelToRemove(null);
  };

  const handleOpenTransactionDetails = (transaction: ITransaction) => {
    setSelectedTransactionForDetails(transaction);
    setIsTransactionDetailsModalOpen(true);
  };

  const handleOpenCreateRule = (transaction: ITransaction) => {
    setSelectedTransactionForRule(transaction);
    setIsCreateRuleModalOpen(true);
  };

  const handleOpenBulkLabel = () => {
    if (selectedItems.length > 0) {
      setIsBulkLabelModalOpen(true);
    }
  };

  const handleLabelClick = (label: ILabelTransaction) => {
    const labelFilter = `${label.key}=${label.value}`;
    
    // Check if the label is already in the filter
    if (!selectedLabels.includes(labelFilter)) {
      // Add the label to the existing filters
      const newSelectedLabels = [...selectedLabels, labelFilter];
      dispatch(setSelectedLabels(newSelectedLabels));
    }
  };

  // Define table columns
  const columns: Array<EuiBasicTableColumn<ITransaction>> = [
    {
      field: 'date',
      name: 'Date',
      sortable: true,
      render: (date: string) => safeFormatDateString(date),
      width: '120px',
    },
    {
      field: 'account',
      name: 'Account',
      sortable: true,
      render: (account: string) => (
        <EuiBadge 
          color={theme === 'dark' ? getAccountDarkColor(parseInt(account) || 0) : getAccountColor(parseInt(account) || 0)}
          style={{ borderRadius: '12px' }}
        >
          {account}
        </EuiBadge>
      ),
      width: '140px',
    },
    {
      field: 'kind',
      name: 'Type',
      sortable: true,
      render: (kind: string) => (
        <EuiBadge 
          color={kind === 'debit' ? 'danger' : 'success'}
          style={{ borderRadius: '12px' }}
        >
          {kind}
        </EuiBadge>
      ),
      width: '80px',
    },
    {
      field: 'description',
      name: 'Description',
      render: (description: string) => (
        <EuiText 
          size="s" 
          style={{ 
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis'
          }}
        >
          {description}
        </EuiText>
      ),
    },
    {
      field: 'labels',
      name: 'Labels',
      width: '300px',
      render: (labels: ILabelTransaction[], transaction: ITransaction) => (
        <EuiFlexGroup gutterSize="xs" wrap>
          {labels.map((label, index) => (
            <EuiFlexItem grow={false} key={index}>
              <EuiBadge 
                color={getLabelColor(label.key)}
                onClickAriaLabel={`Filter by label ${label.key}=${label.value}`}
                onClick={() => handleLabelClick(label)}
                iconType="cross"
                iconSide="right"
                iconOnClick={(e) => {
                  e.stopPropagation(); // Prevent the main click from firing
                  handleRemoveLabel(transaction, label);
                }}
                iconOnClickAriaLabel={`Remove label ${label.key}=${label.value}`}
                style={{ borderRadius: '12px', cursor: 'pointer' }}
              >
                {label.key}={label.value}
              </EuiBadge>
            </EuiFlexItem>
          ))}
        </EuiFlexGroup>
      ),
    },
    {
      field: 'amount',
      name: 'Amount',
      sortable: true,
      render: (amount: number) => (
        <EuiText size="s" style={{ fontFamily: 'monospace' }}>
          {amount.toFixed(2)}
        </EuiText>
      ),
      width: '100px',
    },
  ];

  // Actions column
  const actions: EuiTableActionsColumnType<ITransaction>['actions'] = [
    {
      name: 'Edit',
      description: 'Edit transaction details',
      icon: 'pencil',
      type: 'icon',
      onClick: (transaction) => handleOpenTransactionDetails(transaction),
    },
    {
      name: 'Add Label',
      description: 'Add label to transaction',
      icon: 'tag',
      type: 'icon',
      onClick: (transaction) => {
        setSelectedTransactionForDetails(transaction);
        setIsTransactionDetailsModalOpen(true);
      },
    },
    {
      name: 'Create Rule',
      description: 'Create rule from transaction',
      icon: 'gear',
      type: 'icon',
      onClick: (transaction) => handleOpenCreateRule(transaction),
    },
  ];

  const actionsColumn: EuiTableActionsColumnType<ITransaction> = {
    actions,
    width: '80px',
  };

  const allColumns = [...columns, actionsColumn];

  // Selection configuration
  const selection = {
    onSelectionChange: (selectedTransactions: ITransaction[]) => {
      setSelectedItems(selectedTransactions);
    },
  };

  // Toolbar with filters
  const renderToolbar = () => (
    <EuiPanel paddingSize="s">
      <EuiFlexGroup gutterSize="s" alignItems="flexEnd" wrap>
        <EuiFlexItem style={{ minWidth: '200px' }}>
          <EuiFormRow label="Description Filter">
            <EuiFieldText
              placeholder="Filter by description..."
              value={descriptionFilter}
              onChange={(e) => dispatch(setDescriptionFilter(e.target.value))}
              compressed
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem style={{ minWidth: '200px' }}>
          <EuiFormRow label="Labels">
            <LabelFilter
              availableLabels={availableLabels}
              selectedLabels={selectedLabels}
              onLabelsChange={(labels) => dispatch(setSelectedLabels(labels))}
              placeholder="Filter by labels..."
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem style={{ minWidth: '150px' }}>
          <EuiFormRow label="Account">
            <EuiSelect
              options={[
                { value: '', text: 'All accounts' },
                ...availableAccounts.map(account => ({ value: account, text: account }))
              ]}
              value={selectedAccounts[0]?.toString() || ''}
              onChange={(e) => dispatch(setSelectedAccounts(e.target.value ? [parseInt(e.target.value) || 0] : []))}
              compressed
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem style={{ minWidth: '120px' }}>
          <EuiFormRow label="Type">
            <EuiSelect
              options={[
                { value: '', text: 'All types' },
                ...availableTransactionTypes.map(type => ({ value: type, text: type }))
              ]}
              value={selectedTransactionTypes[0] || ''}
              onChange={(e) => dispatch(setSelectedTransactionTypes(e.target.value ? [e.target.value] : []))}
              compressed
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            <EuiSwitch
              label="Unlabeled only"
              checked={showOnlyUnlabeled}
              onChange={(e) => dispatch(setShowOnlyUnlabeled(e.target.checked))}
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            <EuiButton size="s" onClick={() => dispatch(clearAllFilters())}>
              Clear Filters
            </EuiButton>
          </EuiFormRow>
        </EuiFlexItem>

        {selectedItems.length > 0 && (
          <EuiFlexItem grow={false}>
            <EuiFormRow hasEmptyLabelSpace>
              <EuiButton size="s" fill onClick={handleOpenBulkLabel}>
                Add Labels to {selectedItems.length} Transaction{selectedItems.length !== 1 ? 's' : ''}
              </EuiButton>
            </EuiFormRow>
          </EuiFlexItem>
        )}
      </EuiFlexGroup>
    </EuiPanel>
  );

  const sorting = {
    sort: {
      field: ['date', 'account', 'kind', 'amount'][sortIndex] as keyof ITransaction,
      direction: sortDirection as 'asc' | 'desc',
    },
  };

  const pagination = {
    pageIndex: page - 1,
    pageSize: perPage,
    totalItemCount: filteredTransactions.length,
    showPerPageOptions: true,
    pageSizeOptions: [10, 25, 50, 100],
  };

  const onTableChange = ({ page: newPage, sort }: any) => {
    if (newPage) {
      dispatch(setPage(newPage.index + 1));
      dispatch(setPerPage(newPage.size));
    }
    if (sort) {
      const fieldToIndex = {
        'date': 0,
        'account': 1,
        'kind': 2,
        'amount': 3,
      };
      const newSortIndex = fieldToIndex[sort.field as keyof typeof fieldToIndex] ?? 0;
      dispatch(setSorting({ sortIndex: newSortIndex, sortDirection: sort.direction }));
    }
  };

  if (filtering) {
    return (
      <EuiFlexGroup justifyContent="center" alignItems="center" style={{ minHeight: '200px' }}>
        <EuiFlexItem grow={false}>
          <EuiLoadingSpinner size="l" />
        </EuiFlexItem>
      </EuiFlexGroup>
    );
  }

  return (
    <>
      {renderToolbar()}
      <EuiSpacer size="m" />
      
      {errorMessage && (
        <>
          <EuiCallOut title="Error" color="danger" iconType="alert">
            {errorMessage}
          </EuiCallOut>
          <EuiSpacer size="m" />
        </>
      )}

      <EuiInMemoryTable
        items={filteredTransactions}
        columns={allColumns}
        selection={selection}
        sorting={sorting}
        pagination={pagination}
        onChange={onTableChange}
        loading={filtering}
        message={filteredTransactions.length === 0 ? "No transactions found" : undefined}
      />

      {/* Modals */}
      <TransactionDetailsModal
        isOpen={isTransactionDetailsModalOpen}
        onClose={() => {
          setIsTransactionDetailsModalOpen(false);
          setSelectedTransactionForDetails(null);
        }}
        transaction={selectedTransactionForDetails || undefined}
        onSuccess={() => {
          dispatch(clearAddLabelToTransactionSuccess());
        }}
      />

      <AddLabelModal
        isOpen={isBulkLabelModalOpen}
        onClose={() => setIsBulkLabelModalOpen(false)}
        transactionHrefs={selectedItems.map(t => t.href)}
        onSuccess={() => {
          setIsBulkLabelModalOpen(false);
          setSelectedItems([]);
          dispatch(clearAddLabelToTransactionSuccess());
        }}
      />

      <RemoveLabelModal
        isOpen={isRemoveLabelModalOpen}
        onClose={() => {
          setIsRemoveLabelModalOpen(false);
          setLabelToRemove(null);
        }}
        onConfirm={handleConfirmRemoveLabel}
        transaction={labelToRemove?.transaction}
        label={labelToRemove?.label}
      />

      <CreateRuleModal
        isOpen={isCreateRuleModalOpen}
        onClose={() => {
          setIsCreateRuleModalOpen(false);
          setSelectedTransactionForRule(null);
        }}
        transaction={selectedTransactionForRule || undefined}
      />
    </>
  );
};

export { TransactionList };