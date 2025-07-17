import axios from 'axios';

import { IRule, IRuleForm, IRules, IUpdateRuleForm } from '@app/shared/models/rule';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '../../shared/reducers/reducer.utils';
import { createAxiosDateTransformer } from 'axios-date-transformer';

const ruleApiUrl = 'api/v1/rules';

const initialState = {
  loading: false,
  errorMessage: '',
  creating: false,
  createSuccess: false,
  updating: false,
  updateSuccess: false,
  deleting: false,
  deleteSuccess: false,
  syncing: false,
  syncSuccess: false,
  rules: [] as Array<IRule>,
  totalItems: 0,
};

export const getRules = createAsyncThunk(
  'rules/get',
  async () => {
    return createAxiosDateTransformer().get<IRules>(ruleApiUrl);
  },
  { serializeError: serializeAxiosError },
);

export const createRule = createAsyncThunk(
  'rules/create',
  async (rule: IRuleForm, thunkAPI) => {
    const result = axios.post<IRuleForm>(ruleApiUrl, rule).then(() => {
      thunkAPI.dispatch(getRules());
    });
    return result;
  },
  { serializeError: serializeAxiosError },
);

export const updateRule = createAsyncThunk(
  'rules/update',
  async (rule: IRuleForm, thunkAPI) => {
    const url = `${ruleApiUrl}/${rule.name}`;
    const newRule: IUpdateRuleForm = {
      pattern: rule.pattern,
      labels: rule.labels,
    };
    const result = axios.put<IUpdateRuleForm>(url, newRule).then(() => thunkAPI.dispatch(getRules()));
    return result;
  },
  { serializeError: serializeAxiosError },
);

export const deleteRule = createAsyncThunk(
  'rules/delete',
  async (name: string, thunkAPI) => {
    const url = `${ruleApiUrl}/${name}`;
    const result = axios.delete<void>(url).then(() => thunkAPI.dispatch(getRules()));
    return result;
  },
  { serializeError: serializeAxiosError },
);

export const syncRule = createAsyncThunk(
  'rules/sync',
  async (ruleName: string, thunkAPI) => {
    // For now, we'll use a placeholder endpoint that re-applies the rule to transactions
    // In the future, this could be a dedicated sync endpoint
    const url = `${ruleApiUrl}/${ruleName}/process`;
    const result = axios.post<void>(url).then(() => thunkAPI.dispatch(getRules()));
    return result;
  },
  { serializeError: serializeAxiosError },
);

export type RuleState = Readonly<typeof initialState>;

export const RuleManagementSlice = createSlice({
  name: 'rules',
  initialState: initialState as RuleState,
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
        state.errorMessage = action.error.message || 'failed to load rules';
      })
      .addCase(getRules.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.rules = action.payload.data.rules;
        state.totalItems = action.payload.data.total;
      })
      .addCase(createRule.pending, (state) => {
        state.creating = true;
        state.errorMessage = '';
        state.createSuccess = false;
      })
      .addCase(updateRule.pending, (state) => {
        state.updating = true;
        state.errorMessage = '';
        state.updateSuccess = false;
      })
      .addCase(deleteRule.pending, (state) => {
        state.deleting = true;
        state.deleteSuccess = false;
        state.errorMessage = '';
      })
      .addCase(syncRule.pending, (state) => {
        state.syncing = true;
        state.syncSuccess = false;
        state.errorMessage = '';
      })
      .addCase(createRule.fulfilled, (state) => {
        state.creating = false;
        state.errorMessage = '';
        state.createSuccess = true;
      })
      .addCase(updateRule.fulfilled, (state) => {
        state.updating = false;
        state.updateSuccess = true;
        state.errorMessage = '';
      })
      .addCase(deleteRule.fulfilled, (state) => {
        state.deleting = false;
        state.deleteSuccess = true;
        state.errorMessage = '';
      })
      .addCase(syncRule.fulfilled, (state) => {
        state.syncing = false;
        state.syncSuccess = true;
        state.errorMessage = '';
      })
      .addCase(createRule.rejected, (state, action) => {
        state.creating = false;
        state.errorMessage = action.error.message || 'error creating rule';
        state.createSuccess = false;
      })
      .addCase(updateRule.rejected, (state, action) => {
        state.updating = false;
        state.errorMessage = action.error.message || 'error updating rule';
        state.updateSuccess = false;
      })
      .addCase(deleteRule.rejected, (state, action) => {
        state.deleting = false;
        state.errorMessage = action.error.message || 'error deleting rule';
        state.deleteSuccess = false;
      })
      .addCase(syncRule.rejected, (state, action) => {
        state.syncing = false;
        state.errorMessage = action.error.message || 'error syncing rule';
        state.syncSuccess = false;
      });
  },
});

export const { reset } = RuleManagementSlice.actions;
export default RuleManagementSlice.reducer;
