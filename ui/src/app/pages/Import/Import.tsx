import * as React from 'react';
import {
  PageSection,
  Title,
  Card,
  CardBody,
  Button,
  Alert,
  Label,
  LabelGroup,
  Content,
  Flex,
  FlexItem,
  MultipleFileUpload,
  MultipleFileUploadMain,
  DropEvent,
  DescriptionList,
  DescriptionListGroup,
  DescriptionListTerm,
  DescriptionListDescription,
  Spinner,
  Bullseye,
} from '@patternfly/react-core';
import UploadIcon from '@patternfly/react-icons/dist/esm/icons/upload-icon';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { importFiles, resetImport } from '@app/shared/reducers/import.reducer';

const Import: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const { loading, errorMessage, results, message } = useAppSelector((state) => state.fileImport);
  const [files, setFiles] = React.useState<File[]>([]);

  const handleFileDrop = (_event: DropEvent, droppedFiles: File[]) => {
    const accepted = droppedFiles.filter((f) => {
      const ext = f.name.toLowerCase();
      return ext.endsWith('.csv') || ext.endsWith('.xlsx') || ext.endsWith('.xls');
    });
    setFiles((prev) => [...prev, ...accepted]);
  };

  const handleFileInputChange = (_event: DropEvent, inputFiles: File[]) => {
    setFiles((prev) => [...prev, ...inputFiles]);
  };

  const handleRemoveFile = (name: string) => {
    setFiles((prev) => prev.filter((f) => f.name !== name));
  };

  const handleImport = () => {
    if (files.length === 0) return;
    dispatch(importFiles({ files, account: 0 }));
  };

  const handleClearAll = () => {
    setFiles([]);
    dispatch(resetImport());
  };

  return (
    <PageSection>
      <Title headingLevel="h1" size="lg" style={{ marginBottom: '1.5rem' }}>
        Import Transactions
      </Title>

      <Card>
        <CardBody>
          <MultipleFileUpload
            onFileDrop={handleFileDrop}
            dropzoneProps={{
              accept: {
                'text/csv': ['.csv'],
                'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': ['.xlsx'],
                'application/vnd.ms-excel': ['.xls'],
              },
            }}
            isHorizontal={false}
          >
            <MultipleFileUploadMain
              titleIcon={<UploadIcon />}
              titleText="Drag and drop files here"
              titleTextSeparator="or"
              infoText="Accepted file types: .csv, .xlsx, .xls"
            />
          </MultipleFileUpload>

          <input
            type="file"
            id="file-upload-input"
            multiple
            accept=".csv,.xlsx,.xls"
            style={{ display: 'none' }}
            onChange={(e) => {
              if (e.target.files) {
                handleFileInputChange(e, Array.from(e.target.files));
                e.target.value = '';
              }
            }}
          />

          {files.length > 0 && (
            <div style={{ marginTop: '1rem' }}>
              <Content component="p" style={{ marginBottom: '0.5rem' }}>
                <strong>{files.length} file{files.length !== 1 ? 's' : ''} ready for import:</strong>
              </Content>

              <LabelGroup>
                {files.map((file) => (
                  <Label
                    key={file.name}
                    variant="outline"
                    onClose={() => handleRemoveFile(file.name)}
                  >
                    {file.name}
                  </Label>
                ))}
              </LabelGroup>

              <Flex style={{ marginTop: '1rem' }}>
                <FlexItem>
                  <Button
                    variant="primary"
                    onClick={handleImport}
                    isLoading={loading}
                    isDisabled={loading || files.length === 0}
                    icon={loading ? undefined : <UploadIcon />}
                  >
                    {loading ? 'Importing...' : 'Import Files'}
                  </Button>
                </FlexItem>
                <FlexItem>
                  <Button variant="secondary" onClick={handleClearAll} isDisabled={loading}>
                    Clear All
                  </Button>
                </FlexItem>
              </Flex>
            </div>
          )}
        </CardBody>
      </Card>

      {errorMessage && (
        <Alert
          variant="danger"
          title="Import Error"
          isInline
          style={{ marginTop: '1rem' }}
        >
          {errorMessage}
        </Alert>
      )}

      {message && results.length > 0 && (
        <Card style={{ marginTop: '1rem' }}>
          <CardBody>
            <Title headingLevel="h2" size="md" style={{ marginBottom: '1rem' }}>
              Import Results
            </Title>

            <Alert variant="success" title={message} isInline style={{ marginBottom: '1rem' }} />

            {results.map((r, i) => (
              <Card key={i} isCompact isPlain style={{ marginBottom: '0.5rem' }}>
                <CardBody>
                  <Content component="p" style={{ fontWeight: 'bold', marginBottom: '0.5rem' }}>
                    {r.filename}
                  </Content>
                  <DescriptionList isHorizontal isCompact>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Created</DescriptionListTerm>
                      <DescriptionListDescription>{r.created}</DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Skipped (duplicates)</DescriptionListTerm>
                      <DescriptionListDescription>{r.skipped}</DescriptionListDescription>
                    </DescriptionListGroup>
                    <DescriptionListGroup>
                      <DescriptionListTerm>Errors</DescriptionListTerm>
                      <DescriptionListDescription>{r.errors}</DescriptionListDescription>
                    </DescriptionListGroup>
                  </DescriptionList>
                </CardBody>
              </Card>
            ))}
          </CardBody>
        </Card>
      )}
    </PageSection>
  );
};

export { Import };
