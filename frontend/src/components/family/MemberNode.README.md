# MemberNode Component

Custom React Flow node component for displaying family tree members in VamsaSetu.

## Features

- **Gender-based border colors**: Blue (#3b82f6) for male, pink (#ec4899) for female
- **Avatar display**: Shows member avatar or fallback to name initial
- **Relation badge**: Displays relationship label (e.g., "Father", "Mother")
- **Event indicator**: Glowing amber indicator for members with upcoming events
- **Hover effects**: Smooth scale and shadow animation on hover
- **Click handler**: Emits click event with member ID for opening details panel

## Usage

### Basic Setup

```tsx
import ReactFlow, { Node } from 'reactflow';
import MemberNode, { MemberNodeData } from './components/family/MemberNode';
import 'reactflow/dist/style.css';

// Register custom node type
const nodeTypes = {
  memberNode: MemberNode,
};

// Create nodes
const nodes: Node<MemberNodeData>[] = [
  {
    id: '1',
    type: 'memberNode',
    position: { x: 250, y: 0 },
    data: {
      id: '1',
      name: 'Rajesh Kumar',
      gender: 'male',
      avatarUrl: 'https://example.com/avatar.jpg',
      relationBadge: 'Father',
      hasUpcomingEvent: true,
      onNodeClick: (id) => console.log('Clicked:', id),
    },
  },
];

// Render
<ReactFlow nodes={nodes} nodeTypes={nodeTypes} />
```

### Props (MemberNodeData)

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `id` | `string` | Yes | Unique member identifier |
| `name` | `string` | Yes | Member's full name |
| `gender` | `'male' \| 'female' \| 'other'` | Yes | Member's gender (determines border color) |
| `avatarUrl` | `string` | No | URL to member's avatar image |
| `relationBadge` | `string` | No | Relationship label to display |
| `hasUpcomingEvent` | `boolean` | No | Shows glowing amber indicator if true |
| `onNodeClick` | `(id: string) => void` | No | Callback when node is clicked |

## Styling

The component uses Tailwind CSS with VamsaSetu theme colors:
- `saffron`: #E8650A
- `turmeric`: #F5A623
- `teal`: #0D4A52
- `charcoal`: #2C2420
- `amber`: #F59E0B

## Requirements Validation

**Validates: Requirements 3.2, 3.4, 3.6, 10.6**

- ✅ 3.2: Interactive family tree visualization with custom nodes
- ✅ 3.4: Member profile display with avatar and details
- ✅ 3.6: Visual indicators for upcoming events
- ✅ 10.6: Responsive UI with hover effects

## Example

See `MemberNode.example.tsx` for a complete working example with React Flow integration.

## Testing

To run tests (requires vitest setup):
```bash
npm test MemberNode.test.tsx
```

## Integration

This component is designed to be used with:
- React Flow for tree layout
- Framer Motion for animations
- Tailwind CSS for styling
- VamsaSetu theme colors and fonts
