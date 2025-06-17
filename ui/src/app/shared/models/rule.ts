import { ITag } from '@app/shared/models/tag';

export interface IRule {
  href: string;
  name: string;
  pattern: string;
  created_at: Date;
  transactions: number;
  tags: ITag[];
}

export interface IRules {
  rules: Array<IRule>;
  total: number;
}

export interface IUpdateRuleForm {
  pattern: string;
  tags: Array<string>;
}

export interface IRuleForm extends IUpdateRuleForm {
  name: string;
}
