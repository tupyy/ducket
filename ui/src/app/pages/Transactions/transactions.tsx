import * as React from 'react';
import { getTransactions } from '@app/shared/reducers/transaction.reducer';
import { useAppDispatch, useAppSelector } from '@app/shared/store';
import { Dropdown, PageSection, MenuToggle, MenuToggleElement, Flex, FlexItem } from '@patternfly/react-core';
import { TransactionList } from './list';
import { CalendarAltIcon, CalendarIcon, CubesIcon, TimesIcon } from '@patternfly/react-icons';
import { Content, EmptyState, EmptyStateBody, EmptyStateVariant } from '@patternfly/react-core';
import { TimePicker } from '@app/shared/components/time-picker';

const Transactions: React.FunctionComponent = () => {
  const dispatch = useAppDispatch();
  const transactions = useAppSelector((state) => state.transactions);
  const [isOpen, setIsOpen] = React.useState(false);

  const onToggleClick = () => {
    setIsOpen(!isOpen);
  };

  React.useEffect(() => {
    dispatch(getTransactions());
  }, []);

  const emptyState = (
    <EmptyState variant={EmptyStateVariant.full} titleText="No transactions" icon={CubesIcon}>
      <EmptyStateBody>
        <Content>
          <Content component="p">Please add some transactions</Content>
        </Content>
      </EmptyStateBody>
    </EmptyState>
  );

  return (
    <PageSection hasBodyWrapper={false}>
      {transactions.transactions.length == 0 ? (
        emptyState
      ) : (
        <>
          <Flex
            direction={{ default: 'row' }}
            spaceItems={{ default: 'spaceItemsSm' }}
            alignItems={{ default: 'alignItemsFlexEnd' }}
            justifyContent={{default: 'justifyContentFlexEnd'}}

          >
            <FlexItem>
              <Dropdown
                isOpen={isOpen}
                onOpenChange={(isOpen: boolean) => setIsOpen(isOpen)}
                toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                  <MenuToggle ref={toggleRef} onClick={onToggleClick} isExpanded={isOpen}>
                    <Flex direction={{ default: 'row' }} spaceItems={{ default: 'spaceItemsSm' }}>
                      <FlexItem>
                        <CalendarAltIcon />
                      </FlexItem>
                      <FlexItem>Dropdown</FlexItem>
                    </Flex>
                  </MenuToggle>
                )}
                ouiaId="BasicDropdown"
                shouldFocusToggleOnSelect
              >
                <TimePicker />
              </Dropdown>
            </FlexItem>
          </Flex>
          <TransactionList transactions={transactions.transactions} />
        </>
      )}
    </PageSection>
  );
};

export { Transactions };
