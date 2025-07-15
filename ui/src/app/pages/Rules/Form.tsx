import { css } from '@emotion/css';
import {
  ActionGroup,
  Button,
  Form,
  FormGroup,
  TextInput,
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
  const [inputs, setInputs] = React.useState<IForm>(initialState);
  const [creating, setIsCreating] = React.useState<boolean>(false);
  const [isLabelSelectOpen, setIsLabelSelectOpen] = React.useState<boolean>(false);
  const [labelInputValue, setLabelInputValue] = React.useState<string>('');

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

  const handleChange = (event: React.FormEvent<HTMLInputElement>, value: string) => {
    const name = (event.target as HTMLInputElement).name;
    setInputs((prevState) => ({ ...prevState, [name]: value }));
  };

  const handleLabelInputChange = (_event: React.FormEvent<HTMLInputElement>, value: string) => {
    setLabelInputValue(value);
  };

  const handleLabelSelect = (label: string) => {
    if (!inputs.labels.includes(label)) {
      setInputs((prevState) => ({
        ...prevState,
        labels: [...prevState.labels, label],
      }));
    }
    setLabelInputValue('');
    setIsLabelSelectOpen(false);
  };

  const handleLabelRemove = (labelToRemove: string) => {
    setInputs((prevState) => ({
      ...prevState,
      labels: prevState.labels.filter((label) => label !== labelToRemove),
    }));
  };

  const handleLabelInputKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' && labelInputValue.trim()) {
      event.preventDefault();
      handleLabelSelect(labelInputValue.trim());
    }
  };

  const handleSubmit = async () => {
    setIsCreating(true);
    try {
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

      if (isEditing) {
        await dispatch(updateRule(ruleData));
      } else {
        await dispatch(createRule(ruleData));
      }

      closeFormCB();
    } catch (error) {
      console.error('Error submitting rule:', error);
    } finally {
      setIsCreating(false);
    }
  };

  // Get available label options from the store, filtered by input
  const availableLabelOptions = labels.labels
    .map((label) => `${label.key}=${label.value}`)
    .filter((label) => !inputs.labels.includes(label) && label.toLowerCase().includes(labelInputValue.toLowerCase()));

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
    <Form className={classes.form}>
      <FormGroup label="Name" isRequired fieldId="rule-form-name">
        <TextInput
          isRequired
          type="text"
          id="rule-form-name"
          name="name"
          value={inputs.name}
          onChange={handleChange}
          isDisabled={isEditing} // Name cannot be changed when editing
        />
      </FormGroup>
      <FormGroup label="Pattern" isRequired fieldId="rule-form-pattern">
        <TextInput
          isRequired
          type="text"
          id="rule-form-pattern"
          name="pattern"
          value={inputs.pattern}
          onChange={handleChange}
        />
      </FormGroup>
      <FormGroup label="Labels" isRequired fieldId="rule-form-labels">
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
          {inputs.labels.length > 0 && (
            <FlexItem>
              <Flex spaceItems={{ default: 'spaceItemsXs' }} style={{ marginTop: '8px' }}>
                {inputs.labels.map((label, index) => (
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
      <ActionGroup>
        <Button
          variant="primary"
          isLoading={creating}
          onClick={handleSubmit}
          isDisabled={!inputs.name || !inputs.pattern}
        >
          {isEditing ? 'Update' : 'Create'}
        </Button>
        <Button variant="link" onClick={closeFormCB} isDisabled={creating}>
          Cancel
        </Button>
      </ActionGroup>
    </Form>
  );
};

export { RuleForm };
