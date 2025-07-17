import axios from 'axios';

import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';
import { createAxiosDateTransformer } from 'axios-date-transformer';
import { ITransaction, ITransactionForm, ITransactionUpdateForm, ITransactions } from '@app/shared/models/transaction';

const transactionApiUrl = 'api/v1/transactions';

const initialState = {
  loading: false,
  errorMessage: '',
  creating: false,
  createSuccess: false,
  updating: false,
  updateSuccess: false,
  deleting: false,
  deleteSuccess: false,
  transactions: [] as Array<ITransaction>,
  totalItems: 0,
};

export const getTransactions = createAsyncThunk(
  'transactions/get',
  async (params?: { startDate?: string; endDate?: string; labels?: string[] }) => {
    let url = transactionApiUrl;

    if (params?.startDate || params?.endDate || params?.labels?.length) {
      const searchParams = new URLSearchParams();
      if (params.startDate) {
        const startTimestamp = new Date(params.startDate).getTime().toString();
        searchParams.append('startDate', startTimestamp);
      }
      if (params.endDate) {
        const endTimestamp = new Date(params.endDate).getTime().toString();
        searchParams.append('endDate', endTimestamp);
      }
      if (params.labels && params.labels.length > 0) {
        params.labels.forEach((label) => searchParams.append('labels', label));
      }
      url = `${transactionApiUrl}?${searchParams.toString()}`;
    }

    return axios.get<ITransactions>(url);
  },
  { serializeError: serializeAxiosError },
);

export const createTransaction = createAsyncThunk(
  'transactions/create',
  async (transaction: ITransactionForm, thunkAPI) => {
    const result = axios.post<ITransactionForm>(transactionApiUrl, transaction).then(() => {
      thunkAPI.dispatch(getTransactions());
    });
    return result;
  },
  { serializeError: serializeAxiosError },
);

export const updateTransaction = createAsyncThunk(
  'transactions/update',
  async (form: ITransactionUpdateForm, thunkAPI) => {
    const url = `${transactionApiUrl}/${form.name}`;
    const newTransaction: ITransactionForm = {
      content: form.content,
      kind: form.kind,
      date: form.date,
      amount: form.amount,
      labels: form.labels,
    };
    const result = axios.put<ITransactionForm>(url, newTransaction).then(() => thunkAPI.dispatch(getTransactions()));
    return result;
  },
  { serializeError: serializeAxiosError },
);

export const deleteTransaction = createAsyncThunk(
  'transactions/delete',
  async (name: string, thunkAPI) => {
    const url = `${transactionApiUrl}/${name}`;
    const result = axios.delete<void>(url).then(() => thunkAPI.dispatch(getTransactions()));
    return result;
  },
  { serializeError: serializeAxiosError },
);

export type transactionState = Readonly<typeof initialState>;

export const TransactionManagementSlice = createSlice({
  name: 'transactions',
  initialState: initialState as transactionState,
  reducers: {
    reset() {
      return initialState;
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
        state.errorMessage = action.error.message || 'failed to load transactions';
      })
      .addCase(getTransactions.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.transactions = action.payload.data.items;
        state.totalItems = action.payload.data.total;
      })
      .addCase(createTransaction.pending, (state) => {
        state.creating = true;
        state.errorMessage = '';
        state.createSuccess = false;
      })
      .addCase(updateTransaction.pending, (state) => {
        state.updating = true;
        state.errorMessage = '';
        state.updateSuccess = false;
      })
      .addCase(deleteTransaction.pending, (state) => {
        state.deleting = true;
        state.deleteSuccess = false;
        state.errorMessage = '';
      })
      .addCase(createTransaction.fulfilled, (state) => {
        state.creating = false;
        state.errorMessage = '';
        state.createSuccess = true;
      })
      .addCase(updateTransaction.fulfilled, (state) => {
        state.updating = false;
        state.updateSuccess = true;
        state.errorMessage = '';
      })
      .addCase(deleteTransaction.fulfilled, (state) => {
        state.deleting = false;
        state.deleteSuccess = true;
        state.errorMessage = '';
      })
      .addCase(createTransaction.rejected, (state, action) => {
        state.creating = false;
        state.errorMessage = action.error.message || 'error creating transaction';
        state.createSuccess = false;
      })
      .addCase(updateTransaction.rejected, (state, action) => {
        state.updating = false;
        state.errorMessage = action.error.message || 'error updating transaction';
        state.updateSuccess = false;
      })
      .addCase(deleteTransaction.rejected, (state, action) => {
        state.deleting = false;
        state.errorMessage = action.error.message || 'error deleting transaction';
        state.deleteSuccess = false;
      });
  },
});

export const { reset } = TransactionManagementSlice.actions;
export default TransactionManagementSlice.reducer;
