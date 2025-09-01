import React, { useState, useEffect } from 'react';
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
  EuiComboBox,
  EuiComboBoxOptionOption,
  EuiBadge,
  EuiFlexGroup,
  EuiFlexItem,
  EuiSpacer,
  EuiCallOut,
  EuiText,
} from '@elastic/eui';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { addLabelToTransaction, updateTransactionInfo } from '@app/shared/reducers/transaction.reducer';
import { getLabels } from '@app/shared/reducers/label.reducer';
import { ITransaction } from '@app/shared/models/transaction';

interface TransactionDetailsModalProps {
  isOpen: boolean;
  onClose: () => void;
  transaction?: ITransaction;
  onSuccess?: () => void;
}

export const TransactionDetailsModal: React.FC<TransactionDetailsModalProps> = ({
  isOpen,
  onClose,
  transaction,
  onSuccess,
}) => {
  const dispatch = useAppDispatch();
  const { addingLabel, addLabelSuccess, updatingInfo, errorMessage } = useAppSelector((state) => state.transactions);
  const labels = useAppSelector((state) => state.labels);

  // Local state for form data
  const [info, setInfo] = useState('');
  const [selectedLabels, setSelectedLabels] = useState<Array<{ key: string; value: string }>>([]);

  // Initialize form data when modal opens or transaction changes
  useEffect(() => {
    if (transaction) {
      setInfo(transaction.info || '');
      setSelectedLabels([]);
    }
  }, [transaction]);

  // Load labels when modal opens
  useEffect(() => {
    if (isOpen && labels.labels.length === 0) {
      dispatch(getLabels());
    }
  }, [isOpen, dispatch, labels.labels.length]);

  const handleClose = () => {
    setInfo('');
    setSelectedLabels([]);
    onClose();
  };

  const handleLabelChange = (selectedOptions: EuiComboBoxOptionOption[]) => {
    const newLabels = selectedOptions.map(option => {
      const [key, value] = option.label.split('=', 2);
      return { key: key.trim(), value: value.trim() };
    });
    setSelectedLabels(newLabels);
  };

  const handleSave = async () => {
    if (!transaction) return;

    try {
      // Save transaction info if changed
      if (info !== (transaction.info || '')) {
        await dispatch(updateTransactionInfo({ 
          transactionHref: transaction.href, 
          info 
        }));
      }

      // Add new labels
      for (const label of selectedLabels) {
        await dispatch(addLabelToTransaction({
          transactionHref: transaction.href,
          key: label.key,
          value: label.value,
        }));
      }

      if (onSuccess) {
        onSuccess();
      }
      handleClose();
    } catch (error) {
      console.error('Failed to save transaction details:', error);
    }
  };

  const isLoading = addingLabel || updatingInfo;

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
    <EuiModal onClose={handleClose} style={{ width: '600px' }}>
      <EuiModalHeader>
        <EuiModalHeaderTitle>Edit Transaction Details</EuiModalHeaderTitle>
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

        {errorMessage && (
          <>
            <EuiCallOut title="Error" color="danger" iconType="alert">
              {errorMessage}
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
              rows={3}
              resize="vertical"
              fullWidth
            />
          </EuiFormRow>

          <EuiSpacer size="l" />

          <EuiFormRow label="Add Labels" fullWidth>
            <EuiComboBox
              placeholder="Enter label (key=value) or select from dropdown"
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
                }
              }}
              isClearable={true}
              fullWidth
            />
          </EuiFormRow>

          {/* Current Labels Display */}
          {transaction && transaction.labels && transaction.labels.length > 0 && (
            <>
              <EuiSpacer size="m" />
              <EuiFormRow label="Current Labels" fullWidth>
                <EuiFlexGroup gutterSize="s" wrap>
                  {transaction.labels.map((label, index) => (
                    <EuiFlexItem grow={false} key={index}>
                      <EuiBadge color="hollow">
                        {label.key}={label.value}
                      </EuiBadge>
                    </EuiFlexItem>
                  ))}
                </EuiFlexGroup>
              </EuiFormRow>
            </>
          )}
        </EuiForm>
      </EuiModalBody>

      <EuiModalFooter>
        <EuiButton onClick={handleClose} isDisabled={isLoading}>
          Cancel
        </EuiButton>
        <EuiButton
          fill
          color="primary"
          onClick={handleSave}
          isLoading={isLoading}
          isDisabled={isLoading}
        >
          Save Changes
        </EuiButton>
      </EuiModalFooter>
    </EuiModal>
  );
};