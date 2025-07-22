// Dashboard-related reducers that don't make backend requests
export { default as labelReportReducer } from './label-report.reducer';
export { default as monthlyLabelReportReducer } from './monthly-label-report.reducer';

// Export specific actions to avoid conflicts
export {
  calculateLabelReport,
  calculateTransactionTypeReport,
  calculateAccountTransactionTypeReport,
  reset as resetLabelReport,
} from './label-report.reducer';

export {
  calculateMonthlyLabelReport,
  calculateMonthlyLabelSummaries,
  getLabelAmountsByMonth,
  reset as resetMonthlyLabelReport,
} from './monthly-label-report.reducer'; 