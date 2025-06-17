import { ReducersMapObject } from '@reduxjs/toolkit';
import tags from '@app/shared/reducers/tag.reducer';
import rules from '@app/shared/reducers/rule.reducer';

const rootReducer: ReducersMapObject = {
  tags,
  rules,
};

export default rootReducer;
