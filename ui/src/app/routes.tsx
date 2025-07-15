import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { Dashboard } from '@app/pages/Dashboard/Dashboard';
import { NotFound } from '@app/pages/NotFound/NotFound';
import { Rules } from '@app/pages/Rules/Rules';
import { LabelsPage } from './pages/Labels/Page';
import { Transactions } from './pages/Transactions/transactions';
import { FileUpload } from './pages/FileUpload/FileUpload';

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
    label: 'Transactions',
    routes: [
      {
        element: <Transactions />,
        exact: true,
        label: 'View Transactions',
        path: '/transactions',
        title: 'Finance | Transactions Page',
      },
      {
        element: <FileUpload />,
        exact: true,
        label: 'Upload Files',
        path: '/transactions/upload',
        title: 'Finance | File Upload',
      },
    ],
  },
  {
    element: <LabelsPage />,
    exact: true,
    label: 'Labels',
    path: '/tags',
    title: 'Finance | Labels Page',
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
  [] as IAppRoute[],
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
