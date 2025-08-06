import * as React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import {
  Button,
  Masthead,
  MastheadContent,
  MastheadMain,
  MastheadToggle,
  Nav,
  NavExpandable,
  NavItem,
  NavList,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
  Page,
  PageSidebar,
  PageSidebarBody,
  SkipToContent,
} from '@patternfly/react-core';
import { IAppRoute, IAppRouteGroup, routes } from '@app/routes';
import { BarsIcon } from '@patternfly/react-icons';
import { ThemeToggle } from '@app/shared/components/ThemeToggle';
import { BuildInfo } from '@app/shared/components/BuildInfo';
import { useTheme } from '@app/shared/contexts/ThemeContext';

// Create a context for page styling
const PageStyleContext = React.createContext<{
  pageBackgroundColor: string;
}>({
  pageBackgroundColor: 'var(--pf-v6-global--BackgroundColor--100)',
});

export const usePageStyle = () => React.useContext(PageStyleContext);

interface IAppLayout {
  children: React.ReactNode;
}

const AppLayout: React.FunctionComponent<IAppLayout> = ({ children }) => {
  const [sidebarOpen, setSidebarOpen] = React.useState(true);
  const { theme } = useTheme();

  React.useEffect(() => {
    // Apply theme to document html element
    if (theme === 'dark') {
      document.documentElement.classList.add('pf-v6-theme-dark');
      document.documentElement.classList.remove('pf-v6-theme-light');
    } else {
      document.documentElement.classList.add('pf-v6-theme-light');
      document.documentElement.classList.remove('pf-v6-theme-dark');
    }
  }, [theme]);

  // Define the page background color to pass to children
  const pageBackgroundColor = 'var(--pf-v6-global--BackgroundColor--200)';

  const masthead = (
    <Masthead>
      <MastheadMain>
        <MastheadToggle>
          <Button
            icon={<BarsIcon />}
            variant="plain"
            onClick={() => setSidebarOpen(!sidebarOpen)}
            aria-label="Global navigation"
          />
        </MastheadToggle>
      </MastheadMain>
      <MastheadContent>
        <Toolbar id="masthead-toolbar">
          <ToolbarContent>
            <ToolbarGroup align={{ default: 'alignEnd' }}>
              <ToolbarItem>
                <BuildInfo variant="plain" />
              </ToolbarItem>
              <ToolbarItem>
                <ThemeToggle variant="plain" />
              </ToolbarItem>
            </ToolbarGroup>
          </ToolbarContent>
        </Toolbar>
      </MastheadContent>
    </Masthead>
  );

  const location = useLocation();

  const renderNavItem = (route: IAppRoute, index: number) => (
    <NavItem key={`${route.label}-${index}`} id={`${route.label}-${index}`} isActive={route.path === location.pathname}>
      <NavLink to={route.path}>{route.label}</NavLink>
    </NavItem>
  );

  const renderNavGroup = (group: IAppRouteGroup, groupIndex: number) => (
    <NavExpandable
      isExpanded={true}
      key={`${group.label}-${groupIndex}`}
      id={`${group.label}-${groupIndex}`}
      title={group.label}
      isActive={group.routes.some((route) => route.path === location.pathname)}
    >
      {group.routes.map((route, idx) => route.label && renderNavItem(route, idx))}
    </NavExpandable>
  );

  const Navigation = (
    <Nav id="nav-primary-simple">
      <NavList id="nav-list-simple">
        {routes.map(
          (route, idx) => route.label && (!route.routes ? renderNavItem(route, idx) : renderNavGroup(route, idx))
        )}
      </NavList>
    </Nav>
  );

  const Sidebar = (
    <PageSidebar>
      <PageSidebarBody>{Navigation}</PageSidebarBody>
    </PageSidebar>
  );

  const pageId = 'primary-app-container';

  const PageSkipToContent = (
    <SkipToContent
      onClick={(event) => {
        event.preventDefault();
        const primaryContentContainer = document.getElementById(pageId);
        primaryContentContainer?.focus();
      }}
      href={`#${pageId}`}
    >
      Skip to Content
    </SkipToContent>
  );
  return (
    <PageStyleContext.Provider value={{ pageBackgroundColor }}>
      <Page
        mainContainerId={pageId}
        masthead={masthead}
        sidebar={sidebarOpen && Sidebar}
        skipToContent={PageSkipToContent}
      >
        {children}
      </Page>
    </PageStyleContext.Provider>
  );
};

export { AppLayout };
