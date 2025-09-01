import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit';
import { ITransaction } from '@app/shared/models/transaction';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';

export interface ITransactionFilterState {
  sourceTransactions: ITransaction[];
  filteredTransactions: ITransaction[];
  // Filter states
  selectedLabels: string[];
  selectedTransactionTypes: string[];
  selectedAccounts: number[];
  descriptionFilter: string;
  showOnlyUnlabeled: boolean;
  // Date range state
  dateRange: {
    startDate: string;
    endDate: string;
  };
  // Pagination states
  page: number;
  perPage: number;
  // Sorting states
  sortDirection: 'asc' | 'desc' | null;
  sortIndex: number | null;
  // Row expansion states
  expandedTransactions: string[];
  allExpanded: boolean;
  // Async state for filtering
  filtering: boolean;
  filterError: string;
  // Selection state
  selectedTransactions: string[];
}

const initialState: ITransactionFilterState = {
  sourceTransactions: [],
  filteredTransactions: [],
  // Filter states
  selectedLabels: [],
  selectedTransactionTypes: [],
  selectedAccounts: [],
  descriptionFilter: '',
  showOnlyUnlabeled: false,
  // Date range state - default to last 30 days
  dateRange: {
    startDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    endDate: new Date().toISOString().split('T')[0],
  },
  // Pagination states
  page: 1,
  perPage: 50,
  // Sorting states
  sortDirection: 'desc',
  sortIndex: 0, // 0 = date column
  // Row expansion states
  expandedTransactions: [],
  allExpanded: false,
  // Async state for filtering
  filtering: false,
  filterError: '',
  // Selection state
  selectedTransactions: [],
};

