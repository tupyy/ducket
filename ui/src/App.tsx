import React from 'react';
import { EuiProvider, EuiText } from '@elastic/eui';

function App() {
  return (
    <EuiProvider>
      <EuiText>
        <p>Hello</p>
      </EuiText>
    </EuiProvider>
  );
}

export default App;
