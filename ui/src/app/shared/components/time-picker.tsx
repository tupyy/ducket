import * as React from 'react';
import { ActionGroup, Button, Form, FormGroup, TextInput } from '@patternfly/react-core';

const TimePicker: React.FC = () => {
  const [value, setValue] = React.useState('');

  const handleValueChange = (_event, val: string) => {
    setValue(val);
  };

  return (
    <Form>
      <FormGroup label="Value" isRequired fieldId="tag-form-value">
        <TextInput isRequired type="text" id="tag-form-value" value={value} onChange={handleValueChange} />
      </FormGroup>
      <ActionGroup>
        <Button variant="secondary" isDisabled={value.length === 0}>
          Submit
        </Button>
        <Button variant="link">
          Cancel
        </Button>
      </ActionGroup>
    </Form>
  );
};

export { TimePicker };
