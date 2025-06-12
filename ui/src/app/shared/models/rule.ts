import { ITag } from "@app/shared/models/tag";

export interface IRule {
    href: string;
    name: string;
    pattern?: string;
    tags: ITag[];
}
