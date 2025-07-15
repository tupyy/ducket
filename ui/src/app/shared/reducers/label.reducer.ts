import axios from 'axios';

import { ILabel, ILabelForm, ILabelUpdateForm, ILabels } from '@app/shared/models/label';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '../../shared/reducers/reducer.utils';
import { createAxiosDateTransformer } from 'axios-date-transformer';

const labelApiUrl = 'api/v1/labels';

const initialState = {
  loading: false,
  errorMessage: '',
  creating: false,
  createSuccess: false,
  updating: false,
  updateSuccess: false,
  deleting: false,
  deleteSuccess: false,
  labels: [] as Array<ILabel>,
  totalItems: 0,
};

export const getLabels = createAsyncThunk(
  'labels/get',
  async () => {
    return createAxiosDateTransformer().get<ILabels>(labelApiUrl);
  },
  { serializeError: serializeAxiosError },
);

export const createLabel = createAsyncThunk(
  'labels/create',
  async (label: ILabelForm, thunkAPI) => {
    const result = axios.post<ILabelForm>(labelApiUrl, label).then(() => {
      thunkAPI.dispatch(getLabels());
    });
    return result;
  },
  { serializeError: serializeAxiosError },
);

export const updateLabel = createAsyncThunk(
  'labels/update',
  async (label: ILabelUpdateForm, thunkAPI) => {
    const url = `${labelApiUrl}/${label.id}`;
    const newLabel: ILabelForm = {
      key: label.key,
      value: label.value,
    };
    const result = axios.put<ILabelForm>(url, newLabel);
    thunkAPI.dispatch(getLabels());
    return result;
  },
  { serializeError: serializeAxiosError },
);

export const deleteLabel = createAsyncThunk(
  'labels/delete',
  async (id: string, thunkAPI) => {
    const url = `${labelApiUrl}/${id}`;
    const result = axios.delete<void>(url).then(() => {
      thunkAPI.dispatch(getLabels());
      return result;
    });
    return result;
  },
  { serializeError: serializeAxiosError },
);

export type LabelState = Readonly<typeof initialState>;

export const LabelManagementSlice = createSlice({
  name: 'labels',
  initialState: initialState as LabelState,
  reducers: {
    reset() {
      return initialState;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(getLabels.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(getLabels.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'failed to load labels';
      })
      .addCase(getLabels.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.labels = action.payload.data.labels;
        state.totalItems = action.payload.data.total;
      })
      .addCase(createLabel.pending, (state) => {
        state.creating = true;
        state.errorMessage = '';
        state.createSuccess = false;
      })
      .addCase(updateLabel.pending, (state) => {
        state.updating = true;
        state.errorMessage = '';
        state.updateSuccess = false;
      })
      .addCase(deleteLabel.pending, (state) => {
        state.deleting = true;
        state.deleteSuccess = false;
        state.errorMessage = '';
      })
      .addCase(createLabel.fulfilled, (state) => {
        state.creating = false;
        state.errorMessage = '';
        state.createSuccess = true;
      })
      .addCase(updateLabel.fulfilled, (state) => {
        state.updating = false;
        state.updateSuccess = true;
        state.errorMessage = '';
      })
      .addCase(deleteLabel.fulfilled, (state) => {
        state.deleting = false;
        state.deleteSuccess = true;
        state.errorMessage = '';
      })
      .addCase(createLabel.rejected, (state, action) => {
        state.creating = false;
        state.errorMessage = action.error.message || 'error creating label';
        state.createSuccess = false;
      })
      .addCase(updateLabel.rejected, (state, action) => {
        state.updating = false;
        state.errorMessage = action.error.message || 'error updating label';
        state.updateSuccess = false;
      })
      .addCase(deleteLabel.rejected, (state, action) => {
        state.deleting = false;
        state.errorMessage = action.error.message || 'error deleting label';
        state.deleteSuccess = false;
      });
  },
});

export const { reset } = LabelManagementSlice.actions;
export default LabelManagementSlice.reducer; 