import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { PageSection } from '@patternfly/react-core';
import * as React from 'react';
import { LabelsList } from './Labels';
import { createLabel, deleteLabel, getLabels } from '@app/shared/reducers/label.reducer';
import { LabelForm } from './LabelForm';
import { CreateLabelError, Result } from '@app/shared/models/result';
import { ILabel } from '@app/shared/models/label';

const LabelsPage: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const labels = useAppSelector((state) => state.labels);
  const [isCreateFormActive, setIsCreateFormActive] = React.useState<boolean>(false);
  const [creatingLabel, setIsLabelCreating] = React.useState<boolean>(false);
  const [createLabelResult, setCreateLabelResult] = React.useState<Result<string, CreateLabelError>>(Result.succeed('label'));

  const getLabelsFromProps = () => {
    dispatch(getLabels());
  };

  const showCreateForm = () => {
    setIsCreateFormActive(true);
  };

  const closeCreateForm = () => {
    setIsCreateFormActive(false);
  };

  const _deleteLabel = (id: string) => {
    dispatch(deleteLabel(id));
  };

  const handleSubmitForm = (key: string, value: string) => {
    setIsLabelCreating(true);
    dispatch(createLabel({ key: key, value: value }));
    setIsLabelCreating(false);
    if (!labels.creating && !labels.createSuccess) {
      setCreateLabelResult(Result.fail(labels.errorMessage));
    }
    closeCreateForm();
  };

  const renderLabelList = (loading: boolean, labels: ReadonlyArray<ILabel>) => {
    const sortLabels = (labels: Array<ILabel> | []) => {
      return labels.sort((label1, label2) => {
        if (label1.created_at > label2.created_at) {
          return -1;
        }
        if (label1.created_at < label2.created_at) {
          return 1;
        }
        return 0;
      });
    };

    if (loading) {
      return <div>loading</div>;
    }
    return <LabelsList labels={sortLabels(labels.slice())} showCreateLabelFormCB={showCreateForm} deleteLabelCB={_deleteLabel} />;
  };

  React.useEffect(() => {
    getLabelsFromProps();
  }, []);

  return (
    <PageSection>
      {isCreateFormActive ? (
        <LabelForm
          closeFormCB={closeCreateForm}
          submitValue={handleSubmitForm}
          creating={creatingLabel}
          createResult={createLabelResult}
        />
      ) : (
        renderLabelList(labels.loading, labels.labels)
      )}
    </PageSection>
  );
};

export { LabelsPage };
