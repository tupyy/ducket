import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { PageSection } from '@patternfly/react-core';
import * as React from 'react';
import { TagsList } from './Tags';
import { getTags } from '@app/shared/reducers/tag.reducer';

const TagsPage: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const tags = useAppSelector((state) => state.tags.tags);
  const [isCreateFormActive, setIsCreateFormActive] = React.useState<boolean>(false);

  const getTagsFromProps = () => {
    dispatch(getTags());
  };

  const showCreateForm = () => {
    setIsCreateFormActive(true);
  };

  React.useEffect(() => {
    getTagsFromProps();
  }, []);

  return <PageSection>{isCreateFormActive ? <div>form</div> : <TagsList tags={tags} showCreateTagFormCB={showCreateForm} />}</PageSection>;
};

export { TagsPage };
