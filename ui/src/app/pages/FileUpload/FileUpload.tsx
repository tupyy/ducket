import * as React from 'react';
import {
  EuiPageSection,
  EuiTitle,
  EuiPanel,
  EuiFilePicker,
  EuiFlexGroup,
  EuiFlexItem,
  EuiButton,
  EuiCallOut,
  EuiLoadingSpinner,
  EuiText,
  EuiSpacer,
  EuiBadge,
  EuiButtonIcon,
} from '@elastic/eui';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { importFiles, clearResults } from '@app/shared/reducers/import.reducer';

interface UploadedFile {
  file: File;
  progress: number;
  status: 'uploading' | 'complete' | 'error';
  id: string;
}

const FileUpload: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const { loading, errorMessage, importSuccess, results, summary, lastImportMessage } = useAppSelector(
    (state) => state.import
  );

  const [uploadedFiles, setUploadedFiles] = React.useState<UploadedFile[]>([]);

  const handleFileChange = (files: FileList | null) => {
    if (files) {
      const newFiles: UploadedFile[] = Array.from(files).map((file) => ({
        file,
        progress: 0,
        status: 'uploading' as const,
        id: `${file.name}-${Date.now()}`,
      }));

      setUploadedFiles((prev) => [...prev, ...newFiles]);
    }
  };

  const handleImport = async () => {
    const filesToImport = uploadedFiles.map((uploadedFile) => uploadedFile.file);
    if (filesToImport.length > 0) {
      await dispatch(importFiles(filesToImport));
    }
  };

  const handleFileRemove = (fileId: string) => {
    setUploadedFiles((prev) => prev.filter((file) => file.id !== fileId));
  };

  const clearAllFiles = () => {
    setUploadedFiles([]);
    dispatch(clearResults());
  };

  return (
    <EuiPageSection>
      <EuiTitle size="l">
        <h1>File Upload</h1>
      </EuiTitle>
      
      <EuiSpacer size="m" />
      
      <EuiPanel paddingSize="l">
        <EuiFilePicker
          id="file-picker"
          multiple
          onChange={handleFileChange}
          display="large"
          accept=".xlsx,.xls,.csv"
          aria-label="Upload transaction files"
        />
        
        <EuiSpacer size="s" />
        
        <EuiText size="s" color="subdued">
          Accepted file types: Excel (.xlsx, .xls) and CSV (.csv) files up to 10MB each
        </EuiText>

        {uploadedFiles.length > 0 && (
          <>
            <EuiSpacer size="m" />
            
            <EuiText size="s">
              <strong>{uploadedFiles.length} files ready for import:</strong>
            </EuiText>
            
            <EuiSpacer size="s" />
            
            <EuiFlexGroup wrap gutterSize="s">
              {uploadedFiles.map((uploadedFile) => (
                <EuiFlexItem grow={false} key={uploadedFile.id}>
                  <EuiBadge 
                    color="hollow"
                    iconType="cross"
                    iconSide="right"
                    iconOnClick={() => handleFileRemove(uploadedFile.id)}
                    iconOnClickAriaLabel={`Remove ${uploadedFile.file.name}`}
                  >
                    {uploadedFile.file.name}
                  </EuiBadge>
                </EuiFlexItem>
              ))}
            </EuiFlexGroup>
            
            <EuiSpacer size="m" />
            
            <EuiFlexGroup gutterSize="s" alignItems="center">
              <EuiFlexItem grow={false}>
                <EuiButton 
                  fill 
                  onClick={handleImport} 
                  isDisabled={loading || uploadedFiles.length === 0}
                  iconType={loading ? undefined : "importAction"}
                >
                  {loading && <EuiLoadingSpinner size="s" />}
                  {loading ? ' Importing...' : 'Import Files'}
                </EuiButton>
              </EuiFlexItem>
              <EuiFlexItem grow={false}>
                <EuiButton onClick={clearAllFiles}>
                  Clear All
                </EuiButton>
              </EuiFlexItem>
            </EuiFlexGroup>
          </>
        )}
      </EuiPanel>

      {/* Import Results */}
      {errorMessage && (
        <>
          <EuiSpacer size="m" />
          <EuiCallOut title="Import Error" color="danger" iconType="alert">
            {errorMessage}
          </EuiCallOut>
        </>
      )}

      {importSuccess && summary && (
        <>
          <EuiSpacer size="m" />
          <EuiPanel paddingSize="l">
            <EuiTitle size="m">
              <h2>Import Results</h2>
            </EuiTitle>
            
            <EuiSpacer size="m" />
            
            <EuiCallOut title="Import Completed" color="success" iconType="check">
              {lastImportMessage}
            </EuiCallOut>

            <EuiSpacer size="m" />
            
            <EuiTitle size="s">
              <h3>Summary</h3>
            </EuiTitle>
            
            <EuiSpacer size="s" />
            
            <EuiText size="s">
              <ul>
                <li>Files processed: {summary.files_processed}</li>
                <li>Total rows: {summary.total_rows}</li>
                <li>Processed rows: {summary.total_processed}</li>
                <li>Created transactions: {summary.total_created}</li>
                <li>Updated transactions: {summary.total_updated}</li>
                <li>Errors: {summary.total_errors}</li>
              </ul>
            </EuiText>

            {results && results.length > 0 && (
              <>
                <EuiSpacer size="m" />
                <EuiTitle size="s">
                  <h3>File Details</h3>
                </EuiTitle>
                <EuiSpacer size="s" />
                
                {results.map((result, index) => (
                  <EuiPanel 
                    key={index} 
                    paddingSize="m" 
                    style={{ marginBottom: '1rem' }}
                    color="subdued"
                  >
                    <EuiTitle size="xs">
                      <h4>{result.filename}</h4>
                    </EuiTitle>
                    <EuiSpacer size="xs" />
                    <EuiText size="s">
                      Rows: {result.total_rows} | Processed: {result.processed_rows} | Created: {result.created_count} |
                      Ignored: {result.ignored_count} | Errors: {result.error_count}
                    </EuiText>
                    {result.errors && result.errors.length > 0 && (
                      <>
                        <EuiSpacer size="s" />
                        <EuiText size="s">
                          <strong>Errors:</strong>
                          <ul>
                            {result.errors.map((error, errorIndex) => (
                              <li key={errorIndex}>{error}</li>
                            ))}
                          </ul>
                        </EuiText>
                      </>
                    )}
                  </EuiPanel>
                ))}
              </>
            )}
          </EuiPanel>
        </>
      )}
    </EuiPageSection>
  );
};

export { FileUpload };
