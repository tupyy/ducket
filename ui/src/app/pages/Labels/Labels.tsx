import * as React from 'react';
import {
  EuiInMemoryTable,
  EuiBasicTableColumn,
  EuiFlexGroup,
  EuiFlexItem,
  EuiBadge,
  EuiText,
  EuiEmptyPrompt,
  EuiFieldText,
  EuiSpacer,
  EuiPanel,
  EuiFormRow,
} from '@elastic/eui';
import { ILabel } from '@app/shared/models/label';
import { IRule } from '@app/shared/models/rule';
import { useTheme } from '@app/shared/contexts/ThemeContext';

export interface ILabelListProps {
  labels: Array<ILabel> | [];
}

interface LabelFilters {
  key: string;
  value: string;
}

const LabelsList: React.FunctionComponent<ILabelListProps> = ({ labels }) => {
  const { theme } = useTheme();
  const [filters, setFilters] = React.useState<LabelFilters>({ key: '', value: '' });

  const filteredLabels = React.useMemo(() => {
    return labels.filter((label) => {
      const keyMatch = !filters.key || label.key.toLowerCase().includes(filters.key.toLowerCase());
      const valueMatch = !filters.value || label.value.toLowerCase().includes(filters.value.toLowerCase());
      return keyMatch && valueMatch;
    });
  }, [labels, filters]);

  const renderRuleCell = (rules: IRule[] | undefined) => {
    if (!rules || rules.length === 0) {
      return <EuiText size="s" color="subdued">No rules</EuiText>;
    }

    return (
      <EuiFlexGroup gutterSize="xs" wrap>
        {rules.map((rule: IRule, idx: number) => (
          <EuiFlexItem grow={false} key={`rule-${idx}`}>
            <EuiBadge 
              color="success"
              href={`/api/rules/${rule.name}`}
            >
              {rule.name}
            </EuiBadge>
          </EuiFlexItem>
        ))}
      </EuiFlexGroup>
    );
  };

  const columns: Array<EuiBasicTableColumn<ILabel>> = [
    {
      field: 'key',
      name: 'Key',
      sortable: true,
      render: (key: string) => (
        <EuiText size="s" style={{ fontWeight: 'bold' }}>
          {key}
        </EuiText>
      ),
      width: '200px',
    },
    {
      field: 'value',
      name: 'Value',
      sortable: true,
      render: (value: string) => (
        <EuiText size="s">
          {value}
        </EuiText>
      ),
      width: '200px',
    },
    {
      field: 'rules',
      name: 'Rules',
      render: (rules: IRule[] | undefined) => renderRuleCell(rules),
    },
  ];

  const renderFilters = () => (
    <EuiPanel paddingSize="s">
      <EuiFlexGroup gutterSize="s" alignItems="flexEnd">
        <EuiFlexItem style={{ minWidth: '200px' }}>
          <EuiFormRow label="Filter by Key">
            <EuiFieldText
              placeholder="Enter key to filter..."
              value={filters.key}
              onChange={(e) => setFilters(prev => ({ ...prev, key: e.target.value }))}
              compressed
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem style={{ minWidth: '200px' }}>
          <EuiFormRow label="Filter by Value">
            <EuiFieldText
              placeholder="Enter value to filter..."
              value={filters.value}
              onChange={(e) => setFilters(prev => ({ ...prev, value: e.target.value }))}
              compressed
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            <EuiBadge 
              color="hollow"
              onClick={() => setFilters({ key: '', value: '' })}
              onClickAriaLabel="Clear all filters"
              iconType="cross"
              iconSide="right"
            >
              Clear Filters
            </EuiBadge>
          </EuiFormRow>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPanel>
  );

  const pagination = {
    pageIndex: 0,
    pageSize: 20,
    totalItemCount: filteredLabels.length,
    showPerPageOptions: true,
    pageSizeOptions: [10, 20, 50, 100],
  };

  const sorting = {
    sort: {
      field: 'key' as keyof ILabel,
      direction: 'asc' as const,
    },
  };

  if (labels.length === 0) {
    return (
      <EuiEmptyPrompt
        icon="tag"
        title={<h2>No labels</h2>}
        body={<p>Please add some labels</p>}
      />
    );
  }

  return (
    <>
      {renderFilters()}
      <EuiSpacer size="m" />
      
      <EuiInMemoryTable
        items={filteredLabels}
        columns={columns}
        pagination={pagination}
        sorting={sorting}
        message={filteredLabels.length === 0 ? "No labels match your filters" : undefined}
      />
    </>
  );
};

export { LabelsList };