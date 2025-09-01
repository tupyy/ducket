import * as React from 'react';
import { 
  EuiPageSection, 
  EuiTitle, 
  EuiPanel, 
  EuiFlexGroup, 
  EuiFlexItem, 
  EuiSpacer, 
  EuiText,
  EuiStat,
  EuiCallOut
} from '@elastic/eui';
import { TimePicker } from '@app/shared/components/time-picker';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { calculateDateRange, getRelativeTimeRange } from '@app/utils/dateUtils';

const Dashboard: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const { transactions, loading } = useAppSelector((state) => state.transactions);

  // Initialize with 1 year default date range
  const defaultDateRange = React.useMemo(() => calculateDateRange('last 1 year'), []);
  const [startDate, setStartDate] = React.useState<string>(defaultDateRange.startDateValue);
  const [endDate, setEndDate] = React.useState<string>(defaultDateRange.endDateValue);

  const handleDateChange = (start: string, end: string) => {
    setStartDate(start);
    setEndDate(end);
    dispatch(getTransactions({ startDate: start, endDate: end }));
  };

  // Calculate basic statistics
  const stats = React.useMemo(() => {
    const totalTransactions = transactions.length;
    const totalDebit = transactions.filter(t => t.kind === 'debit').reduce((sum, t) => sum + t.amount, 0);
    const totalCredit = transactions.filter(t => t.kind === 'credit').reduce((sum, t) => sum + t.amount, 0);
    const uniqueAccounts = new Set(transactions.map(t => t.account)).size;
    const uniqueLabels = new Set(transactions.flatMap(t => t.labels.map(l => l.key))).size;

    return {
      totalTransactions,
      totalDebit,
      totalCredit,
      balance: totalCredit - totalDebit,
      uniqueAccounts,
      uniqueLabels
    };
  }, [transactions]);

  // Fetch initial data on component mount with default date range
  React.useEffect(() => {
    dispatch(getTransactions({ startDate, endDate }));
  }, [dispatch, startDate, endDate]);

  return (
    <EuiPageSection>
      <EuiTitle size="l">
        <h1>Dashboard</h1>
      </EuiTitle>

      <EuiSpacer size="l" />

      <EuiFlexGroup gutterSize="s" alignItems="center">
        <EuiFlexItem grow={false}>
          <TimePicker 
            onDateChange={handleDateChange}
            initialStartDate={startDate}
            initialEndDate={endDate}
          />
        </EuiFlexItem>
      </EuiFlexGroup>

      {startDate && endDate && (
        <>
          <EuiSpacer size="m" />
          <EuiCallOut size="s" title="Selected Range" iconType="calendar">
            <EuiText size="s">{getRelativeTimeRange(startDate, endDate)}</EuiText>
          </EuiCallOut>
        </>
      )}

      <EuiSpacer size="l" />

      <EuiFlexGroup gutterSize="l" wrap>
        <EuiFlexItem>
          <EuiStat
            title={stats.totalTransactions.toString()}
            description="Total Transactions"
            color="primary"
            titleSize="l"
          />
        </EuiFlexItem>
        
        <EuiFlexItem>
          <EuiStat
            title={`€${stats.totalCredit.toFixed(2)}`}
            description="Total Credits"
            color="success"
            titleSize="l"
          />
        </EuiFlexItem>
        
        <EuiFlexItem>
          <EuiStat
            title={`€${stats.totalDebit.toFixed(2)}`}
            description="Total Debits"
            color="danger"
            titleSize="l"
          />
        </EuiFlexItem>
        
        <EuiFlexItem>
          <EuiStat
            title={`€${stats.balance.toFixed(2)}`}
            description="Net Balance"
            color={stats.balance >= 0 ? "success" : "danger"}
            titleSize="l"
          />
        </EuiFlexItem>
      </EuiFlexGroup>

      <EuiSpacer size="l" />

      <EuiFlexGroup gutterSize="l" wrap>
        <EuiFlexItem>
          <EuiStat
            title={stats.uniqueAccounts.toString()}
            description="Unique Accounts"
            color="accent"
            titleSize="m"
          />
        </EuiFlexItem>
        
        <EuiFlexItem>
          <EuiStat
            title={stats.uniqueLabels.toString()}
            description="Unique Labels"
            color="accent"
            titleSize="m"
          />
        </EuiFlexItem>
      </EuiFlexGroup>

      <EuiSpacer size="l" />

      <EuiPanel paddingSize="l">
        <EuiTitle size="m">
          <h2>Quick Actions</h2>
        </EuiTitle>
        <EuiSpacer size="m" />
        <EuiText>
          <p>Use the navigation above to:</p>
          <ul>
            <li>View and manage <strong>Transactions</strong></li>
            <li>Create and edit <strong>Rules</strong> for automatic labeling</li>
            <li>Browse existing <strong>Labels</strong></li>
          </ul>
        </EuiText>
      </EuiPanel>
    </EuiPageSection>
  );
};

export { Dashboard };
