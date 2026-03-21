import * as React from 'react';
import {
  Form,
  FormGroup,
  TextInput,
  TextInputGroup,
  TextInputGroupMain,
  TextInputGroupUtilities,
  Button,
  ActionGroup,
  Alert,
  Title,
  Label,
  LabelGroup,
  Content,
} from '@patternfly/react-core';
import TimesIcon from '@patternfly/react-icons/dist/esm/icons/times-icon';
import { Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import axios from 'axios';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { createRule, updateRule } from '@app/shared/reducers/rule.reducer';
import { IRule } from '@app/shared/models/rule';
import { ITransaction } from '@app/shared/models/transaction';

interface RuleFormProps {
  onClose: () => void;
  editingRule?: IRule | null;
}

const RuleForm: React.FunctionComponent<RuleFormProps> = ({ onClose, editingRule }) => {
  const dispatch = useAppDispatch();
  const { errorMessage, creating, updating } = useAppSelector((state) => state.rules);
  const [name, setName] = React.useState(editingRule?.name || '');
  const [filter, setFilter] = React.useState(editingRule?.filter || '');
  const [tags, setTags] = React.useState<string[]>(editingRule?.tags || []);
  const [tagInputValue, setTagInputValue] = React.useState('');

  const [testLoading, setTestLoading] = React.useState(false);
  const [testError, setTestError] = React.useState('');
  const [testResults, setTestResults] = React.useState<ITransaction[] | null>(null);
  const [testTotal, setTestTotal] = React.useState(0);

  const isEditing = !!editingRule;
  const isSubmitting = creating || updating;

  const addTag = (value: string) => {
    const trimmed = value.trim();
    if (trimmed && !tags.includes(trimmed)) {
      setTags((prev) => [...prev, trimmed]);
    }
    setTagInputValue('');
  };

  const removeTag = (tag: string) => {
    setTags((prev) => prev.filter((t) => t !== tag));
  };

  const handleTagKeyDown = (event: React.KeyboardEvent) => {
    if ((event.key === ' ' || event.key === 'Enter' || event.key === ',') && tagInputValue.trim()) {
      event.preventDefault();
      addTag(tagInputValue);
    }
    if (event.key === 'Backspace' && !tagInputValue && tags.length > 0) {
      setTags((prev) => prev.slice(0, -1));
    }
  };

  const handleSubmit = async () => {
    const result = isEditing
      ? await dispatch(updateRule({ id: editingRule!.id, name, filter, tags }))
      : await dispatch(createRule({ name, filter, tags }));

    if (updateRule.fulfilled.match(result) || createRule.fulfilled.match(result)) {
      onClose();
    }
  };

  const handleTestRule = async () => {
    if (!filter) return;
    setTestLoading(true);
    setTestError('');
    setTestResults(null);
    try {
      const params = new URLSearchParams();
      params.append('filter', filter);
      params.append('limit', '10');
      const res = await axios.get<{ items: ITransaction[]; total: number }>(
        `api/v1/transactions?${params.toString()}`,
      );
      setTestResults(res.data.items);
      setTestTotal(res.data.total);
    } catch (err: any) {
      setTestError(err?.response?.data?.error || err?.message || 'Failed to test rule');
    } finally {
      setTestLoading(false);
    }
  };

  const formatDate = (dateStr: string) => {
    try {
      const d = new Date(dateStr);
      const day = String(d.getDate()).padStart(2, '0');
      const month = String(d.getMonth() + 1).padStart(2, '0');
      const year = d.getFullYear();
      return `${day}-${month}-${year}`;
    } catch {
      return dateStr;
    }
  };

  return (
    <>
      <Title headingLevel="h1" size="lg" style={{ marginBottom: '1rem' }}>
        {isEditing ? 'Edit Rule' : 'Create Rule'}
      </Title>

      {errorMessage && (
        <Alert variant="danger" title="Error" isInline style={{ marginBottom: '1rem' }}>
          {errorMessage}
        </Alert>
      )}

      <Form style={{ maxWidth: '500px' }}>
        <FormGroup label="Name" isRequired fieldId="rule-name">
          <TextInput
            isRequired
            id="rule-name"
            value={name}
            onChange={(_evt, val) => setName(val)}
            isDisabled={isEditing}
          />
        </FormGroup>

        <FormGroup label="Filter" isRequired fieldId="rule-filter">
          <TextInput
            isRequired
            id="rule-filter"
            value={filter}
            onChange={(_evt, val) => { setFilter(val); setTestResults(null); }}
            placeholder="e.g. content ~ /kaufland/"
          />
        </FormGroup>

        <FormGroup label="Tags" fieldId="rule-tags">
          <div style={{ border: '1px solid var(--pf-t--global--border--color--default)', borderRadius: '3px', padding: '4px 8px', display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: '4px', minHeight: '36px' }}>
            {tags.map((tag) => (
              <Label key={tag} color="teal" onClose={() => removeTag(tag)}>
                {tag}
              </Label>
            ))}
            <input
              id="rule-tags"
              value={tagInputValue}
              onChange={(e) => setTagInputValue(e.target.value)}
              onKeyDown={handleTagKeyDown}
              onBlur={() => { if (tagInputValue.trim()) addTag(tagInputValue); }}
              placeholder={tags.length === 0 ? 'Type a tag and press space...' : ''}
              style={{ border: 'none', outline: 'none', flex: 1, minWidth: '120px', background: 'transparent' }}
            />
          </div>
        </FormGroup>

        <ActionGroup>
          <Button
            variant="primary"
            onClick={handleSubmit}
            isLoading={isSubmitting}
            isDisabled={!name || !filter || isSubmitting}
          >
            {isEditing ? 'Update' : 'Create'}
          </Button>
          <Button
            variant="secondary"
            onClick={handleTestRule}
            isLoading={testLoading}
            isDisabled={!filter || testLoading}
          >
            Test Rule
          </Button>
          <Button variant="link" onClick={onClose} isDisabled={isSubmitting}>
            Cancel
          </Button>
        </ActionGroup>
      </Form>

      {testError && (
        <Alert variant="danger" title="Test failed" isInline style={{ marginTop: '1rem', maxWidth: '700px' }}>
          {testError}
        </Alert>
      )}

      {testResults !== null && (
        <div style={{ marginTop: '1rem', maxWidth: '700px' }}>
          <Alert
            variant={testTotal > 0 ? 'success' : 'warning'}
            title={`${testTotal} transaction${testTotal !== 1 ? 's' : ''} matched`}
            isInline
            style={{ marginBottom: '0.5rem' }}
          />

          {testResults.length > 0 && (
            <Table aria-label="Test results" variant="compact">
              <Thead>
                <Tr>
                  <Th>Date</Th>
                  <Th>Type</Th>
                  <Th>Content</Th>
                  <Th>Amount</Th>
                </Tr>
              </Thead>
              <Tbody>
                {testResults.map((txn) => (
                  <Tr key={txn.id}>
                    <Td>{formatDate(txn.date)}</Td>
                    <Td>
                      <Label color={txn.kind === 'debit' ? 'red' : 'green'}>{txn.kind}</Label>
                    </Td>
                    <Td>
                      <Content component="p" style={{ maxWidth: '300px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                        {txn.content}
                      </Content>
                    </Td>
                    <Td style={{ fontFamily: 'monospace' }}>€{txn.amount.toFixed(2)}</Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          )}

          {testTotal > 10 && (
            <Content component="p" style={{ marginTop: '0.5rem', fontStyle: 'italic' }}>
              Showing 10 of {testTotal} matches
            </Content>
          )}
        </div>
      )}
    </>
  );
};

export { RuleForm };
