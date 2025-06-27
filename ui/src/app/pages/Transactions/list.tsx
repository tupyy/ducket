import * as React from 'react';
import {
  Content,
  Flex,
  FlexItem,
  Label,
  PageSection,
  Pagination,
  PaginationVariant,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
} from '@patternfly/react-core';
import { DataView, DataViewToolbar } from '@patternfly/react-data-view';
import { ExpandableRowContent, Table, Tbody, Td, Th, Thead, ThProps, Tr } from '@patternfly/react-table';
import { ITagTransaction, ITransaction } from '@app/shared/models/transaction';
import { DateRangePicker } from '@app/shared/components/time-picker';

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

  React.useEffect(() => {
    setPaginatedRows(sortedTransactions?.slice(0, perPage));
    setPage(1);
  }, [sortedTransactions, perPage]);

  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: activeSortIndex || undefined,
      direction: activeSortDirection || undefined,
      defaultDirection: 'asc', // starting sort direction when first sorting a column. Defaults to 'asc'
    },
    onSort: (_event, index, direction) => {
      setSortedTransactions(
        sortedTransactions.sort((a, b) => {
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
            return (aValue as Date).getDate() - (bValue as Date).getDate();
          }
          return (bValue as Date).getDate() - (aValue as Date).getDate();
        })
      );
      setActiveSortIndex(index);
      setActiveSortDirection(direction);
      setPaginatedRows(sortedTransactions.slice(0, perPage));
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
    setPaginatedRows(sortedTransactions?.slice(startIdx, endIdx));
    setPage(newPage);
  };

  const handlePerPageSelect = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPerPage: number,
    newPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPaginatedRows(sortedTransactions.slice(startIdx, endIdx));
    setPage(newPage);
    setPerPage(newPerPage);
  };

  const renderPagination = (
    transactions: Array<ITransaction>,
    variant: PaginationVariant,
    isCompact: boolean,
    isSticky: boolean,
    isStatic: boolean
  ) => (
    <Pagination
      id={`datalist-${variant}-pagination`}
      variant={variant}
      itemCount={transactions.length}
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

  const renderList = (
    <React.Fragment>
      <Table aria-label="rule-list">
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
                      <Label variant="filled" color="green" href={`${tag.href}`}>
                        <Content component="p">{tag.value}</Content>
                      </Label>
                    </FlexItem>
                  ))}
                </Flex>
              </Td>
              <Td dataLabel="{columns.rules">
                <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  {t.tags.map((tag: ITagTransaction, idx: number) => (
                    <FlexItem key={`rule-${idx}`}>
                      <Label variant="filled" color="blue" href={`/api/rules/${tag.rule}`}>
                        <Content component="p">{tag.rule}</Content>
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
        {renderList}
        {renderPagination(sortedTransactions, PaginationVariant.bottom, false, false, true)}
      </DataView>
    </PageSection>
  );
};

export { TransactionList };
