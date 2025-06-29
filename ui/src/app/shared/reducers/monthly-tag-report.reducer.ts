import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { IMonthlyTagReport, IMonthlyTagSummary } from '@app/shared/models/tag';
import { ITransaction } from '@app/shared/models/transaction';
import { serializeAxiosError } from './reducer.utils';

const initialState = {
  loading: false,
  errorMessage: '',
  monthlyTagReports: [] as Array<IMonthlyTagReport>,
  monthlyTagSummaries: [] as Array<IMonthlyTagSummary>,
};

export const calculateMonthlyTagReport = createAsyncThunk(
  'monthlyTagReport/calculate',
  async (params: { transactions: ITransaction[]; excludeCredits?: boolean }) => {
    const { transactions, excludeCredits = false } = params;
    
    // Group transactions by month-year and tag
    const monthlyTagData: { [key: string]: { [tag: string]: { amount: number; count: number } } } = {};

    transactions.forEach((transaction: ITransaction) => {
      // Skip credit transactions if excludeCredits is true
      if (excludeCredits && transaction.kind === 'credit') {
        return;
      }

      const transactionDate = new Date(transaction.date);
      const monthYear = `${transactionDate.getFullYear()}-${String(transactionDate.getMonth() + 1).padStart(2, '0')}`;
      
      if (!monthlyTagData[monthYear]) {
        monthlyTagData[monthYear] = {};
      }

      transaction.tags.forEach((tag) => {
        const tagValue = tag.value;
        if (!monthlyTagData[monthYear][tagValue]) {
          monthlyTagData[monthYear][tagValue] = { amount: 0, count: 0 };
        }
        monthlyTagData[monthYear][tagValue].amount += Math.abs(transaction.amount);
        monthlyTagData[monthYear][tagValue].count += 1;
      });
    });

    // Convert to IMonthlyTagReport format
    const monthlyTagReports: IMonthlyTagReport[] = [];
    Object.entries(monthlyTagData).forEach(([monthYear, tagData]) => {
      const [year, month] = monthYear.split('-').map(Number);
      Object.entries(tagData).forEach(([tag, data]) => {
        monthlyTagReports.push({
          tag,
          month,
          year,
          amount: data.amount,
          transactionCount: data.count,
        });
      });
    });

    return monthlyTagReports;
  },
  { serializeError: serializeAxiosError },
);

export const calculateMonthlyTagSummaries = createAsyncThunk(
  'monthlyTagReport/calculateSummaries',
  async (monthlyTagReports: IMonthlyTagReport[]) => {
    // Group by month-year and create summaries
    const summariesMap: { [key: string]: IMonthlyTagSummary } = {};

    monthlyTagReports.forEach((report) => {
      const monthYear = `${report.year}-${String(report.month).padStart(2, '0')}`;
      
      if (!summariesMap[monthYear]) {
        summariesMap[monthYear] = {
          monthYear,
          tags: [],
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

export const getTagAmountsByMonth = createAsyncThunk(
  'monthlyTagReport/getTagAmountsByMonth',
  async (params: { tagName: string; monthlyTagReports: IMonthlyTagReport[] }) => {
    const { tagName, monthlyTagReports } = params;
    
    const tagReports = monthlyTagReports
      .filter((report) => report.tag === tagName)
      .sort((a, b) => {
        // Sort by year, then by month
        if (a.year !== b.year) return a.year - b.year;
        return a.month - b.month;
      });

    return tagReports;
  },
  { serializeError: serializeAxiosError },
);

export type MonthlyTagReportState = Readonly<typeof initialState>;

export const MonthlyTagReportSlice = createSlice({
  name: 'monthlyTagReport',
  initialState: initialState as MonthlyTagReportState,
  reducers: {
    reset() {
      return initialState;
    },
    clearData(state) {
      state.monthlyTagReports = [];
      state.monthlyTagSummaries = [];
      state.errorMessage = '';
    },
  },
  extraReducers(builder) {
    builder
      .addCase(calculateMonthlyTagReport.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(calculateMonthlyTagReport.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to calculate monthly tag report';
      })
      .addCase(calculateMonthlyTagReport.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.monthlyTagReports = action.payload;
      })
      .addCase(calculateMonthlyTagSummaries.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(calculateMonthlyTagSummaries.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to calculate monthly tag summaries';
      })
      .addCase(calculateMonthlyTagSummaries.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.monthlyTagSummaries = action.payload;
      })
      .addCase(getTagAmountsByMonth.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(getTagAmountsByMonth.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to get tag amounts by month';
      })
      .addCase(getTagAmountsByMonth.fulfilled, (state) => {
        state.loading = false;
        state.errorMessage = '';
        // This action doesn't modify state, it just returns filtered data
      });
  },
});

export const { reset, clearData } = MonthlyTagReportSlice.actions;
export default MonthlyTagReportSlice.reducer; 