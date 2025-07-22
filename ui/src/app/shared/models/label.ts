import { IRule } from '@app/shared/models/rule';

export interface ILabel {
  href: string;
  key: string;
  value: string;
  created_at: Date;
  transactions: number;
  rules: IRule[];
}

export interface ILabels {
  total: number;
  labels: ILabel[];
}

export interface ILabelForm {
  key: string;
  value: string;
}

export interface ILabelUpdateForm {
  id: string;
  key: string;
  value: string;
}

// Legacy aliases for backward compatibility during migration
export interface ITagForm {
  key: string;
  value: string;
}

export interface ITagUpdateForm {
  id: string;
  key: string;
  value: string;
}

export interface ITag {
  href: string;
  key: string;
  value: string;
  created_at: Date;
  transactions: number;
  rules: IRule[];
}

export interface ITags {
  total: number;
  tags: ITag[];
}

export interface ITagReport {
  tag: string;
  amount: number;
}

export interface ITransactionTypeReport {
  type: 'debit' | 'credit';
  amount: number;
}

export interface IAccountTransactionTypeReport {
  account: number;
  type: 'debit' | 'credit';
  amount: number;
}

export interface IMonthlyTagReport {
  tag: string;
  month: number; // 1-12
  year: number;
  amount: number;
  transactionCount: number;
}

export interface IMonthlyTagSummary {
  monthYear: string; // Format: "YYYY-MM"
  tags: IMonthlyTagReport[];
  totalAmount: number;
}

export function isKeyOfObject<T extends object>(key: string , obj: T):boolean  {
  return key in obj;
}
