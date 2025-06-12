import { ITag } from "./tag";

export interface IRule {
    href: string;
    name: string;
    pattern?: string;
    tags: ITag[];
}
