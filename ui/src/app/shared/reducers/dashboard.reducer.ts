import axios from 'axios';
import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';
import { ITransaction } from '@app/shared/models/transaction';

const summaryUrl = 'api/v1/summary';
const txnUrl = 'api/v1/transactions';

export interface SummaryOverview {
  total_transactions: number;
  total_debit: number;
  total_credit: number;
  balance: number;
  unique_accounts: number;
  unique_tags: number;
}

export interface TagSummary {
  tag: string;
  total_debit: number;
  total_credit: number;
  count: number;
}

export interface BalanceTrendPoint {
  month: string;
  debit: number;
  credit: number;
  balance: number;
}

interface DashboardState {
  loading: boolean;
  errorMessage: string;
  overview: SummaryOverview | null;
  tagSummary: TagSummary[];
  balanceTrend: BalanceTrendPoint[];
  topExpenses: ITransaction[];
  dateFrom: string;
  dateTo: string;
}

const initialState: DashboardState = {
  loading: false,
  errorMessage: '',
  overview: null,
  tagSummary: [],
  balanceTrend: [],
  topExpenses: [],
  dateFrom: '',
  dateTo: '',
};

function buildFilter(dateFrom: string, dateTo: string): string | undefined {
  const parts: string[] = [];
  if (dateFrom) parts.push(`date >= '${dateFrom}'`);
  if (dateTo) parts.push(`date <= '${dateTo}'`);
  return parts.length > 0 ? parts.join(' and ') : undefined;
}

export const fetchOverview = createAsyncThunk(
  'dashboard/overview',
  async (params: { dateFrom: string; dateTo: string }) => {
    const filter = buildFilter(params.dateFrom, params.dateTo);
    const qp = filter ? `?filter=${encodeURIComponent(filter)}` : '';
    return axios.get<SummaryOverview>(`${summaryUrl}/overview${qp}`);
  },
  { serializeError: serializeAxiosError },
);

export const fetchTagSummary = createAsyncThunk(
  'dashboard/tagSummary',
  async (params: { dateFrom: string; dateTo: string }) => {
    const filter = buildFilter(params.dateFrom, params.dateTo);
    const qp = filter ? `?filter=${encodeURIComponent(filter)}` : '';
    return axios.get<TagSummary[]>(`${summaryUrl}/by-tag${qp}`);
  },
  { serializeError: serializeAxiosError },
);

export const fetchBalanceTrend = createAsyncThunk(
  'dashboard/balanceTrend',
  async (params: { dateFrom: string; dateTo: string }) => {
    const filter = buildFilter(params.dateFrom, params.dateTo);
    const qp = filter ? `?filter=${encodeURIComponent(filter)}` : '';
    return axios.get<BalanceTrendPoint[]>(`${summaryUrl}/balance-trend${qp}`);
  },
  { serializeError: serializeAxiosError },
);

export const fetchTopExpenses = createAsyncThunk(
  'dashboard/topExpenses',
  async (params: { dateFrom: string; dateTo: string }) => {
    const filter = buildFilter(params.dateFrom, params.dateTo);
    const parts = ["kind = 'debit'"];
    if (filter) parts.push(filter);
    const combined = parts.join(' and ');
    const qp = new URLSearchParams();
    qp.append('filter', combined);
    qp.append('sort', 'amount:desc');
    qp.append('limit', '10');
    return axios.get<{ items: ITransaction[]; total: number }>(`${txnUrl}?${qp.toString()}`);
  },
  { serializeError: serializeAxiosError },
);

export const DashboardSlice = createSlice({
  name: 'dashboard',
  initialState,
  reducers: {
    setDateRange(state, action: PayloadAction<{ dateFrom: string; dateTo: string }>) {
      state.dateFrom = action.payload.dateFrom;
      state.dateTo = action.payload.dateTo;
    },
    resetDashboard() {
      return initialState;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(fetchOverview.pending, (state) => { state.loading = true; })
      .addCase(fetchOverview.fulfilled, (state, action) => {
        state.overview = action.payload.data;
        state.loading = false;
      })
      .addCase(fetchOverview.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to load overview';
      })
      .addCase(fetchTagSummary.fulfilled, (state, action) => {
        state.tagSummary = action.payload.data;
      })
      .addCase(fetchBalanceTrend.fulfilled, (state, action) => {
        state.balanceTrend = action.payload.data;
      })
      .addCase(fetchTopExpenses.fulfilled, (state, action) => {
        state.topExpenses = action.payload.data.items;
      });
  },
});

export const { setDateRange, resetDashboard } = DashboardSlice.actions;
export default DashboardSlice.reducer;
