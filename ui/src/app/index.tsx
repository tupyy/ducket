import * as React from 'react';
import '@patternfly/patternfly/patternfly.css';
import { BrowserRouter as Router } from 'react-router-dom';
import { Provider } from 'react-redux';
import store from '@app/shared/store';
import { AppLayout } from '@app/AppLayout/AppLayout';
import { AppRoutes } from '@app/routes';

const App: React.FunctionComponent = () => (
  <Provider store={store}>
    <Router>
      <AppLayout>
        <AppRoutes />
      </AppLayout>
    </Router>
  </Provider>
);

export default App;
