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

export interface ITransactionListProps {
  transactions: Array<ITransaction> | [];
}

const columns = {
  date: 'Date',
  kind: 'Type',
  amount: 'Amount',
  tags: 'Tags',
  rules: 'Rules',
};

const TransactionList: React.FunctionComponent<ITransactionListProps> = ({ transactions }) => {
  const [sortedTransactions, setSortedTransactions] = React.useState<Array<ITransaction>>(Array.from(transactions));
  const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
  const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(10);
  const [paginatedRows, setPaginatedRows] = React.useState(sortedTransactions.slice(0, 10));
  const [selectedTags, setSelectedTags] = React.useState<string[]>([]);
  const [selectedTransactionTypes, setSelectedTransactionTypes] = React.useState<string[]>([]);
  const [filteredTransactions, setFilteredTransactions] = React.useState<Array<ITransaction>>([]);
  const [isTransactionTypeSelectOpen, setIsTransactionTypeSelectOpen] = React.useState(false);

  // Date filter state
  const [startDate, setStartDate] = React.useState<string>('');
  const [endDate, setEndDate] = React.useState<string>('');

  const [expandedTransactions, setExpandedTransactions] = React.useState<string[]>([]);
  const setTransactionExpanded = (t: ITransaction, isExpanding = true) =>
    setExpandedTransactions((prevExpanded) => {
      const otherExpandedRepoNames = prevExpanded.filter((tt) => tt !== t.href);
      return isExpanding ? [...otherExpandedRepoNames, t.href] : otherExpandedRepoNames;
    });
  const isTransactionExpanded = (t: ITransaction) => expandedTransactions.includes(t.href);

  // Filter transactions based on date range
  React.useEffect(() => {
    setSortedTransactions(Array.from(transactions));
  }, [transactions]);

  // Client-side filtering by tags and transaction types
  React.useEffect(() => {
    let filtered = sortedTransactions;

    // Filter by tags
    if (selectedTags.length > 0) {
      filtered = filtered.filter(transaction =>
        selectedTags.some(selectedTag =>
          transaction.tags.some(tag => tag.value === selectedTag)
        )
      );
    }

    // Filter by transaction types
    if (selectedTransactionTypes.length > 0) {
      filtered = filtered.filter(transaction =>
        selectedTransactionTypes.includes(transaction.kind)
      );
    }

    setFilteredTransactions(filtered);
  }, [sortedTransactions, selectedTags, selectedTransactionTypes]);

  React.useEffect(() => {
    setPaginatedRows(filteredTransactions?.slice(0, perPage));
    setPage(1);
  }, [filteredTransactions, perPage]);

  // Get available tag values from the actual transactions
  const availableTags = React.useMemo(() => {
    const tagSet = new Set<string>();
    transactions.forEach(transaction => {
      transaction.tags.forEach(tag => {
        tagSet.add(tag.value);
      });
    });
    return Array.from(tagSet).sort();
  }, [transactions]);

  // Get available transaction types
  const availableTransactionTypes = React.useMemo(() => {
    const typeSet = new Set<string>();
    transactions.forEach(transaction => {
      typeSet.add(transaction.kind);
    });
    return Array.from(typeSet).sort();
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
    value: string | number | undefined
  ) => {
    if (typeof value === 'string') {
      setSelectedTransactionTypes(prev => {
        if (prev.includes(value)) {
          return prev.filter(type => type !== value);
        } else {
          return [...prev, value];
        }
      });
    }
  };

  const handleTransactionTypeRemove = (typeToRemove: string) => {
    setSelectedTransactionTypes(prev => prev.filter(type => type !== typeToRemove));
  };

  const handleClearAllFilters = () => {
    setSelectedTags([]);
    setSelectedTransactionTypes([]);
  };

  const handleTagClick = (tagValue: string) => {
    if (!selectedTags.includes(tagValue)) {
      setSelectedTags(prev => [...prev, tagValue]);
    }
  };

  const handleRuleClick = (ruleId: string) => {
    // Find all tags associated with this rule across all transactions
    const tagsForRule = new Set<string>();
    transactions.forEach(transaction => {
      transaction.tags.forEach(tag => {
        if (tag.rule === ruleId) {
          tagsForRule.add(tag.value);
        }
      });
    });

    // Add all tags for this rule to the selected tags (if not already selected)
    const newTags = Array.from(tagsForRule).filter(tag => !selectedTags.includes(tag));
    if (newTags.length > 0) {
      setSelectedTags(prev => [...prev, ...newTags]);
    }
  };

  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: activeSortIndex || undefined,
      direction: activeSortDirection || undefined,
      defaultDirection: 'asc',
    },
    onSort: (_event, index, direction) => {
      const sorted = [...filteredTransactions].sort((a, b) => {
        let aValue: Date | number = a.date;
        let bValue: Date | number = b.date;
        if (index == 5) {
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

  // Date picker handlers
  const handleStartDateChange = (value: string) => {
    setStartDate(value);
  };

  const handleEndDateChange = (value: string) => {
    setEndDate(value);
  };

  const handleSetPage = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPage: number,
    _perPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPaginatedRows(filteredTransactions?.slice(startIdx, endIdx));
    setPage(newPage);
  };

  const handlePerPageSelect = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPerPage: number,
    newPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPaginatedRows(filteredTransactions.slice(startIdx, endIdx));
    setPage(newPage);
    setPerPage(newPerPage);
  };

  const renderPagination = (
    variant: PaginationVariant,
    isCompact: boolean,
    isSticky: boolean,
    isStatic: boolean
  ) => (
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
                              variant="filled" 
                              color="orange"
                              onClose={() => handleTransactionTypeRemove(type)}
                              closeBtnAriaLabel={`Remove ${type} filter`}
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
            </Flex>
          </FlexItem>
          {(selectedTags.length > 0 || selectedTransactionTypes.length > 0) && (
            <FlexItem>
              <Button 
                variant="link" 
                onClick={handleClearAllFilters}
                isInline
              >
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
            <Th width={10} sort={getSortParams(5)}>
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
              <Td dataLabel="{columns.kind}">{t.kind}</Td>
              <Td dataLabel="{columns.tags}">
                <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  {t.tags.map((tag: ITagTransaction, idx: number) => (
                    <FlexItem key={`tag-${idx}`}>
                      <Label 
                        variant="filled" 
                        color="green"
                        onClick={() => handleTagClick(tag.value)}
                        style={{ cursor: 'pointer' }}
                        aria-label={`Filter by ${tag.value} tag`}
                      >
                        <Content component="p">{tag.value}</Content>
                      </Label>
                    </FlexItem>
                  ))}
                </Flex>
              </Td>
              <Td dataLabel="{columns.rules">
                <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  {Array.from(new Set(t.tags.map(tag => tag.rule))).map((rule: string, idx: number) => (
                    <FlexItem key={`rule-${idx}`}>
                      <Label 
                        variant="filled" 
                        color="blue"
                        onClick={() => handleRuleClick(rule)}
                        style={{ cursor: 'pointer' }}
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
              <Td noPadding={false} colSpan={5}>
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
