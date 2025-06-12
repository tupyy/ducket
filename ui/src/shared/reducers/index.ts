import { ReducersMapObject } from "@reduxjs/toolkit";
import tags from '../../modules/tag/tag.reducer';

const rootReducer: ReducersMapObject = {
    tags,
};

export default rootReducer;
