import axios from 'axios';

import { ITag, ITags } from '../../shared/models/tag';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '../../shared/reducers/reducer.utils';

const apiUrl = "api/v1/tags";

const initialState = {
    loading: false,
    errorMessage: null,
    updateSuccess: false,
    failure: false,
    tags: [] as Array<ITag>,
    totalItems: 0,
}

export const getTags = createAsyncThunk('tags/get', async () => {
    return axios.get<ITags>(apiUrl);
},
{ serializeError: serializeAxiosError});

export type TagState = Readonly<typeof initialState>;

export const TagManagementSlice = createSlice({
    name: 'tags',
    initialState: initialState as TagState,
    reducers: {
        reset() {
            return initialState;
        },
    },
    extraReducers(builder) {
        builder.
            addCase(getTags.pending, state => {
                state.loading = true;
                state.errorMessage = null;
        })
        .addCase(getTags.rejected, state => {
            state.loading = false;
            state.failure = true;
        })
        .addCase(getTags.fulfilled, (state, action) => {
            state.loading = false;
            state.failure = false;
            state.tags = action.payload.data.tags;
            state.totalItems = action.payload.data.total;
        })
    },
});

export const {reset} = TagManagementSlice.actions;
export default TagManagementSlice.reducer;
    
    
