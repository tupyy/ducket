import * as React from 'react';
import { CodeBranchIcon, CubesIcon } from '@patternfly/react-icons';
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
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  Tooltip,
} from '@patternfly/react-core';
import { ITag } from '@app/shared/models/tag';
import { IRule } from '@app/shared/models/rule';

export interface ITagListProps {
  tags: ReadonlyArray<ITag> | [];
  showCreateTagFormCB: () => void;
}

// eslint-disable-next-line prefer-const
const TagsList: React.FunctionComponent<ITagListProps> = ({ tags, showCreateTagFormCB }) => {
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(10);
  const [paginatedRows, setPaginatedRows] = React.useState(tags.slice(0, 10));

  const handleSetPage = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPage: number,
    _perPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPaginatedRows(tags?.slice(startIdx, endIdx));
    setPage(newPage);
  };

  const handlePerPageSelect = (
    _evt: React.MouseEvent | React.KeyboardEvent | MouseEvent,
    newPerPage: number,
    newPage: number | undefined,
    startIdx: number | undefined,
    endIdx: number | undefined
  ) => {
    setPaginatedRows(tags.slice(startIdx, endIdx));
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

  const renderRow = (tag: ITag) => {
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
                <strong>{rule.name}</strong>
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
      itemCount={tags.length}
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

  const toolbarItems = (
    <React.Fragment>
      <ToolbarItem>
        <Button variant="secondary" onClick={() => showCreateTagFormCB()}>
          Create tag
        </Button>
      </ToolbarItem>
      <ToolbarItem variant="pagination" align={{ default: 'alignEnd' }}>
        {renderPagination(PaginationVariant.top, true, false, false)}
      </ToolbarItem>
    </React.Fragment>
  );

  const renderTagList = (
    <React.Fragment>
      <DataList aria-label="tag list">
        {tags.map((tag: ITag, i: number) => (
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
                                <CodeBranchIcon /> 3
                              </div>
                            </Tooltip>
                          </FlexItem>
                          <FlexItem>Created date-of-creation</FlexItem>
                        </Flex>
                      </FlexItem>
                    </Flex>
                  </DataListCell>,
                  renderRow(tag),
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
                  //     setIsDeleted(true);
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
      {tags.length == 0 ? (
        emptyState
      ) : (
        <React.Fragment>
          <Toolbar id="tags-toolbar">
            <ToolbarContent>{toolbarItems}</ToolbarContent>
          </Toolbar>
          {renderTagList}
          {renderPagination(PaginationVariant.bottom, false, false, true)}
        </React.Fragment>
      )}
    </PageSection>
  );
};

export { TagsList };
