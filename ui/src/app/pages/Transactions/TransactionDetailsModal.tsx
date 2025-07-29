import React, { useState, useEffect } from 'react';
import {
  Button,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalVariant,
  Form,
  FormGroup,
  Alert,
  AlertVariant,
  MenuToggle,
  MenuToggleElement,
  Dropdown,
  DropdownList,
  DropdownItem,
  TextInputGroup,
  TextInputGroupMain,
  TextInputGroupUtilities,
  TextArea,
  Label,
  Flex,
  FlexItem,
  Divider,
} from '@patternfly/react-core';
import { TimesIcon } from '@patternfly/react-icons';
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
  const [labelInputValue, setLabelInputValue] = useState('');
  const [selectedLabels, setSelectedLabels] = useState<Array<{ key: string; value: string }>>([]);
  const [isLabelDropdownOpen, setIsLabelDropdownOpen] = useState(false);
  const [filteredLabels, setFilteredLabels] = useState<Array<{ key: string; value: string }>>([]);

  // Initialize form data when modal opens or transaction changes
  useEffect(() => {
    if (transaction) {
      setInfo(transaction.info || '');
      setSelectedLabels([]);
      setLabelInputValue('');
    }
  }, [transaction]);

  // Load labels when modal opens
  useEffect(() => {
    if (isOpen && labels.labels.length === 0) {
      dispatch(getLabels());
    }
  }, [isOpen, dispatch, labels.labels.length]);

  // Filter labels based on input
  useEffect(() => {
    if (labelInputValue) {
      const availableLabels = labels.labels.map(label => ({
        key: label.key,
        value: label.value
      }));
      
      const filtered = availableLabels.filter(
        label =>
          label.key.toLowerCase().includes(labelInputValue.toLowerCase()) ||
          label.value.toLowerCase().includes(labelInputValue.toLowerCase())
      );
      setFilteredLabels(filtered);
    } else {
      setFilteredLabels([]);
    }
  }, [labelInputValue, labels.labels]);

  const handleClose = () => {
    setInfo('');
    setSelectedLabels([]);
    setLabelInputValue('');
    setIsLabelDropdownOpen(false);
    onClose();
  };

  const handleAddLabel = (key: string, value: string) => {
    const newLabel = { key, value };
    if (!selectedLabels.some(label => label.key === key && label.value === value)) {
      setSelectedLabels([...selectedLabels, newLabel]);
    }
    setLabelInputValue('');
    setIsLabelDropdownOpen(false);
  };

  const handleRemoveLabel = (index: number) => {
    setSelectedLabels(selectedLabels.filter((_, i) => i !== index));
  };

  const handleLabelInputKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' && labelInputValue.includes('=')) {
      event.preventDefault();
      const [key, value] = labelInputValue.split('=', 2);
      if (key.trim() && value.trim()) {
        handleAddLabel(key.trim(), value.trim());
      }
    }
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

  const labelToggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      variant="typeahead"
      onClick={() => setIsLabelDropdownOpen(!isLabelDropdownOpen)}
      isExpanded={isLabelDropdownOpen}
      isFullWidth
    >
      <TextInputGroup isPlain>
        <TextInputGroupMain
          value={labelInputValue}
          onClick={() => setIsLabelDropdownOpen(!isLabelDropdownOpen)}
          onChange={(_event, value) => setLabelInputValue(value)}
          onKeyDown={handleLabelInputKeyDown}
          placeholder="Enter label (key=value) or select from dropdown"
          role="combobox"
          isExpanded={isLabelDropdownOpen}
          aria-controls="label-listbox"
        />

        {labelInputValue && (
          <TextInputGroupUtilities>
            <Button variant="plain" onClick={() => setLabelInputValue('')} aria-label="Clear input">
              <TimesIcon />
            </Button>
          </TextInputGroupUtilities>
        )}
      </TextInputGroup>
    </MenuToggle>
  );

  return (
    <Modal
      variant={ModalVariant.medium}
      isOpen={isOpen}
      onClose={handleClose}
      aria-labelledby="transaction-details-modal-title"
      aria-describedby="transaction-details-modal-description"
    >
      <ModalHeader
        title="Edit Transaction Details"
        labelId="transaction-details-modal-title"
        description={
          transaction?.description 
            ? `Transaction: ${transaction.description}` 
            : undefined
        }
        descriptorId="transaction-details-modal-description"
      />
      <ModalBody>
        {errorMessage && (
          <Alert variant={AlertVariant.danger} title="Error" style={{ marginBottom: '1rem' }}>
            {errorMessage}
          </Alert>
        )}

        <Form>
          {/* Transaction Info Section */}
          <FormGroup label="Transaction Info" fieldId="info-input">
            <TextArea
              id="info-input"
              value={info}
              onChange={(_event, value) => setInfo(value)}
              placeholder="Enter additional transaction information..."
              rows={3}
              resizeOrientation="vertical"
            />
          </FormGroup>

          <Divider style={{ margin: '1.5rem 0' }} />

          {/* Labels Section */}
          <FormGroup label="Add Labels" fieldId="label-input">
            <Dropdown
              isOpen={isLabelDropdownOpen}
              onSelect={(_event, value) => {
                const selectedValue = value as string;
                if (selectedValue.includes('=')) {
                  const [key, val] = selectedValue.split('=', 2);
                  handleAddLabel(key.trim(), val.trim());
                }
              }}
              onOpenChange={(isOpen) => setIsLabelDropdownOpen(isOpen)}
              toggle={labelToggle}
              shouldFocusToggleOnSelect
            >
              <DropdownList>
                {filteredLabels.map((label, index) => (
                  <DropdownItem
                    key={index}
                    value={`${label.key}=${label.value}`}
                  >
                    {label.key}={label.value}
                  </DropdownItem>
                ))}
                {labelInputValue && labelInputValue.includes('=') && (
                  <DropdownItem value={labelInputValue}>
                    Create: {labelInputValue}
                  </DropdownItem>
                )}
              </DropdownList>
            </Dropdown>
          </FormGroup>

          {/* Selected Labels Display */}
          {selectedLabels.length > 0 && (
            <FormGroup label="Labels to Add" fieldId="selected-labels">
              <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }} flexWrap={{ default: 'wrap' }}>
                {selectedLabels.map((label, index) => (
                  <FlexItem key={index}>
                    <Label
                      variant="outline"
                      color="blue"
                      onClose={() => handleRemoveLabel(index)}
                      closeBtnAriaLabel={`Remove ${label.key}=${label.value} label`}
                    >
                      {label.key}={label.value}
                    </Label>
                  </FlexItem>
                ))}
              </Flex>
            </FormGroup>
          )}

          {/* Current Labels Display */}
          {transaction && transaction.labels && transaction.labels.length > 0 && (
            <FormGroup label="Current Labels" fieldId="current-labels">
              <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }} flexWrap={{ default: 'wrap' }}>
                {transaction.labels.map((label, index) => (
                  <FlexItem key={index}>
                    <Label variant="filled" color="grey">
                      {label.key}={label.value}
                    </Label>
                  </FlexItem>
                ))}
              </Flex>
            </FormGroup>
          )}
        </Form>
      </ModalBody>
      <ModalFooter>
        <Button
          variant="primary"
          onClick={handleSave}
          isLoading={isLoading}
          isDisabled={isLoading}
        >
          Save Changes
        </Button>
        <Button variant="link" onClick={handleClose} isDisabled={isLoading}>
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
}; 