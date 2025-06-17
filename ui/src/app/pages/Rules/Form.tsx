import { css } from '@emotion/css';
import { ActionGroup, Button, Form, FormGroup, TextInput } from '@patternfly/react-core';
import * as React from 'react';

export interface IRuleForm {
  closeFormCB: () => void;
}

const classes = {
  form: css({
    width: '22em',
  }),
};

interface IForm {
  name: string;
  pattern: string;
  tags: string;
}

const initialState: IForm = {
  name: '',
  pattern: '',
  tags: '',
};

const RuleForm: React.FunctionComponent<IRuleForm> = ({ closeFormCB }) => {
  const [inputs, setInputs] = React.useState<IForm>(initialState);
  const [creating, setIsCreating] = React.useState<boolean>(false);
  const handleChange = (e) => setInputs((prevState) => ({ ...prevState, [e.target.name]: e.target.value }));

  const handleSubmit = () => {};

  return (
    <Form className={classes.form}>
      <FormGroup label="Value" isRequired fieldId="tag-form-value">
        <TextInput isRequired type="text" id="tag-form-value" value={inputs.name} onChange={handleChange} />
      </FormGroup>
      <FormGroup label="Tags" isRequired fieldId="tags-form-value">
        <TextInput isRequired type="text" id="tags-form-value" value={inputs.tags} onChange={handleChange} />
      </FormGroup>
      <ActionGroup>
        <Button variant="secondary" isLoading={creating} onClick={handleSubmit} isDisabled={!inputs.name}>
          Submit
        </Button>
        <Button variant="link" onClick={closeFormCB} isLoading={creating}>
          Cancel
        </Button>
      </ActionGroup>
    </Form>
  );
};

export { RuleForm };
