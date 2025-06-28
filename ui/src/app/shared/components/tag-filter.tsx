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
  Button
} from '@patternfly/react-core';
import { TimesIcon } from '@patternfly/react-icons';

interface TagFilterProps {
  availableTags: string[];
  selectedTags: string[];
  onTagsChange: (tags: string[]) => void;
  placeholder?: string;
}

const TagFilter: React.FC<TagFilterProps> = ({ 
  availableTags, 
  selectedTags, 
  onTagsChange,
  placeholder = "Filter by tags..."
}) => {
  const [isOpen, setIsOpen] = React.useState<boolean>(false);
  const [inputValue, setInputValue] = React.useState<string>('');

  const handleInputChange = (_event: React.FormEvent<HTMLInputElement>, value: string) => {
    setInputValue(value);
  };

  const handleTagSelect = (tag: string) => {
    if (!selectedTags.includes(tag)) {
      onTagsChange([...selectedTags, tag]);
    }
    setInputValue('');
    setIsOpen(false);
  };

  const handleTagRemove = (tagToRemove: string) => {
    onTagsChange(selectedTags.filter(tag => tag !== tagToRemove));
  };

  const handleClearAll = () => {
    onTagsChange([]);
    setInputValue('');
  };

  // Filter available tags based on input and exclude already selected ones
  const filteredTags = availableTags.filter(tag => 
    !selectedTags.includes(tag) && 
    tag.toLowerCase().includes(inputValue.toLowerCase())
  );

  const toggle = (toggleRef: React.Ref<MenuToggleElement>) => (
    <MenuToggle 
      ref={toggleRef} 
      onClick={() => setIsOpen(!isOpen)}
      isExpanded={isOpen}
      style={{ width: '300px' }}
    >
      <TextInputGroup>
        <TextInputGroupMain
          value={inputValue}
          placeholder={placeholder}
          onChange={handleInputChange}
        />
        {(inputValue || selectedTags.length > 0) && (
          <TextInputGroupUtilities>
            <Button
              variant="plain"
              onClick={() => {
                setInputValue('');
                if (selectedTags.length > 0) {
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
          ouiaId="TagFilterDropdown"
          shouldFocusToggleOnSelect
        >
          <DropdownList>
            {filteredTags.length > 0 ? (
              filteredTags.map((tag, index) => (
                <DropdownItem 
                  key={index}
                  value={tag}
                  onClick={() => handleTagSelect(tag)}
                >
                  {tag}
                </DropdownItem>
              ))
            ) : (
              <DropdownItem isDisabled>
                {inputValue ? 'No matching tags found' : 'No more tags available'}
              </DropdownItem>
            )}
          </DropdownList>
        </Dropdown>
      </FlexItem>
      {selectedTags.length > 0 && (
        <FlexItem>
          <Flex spaceItems={{ default: 'spaceItemsXs' }} style={{ marginTop: '8px' }}>
            {selectedTags.map((tag, index) => (
              <FlexItem key={index}>
                <Label 
                  variant="filled" 
                  color="blue"
                  onClose={() => handleTagRemove(tag)}
                  closeBtnAriaLabel={`Remove ${tag} filter`}
                >
                  {tag}
                </Label>
              </FlexItem>
            ))}
          </Flex>
        </FlexItem>
      )}
    </Flex>
  );
};

export { TagFilter }; 