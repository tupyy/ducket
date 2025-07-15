import { CreateLabelError, Result } from '@app/shared/models/result';
import { css } from '@emotion/css';
import { ActionGroup, Button, Form, FormGroup, TextInput } from '@patternfly/react-core';
import * as React from 'react';

export interface ILabelForm {
  closeFormCB: () => void;
  creating: boolean;
  submitValue: (key: string, value: string) => void;
  createResult: Result<string, CreateLabelError>;
}

const classes = {
  form: css({
    width: '22em',
  }),
};

const LabelForm: React.FunctionComponent<ILabelForm> = ({ closeFormCB, submitValue, creating }) => {
  const [key, setKey] = React.useState('');
  const [value, setValue] = React.useState('');

  const handleKeyChange = (_event, val: string) => {
    setKey(val);
  };

  const handleValueChange = (_event, val: string) => {
    setValue(val);
  };

  const handleSubmit = () => {
    submitValue(key, value);
  };

  return (
    <Form className={classes.form}>
      <FormGroup label="Key" isRequired fieldId="label-form-key">
        <TextInput isRequired type="text" id="label-form-key" value={key} onChange={handleKeyChange} placeholder="e.g., category" />
      </FormGroup>
      <FormGroup label="Value" isRequired fieldId="label-form-value">
        <TextInput isRequired type="text" id="label-form-value" value={value} onChange={handleValueChange} placeholder="e.g., food" />
      </FormGroup>
      <ActionGroup>
        <Button variant="secondary" isLoading={creating} onClick={handleSubmit} isDisabled={key.length === 0 || value.length === 0}>
          Submit
        </Button>
        <Button variant="link" onClick={closeFormCB} isLoading={creating}>
          Cancel
        </Button>
      </ActionGroup>
    </Form>
  );
};

export { LabelForm };
