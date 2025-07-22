import * as React from 'react';
import {
  Modal,
  ModalVariant,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
} from '@patternfly/react-core';
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
  return (
    <Modal
      variant={ModalVariant.small}
      isOpen={isOpen}
      onClose={onClose}
    >
      <ModalHeader title="Remove Label" />
      <ModalBody>
        {transaction && label && (
          <p>
            Are you sure you want to remove the label <strong>{label.key}={label.value}</strong> from the transaction "<strong>{transaction.description}</strong>"?
          </p>
        )}
      </ModalBody>
      <ModalFooter>
        <Button variant="danger" onClick={onConfirm}>
          Remove
        </Button>
        <Button variant="link" onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
}; 