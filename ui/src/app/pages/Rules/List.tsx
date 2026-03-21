import * as React from 'react';
import {
  Title,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  Button,
  Label,
  EmptyState,
  EmptyStateBody,
  EmptyStateActions,
  EmptyStateFooter,
  Content,
} from '@patternfly/react-core';
import { PencilAltIcon, TrashIcon } from '@patternfly/react-icons';
import { Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import { IRule } from '@app/shared/models/rule';

interface RulesListProps {
  rules: IRule[];
  loading: boolean;
  onCreateRule: () => void;
  onEditRule: (rule: IRule) => void;
  onDeleteRule: (id: number) => void;
}

const RulesList: React.FunctionComponent<RulesListProps> = ({
  rules,
  loading,
  onCreateRule,
  onEditRule,
  onDeleteRule,
}) => {
  if (rules.length === 0 && !loading) {
    return (
      <EmptyState titleText="No rules" headingLevel="h2">
        <EmptyStateBody>Create your first rule to automatically tag transactions.</EmptyStateBody>
        <EmptyStateFooter>
          <EmptyStateActions>
            <Button variant="primary" onClick={onCreateRule}>
              Create Rule
            </Button>
          </EmptyStateActions>
        </EmptyStateFooter>
      </EmptyState>
    );
  }

  return (
    <>
      <Title headingLevel="h1" size="lg" style={{ marginBottom: '1rem' }}>
        Rules
      </Title>

      <Toolbar>
        <ToolbarContent>
          <ToolbarItem>
            <Button variant="primary" onClick={onCreateRule}>
              Create Rule
            </Button>
          </ToolbarItem>
        </ToolbarContent>
      </Toolbar>

      <Table aria-label="Rules table" variant="compact">
        <Thead>
          <Tr>
            <Th width={20}>Name</Th>
            <Th width={30}>Filter</Th>
            <Th>Tags</Th>
            <Th width={15}>Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {rules.map((rule) => (
            <Tr key={rule.id}>
              <Td dataLabel="Name">
                <Content component="p" style={{ fontWeight: 'bold' }}>{rule.name}</Content>
              </Td>
              <Td dataLabel="Filter">
                <code>{rule.filter}</code>
              </Td>
              <Td dataLabel="Tags">
                {rule.tags.map((tag, i) => (
                  <Label key={i} color="teal" isCompact style={{ marginRight: 4, marginBottom: 2 }}>
                    {tag}
                  </Label>
                ))}
              </Td>
              <Td dataLabel="Actions">
                <Button variant="plain" aria-label="Edit" onClick={() => onEditRule(rule)}>
                  <PencilAltIcon />
                </Button>
                <Button variant="plain" aria-label="Delete" onClick={() => onDeleteRule(rule.id)}>
                  <TrashIcon />
                </Button>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </>
  );
};

export { RulesList };
