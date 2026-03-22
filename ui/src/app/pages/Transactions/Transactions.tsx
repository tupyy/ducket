import * as React from 'react';
import {
  PageSection,
  Title,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  Pagination,
  EmptyState,
  EmptyStateBody,
  Spinner,
  Bullseye,
  Label,
  LabelGroup,
  Content,
  Dropdown,
  DropdownItem,
  DropdownList,
  MenuToggle,
  type MenuToggleElement,
  Modal,
  ModalVariant,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  Form,
  FormGroup,
  TextArea,
} from '@patternfly/react-core';
import { EllipsisVIcon } from '@patternfly/react-icons';
import { Table, Thead, Tr, Th, Tbody, Td, ThProps } from '@patternfly/react-table';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import {
  getTransactions, setPage, setPerPage, setFilter, setSort,
  addTag, removeTag,
  addAccount, removeAccount, toggleKind,
  buildCompositeFilter, updateTransactionInfo,
} from '@app/shared/reducers/transaction.reducer';
import { ITransaction } from '@app/shared/models/transaction';
import { TransactionFilter } from './TransactionFilter';

const Transactions: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const [openMenuId, setOpenMenuId] = React.useState<number | null>(null);
  const [editTransaction, setEditTransaction] = React.useState<ITransaction | null>(null);
  const [editInfo, setEditInfo] = React.useState('');
  const { transactions, loading, total, page, perPage, filter, selectedTags, selectedAccounts, selectedKind, sort, errorMessage } = useAppSelector(
    (state) => state.transactions,
  );

  React.useEffect(() => {
    const compositeFilter = buildCompositeFilter(filter, selectedAccounts, selectedKind);
    dispatch(
      getTransactions({
        filter: compositeFilter || undefined,
        tags: selectedTags.length > 0 ? selectedTags : undefined,
        sort: sort.length > 0 ? sort : undefined,
        limit: perPage,
        offset: (page - 1) * perPage,
      }),
    );
  }, [dispatch, page, perPage, filter, selectedAccounts, selectedKind, selectedTags, sort]);

  const handleFilterChange = (newFilter: string) => {
    dispatch(setFilter(newFilter));
  };

  const onSetPage = (_evt: React.MouseEvent | React.KeyboardEvent | MouseEvent, newPage: number) => {
    dispatch(setPage(newPage));
  };

  const onPerPageSelect = (_evt: React.MouseEvent | React.KeyboardEvent | MouseEvent, newPerPage: number) => {
    dispatch(setPerPage(newPerPage));
  };

  const getActiveSortDirection = (field: string): 'asc' | 'desc' | undefined => {
    const active = sort.find((s) => s.field === field);
    return active?.direction;
  };

  const handleSort = (field: string) => {
    const current = getActiveSortDirection(field);
    let newDirection: 'asc' | 'desc';
    if (!current) {
      newDirection = field === 'amount' ? 'desc' : 'asc';
    } else {
      newDirection = current === 'asc' ? 'desc' : 'asc';
    }
    dispatch(setSort([{ field, direction: newDirection }]));
  };

  const getSortParams = (field: string): ThProps['sort'] => ({
    sortBy: {
      index: sort.length > 0 && sort[0].field === field ? 0 : undefined,
      direction: getActiveSortDirection(field),
      defaultDirection: field === 'amount' ? 'desc' : 'asc',
    },
    onSort: () => handleSort(field),
    columnIndex: 0,
  });

  const formatDate = (dateStr: string) => {
    try {
      const d = new Date(dateStr);
      const day = String(d.getDate()).padStart(2, '0');
      const month = String(d.getMonth() + 1).padStart(2, '0');
      const year = d.getFullYear();
      return `${day}-${month}-${year}`;
    } catch {
      return dateStr;
    }
  };

  if (errorMessage && transactions.length === 0) {
    return (
      <PageSection>
        <EmptyState titleText="Error loading transactions" headingLevel="h2">
          <EmptyStateBody>{errorMessage}</EmptyStateBody>
        </EmptyState>
      </PageSection>
    );
  }

  return (
    <PageSection>
      <Title headingLevel="h1" size="lg" style={{ marginBottom: '1rem' }}>
        Transactions
      </Title>

      <TransactionFilter
        filter={filter}
        onFilterChange={handleFilterChange}
        total={total}
        page={page}
        perPage={perPage}
        onSetPage={onSetPage}
        onPerPageSelect={onPerPageSelect}
      />

      {(selectedTags.length > 0 || selectedAccounts.length > 0 || selectedKind) && (
        <div style={{ marginBottom: '0.5rem', display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
          {selectedKind && (
            <LabelGroup categoryName="Type">
              <Label
                color={selectedKind === 'debit' ? 'red' : 'green'}
                onClose={() => dispatch(toggleKind(selectedKind))}
              >
                {selectedKind}
              </Label>
            </LabelGroup>
          )}
          {selectedAccounts.length > 0 && (
            <LabelGroup categoryName="Account">
              {selectedAccounts.map((acct) => (
                <Label key={acct} variant="outline" color="blue" onClose={() => dispatch(removeAccount(acct))}>
                  {acct}
                </Label>
              ))}
            </LabelGroup>
          )}
          {selectedTags.length > 0 && (
            <LabelGroup categoryName="Tags">
              {selectedTags.map((tag) => (
                <Label key={tag} variant="outline" color="teal" onClose={() => dispatch(removeTag(tag))}>
                  {tag}
                </Label>
              ))}
            </LabelGroup>
          )}
        </div>
      )}

      {loading ? (
        <Bullseye style={{ minHeight: '200px' }}>
          <Spinner size="xl" />
        </Bullseye>
      ) : transactions.length === 0 ? (
        <EmptyState titleText="No transactions" headingLevel="h2">
          <EmptyStateBody>
            {filter ? 'No transactions match the current filter.' : 'Import some transactions to get started.'}
          </EmptyStateBody>
        </EmptyState>
      ) : (
        <>
          <Table aria-label="Transactions table" variant="compact">
            <Thead>
              <Tr>
                <Th width={10} sort={getSortParams('date')}>Date</Th>
                <Th width={10} sort={getSortParams('account')}>Account</Th>
                <Th width={10} sort={getSortParams('kind')}>Type</Th>
                <Th>Content</Th>
                <Th>Info</Th>
                <Th width={15}>Tags</Th>
                <Th width={10} sort={getSortParams('amount')}>Amount</Th>
                <Th width={10}></Th>
              </Tr>
            </Thead>
            <Tbody>
              {transactions.map((txn) => (
                <Tr key={txn.id}>
                  <Td dataLabel="Date">{formatDate(txn.date)}</Td>
                  <Td dataLabel="Account">
                    <Label
                      color="blue"
                      onClick={() => dispatch(addAccount(txn.account))}
                      style={{ cursor: 'pointer' }}
                    >
                      {txn.account}
                    </Label>
                  </Td>
                  <Td dataLabel="Type">
                    <Label
                  color={txn.kind === 'debit' ? 'red' : 'green'}
                  onClick={() => dispatch(toggleKind(txn.kind))}
                  style={{ cursor: 'pointer' }}
                >{txn.kind}</Label>
                  </Td>
                  <Td dataLabel="Content">
                    <Content
                      component="p"
                      style={{ maxWidth: '400px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}
                    >
                      {txn.content}
                    </Content>
                  </Td>
                  <Td dataLabel="Info">
                    {txn.info && (
                      <Content
                        component="p"
                        style={{ maxWidth: '300px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}
                      >
                        {txn.info}
                      </Content>
                    )}
                  </Td>
                  <Td dataLabel="Tags">
                    {txn.tags.map((tag, i) => (
                      <Label
                        key={i}
                        color="teal"
                        onClick={() => dispatch(addTag(tag))}
                        style={{ marginRight: 4, marginBottom: 2, cursor: 'pointer' }}
                      >
                        {tag}
                      </Label>
                    ))}
                  </Td>
                  <Td dataLabel="Amount" style={{ fontFamily: 'monospace' }}>
                    €{txn.amount.toFixed(2)}
                  </Td>
                  <Td isActionCell>
                    <Dropdown
                      isOpen={openMenuId === txn.id}
                      onOpenChange={(isOpen) => setOpenMenuId(isOpen ? txn.id : null)}
                      toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                        <MenuToggle
                          ref={toggleRef}
                          variant="plain"
                          onClick={() => setOpenMenuId(openMenuId === txn.id ? null : txn.id)}
                          isExpanded={openMenuId === txn.id}
                        >
                          <EllipsisVIcon />
                        </MenuToggle>
                      )}
                      popperProps={{ position: 'right' }}
                    >
                      <DropdownList>
                        <DropdownItem
                          key="edit-info"
                          onClick={() => {
                            setEditTransaction(txn);
                            setEditInfo(txn.info || '');
                            setOpenMenuId(null);
                          }}
                        >
                          Edit info
                        </DropdownItem>
                      </DropdownList>
                    </Dropdown>
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>

          <Toolbar>
            <ToolbarContent>
              <ToolbarItem variant="pagination" align={{ default: 'alignEnd' }}>
                <Pagination
                  itemCount={total}
                  perPage={perPage}
                  page={page}
                  onSetPage={onSetPage}
                  onPerPageSelect={onPerPageSelect}
                  variant="bottom"
                />
              </ToolbarItem>
            </ToolbarContent>
          </Toolbar>
        </>
      )}

      <Modal
        variant={ModalVariant.small}
        isOpen={editTransaction !== null}
        onClose={() => setEditTransaction(null)}
      >
        <ModalHeader title="Edit info" />
        <ModalBody>
          <Form>
            <FormGroup label="Info" fieldId="edit-info">
              <TextArea
                id="edit-info"
                value={editInfo}
                onChange={(_event, value) => setEditInfo(value)}
                rows={4}
              />
            </FormGroup>
          </Form>
        </ModalBody>
        <ModalFooter>
          <Button
            variant="primary"
            onClick={() => {
              if (editTransaction) {
                dispatch(updateTransactionInfo({ id: editTransaction.id, info: editInfo }));
                setEditTransaction(null);
              }
            }}
          >
            Save
          </Button>
          <Button variant="link" onClick={() => setEditTransaction(null)}>
            Cancel
          </Button>
        </ModalFooter>
      </Modal>
    </PageSection>
  );
};

export { Transactions };
