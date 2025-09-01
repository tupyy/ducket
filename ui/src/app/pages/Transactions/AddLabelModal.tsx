import React, { useState } from 'react';
import {
  EuiModal,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiModalBody,
  EuiModalFooter,
  EuiButton,
  EuiForm,
  EuiFormRow,
  EuiCallOut,
  EuiComboBox,
  EuiComboBoxOptionOption,
  EuiBadge,
  EuiFlexGroup,
  EuiFlexItem,
  EuiText,
  EuiSpacer,
} from '@elastic/eui';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { addLabelToTransaction } from '@app/shared/reducers/transaction.reducer';
import { getLabels } from '@app/shared/reducers/label.reducer';

interface AddLabelModalProps {
  isOpen: boolean;
  onClose: () => void;
  transactionHref?: string;
  transactionHrefs?: string[];
  transactionDescription?: string;
  onSuccess?: () => void;
}

export const AddLabelModal: React.FunctionComponent<AddLabelModalProps> = ({
  isOpen,
  onClose,
  transactionHref,
  transactionHrefs,
  transactionDescription = '',
  onSuccess,
}) => {
  const dispatch = useAppDispatch();
  const { addingLabel, addLabelSuccess, errorMessage } = useAppSelector((state) => state.transactions);
  const labels = useAppSelector((state) => state.labels);

  const [selectedLabels, setSelectedLabels] = useState<Array<{ key: string; value: string }>>([]);
  const [validationError, setValidationError] = useState('');

  // Determine if this is a bulk operation
  const isBulkOperation = Boolean(transactionHrefs?.length);
  const transactionCount = isBulkOperation ? transactionHrefs!.length : 1;
  const targetTransactionHrefs = isBulkOperation 
    ? transactionHrefs! 
    : (transactionHref ? [transactionHref] : []);

  // Load available labels when modal opens
  React.useEffect(() => {
    if (isOpen) {
      dispatch(getLabels());
    }
  }, [isOpen, dispatch]);

  // Reset form when modal opens/closes
  React.useEffect(() => {
    if (isOpen) {
      setSelectedLabels([]);
      setValidationError('');
    }
  }, [isOpen]);

  // Close modal when labels are successfully added
  React.useEffect(() => {
    if (addLabelSuccess) {
      handleClose();
      onSuccess?.();
    }
  }, [addLabelSuccess, onSuccess]);

  const handleSubmit = async (event?: React.FormEvent<HTMLFormElement>) => {
    if (event) {
      event.preventDefault();
    }

    if (selectedLabels.length === 0) {
      setValidationError('Please add at least one label');
      return;
    }

    if (targetTransactionHrefs.length === 0) {
      setValidationError('No transactions selected');
      return;
    }

    // Submit all selected labels to all target transactions
    for (const href of targetTransactionHrefs) {
      for (const label of selectedLabels) {
        await dispatch(
          addLabelToTransaction({
            transactionHref: href,
            key: label.key,
            value: label.value,
          })
        );
      }
    }
  };

  const handleClose = () => {
    if (!addingLabel) {
      onClose();
    }
  };

  const handleLabelChange = (selectedOptions: EuiComboBoxOptionOption[]) => {
    const newLabels = selectedOptions.map(option => {
      const [key, value] = option.label.split('=', 2);
      return { key: key.trim(), value: value.trim() };
    });
    setSelectedLabels(newLabels);
    if (validationError) {
      setValidationError('');
    }
  };

  // Convert labels to ComboBox options
  const availableLabelOptions: EuiComboBoxOptionOption[] = labels.labels
    .map((label) => ({
      label: `${label.key}=${label.value}`,
      value: `${label.key}=${label.value}`,
    }));

  const selectedLabelOptions: EuiComboBoxOptionOption[] = selectedLabels.map((label) => ({
    label: `${label.key}=${label.value}`,
    value: `${label.key}=${label.value}`,
  }));

  if (!isOpen) return null;

  return (
    <EuiModal onClose={handleClose} style={{ width: '500px' }}>
      <EuiModalHeader>
        <EuiModalHeaderTitle>
          {isBulkOperation ? `Add Labels to ${transactionCount} Transactions` : "Add Custom Labels"}
        </EuiModalHeaderTitle>
      </EuiModalHeader>

      <EuiModalBody>
        {(isBulkOperation || transactionDescription) && (
          <>
            <EuiText size="s" color="subdued">
              {isBulkOperation 
                ? `Add labels to ${transactionCount} selected transaction${transactionCount !== 1 ? 's' : ''}`
                : `Transaction: ${transactionDescription}`
              }
            </EuiText>
            <EuiSpacer size="m" />
          </>
        )}

        {errorMessage && (
          <>
            <EuiCallOut title="Error" color="danger" iconType="alert">
              {errorMessage}
            </EuiCallOut>
            <EuiSpacer size="m" />
          </>
        )}

        <EuiForm>
          <EuiFormRow 
            label="Labels" 
            isInvalid={!!validationError} 
            error={validationError ? [validationError] : []}
            helpText={
              <EuiText size="xs" color="subdued">
                Enter labels in key=value format (e.g., category=food, type=expense).
                {isBulkOperation 
                  ? ` These labels will be added to all ${transactionCount} selected transactions.`
                  : ' You can add multiple labels.'
                }
              </EuiText>
            }
            fullWidth
          >
            <EuiComboBox
              placeholder="Type label in format: key=value (e.g., category=food)"
              options={availableLabelOptions}
              selectedOptions={selectedLabelOptions}
              onChange={handleLabelChange}
              onCreateOption={(searchValue: string) => {
                if (searchValue.includes('=')) {
                  const [key, value] = searchValue.split('=', 2);
                  if (key.trim() && value.trim()) {
                    const newLabel = { key: key.trim(), value: value.trim() };
                    setSelectedLabels([...selectedLabels, newLabel]);
                  }
                } else {
                  setValidationError('Label must be in format: key=value');
                }
              }}
              isClearable={true}
              isDisabled={addingLabel}
              fullWidth
            />
          </EuiFormRow>
        </EuiForm>
      </EuiModalBody>

      <EuiModalFooter>
        <EuiButton onClick={handleClose} isDisabled={addingLabel}>
          Cancel
        </EuiButton>
        <EuiButton
          fill
          color="primary"
          onClick={() => handleSubmit()}
          isLoading={addingLabel}
          isDisabled={addingLabel || selectedLabels.length === 0}
        >
          Add {selectedLabels.length} Label{selectedLabels.length !== 1 ? 's' : ''}
          {isBulkOperation ? ` to ${transactionCount} Transaction${transactionCount !== 1 ? 's' : ''}` : ''}
        </EuiButton>
      </EuiModalFooter>
    </EuiModal>
  );
};