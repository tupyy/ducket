import * as React from 'react';
import {CubesIcon} from '@patternfly/react-icons';
import {
  Button,
  Content,
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
import { IRule } from '@app/shared/models/rule';
import { DataView, DataViewToolbar, useDataViewFilters } from '@patternfly/react-data-view';
import { DataViewFilters } from '@patternfly/react-data-view/dist/dynamic/DataViewFilters';
import { DataViewTextFilter } from '@patternfly/react-data-view/dist/dynamic/DataViewTextFilter';
import { Table, Tbody, Td, Th, Thead, Tr, ActionsColumn, IAction } from '@patternfly/react-table';
import { ITag } from '@app/shared/models/tag';

export interface IRuleListProps {
  rules: Array<IRule> | [];
  showCreateRuleFormCB: () => void;
  showEditRuleFormCB: (rule: IRule) => void;
  onSyncRule: (ruleName: string) => void;
  onDeleteRule: (ruleName: string) => void;
}

interface RepositoryFilters {
  name: string;
  pattern: string;
}

const columns = {
  name: 'Name',
  pattern: 'Pattern',
  tags: 'Tags',
  transactions: 'Transactions',
  createdAt: 'Created at',
  action: 'action',
};

// eslint-disable-next-line prefer-const
const RulesList: React.FunctionComponent<IRuleListProps> = ({ rules, showCreateRuleFormCB, showEditRuleFormCB, onSyncRule, onDeleteRule }) => {
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(10);
  const { filters, onSetFilters, clearAllFilters } = useDataViewFilters<RepositoryFilters>({
    initialFilters: { name: '', pattern: '' },
  });

  const filteredRows = React.useMemo(
    () =>
      rules.filter(
        (rule) => !filters.name || rule.name.toLocaleLowerCase().includes(filters.name?.toLocaleLowerCase())
      ),
    [filters, rules]
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
    <EmptyState variant={EmptyStateVariant.full} titleText="No rules" icon={CubesIcon}>
      <EmptyStateBody>
        <Content>
          <Content component="p">Please add some rules</Content>
        </Content>
      </EmptyStateBody>
      <EmptyStateFooter>
        <Button variant="primary" onClick={showCreateRuleFormCB}>
          Add rule
        </Button>
      </EmptyStateFooter>
    </EmptyState>
  );

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
        <Button onClick={showCreateRuleFormCB} variant="control">
          Create rule
        </Button>
      }
      clearAllFilters={clearAllFilters}
      filters={
        <DataViewFilters onChange={(_e, values) => onSetFilters(values)} values={filters}>
          <DataViewTextFilter filterId="name" title="Name" placeholder="Filter by name" />
        </DataViewFilters>
      }
      pagination={renderPagination(PaginationVariant.top, true, false, false)}
    />
  );

  const renderList = (
    <React.Fragment>
      <Table aria-label="rule-list">
        <Thead>
          <Tr>
            <Th>
              <Content component="p">
                <strong>{columns.name}</strong>
              </Content>
            </Th>
            <Th>
              <Content component="p">
                <strong>{columns.pattern}</strong>
              </Content>
            </Th>
            <Th>
              <Content component="p">
                <strong>{columns.tags}</strong>
              </Content>
            </Th>
            <Th width={20}>
              <Content component="p">
                <strong>{columns.transactions}</strong>
              </Content>
            </Th>
            <Th width={20}>
              <Content component="p">
                <strong>{columns.createdAt}</strong>
              </Content>
            </Th>
            <Th />
          </Tr>
        </Thead>
        <Tbody>
          {paginatedRows.map((rule: IRule, i: number) => (
            <Tr key={`row-${i}`}>
              <Td dataLabel="{columns.name}">{rule.name}</Td>
              <Td dataLabel="{columns.pattern">/{rule.pattern}/</Td>
              <Td dataLabel="{columns.tags}">
                <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  {rule.tags.map((tag: ITag, idx: number) => (
                    <FlexItem key={`tag-${idx}`}>
                      <Label variant="filled" color="green" href={`/api/tags/${tag.value}`}>
                        <Content component="p">{tag.value}</Content>
                      </Label>
                    </FlexItem>
                  ))}
                </Flex>
              </Td>
              <Td dataLabel="{columns.transactions">{rule.transactions}</Td>
              <Td dataLabel={columns.createdAt}>
                {rule.created_at.toLocaleDateString('fr-FR', {
                  weekday: 'long',
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </Td>
              <Td isActionCell dataLabel={columns.action}>
                <ActionsColumn
                  items={[
                    {
                      title: 'Sync',
                      onClick: () => onSyncRule(rule.name)
                    },
                    {
                      title: 'Edit',
                      onClick: () => showEditRuleFormCB(rule)
                    },
                    {
                      isSeparator: true
                    },
                    {
                      title: 'Delete',
                      onClick: () => {
                        if (confirm('Are you sure you want to delete this rule?')) {
                          onDeleteRule(rule.name);
                        }
                      }
                    }
                  ]}
                />
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </React.Fragment>
  );

  return (
    <PageSection hasBodyWrapper={false}>
      {rules.length == 0 ? (
        emptyState
      ) : (
        <DataView>
          {renderToolbar}
          {renderList}
          {renderPagination(PaginationVariant.bottom, false, false, true)}
        </DataView>
      )}
    </PageSection>
  );
};

export { RulesList };
