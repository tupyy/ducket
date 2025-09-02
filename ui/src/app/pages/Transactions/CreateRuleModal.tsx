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
  EuiFieldText,
  EuiTextArea,
  EuiComboBox,
  EuiComboBoxOptionOption,
  EuiCallOut,
  EuiText,
  EuiSpacer,
} from '@elastic/eui';
import { ITransaction } from '@app/shared/models/transaction';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { createRule, reset } from '@app/shared/reducers/rule.reducer';
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
  const rules = useAppSelector((state) => state.rules);
  const [error, setError] = React.useState<string | null>(null);
  
  // Form state
  const [ruleName, setRuleName] = React.useState('');
  const [ruleDescription, setRuleDescription] = React.useState('');
  const [pattern, setPattern] = React.useState('');
  const [selectedLabels, setSelectedLabels] = React.useState<Array<{ key: string; value: string }>>([]);

  // Fetch labels when component mounts
  React.useEffect(() => {
    if (isOpen) {
      dispatch(getLabels());
      dispatch(reset()); // Clear any previous errors
      // Initialize form with transaction data if available
      if (transaction) {
        setPattern(transaction.description || '');
        setRuleName(`Rule for ${transaction.description?.substring(0, 20)}...`);
      }
    }
  }, [isOpen, transaction, dispatch]);

  // Reset form when modal closes
  React.useEffect(() => {
    if (!isOpen) {
      setRuleName('');
      setRuleDescription('');
      setPattern('');
      setSelectedLabels([]);
      setError(null);
    }
  }, [isOpen]);

  const handleSubmit = async () => {
    if (!ruleName.trim() || !pattern.trim()) {
      setError('Rule name and pattern are required');
      return;
    }

    setError(null);

    const ruleData = {
      name: ruleName.trim(),
      pattern: pattern.trim(),
      labels: selectedLabels.reduce((acc, label) => {
        acc[label.key] = label.value;
        return acc;
      }, {} as { [key: string]: string }),
    };

    const result = await dispatch(createRule(ruleData));
    if (createRule.fulfilled.match(result)) {
      onClose();
    }
  };

  const handleLabelChange = (selectedOptions: EuiComboBoxOptionOption[]) => {
    const newLabels = selectedOptions.map(option => {
      const [key, value] = option.label.split('=', 2);
      return { key: key.trim(), value: value.trim() };
    });
    setSelectedLabels(newLabels);
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
    <EuiModal onClose={onClose} style={{ width: '500px' }}>
      <EuiModalHeader>
        <EuiModalHeaderTitle>Create Rule</EuiModalHeaderTitle>
      </EuiModalHeader>

      <EuiModalBody>
        {transaction && (
          <>
            <EuiText size="s" color="subdued">
              Create a rule based on transaction: {transaction.description}
            </EuiText>
            <EuiSpacer size="m" />
          </>
        )}

        {(error || rules.errorMessage) && (
          <>
            <EuiCallOut title="Error" color="danger" iconType="alert">
              {error || rules.errorMessage}
            </EuiCallOut>
            <EuiSpacer size="m" />
          </>
        )}

        <EuiForm>
          <EuiFormRow label="Rule Name" isInvalid={false} error={[]} fullWidth>
            <EuiFieldText
              value={ruleName}
              onChange={(e) => setRuleName(e.target.value)}
              placeholder="Enter rule name"
              fullWidth
            />
          </EuiFormRow>

          <EuiFormRow label="Pattern" isInvalid={false} error={[]} fullWidth>
            <EuiFieldText
              value={pattern}
              onChange={(e) => setPattern(e.target.value)}
              placeholder="Enter pattern to match transactions"
              fullWidth
            />
          </EuiFormRow>

          <EuiFormRow label="Description (Optional)" fullWidth>
            <EuiTextArea
              value={ruleDescription}
              onChange={(e) => setRuleDescription(e.target.value)}
              placeholder="Enter rule description"
              rows={3}
              fullWidth
            />
          </EuiFormRow>

          <EuiFormRow label="Labels to Apply" fullWidth>
            <EuiComboBox
              placeholder="Select or create labels (key=value format)"
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
        </EuiForm>
      </EuiModalBody>

      <EuiModalFooter>
        <EuiButton onClick={onClose} isDisabled={rules.creating}>
          Cancel
        </EuiButton>
        <EuiButton
          fill
          color="primary"
          onClick={handleSubmit}
          isLoading={rules.creating}
          isDisabled={rules.creating || !ruleName.trim() || !pattern.trim()}
        >
          Create Rule
        </EuiButton>
      </EuiModalFooter>
    </EuiModal>
  );
};