import { ILabel } from '@app/shared/models/label';

export interface IRule {
  href: string;
  name: string;
  pattern: string;
  created_at: Date;
  transactions: number;
  labels: ILabel[];
}

export interface IRules {
  rules: Array<IRule>;
  total: number;
}

export interface IUpdateRuleForm {
  pattern: string;
  labels: { [key: string]: string };
}

export interface IRuleForm extends IUpdateRuleForm {
  name: string;
}
