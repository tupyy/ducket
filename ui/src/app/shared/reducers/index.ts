import { ReducersMapObject } from '@reduxjs/toolkit';
import tags from '@app/shared/reducers/tag.reducer';

const rootReducer: ReducersMapObject = {
  tags,
};

export default rootReducer;
