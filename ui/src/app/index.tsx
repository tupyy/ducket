import * as React from 'react';
import '@elastic/eui/dist/eui_theme_light.css';
import '@elastic/eui/dist/eui_theme_dark.css';
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