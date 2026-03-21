import axios from 'axios';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { serializeAxiosError } from '@app/shared/reducers/reducer.utils';

const apiUrl = 'api/v1/transactions/import';

interface FileResult {
  filename: string;
  created: number;
  skipped: number;
  errors: number;
}

interface ImportResponse {
  files: FileResult[];
  message: string;
}

interface ImportState {
  loading: boolean;
  errorMessage: string;
  results: FileResult[];
  message: string;
}

const initialState: ImportState = {
  loading: false,
  errorMessage: '',
  results: [],
  message: '',
};

export const importFiles = createAsyncThunk(
  'import/files',
  async (params: { files: File[]; account: number }) => {
    const formData = new FormData();
    params.files.forEach((file) => formData.append('files', file));
    formData.append('account', params.account.toString());

    return axios.post<ImportResponse>(apiUrl, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
  { serializeError: serializeAxiosError },
);

export const ImportSlice = createSlice({
  name: 'import',
  initialState,
  reducers: {
    resetImport() {
      return initialState;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(importFiles.pending, (state) => {
        state.loading = true;
        state.errorMessage = '';
        state.results = [];
        state.message = '';
      })
      .addCase(importFiles.rejected, (state, action) => {
        state.loading = false;
        state.errorMessage = action.error.message || 'Failed to import files';
      })
      .addCase(importFiles.fulfilled, (state, action) => {
        state.loading = false;
        state.errorMessage = '';
        state.results = action.payload.data.files;
        state.message = action.payload.data.message;
      });
  },
});

export const { resetImport } = ImportSlice.actions;
export default ImportSlice.reducer;
