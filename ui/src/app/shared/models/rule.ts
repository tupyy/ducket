export interface IRule {
  id: number;
  name: string;
  filter: string;
  tags: string[];
  created_at?: string;
}
