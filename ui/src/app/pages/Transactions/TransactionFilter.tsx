import * as React from 'react';
import {
  Button,
  Checkbox,
  DatePicker,
  Dropdown,
  Label,
  LabelGroup,
  MenuToggle,
  type MenuToggleElement,
  Pagination,
  SearchInput,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
} from '@patternfly/react-core';
import { FilterIcon } from '@patternfly/react-icons';

interface TransactionFilterProps {
  filter: string;
  onFilterChange: (filter: string) => void;
  total: number;
  page: number;
  perPage: number;
  onSetPage: (event: React.MouseEvent | React.KeyboardEvent | MouseEvent, page: number) => void;
  onPerPageSelect: (event: React.MouseEvent | React.KeyboardEvent | MouseEvent, perPage: number) => void;
}

interface FilterState {
  kind: string[];
  amountTier: string[];
  dateFrom: string;
  dateTo: string;
}

interface AppliedFilter {
  category: string;
  label: string;
  key: string;
}

const amountTiers = [
  { label: '€0 - €50', expr: 'amount <= 50' },
  { label: '€50 - €200', expr: 'amount > 50 and amount <= 200' },
  { label: '€200 - €500', expr: 'amount > 200 and amount <= 500' },
  { label: '€500+', expr: 'amount > 500' },
];

const emptyFilters: FilterState = { kind: [], amountTier: [], dateFrom: '', dateTo: '' };

function buildFilterExpr(state: FilterState): string {
  const parts: string[] = [];

  if (state.kind.length === 1) {
    parts.push(`kind = '${state.kind[0]}'`);
  }

  if (state.amountTier.length > 0) {
    const tierExprs = state.amountTier.map((t) => {
      const tier = amountTiers.find((at) => at.label === t);
      return tier ? `(${tier.expr})` : '';
    }).filter(Boolean);

    if (tierExprs.length === 1) {
      parts.push(tierExprs[0]);
    } else if (tierExprs.length > 1) {
      parts.push(`(${tierExprs.join(' or ')})`);
    }
  }

  if (state.dateFrom) {
    parts.push(`date >= '${state.dateFrom}'`);
  }
  if (state.dateTo) {
    parts.push(`date <= '${state.dateTo}'`);
  }

  return parts.join(' and ');
}

function getAppliedFilters(state: FilterState): AppliedFilter[] {
  const filters: AppliedFilter[] = [];
  state.kind.forEach((k) => filters.push({ category: 'Type', label: k, key: `kind-${k}` }));
  state.amountTier.forEach((t) => filters.push({ category: 'Amount', label: t, key: `tier-${t}` }));
  if (state.dateFrom && state.dateTo) {
    filters.push({ category: 'Date', label: `${state.dateFrom} to ${state.dateTo}`, key: 'date' });
  } else if (state.dateFrom) {
    filters.push({ category: 'Date', label: `from ${state.dateFrom}`, key: 'date' });
  } else if (state.dateTo) {
    filters.push({ category: 'Date', label: `to ${state.dateTo}`, key: 'date' });
  }
  return filters;
}

const columnTitleStyle: React.CSSProperties = {
  fontSize: '13px',
  fontWeight: 700,
  marginBottom: '16px',
  color: 'var(--pf-t--global--text--color--regular)',
};

const checkboxListStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '12px',
};

