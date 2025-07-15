import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { ITagReport, ITransactionTypeReport, IAccountTransactionTypeReport } from '@app/shared/models/label';
import { ITransaction } from '@app/shared/models/transaction';
import { serializeAxiosError } from './reducer.utils';

const initialState = {
  loading: false,
  errorMessage: '',
  labelReportData: [] as Array<ITagReport>,
  transactionTypeData: [] as Array<ITransactionTypeReport>,
  accountTransactionTypeData: [] as Array<IAccountTransactionTypeReport>,
};

export const calculateLabelReport = createAsyncThunk(
  'labelReport/calculate',
  async (params: { transactions: ITransaction[]; excludeCredits?: boolean }) => {
    const { transactions, excludeCredits = false } = params;
    const labelAmounts: { [key: string]: number } = {};

    transactions.forEach((transaction: ITransaction) => {
      // Skip credit transactions if excludeCredits is true
      if (excludeCredits && transaction.kind === 'credit') {
        return;
      }

      transaction.labels.forEach((label) => {
        // Create a combined label identifier with key:value format
        const labelKey = `${label.key}:${label.value}`;
        if (!labelAmounts[labelKey]) {
          labelAmounts[labelKey] = 0;
        }
        labelAmounts[labelKey] += Math.abs(transaction.amount);
      });
    });

    // Convert to LabelReport format
    const labelReportData: ITagReport[] = Object.entries(labelAmounts).map(([label, amount]) => ({
      tag: label, // Keep as 'tag' for backward compatibility with existing interfaces
      amount,
    }));

    return labelReportData;
  },
  { serializeError: serializeAxiosError },
);

export const calculateTransactionTypeReport = createAsyncThunk(
  'labelReport/calculateTransactionType',
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

export const calculateAccountTransactionTypeReport = createAsyncThunk(
  'labelReport/calculateAccountTransactionType',
  async (transactions: ITransaction[]) => {
    const accountData: { [key: number]: { debit: number; credit: number } } = {};

    transactions.forEach((transaction: ITransaction) => {
      const account = transaction.account;
      if (!accountData[account]) {
        accountData[account] = { debit: 0, credit: 0 };
      }

      if (transaction.kind === 'debit') {
        accountData[account].debit += transaction.amount;
      } else if (transaction.kind === 'credit') {
        accountData[account].credit += transaction.amount;
      }
    });

    const result: IAccountTransactionTypeReport[] = [];
    Object.entries(accountData).forEach(([account, amounts]) => {
      const accountNumber = parseInt(account, 10);
      if (amounts.debit > 0) {
        result.push({ account: accountNumber, type: 'debit', amount: amounts.debit });
      }
      if (amounts.credit > 0) {
        result.push({ account: accountNumber, type: 'credit', amount: amounts.credit });
      }
    });

    // Sort by account number for consistent display
    result.sort((a, b) => a.account - b.account || a.type.localeCompare(b.type));

    return result;
  },
  { serializeError: serializeAxiosError },
);

export type LabelReportState = Readonly<typeof initialState>;

export const LabelReportSlice = createSlice({
  name: 'labelReport',
  initialState: initialState as LabelReportState,
  reducers: {
    reset() {
      return initialState;
    },
    clearData(state) {
      state.labelReportData = [];
      state.transactionTypeData = [];
      state.accountTransactionTypeData = [];
      state.errorMessage = '';
    },
  },
  extraReducers(builder) {
    builder
      .addCase(calculateLabelReport.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(calculateLabelReport.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to calculate label report';
      })
      .addCase(calculateLabelReport.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.labelReportData = action.payload;
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
      })
      .addCase(calculateAccountTransactionTypeReport.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(calculateAccountTransactionTypeReport.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to calculate account transaction type report';
      })
      .addCase(calculateAccountTransactionTypeReport.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.accountTransactionTypeData = action.payload;
      });
  },
});

export const { reset, clearData } = LabelReportSlice.actions;
export default LabelReportSlice.reducer; 