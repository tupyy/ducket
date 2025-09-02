import { css } from '@emotion/css';
import {
  EuiForm,
  EuiFormRow,
  EuiFieldText,
  EuiButton,
  EuiFlexGroup,
  EuiFlexItem,
  EuiComboBox,
  EuiComboBoxOptionOption,
  EuiBadge,
  EuiCallOut,
  EuiSpacer,
} from '@elastic/eui';
import * as React from 'react';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { createRule, updateRule } from '@app/shared/reducers/rule.reducer';
import { getLabels } from '@app/shared/reducers/label.reducer';
import { IRule } from '@app/shared/models/rule';

export interface IRuleForm {
  closeFormCB: () => void;
  editingRule?: IRule;
}

const classes = {
  form: css({
    width: '22em',
  }),
};

interface IForm {
  name: string;
  pattern: string;
  labels: string[];
}

const initialState: IForm = {
  name: '',
  pattern: '',
  labels: [],
};

const RuleForm: React.FunctionComponent<IRuleForm> = ({ closeFormCB, editingRule }) => {
  const dispatch = useAppDispatch();
  const labels = useAppSelector((state) => state.labels);
  const rules = useAppSelector((state) => state.rules);
  const [inputs, setInputs] = React.useState<IForm>(initialState);
  const [creating, setIsCreating] = React.useState<boolean>(false);

  const isEditing = !!editingRule;

  React.useEffect(() => {
    dispatch(getLabels());
  }, [dispatch]);

  React.useEffect(() => {
    if (editingRule) {
      setInputs({
        name: editingRule.name,
        pattern: editingRule.pattern,
        labels: editingRule.labels.map((label) => `${label.key}=${label.value}`),
      });
    } else {
      setInputs(initialState);
    }
  }, [editingRule]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setInputs((prevState) => ({ ...prevState, [name]: value }));
  };

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputs((prevState) => ({ ...prevState, name: e.target.value }));
  };

  const handlePatternChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputs((prevState) => ({ ...prevState, pattern: e.target.value }));
  };

  const handleLabelChange = (selectedOptions: EuiComboBoxOptionOption[]) => {
    const selectedLabels = selectedOptions.map((option) => option.label);
    setInputs((prevState) => ({
      ...prevState,
      labels: selectedLabels,
    }));
  };

  const handleSubmit = async () => {
    setIsCreating(true);
    
    const ruleData = {
      name: inputs.name,
      pattern: inputs.pattern,
      labels: inputs.labels.filter((label) => label.length > 0).reduce((acc, label) => {
        const [key, value] = label.split('=');
        if (key && value) {
          acc[key] = value;
        }
        return acc;
      }, {} as { [key: string]: string }),
    };

    try {
      const result = isEditing 
        ? await dispatch(updateRule(ruleData))
        : await dispatch(createRule(ruleData));

      if (updateRule.fulfilled.match(result) || createRule.fulfilled.match(result)) {
        closeFormCB();
      }
    } finally {
      setIsCreating(false);
    }
  };

  // Convert labels to ComboBox options
  const availableLabelOptions: EuiComboBoxOptionOption[] = labels.labels
    .map((label) => ({
      label: `${label.key}=${label.value}`,
      value: `${label.key}=${label.value}`,
    }));

  const selectedLabelOptions: EuiComboBoxOptionOption[] = inputs.labels.map((label) => ({
    label,
    value: label,
  }));

  return (
    <EuiForm className={classes.form}>
      {rules.errorMessage && (
        <>
          <EuiCallOut title="Error" color="danger" iconType="alert">
            {rules.errorMessage}
          </EuiCallOut>
          <EuiSpacer size="m" />
        </>
      )}
      <EuiFormRow label="Name" isInvalid={false} error={[]}>
        <EuiFieldText
          required
          id="rule-form-name"
          name="name"
          value={inputs.name}
          onChange={handleNameChange}
          disabled={isEditing} // Name cannot be changed when editing
        />
      </EuiFormRow>
      
      <EuiFormRow label="Pattern" isInvalid={false} error={[]}>
        <EuiFieldText
          required
          id="rule-form-pattern"
          name="pattern"
          value={inputs.pattern}
          onChange={handlePatternChange}
        />
      </EuiFormRow>
      
      <EuiFormRow label="Labels" isInvalid={false} error={[]}>
        <EuiComboBox
          placeholder="Type to search or add labels (key=value)..."
          options={availableLabelOptions}
          selectedOptions={selectedLabelOptions}
          onChange={handleLabelChange}
          onCreateOption={(searchValue: string) => {
            const newOption = { label: searchValue, value: searchValue };
            handleLabelChange([...selectedLabelOptions, newOption]);
          }}
          isClearable={true}
          data-test-subj="rule-labels-combo-box"
        />
      </EuiFormRow>
      
      <EuiFlexGroup gutterSize="s">
        <EuiFlexItem grow={false}>
          <EuiButton
            fill
            color="primary"
            isLoading={creating}
            onClick={handleSubmit}
            isDisabled={!inputs.name || !inputs.pattern}
          >
            {isEditing ? 'Update' : 'Create'}
          </EuiButton>
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          <EuiButton color="text" onClick={closeFormCB} isDisabled={creating}>
            Cancel
          </EuiButton>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiForm>
  );
};

export { RuleForm };