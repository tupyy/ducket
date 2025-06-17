import * as React from 'react';
import { CodeBranchIcon, CubesIcon, OutlinedCalendarIcon } from '@patternfly/react-icons';
import {
  Button,
  Content,
  ContentVariants,
  DataList,
  DataListAction,
  DataListCell,
  DataListItem,
  DataListItemCells,
  DataListItemRow,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateVariant,
  Flex,
  FlexItem,
  PageSection,
  Pagination,
  PaginationVariant,
  Tooltip,
} from '@patternfly/react-core';
import { IRule } from '@app/shared/models/rule';
import { DataView, DataViewToolbar, useDataViewFilters } from '@patternfly/react-data-view';
import { DataViewFilters } from '@patternfly/react-data-view/dist/dynamic/DataViewFilters';
import { DataViewTextFilter } from '@patternfly/react-data-view/dist/dynamic/DataViewTextFilter';

export interface IRuleListProps {
  rules: Array<IRule> | [];
  showCreateRuleFormCB: () => void;
  // deleteTagCB: (name: string) => void;
}

interface RepositoryFilters {
  name: string;
  pattern: string;
}

// eslint-disable-next-line prefer-const
const RulesList: React.FunctionComponent<IRuleListProps> = ({ rules, showCreateRuleFormCB }) => {
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
          <DataViewTextFilter filterId="pattern" title="Pattern" placeholder="Filter by pattern" />
        </DataViewFilters>
      }
      pagination={renderPagination(PaginationVariant.top, true, false, false)}
    />
  );

  const renderTagList = (
    <React.Fragment>
      <DataList aria-label="tag list">
        {paginatedRows.map((rule: IRule, i: number) => (
          <DataListItem key={`tag-${i}`}>
            <DataListItemRow>
              <DataListItemCells
                dataListCells={[
                  <DataListCell key="tag name">
                    <Flex direction={{ default: 'column' }}>
                      <FlexItem>
                        <Content component={ContentVariants.p}>
                          <strong>{rule.name}</strong>
                        </Content>
                      </FlexItem>
                      <FlexItem>
                        <Flex flexWrap={{ default: 'wrap' }} key="info">
                          <FlexItem>
                            <Tooltip content={<div>Number of transactions on which this tag is applied</div>}>
                              <div>
                                <CodeBranchIcon /> {rule.transactions}
                              </div>
                            </Tooltip>
                          </FlexItem>
                          <FlexItem>
                            <OutlinedCalendarIcon />
                          </FlexItem>
                        </Flex>
                      </FlexItem>
                    </Flex>
                  </DataListCell>,
                  <DataListCell key="pattern">
                    <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsMd', sm: 'spaceItemsXs' }}>
                      <Content component="p">
                        /<strong>{rule.pattern}</strong>/
                      </Content>
                    </Flex>
                  </DataListCell>,
                ]}
              />
              <DataListAction
                aria-labelledby="single-action-item1 single-action-action1"
                id="single-action-action1"
                aria-label="Actions"
              >
                <Button
                  // onClick={() => {
                  //   if (confirm('Are you sure?')) {
                  //     deleteTagCB(rule.name);
                  //   }
                  // }}
                  variant="secondary"
                  key="delete-action"
                  size="sm"
                >
                  Delete
                </Button>
              </DataListAction>
            </DataListItemRow>
          </DataListItem>
        ))}
      </DataList>
    </React.Fragment>
  );

  return (
    <PageSection hasBodyWrapper={false}>
      {rules.length == 0 ? (
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

export { RulesList };
