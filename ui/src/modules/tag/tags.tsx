import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '../../config/store';
import { getTags } from './tag.reducer';
import { ITag } from '../../shared/models/tag';

export const TagsPage = () => {
  const dispatch = useAppDispatch();

  const getTagsFromProps = () => {
      dispatch(getTags());
  };

  useEffect(() => {
      getTagsFromProps();
  },[]);

  const tags = useAppSelector(state => state.tags.tags);
  return (
      <div>
        {tags.map((tag: ITag, i: number) => (
            <div key={`tag-${i}`}>
                {tag.value}
            </div>
        ))}
      </div>
  )
};

export default TagsPage;
