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
  TextInput,
  TextArea,
  Alert,
  MenuToggle,
  MenuToggleElement,
  Dropdown,
  DropdownList,
  DropdownItem,
  Label,
  TextInputGroup,
  TextInputGroupMain,
  TextInputGroupUtilities,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import { TimesIcon } from '@patternfly/react-icons';
import { ITransaction } from '@app/shared/models/transaction';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { createRule } from '@app/shared/reducers/rule.reducer';
import { getLabels } from '@app/shared/reducers/label.reducer';

interface CreateRuleModalProps {
  isOpen: boolean;
  onClose: () => void;
  transaction?: ITransaction;
}

export const CreateRuleModal: React.FC<CreateRuleModalProps> = ({
  isOpen,
  onClose,
  transaction,
}) => {
  const dispatch = useAppDispatch();
  const labels = useAppSelector((state) => state.labels);
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  
  // Form state
  const [ruleName, setRuleName] = React.useState('');
  const [ruleDescription, setRuleDescription] = React.useState('');
  const [pattern, setPattern] = React.useState('');
  const [selectedLabels, setSelectedLabels] = React.useState<string[]>([]);
  
  // Label dropdown state
  const [isLabelSelectOpen, setIsLabelSelectOpen] = React.useState<boolean>(false);
  const [labelInputValue, setLabelInputValue] = React.useState<string>('');

  // Fetch labels when component mounts
  React.useEffect(() => {
    if (isOpen) {
      dispatch(getLabels());
    }
  }, [dispatch, isOpen]);

  // Initialize form when transaction changes
  React.useEffect(() => {
    if (transaction && isOpen) {
      setRuleName(`Rule for ${transaction.description.substring(0, 30)}...`);
      setRuleDescription(`Auto-generated rule for transactions similar to: ${transaction.description}`);
      setPattern(transaction.description);
      setSelectedLabels([]);
      setError(null);
    }
  }, [transaction, isOpen]);

  // Reset form when modal closes
  React.useEffect(() => {
    if (!isOpen) {
      setRuleName('');
      setRuleDescription('');
      setPattern('');
      setSelectedLabels([]);
      setLabelInputValue('');
      setError(null);
      setIsLoading(false);
      setIsLabelSelectOpen(false);
    }
  }, [isOpen]);

  // Label handling functions
  const handleLabelInputChange = (_event: React.FormEvent<HTMLInputElement>, value: string) => {
    setLabelInputValue(value);
  };

  const handleLabelSelect = (label: string) => {
    if (!selectedLabels.includes(label)) {
      setSelectedLabels((prev) => [...prev, label]);
    }
    setLabelInputValue('');
    setIsLabelSelectOpen(false);
  };

  const handleLabelRemove = (labelToRemove: string) => {
    setSelectedLabels((prev) => prev.filter((label) => label !== labelToRemove));
  };

  const handleLabelInputKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' && labelInputValue.trim()) {
      event.preventDefault();
      handleLabelSelect(labelInputValue.trim());
    }
  };

  const handleSubmit = async () => {
    if (!transaction) return;

    setIsLoading(true);
    setError(null);

    try {
      // Convert labels array to object format expected by API
      const labelsObject = selectedLabels.filter((label) => label.length > 0).reduce((acc, label) => {
        const [key, value] = label.split('=');
        if (key && value) {
          acc[key] = value;
        }
        return acc;
      }, {} as { [key: string]: string });

      await dispatch(createRule({
        name: ruleName.trim(),
        pattern: pattern.trim(),
        labels: labelsObject,
      }));
      
      onClose();
    } catch (err: any) {
      setError(err.message || 'Failed to create rule');
    } finally {
      setIsLoading(false);
    }
  };

  // Get available label options from the store, filtered by input
  const availableLabelOptions = labels.labels
    .map((label) => `${label.key}=${label.value}`)
    .filter((label) => !selectedLabels.includes(label) && label.toLowerCase().includes(labelInputValue.toLowerCase()));

  const isFormValid = ruleName.trim() && pattern.trim();

  const labelToggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setIsLabelSelectOpen(!isLabelSelectOpen)}
      isExpanded={isLabelSelectOpen}
      style={{ width: '100%' }}
    >
      <TextInputGroup>
        <TextInputGroupMain
          value={labelInputValue}
          placeholder="Type to search or add labels (key=value)..."
          onChange={handleLabelInputChange}
          onKeyDown={handleLabelInputKeyDown}
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
      onClose={onClose}
    >
      <ModalHeader title="Create Rule from Transaction" />
      <ModalBody>
        {error && (
          <Alert
            variant="danger"
            title="Error creating rule"
            style={{ marginBottom: '1rem' }}
          >
            {error}
          </Alert>
        )}
        
        {transaction && (
          <div style={{ marginBottom: '1rem' }}>
            <strong>Based on transaction:</strong> {transaction.description}
          </div>
        )}

        <Form>
          <FormGroup label="Rule Name" isRequired fieldId="rule-name">
            <TextInput
              id="rule-name"
              type="text"
              value={ruleName}
              onChange={(_event, value) => setRuleName(value)}
              placeholder="Enter rule name"
            />
          </FormGroup>

          <FormGroup label="Description" fieldId="rule-description">
            <TextArea
              id="rule-description"
              value={ruleDescription}
              onChange={(_event, value) => setRuleDescription(value)}
              placeholder="Enter rule description"
              rows={3}
            />
          </FormGroup>

          <FormGroup label="Pattern" isRequired fieldId="rule-pattern">
            <TextInput
              id="rule-pattern"
              type="text"
              value={pattern}
              onChange={(_event, value) => setPattern(value)}
              placeholder="Enter pattern to match transactions"
            />
          </FormGroup>

          <FormGroup label="Labels" fieldId="rule-labels">
            <Flex direction={{ default: 'column' }}>
              <FlexItem>
                <Dropdown
                  isOpen={isLabelSelectOpen}
                  onOpenChange={(isOpen: boolean) => setIsLabelSelectOpen(isOpen)}
                  toggle={labelToggle}
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
                    ) : labelInputValue.trim() ? (
                      <DropdownItem onClick={() => handleLabelSelect(labelInputValue.trim())}>
                        Create "{labelInputValue.trim()}"
                      </DropdownItem>
                    ) : (
                      <DropdownItem isDisabled>No matching labels found</DropdownItem>
                    )}
                  </DropdownList>
                </Dropdown>
              </FlexItem>
              {selectedLabels.length > 0 && (
                <FlexItem>
                  <Flex spaceItems={{ default: 'spaceItemsXs' }} style={{ marginTop: '8px' }}>
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
          </FormGroup>
        </Form>
      </ModalBody>
      <ModalFooter>
        <Button
          variant="primary"
          onClick={handleSubmit}
          isDisabled={!isFormValid || isLoading}
          isLoading={isLoading}
        >
          Create Rule
        </Button>
        <Button variant="link" onClick={onClose} isDisabled={isLoading}>
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
}; 