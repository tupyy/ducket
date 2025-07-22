import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { ITransaction } from '@app/shared/models/transaction';

export interface ITransactionSummaryData {
  label: string;
  count: number;
  debitAmount: number;
  creditAmount: number;
}

export interface ITransactionSummary {
  type: 'filtered' | 'byLabelKey';
  data: ITransactionSummaryData[];
  totals: {
    count: number;
    debitAmount: number;
    creditAmount: number;
  };
}

export interface ITransactionSummaryState {
  summary: ITransactionSummary | null;
  loading: boolean;
}

const initialState: ITransactionSummaryState = {
  summary: null,
  loading: false,
};

// Helper function to calculate summary
const calculateSummary = (
  allTransactions: ITransaction[],
  filteredTransactions: ITransaction[]
): ITransactionSummary => {
  // Check if there are active filters by comparing source vs filtered transactions
  const hasActiveFilters = filteredTransactions.length !== allTransactions.length;

  if (hasActiveFilters && filteredTransactions.length > 0) {
    // Show totals of filtered transactions
    const transactionCount = filteredTransactions.length;
    const debitTransactions = filteredTransactions.filter((t) => t.kind === 'debit');
    const creditTransactions = filteredTransactions.filter((t) => t.kind === 'credit');
    const debitTotal = debitTransactions.reduce((sum, t) => sum + Math.abs(t.amount), 0);
    const creditTotal = creditTransactions.reduce((sum, t) => sum + Math.abs(t.amount), 0);

    return {
      type: 'filtered',
      data: [
        {
          label: 'Total Transactions',
          count: transactionCount,
          debitAmount: debitTotal,
          creditAmount: creditTotal,
        },
      ],
      totals: {
        count: transactionCount,
        debitAmount: debitTotal,
        creditAmount: creditTotal,
      },
    };
  } else {
    // Show totals by label key only
    const labelKeyTotals: { [key: string]: { count: number; debitAmount: number; creditAmount: number } } = {};

    allTransactions.forEach((transaction) => {
      transaction.labels.forEach((label) => {
        if (!labelKeyTotals[label.key]) {
          labelKeyTotals[label.key] = { count: 0, debitAmount: 0, creditAmount: 0 };
        }
        labelKeyTotals[label.key].count += 1;

        if (transaction.kind === 'debit') {
          labelKeyTotals[label.key].debitAmount += Math.abs(transaction.amount);
        } else if (transaction.kind === 'credit') {
          labelKeyTotals[label.key].creditAmount += Math.abs(transaction.amount);
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
  }
};

export const transactionSummarySlice = createSlice({
  name: 'transactionSummary',
  initialState,
  reducers: {
    calculateSummary: (
      state,
      action: PayloadAction<{ allTransactions: ITransaction[]; filteredTransactions: ITransaction[] }>
    ) => {
      state.loading = false;
      state.summary = calculateSummary(action.payload.allTransactions, action.payload.filteredTransactions);
    },

    clearSummary: (state) => {
      state.summary = null;
      state.loading = false;
    },

    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },

    reset: () => initialState,
  },
});

export const { calculateSummary: calculateTransactionSummary, clearSummary, setLoading, reset } =
  transactionSummarySlice.actions;

export default transactionSummarySlice.reducer; 