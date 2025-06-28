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
  FlexItem
} from '@patternfly/react-core';
import { TimesIcon } from '@patternfly/react-icons';
import * as React from 'react';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { createRule, updateRule } from '@app/shared/reducers/rule.reducer';
import { getTags } from '@app/shared/reducers/tag.reducer';
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
  tags: string[];
}

const initialState: IForm = {
  name: '',
  pattern: '',
  tags: [],
};

const RuleForm: React.FunctionComponent<IRuleForm> = ({ closeFormCB, editingRule }) => {
  const dispatch = useAppDispatch();
  const tags = useAppSelector((state) => state.tags);
  const [inputs, setInputs] = React.useState<IForm>(initialState);
  const [creating, setIsCreating] = React.useState<boolean>(false);
  const [isTagSelectOpen, setIsTagSelectOpen] = React.useState<boolean>(false);
  const [tagInputValue, setTagInputValue] = React.useState<string>('');
  
  const isEditing = !!editingRule;

  React.useEffect(() => {
    dispatch(getTags());
  }, [dispatch]);

  React.useEffect(() => {
    if (editingRule) {
      setInputs({
        name: editingRule.name,
        pattern: editingRule.pattern,
        tags: editingRule.tags.map(tag => tag.value)
      });
    } else {
      setInputs(initialState);
    }
  }, [editingRule]);

  const handleChange = (event: React.FormEvent<HTMLInputElement>, value: string) => {
    const name = (event.target as HTMLInputElement).name;
    setInputs((prevState) => ({ ...prevState, [name]: value }));
  };

  const handleTagInputChange = (_event: React.FormEvent<HTMLInputElement>, value: string) => {
    setTagInputValue(value);
  };

  const handleTagSelect = (tag: string) => {
    if (!inputs.tags.includes(tag)) {
      setInputs((prevState) => ({ 
        ...prevState, 
        tags: [...prevState.tags, tag] 
      }));
    }
    setTagInputValue('');
    setIsTagSelectOpen(false);
  };

  const handleTagRemove = (tagToRemove: string) => {
    setInputs((prevState) => ({ 
      ...prevState, 
      tags: prevState.tags.filter(tag => tag !== tagToRemove) 
    }));
  };

  const handleTagInputKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' && tagInputValue.trim()) {
      event.preventDefault();
      handleTagSelect(tagInputValue.trim());
    }
  };

  const handleSubmit = async () => {
    setIsCreating(true);
    try {
      const ruleData = {
        name: inputs.name,
        pattern: inputs.pattern,
        tags: inputs.tags.filter(tag => tag.length > 0)
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

  // Get available tag options from the store, filtered by input
  const availableTagOptions = tags.tags
    .map(tag => tag.value)
    .filter(tag => 
      !inputs.tags.includes(tag) && 
      tag.toLowerCase().includes(tagInputValue.toLowerCase())
    );

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle 
      ref={toggleRef} 
      onClick={() => setIsTagSelectOpen(!isTagSelectOpen)}
      isExpanded={isTagSelectOpen}
      style={{ width: '100%' }}
    >
      <TextInputGroup>
        <TextInputGroupMain
          value={tagInputValue}
          placeholder="Type to search or add tags..."
          onChange={handleTagInputChange}
          onKeyDown={handleTagInputKeyDown}
        />
        {tagInputValue && (
          <TextInputGroupUtilities>
            <Button
              variant="plain"
              onClick={() => setTagInputValue('')}
              aria-label="Clear input"
            >
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
      <FormGroup label="Tags" isRequired fieldId="rule-form-tags">
        <Flex direction={{ default: 'column' }}>
          <FlexItem>
            <Dropdown
              isOpen={isTagSelectOpen}
              onOpenChange={(isOpen: boolean) => setIsTagSelectOpen(isOpen)}
              toggle={toggle}
              ouiaId="TagDropdown"
              shouldFocusToggleOnSelect
            >
              <DropdownList>
                {availableTagOptions.length > 0 ? (
                  availableTagOptions.map((tag, index) => (
                    <DropdownItem 
                      key={index}
                      value={tag}
                      onClick={() => handleTagSelect(tag)}
                    >
                      {tag}
                    </DropdownItem>
                  ))
                ) : tagInputValue.trim() ? (
                  <DropdownItem 
                    onClick={() => handleTagSelect(tagInputValue.trim())}
                  >
                    Create "{tagInputValue.trim()}"
                  </DropdownItem>
                ) : (
                  <DropdownItem isDisabled>
                    No matching tags found
                  </DropdownItem>
                )}
              </DropdownList>
            </Dropdown>
          </FlexItem>
          {inputs.tags.length > 0 && (
            <FlexItem>
              <Flex spaceItems={{ default: 'spaceItemsXs' }} style={{ marginTop: '8px' }}>
                {inputs.tags.map((tag, index) => (
                  <FlexItem key={index}>
                    <Label 
                      variant="filled" 
                      color="blue"
                      onClose={() => handleTagRemove(tag)}
                      closeBtnAriaLabel={`Remove ${tag} tag`}
                    >
                      {tag}
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
