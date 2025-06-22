import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { Dashboard } from '@app/pages/Dashboard/Dashboard';
import { NotFound } from '@app/pages/NotFound/NotFound';
import { Rules } from '@app/pages/Rules/Rules';
import { TagsPage } from './pages/Tags/Page';
import { Transactions } from './pages/Transactions/transactions';

export interface IAppRoute {
  label?: string; // Excluding the label will exclude the route from the nav sidebar in AppLayout
  /* eslint-disable @typescript-eslint/no-explicit-any */
  element: React.ReactElement;
  /* eslint-enable @typescript-eslint/no-explicit-any */
  exact?: boolean;
  path: string;
  title: string;
  routes?: undefined;
}

export interface IAppRouteGroup {
  label: string;
  routes: IAppRoute[];
}

export type AppRouteConfig = IAppRoute | IAppRouteGroup;

const routes: AppRouteConfig[] = [
  {
    element: <Dashboard />,
    exact: true,
    label: 'Dashboard',
    path: '/',
    title: 'PatternFly Seed | Main Dashboard',
  },
  {
    element: <Transactions />,
    exact: true,
    label: 'Transactions',
    path: '/transactions',
    title: 'Finance | Transactions Page',
  },
  {
    element: <TagsPage />,
    exact: true,
    label: 'Tags',
    path: '/tags',
    title: 'Finance | Tags Page',
  },
  {
    element: <Rules />,
    exact: true,
    label: 'Rules',
    path: '/rules',
    title: 'Finance | Rules Page',
  },
];

const flattenedRoutes: IAppRoute[] = routes.reduce(
  (flattened, route) => [...flattened, ...(route.routes ? route.routes : [route])],
  [] as IAppRoute[]
);

const AppRoutes = (): React.ReactElement => (
  <Routes>
    {flattenedRoutes.map(({ path, element }, idx) => (
      <Route path={path} element={element} key={idx} />
    ))}
    <Route element={<NotFound />} />
  </Routes>
);

export { AppRoutes, routes };
