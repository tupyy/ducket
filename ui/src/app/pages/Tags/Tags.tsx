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
  Label,
  PageSection,
  Pagination,
  PaginationVariant,
  Tooltip,
} from '@patternfly/react-core';
import { ITag } from '@app/shared/models/tag';
import { IRule } from '@app/shared/models/rule';
import { DataView, DataViewToolbar, useDataViewFilters } from '@patternfly/react-data-view';
import { DataViewFilters } from '@patternfly/react-data-view/dist/dynamic/DataViewFilters';
import { DataViewTextFilter } from '@patternfly/react-data-view/dist/dynamic/DataViewTextFilter';

export interface ITagListProps {
  tags: Array<ITag> | [];
  showCreateTagFormCB: () => void;
  deleteTagCB: (name: string) => void;
}

interface RepositoryFilters {
  name: string;
  rules: string;
}

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
      return (
        <DataListCell key="rules">
          <span>-</span>
        </DataListCell>
      );
    }
    return (
      <DataListCell key="rules">
        <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsMd', sm: 'spaceItemsXs' }}>
          {tag.rules.map((rule: IRule, idx: number) => (
            <FlexItem key={`rule-${idx}`}>
              <Label variant="filled" color="green" href={`/api/rules/${rule.name}`}>
                <Content component="p">
                  <strong>{rule.name}</strong>
                </Content>
              </Label>
            </FlexItem>
          ))}
        </Flex>
      </DataListCell>
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
        <Button onClick={showCreateTagFormCB} variant="secondary">
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
      <DataList aria-label="tag list">
        {paginatedRows.map((tag: ITag, i: number) => (
          <DataListItem key={`tag-${i}`}>
            <DataListItemRow>
              <DataListItemCells
                dataListCells={[
                  <DataListCell key="tag name">
                    <Flex direction={{ default: 'column' }}>
                      <FlexItem>
                        <Content component={ContentVariants.p}>
                          <strong>{tag.value}</strong>
                        </Content>
                      </FlexItem>
                      <FlexItem>
                        <Flex flexWrap={{ default: 'wrap' }} key="info">
                          <FlexItem>
                            <Tooltip content={<div>Number of transactions on which this tag is applied</div>}>
                              <div>
                                <CodeBranchIcon /> {tag.transactions}
                              </div>
                            </Tooltip>
                          </FlexItem>
                          <FlexItem>
                            <OutlinedCalendarIcon />
                            {` ` +
                              tag.created_at.toLocaleDateString('fr-FR', {
                                weekday: 'long',
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                              })}
                          </FlexItem>
                        </Flex>
                      </FlexItem>
                    </Flex>
                  </DataListCell>,
                  renderRuleCell(tag),
                ]}
              />
              <DataListAction
                aria-labelledby="single-action-item1 single-action-action1"
                id="single-action-action1"
                aria-label="Actions"
              >
                <Button
                  onClick={() => {
                    if (confirm('Are you sure?')) {
                      deleteTagCB(tag.value);
                    }
                  }}
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
