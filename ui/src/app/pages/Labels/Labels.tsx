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
  key: string;
  value: string;
}

const columns = {
  Label: 'Label',
  Rule: 'Rule',
};

// eslint-disable-next-line prefer-const
const LabelsList: React.FunctionComponent<ILabelListProps> = ({ labels }) => {
  const { theme } = useTheme();
  const [page, setPage] = React.useState<number | undefined>(1);
  const [perPage, setPerPage] = React.useState<number>(20);
  const { filters, onSetFilters, clearAllFilters } = useDataViewFilters<RepositoryFilters>({
    initialFilters: { key: '', value: '' },
  });

  const filteredRows = React.useMemo(
    () =>
      labels.filter((label) => {
        let keyMatch = true;
        let valueMatch = true;
        
        // Filter by key if key filter is provided
        if (filters.key) {
          keyMatch = label.key.toLowerCase().includes(filters.key.toLowerCase());
        }
        
        // Filter by value if value filter is provided
        if (filters.value) {
          valueMatch = label.value.toLowerCase().includes(filters.value.toLowerCase());
        }
        
        // Both filters must match (AND logic)
        return keyMatch && valueMatch;
      }),
    [filters, labels]
  );
  const [paginatedRows, setPaginatedRows] = React.useState(filteredRows.slice(0, 20));

  React.useEffect(() => {
    setPaginatedRows(filteredRows?.slice(0, perPage));
    setPage(1);
  }, [filteredRows, perPage]);

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
          <DataViewTextFilter filterId="key" title="Key" placeholder="Filter by key" />
          <DataViewTextFilter filterId="value" title="Value" placeholder="Filter by value" />
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
                <strong>{columns.Label}</strong>
              </Content>
            </Th>
            <Th>
              <Content component={ContentVariants.p}>
                <strong>{columns.Rule}</strong>
              </Content>
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          {paginatedRows.map((label: ILabel, i: number) => (
            <Tr key={`label-${i}`}>
              <Td dataLabel={columns.Label}>
                <Content component={ContentVariants.p}>
                  {label.key}={label.value}
                </Content>
              </Td>
              <Td dataLabel={columns.Rule}>{renderRuleCell(label)}</Td>
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
