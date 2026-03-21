import * as React from 'react';
import { Route, Routes } from 'react-router-dom';
import { Dashboard } from '@app/pages/Dashboard/Dashboard';
import { Transactions } from '@app/pages/Transactions/Transactions';
import { Rules } from '@app/pages/Rules/Rules';
import { Import } from '@app/pages/Import/Import';

const AppRoutes: React.FunctionComponent = () => (
  <Routes>
    <Route path="/" element={<Dashboard />} />
    <Route path="/transactions" element={<Transactions />} />
    <Route path="/rules" element={<Rules />} />
    <Route path="/import" element={<Import />} />
  </Routes>
);

export { AppRoutes };
