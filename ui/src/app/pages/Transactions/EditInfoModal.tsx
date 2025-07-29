import * as React from 'react';
import {
  Modal,
  ModalVariant,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  Form,
  FormGroup,
  TextArea,
  Alert,
} from '@patternfly/react-core';
import { ITransaction } from '@app/shared/models/transaction';

interface EditInfoModalProps {
  isOpen: boolean;
  onClose: () => void;
  transaction?: ITransaction;
  onSave: (transactionHref: string, info: string) => void;
  loading?: boolean;
  error?: string;
}

export const EditInfoModal: React.FC<EditInfoModalProps> = ({
  isOpen,
  onClose,
  transaction,
  onSave,
  loading = false,
  error,
}) => {
  const [info, setInfo] = React.useState('');

  // Initialize info field when modal opens or transaction changes
  React.useEffect(() => {
    if (transaction) {
      setInfo(transaction.info || '');
    }
  }, [transaction]);

  const handleClose = () => {
    setInfo('');
    onClose();
  };

  const handleSave = () => {
    if (transaction) {
      onSave(transaction.href, info);
    }
  };

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    handleSave();
  };

  return (
    <Modal
      variant={ModalVariant.small}
      isOpen={isOpen}
      onClose={handleClose}
      aria-labelledby="edit-info-modal-title"
      aria-describedby="edit-info-modal-description"
    >
      <ModalHeader
        title="Edit Transaction Info"
        labelId="edit-info-modal-title"
        description={
          transaction?.description 
            ? `Transaction: ${transaction.description}` 
            : undefined
        }
        descriptorId="edit-info-modal-description"
      />
      <ModalBody>
        {error && (
          <Alert variant="danger" title="Error" style={{ marginBottom: '1rem' }}>
            {error}
          </Alert>
        )}

        <Form id="edit-info-form" onSubmit={handleSubmit}>
          <FormGroup label="Transaction Info" fieldId="info-input">
            <TextArea
              id="info-input"
              value={info}
              onChange={(_event, value) => setInfo(value)}
              placeholder="Enter additional transaction information..."
              rows={4}
              resizeOrientation="vertical"
            />
          </FormGroup>
        </Form>
      </ModalBody>
      <ModalFooter>
        <Button
          variant="primary"
          onClick={handleSave}
          isLoading={loading}
          isDisabled={loading}
        >
          Save
        </Button>
        <Button variant="link" onClick={handleClose} isDisabled={loading}>
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
}; 