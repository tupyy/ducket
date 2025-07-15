import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { ITransaction } from '@app/shared/models/transaction';

export interface ITransactionFilterState {
  sourceTransactions: ITransaction[];
  filteredTransactions: ITransaction[];
}

const initialState: ITransactionFilterState = {
  sourceTransactions: [],
  filteredTransactions: [],
};

// Helper function for filtering transactions
const filterTransactions = (
  transactions: ITransaction[],
  selectedLabels: string[],
  selectedTransactionTypes: string[],
  selectedAccounts: number[]
): ITransaction[] => {
  let filtered = transactions;

  // Filter by labels
  if (selectedLabels.length > 0) {
    filtered = filtered.filter((transaction) =>
      selectedLabels.some((selectedLabel) => transaction.labels.some((label) => `${label.key}=${label.value}` === selectedLabel))
    );
  }

  // Filter by transaction types
  if (selectedTransactionTypes.length > 0) {
    filtered = filtered.filter((transaction) => selectedTransactionTypes.includes(transaction.kind));
  }

  // Filter by accounts
  if (selectedAccounts.length > 0) {
    filtered = filtered.filter((transaction) => selectedAccounts.includes(transaction.account));
  }

  return filtered;
};

export const transactionFilterSlice = createSlice({
  name: 'transactionFilter',
  initialState,
  reducers: {
    setSourceTransactions: (state, action: PayloadAction<ITransaction[]>) => {
      state.sourceTransactions = action.payload;
      // Reset filtered transactions when source changes
      state.filteredTransactions = action.payload;
    },
    
    applyFilters: (state, action: PayloadAction<{
      selectedLabels: string[];
      selectedTransactionTypes: string[];
      selectedAccounts: number[];
    }>) => {
      const { selectedLabels, selectedTransactionTypes, selectedAccounts } = action.payload;
      state.filteredTransactions = filterTransactions(
        state.sourceTransactions,
        selectedLabels,
        selectedTransactionTypes,
        selectedAccounts
      );
    },
    
    reset: () => initialState,
  },
});

export const {
  setSourceTransactions,
  applyFilters,
  reset,
} = transactionFilterSlice.actions;

export default transactionFilterSlice.reducer; 