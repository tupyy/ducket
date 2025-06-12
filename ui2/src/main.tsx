import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import getStore from './src/config/store';
import { Provider } from 'react-redux';

const store = getStore();

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </React.StrictMode>
);
