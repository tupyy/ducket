import * as React from 'react';
import {
  MenuToggle,
  MenuToggleElement,
  Dropdown,
  DropdownList,
  DropdownItem,
  Label,
  TextInputGroup,
  TextInputGroupMain,
  TextInputGroupUtilities,
  Flex,
  FlexItem,
  Button,
} from '@patternfly/react-core';
import { TimesIcon } from '@patternfly/react-icons';
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
  const [isOpen, setIsOpen] = React.useState<boolean>(false);
  const [inputValue, setInputValue] = React.useState<string>('');

  const handleInputChange = (_event: React.FormEvent<HTMLInputElement>, value: string) => {
    setInputValue(value);
  };

  const handleLabelSelect = (label: string) => {
    if (!selectedLabels.includes(label)) {
      onLabelsChange([...selectedLabels, label]);
    }
    setInputValue('');
    setIsOpen(false);
  };

  const handleLabelRemove = (labelToRemove: string) => {
    onLabelsChange(selectedLabels.filter((label) => label !== labelToRemove));
  };

  const handleClearAll = () => {
    onLabelsChange([]);
    setInputValue('');
  };

  // Filter available labels based on input and exclude already selected ones
  // Filter by either key or value
  const filteredLabels = availableLabels.filter((label) => {
    if (selectedLabels.includes(label)) {
      return false;
    }
    
    if (!inputValue) {
      return true;
    }
    
    const searchTerm = inputValue.toLowerCase();
    
    // Check if the entire label contains the search term (backward compatibility)
    if (label.toLowerCase().includes(searchTerm)) {
      return true;
    }
    
    // Split label by "=" and check key and value separately
    const [key, value] = label.split('=');
    if (key && key.toLowerCase().includes(searchTerm)) {
      return true;
    }
    if (value && value.toLowerCase().includes(searchTerm)) {
      return true;
    }
    
    return false;
  });

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle ref={toggleRef} onClick={() => setIsOpen(!isOpen)} isExpanded={isOpen} style={{ width: '300px' }}>
      <TextInputGroup>
        <TextInputGroupMain value={inputValue} placeholder={placeholder} onChange={handleInputChange} />
        {(inputValue || selectedLabels.length > 0) && (
          <TextInputGroupUtilities>
            <Button
              variant="plain"
              onClick={() => {
                setInputValue('');
                if (selectedLabels.length > 0) {
                  handleClearAll();
                }
              }}
              aria-label="Clear all"
            >
              <TimesIcon />
            </Button>
          </TextInputGroupUtilities>
        )}
      </TextInputGroup>
    </MenuToggle>
  );

  return (
    <Flex direction={{ default: 'column' }}>
      <FlexItem>
        <Dropdown
          isOpen={isOpen}
          onOpenChange={(isOpen: boolean) => setIsOpen(isOpen)}
          toggle={toggle}
          ouiaId="LabelFilterDropdown"
          shouldFocusToggleOnSelect
        >
          <DropdownList>
            {filteredLabels.length > 0 ? (
              filteredLabels.map((label, index) => (
                <DropdownItem key={index} value={label} onClick={() => handleLabelSelect(label)}>
                  {label}
                </DropdownItem>
              ))
            ) : (
              <DropdownItem isDisabled>{inputValue ? 'No matching labels found' : 'No more labels available'}</DropdownItem>
            )}
          </DropdownList>
        </Dropdown>
      </FlexItem>
      {selectedLabels.length > 0 && (
        <FlexItem>
          <Flex spaceItems={{ default: 'spaceItemsXs' }} style={{ marginTop: '8px' }}>
            {selectedLabels.map((label, index) => (
              <FlexItem key={index}>
                <Label
                  variant={theme === 'dark' ? 'outline' : 'filled'}
                  color="blue"
                  onClose={() => handleLabelRemove(label)}
                  closeBtnAriaLabel={`Remove ${label} filter`}
                  style={theme === 'dark' ? { color: '#73bcf7' } : {}}
                >
                  {label}
                </Label>
              </FlexItem>
            ))}
          </Flex>
        </FlexItem>
      )}
    </Flex>
  );
};

export { LabelFilter }; 