import * as React from 'react';
import {CubesIcon} from '@patternfly/react-icons';
import {
  Button,
  Content,
  ContentVariants,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateVariant,
  Flex,
  FlexItem,
  Label,
  PageSection,
  Pagination,
  PaginationVariant,
} from '@patternfly/react-core';
import { ITag } from '@app/shared/models/tag';
import { IRule } from '@app/shared/models/rule';
import { DataView, DataViewToolbar, useDataViewFilters } from '@patternfly/react-data-view';
import { DataViewFilters } from '@patternfly/react-data-view/dist/dynamic/DataViewFilters';
import { DataViewTextFilter } from '@patternfly/react-data-view/dist/dynamic/DataViewTextFilter';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';

export interface ITagListProps {
  tags: Array<ITag> | [];
  showCreateTagFormCB: () => void;
  deleteTagCB: (name: string) => void;
}

interface RepositoryFilters {
  name: string;
  rules: string;
}

const columns = {
  name: 'Name',
  rules: 'Rules',
  createdAt: 'Created at',
  transactions: 'Transactions',
  action: 'action',
};

// eslint-disable-next-line prefer-const
const TagsList: React.FunctionComponent<ITagListProps> = ({ tags, showCreateTagFormCB, deleteTagCB }) => {
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(10);
  const { filters, onSetFilters, clearAllFilters } = useDataViewFilters<RepositoryFilters>({
    initialFilters: { name: '', rules: '' },
  });

  const filteredRows = React.useMemo(
    () =>
      tags.filter((tag) => !filters.name || tag.value.toLocaleLowerCase().includes(filters.name?.toLocaleLowerCase())),
    [filters, tags]
  );
  const [paginatedRows, setPaginatedRows] = React.useState(filteredRows.slice(0, 10));

  React.useEffect(() => {
    setPaginatedRows(filteredRows?.slice(0, 10));
    setPage(1);
  }, [filteredRows]);

  const handleSetPage = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPage: number,
    _perPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPaginatedRows(filteredRows?.slice(startIdx, endIdx));
    setPage(newPage);
  };

  const handlePerPageSelect = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPerPage: number,
    newPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPaginatedRows(filteredRows.slice(startIdx, endIdx));
    setPage(newPage);
    setPerPage(newPerPage);
  };

  const emptyState = (
    <EmptyState variant={EmptyStateVariant.full} titleText="No tags" icon={CubesIcon}>
      <EmptyStateBody>
        <Content>
          <Content component="p">Please add some tags</Content>
        </Content>
      </EmptyStateBody>
      <EmptyStateFooter>
        <Button variant="primary">Add tag</Button>
      </EmptyStateFooter>
    </EmptyState>
  );

  const renderRuleCell = (tag: ITag) => {
    if (tag.rules === undefined) {
      return <span></span>;
    }
    return (
      <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsMd', sm: 'spaceItemsXs' }}>
        {tag.rules.map((rule: IRule, idx: number) => (
          <FlexItem key={`rule-${idx}`}>
            <Label variant="filled" color="green" href={`/api/rules/${rule.name}`}>
              <Content component="p">{rule.name}</Content>
            </Label>
          </FlexItem>
        ))}
      </Flex>
    );
  };

  const renderPagination = (variant: PaginationVariant, isCompact: boolean, isSticky: boolean, isStatic: boolean) => (
    <Pagination
      id={`datalist-${variant}-pagination`}
      variant={variant}
      itemCount={filteredRows.length}
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
      bulkSelect={
        <Button onClick={showCreateTagFormCB} variant="control">
          Create tag
        </Button>
      }
      clearAllFilters={clearAllFilters}
      filters={
        <DataViewFilters onChange={(_e, values) => onSetFilters(values)} values={filters}>
          <DataViewTextFilter filterId="name" title="Name" placeholder="Filter by name" />
          <DataViewTextFilter filterId="rules" title="Rules" placeholder="Filter by rules" />
        </DataViewFilters>
      }
      pagination={renderPagination(PaginationVariant.top, true, false, false)}
    />
  );

  const renderTagList = (
    <React.Fragment>
      <Table aria-label="tag list">
        <Thead>
          <Tr>
            <Th>
              <Content component={ContentVariants.p}>
                <strong>{columns.name}</strong>
              </Content>
            </Th>
            <Th>
              <Content component={ContentVariants.p}>
                <strong>{columns.rules}</strong>
              </Content>
            </Th>
            <Th width={10}>
              <Content component={ContentVariants.p}>
                <strong>{columns.transactions}</strong>
              </Content>
            </Th>
            <Th width={10}>
              <Content component={ContentVariants.p}>
                <strong>{columns.createdAt}</strong>
              </Content>
            </Th>
            <Th screenReaderText="action" />
          </Tr>
        </Thead>
        <Tbody>
          {paginatedRows.map((tag: ITag, i: number) => (
            <Tr key={`tag-${i}`}>
              <Td dataLabel={columns.name}>
                <Content component={ContentVariants.p}>{tag.value}</Content>
              </Td>
              <Td dataLabel={columns.rules}>{renderRuleCell(tag)}</Td>
              <Td dataLabel={columns.transactions}>{tag.transactions}</Td>
              <Td dataLabel={columns.createdAt}>
                {tag.created_at.toLocaleDateString('fr-FR', {
                  weekday: 'long',
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </Td>
              <Td isActionCell dataLabel={columns.action}>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => {
                    if (confirm('Are you sure?')) {
                      deleteTagCB(tag.value);
                    }
                  }}
                >
                  Delete
                </Button>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </React.Fragment>
  );

  return (
    <PageSection hasBodyWrapper={false}>
      {tags.length == 0 ? (
        emptyState
      ) : (
        <DataView>
          {renderToolbar}
          {renderTagList}
          {renderPagination(PaginationVariant.bottom, false, false, true)}
        </DataView>
      )}
    </PageSection>
  );
};

export { TagsList };
