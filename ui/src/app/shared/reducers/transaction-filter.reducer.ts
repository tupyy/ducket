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
  selectedTags: string[],
  selectedTransactionTypes: string[],
  selectedAccounts: number[]
): ITransaction[] => {
  let filtered = transactions;

  // Filter by tags
  if (selectedTags.length > 0) {
    filtered = filtered.filter((transaction) =>
      selectedTags.some((selectedTag) => transaction.tags.some((tag) => tag.value === selectedTag))
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
      selectedTags: string[];
      selectedTransactionTypes: string[];
      selectedAccounts: number[];
    }>) => {
      const { selectedTags, selectedTransactionTypes, selectedAccounts } = action.payload;
      state.filteredTransactions = filterTransactions(
        state.sourceTransactions,
        selectedTags,
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