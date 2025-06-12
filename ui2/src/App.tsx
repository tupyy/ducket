import { BrowserRouter as Router } from 'react-router';
import { AppRoutes } from './shared/routes/routes';
import { AppLayout } from './modules/home/home';

const App: React.FC = () => (
  <Router>
    <AppLayout>
      <AppRoutes />
    </AppLayout>
  </Router>
);

export default App;