const TransactionFilter: React.FunctionComponent<TransactionFilterProps> = ({
  filter, onFilterChange, total, page, perPage, onSetPage, onPerPageSelect,
}) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [expression, setExpression] = React.useState(filter);
  const [namedFilters, setNamedFilters] = React.useState<FilterState>(emptyFilters);

  // Temp state for modal
  const [tempFilters, setTempFilters] = React.useState<FilterState>(emptyFilters);

  const appliedFilters = getAppliedFilters(namedFilters);

  React.useEffect(() => {
    if (isOpen) {
      setTempFilters({ ...namedFilters });
    }
  }, [isOpen]);

  const toggleTempKind = (kind: string) => {
    setTempFilters((prev) => ({
      ...prev,
      kind: prev.kind.includes(kind) ? prev.kind.filter((k) => k !== kind) : [...prev.kind, kind],
    }));
  };

  const toggleTempTier = (tier: string) => {
    setTempFilters((prev) => ({
      ...prev,
      amountTier: prev.amountTier.includes(tier) ? prev.amountTier.filter((t) => t !== tier) : [...prev.amountTier, tier],
    }));
  };

  const applyFilters = () => {
    setNamedFilters(tempFilters);
    const namedExpr = buildFilterExpr(tempFilters);
    const combined = [expression.trim(), namedExpr].filter(Boolean).join(' and ');
    onFilterChange(combined);
    setIsOpen(false);
  };

  const cancelFilters = () => {
    setIsOpen(false);
  };

  const removeFilter = (filterKey: string) => {
    let updated: FilterState;
    if (filterKey.startsWith('kind-')) {
      const value = filterKey.replace('kind-', '');
      updated = { ...namedFilters, kind: namedFilters.kind.filter((k) => k !== value) };
    } else if (filterKey.startsWith('tier-')) {
      const value = filterKey.replace('tier-', '');
      updated = { ...namedFilters, amountTier: namedFilters.amountTier.filter((t) => t !== value) };
    } else if (filterKey === 'date') {
      updated = { ...namedFilters, dateFrom: '', dateTo: '' };
    } else {
      return;
    }
    setNamedFilters(updated);
    const namedExpr = buildFilterExpr(updated);
    const combined = [expression.trim(), namedExpr].filter(Boolean).join(' and ');
    onFilterChange(combined);
  };

  const clearAllFilters = () => {
    setExpression('');
    setNamedFilters(emptyFilters);
    onFilterChange('');
  };

  const handleSearchChange = (_event: React.FormEvent, value: string) => {
    setExpression(value);
  };

  const handleSearchClear = () => {
    setExpression('');
    const namedExpr = buildFilterExpr(namedFilters);
    if (!namedExpr) onFilterChange('');
  };

  const handleSearchSubmit = () => {
    const namedExpr = buildFilterExpr(namedFilters);
    const combined = [expression.trim(), namedExpr].filter(Boolean).join(' and ');
    onFilterChange(combined);
  };

  return (
    <Toolbar>
      <ToolbarContent>
        <ToolbarGroup variant="filter-group">
          <ToolbarItem>
            <SearchInput
              placeholder="Filter expression..."
              value={expression}
              onChange={handleSearchChange}
              onClear={handleSearchClear}
              onSearch={handleSearchSubmit}
            />
          </ToolbarItem>

          <ToolbarItem>
            <Dropdown
              isOpen={isOpen}
              onOpenChange={setIsOpen}
              toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                <MenuToggle
                  ref={toggleRef}
                  onClick={() => setIsOpen(!isOpen)}
                  isExpanded={isOpen}
                  variant="default"
                >
                  <FilterIcon /> Filters
                </MenuToggle>
              )}
              popperProps={{ maxWidth: '95vw' }}
            >
              <div style={{ padding: '24px', width: '600px', maxWidth: '95vw', overflow: 'visible' }}>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '24px' }}>
                  {/* Transaction type column */}
                  <div>
                    <h3 style={columnTitleStyle}>Transaction type</h3>
                    <div style={checkboxListStyle}>
                      <Checkbox
                        id="filter-kind-debit"
                        label="Debit"
                        isChecked={tempFilters.kind.includes('debit')}
                        onChange={() => toggleTempKind('debit')}
                      />
                      <Checkbox
                        id="filter-kind-credit"
                        label="Credit"
                        isChecked={tempFilters.kind.includes('credit')}
                        onChange={() => toggleTempKind('credit')}
                      />
                    </div>
                  </div>

                  {/* Amount tier column */}
                  <div>
                    <h3 style={columnTitleStyle}>Amount</h3>
                    <div style={checkboxListStyle}>
                      {amountTiers.map((tier) => (
                        <Checkbox
                          key={tier.label}
                          id={`filter-tier-${tier.label}`}
                          label={tier.label}
                          isChecked={tempFilters.amountTier.includes(tier.label)}
                          onChange={() => toggleTempTier(tier.label)}
                        />
                      ))}
                    </div>
                  </div>

                  {/* Date range column */}
                  <div>
                    <h3 style={columnTitleStyle}>Date range</h3>
                    <div style={checkboxListStyle}>
                      <div>
                        <label style={{ display: 'block', fontSize: '13px', marginBottom: '4px', color: 'var(--pf-t--global--text--color--regular)' }}>From</label>
                        <DatePicker
                          value={tempFilters.dateFrom}
                          onChange={(_evt, value) => setTempFilters((prev) => ({ ...prev, dateFrom: value }))}
                          placeholder="YYYY-MM-DD"
                          aria-label="Start date"
                          appendTo={() => document.body}
                        />
                      </div>
                      <div>
                        <label style={{ display: 'block', fontSize: '13px', marginBottom: '4px', color: 'var(--pf-t--global--text--color--regular)' }}>To</label>
                        <DatePicker
                          value={tempFilters.dateTo}
                          onChange={(_evt, value) => setTempFilters((prev) => ({ ...prev, dateTo: value }))}
                          placeholder="YYYY-MM-DD"
                          aria-label="End date"
                          rangeStart={tempFilters.dateFrom ? new Date(`${tempFilters.dateFrom}T00:00:00`) : undefined}
                          appendTo={() => document.body}
                        />
                      </div>
                    </div>
                  </div>
                </div>

                {/* Footer */}
                <div style={{ display: 'flex', justifyContent: 'flex-start', gap: '16px', marginTop: '32px', paddingTop: '20px' }}>
                  <Button variant="primary" onClick={applyFilters}>
                    Apply filters
                  </Button>
                  <Button variant="link" onClick={cancelFilters}>
                    Cancel
                  </Button>
                </div>
              </div>
            </Dropdown>
          </ToolbarItem>
        </ToolbarGroup>

        <ToolbarItem variant="pagination" align={{ default: 'alignEnd' }}>
          <Pagination
            itemCount={total}
            perPage={perPage}
            page={page}
            onSetPage={onSetPage}
            onPerPageSelect={onPerPageSelect}
            perPageOptions={[
              { title: '25', value: 25 },
              { title: '50', value: 50 },
              { title: '100', value: 100 },
            ]}
            isCompact
          />
        </ToolbarItem>
      </ToolbarContent>

      {/* Applied filters chips */}
      {appliedFilters.length > 0 && (
        <ToolbarContent alignItems="center">
          <ToolbarItem>
            <LabelGroup categoryName="Filters">
              {appliedFilters.map((f) => (
                <Label key={f.key} onClose={() => removeFilter(f.key)}>
                  {f.label}
                </Label>
              ))}
            </LabelGroup>
          </ToolbarItem>
          <ToolbarItem>
            <span>{appliedFilters.length} filter{appliedFilters.length !== 1 ? 's' : ''} applied</span>
          </ToolbarItem>
          <ToolbarItem>
            <Button variant="link" onClick={clearAllFilters}>
              Clear all filters
            </Button>
          </ToolbarItem>
        </ToolbarContent>
      )}
    </Toolbar>
  );
};

export { TransactionFilter };
