import * as React from 'react';
import { CubesIcon } from '@patternfly/react-icons';
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
  TextInput,
} from '@patternfly/react-core';
import { IRule } from '@app/shared/models/rule';
import { DataView, DataViewToolbar } from '@patternfly/react-data-view';
import { Table, Tbody, Td, Th, Thead, Tr, ActionsColumn, IAction } from '@patternfly/react-table';
import { ILabel } from '@app/shared/models/label';
import { useTheme } from '@app/shared/contexts/ThemeContext';

export interface IRuleListProps {
  rules: Array<IRule> | [];
  showCreateRuleFormCB: () => void;
  showEditRuleFormCB: (rule: IRule) => void;
  onSyncRule: (ruleName: string) => void;
  onSyncAllRules: () => void;
  onDeleteRule: (ruleName: string) => void;
  syncing?: boolean;
  syncingAll?: boolean;
}

interface RepositoryFilters {
  name: string;
  pattern: string;
}

const columns = {
  name: 'Name',
  pattern: 'Pattern',
  labels: 'Labels',
  transactions: 'Transactions',
  action: 'action',
};

// eslint-disable-next-line prefer-const
const RulesList: React.FunctionComponent<IRuleListProps> = ({
  rules,
  showCreateRuleFormCB,
  showEditRuleFormCB,
  onSyncRule,
  onSyncAllRules,
  onDeleteRule,
  syncing = false,
  syncingAll = false,
}) => {
  const { theme } = useTheme();
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(10);
  const [filters, setFilters] = React.useState<RepositoryFilters>({ name: '', pattern: '' });

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

  const handleFilterChange = (
    _evt: React.FormEvent<HTMLInputElement> | React.ChangeEvent<HTMLTextAreaElement>,
    value: string
  ) => {
    setFilters((prev) => ({ ...prev, name: value }));
  };

  const clearAllFilters = () => {
    setFilters({ name: '', pattern: '' });
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

  const renderButtons = (
    <div style={{ padding: '1rem 0', marginBottom: '1rem' }}>
      <div style={{ display: 'flex', gap: '0.5rem' }}>
        <Button onClick={showCreateRuleFormCB} variant="control">
          Create rule
        </Button>
        <Button onClick={onSyncAllRules} variant="control" isLoading={syncingAll} isDisabled={syncingAll}>
          {syncingAll ? 'Syncing all rules...' : 'Sync all rules'}
        </Button>
      </div>
    </div>
  );

  const renderToolbar = (
    <DataViewToolbar
      filters={
        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
          <TextInput
            type="text"
            placeholder="Filter by name"
            value={filters.name}
            onChange={handleFilterChange}
            aria-label="Name filter"
            style={{ width: '300px' }}
          />
          {(filters.name) && (
            <Button
              variant="link"
              onClick={clearAllFilters}
              size="sm"
            >
              Clear filter
            </Button>
          )}
        </div>
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
                <strong>{columns.labels}</strong>
              </Content>
            </Th>
            <Th width={20}>
              <Content component="p">
                <strong>{columns.transactions}</strong>
              </Content>
            </Th>
            <Th screenReaderText="actions" />
          </Tr>
        </Thead>
        <Tbody>
          {paginatedRows.map((rule: IRule, i: number) => (
            <Tr key={`row-${i}`}>
              <Td dataLabel={columns.name}>{rule.name}</Td>
              <Td dataLabel={columns.pattern}>/{rule.pattern}/</Td>
              <Td dataLabel={columns.labels}>
                <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                  {rule.labels.map((label: ILabel, idx: number) => (
                    <FlexItem key={`label-${idx}`}>
                      <Label
                        variant={theme === 'dark' ? 'outline' : 'filled'}
                        color="green"
                        href={`/api/labels/${label.key}/${label.value}`}
                        style={theme === 'dark' ? { color: '#3e8635' } : {}}
                      >
                        <Content component="p">
                          {label.key}={label.value}
                        </Content>
                      </Label>
                    </FlexItem>
                  ))}
                </Flex>
              </Td>
              <Td dataLabel={columns.transactions}>{rule.transactions}</Td>
              <Td isActionCell dataLabel={columns.action}>
                <ActionsColumn
                  items={[
                    {
                      title: syncing ? 'Syncing...' : 'Sync',
                      onClick: () => onSyncRule(rule.name),
                      isDisabled: syncing,
                    },
                    {
                      title: 'Edit',
                      onClick: () => showEditRuleFormCB(rule),
                    },
                    {
                      isSeparator: true,
                    },
                    {
                      title: 'Delete',
                      onClick: () => {
                        if (confirm('Are you sure you want to delete this rule?')) {
                          onDeleteRule(rule.name);
                        }
                      },
                    },
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
      {rules.length === 0 ? (
        emptyState
      ) : (
        <React.Fragment>
          {renderButtons}
          <DataView>
            {renderToolbar}
            {renderList}
            {renderPagination(PaginationVariant.bottom, false, false, true)}
          </DataView>
        </React.Fragment>
      )}
    </PageSection>
  );
};

export { RulesList };
