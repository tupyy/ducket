import { Route, BrowserRouter as Router, Routes } from 'react-router';
import Home from './modules/home/home';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </Router>
  );
}

export default App;
