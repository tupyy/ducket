import { ITransaction } from '@app/shared/models/transaction';
import { ITransactionSummary } from '../reducers/transactionSummary.reducer';

export interface IHierarchicalSummaryData {
  label: string;
  count: number;
  debitAmount: number;
  creditAmount: number;
  isParent: boolean;
  children?: IHierarchicalSummaryData[];
}

export interface IHierarchicalSummary {
  type: 'byLabelKey' | 'hierarchical';
  data: IHierarchicalSummaryData[];
  totals: {
    count: number;
    debitAmount: number;
    creditAmount: number;
  };
}

export const calculateSummary = (filteredTransactions: ITransaction[]): ITransactionSummary => {
  // Show totals by label key only
  const labelKeyTotals: { [key: string]: { count: number; debitAmount: number; creditAmount: number } } = {};

  filteredTransactions.forEach((transaction) => {
    transaction.labels.forEach((label) => {
      const key = label.key + "=" + label.value;
      if (!labelKeyTotals[key]) {
        labelKeyTotals[key] = { count: 0, debitAmount: 0, creditAmount: 0 };
      }
      labelKeyTotals[key].count += 1;

      if (transaction.kind === 'debit') {
        labelKeyTotals[key].debitAmount += Math.abs(transaction.amount);
      } else if (transaction.kind === 'credit') {
        labelKeyTotals[key].creditAmount += Math.abs(transaction.amount);
      }
    });
  });

  const data = Object.entries(labelKeyTotals)
    .map(([key, totals]) => ({
      label: key,
      count: totals.count,
      debitAmount: totals.debitAmount,
      creditAmount: totals.creditAmount,
    }))
    .sort(
      (a, b) =>
        Math.abs(b.debitAmount) + Math.abs(b.creditAmount) - (Math.abs(a.debitAmount) + Math.abs(a.creditAmount))
    ); // Sort by total absolute amount descending

  // Calculate totals from the data
  const totals = data.reduce(
    (acc, row) => ({
      count: acc.count + row.count,
      debitAmount: acc.debitAmount + row.debitAmount,
      creditAmount: acc.creditAmount + row.creditAmount,
    }),
    { count: 0, debitAmount: 0, creditAmount: 0 }
  );

  return {
    type: 'byLabelKey',
    data,
    totals,
  };
};

export const calculateHierarchicalSummary = (filteredTransactions: ITransaction[]): IHierarchicalSummary => {
  // First, group by label key and collect all key=value pairs
  const labelGroups: { [key: string]: { [keyValue: string]: { count: number; debitAmount: number; creditAmount: number } } } = {};

  filteredTransactions.forEach((transaction) => {
    transaction.labels.forEach((label) => {
      const key = label.key;
      const keyValue = `${label.key}=${label.value}`;
      
      if (!labelGroups[key]) {
        labelGroups[key] = {};
      }
      
      if (!labelGroups[key][keyValue]) {
        labelGroups[key][keyValue] = { count: 0, debitAmount: 0, creditAmount: 0 };
      }
      
      labelGroups[key][keyValue].count += 1;

      if (transaction.kind === 'debit') {
        labelGroups[key][keyValue].debitAmount += Math.abs(transaction.amount);
      } else if (transaction.kind === 'credit') {
        labelGroups[key][keyValue].creditAmount += Math.abs(transaction.amount);
      }
    });
  });

  // Convert to hierarchical structure
  const data: IHierarchicalSummaryData[] = Object.entries(labelGroups)
    .map(([key, keyValues]) => {
      // Calculate totals for this key
      const keyTotals = Object.values(keyValues).reduce(
        (acc, item) => ({
          count: acc.count + item.count,
          debitAmount: acc.debitAmount + item.debitAmount,
          creditAmount: acc.creditAmount + item.creditAmount,
        }),
        { count: 0, debitAmount: 0, creditAmount: 0 }
      );

      // Create children for each key=value pair
      const children: IHierarchicalSummaryData[] = Object.entries(keyValues)
        .map(([keyValue, totals]) => ({
          label: keyValue,
          count: totals.count,
          debitAmount: totals.debitAmount,
          creditAmount: totals.creditAmount,
          isParent: false,
        }))
        .sort(
          (a, b) =>
            Math.abs(b.debitAmount) + Math.abs(b.creditAmount) - (Math.abs(a.debitAmount) + Math.abs(a.creditAmount))
        );

      return {
        label: key,
        count: keyTotals.count,
        debitAmount: keyTotals.debitAmount,
        creditAmount: keyTotals.creditAmount,
        isParent: true,
        children,
      };
    })
    .sort(
      (a, b) =>
        Math.abs(b.debitAmount) + Math.abs(b.creditAmount) - (Math.abs(a.debitAmount) + Math.abs(a.creditAmount))
    );

  // Calculate totals from the data
  const totals = data.reduce(
    (acc, row) => ({
      count: acc.count + row.count,
      debitAmount: acc.debitAmount + row.debitAmount,
      creditAmount: acc.creditAmount + row.creditAmount,
    }),
    { count: 0, debitAmount: 0, creditAmount: 0 }
  );

  return {
    type: 'hierarchical',
    data,
    totals,
  };
};