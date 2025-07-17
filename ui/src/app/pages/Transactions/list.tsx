import * as React from 'react';
import {
  Content,
  Flex,
  FlexItem,
  Label,
  PageSection,
  Pagination,
  PaginationVariant,
  Select,
  SelectOption,
  SelectList,
  MenuToggle,
  MenuToggleElement,
  Button,
} from '@patternfly/react-core';
import { DataView, DataViewToolbar } from '@patternfly/react-data-view';
import { ExpandableRowContent, Table, Tbody, Td, Th, Thead, ThProps, Tr } from '@patternfly/react-table';
import { ILabelTransaction, ITransaction } from '@app/shared/models/transaction';
import { LabelFilter } from '@app/shared/components/label-filter';
import { useTheme } from '@app/shared/contexts/ThemeContext';
import { getAccountColor, getAccountDarkColor } from '@app/utils/colorUtils';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { setSourceTransactions, applyFilters } from '@app/shared/reducers/transaction-filter.reducer';

export interface ITransactionListProps {
  transactions: Array<ITransaction> | [];
}

// Table column definitions for consistent referencing
const columns = {
  date: 'Date',
  account: 'Account',
  kind: 'Type',
  amount: 'Amount',
  labels: 'Labels',
  rules: 'Rules',
};

/**
 * TransactionList Component
 *
 * This component displays a paginated, filterable, and sortable table of transactions.
 * Features:
 * - Filtering by labels, transaction types, and accounts
 * - Sorting by various columns
 * - Pagination support
 * - Expandable rows showing transaction descriptions
 * - Expand all/collapse all functionality
 * - Interactive labels and filters for quick filtering
 */
