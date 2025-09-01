import * as React from 'react';
import {
  EuiComboBox,
  EuiComboBoxOptionOption,
  EuiFlexGroup,
  EuiFlexItem,
  EuiBadge,
} from '@elastic/eui';
import { useTheme } from '@app/shared/contexts/ThemeContext';

interface LabelFilterProps {
  availableLabels: string[];
  selectedLabels: string[];
  onLabelsChange: (labels: string[]) => void;
  placeholder?: string;
}

const LabelFilter: React.FC<LabelFilterProps> = ({
  availableLabels,
  selectedLabels,
  onLabelsChange,
  placeholder = 'Filter by key or value...',
}) => {
  const { theme } = useTheme();

  const handleLabelChange = (selectedOptions: EuiComboBoxOptionOption[]) => {
    const labels = selectedOptions.map(option => option.label);
    onLabelsChange(labels);
  };

  // Convert available labels to ComboBox options
  const labelOptions: EuiComboBoxOptionOption[] = availableLabels.map(label => ({
    label,
    value: label,
  }));

  // Convert selected labels to ComboBox options
  const selectedOptions: EuiComboBoxOptionOption[] = selectedLabels.map(label => ({
    label,
    value: label,
  }));

  return (
    <EuiComboBox
      placeholder={placeholder}
      options={labelOptions}
      selectedOptions={selectedOptions}
      onChange={handleLabelChange}
      isClearable={true}
      compressed
    />
  );
};

export { LabelFilter };