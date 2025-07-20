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
  selectedAccounts: number[],
  startDate?: string,
  endDate?: string
): ITransaction[] => {
  let filtered = transactions;

  // Filter by date range (parse date strings only when needed)
  if (startDate || endDate) {
    filtered = filtered.filter((transaction) => {
      const transactionDate = new Date(transaction.date);
      
      if (startDate) {
        const start = new Date(startDate);
        if (transactionDate < start) {
          return false;
        }
      }
      
      if (endDate) {
        const end = new Date(endDate);
        // Set end date to end of day for inclusive filtering
        end.setHours(23, 59, 59, 999);
        if (transactionDate > end) {
          return false;
        }
      }
      
      return true;
    });
  }

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
      startDate?: string;
      endDate?: string;
    }>) => {
      const { selectedLabels, selectedTransactionTypes, selectedAccounts, startDate, endDate } = action.payload;
      state.filteredTransactions = filterTransactions(
        state.sourceTransactions,
        selectedLabels,
        selectedTransactionTypes,
        selectedAccounts,
        startDate,
        endDate
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