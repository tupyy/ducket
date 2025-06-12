import { IRule } from "./rule";

export interface ITag {
    href: string;
    value: string;
    rules: IRule[];
}

export interface ITags {
    total: number;
    tags: ITag[];
}



