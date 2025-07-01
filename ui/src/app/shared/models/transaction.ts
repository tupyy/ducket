export interface ITagTransaction {
  href: string;
  value: string;
  rule: string;
}

export interface ITransaction {
  href: string;
  kind: string;
  date: Date;
  amount: number;
  account: string;
  description: string;
  tags: ITagTransaction[];
}

export interface ITransactions {
  total: number;
  items: ITransaction[];
}

export interface ITransactionForm {
  kind: string;
  date: string;
  content: string;
  amount: number;
  tags: Map<string, string>;
}

export interface ITransactionUpdateForm {
  name: string;
  kind: string;
  date: string;
  content: string;
  amount: number;
  tags: Map<string, string>;
}
