import * as React from 'react';
import {
  EuiModal,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiModalBody,
  EuiModalFooter,
  EuiButton,
  EuiText,
} from '@elastic/eui';
import { ILabelTransaction, ITransaction } from '@app/shared/models/transaction';

interface RemoveLabelModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  transaction?: ITransaction;
  label?: ILabelTransaction;
}

export const RemoveLabelModal: React.FC<RemoveLabelModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  transaction,
  label,
}) => {
  if (!isOpen) return null;

  return (
    <EuiModal onClose={onClose} style={{ width: '400px' }}>
      <EuiModalHeader>
        <EuiModalHeaderTitle>Remove Label</EuiModalHeaderTitle>
      </EuiModalHeader>
      
      <EuiModalBody>
        {transaction && label && (
          <EuiText>
            <p>
              Are you sure you want to remove the label <strong>{label.key}={label.value}</strong> from the transaction "<strong>{transaction.description}</strong>"?
            </p>
          </EuiText>
        )}
      </EuiModalBody>
      
      <EuiModalFooter>
        <EuiButton onClick={onClose}>
          Cancel
        </EuiButton>
        <EuiButton fill color="danger" onClick={onConfirm}>
          Remove
        </EuiButton>
      </EuiModalFooter>
    </EuiModal>
  );
};