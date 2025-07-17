export interface ILabelTransaction {
  href: string;
  key: string;
  value: string;
  ruleHref: string;
}

export interface ITransaction {
  href: string;
  kind: string;
  date: string;
  amount: number;
  account: number;
  description: string;
  labels: ILabelTransaction[];
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
  labels: Map<string, string>;
}

export interface ITransactionUpdateForm {
  name: string;
  kind: string;
  date: string;
  content: string;
  amount: number;
  labels: Map<string, string>;
}
