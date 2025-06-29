import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { ITagReport, ITransactionTypeReport } from '@app/shared/models/tag';
import { ITransaction } from '@app/shared/models/transaction';
import { serializeAxiosError } from './reducer.utils';

const initialState = {
  loading: false,
  errorMessage: '',
  tagReportData: [] as Array<ITagReport>,
  transactionTypeData: [] as Array<ITransactionTypeReport>,
};

export const calculateTagReport = createAsyncThunk(
  'tagReport/calculate',
  async (params: { transactions: ITransaction[]; excludeCredits?: boolean }) => {
    const { transactions, excludeCredits = false } = params;
    const tagAmounts: { [key: string]: number } = {};

    transactions.forEach((transaction: ITransaction) => {
      // Skip credit transactions if excludeCredits is true
      if (excludeCredits && transaction.kind === 'credit') {
        return;
      }

      transaction.tags.forEach((tag) => {
        const tagValue = tag.value;
        if (!tagAmounts[tagValue]) {
          tagAmounts[tagValue] = 0;
        }
        tagAmounts[tagValue] += Math.abs(transaction.amount);
      });
    });

    // Convert to TagReport format
    const tagReportData: ITagReport[] = Object.entries(tagAmounts).map(([tag, amount]) => ({
      tag,
      amount,
    }));

    return tagReportData;
  },
  { serializeError: serializeAxiosError },
);

export const calculateTransactionTypeReport = createAsyncThunk(
  'tagReport/calculateTransactionType',
  async (transactions: ITransaction[]) => {
    const debitData = { amount: 0, count: 0 };
    const creditData = { amount: 0, count: 0 };

    transactions.forEach((transaction: ITransaction) => {
      if (transaction.kind === 'debit') {
        debitData.amount += transaction.amount;
        debitData.count += 1;
      } else if (transaction.kind === 'credit') {
        creditData.amount += transaction.amount;
        creditData.count += 1;
      }
    });

    const result: ITransactionTypeReport[] = [];
    if (debitData.count > 0) {
      result.push({ type: 'debit', amount: debitData.amount });
    }
    if (creditData.count > 0) {
      result.push({ type: 'credit', amount: creditData.amount });
    }

    return result;
  },
  { serializeError: serializeAxiosError },
);

export type TagReportState = Readonly<typeof initialState>;

export const TagReportSlice = createSlice({
  name: 'tagReport',
  initialState: initialState as TagReportState,
  reducers: {
    reset() {
      return initialState;
    },
    clearData(state) {
      state.tagReportData = [];
      state.transactionTypeData = [];
      state.errorMessage = '';
    },
  },
  extraReducers(builder) {
    builder
      .addCase(calculateTagReport.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(calculateTagReport.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to calculate tag report';
      })
      .addCase(calculateTagReport.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.tagReportData = action.payload;
      })
      .addCase(calculateTransactionTypeReport.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(calculateTransactionTypeReport.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to calculate transaction type report';
      })
      .addCase(calculateTransactionTypeReport.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.transactionTypeData = action.payload;
      });
  },
});

export const { reset, clearData } = TagReportSlice.actions;
export default TagReportSlice.reducer;
