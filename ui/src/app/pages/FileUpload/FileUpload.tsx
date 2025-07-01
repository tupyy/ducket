import * as React from 'react';
import {
  PageSection,
  Title,
  Card,
  CardBody,
  MultipleFileUpload,
  MultipleFileUploadMain,
  MultipleFileUploadStatus,
  MultipleFileUploadStatusItem,
  HelperText,
  HelperTextItem,
  Button,
  Alert,
  AlertVariant,
  Spinner,
} from '@patternfly/react-core';
import { UploadIcon } from '@patternfly/react-icons';
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

  const handleFileDrop = (_event: any, files: File[]) => {
    const newFiles: UploadedFile[] = files.map((file) => ({
      file,
      progress: 0,
      status: 'uploading' as const,
      id: `${file.name}-${Date.now()}`,
    }));

    setUploadedFiles((prev) => [...prev, ...newFiles]);
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
    <PageSection>
      <Title headingLevel="h1" size="lg">
        File Upload
      </Title>
      <Card>
        <CardBody>
          <MultipleFileUpload
            onFileDrop={handleFileDrop}
            dropzoneProps={{
              accept: {
                'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': ['.xlsx'],
                'application/vnd.ms-excel': ['.xls'],
                'text/csv': ['.csv'],
              },
              maxSize: 10 * 1024 * 1024, // 10MB
            }}
          >
            <MultipleFileUploadMain
              titleIcon={<UploadIcon />}
              titleText="Drag and drop files here"
              titleTextSeparator="or"
              infoText="Upload transaction files"
            />
            {uploadedFiles.length > 0 && (
              <MultipleFileUploadStatus
                statusToggleText={`${uploadedFiles.length} files ready for import`}
                statusToggleIcon="info"
              >
                {uploadedFiles.map((uploadedFile) => (
                  <MultipleFileUploadStatusItem
                    key={uploadedFile.id}
                    file={uploadedFile.file}
                    onClearClick={() => handleFileRemove(uploadedFile.id)}
                  />
                ))}
              </MultipleFileUploadStatus>
            )}
            <HelperText>
              <HelperTextItem>
                Accepted file types: Excel (.xlsx, .xls) and CSV (.csv) files up to 10MB each
              </HelperTextItem>
            </HelperText>
          </MultipleFileUpload>

          {uploadedFiles.length > 0 && (
            <div style={{ marginTop: '1rem', display: 'flex', gap: '1rem' }}>
              <Button variant="primary" onClick={handleImport} isDisabled={loading || uploadedFiles.length === 0}>
                {loading ? <Spinner size="sm" /> : null}
                {loading ? ' Importing...' : 'Import Files'}
              </Button>
              <Button variant="secondary" onClick={clearAllFiles}>
                Clear All
              </Button>
            </div>
          )}
        </CardBody>
      </Card>

      {/* Import Results */}
      {errorMessage && (
        <Alert variant={AlertVariant.danger} title="Import Error" style={{ marginTop: '1rem' }}>
          {errorMessage}
        </Alert>
      )}

      {importSuccess && summary && (
        <Card style={{ marginTop: '1rem' }}>
          <CardBody>
            <Title headingLevel="h2" size="md">
              Import Results
            </Title>
            <Alert variant={AlertVariant.success} title="Import Completed" style={{ marginTop: '1rem' }}>
              {lastImportMessage}
            </Alert>

            <div style={{ marginTop: '1rem' }}>
              <h3>Summary</h3>
              <ul>
                <li>Files processed: {summary.files_processed}</li>
                <li>Total rows: {summary.total_rows}</li>
                <li>Processed rows: {summary.total_processed}</li>
                <li>Created transactions: {summary.total_created}</li>
                <li>Updated transactions: {summary.total_updated}</li>
                <li>Errors: {summary.total_errors}</li>
              </ul>
            </div>

            {results && results.length > 0 && (
              <div style={{ marginTop: '1rem' }}>
                <h3>File Details</h3>
                {results.map((result, index) => (
                  <div
                    key={index}
                    style={{ marginBottom: '1rem', padding: '1rem', border: '1px solid #ccc', borderRadius: '4px' }}
                  >
                    <h4>{result.filename}</h4>
                    <p>
                      Rows: {result.total_rows} | Processed: {result.processed_rows} | Created: {result.created_count} |
                      Ignored: {result.ignored_count} | Errors: {result.error_count}
                    </p>
                    {result.errors && result.errors.length > 0 && (
                      <div>
                        <strong>Errors:</strong>
                        <ul>
                          {result.errors.map((error, errorIndex) => (
                            <li key={errorIndex}>{error}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </CardBody>
        </Card>
      )}
    </PageSection>
  );
};

export { FileUpload };
