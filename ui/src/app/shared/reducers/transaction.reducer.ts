import axios from 'axios';

import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';
import { ITransaction, ITransactionForm, ITransactionUpdateForm, ITransactions } from '@app/shared/models/transaction';
import { applyFilters } from '@app/pages/Transactions/reducers/transaction-filter.reducer';
import { createAxiosDateTransformer } from 'axios-date-transformer';

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
  addingLabel: false,
  addLabelSuccess: false,
  updatingInfo: false,
  updateInfoSuccess: false,
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

    return createAxiosDateTransformer().get<ITransactions>(url);
  },
  { serializeError: serializeAxiosError }
);

export const createTransaction = createAsyncThunk(
  'transactions/create',
  async (transaction: ITransactionForm, thunkAPI) => {
    const result = axios.post<ITransactionForm>(transactionApiUrl, transaction).then(() => {
      thunkAPI.dispatch(getTransactions());
    });
    return result;
  },
  { serializeError: serializeAxiosError }
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
  { serializeError: serializeAxiosError }
);

export const addLabelToTransaction = createAsyncThunk(
  'transactions/addLabel',
  async (params: { transactionHref: string; key: string; value: string }, thunkAPI) => {
    // Extract transaction ID from href (assuming href is like "/api/v1/transactions/123")
    const transactionId = params.transactionHref.split('/').pop();
    const url = `${transactionApiUrl}/${transactionId}/labels`;

    const labelData = {
      key: params.key,
      value: params.value,
    };

    const result = await axios.post(url, labelData);

    // Refresh transactions after adding label
    await thunkAPI.dispatch(getTransactions());

    // Recalculate filters after transactions are refreshed
    const state = thunkAPI.getState() as any;
    const filterState = state.transactionFilter;

    await thunkAPI.dispatch(
      applyFilters({
        selectedLabels: filterState.selectedLabels,
        selectedTransactionTypes: filterState.selectedTransactionTypes,
        selectedAccounts: filterState.selectedAccounts,
        descriptionFilter: filterState.descriptionFilter,
        showOnlyUnlabeled: filterState.showOnlyUnlabeled,
      })
    );

    return result;
  },
  { serializeError: serializeAxiosError }
);

export const removeLabelFromTransaction = createAsyncThunk(
  'transactions/removeLabel',
  async (params: { transactionHref: string; key: string; value: string }, thunkAPI) => {
    // Extract transaction ID from href (assuming href is like "/api/v1/transactions/123")
    const transactionId = params.transactionHref.split('/').pop();
    const url = `${transactionApiUrl}/${transactionId}/labels`;

    const deleteData = {
      key: params.key,
      value: params.value,
    };

    const result = await axios.delete(url, { data: deleteData });

    // Refresh transactions after removing label
    await thunkAPI.dispatch(getTransactions());

    // Recalculate filters after transactions are refreshed
    const state = thunkAPI.getState() as any;
    const filterState = state.transactionFilter;

    await thunkAPI.dispatch(
      applyFilters({
        selectedLabels: filterState.selectedLabels,
        selectedTransactionTypes: filterState.selectedTransactionTypes,
        selectedAccounts: filterState.selectedAccounts,
        descriptionFilter: filterState.descriptionFilter,
        showOnlyUnlabeled: filterState.showOnlyUnlabeled,
      })
    );

    return result;
  },
  { serializeError: serializeAxiosError }
);

export const updateTransactionInfo = createAsyncThunk(
  'transactions/updateInfo',
  async (params: { transactionHref: string; info: string }, thunkAPI) => {
    // Extract transaction ID from href (assuming href is like "/api/v1/transactions/123")
    const transactionId = params.transactionHref.split('/').pop();
    const url = `${transactionApiUrl}/${transactionId}`;

    const updateData = {
      info: params.info,
    };

    const result = await axios.patch(url, updateData);

    // Refresh transactions after updating info
    await thunkAPI.dispatch(getTransactions());

    // Recalculate filters after transactions are refreshed
    const state = thunkAPI.getState() as any;
    const filterState = state.transactionFilter;

    await thunkAPI.dispatch(
      applyFilters({
        selectedLabels: filterState.selectedLabels,
        selectedTransactionTypes: filterState.selectedTransactionTypes,
        selectedAccounts: filterState.selectedAccounts,
        descriptionFilter: filterState.descriptionFilter,
        showOnlyUnlabeled: filterState.showOnlyUnlabeled,
      })
    );

    return result;
  },
  { serializeError: serializeAxiosError }
);

export const deleteTransaction = createAsyncThunk(
  'transactions/delete',
  async (name: string, thunkAPI) => {
    const url = `${transactionApiUrl}/${name}`;
    const result = axios.delete<void>(url).then(() => thunkAPI.dispatch(getTransactions()));
    return result;
  },
  { serializeError: serializeAxiosError }
);

export type transactionState = Readonly<typeof initialState>;

export const TransactionManagementSlice = createSlice({
  name: 'transactions',
  initialState: initialState as transactionState,
  reducers: {
    reset() {
      return initialState;
    },
    clearAddLabelToTransactionSuccess(state) {
      if (state.addLabelSuccess) {
        state.addLabelSuccess = false;
      }
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
      .addCase(addLabelToTransaction.pending, (state) => {
        state.addingLabel = true;
        state.addLabelSuccess = false;
        state.errorMessage = '';
      })
      .addCase(updateTransactionInfo.pending, (state) => {
        state.updatingInfo = true;
        state.updateInfoSuccess = false;
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
      .addCase(addLabelToTransaction.fulfilled, (state) => {
        state.addingLabel = false;
        state.addLabelSuccess = true;
        state.errorMessage = '';
      })
      .addCase(updateTransactionInfo.fulfilled, (state) => {
        state.updatingInfo = false;
        state.updateInfoSuccess = true;
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
      })
      .addCase(addLabelToTransaction.rejected, (state, action) => {
        state.addingLabel = false;
        const errorPayload = action.payload as any;
        state.errorMessage = errorPayload?.response?.data?.message || 
                           errorPayload?.message || 
                           action.error.message || 
                           'error adding label to transaction';
        state.addLabelSuccess = false;
      })
      .addCase(updateTransactionInfo.rejected, (state, action) => {
        state.updatingInfo = false;
        const errorPayload = action.payload as any;
        state.errorMessage = errorPayload?.response?.data?.message || 
                           errorPayload?.message || 
                           action.error.message || 
                           'error updating transaction info';
        state.updateInfoSuccess = false;
      });
  },
});

export const { clearAddLabelToTransactionSuccess, reset } = TransactionManagementSlice.actions;
export default TransactionManagementSlice.reducer;
