import { CreateTagError, Result } from '@app/shared/models/result';
import { css } from '@emotion/css';
import { ActionGroup, Button, Form, FormGroup, TextInput } from '@patternfly/react-core';
import * as React from 'react';

export interface ITagForm {
  closeFormCB: () => void;
  creating: boolean;
  submitValue: (value: string) => void;
  createResult: Result<string, CreateTagError>;
}

const classes = {
  form: css({
    width: '22em',
  }),
};

const TagForm: React.FunctionComponent<ITagForm> = ({ closeFormCB, submitValue, creating }) => {
  const [value, setValue] = React.useState('');

  const handleValueChange = (_event, val: string) => {
    setValue(val);
  };

  const handleSubmit = () => {
    submitValue(value);
  };

  return (
    <Form className={classes.form}>
      <FormGroup label="Value" isRequired fieldId="tag-form-value">
        <TextInput isRequired type="text" id="tag-form-value" value={value} onChange={handleValueChange} />
      </FormGroup>
      <ActionGroup>
        <Button variant="secondary" isLoading={creating} onClick={handleSubmit} isDisabled={value.length === 0}>
          Submit
        </Button>
        <Button variant="link" onClick={closeFormCB} isLoading={creating}>
          Cancel
        </Button>
      </ActionGroup>
    </Form>
  );
};

export { TagForm };
