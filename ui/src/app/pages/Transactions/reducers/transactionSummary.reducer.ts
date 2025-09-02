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
const calculateSummary = (filteredTransactions: ITransaction[]): ITransactionSummary => {
  // Show totals by label key only
  const labelKeyTotals: { [key: string]: { count: number; debitAmount: number; creditAmount: number } } = {};

  filteredTransactions.forEach((transaction) => {
    transaction.labels.forEach((label) => {
      const key = label.key +"="+ label.value;
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

export const transactionSummarySlice = createSlice({
  name: 'transactionSummary',
  initialState,
  reducers: {
    calculateSummary: (state, action: PayloadAction<{ filteredTransactions: ITransaction[] }>) => {
      state.loading = false;
      state.summary = calculateSummary(action.payload.filteredTransactions);
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

export const {
  calculateSummary: calculateTransactionSummary,
  clearSummary,
  setLoading,
  reset,
} = transactionSummarySlice.actions;

export default transactionSummarySlice.reducer;
