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
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { getTags } from '@app/shared/reducers/tag.reducer';
import { ITag } from '@app/shared/models/tag';

export interface ISupportProps {
  sampleProp?: string;
}

// eslint-disable-next-line prefer-const
const Tags = () => {
  const dispatch = useAppDispatch();

  const getTagsFromProps = () => {
    dispatch(getTags());
  };

  React.useEffect(() => {
    getTagsFromProps();
  }, []);

  const tags = useAppSelector((state) => state.tags.tags);
  const count = useAppSelector((state) => state.tags.total);

  const emptyState = (
    <EmptyState variant={EmptyStateVariant.full} titleText="No tags" icon={CubesIcon}>
      <EmptyStateBody>
        <Content>
          <Content component="p">Please add some tags</Content>
        </Content>
      </EmptyStateBody>
      <EmptyStateFooter>
        <Button variant="primary">Add tag</Button>
      </EmptyStateFooter>
    </EmptyState>
  );

  const renderTagList = (
    <div>
      {tags.map((tag: ITag, i: number) => (
        <div key={`tag-${i}`}>{tag.value}</div>
      ))}
    </div>
  );

  return <PageSection hasBodyWrapper={false}>{
    count == 0 ? emptyState: renderTagList
  }</PageSection>;
};

export { Tags };
