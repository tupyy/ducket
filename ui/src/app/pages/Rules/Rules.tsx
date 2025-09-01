import * as React from 'react';
import { EuiPageSection } from '@elastic/eui';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getRules, syncRule, syncAllRules, deleteRule } from '@app/shared/reducers/rule.reducer';
import { RulesList } from '@app/pages/Rules/List';
import { IRule } from '@app/shared/models/rule';
import { RuleForm } from './Form';

const Rules: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const rules = useAppSelector((state) => state.rules);
  const [isCreateFormActive, setIsCreateFormActive] = React.useState<boolean>(false);
  const [isEditFormActive, setIsEditFormActive] = React.useState<boolean>(false);
  const [editingRule, setEditingRule] = React.useState<IRule | null>(null);

  const showCreateForm = () => {
    setIsCreateFormActive(true);
    setIsEditFormActive(false);
    setEditingRule(null);
  };

  const showEditForm = (rule: IRule) => {
    setIsEditFormActive(true);
    setIsCreateFormActive(false);
    setEditingRule(rule);
  };

  const closeForm = () => {
    setIsCreateFormActive(false);
    setIsEditFormActive(false);
    setEditingRule(null);
  };

  const handleSyncRule = (ruleName: string) => {
    dispatch(syncRule(ruleName));
  };

  const handleSyncAllRules = () => {
    dispatch(syncAllRules());
  };

  const handleDeleteRule = (ruleName: string) => {
    dispatch(deleteRule(ruleName));
  };

  React.useEffect(() => {
    dispatch(getRules());
  }, [dispatch]);

  const renderList = (loading: boolean, rulesArray: Array<IRule>) => {
    const sortRules = (rules: Array<IRule> | []) => {
      return rules.sort((rule1, rule2) => {
        if (rule1.name < rule2.name) {
          return -1;
        }
        if (rule1.name > rule2.name) {
          return 1;
        }
        return 0;
      });
    };
    return (
      <RulesList
        rules={sortRules(rulesArray.slice())}
        showCreateRuleFormCB={showCreateForm}
        showEditRuleFormCB={showEditForm}
        onSyncRule={handleSyncRule}
        onSyncAllRules={handleSyncAllRules}
        onDeleteRule={handleDeleteRule}
        syncing={rules.syncing}
        syncingAll={rules.syncingAll}
      />
    );
  };

  const renderForm = () => {
    if (isCreateFormActive) {
      return <RuleForm closeFormCB={closeForm} />;
    }
    if (isEditFormActive && editingRule) {
      return <RuleForm closeFormCB={closeForm} editingRule={editingRule} />;
    }
    return null;
  };

  return (
    <EuiPageSection>
      {isCreateFormActive || isEditFormActive ? renderForm() : renderList(rules.loading, rules.rules)}
    </EuiPageSection>
  );
};

export { Rules };
