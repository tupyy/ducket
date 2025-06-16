import axios from 'axios';

import { ITag, ITagForm, ITagUpdateForm, ITags } from '@app/shared/models/tag';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '../../shared/reducers/reducer.utils';
import { createAxiosDateTransformer } from 'axios-date-transformer';

const tagApiUrl = 'api/v1/tags';

const initialState = {
  loading: false,
  errorMessage: '',
  creating: false,
  createSuccess: false,
  updating: false,
  updateSuccess: false,
  deleting: false,
  deleteSuccess: false,
  tags: [] as Array<ITag>,
  totalItems: 0,
};

export const getTags = createAsyncThunk(
  'tags/get',
  async () => {
    return createAxiosDateTransformer().get<ITags>(tagApiUrl);
  },
  { serializeError: serializeAxiosError }
);

export const createTag = createAsyncThunk(
  'tags/create',
  async (tag: ITagForm, thunkAPI) => {
    const result = axios.post<ITagForm>(tagApiUrl, tag).then(() => {
      thunkAPI.dispatch(getTags());
    });
    return result;
  },
  { serializeError: serializeAxiosError }
);

export const updateTag = createAsyncThunk(
  'tags/update',
  async (tag: ITagUpdateForm, thunkAPI) => {
    const url = `${tagApiUrl}/${tag.id}`;
    const newTag: ITagForm = {
      value: tag.value,
    };
    const result = axios.put<ITagForm>(url, newTag);
    thunkAPI.dispatch(getTags());
    return result;
  },
  { serializeError: serializeAxiosError }
);

export const deleteTag = createAsyncThunk(
  'tags/delete',
  async (name: string, thunkAPI) => {
    const url = `${tagApiUrl}/${name}`;
    const result = axios.delete<void>(url).then(() => {
      thunkAPI.dispatch(getTags());
      return result;
    });
    return result;
  },
  { serializeError: serializeAxiosError }
);

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
    builder
      .addCase(getTags.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
      })
      .addCase(getTags.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'failed to load tags';
      })
      .addCase(getTags.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.tags = action.payload.data.tags;
        state.totalItems = action.payload.data.total;
      })
      .addCase(createTag.pending, (state) => {
        state.creating = true;
        state.errorMessage = '';
        state.createSuccess = false;
      })
      .addCase(updateTag.pending, (state) => {
        state.updating = true;
        state.errorMessage = '';
        state.updateSuccess = false;
      })
      .addCase(deleteTag.pending, (state) => {
        state.deleting = true;
        state.deleteSuccess = false;
        state.errorMessage = '';
      })
      .addCase(createTag.fulfilled, (state) => {
        state.creating = false;
        state.errorMessage = '';
        state.createSuccess = true;
      })
      .addCase(updateTag.fulfilled, (state) => {
        state.updating = false;
        state.updateSuccess = true;
        state.errorMessage = '';
      })
      .addCase(deleteTag.fulfilled, (state) => {
        state.deleting = false;
        state.deleteSuccess = true;
        state.errorMessage = '';
      })
      .addCase(createTag.rejected, (state, action) => {
        state.creating = false;
        state.errorMessage = action.error.message || 'error creating tag';
        state.createSuccess = false;
      })
      .addCase(updateTag.rejected, (state, action) => {
        state.updating = false;
        state.errorMessage = action.error.message || 'error updating tag';
        state.updateSuccess = false;
      })
      .addCase(deleteTag.rejected, (state, action) => {
        state.deleting = false;
        state.errorMessage = action.error.message || 'error deleting tag';
        state.deleteSuccess = false;
      });
  },
});

export const { reset } = TagManagementSlice.actions;
export default TagManagementSlice.reducer;
