import * as React from 'react';
import { PageSection } from '@patternfly/react-core';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getRules, deleteRule } from '@app/shared/reducers/rule.reducer';
import { RulesList } from './List';
import { RuleForm } from './Form';
import { IRule } from '@app/shared/models/rule';

const Rules: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const { rules, loading, creating, updating } = useAppSelector((state) => state.rules);
  const [editingRule, setEditingRule] = React.useState<IRule | null>(null);
  const [showForm, setShowForm] = React.useState(false);

  React.useEffect(() => {
    dispatch(getRules());
  }, [dispatch]);

  const handleCreate = () => {
    setEditingRule(null);
    setShowForm(true);
  };

  const handleEdit = (rule: IRule) => {
    setEditingRule(rule);
    setShowForm(true);
  };

  const handleCloseForm = () => {
    setShowForm(false);
    setEditingRule(null);
  };

  const handleDelete = (id: number) => {
    dispatch(deleteRule(id));
  };

  return (
    <PageSection>
      {showForm ? (
        <RuleForm onClose={handleCloseForm} editingRule={editingRule} />
      ) : (
        <RulesList
          rules={rules}
          loading={loading}
          onCreateRule={handleCreate}
          onEditRule={handleEdit}
          onDeleteRule={handleDelete}
        />
      )}
    </PageSection>
  );
};

export { Rules };
