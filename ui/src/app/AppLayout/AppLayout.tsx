import * as React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import {
  EuiPage,
  EuiPageBody,
  EuiHeader,
  EuiHeaderSection,
  EuiHeaderSectionItem,
  EuiHeaderLinks,
  EuiHeaderLink,
  EuiFlexGroup,
  EuiFlexItem,
  EuiButtonIcon,
} from '@elastic/eui';
import { ThemeToggle } from '@app/shared/components/ThemeToggle';
import { BuildInfo } from '@app/shared/components/BuildInfo';
import { useTheme } from '@app/shared/contexts/ThemeContext';

interface IAppLayout {
  children: React.ReactNode;
}

const AppLayout: React.FunctionComponent<IAppLayout> = ({ children }) => {
  const { theme } = useTheme();
  const location = useLocation();

  const navigation = [
    { path: '/transactions', label: 'Transactions' },
    { path: '/transactions/upload', label: 'Upload' },
    { path: '/rules', label: 'Rules' },
    { path: '/labels', label: 'Labels' },
  ];

  const header = (
    <EuiHeader position="fixed">
      <EuiHeaderSection grow={false}>
        <EuiHeaderSectionItem>
          <EuiHeaderLinks>
            {navigation.map((item) => (
              <EuiHeaderLink
                key={item.path}
                href={item.path}
                isActive={location.pathname.startsWith(item.path)}
              >
                <NavLink 
                  to={item.path} 
                  style={{ textDecoration: 'none', color: 'inherit' }}
                >
                  {item.label}
                </NavLink>
              </EuiHeaderLink>
            ))}
          </EuiHeaderLinks>
        </EuiHeaderSectionItem>
      </EuiHeaderSection>

      <EuiHeaderSection side="right">
        <EuiHeaderSectionItem>
          <EuiFlexGroup gutterSize="s" alignItems="center">
            <EuiFlexItem grow={false}>
              <BuildInfo />
            </EuiFlexItem>
            <EuiFlexItem grow={false}>
              <ThemeToggle iconOnly={true} />
            </EuiFlexItem>
          </EuiFlexGroup>
        </EuiHeaderSectionItem>
      </EuiHeaderSection>
    </EuiHeader>
  );

  return (
    <EuiPage paddingSize="none">
      {header}
      <EuiPageBody 
        panelled={false} 
        paddingSize="none" 
        style={{ paddingTop: '48px' }} // Add space for fixed header
      >
        <main>
          {children}
        </main>
      </EuiPageBody>
    </EuiPage>
  );
};

export { AppLayout };