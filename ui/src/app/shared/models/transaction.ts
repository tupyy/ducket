export interface ITransaction {
  id: number;
  hash: string;
  date: string;
  account: number;
  kind: 'debit' | 'credit';
  amount: number;
  content: string;
  info?: string;
  recipient?: string;
  tags: string[];
}
