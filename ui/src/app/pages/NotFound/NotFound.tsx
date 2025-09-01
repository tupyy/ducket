import * as React from 'react';
import { EuiPageSection, EuiEmptyPrompt, EuiButton } from '@elastic/eui';
import { useNavigate } from 'react-router-dom';

const NotFound: React.FunctionComponent = () => {
  const navigate = useNavigate();
  
  const handleGoHome = () => {
    navigate('/');
  };

  return (
    <EuiPageSection>
      <EuiEmptyPrompt
        icon="alert"
        title={<h2>404 Page not found</h2>}
        body="We didn't find a page that matches the address you navigated to."
        actions={
          <EuiButton fill onClick={handleGoHome}>
            Take me home
          </EuiButton>
        }
      />
    </EuiPageSection>
  );
};

export { NotFound };
