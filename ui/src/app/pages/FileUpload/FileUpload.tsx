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
} from '@patternfly/react-core';
import { UploadIcon } from '@patternfly/react-icons';

interface UploadedFile {
  file: File;
  progress: number;
  status: 'uploading' | 'complete' | 'error';
  id: string;
}

const FileUpload: React.FunctionComponent = () => {
  const [uploadedFiles, setUploadedFiles] = React.useState<UploadedFile[]>([]);
  const [isUploading, setIsUploading] = React.useState(false);

  const handleFileDrop = (_event: any, files: File[]) => {
    const newFiles: UploadedFile[] = files.map((file) => ({
      file,
      progress: 0,
      status: 'uploading' as const,
      id: `${file.name}-${Date.now()}`,
    }));

    setUploadedFiles((prev) => [...prev, ...newFiles]);
    setIsUploading(true);

    // Simulate file upload progress
    newFiles.forEach((uploadedFile) => {
      simulateUpload(uploadedFile.id);
    });
  };

  const simulateUpload = (fileId: string) => {
    const interval = setInterval(() => {
      setUploadedFiles((prev) =>
        prev.map((file) => {
          if (file.id === fileId) {
            const newProgress = Math.min(file.progress + 10, 100);
            const newStatus = newProgress === 100 ? 'complete' : 'uploading';

            if (newProgress === 100) {
              // Check if all uploads are complete
              setTimeout(() => {
                setUploadedFiles((current) => {
                  const allComplete = current.every((f) => f.status === 'complete' || f.status === 'error');
                  if (allComplete) {
                    setIsUploading(false);
                  }
                  return current;
                });
              }, 100);
            }

            return { ...file, progress: newProgress, status: newStatus };
          }
          return file;
        }),
      );

      if (uploadedFiles.find((f) => f.id === fileId)?.progress === 100) {
        clearInterval(interval);
      }
    }, 200);
  };

  const handleFileRemove = (fileId: string) => {
    setUploadedFiles((prev) => prev.filter((file) => file.id !== fileId));
  };

  const clearAllFiles = () => {
    setUploadedFiles([]);
    setIsUploading(false);
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
                statusToggleText={`${uploadedFiles.length} of ${uploadedFiles.length} files uploaded`}
                statusToggleIcon="success"
              >
                {uploadedFiles.map((uploadedFile) => (
                  <MultipleFileUploadStatusItem
                    key={uploadedFile.id}
                    file={uploadedFile.file}
                    onClearClick={() => handleFileRemove(uploadedFile.id)}
                    progressValue={uploadedFile.progress}
                    progressVariant={
                      uploadedFile.status === 'error'
                        ? 'danger'
                        : uploadedFile.status === 'complete'
                          ? 'success'
                          : undefined
                    }
                  />
                ))}
              </MultipleFileUploadStatus>
            )}
          </MultipleFileUpload>
          <HelperText>
            <HelperTextItem>
              Accepted file types: Excel (.xlsx, .xls) and CSV (.csv) files up to 10MB each
            </HelperTextItem>
          </HelperText>
        </CardBody>
      </Card>
    </PageSection>
  );
};

export { FileUpload };
