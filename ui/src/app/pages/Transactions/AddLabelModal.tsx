import React, { useState } from 'react';
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
  Label,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import { TimesIcon } from '@patternfly/react-icons';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { addLabelToTransaction } from '@app/shared/reducers/transaction.reducer';
import { getLabels } from '@app/shared/reducers/label.reducer';

interface AddLabelModalProps {
  isOpen: boolean;
  onClose: () => void;
  transactionHref: string;
  transactionDescription?: string;
}

export const AddLabelModal: React.FunctionComponent<AddLabelModalProps> = ({
  isOpen,
  onClose,
  transactionHref,
  transactionDescription = '',
}) => {
  const dispatch = useAppDispatch();
  const { addingLabel, addLabelSuccess, errorMessage } = useAppSelector((state) => state.transactions);
  const labels = useAppSelector((state) => state.labels);

  const [labelInputValue, setLabelInputValue] = useState('');
  const [isLabelSelectOpen, setIsLabelSelectOpen] = useState(false);
  const [selectedLabels, setSelectedLabels] = useState<string[]>([]);
  const [validationError, setValidationError] = useState('');

  // Load available labels when modal opens
  React.useEffect(() => {
    if (isOpen) {
      dispatch(getLabels());
    }
  }, [isOpen, dispatch]);

  // Reset form when modal opens/closes
  React.useEffect(() => {
    if (isOpen) {
      setLabelInputValue('');
      setSelectedLabels([]);
      setValidationError('');
      setIsLabelSelectOpen(false);
    }
  }, [isOpen]);

  // Close modal when labels are successfully added
  React.useEffect(() => {
    if (addLabelSuccess) {
      handleClose();
    }
  });

  const validateLabelInput = (input: string) => {
    const trimmedInput = input.trim();
    if (!trimmedInput) {
      return 'Label is required';
    }

    const parts = trimmedInput.split('=');
    if (parts.length !== 2 || !parts[0].trim() || !parts[1].trim()) {
      return 'Label must be in format: key=value';
    }

    return '';
  };

  const handleSubmit = async (event?: React.FormEvent<HTMLFormElement>) => {
    if (event) {
      event.preventDefault();
    }

    if (selectedLabels.length === 0) {
      setValidationError('Please add at least one label');
      return;
    }

    // Submit all selected labels
    for (const labelString of selectedLabels) {
      const [key, value] = labelString.split('=');
      await dispatch(
        addLabelToTransaction({
          transactionHref,
          key: key.trim(),
          value: value.trim(),
        })
      );
    }
  };

  const handleClose = () => {
    if (!addingLabel) {
      onClose();
    }
  };

  const handleLabelInputChange = (_event: React.FormEvent<HTMLInputElement>, value: string) => {
    setLabelInputValue(value);
    if (validationError && value.trim()) {
      setValidationError('');
    }
  };

  const handleLabelSelect = (label: string) => {
    if (!selectedLabels.includes(label)) {
      setSelectedLabels((prev) => [...prev, label]);
    }
    setLabelInputValue('');
    setIsLabelSelectOpen(false);
    if (validationError) {
      setValidationError('');
    }
  };

  const handleLabelRemove = (labelToRemove: string) => {
    setSelectedLabels((prev) => prev.filter((label) => label !== labelToRemove));
  };

  const handleLabelInputKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' && labelInputValue.trim()) {
      event.preventDefault();
      const error = validateLabelInput(labelInputValue);
      if (!error) {
        handleLabelSelect(labelInputValue.trim());
      } else {
        setValidationError(error);
      }
    }
  };

  // Get available label options from the store, filtered by input and not already selected
  const availableLabelOptions = labels.labels
    .map((label) => `${label.key}=${label.value}`)
    .filter((label) => !selectedLabels.includes(label) && label.toLowerCase().includes(labelInputValue.toLowerCase()))
    .slice(0, 10); // Limit to 10 suggestions

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setIsLabelSelectOpen(!isLabelSelectOpen)}
      isExpanded={isLabelSelectOpen}
      style={{ width: '100%' }}
    >
      <TextInputGroup>
        <TextInputGroupMain
          value={labelInputValue}
          placeholder="Type label in format: key=value (e.g., category=food)"
          onChange={handleLabelInputChange}
          onKeyDown={handleLabelInputKeyDown}
          disabled={addingLabel}
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

  const formId = 'add-labels-form';

  return (
    <Modal
      variant={ModalVariant.small}
      isOpen={isOpen}
      onClose={handleClose}
      aria-labelledby="add-labels-modal-title"
      aria-describedby="add-labels-modal-description"
    >
      <ModalHeader
        title="Add Custom Labels"
        labelId="add-labels-modal-title"
        description={transactionDescription ? `Transaction: ${transactionDescription}` : undefined}
        descriptorId="add-labels-modal-description"
      />
      <ModalBody>
        {errorMessage && (
          <Alert variant={AlertVariant.danger} title="Error" style={{ marginBottom: '1rem' }}>
            {errorMessage}
          </Alert>
        )}

        <Form id={formId} onSubmit={handleSubmit}>
          <FormGroup label="Labels" isRequired fieldId="label-input">
            <Flex direction={{ default: 'column' }}>
              <FlexItem>
                <Dropdown
                  isOpen={isLabelSelectOpen}
                  onOpenChange={(isOpen: boolean) => setIsLabelSelectOpen(isOpen)}
                  toggle={toggle}
                  ouiaId="LabelDropdown"
                  shouldFocusToggleOnSelect
                >
                  <DropdownList>
                    {availableLabelOptions.length > 0 ? (
                      availableLabelOptions.map((label, index) => (
                        <DropdownItem key={index} value={label} onClick={() => handleLabelSelect(label)}>
                          {label}
                        </DropdownItem>
                      ))
                    ) : labelInputValue.trim() &&
                      !validateLabelInput(labelInputValue) &&
                      !selectedLabels.includes(labelInputValue.trim()) ? (
                      <DropdownItem onClick={() => handleLabelSelect(labelInputValue.trim())}>
                        Add "{labelInputValue.trim()}"
                      </DropdownItem>
                    ) : (
                      <DropdownItem isDisabled>
                        {labelInputValue.trim()
                          ? selectedLabels.includes(labelInputValue.trim())
                            ? 'Label already added'
                            : 'Invalid format. Use: key=value'
                          : 'Type to search existing labels or create new one'}
                      </DropdownItem>
                    )}
                  </DropdownList>
                </Dropdown>
              </FlexItem>
              {selectedLabels.length > 0 && (
                <FlexItem>
                  <Flex spaceItems={{ default: 'spaceItemsXs' }} style={{ marginTop: '0.5rem' }}>
                    {selectedLabels.map((label, index) => (
                      <FlexItem key={index}>
                        <Label
                          variant="filled"
                          color="blue"
                          onClose={() => handleLabelRemove(label)}
                          closeBtnAriaLabel={`Remove ${label} label`}
                        >
                          {label}
                        </Label>
                      </FlexItem>
                    ))}
                  </Flex>
                </FlexItem>
              )}
            </Flex>
            {validationError && (
              <div
                style={{ color: 'var(--pf-v6-global--danger-color--100)', fontSize: '0.875rem', marginTop: '0.25rem' }}
              >
                {validationError}
              </div>
            )}
            <div style={{ color: 'var(--pf-v6-global--Color--200)', fontSize: '0.875rem', marginTop: '0.25rem' }}>
              Enter labels in key=value format (e.g., category=food, type=expense). You can add multiple labels.
            </div>
          </FormGroup>
        </Form>
      </ModalBody>
      <ModalFooter>
        <Button
          variant="primary"
          form={formId}
          onClick={() => handleSubmit()}
          isLoading={addingLabel}
          isDisabled={addingLabel || selectedLabels.length === 0}
        >
          Add {selectedLabels.length} Label{selectedLabels.length !== 1 ? 's' : ''}
        </Button>
        <Button variant="link" onClick={handleClose} isDisabled={addingLabel}>
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
};
