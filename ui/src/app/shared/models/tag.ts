import { IRule } from '@app/shared/models/rule';

export interface ITag {
  href: string;
  value: string;
  created_at: Date;
  transactions: number;
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

export interface ITagReport {
  tag: string;
  amount: number;
}

export interface ITransactionTypeReport {
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
