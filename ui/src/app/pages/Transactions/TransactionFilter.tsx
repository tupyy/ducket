import * as React from 'react';
import {
  Button,
  TextInputGroup,
  TextInputGroupMain,
  TextInputGroupUtilities,
  Popper,
  Panel,
  PanelMain,
  PanelMainBody,
  PanelHeader,
  PanelFooter,
  Checkbox,
  Title,
  Divider,
  Flex,
  FlexItem,
  DatePicker,
  LabelGroup,
  Label,
} from '@patternfly/react-core';
import FilterIcon from '@patternfly/react-icons/dist/esm/icons/filter-icon';
import TimesIcon from '@patternfly/react-icons/dist/esm/icons/times-icon';

interface TransactionFilterProps {
  filter: string;
  onFilterChange: (filter: string) => void;
}

interface FilterState {
  kind: string[];
  amountTier: string[];
  dateFrom: string;
  dateTo: string;
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

function describeFilter(state: FilterState): string[] {
  const chips: string[] = [];
  state.kind.forEach((k) => chips.push(`Type: ${k}`));
  state.amountTier.forEach((t) => chips.push(`Amount: ${t}`));
  if (state.dateFrom && state.dateTo) {
    chips.push(`Date: ${state.dateFrom} to ${state.dateTo}`);
  } else if (state.dateFrom) {
    chips.push(`Date: from ${state.dateFrom}`);
  } else if (state.dateTo) {
    chips.push(`Date: to ${state.dateTo}`);
  }
  return chips;
}

const TransactionFilter: React.FunctionComponent<TransactionFilterProps> = ({ filter, onFilterChange }) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [expression, setExpression] = React.useState(filter);
  const [namedFilters, setNamedFilters] = React.useState<FilterState>(emptyFilters);

  const toggleRef = React.useRef<HTMLDivElement>(null);
  const menuRef = React.useRef<HTMLDivElement>(null);

  const activeChips = describeFilter(namedFilters);

  React.useEffect(() => {
    if (!isOpen) return;
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Node;
      if (
        toggleRef.current && !toggleRef.current.contains(target) &&
        menuRef.current && !menuRef.current.contains(target)
      ) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen]);

  const handleToggleKind = (kind: string, checked: boolean) => {
    setNamedFilters((prev) => ({
      ...prev,
      kind: checked ? [...prev.kind, kind] : prev.kind.filter((k) => k !== kind),
    }));
  };

  const handleToggleTier = (tier: string, checked: boolean) => {
    setNamedFilters((prev) => ({
      ...prev,
      amountTier: checked ? [...prev.amountTier, tier] : prev.amountTier.filter((t) => t !== tier),
    }));
  };

  const handleApply = () => {
    const namedExpr = buildFilterExpr(namedFilters);
    const combined = [expression.trim(), namedExpr].filter(Boolean).join(' and ');
    onFilterChange(combined);
    setIsOpen(false);
  };

  const handleClear = () => {
    setExpression('');
    setNamedFilters(emptyFilters);
    onFilterChange('');
    setIsOpen(false);
  };

  const handleRemoveChip = (chip: string) => {
    let updated: FilterState;
    if (chip.startsWith('Type: ')) {
      const value = chip.replace('Type: ', '');
      updated = { ...namedFilters, kind: namedFilters.kind.filter((k) => k !== value) };
    } else if (chip.startsWith('Amount: ')) {
      const value = chip.replace('Amount: ', '');
      updated = { ...namedFilters, amountTier: namedFilters.amountTier.filter((t) => t !== value) };
    } else if (chip.startsWith('Date:')) {
      updated = { ...namedFilters, dateFrom: '', dateTo: '' };
    } else {
      return;
    }
    setNamedFilters(updated);
    const namedExpr = buildFilterExpr(updated);
    const combined = [expression.trim(), namedExpr].filter(Boolean).join(' and ');
    onFilterChange(combined);
  };

