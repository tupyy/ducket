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
import { ITagTransaction, ITransaction } from '@app/shared/models/transaction';
import { TagFilter } from '@app/shared/components/tag-filter';
import { useTheme } from '@app/shared/contexts/ThemeContext';
import { getAccountColor, getAccountDarkColor } from '@app/utils/colorUtils';

export interface ITransactionListProps {
  transactions: Array<ITransaction> | [];
}

const columns = {
  date: 'Date',
  account: 'Account',
  kind: 'Type',
  amount: 'Amount',
  tags: 'Tags',
  rules: 'Rules',
};

const TransactionList: React.FunctionComponent<ITransactionListProps> = ({ transactions }) => {
  const { theme } = useTheme();
  const [sortedTransactions, setSortedTransactions] = React.useState<Array<ITransaction>>(Array.from(transactions));
  const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(1);
  const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>('desc');
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(10);
  const [paginatedRows, setPaginatedRows] = React.useState(sortedTransactions.slice(0, 10));
  const [selectedTags, setSelectedTags] = React.useState<string[]>([]);
  const [selectedTransactionTypes, setSelectedTransactionTypes] = React.useState<string[]>([]);
  const [selectedAccounts, setSelectedAccounts] = React.useState<number[]>([]);
  const [filteredTransactions, setFilteredTransactions] = React.useState<Array<ITransaction>>([]);
  const [isTransactionTypeSelectOpen, setIsTransactionTypeSelectOpen] = React.useState(false);
  const [isAccountSelectOpen, setIsAccountSelectOpen] = React.useState(false);

  // Helper function to get transaction kind label color
  const getTransactionKindColor = (kind: string): 'red' | 'blue' => {
    return kind.toLowerCase() === 'debit' ? 'red' : 'blue';
  };

  const [expandedTransactions, setExpandedTransactions] = React.useState<string[]>([]);
  const setTransactionExpanded = (t: ITransaction, isExpanding = true) =>
    setExpandedTransactions((prevExpanded) => {
      const otherExpandedRepoNames = prevExpanded.filter((tt) => tt !== t.href);
      return isExpanding ? [...otherExpandedRepoNames, t.href] : otherExpandedRepoNames;
    });
  const isTransactionExpanded = (t: ITransaction) => expandedTransactions.includes(t.href);

  // Sort transactions by initial sort state when component loads
  React.useEffect(() => {
    const initialSorted = [...transactions].sort((a, b) => {
      if (activeSortIndex === 1) { // Date column
        const aValue = a.date.getTime();
        const bValue = b.date.getTime();
        
        if (activeSortDirection === 'desc') {
          return bValue - aValue;
        }
        return aValue - bValue;
      }
      return 0; // No sorting for other columns initially
    });
    
    setSortedTransactions(initialSorted);
  }, [transactions, activeSortIndex, activeSortDirection]);

  // Client-side filtering by tags, transaction types, and accounts
  React.useEffect(() => {
    let filtered = sortedTransactions;

    // Filter by tags
    if (selectedTags.length > 0) {
      filtered = filtered.filter((transaction) =>
        selectedTags.some((selectedTag) => transaction.tags.some((tag) => tag.value === selectedTag)),
      );
    }

    // Filter by transaction types
    if (selectedTransactionTypes.length > 0) {
      filtered = filtered.filter((transaction) => selectedTransactionTypes.includes(transaction.kind));
    }

    // Filter by accounts
    if (selectedAccounts.length > 0) {
      filtered = filtered.filter((transaction) => selectedAccounts.includes(transaction.account));
    }

    setFilteredTransactions(filtered);
  }, [sortedTransactions, selectedTags, selectedTransactionTypes, selectedAccounts]);

  React.useEffect(() => {
    setPaginatedRows(filteredTransactions?.slice(0, perPage));
    setPage(1);
  }, [filteredTransactions, perPage]);

  // Get available tag values from the actual transactions
  const availableTags = React.useMemo(() => {
    const tagSet = new Set<string>();
    transactions.forEach((transaction) => {
      transaction.tags.forEach((tag) => {
        tagSet.add(tag.value);
      });
    });
    return Array.from(tagSet).sort();
  }, [transactions]);

  // Get available transaction types
  const availableTransactionTypes = React.useMemo(() => {
    const typeSet = new Set<string>();
    transactions.forEach((transaction) => {
      typeSet.add(transaction.kind);
    });
    return Array.from(typeSet).sort();
  }, [transactions]);

  // Get available accounts
  const availableAccounts = React.useMemo(() => {
    const accountSet = new Set<number>();
    transactions.forEach((transaction) => {
      if (transaction.account) {
        accountSet.add(transaction.account);
      }
    });
    return Array.from(accountSet).sort((a, b) => a - b);
  }, [transactions]);

  const handleTagsChange = (tags: string[]) => {
    console.log('Selected tags changed:', tags);
    setSelectedTags(tags);
  };

  const handleTransactionTypeToggle = () => {
    setIsTransactionTypeSelectOpen(!isTransactionTypeSelectOpen);
  };

  const handleTransactionTypeSelect = (
    _event: React.MouseEvent<Element, MouseEvent> | undefined,
    value: string | number | undefined,
  ) => {
    if (typeof value === 'string') {
      setSelectedTransactionTypes((prev) => {
        if (prev.includes(value)) {
          return prev.filter((type) => type !== value);
        } else {
          return [...prev, value];
        }
      });
    }
  };

  const handleTransactionTypeRemove = (typeToRemove: string) => {
    setSelectedTransactionTypes((prev) => prev.filter((type) => type !== typeToRemove));
  };

  const handleAccountToggle = () => {
    setIsAccountSelectOpen(!isAccountSelectOpen);
  };

  const handleAccountSelect = (
    _event: React.MouseEvent<Element, MouseEvent> | undefined,
    value: string | number | undefined,
  ) => {
    if (typeof value === 'string') {
      const accountNumber = parseInt(value, 10);
      if (!isNaN(accountNumber)) {
        setSelectedAccounts((prev) => {
          if (prev.includes(accountNumber)) {
            return prev.filter((account) => account !== accountNumber);
          } else {
            return [...prev, accountNumber];
          }
        });
      }
    }
  };

  const handleAccountRemove = (accountToRemove: number) => {
    setSelectedAccounts((prev) => prev.filter((account) => account !== accountToRemove));
  };

  const handleClearAllFilters = () => {
    setSelectedTags([]);
    setSelectedTransactionTypes([]);
    setSelectedAccounts([]);
  };

  const handleTagClick = (tagValue: string) => {
    if (!selectedTags.includes(tagValue)) {
      setSelectedTags((prev) => [...prev, tagValue]);
    }
  };

  const handleRuleClick = (ruleId: string) => {
    // Find all tags associated with this rule across all transactions
    const tagsForRule = new Set<string>();
    transactions.forEach((transaction) => {
      transaction.tags.forEach((tag) => {
        if (tag.rule === ruleId) {
          tagsForRule.add(tag.value);
        }
      });
    });

    // Add all tags for this rule to the selected tags (if not already selected)
    const newTags = Array.from(tagsForRule).filter((tag) => !selectedTags.includes(tag));
    if (newTags.length > 0) {
      setSelectedTags((prev) => [...prev, ...newTags]);
    }
  };

  const handleTransactionTypeClick = (transactionType: string) => {
    if (!selectedTransactionTypes.includes(transactionType)) {
      setSelectedTransactionTypes((prev) => [...prev, transactionType]);
    }
  };

  const handleAccountClick = (accountNumber: number) => {
    if (!selectedAccounts.includes(accountNumber)) {
      setSelectedAccounts((prev) => [...prev, accountNumber]);
    }
  };

  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: activeSortIndex || undefined,
      direction: activeSortDirection || undefined,
      defaultDirection: columnIndex === 1 ? 'desc' : 'asc', // Date column (index 1) defaults to desc
    },
    onSort: (_event, index, direction) => {
      const sorted = [...filteredTransactions].sort((a, b) => {
        let aValue: Date | number = a.date;
        let bValue: Date | number = b.date;
        if (index == 6) {
          aValue = a.amount;
          bValue = b.amount;
        }

        if (typeof aValue === 'number') {
          // Numeric sort
          if (direction === 'asc') {
            return (aValue as number) - (bValue as number);
          }
          return (bValue as number) - (aValue as number);
        }
        // date sort
        if (direction === 'asc') {
          return (aValue as Date).getTime() - (bValue as Date).getTime();
        }
        return (bValue as Date).getTime() - (aValue as Date).getTime();
      });

      setFilteredTransactions(sorted);
      setActiveSortIndex(index);
      setActiveSortDirection(direction);
      setPaginatedRows(sorted.slice(0, perPage));
      setPage(1);
    },
    columnIndex,
  });

  const handleSetPage = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPage: number,
    _perPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined,
  ) => {
    setPaginatedRows(filteredTransactions?.slice(startIdx, endIdx));
    setPage(newPage);
  };

  const handlePerPageSelect = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPerPage: number,
    newPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined,
  ) => {
    setPaginatedRows(filteredTransactions.slice(startIdx, endIdx));
    setPage(newPage);
    setPerPage(newPerPage);
  };

  const renderPagination = (variant: PaginationVariant, isCompact: boolean, isSticky: boolean, isStatic: boolean) => (
    <Pagination
      id={`datalist-${variant}-pagination`}
      variant={variant}
      itemCount={filteredTransactions.length}
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

  const renderToolbar = (
    <DataViewToolbar
      filters={
        <Flex direction={{ default: 'column' }} spaceItems={{ default: 'spaceItemsMd' }}>
          <FlexItem>
            <Flex spaceItems={{ default: 'spaceItemsMd' }}>
              <FlexItem>
                <TagFilter
                  availableTags={availableTags}
                  selectedTags={selectedTags}
                  onTagsChange={handleTagsChange}
                  placeholder="Filter by tags..."
                />
              </FlexItem>
              <FlexItem>
                <Flex direction={{ default: 'column' }}>
                  <FlexItem>
                    <Select
                      isOpen={isTransactionTypeSelectOpen}
                      selected={selectedTransactionTypes}
                      onSelect={handleTransactionTypeSelect}
                      onOpenChange={(isOpen) => setIsTransactionTypeSelectOpen(isOpen)}
                      toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                        <MenuToggle
                          ref={toggleRef}
                          onClick={handleTransactionTypeToggle}
                          isExpanded={isTransactionTypeSelectOpen}
                        >
                          {selectedTransactionTypes.length > 0
                            ? `Transaction Types (${selectedTransactionTypes.length})`
                            : 'Filter by transaction type...'}
                        </MenuToggle>
                      )}
                    >
                      <SelectList>
                        {availableTransactionTypes.map((type) => (
                          <SelectOption
                            key={type}
                            value={type}
                            isSelected={selectedTransactionTypes.includes(type)}
                            hasCheckbox
                          >
                            {type.charAt(0).toUpperCase() + type.slice(1)}
                          </SelectOption>
                        ))}
                      </SelectList>
                    </Select>
                  </FlexItem>
                  {selectedTransactionTypes.length > 0 && (
                    <FlexItem>
                      <Flex spaceItems={{ default: 'spaceItemsXs' }} style={{ marginTop: '8px' }}>
                        {selectedTransactionTypes.map((type, index) => (
                          <FlexItem key={index}>
                            <Label
                              variant={theme === 'dark' ? 'outline' : 'filled'}
                              color="orange"
                              onClose={() => handleTransactionTypeRemove(type)}
                              closeBtnAriaLabel={`Remove ${type} filter`}
                              style={theme === 'dark' ? { color: '#f4c430' } : {}}
                            >
                              {type.charAt(0).toUpperCase() + type.slice(1)}
                            </Label>
                          </FlexItem>
                        ))}
                      </Flex>
                    </FlexItem>
                  )}
                </Flex>
              </FlexItem>
              <FlexItem>
                <Flex direction={{ default: 'column' }}>
                  <FlexItem>
                    <Select
                      isOpen={isAccountSelectOpen}
                      selected={selectedAccounts}
                      onSelect={handleAccountSelect}
                      onOpenChange={(isOpen) => setIsAccountSelectOpen(isOpen)}
                      toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                        <MenuToggle
                          ref={toggleRef}
                          onClick={handleAccountToggle}
                          isExpanded={isAccountSelectOpen}
                        >
                          {selectedAccounts.length > 0
                            ? `Accounts (${selectedAccounts.length})`
                            : 'Filter by account...'}
                        </MenuToggle>
                      )}
                    >
                      <SelectList>
                        {availableAccounts.map((account) => (
                          <SelectOption
                            key={account}
                            value={account.toString()}
                            isSelected={selectedAccounts.includes(account)}
                            hasCheckbox
                          >
                            {account.toString()}
                          </SelectOption>
                        ))}
                      </SelectList>
                    </Select>
                  </FlexItem>
                  {selectedAccounts.length > 0 && (
                    <FlexItem>
                      <Flex spaceItems={{ default: 'spaceItemsXs' }} style={{ marginTop: '8px' }}>
                        {selectedAccounts.map((account, index) => (
                          <FlexItem key={index}>
                            <Label
                              variant={theme === 'dark' ? 'outline' : 'filled'}
                              color={getAccountColor(account)}
                              onClose={() => handleAccountRemove(account)}
                              closeBtnAriaLabel={`Remove ${account} filter`}
                              style={theme === 'dark' ? { color: getAccountDarkColor(account) } : {}}
                            >
                              {account.toString()}
                            </Label>
                          </FlexItem>
                        ))}
                      </Flex>
                    </FlexItem>
                  )}
                </Flex>
              </FlexItem>
            </Flex>
          </FlexItem>
          {(selectedTags.length > 0 || selectedTransactionTypes.length > 0 || selectedAccounts.length > 0) && (
            <FlexItem>
              <Button variant="link" onClick={handleClearAllFilters} isInline>
                Clear all filters
              </Button>
            </FlexItem>
          )}
        </Flex>
      }
      pagination={renderPagination(PaginationVariant.top, true, false, false)}
    />
  );

  const renderList = (
    <React.Fragment>
      <Table aria-label="transaction-list">
        <Thead>
          <Tr>
            <Th screenReaderText="Row expansion" />
            <Th sort={getSortParams(1)}>
              <Content component="p">
                <strong>{columns.date}</strong>
              </Content>
            </Th>
            <Th>
              <Content component="p">
                <strong>{columns.account}</strong>
              </Content>
            </Th>
            <Th>
              <Content component="p">
                <strong>{columns.kind}</strong>
              </Content>
            </Th>
            <Th>
              <Content component="p">
                <strong>{columns.tags}</strong>
              </Content>
            </Th>
            <Th>
              <Content component="p">
                <strong>{columns.rules}</strong>
              </Content>
            </Th>
            <Th width={10} sort={getSortParams(6)}>
              <Content component="p">
                <strong>{columns.amount}</strong>
              </Content>
            </Th>
          </Tr>
        </Thead>
        {paginatedRows.map((t: ITransaction, i: number) => (
          <Tbody key={t.href} isExpanded={isTransactionExpanded(t)}>
            <Tr>
              <Td
                expand={
                  t.description
                    ? {
                        rowIndex: i,
                        isExpanded: isTransactionExpanded(t),
                        onToggle: () => setTransactionExpanded(t, !isTransactionExpanded(t)),
                        expandId: 'composable-expandable-example',
                      }
                    : undefined
                }
              />
              <Td dataLabel={columns.date}>
                {t.date.toLocaleDateString('fr-FR', {
                  weekday: 'long',
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </Td>
              <Td dataLabel={columns.account}>
                <Label
                  variant={theme === 'dark' ? 'outline' : 'filled'}
                  color={getAccountColor(t.account)}
                  onClick={() => handleAccountClick(t.account)}
                  style={{
                    cursor: 'pointer',
                    ...(theme === 'dark' && { color: getAccountDarkColor(t.account) }),
                  }}
                  aria-label={`Filter by account ${t.account}`}
                >
                  {t.account}
                </Label>
              </Td>
              <Td dataLabel="{columns.kind}">
                <Label
                  variant={theme === 'dark' ? 'outline' : 'filled'}
                  color={getTransactionKindColor(t.kind)}
                  onClick={() => handleTransactionTypeClick(t.kind)}
                  style={{
                    cursor: 'pointer',
                  }}
                  aria-label={`Filter by ${t.kind} transactions`}
                >
                  {t.kind.charAt(0).toUpperCase() + t.kind.slice(1)}
                </Label>
              </Td>
              <Td dataLabel="{columns.tags}">
                <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  {t.tags.map((tag: ITagTransaction, idx: number) => (
                    <FlexItem key={`tag-${idx}`}>
                      <Label
                        variant={theme === 'dark' ? 'outline' : 'filled'}
                        color="green"
                        onClick={() => handleTagClick(tag.value)}
                        style={{
                          cursor: 'pointer',
                          ...(theme === 'dark' && { color: '#3e8635' }),
                        }}
                        aria-label={`Filter by ${tag.value} tag`}
                      >
                        <Content component="p" style={theme === 'dark' ? { color: 'white' } : { color: 'black' }}>
                          {tag.value}
                        </Content>
                      </Label>
                    </FlexItem>
                  ))}
                </Flex>
              </Td>
              <Td dataLabel="{columns.rules">
                <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  {Array.from(new Set(t.tags.map((tag) => tag.rule))).map((rule: string, idx: number) => (
                    <FlexItem key={`rule-${idx}`}>
                      <Label
                        variant={theme === 'dark' ? 'outline' : 'filled'}
                        color="blue"
                        onClick={() => handleRuleClick(rule)}
                        style={{
                          cursor: 'pointer',
                          ...(theme === 'dark' && { color: '#73bcf7' }),
                        }}
                        aria-label={`Filter by all tags associated with ${rule} rule`}
                      >
                        <Content component="p">{rule}</Content>
                      </Label>
                    </FlexItem>
                  ))}
                </Flex>
              </Td>
              <Td dataLabel="{columns.amount}">
                {t.amount.toLocaleString('fr-FR', {
                  style: 'currency',
                  currency: 'EUR',
                  minimumFractionDigits: 2,
                  maximumFractionDigits: 2,
                })}
              </Td>
            </Tr>
            <Tr isExpanded={isTransactionExpanded(t)}>
              <Td noPadding={false} colSpan={6}>
                <ExpandableRowContent>{t.description}</ExpandableRowContent>
              </Td>
            </Tr>
          </Tbody>
        ))}
      </Table>
    </React.Fragment>
  );

  return (
    <PageSection hasBodyWrapper={false}>
      <DataView>
        {renderToolbar}
        {renderList}
        {renderPagination(PaginationVariant.bottom, false, false, true)}
      </DataView>
    </PageSection>
  );
};

export { TransactionList };
