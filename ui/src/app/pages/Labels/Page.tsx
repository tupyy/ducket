import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { PageSection } from '@patternfly/react-core';
import * as React from 'react';
import { LabelsList } from './Labels';
import { ILabel } from '@app/shared/models/label';
import { getLabels } from '@app/shared/reducers/label.reducer';

const LabelsPage: React.FunctionComponent = () => {
  const labels = useAppSelector((state) => state.labels);
  const dispatch = useAppDispatch();

  const getLabelsFromProps = () => {
    dispatch(getLabels());
  };

  React.useEffect(() => {
    getLabelsFromProps();
  }, []);

  const renderLabelList = (loading: boolean, labels: ReadonlyArray<ILabel>) => {
    const sortLabels = (labels: Array<ILabel> | []) => {
      return labels.sort((l1, l2) => {
        if (l1.key > l2.key) {
          return 1;
        }
        if (l1.key < l2.key) {
          return -1;
        }
        return 0;
      });
    };

    if (loading) {
      return <div>loading</div>;
    }
    return <LabelsList labels={sortLabels(labels.slice())} />;
  };

  return <PageSection>{renderLabelList(labels.loading, labels.labels)}</PageSection>;
};

export { LabelsPage };