const TransactionList: React.FunctionComponent<ITransactionListProps> = ({ transactions }) => {
  const { theme } = useTheme();
  const dispatch = useAppDispatch();

  // Get filtered transactions from the filter reducer
  const { filteredTransactions } = useAppSelector((state) => state.transactionFilter);

  // ===============================
  // STATE MANAGEMENT
  // ===============================

  // Filter states
  const [selectedLabels, setSelectedLabels] = React.useState<string[]>([]);
  const [selectedTransactionTypes, setSelectedTransactionTypes] = React.useState<string[]>([]);
  const [selectedAccounts, setSelectedAccounts] = React.useState<number[]>([]);

  // UI states for dropdown controls
  const [isTransactionTypeSelectOpen, setIsTransactionTypeSelectOpen] = React.useState(false);
  const [isAccountSelectOpen, setIsAccountSelectOpen] = React.useState(false);

  // Pagination states
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(10);

  // Sorting states
  const [sortDirection, setSortDirection] = React.useState<'asc' | 'desc' | null>(null);
  const [sortIndex, setSortIndex] = React.useState<number | null>(null);

  // Row expansion states
  const [expandedTransactions, setExpandedTransactions] = React.useState<string[]>([]);
  const [allExpanded, setAllExpanded] = React.useState<boolean>(false);

  // ===============================
  // COMPUTED VALUES
  // ===============================

  // Calculate available labels from all transactions for filtering
  const availableLabels = React.useMemo(() => {
    const labelSet = new Set<string>();
    transactions.forEach((transaction) => {
      transaction.labels.forEach((label) => {
        labelSet.add(`${label.key}=${label.value}`);
      });
    });
    return Array.from(labelSet).sort();
  }, [transactions]);

  // Calculate available transaction types from all transactions
  const availableTransactionTypes = React.useMemo(() => {
    const typeSet = new Set<string>();
    transactions.forEach((transaction) => {
      typeSet.add(transaction.kind);
    });
    return Array.from(typeSet).sort();
  }, [transactions]);

  // Calculate available accounts from all transactions
  const availableAccounts = React.useMemo(() => {
    const accountSet = new Set<number>();
    transactions.forEach((transaction) => {
      accountSet.add(transaction.account);
    });
    return Array.from(accountSet).sort();
  }, [transactions]);

  // ===============================
  // UTILITY FUNCTIONS
  // ===============================

  /**
   * Get color for transaction type labels
   * @param kind - Transaction type (EXPENSE/INCOME)
   * @returns Color name for the label
   */
  const getTransactionKindColor = (kind: string): 'red' | 'blue' => {
    return kind === 'EXPENSE' ? 'red' : 'blue';
  };

  /**
   * Set transaction expanded state
   * @param t - Transaction to expand/collapse
   * @param isExpanding - Whether to expand (true) or collapse (false)
   */
  const setTransactionExpanded = (t: ITransaction, isExpanding = true) =>
    setExpandedTransactions((prevExpanded) => {
      const transactionIndex = prevExpanded.findIndex((href) => href === t.href);
      const newExpanded = [...prevExpanded];
      if (isExpanding && transactionIndex === -1) {
        newExpanded.push(t.href);
      } else if (!isExpanding && transactionIndex !== -1) {
        newExpanded.splice(transactionIndex, 1);
      }
      return newExpanded;
    });

  /**
   * Check if a transaction is currently expanded
   * @param t - Transaction to check
   * @returns Whether the transaction is expanded
   */
  const isTransactionExpanded = (t: ITransaction) => expandedTransactions.includes(t.href);

  // ===============================
  // EFFECTS
  // ===============================

  // Apply filters whenever filter selections change
  React.useEffect(() => {
    // Set source transactions for the filter reducer
    dispatch(setSourceTransactions(transactions));

    // Apply filters
    dispatch(
      applyFilters({
        selectedLabels,
        selectedTransactionTypes,
        selectedAccounts,
      })
    );
  }, [selectedLabels, selectedTransactionTypes, selectedAccounts, dispatch]);

  // ===============================
  // SORTING AND PAGINATION
  // ===============================

  /**
   * Sort transactions based on column and direction
   * @param transactionsToSort - Array of transactions to sort
   * @param sortIndex - Column index to sort by
   * @param sortDirection - Sort direction (asc/desc)
   * @returns Sorted array of transactions
   */
  const sortTransactions = (
    transactionsToSort: ITransaction[],
    sortIndex: number | null,
    sortDirection: 'asc' | 'desc' | null
  ): ITransaction[] => {
    if (sortIndex === null || sortDirection === null) {
      return transactionsToSort;
    }

    return [...transactionsToSort].sort((a, b) => {
      let aValue: any;
      let bValue: any;

      switch (sortIndex) {
        case 0: // Date
          aValue = new Date(a.date);
          bValue = new Date(b.date);
          break;
        case 1: // Account
          aValue = a.account;
          bValue = b.account;
          break;
        case 2: // Type
          aValue = a.kind;
          bValue = b.kind;
          break;
        case 3: // Amount
          aValue = Math.abs(a.amount);
          bValue = Math.abs(b.amount);
          break;
        case 4: // Labels
          aValue = a.labels.length;
          bValue = b.labels.length;
          break;
        case 5: // Rules (based on unique rules from labels)
          aValue = Array.from(new Set(a.labels.map((label) => label.ruleHref))).length;
          bValue = Array.from(new Set(b.labels.map((label) => label.ruleHref))).length;
          break;
        default:
          return 0;
      }

      if (aValue < bValue) return sortDirection === 'asc' ? -1 : 1;
      if (aValue > bValue) return sortDirection === 'asc' ? 1 : -1;
      return 0;
    });
  };

  // Apply sorting to filtered transactions
  const sortedTransactions = React.useMemo(() => {
    return sortTransactions(filteredTransactions, sortIndex, sortDirection);
  }, [filteredTransactions, sortIndex, sortDirection]);

  // Apply pagination to sorted transactions
  const paginatedTransactions = React.useMemo(() => {
    const startIdx = ((page || 1) - 1) * perPage;
    const endIdx = startIdx + perPage;
    return sortedTransactions.slice(startIdx, endIdx);
  }, [sortedTransactions, page, perPage]);

  // ===============================
  // EXPAND ALL FUNCTIONALITY
  // ===============================

  // Check if all current page transactions are expanded
  const areAllCurrentPageExpanded = React.useMemo(() => {
    return paginatedTransactions.length > 0 && paginatedTransactions.every(t => expandedTransactions.includes(t.href));
  }, [paginatedTransactions, expandedTransactions]);

  /**
   * Handle expand all / collapse all functionality
   * Only affects transactions on the current page
   */
  const handleExpandAll = () => {
    if (areAllCurrentPageExpanded) {
      // Collapse all current page transactions
      const currentPageHrefs = paginatedTransactions.map(t => t.href);
      setExpandedTransactions(prev => prev.filter(href => !currentPageHrefs.includes(href)));
    } else {
      // Expand all current page transactions
      const currentPageHrefs = paginatedTransactions.map(t => t.href);
      setExpandedTransactions(prev => Array.from(new Set([...prev, ...currentPageHrefs])));
    }
  };

  // ===============================
  // FILTER EVENT HANDLERS
  // ===============================

  /**
   * Handle label filter changes
   * @param labels - Array of selected label strings in key=value format
   */
  const handleLabelsChange = (labels: string[]) => {
    console.log('Selected labels changed:', labels);
    setSelectedLabels(labels);
  };

  // Transaction type filter handlers
  const handleTransactionTypeToggle = () => {
    setIsTransactionTypeSelectOpen(!isTransactionTypeSelectOpen);
  };

  const handleTransactionTypeSelect = (
    _event: React.MouseEvent<Element, MouseEvent> | undefined,
    value: string | number | undefined
  ) => {
    const stringValue = String(value);
    if (selectedTransactionTypes.includes(stringValue)) {
      setSelectedTransactionTypes(selectedTransactionTypes.filter((type) => type !== stringValue));
    } else {
      setSelectedTransactionTypes([...selectedTransactionTypes, stringValue]);
    }
  };

  const handleTransactionTypeRemove = (typeToRemove: string) => {
    setSelectedTransactionTypes(selectedTransactionTypes.filter((type) => type !== typeToRemove));
  };

  // Account filter handlers
  const handleAccountToggle = () => {
    setIsAccountSelectOpen(!isAccountSelectOpen);
  };

  const handleAccountSelect = (
    _event: React.MouseEvent<Element, MouseEvent> | undefined,
    value: string | number | undefined
  ) => {
    const numberValue = Number(value);
    if (selectedAccounts.includes(numberValue)) {
      setSelectedAccounts(selectedAccounts.filter((account) => account !== numberValue));
    } else {
      setSelectedAccounts([...selectedAccounts, numberValue]);
    }
  };

  const handleAccountRemove = (accountToRemove: number) => {
    setSelectedAccounts(selectedAccounts.filter((account) => account !== accountToRemove));
  };

  /**
   * Clear all active filters
   */
  const handleClearAllFilters = () => {
    setSelectedLabels([]);
    setSelectedTransactionTypes([]);
    setSelectedAccounts([]);
  };

  // ===============================
  // INTERACTIVE CLICK HANDLERS
  // ===============================

  /**
   * Handle clicking on a label to add/remove it from filters
   * @param labelValue - Label string in key=value format
   */
  const handleLabelClick = (labelValue: string) => {
    const newSelectedLabels = selectedLabels.includes(labelValue)
      ? selectedLabels.filter((label) => label !== labelValue)
      : [...selectedLabels, labelValue];
    setSelectedLabels(newSelectedLabels);
  };

  /**
   * Handle clicking on a rule to filter by all labels associated with that rule
   * @param ruleId - Rule identifier
   */
  const handleRuleClick = (ruleId: string) => {
    // Find all labels associated with this rule
    const ruleLabels = new Set<string>();
    transactions.forEach((transaction) => {
      transaction.labels.forEach((label) => {
        if (label.ruleHref === ruleId) {
          ruleLabels.add(`${label.key}=${label.value}`);
        }
      });
    });

    const ruleLabelsArray = Array.from(ruleLabels);
    const newSelectedLabels = Array.from(new Set([...selectedLabels, ...ruleLabelsArray]));
    setSelectedLabels(newSelectedLabels);
  };

  /**
   * Handle clicking on transaction type to add/remove it from filters
   * @param transactionType - Transaction type string
   */
  const handleTransactionTypeClick = (transactionType: string) => {
    const newSelectedTypes = selectedTransactionTypes.includes(transactionType)
      ? selectedTransactionTypes.filter((type) => type !== transactionType)
      : [...selectedTransactionTypes, transactionType];
    setSelectedTransactionTypes(newSelectedTypes);
  };

  /**
   * Handle clicking on account to add/remove it from filters
   * @param accountNumber - Account number
   */
  const handleAccountClick = (accountNumber: number) => {
    const newSelectedAccounts = selectedAccounts.includes(accountNumber)
      ? selectedAccounts.filter((account) => account !== accountNumber)
      : [...selectedAccounts, accountNumber];
    setSelectedAccounts(newSelectedAccounts);
  };

  // ===============================
  // SORTING HELPERS
  // ===============================

  /**
   * Generate sorting parameters for table headers
   * @param columnIndex - Index of the column
   * @returns Sort parameters object
   */
  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: sortIndex || undefined,
      direction: sortDirection || undefined,
    },
    onSort: (_event, index, direction) => {
      setSortIndex(index);
      setSortDirection(direction);
    },
    columnIndex,
  });

  // ===============================
  // PAGINATION HANDLERS
  // ===============================

  const handleSetPage = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPage: number,
    _perPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPage(newPage);
  };

  const handlePerPageSelect = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPerPage: number,
    newPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPerPage(newPerPage);
    setPage(newPage);
  };

  // ===============================
  // RENDER HELPERS
  // ===============================

  /**
   * Render pagination component
   * @param variant - Pagination variant (top/bottom)
   * @param isCompact - Whether to use compact layout
   * @param isSticky - Whether pagination should be sticky
   * @param isStatic - Whether pagination should be static
   * @returns Pagination component
   */
  const renderPagination = (variant: PaginationVariant, isCompact: boolean, isSticky: boolean, isStatic: boolean) => (
    <Pagination
      id={`transaction-table-${variant}-pagination`}
      variant={variant}
      itemCount={sortedTransactions.length}
      page={page}
      perPage={perPage}
      isCompact={isCompact}
      isSticky={isSticky}
      isStatic={isStatic}
      onSetPage={handleSetPage}
      onPerPageSelect={handlePerPageSelect}
      titles={{
        paginationAriaLabel: `${variant} pagination`,
      }}
    />
  );

  /**
   * Render the toolbar with filters and controls
   */
  const renderToolbar = (
    <DataViewToolbar
      clearAllFilters={handleClearAllFilters}
      pagination={renderPagination(PaginationVariant.top, true, false, false)}
      filters={
        <React.Fragment>
          <LabelFilter
            availableLabels={availableLabels}
            selectedLabels={selectedLabels}
            onLabelsChange={handleLabelsChange}
            placeholder="Filter by labels..."
          />
          <Select
            id="transaction-type-select"
            isOpen={isTransactionTypeSelectOpen}
            selected={selectedTransactionTypes}
            onSelect={handleTransactionTypeSelect}
            onOpenChange={(isOpen) => setIsTransactionTypeSelectOpen(isOpen)}
            toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
              <MenuToggle
                ref={toggleRef}
                onClick={handleTransactionTypeToggle}
                isExpanded={isTransactionTypeSelectOpen}
                style={{ width: '200px' }}
              >
                Transaction Type
              </MenuToggle>
            )}
            shouldFocusToggleOnSelect
          >
            <SelectList>
              {availableTransactionTypes.map((type, index) => (
                <SelectOption key={index} value={type}>
                  {type}
                </SelectOption>
              ))}
            </SelectList>
          </Select>
          <Select
            id="account-select"
            isOpen={isAccountSelectOpen}
            selected={selectedAccounts}
            onSelect={handleAccountSelect}
            onOpenChange={(isOpen) => setIsAccountSelectOpen(isOpen)}
            toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
              <MenuToggle
                ref={toggleRef}
                onClick={handleAccountToggle}
                isExpanded={isAccountSelectOpen}
                style={{ width: '200px' }}
              >
                Account
              </MenuToggle>
            )}
            shouldFocusToggleOnSelect
          >
            <SelectList>
              {availableAccounts.map((account, index) => (
                <SelectOption key={index} value={account}>
                  {account}
                </SelectOption>
              ))}
            </SelectList>
          </Select>
        </React.Fragment>
      }
    />
  );

  /**
   * Render the active filters section
   */
  const renderSelectedFilters = () => (
    <PageSection>
      <Flex direction={{ default: 'column' }}>
        {(selectedLabels.length > 0 || selectedTransactionTypes.length > 0 || selectedAccounts.length > 0) && (
          <FlexItem>
            <Content>
              <strong>Active Filters:</strong>
            </Content>
          </FlexItem>
        )}

        {/* Selected Labels */}
        {selectedLabels.length > 0 && (
          <FlexItem>
            <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
              <FlexItem>
                <Content>Labels:</Content>
              </FlexItem>
              {selectedLabels.map((label, index) => (
                <FlexItem key={index}>
                  <Label
                    variant="filled"
                    color="blue"
                    onClose={() => handleLabelClick(label)}
                    closeBtnAriaLabel={`Remove ${label} filter`}
                  >
                    {label}
                  </Label>
                </FlexItem>
              ))}
            </Flex>
          </FlexItem>
        )}

        {/* Selected Transaction Types */}
        {selectedTransactionTypes.length > 0 && (
          <FlexItem>
            <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
              <FlexItem>
                <Content>Transaction Types:</Content>
              </FlexItem>
              {selectedTransactionTypes.map((type, index) => (
                <FlexItem key={index}>
                  <Label
                    variant="filled"
                    color="green"
                    onClose={() => handleTransactionTypeRemove(type)}
                    closeBtnAriaLabel={`Remove ${type} filter`}
                  >
                    {type}
                  </Label>
                </FlexItem>
              ))}
            </Flex>
          </FlexItem>
        )}

        {/* Selected Accounts */}
        {selectedAccounts.length > 0 && (
          <FlexItem>
            <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
              <FlexItem>
                <Content>Accounts:</Content>
              </FlexItem>
              {selectedAccounts.map((account, index) => (
                <FlexItem key={index}>
                  <Label
                    variant="filled"
                    color="purple"
                    onClose={() => handleAccountRemove(account)}
                    closeBtnAriaLabel={`Remove account ${account} filter`}
                  >
                    {account}
                  </Label>
                </FlexItem>
              ))}
            </Flex>
          </FlexItem>
        )}
      </Flex>
    </PageSection>
  );

  /**
   * Render the main transaction table
   */
  const renderList = (
    <React.Fragment>
      <Table aria-label="transaction-list">
        <Thead>
          <Tr>
            {/* Expand All Button */}
            <Th>
              <Button
                variant="plain"
                onClick={handleExpandAll}
                size="sm"
                aria-label={areAllCurrentPageExpanded ? 'Collapse all rows' : 'Expand all rows'}
              >
                {areAllCurrentPageExpanded ? '▼' : '▶'}
              </Button>
            </Th>
            {/* Sortable Columns */}
            <Th sort={getSortParams(0)}>
              <strong>{columns.date}</strong>
            </Th>
            <Th sort={getSortParams(1)}>
              <strong>{columns.account}</strong>
            </Th>
            <Th sort={getSortParams(2)}>
              <strong>{columns.kind}</strong>
            </Th>
            <Th sort={getSortParams(4)}>
              <strong>{columns.labels}</strong>
            </Th>
            <Th sort={getSortParams(3)}>
              <strong>{columns.amount}</strong>
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          {paginatedTransactions.map((t: ITransaction, rowIndex: number) => (
            <React.Fragment key={t.href}>
              {/* Main Transaction Row */}
              <Tr>
                {/* Expand/Collapse Button */}
                <Td
                  expand={{
                    rowIndex,
                    isExpanded: isTransactionExpanded(t),
                    onToggle: () => setTransactionExpanded(t, !isTransactionExpanded(t)),
                  }}
                />
                {/* Date Column */}
                <Td dataLabel={columns.date}>
                  {new Date(t.date).toLocaleDateString('en-US', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                  })}
                </Td>
                {/* Account Column - Clickable label for filtering */}
                <Td dataLabel={columns.account}>
                  <Label
                    variant={theme === 'dark' ? 'outline' : 'filled'}
                    color="purple"
                    onClick={() => handleAccountClick(t.account)}
                    style={{
                      cursor: 'pointer',
                      color: theme === 'dark' ? getAccountDarkColor(t.account) : 'black',
                    }}
                    aria-label={`Filter by account ${t.account}`}
                  >
                    {t.account}
                  </Label>
                </Td>
                {/* Transaction Type Column - Clickable label for filtering */}
                <Td dataLabel={columns.kind}>
                  <Label
                    variant={theme === 'dark' ? 'outline' : 'filled'}
                    color={getTransactionKindColor(t.kind)}
                    onClick={() => handleTransactionTypeClick(t.kind)}
                    style={{ cursor: 'pointer' }}
                  >
                    {t.kind}
                  </Label>
                </Td>
                {/* Labels Column - Clickable labels for filtering */}
                <Td dataLabel={columns.labels}>
                  <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                    {t.labels.map((label: ILabelTransaction, idx: number) => (
                      <FlexItem key={`label-${idx}`}>
                        <Label
                          variant={theme === 'dark' ? 'outline' : 'filled'}
                          color="green"
                          onClick={() => handleLabelClick(`${label.key}=${label.value}`)}
                          style={{ cursor: 'pointer' }}
                          aria-label={`Filter by ${label.key}=${label.value} label`}
                        >
                          {label.key}={label.value}
                        </Label>
                      </FlexItem>
                    ))}
                  </Flex>
                </Td>
                {/* Amount Column */}
                <Td dataLabel={columns.amount}>
                  <Content>
                    <strong>{t.amount.toFixed(2)}</strong>
                  </Content>
                </Td>
              </Tr>
              {/* Expandable Row Content - Shows transaction description */}
              <Tr isExpanded={isTransactionExpanded(t)}>
                <Td />
                <Td colSpan={6}>
                  <ExpandableRowContent>
                    <Content>{t.description || 'No description'}</Content>
                  </ExpandableRowContent>
                </Td>
              </Tr>
            </React.Fragment>
          ))}
        </Tbody>
      </Table>
    </React.Fragment>
  );

  // ===============================
  // MAIN COMPONENT RENDER
  // ===============================

  return (
    <React.Fragment>
      <DataView>
        {renderToolbar}
        {renderSelectedFilters()}
        {renderList}
        {renderPagination(PaginationVariant.bottom, false, false, true)}
      </DataView>
    </React.Fragment>
  );
};

export { TransactionList };
