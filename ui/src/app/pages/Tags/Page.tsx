import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { PageSection } from '@patternfly/react-core';
import * as React from 'react';
import { TagsList } from './Tags';
import { createTag, getTags, deleteTag } from '@app/shared/reducers/tag.reducer';
import { TagForm } from './TagForm';
import { CreateTagError, Result } from '@app/shared/models/result';
import { ITag } from '@app/shared/models/tag';

const TagsPage: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const tags = useAppSelector((state) => state.tags);
  const [isCreateFormActive, setIsCreateFormActive] = React.useState<boolean>(false);
  const [creatingTag, setIsTagCreating] = React.useState<boolean>(false);
  const [createTagResult, setCreateTagResult] = React.useState<Result<string, CreateTagError>>(Result.succeed('tag'));

  const getTagsFromProps = () => {
    dispatch(getTags());
  };

  const showCreateForm = () => {
    setIsCreateFormActive(true);
  };

  const closeCreateForm = () => {
    setIsCreateFormActive(false);
  };

  const _deleteTag = (name: string) => {
    dispatch(deleteTag(name));
  };

  const handleSubmitForm = (value: string) => {
    setIsTagCreating(true);
    dispatch(createTag({ value: value }));
    setIsTagCreating(false);
    if (!tags.creating && !tags.createSuccess) {
      setCreateTagResult(Result.fail(tags.errorMessage));
    }
    closeCreateForm();
  };

  const renderTagList = (loading: boolean, tags: ReadonlyArray<ITag>) => {
    const sortTags = (tags: Array<ITag> | []) => {
      return tags.sort((tag1, tag2) => {
        if (tag1.value < tag2.value) {
          return -1;
        }
        if (tag1.value > tag2.value) {
          return 1;
        }
        return 0;
      });
    };

    if (loading) {
      return <div>loading</div>;
    }
    return <TagsList tags={sortTags(tags.slice())} showCreateTagFormCB={showCreateForm} deleteTagCB={_deleteTag} />;
  };

  React.useEffect(() => {
    getTagsFromProps();
  }, []);

  return (
    <PageSection>
      {isCreateFormActive ? (
        <TagForm
          closeFormCB={closeCreateForm}
          submitValue={handleSubmitForm}
          creating={creatingTag}
          createResult={createTagResult}
        />
      ) : (
        renderTagList(tags.loading, tags.tags)
      )}
    </PageSection>
  );
};

export { TagsPage };
