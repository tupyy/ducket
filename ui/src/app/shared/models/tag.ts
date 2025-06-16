import { IRule } from '@app/shared/models/rule';

export interface ITag {
  href: string;
  value: string;
  created_at: Date;
  rules: IRule[];
}

export interface ITags {
  total: number;
  tags: ITag[];
}

export interface ITagForm {
  value: string;
}

export interface ITagUpdateForm {
  id: string;
  value: string;
}
