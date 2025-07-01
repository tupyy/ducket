export interface IImportResult {
  filename: string;
  total_rows: number;
  processed_rows: number;
  created_count: number;
  ignored_count: number;
  error_count: number;
  errors?: string[];
}

export interface IImportResponse {
  success: boolean;
  message: string;
  summary: IImportSummary;
  results: IImportResult[];
}

export interface IImportSummary {
  files_processed: number;
  total_rows: number;
  total_processed: number;
  total_created: number;
  total_updated: number;
  total_errors: number;
}

export interface IImportRequest {
  files: File[];
} 