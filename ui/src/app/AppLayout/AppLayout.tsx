import * as React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import {
  Brand,
  Masthead,
  MastheadBrand,
  MastheadContent,
  MastheadMain,
  MastheadToggle,
  Nav,
  NavItem,
  NavList,
  Page,
  PageSidebar,
  PageSidebarBody,
  PageToggleButton,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  Toolbar,
  Label,
} from '@patternfly/react-core';
import BarsIcon from '@patternfly/react-icons/dist/esm/icons/bars-icon';

interface AppLayoutProps {
  children: React.ReactNode;
}

const AppLayout: React.FunctionComponent<AppLayoutProps> = ({ children }) => {
  const [isSidebarOpen, setIsSidebarOpen] = React.useState(true);
  const location = useLocation();

  const gitCommit = process.env.GIT_COMMIT || 'dev';

  const routes = [
    { path: '/', label: 'Dashboard' },
    { path: '/transactions', label: 'Transactions' },
    { path: '/rules', label: 'Rules' },
    { path: '/import', label: 'Import' },
  ];

  const headerToolbar = (
    <Toolbar id="header-toolbar" isFullHeight isStatic>
      <ToolbarContent>
        <ToolbarGroup align={{ default: 'alignEnd' }}>
          <ToolbarItem>
            <Label isCompact>{gitCommit}</Label>
          </ToolbarItem>
        </ToolbarGroup>
      </ToolbarContent>
    </Toolbar>
  );

  const masthead = (
    <Masthead>
      <MastheadToggle>
        <PageToggleButton
          variant="plain"
          aria-label="Global navigation"
          isSidebarOpen={isSidebarOpen}
          onSidebarToggle={() => setIsSidebarOpen(!isSidebarOpen)}
        >
          <BarsIcon />
        </PageToggleButton>
      </MastheadToggle>
      <MastheadMain>
        <MastheadBrand data-codemods>Finante</MastheadBrand>
      </MastheadMain>
      <MastheadContent>{headerToolbar}</MastheadContent>
    </Masthead>
  );

  const navigation = (
    <Nav aria-label="Global">
      <NavList>
        {routes.map((route) => (
          <NavItem key={route.path} isActive={location.pathname === route.path}>
            <NavLink to={route.path}>{route.label}</NavLink>
          </NavItem>
        ))}
      </NavList>
    </Nav>
  );

  const sidebar = (
    <PageSidebar isSidebarOpen={isSidebarOpen}>
      <PageSidebarBody>{navigation}</PageSidebarBody>
    </PageSidebar>
  );

  return (
    <Page masthead={masthead} sidebar={sidebar}>
      {children}
    </Page>
  );
};

export { AppLayout };
