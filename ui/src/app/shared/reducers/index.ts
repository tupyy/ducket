import { ReducersMapObject } from '@reduxjs/toolkit';
import tags from '@app/shared/reducers/tag.reducer';
import rules from '@app/shared/reducers/rule.reducer';
import transactions from '@app/shared/reducers/transaction.reducer';
import transactionFilter from '@app/shared/reducers/transaction-filter.reducer';
import tagReport from '@app/shared/reducers/tag-report.reducer';
import monthlyTagReport from '@app/shared/reducers/monthly-tag-report.reducer';
import importReducer from '@app/shared/reducers/import.reducer';

const rootReducer: ReducersMapObject = {
  tags,
  rules,
  transactions,
  transactionFilter,
  tagReport,
  monthlyTagReport,
  import: importReducer,
};

export default rootReducer;
