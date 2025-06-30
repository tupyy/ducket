import axios from 'axios';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';
import { IImportResponse, IImportResult, IImportSummary } from '@app/shared/models/import';

const importApiUrl = 'api/v1/import';

const initialState = {
  loading: false,
  errorMessage: '',
  importSuccess: false,
  results: [] as IImportResult[],
  summary: null as IImportSummary | null,
  lastImportMessage: '',
};

export const importFiles = createAsyncThunk(
  'import/files',
  async (files: File[]) => {
    const formData = new FormData();
    
    // Append all files to the form data
    files.forEach((file) => {
      formData.append('files', file);
    });

    const response = await axios.post<IImportResponse>(importApiUrl, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    return response.data;
  },
  { serializeError: serializeAxiosError },
);

export type ImportState = Readonly<typeof initialState>;

export const ImportSlice = createSlice({
  name: 'import',
  initialState: initialState as ImportState,
  reducers: {
    reset() {
      return initialState;
    },
    clearResults(state) {
      state.results = [];
      state.summary = null;
      state.lastImportMessage = '';
      state.importSuccess = false;
      state.errorMessage = '';
    },
  },
  extraReducers(builder) {
    builder
      .addCase(importFiles.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
        state.importSuccess = false;
      })
      .addCase(importFiles.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to import files';
        state.importSuccess = false;
      })
      .addCase(importFiles.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.importSuccess = action.payload.success;
        state.results = action.payload.results;
        state.summary = action.payload.summary;
        state.lastImportMessage = action.payload.message;
      });
  },
});

export const { reset, clearResults } = ImportSlice.actions;
export default ImportSlice.reducer; 