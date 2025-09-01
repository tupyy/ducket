import { CreateLabelError, Result } from '@app/shared/models/result';
import { css } from '@emotion/css';
import { EuiForm, EuiFormRow, EuiFieldText, EuiButton, EuiFlexGroup, EuiFlexItem } from '@elastic/eui';
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

  const handleKeyChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setKey(e.target.value);
  };

  const handleValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  const handleSubmit = () => {
    submitValue(key, value);
  };

  return (
    <EuiForm className={classes.form}>
      <EuiFormRow label="Key" isInvalid={false} error={[]}>
        <EuiFieldText
          required
          id="label-form-key"
          value={key}
          onChange={handleKeyChange}
          placeholder="e.g., category"
        />
      </EuiFormRow>
      <EuiFormRow label="Value" isInvalid={false} error={[]}>
        <EuiFieldText
          required
          id="label-form-value"
          value={value}
          onChange={handleValueChange}
          placeholder="e.g., food"
        />
      </EuiFormRow>
      <EuiFlexGroup gutterSize="s">
        <EuiFlexItem grow={false}>
          <EuiButton
            fill
            color="primary"
            isLoading={creating}
            onClick={handleSubmit}
            isDisabled={key.length === 0 || value.length === 0}
          >
            Submit
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

export { LabelForm };