import * as React from 'react';
import {
  EuiModal,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiModalBody,
  EuiModalFooter,
  EuiButton,
  EuiForm,
  EuiFormRow,
  EuiTextArea,
  EuiCallOut,
  EuiText,
  EuiSpacer,
} from '@elastic/eui';
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

  if (!isOpen) return null;

  return (
    <EuiModal onClose={handleClose} style={{ width: '500px' }}>
      <EuiModalHeader>
        <EuiModalHeaderTitle>Edit Transaction Info</EuiModalHeaderTitle>
      </EuiModalHeader>

      <EuiModalBody>
        {transaction?.description && (
          <>
            <EuiText size="s" color="subdued">
              Transaction: {transaction.description}
            </EuiText>
            <EuiSpacer size="m" />
          </>
        )}

        {error && (
          <>
            <EuiCallOut title="Error" color="danger" iconType="alert">
              {error}
            </EuiCallOut>
            <EuiSpacer size="m" />
          </>
        )}

        <EuiForm>
          <EuiFormRow label="Transaction Info" fullWidth>
            <EuiTextArea
              value={info}
              onChange={(e) => setInfo(e.target.value)}
              placeholder="Enter additional transaction information..."
              rows={4}
              resize="vertical"
              fullWidth
            />
          </EuiFormRow>
        </EuiForm>
      </EuiModalBody>

      <EuiModalFooter>
        <EuiButton onClick={handleClose} isDisabled={loading}>
          Cancel
        </EuiButton>
        <EuiButton
          fill
          color="primary"
          onClick={handleSave}
          isLoading={loading}
          isDisabled={loading}
        >
          Save
        </EuiButton>
      </EuiModalFooter>
    </EuiModal>
  );
};