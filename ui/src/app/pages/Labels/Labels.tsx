import * as React from 'react';
import { CubesIcon } from '@patternfly/react-icons';
import {
  Content,
  ContentVariants,
  EmptyState,
  EmptyStateBody,
  EmptyStateVariant,
  Flex,
  FlexItem,
  Label,
  PageSection,
  Pagination,
  PaginationVariant,
} from '@patternfly/react-core';
import { ILabel } from '@app/shared/models/label';
import { IRule } from '@app/shared/models/rule';
import { DataView, DataViewToolbar, useDataViewFilters } from '@patternfly/react-data-view';
import { DataViewFilters } from '@patternfly/react-data-view/dist/dynamic/DataViewFilters';
import { DataViewTextFilter } from '@patternfly/react-data-view/dist/dynamic/DataViewTextFilter';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { useTheme } from '@app/shared/contexts/ThemeContext';

export interface ILabelListProps {
  labels: Array<ILabel> | [];
}

interface RepositoryFilters {
  name: string;
  rules: string;
}

const columns = {
  name: 'Name',
  rules: 'Rules',
  action: 'action',
};

// eslint-disable-next-line prefer-const
const LabelsList: React.FunctionComponent<ILabelListProps> = ({ labels }) => {
  const { theme } = useTheme();
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(10);
  const { filters, onSetFilters, clearAllFilters } = useDataViewFilters<RepositoryFilters>({
    initialFilters: { name: '', rules: '' },
  });

  const filteredRows = React.useMemo(
    () =>
      labels.filter(
        (label) =>
          !filters.name || `${label.key}=${label.value}`.toLocaleLowerCase().includes(filters.name?.toLocaleLowerCase())
      ),
    [filters, labels]
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
    <EmptyState variant={EmptyStateVariant.full} titleText="No labels" icon={CubesIcon}>
      <EmptyStateBody>
        <Content>
          <Content component="p">Please add some labels</Content>
        </Content>
      </EmptyStateBody>
    </EmptyState>
  );

  const renderRuleCell = (label: ILabel) => {
    if (label.rules === undefined) {
      return <span></span>;
    }
    return (
      <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsMd', sm: 'spaceItemsXs' }}>
        {label.rules.map((rule: IRule, idx: number) => (
          <FlexItem key={`rule-${idx}`}>
            <Label
              variant={theme === 'dark' ? 'outline' : 'filled'}
              color="green"
              href={`/api/rules/${rule.name}`}
              style={theme === 'dark' ? { color: '#3e8635' } : {}}
            >
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

  const renderLabelList = (
    <React.Fragment>
      <Table aria-label="label list">
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
          </Tr>
        </Thead>
        <Tbody>
          {paginatedRows.map((label: ILabel, i: number) => (
            <Tr key={`label-${i}`}>
              <Td dataLabel={columns.name}>
                <Content component={ContentVariants.p}>
                  {label.key}={label.value}
                </Content>
              </Td>
              <Td dataLabel={columns.rules}>{renderRuleCell(label)}</Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </React.Fragment>
  );

  return (
    <PageSection hasBodyWrapper={false}>
      {labels.length === 0 ? (
        emptyState
      ) : (
        <DataView>
          {renderToolbar}
          {renderLabelList}
          {renderPagination(PaginationVariant.bottom, false, false, true)}
        </DataView>
      )}
    </PageSection>
  );
};

export { LabelsList };
