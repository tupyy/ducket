import * as React from 'react';
import { EuiToolTip, EuiIcon, EuiFlexGroup, EuiFlexItem, EuiText } from '@elastic/eui';

interface IBuildInfo {
  variant?: 'default' | 'subdued';
}

const BuildInfo: React.FunctionComponent<IBuildInfo> = ({ variant = 'subdued' }) => {
  const gitCommit = process.env.GIT_COMMIT || 'unknown';
  
  return (
    <EuiToolTip content={`Build commit: ${gitCommit}`}>
      <EuiFlexGroup gutterSize="xs" alignItems="center" responsive={false}>
        <EuiFlexItem grow={false}>
          <EuiIcon type="console" size="s" />
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          <EuiText size="xs" color={variant}>
            {gitCommit.substring(0, 8)}
          </EuiText>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiToolTip>
  );
};

export { BuildInfo };