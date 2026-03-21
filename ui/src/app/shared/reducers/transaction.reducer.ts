import axios from 'axios';
import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';
import { ITransaction } from '@app/shared/models/transaction';

const apiUrl = 'api/v1/transactions';

export interface SortField {
  field: string;
  direction: 'asc' | 'desc';
}

interface TransactionState {
  loading: boolean;
  errorMessage: string;
  transactions: ITransaction[];
  total: number;
  page: number;
  perPage: number;
  filter: string;
  selectedTags: string[];
  selectedAccounts: number[];
  selectedKind: string | null;
  sort: SortField[];
}

const initialState: TransactionState = {
  loading: false,
  errorMessage: '',
  transactions: [],
  total: 0,
  page: 1,
  perPage: 50,
  filter: '',
  selectedTags: [],
  selectedAccounts: [],
  selectedKind: null,
  sort: [{ field: 'date', direction: 'desc' }],
};

export function buildCompositeFilter(filter: string, accounts: number[], kind: string | null): string {
  const parts: string[] = [];
  if (filter) parts.push(filter);
  if (accounts.length === 1) {
    parts.push(`account = ${accounts[0]}`);
  } else if (accounts.length > 1) {
    const exprs = accounts.map((a) => `account = ${a}`);
    parts.push(`(${exprs.join(' or ')})`);
  }
  if (kind) {
    parts.push(`kind = '${kind}'`);
  }
  return parts.join(' and ');
}

export const getTransactions = createAsyncThunk(
  'transactions/get',
  async (params: { filter?: string; tags?: string[]; sort?: SortField[]; limit?: number; offset?: number } | void) => {
    params = params || {};
    const searchParams = new URLSearchParams();
    if (params.filter) searchParams.append('filter', params.filter);
    if (params.tags) {
      params.tags.forEach((t) => searchParams.append('tags', t));
    }
    if (params.sort) {
      params.sort.forEach((s) => searchParams.append('sort', `${s.field}:${s.direction}`));
    }
    if (params.limit) searchParams.append('limit', params.limit.toString());
    if (params.offset !== undefined) searchParams.append('offset', params.offset.toString());

    const query = searchParams.toString();
    const url = query ? `${apiUrl}?${query}` : apiUrl;
    return axios.get<{ items: ITransaction[]; total: number }>(url);
  },
  { serializeError: serializeAxiosError },
);

export const deleteTransaction = createAsyncThunk(
  'transactions/delete',
  async (id: number, thunkAPI) => {
    await axios.delete(`${apiUrl}/${id}`);
    thunkAPI.dispatch(getTransactions(undefined));
  },
  { serializeError: serializeAxiosError },
);

export const TransactionSlice = createSlice({
  name: 'transactions',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    setPage(state, action: PayloadAction<number>) {
      state.page = action.payload;
    },
    setPerPage(state, action: PayloadAction<number>) {
      state.perPage = action.payload;
      state.page = 1;
    },
    setFilter(state, action: PayloadAction<string>) {
      state.filter = action.payload;
      state.page = 1;
    },
    setSort(state, action: PayloadAction<SortField[]>) {
      state.sort = action.payload;
      state.page = 1;
    },
    addTag(state, action: PayloadAction<string>) {
      if (!state.selectedTags.includes(action.payload)) {
        state.selectedTags.push(action.payload);
        state.page = 1;
      }
    },
    removeTag(state, action: PayloadAction<string>) {
      state.selectedTags = state.selectedTags.filter((t) => t !== action.payload);
      state.page = 1;
    },
    clearTags(state) {
      state.selectedTags = [];
      state.page = 1;
    },
    addAccount(state, action: PayloadAction<number>) {
      if (!state.selectedAccounts.includes(action.payload)) {
        state.selectedAccounts.push(action.payload);
        state.page = 1;
      }
    },
    removeAccount(state, action: PayloadAction<number>) {
      state.selectedAccounts = state.selectedAccounts.filter((a) => a !== action.payload);
      state.page = 1;
    },
    toggleKind(state, action: PayloadAction<string>) {
      state.selectedKind = state.selectedKind === action.payload ? null : action.payload;
      state.page = 1;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(getTransactions.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(getTransactions.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to load transactions';
      })
      .addCase(getTransactions.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.transactions = action.payload.data.items;
        state.total = action.payload.data.total;
      });
  },
});

export const {
  reset, setPage, setPerPage, setFilter, setSort,
  addTag, removeTag, clearTags,
  addAccount, removeAccount, toggleKind,
} = TransactionSlice.actions;
export default TransactionSlice.reducer;
