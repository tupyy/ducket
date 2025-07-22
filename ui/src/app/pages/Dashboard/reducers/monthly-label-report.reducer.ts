import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { IMonthlyTagReport, IMonthlyTagSummary } from '@app/shared/models/label';
import { ITransaction } from '@app/shared/models/transaction';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';

const initialState = {
  loading: false,
  errorMessage: '',
  monthlyLabelReports: [] as Array<IMonthlyTagReport>,
  monthlyLabelSummaries: [] as Array<IMonthlyTagSummary>,
};

export const calculateMonthlyLabelReport = createAsyncThunk(
  'monthlyLabelReport/calculate',
  async (params: { transactions: ITransaction[]; excludeCredits?: boolean }) => {
    const { transactions, excludeCredits = false } = params;
    
    // Group transactions by month-year and label
    const monthlyLabelData: { [key: string]: { [label: string]: { amount: number; count: number } } } = {};

    transactions.forEach((transaction: ITransaction) => {
      // Skip credit transactions if excludeCredits is true
      if (excludeCredits && transaction.kind === 'credit') {
        return;
      }

      const transactionDate = new Date(transaction.date);
      const monthYear = `${transactionDate.getFullYear()}-${String(transactionDate.getMonth() + 1).padStart(2, '0')}`;
      
      if (!monthlyLabelData[monthYear]) {
        monthlyLabelData[monthYear] = {};
      }

      transaction.labels.forEach((label) => {
        // Create a combined label identifier with key:value format
        const labelKey = `${label.key}:${label.value}`;
        if (!monthlyLabelData[monthYear][labelKey]) {
          monthlyLabelData[monthYear][labelKey] = { amount: 0, count: 0 };
        }
        monthlyLabelData[monthYear][labelKey].amount += Math.abs(transaction.amount);
        monthlyLabelData[monthYear][labelKey].count += 1;
      });
    });

    // Convert to IMonthlyTagReport format (keeping interface name for backward compatibility)
    const monthlyLabelReports: IMonthlyTagReport[] = [];
    Object.entries(monthlyLabelData).forEach(([monthYear, labelData]) => {
      const [year, month] = monthYear.split('-').map(Number);
      Object.entries(labelData).forEach(([label, data]) => {
        monthlyLabelReports.push({
          tag: label, // Keep as 'tag' for backward compatibility with existing interfaces
          month,
          year,
          amount: data.amount,
          transactionCount: data.count,
        });
      });
    });

    return monthlyLabelReports;
  },
  { serializeError: serializeAxiosError },
);

export const calculateMonthlyLabelSummaries = createAsyncThunk(
  'monthlyLabelReport/calculateSummaries',
  async (monthlyLabelReports: IMonthlyTagReport[]) => {
    // Group by month-year and create summaries
    const summariesMap: { [key: string]: IMonthlyTagSummary } = {};

    monthlyLabelReports.forEach((report) => {
      const monthYear = `${report.year}-${String(report.month).padStart(2, '0')}`;
      
      if (!summariesMap[monthYear]) {
        summariesMap[monthYear] = {
          monthYear,
          tags: [], // Keep as 'tags' for backward compatibility
          totalAmount: 0,
        };
      }

      summariesMap[monthYear].tags.push(report);
      summariesMap[monthYear].totalAmount += report.amount;
    });

    // Sort summaries by month-year (newest first)
    const summaries = Object.values(summariesMap).sort((a, b) => 
      b.monthYear.localeCompare(a.monthYear)
    );

    // Sort tags within each summary by amount (highest first)
    summaries.forEach((summary) => {
      summary.tags.sort((a, b) => b.amount - a.amount);
    });

    return summaries;
  },
  { serializeError: serializeAxiosError },
);

export const getLabelAmountsByMonth = createAsyncThunk(
  'monthlyLabelReport/getLabelAmountsByMonth',
  async (params: { labelName: string; monthlyLabelReports: IMonthlyTagReport[] }) => {
    const { labelName, monthlyLabelReports } = params;
    
    const labelReports = monthlyLabelReports
      .filter((report) => report.tag === labelName)
      .sort((a, b) => {
        // Sort by year, then by month
        if (a.year !== b.year) return a.year - b.year;
        return a.month - b.month;
      });

    return labelReports;
  },
  { serializeError: serializeAxiosError },
);

export type MonthlyLabelReportState = Readonly<typeof initialState>;

export const MonthlyLabelReportSlice = createSlice({
  name: 'monthlyLabelReport',
  initialState: initialState as MonthlyLabelReportState,
  reducers: {
    reset() {
      return initialState;
    },
    clearData(state) {
      state.monthlyLabelReports = [];
      state.monthlyLabelSummaries = [];
      state.errorMessage = '';
    },
  },
  extraReducers(builder) {
    builder
      .addCase(calculateMonthlyLabelReport.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(calculateMonthlyLabelReport.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to calculate monthly label report';
      })
      .addCase(calculateMonthlyLabelReport.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.monthlyLabelReports = action.payload;
      })
      .addCase(calculateMonthlyLabelSummaries.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(calculateMonthlyLabelSummaries.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to calculate monthly label summaries';
      })
      .addCase(calculateMonthlyLabelSummaries.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.monthlyLabelSummaries = action.payload;
      })
      .addCase(getLabelAmountsByMonth.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(getLabelAmountsByMonth.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to get label amounts by month';
      })
      .addCase(getLabelAmountsByMonth.fulfilled, (state) => {
        state.loading = false;
        state.errorMessage = '';
        // This action doesn't modify state, it just returns filtered data
      });
  },
});

export const { reset, clearData } = MonthlyLabelReportSlice.actions;
export default MonthlyLabelReportSlice.reducer; 