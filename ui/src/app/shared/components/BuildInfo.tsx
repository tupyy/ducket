import * as React from 'react';
import { Tooltip } from '@patternfly/react-core';
import { CodeIcon } from '@patternfly/react-icons';

interface IBuildInfo {
  variant?: 'plain' | 'link';
}

const BuildInfo: React.FunctionComponent<IBuildInfo> = ({ variant = 'plain' }) => {
  const gitCommit = process.env.GIT_COMMIT || 'unknown';
  
  return (
    <Tooltip content={`Build commit: ${gitCommit}`}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '4px', opacity: 0.7 }}>
        <CodeIcon />
        <span style={{ fontSize: '0.75rem' }}>
          {gitCommit}
        </span>
      </div>
    </Tooltip>
  );
};

export { BuildInfo };