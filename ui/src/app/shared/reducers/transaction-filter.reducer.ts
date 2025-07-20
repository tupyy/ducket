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
      
      // Start with all source transactions from the store
      let filtered = state.sourceTransactions;

      // Filter by labels
      if (selectedLabels.length > 0) {
        filtered = filtered.filter((transaction) =>
          selectedLabels.some((selectedLabel) => {
            // Check if this is a wildcard filter (e.g., "income=*")
            if (selectedLabel.endsWith('=*')) {
              // Extract the key part (everything before "=*")
              const labelKey = selectedLabel.slice(0, -2);
              // Match any transaction that has this label key, regardless of value
              return transaction.labels.some((label) => label.key === labelKey);
            } else {
              // Exact match for specific key=value pairs
              return transaction.labels.some((label) => `${label.key}=${label.value}` === selectedLabel);
            }
          })
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

      // Update the filtered transactions in state
      state.filteredTransactions = filtered;
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