  const handleExpressionKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      handleApply();
    }
  };

  const filterPanel = (
    <Panel ref={menuRef} variant="raised" style={{ minWidth: '500px' }}>
      <PanelHeader>
        <Title headingLevel="h4" size="md">Filters</Title>
      </PanelHeader>
      <Divider />
      <PanelMain>
        <PanelMainBody>
          <Flex direction={{ default: 'row' }} gap={{ default: 'gap2xl' }}>
            <FlexItem>
              <Title headingLevel="h5" size="sm" style={{ marginBottom: '0.5rem' }}>
                Transaction type
              </Title>
              <Checkbox
                id="filter-kind-debit"
                label="Debit"
                isChecked={namedFilters.kind.includes('debit')}
                onChange={(_evt, checked) => handleToggleKind('debit', checked)}
              />
              <Checkbox
                id="filter-kind-credit"
                label="Credit"
                isChecked={namedFilters.kind.includes('credit')}
                onChange={(_evt, checked) => handleToggleKind('credit', checked)}
              />
            </FlexItem>
            <FlexItem>
              <Title headingLevel="h5" size="sm" style={{ marginBottom: '0.5rem' }}>
                Amount tier
              </Title>
              {amountTiers.map((tier) => (
                <Checkbox
                  key={tier.label}
                  id={`filter-tier-${tier.label}`}
                  label={tier.label}
                  isChecked={namedFilters.amountTier.includes(tier.label)}
                  onChange={(_evt, checked) => handleToggleTier(tier.label, checked)}
                />
              ))}
            </FlexItem>
            <FlexItem>
              <Title headingLevel="h5" size="sm" style={{ marginBottom: '0.5rem' }}>
                Date range
              </Title>
              <div style={{ marginBottom: '0.5rem' }}>
                <label htmlFor="filter-date-from" style={{ display: 'block', fontSize: '0.875rem', marginBottom: '0.25rem' }}>From</label>
                <DatePicker
                  value={namedFilters.dateFrom}
                  onChange={(_evt, value) => setNamedFilters((prev) => ({ ...prev, dateFrom: value }))}
                  placeholder="YYYY-MM-DD"
                  aria-label="Start date"
                />
              </div>
              <div>
                <label htmlFor="filter-date-to" style={{ display: 'block', fontSize: '0.875rem', marginBottom: '0.25rem' }}>To</label>
                <DatePicker
                  value={namedFilters.dateTo}
                  onChange={(_evt, value) => setNamedFilters((prev) => ({ ...prev, dateTo: value }))}
                  placeholder="YYYY-MM-DD"
                  aria-label="End date"
                  rangeStart={namedFilters.dateFrom ? new Date(`${namedFilters.dateFrom}T00:00:00`) : undefined}
                />
              </div>
            </FlexItem>
          </Flex>
        </PanelMainBody>
      </PanelMain>
      <Divider />
      <PanelFooter>
        <Flex>
          <FlexItem>
            <Button variant="primary" size="sm" onClick={handleApply}>
              Apply filters
            </Button>
          </FlexItem>
          <FlexItem>
            <Button variant="link" size="sm" onClick={handleClear}>
              Clear
            </Button>
          </FlexItem>
        </Flex>
      </PanelFooter>
    </Panel>
  );

  return (
    <Flex direction={{ default: 'column' }} gap={{ default: 'gapSm' }}>
      <Flex gap={{ default: 'gapSm' }} alignItems={{ default: 'alignItemsCenter' }}>
        <FlexItem grow={{ default: 'grow' }} style={{ maxWidth: '400px' }}>
          <TextInputGroup>
            <TextInputGroupMain
              placeholder="Filter expression..."
              value={expression}
              onChange={(_evt, val) => setExpression(val)}
              onKeyDown={handleExpressionKeyDown}
            />
            {expression && (
              <TextInputGroupUtilities>
                <Button
                  variant="plain"
                  aria-label="Clear expression"
                  onClick={() => {
                    setExpression('');
                    if (!activeChips.length) onFilterChange('');
                  }}
                >
                  <TimesIcon />
                </Button>
              </TextInputGroupUtilities>
            )}
          </TextInputGroup>
        </FlexItem>
        <FlexItem>
          <div ref={toggleRef}>
            <Button
              variant="secondary"
              icon={<FilterIcon />}
              onClick={() => setIsOpen(!isOpen)}
            >
              Filters
            </Button>
          </div>
          <Popper
            triggerRef={toggleRef}
            popper={filterPanel}
            popperRef={menuRef}
            isVisible={isOpen}
            appendTo={() => document.body}
          />
        </FlexItem>
      </Flex>

      {activeChips.length > 0 && (
        <FlexItem>
          <LabelGroup>
            {activeChips.map((chip) => (
              <Label key={chip} variant="outline" onClose={() => handleRemoveChip(chip)}>
                {chip}
              </Label>
            ))}
          </LabelGroup>
        </FlexItem>
      )}
    </Flex>
  );
};

export { TransactionFilter };
