import * as React from 'react';
import { PageSection } from '@patternfly/react-core';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getRules } from '@app/shared/reducers/rule.reducer';
import { RulesList } from '@app/pages/Rules/List';
import { IRule } from '@app/shared/models/rule';
import { RuleForm } from './Form';

export interface ISupportProps {
  sampleProp?: string;
}

// eslint-disable-next-line prefer-const
let Rules: React.FunctionComponent<ISupportProps> = () => {
  const dispatch = useAppDispatch();
  const rules = useAppSelector((state) => state.rules);
  const [isCreateFormActive, setIsCreateFormActive] = React.useState<boolean>(false);

  const showCreateForm = () => setIsCreateFormActive(true);
  const closeCreateForm = () => setIsCreateFormActive(false);

  React.useEffect(() => {
    dispatch(getRules());
  }, []);

  const renderList = (loading: boolean, rules: Array<IRule>) => {
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
    return <RulesList rules={sortRules(rules.slice())} showCreateRuleFormCB={showCreateForm} />;
  };

  return (
    <PageSection hasBodyWrapper={false}>
      {isCreateFormActive ? <RuleForm closeFormCB={closeCreateForm} /> : renderList(rules.loading, rules.rules)}
    </PageSection>
  );
};

export { Rules };
