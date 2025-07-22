import { ReducersMapObject } from '@reduxjs/toolkit';
import labels from '@app/shared/reducers/label.reducer';
import rules from '@app/shared/reducers/rule.reducer';
import transactions from '@app/shared/reducers/transaction.reducer';
import transactionFilter from '@app/pages/Transactions/reducers/transaction-filter.reducer';
import labelReport from '@app/pages/Dashboard/reducers/label-report.reducer';
import monthlyLabelReport from '@app/pages/Dashboard/reducers/monthly-label-report.reducer';
import importReducer from '@app/shared/reducers/import.reducer';
import transactionSummary from '@app/pages/Transactions/reducers/transactionSummary.reducer';

const rootReducer: ReducersMapObject = {
  labels,
  rules,
  transactions,
  transactionFilter,
  labelReport,
  monthlyLabelReport,
  import: importReducer,
  transactionSummary,
};

export default rootReducer;
