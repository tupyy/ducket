import { ReducersMapObject } from '@reduxjs/toolkit';
import tags from '@app/shared/reducers/tag.reducer';
import rules from '@app/shared/reducers/rule.reducer';
import transactions from '@app/shared/reducers/transaction.reducer';

const rootReducer: ReducersMapObject = {
  tags,
  rules,
  transactions,
};

export default rootReducer;
