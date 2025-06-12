import * as React from 'react';
import { CubesIcon } from '@patternfly/react-icons';
import {
  Button,
  Content,
  EmptyState,
  EmptyStateBody,
  EmptyStateFooter,
  EmptyStateVariant,
  PageSection,
} from '@patternfly/react-core';
import { useAppDispatch } from '@app/shared/store';
import { getTags } from '@app/shared/reducers/tag.reducer';

export interface ISupportProps {
  sampleProp?: string;
}

// eslint-disable-next-line prefer-const
let Rules: React.FunctionComponent<ISupportProps> = () => {
  const dispatch = useAppDispatch();

  const getTagsFromProps = () => {
    dispatch(getTags());
  };

  React.useEffect(() => {
    getTagsFromProps();
  }, []);

  const emptyState = (
    <EmptyState variant={EmptyStateVariant.full} titleText="No rules" icon={CubesIcon}>
      <EmptyStateBody>
        <Content>
          <Content component="p">Please add some rules</Content>
        </Content>
      </EmptyStateBody>
      <EmptyStateFooter>
        <Button variant="primary">Add rules</Button>
      </EmptyStateFooter>
    </EmptyState>
  );

  return <PageSection hasBodyWrapper={false}>{emptyState}</PageSection>;
};

export { Rules };
