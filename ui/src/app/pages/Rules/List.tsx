import * as React from 'react';
import {
  EuiInMemoryTable,
  EuiBasicTableColumn,
  EuiTableActionsColumnType,
  EuiFlexGroup,
  EuiFlexItem,
  EuiBadge,
  EuiText,
  EuiEmptyPrompt,
  EuiFieldText,
  EuiSpacer,
  EuiPanel,
  EuiFormRow,
  EuiButton,
  EuiButtonIcon,
  EuiLoadingSpinner,
  EuiToolTip,
} from '@elastic/eui';
import { IRule } from '@app/shared/models/rule';
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

interface RuleFilters {
  name: string;
  pattern: string;
}

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
  const [filters, setFilters] = React.useState<RuleFilters>({ name: '', pattern: '' });

  const filteredRules = React.useMemo(() => {
    return rules.filter((rule) => {
      const nameMatch = !filters.name || rule.name.toLowerCase().includes(filters.name.toLowerCase());
      const patternMatch = !filters.pattern || rule.pattern.toLowerCase().includes(filters.pattern.toLowerCase());
      return nameMatch && patternMatch;
    });
  }, [rules, filters]);

  const renderLabelsCell = (labels: ILabel[] | undefined) => {
    if (!labels || labels.length === 0) {
      return <EuiText size="s" color="subdued">No labels</EuiText>;
    }

    return (
      <EuiFlexGroup gutterSize="xs" wrap>
        {labels.map((label: ILabel, idx: number) => (
          <EuiFlexItem grow={false} key={`label-${idx}`}>
            <EuiBadge color="hollow">
              {label.key}={label.value}
            </EuiBadge>
          </EuiFlexItem>
        ))}
      </EuiFlexGroup>
    );
  };

  const actions: EuiTableActionsColumnType<IRule>['actions'] = [
    {
      name: 'Edit',
      description: 'Edit this rule',
      icon: 'pencil',
      type: 'icon',
      onClick: (rule) => showEditRuleFormCB(rule),
    },
    {
      name: 'Sync',
      description: 'Sync this rule',
      icon: 'refresh',
      type: 'icon',
      onClick: (rule) => onSyncRule(rule.name),
      enabled: () => !syncing,
    },
    {
      name: 'Delete',
      description: 'Delete this rule',
      icon: 'trash',
      type: 'icon',
      color: 'danger',
      onClick: (rule) => onDeleteRule(rule.name),
      enabled: () => !syncing,
    },
  ];

  const columns: Array<EuiBasicTableColumn<IRule>> = [
    {
      field: 'name',
      name: 'Name',
      sortable: true,
      render: (name: string) => (
        <EuiText size="s" style={{ fontWeight: 'bold' }}>
          {name}
        </EuiText>
      ),
      width: '200px',
    },
    {
      field: 'pattern',
      name: 'Pattern',
      sortable: true,
      render: (pattern: string) => (
        <EuiText size="s" style={{ fontFamily: 'monospace', maxWidth: '300px' }}>
          {pattern}
        </EuiText>
      ),
    },
    {
      field: 'labels',
      name: 'Labels',
      render: (labels: ILabel[] | undefined) => renderLabelsCell(labels),
    },
    {
      field: 'transactionCount',
      name: 'Transactions',
      sortable: true,
      render: (count: number) => (
        <EuiBadge color={count > 0 ? 'success' : 'default'}>
          {count || 0}
        </EuiBadge>
      ),
      width: '120px',
    },
  ];

  const actionsColumn: EuiTableActionsColumnType<IRule> = {
    actions,
    width: '120px',
  };

  const allColumns = [...columns, actionsColumn];

  const renderToolbar = () => (
    <EuiPanel paddingSize="s">
      <EuiFlexGroup gutterSize="s" alignItems="flexEnd">
        <EuiFlexItem style={{ minWidth: '200px' }}>
          <EuiFormRow label="Filter by Name">
            <EuiFieldText
              placeholder="Enter rule name to filter..."
              value={filters.name}
              onChange={(e) => setFilters(prev => ({ ...prev, name: e.target.value }))}
              compressed
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem style={{ minWidth: '200px' }}>
          <EuiFormRow label="Filter by Pattern">
            <EuiFieldText
              placeholder="Enter pattern to filter..."
              value={filters.pattern}
              onChange={(e) => setFilters(prev => ({ ...prev, pattern: e.target.value }))}
              compressed
            />
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            <EuiButton size="s" onClick={() => setFilters({ name: '', pattern: '' })}>
              Clear Filters
            </EuiButton>
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            <EuiButton size="s" fill onClick={showCreateRuleFormCB}>
              Create Rule
            </EuiButton>
          </EuiFormRow>
        </EuiFlexItem>

        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            <EuiButton 
              size="s" 
              onClick={onSyncAllRules}
              isLoading={syncingAll}
              isDisabled={syncingAll || syncing}
            >
              Sync All Rules
            </EuiButton>
          </EuiFormRow>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPanel>
  );

  const pagination = {
    pageIndex: 0,
    pageSize: 10,
    totalItemCount: filteredRules.length,
    showPerPageOptions: true,
    pageSizeOptions: [10, 20, 50],
  };

  const sorting = {
    sort: {
      field: 'name' as keyof IRule,
      direction: 'asc' as const,
    },
  };

  if (rules.length === 0) {
    return (
      <EuiEmptyPrompt
        icon="gear"
        title={<h2>No rules</h2>}
        body={<p>Create your first rule to automatically label transactions</p>}
        actions={
          <EuiButton fill onClick={showCreateRuleFormCB}>
            Create Rule
          </EuiButton>
        }
      />
    );
  }

  return (
    <>
      {renderToolbar()}
      <EuiSpacer size="m" />
      
      <EuiInMemoryTable
        items={filteredRules}
        columns={allColumns}
        pagination={pagination}
        sorting={sorting}
        loading={syncing}
        message={filteredRules.length === 0 ? "No rules match your filters" : undefined}
      />
    </>
  );
};

export { RulesList };