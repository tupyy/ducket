import axios from 'axios';
import { IRule } from '@app/shared/models/rule';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';

const apiUrl = 'api/v1/rules';

interface RuleState {
  loading: boolean;
  errorMessage: string;
  rules: IRule[];
  creating: boolean;
  updating: boolean;
}

const initialState: RuleState = {
  loading: false,
  errorMessage: '',
  rules: [],
  creating: false,
  updating: false,
};

export const getRules = createAsyncThunk(
  'rules/get',
  async () => axios.get<IRule[]>(apiUrl),
  { serializeError: serializeAxiosError },
);

export const createRule = createAsyncThunk(
  'rules/create',
  async (rule: { name: string; filter: string; tags: string[] }, thunkAPI) => {
    await axios.post(apiUrl, rule);
    thunkAPI.dispatch(getRules());
  },
  { serializeError: serializeAxiosError },
);

export const updateRule = createAsyncThunk(
  'rules/update',
  async (rule: { id: number; name: string; filter: string; tags: string[] }, thunkAPI) => {
    await axios.put(`${apiUrl}/${rule.id}`, rule);
    thunkAPI.dispatch(getRules());
  },
  { serializeError: serializeAxiosError },
);

export const deleteRule = createAsyncThunk(
  'rules/delete',
  async (id: number, thunkAPI) => {
    await axios.delete(`${apiUrl}/${id}`);
    thunkAPI.dispatch(getRules());
  },
  { serializeError: serializeAxiosError },
);

export const RuleSlice = createSlice({
  name: 'rules',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(getRules.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(getRules.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to load rules';
      })
      .addCase(getRules.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.rules = action.payload.data;
      })
      .addCase(createRule.pending, (state) => {
        state.creating = true;
        state.errorMessage = '';
      })
      .addCase(createRule.fulfilled, (state) => {
        state.creating = false;
      })
      .addCase(createRule.rejected, (state, action) => {
        state.creating = false;
        state.errorMessage = action.error.message || 'Failed to create rule';
      })
      .addCase(updateRule.pending, (state) => {
        state.updating = true;
        state.errorMessage = '';
      })
      .addCase(updateRule.fulfilled, (state) => {
        state.updating = false;
      })
      .addCase(updateRule.rejected, (state, action) => {
        state.updating = false;
        state.errorMessage = action.error.message || 'Failed to update rule';
      });
  },
});

export const { reset } = RuleSlice.actions;
export default RuleSlice.reducer;