// Async thunk for applying filters
export const applyFilters = createAsyncThunk(
  'transactionFilter/applyFilters',
  async (filterParams: {
    selectedLabels: string[];
    selectedTransactionTypes: string[];
    selectedAccounts: number[];
    descriptionFilter?: string;
    showOnlyUnlabeled?: boolean;
  }, { getState }) => {
    // Simulate async operation - you can add actual async logic here if needed
    await new Promise(resolve => setTimeout(resolve, 100));
    
    const state = getState() as { transactionFilter: ITransactionFilterState };
    const { selectedLabels, selectedTransactionTypes, selectedAccounts, descriptionFilter, showOnlyUnlabeled } = filterParams;

    // Start with all source transactions from the store
    let filtered = state.transactionFilter.sourceTransactions;

    // Filter by unlabeled transactions first if requested
    const currentShowOnlyUnlabeled = showOnlyUnlabeled !== undefined ? showOnlyUnlabeled : state.transactionFilter.showOnlyUnlabeled;
    if (currentShowOnlyUnlabeled) {
      filtered = filtered.filter((transaction) => transaction.labels.length === 0);
    }

    // Filter by labels (only if not showing unlabeled transactions)
    if (!currentShowOnlyUnlabeled && selectedLabels.length > 0) {
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

    // Filter by description
    const currentDescriptionFilter = descriptionFilter !== undefined ? descriptionFilter : state.transactionFilter.descriptionFilter;
    if (currentDescriptionFilter.trim()) {
      const filterText = currentDescriptionFilter.toLowerCase().trim();
      filtered = filtered.filter((transaction) =>
        transaction.description.toLowerCase().includes(filterText)
      );
    }

    return {
      filteredTransactions: filtered,
      selectedLabels,
      selectedTransactionTypes,
      selectedAccounts,
      descriptionFilter: currentDescriptionFilter,
      showOnlyUnlabeled: currentShowOnlyUnlabeled,
    };
  },
  { serializeError: serializeAxiosError }
);

export const transactionFilterSlice = createSlice({
  name: 'transactionFilter',
  initialState,
  reducers: {
    setSourceTransactions: (state, action: PayloadAction<ITransaction[]>) => {
      state.sourceTransactions = action.payload;
      // Only reset filtered transactions when no filters are active
      // If filters are active, let applyFilters handle the filtering
      if (
        state.selectedLabels.length === 0 &&
        state.selectedTransactionTypes.length === 0 &&
        state.selectedAccounts.length === 0 &&
        state.descriptionFilter.trim() === '' &&
        !state.showOnlyUnlabeled
      ) {
        state.filteredTransactions = action.payload;
      }
      // Reset pagination to first page when source data changes
      state.page = 1;
    },

    // Date range actions
    setDateRange: (state, action: PayloadAction<{ startDate: string; endDate: string }>) => {
      state.dateRange = action.payload;
      // Reset pagination when date range changes
      state.page = 1;
    },

    // Filter actions
    setSelectedLabels: (state, action: PayloadAction<string[]>) => {
      state.selectedLabels = action.payload;
      state.page = 1; // Reset pagination when filters change
    },

    setSelectedTransactionTypes: (state, action: PayloadAction<string[]>) => {
      state.selectedTransactionTypes = action.payload;
      state.page = 1; // Reset pagination when filters change
    },

    setSelectedAccounts: (state, action: PayloadAction<number[]>) => {
      state.selectedAccounts = action.payload;
      state.page = 1; // Reset pagination when filters change
    },

    setDescriptionFilter: (state, action: PayloadAction<string>) => {
      state.descriptionFilter = action.payload;
      state.page = 1; // Reset pagination when filters change
    },

    setShowOnlyUnlabeled: (state, action: PayloadAction<boolean>) => {
      state.showOnlyUnlabeled = action.payload;
      // Clear selected labels when showing only unlabeled transactions
      if (action.payload) {
        state.selectedLabels = [];
      }
      state.page = 1; // Reset pagination when filters change
    },

    clearAllFilters: (state) => {
      state.selectedLabels = [];
      state.selectedTransactionTypes = [];
      state.selectedAccounts = [];
      state.descriptionFilter = '';
      state.showOnlyUnlabeled = false;
      state.page = 1; // Reset pagination when filters change
    },

    // Pagination actions
    setPage: (state, action: PayloadAction<number>) => {
      state.page = action.payload;
    },

    setPerPage: (state, action: PayloadAction<{ perPage: number; page?: number }>) => {
      state.perPage = action.payload.perPage;
      if (action.payload.page !== undefined) {
        state.page = action.payload.page;
      }
    },

    // Sorting actions
    setSorting: (state, action: PayloadAction<{ sortIndex: number; sortDirection: 'asc' | 'desc' }>) => {
      state.sortIndex = action.payload.sortIndex;
      state.sortDirection = action.payload.sortDirection;
      state.page = 1; // Reset pagination when sorting changes
    },

    clearSorting: (state) => {
      state.sortIndex = null;
      state.sortDirection = null;
      state.page = 1; // Reset pagination when sorting changes
    },

    // Row expansion actions
    setTransactionExpanded: (state, action: PayloadAction<{ href: string; isExpanding: boolean }>) => {
      const { href, isExpanding } = action.payload;
      const transactionIndex = state.expandedTransactions.findIndex((expandedHref) => expandedHref === href);

      if (isExpanding && transactionIndex === -1) {
        state.expandedTransactions.push(href);
      } else if (!isExpanding && transactionIndex !== -1) {
        state.expandedTransactions.splice(transactionIndex, 1);
      }
    },

    toggleAllExpanded: (state, action: PayloadAction<string[]>) => {
      const currentPageHrefs = action.payload;
      const areAllCurrentPageExpanded = currentPageHrefs.length > 0 &&
        currentPageHrefs.every(href => state.expandedTransactions.includes(href));

      if (areAllCurrentPageExpanded) {
        // Collapse all current page transactions
        state.expandedTransactions = state.expandedTransactions.filter(href => !currentPageHrefs.includes(href));
      } else {
        // Expand all current page transactions
        const newExpanded = Array.from(new Set([...state.expandedTransactions, ...currentPageHrefs]));
        state.expandedTransactions = newExpanded;
      }
    },

    // Selection actions
    setTransactionSelected: (state, action: PayloadAction<{ href: string; isSelected: boolean }>) => {
      const { href, isSelected } = action.payload;
      const transactionIndex = state.selectedTransactions.findIndex((selectedHref) => selectedHref === href);

      if (isSelected && transactionIndex === -1) {
        state.selectedTransactions.push(href);
      } else if (!isSelected && transactionIndex !== -1) {
        state.selectedTransactions.splice(transactionIndex, 1);
      }
    },

    selectAllTransactions: (state, action: PayloadAction<string[]>) => {
      const currentPageHrefs = action.payload;
      const areAllCurrentPageSelected = currentPageHrefs.length > 0 &&
        currentPageHrefs.every(href => state.selectedTransactions.includes(href));

      if (areAllCurrentPageSelected) {
        // Deselect all current page transactions
        state.selectedTransactions = state.selectedTransactions.filter(href => !currentPageHrefs.includes(href));
      } else {
        // Select all current page transactions
        const newSelected = Array.from(new Set([...state.selectedTransactions, ...currentPageHrefs]));
        state.selectedTransactions = newSelected;
      }
    },

    clearSelection: (state) => {
      state.selectedTransactions = [];
    },

    reset: () => initialState,
  },
  extraReducers: (builder) => {
    builder
      .addCase(applyFilters.pending, (state) => {
        state.filtering = true;
        state.filterError = '';
      })
      .addCase(applyFilters.fulfilled, (state, action) => {
        state.filtering = false;
        state.filterError = '';
        state.filteredTransactions = action.payload.filteredTransactions;
        state.selectedLabels = action.payload.selectedLabels;
        state.selectedTransactionTypes = action.payload.selectedTransactionTypes;
        state.selectedAccounts = action.payload.selectedAccounts;
        state.descriptionFilter = action.payload.descriptionFilter;
        state.showOnlyUnlabeled = action.payload.showOnlyUnlabeled;
        // Reset pagination to first page when filters change
        state.page = 1;
      })
      .addCase(applyFilters.rejected, (state, action) => {
        state.filtering = false;
        state.filterError = action.error.message || 'Failed to apply filters';
      });
  },
});

export const {
  setSourceTransactions,
  setDateRange,
  setSelectedLabels,
  setSelectedTransactionTypes,
  setSelectedAccounts,
  setDescriptionFilter,
  setShowOnlyUnlabeled,
  clearAllFilters,
  setPage,
  setPerPage,
  setSorting,
  clearSorting,
  setTransactionExpanded,
  toggleAllExpanded,
  setTransactionSelected,
  selectAllTransactions,
  clearSelection,
  reset,
} = transactionFilterSlice.actions;

export default transactionFilterSlice.reducer;
