# TreeControls Component

## Overview

The `TreeControls` component provides a vertical toolbar with action buttons for controlling the family tree visualization. It includes zoom controls, fit view functionality, and optional buttons for adding members and relationships.

**Validates: Requirements 3.5**

## Features

- **Zoom Controls**: Zoom in and zoom out with smooth animations
- **Fit View**: Automatically fit the entire tree in the viewport
- **Add Member**: Optional button to trigger add member modal
- **Add Relationship**: Optional button to trigger add relationship modal
- **Responsive Design**: Works on both desktop and mobile devices
- **Accessibility**: Full keyboard navigation and ARIA labels
- **VamsaSetu Theme**: Styled with brand colors (saffron, teal, ivory)

## Props

```typescript
interface TreeControlsProps {
  onAddMember?: () => void;
  onAddRelationship?: () => void;
  className?: string;
}
```

### Prop Details

- `onAddMember` (optional): Callback function when "Add Member" button is clicked. If not provided, the button won't be rendered.
- `onAddRelationship` (optional): Callback function when "Add Relationship" button is clicked. If not provided, the button won't be rendered.
- `className` (optional): Additional CSS classes to apply to the container.

## Usage

### Basic Usage

```tsx
import TreeControls from './components/family/TreeControls';
import { ReactFlowProvider } from 'reactflow';

function FamilyTreePage() {
  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen">
        {/* TreeCanvas component */}
        <TreeCanvas />
        
        {/* TreeControls positioned absolutely */}
        <div className="absolute top-4 right-4 z-10">
          <TreeControls />
        </div>
      </div>
    </ReactFlowProvider>
  );
}
```

### With Add Member and Relationship Callbacks

```tsx
import { useState } from 'react';
import TreeControls from './components/family/TreeControls';
import AddMemberModal from './components/family/AddMemberModal';
import AddRelationshipModal from './components/family/AddRelationshipModal';

function FamilyTreePage() {
  const [showAddMember, setShowAddMember] = useState(false);
  const [showAddRelationship, setShowAddRelationship] = useState(false);

  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen">
        <TreeCanvas />
        
        <div className="absolute top-4 right-4 z-10">
          <TreeControls
            onAddMember={() => setShowAddMember(true)}
            onAddRelationship={() => setShowAddRelationship(true)}
          />
        </div>

        {showAddMember && (
          <AddMemberModal onClose={() => setShowAddMember(false)} />
        )}
        
        {showAddRelationship && (
          <AddRelationshipModal onClose={() => setShowAddRelationship(false)} />
        )}
      </div>
    </ReactFlowProvider>
  );
}
```

### Mobile Positioning

```tsx
// Position at bottom-right on mobile, top-right on desktop
<div className="absolute bottom-4 right-4 md:top-4 md:bottom-auto z-10">
  <TreeControls
    onAddMember={() => setShowAddMember(true)}
    onAddRelationship={() => setShowAddRelationship(true)}
  />
</div>
```

## Styling

The component uses VamsaSetu theme colors:
- **Primary buttons** (Add Member): Saffron (#E8650A)
- **Secondary buttons** (Add Relationship): Teal (#0D4A52)
- **Outline buttons** (Zoom controls): White background with saffron border
- **Background**: White with shadow and border

## Accessibility

- All buttons have `aria-label` attributes for screen readers
- Toolbar has `role="toolbar"` and `aria-label="Family tree controls"`
- All buttons have descriptive `title` attributes for tooltips
- Keyboard accessible with proper focus indicators
- Minimum touch target size of 40x40px (meets WCAG guidelines)

## Icons

The component uses emoji icons for simplicity:
- **Zoom In**: `+`
- **Zoom Out**: `−`
- **Fit View**: `⊡`
- **Add Member**: `👤`
- **Add Relationship**: `🔗`

You can replace these with icon library components (e.g., Lucide React, Heroicons) by modifying the button content.

## Dependencies

- `reactflow`: For `useReactFlow` hook (zoom and fit view functionality)
- `../ui/Button`: VamsaSetu Button component

## Integration with TreeCanvas

The TreeControls component must be used within a `ReactFlowProvider` context, as it uses the `useReactFlow` hook to access zoom and fit view functions.

```tsx
import { ReactFlowProvider } from 'reactflow';

// ✅ Correct
<ReactFlowProvider>
  <TreeCanvas />
  <TreeControls />
</ReactFlowProvider>

// ❌ Incorrect - will throw error
<TreeControls />
```

## Responsive Behavior

The component is designed to be positioned absolutely within the tree container. Recommended positioning:

- **Desktop**: Top-right corner (`top-4 right-4`)
- **Mobile**: Bottom-right corner (`bottom-4 right-4`)
- **Z-index**: Use `z-10` or higher to ensure it appears above the tree

## Testing

See `TreeControls.test.tsx` for comprehensive unit tests covering:
- Zoom control functionality
- Fit view functionality
- Conditional rendering of add buttons
- Callback invocation
- Accessibility attributes

## Related Components

- `TreeCanvas`: The main family tree visualization component
- `Button`: VamsaSetu UI button component
- `AddMemberModal`: Modal for adding new family members
- `AddRelationshipModal`: Modal for adding relationships

## Future Enhancements

Potential improvements for future iterations:
- Add undo/redo buttons
- Add export/download tree as image button
- Add fullscreen toggle
- Add layout algorithm selector (hierarchical, force-directed, etc.)
- Add filter/search toggle button
- Add settings/preferences button
