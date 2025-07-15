import { ReducersMapObject } from '@reduxjs/toolkit';
import labels from '@app/shared/reducers/label.reducer';
import rules from '@app/shared/reducers/rule.reducer';
import transactions from '@app/shared/reducers/transaction.reducer';
import transactionFilter from '@app/shared/reducers/transaction-filter.reducer';
import labelReport from '@app/shared/reducers/label-report.reducer';
import monthlyLabelReport from '@app/shared/reducers/monthly-label-report.reducer';
import importReducer from '@app/shared/reducers/import.reducer';

const rootReducer: ReducersMapObject = {
  labels,
  rules,
  transactions,
  transactionFilter,
  labelReport,
  monthlyLabelReport,
  import: importReducer,
};

export default rootReducer;
