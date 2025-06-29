import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import { BrowserRouter as Router } from 'react-router-dom';
import { AppLayout } from '@app/AppLayout/AppLayout';
import { AppRoutes } from '@app/routes';
import { ThemeProvider } from '@app/shared/contexts/ThemeContext';
import '@app/app.css';

const App: React.FunctionComponent = () => (
  <ThemeProvider defaultTheme="light">
    <Router>
      <AppLayout>
        <AppRoutes />
      </AppLayout>
    </Router>
  </ThemeProvider>
);

export default App;
