# RelationshipEdge Component

Custom React Flow edge component for displaying color-coded relationship connections in the VamsaSetu family tree visualization.

## Features

- **Color-coded edges** based on relationship type:
  - `SPOUSE_OF`: Rose color (#E11D48)
  - `PARENT_OF`: Teal color (#14b8a6)
  - `SIBLING_OF`: Amber color (#F59E0B)
- **Bezier curve styling** for smooth, natural-looking connections
- **Consistent stroke width** (2px) for clear visibility
- **Fallback color** (gray) for undefined relationship types

## Usage

### Basic Setup

```tsx
import ReactFlow from 'reactflow';
import RelationshipEdge from './components/family/RelationshipEdge';
import type { RelationshipEdgeData } from './components/family/RelationshipEdge';

const edgeTypes = {
  relationshipEdge: RelationshipEdge,
};

function FamilyTree() {
  const edges = [
    {
      id: 'e1-2',
      source: '1',
      target: '2',
      type: 'relationshipEdge',
      data: { type: 'SPOUSE_OF' },
    },
    {
      id: 'e1-3',
      source: '1',
      target: '3',
      type: 'relationshipEdge',
      data: { type: 'PARENT_OF' },
    },
  ];

  return (
    <ReactFlow
      edges={edges}
      edgeTypes={edgeTypes}
    />
  );
}
```

### Edge Data Interface

```tsx
interface RelationshipEdgeData {
  type: 'SPOUSE_OF' | 'PARENT_OF' | 'SIBLING_OF';
}
```

## Props

The component accepts standard React Flow `EdgeProps` with custom `RelationshipEdgeData`:

| Prop | Type | Description |
|------|------|-------------|
| `id` | `string` | Unique identifier for the edge |
| `sourceX` | `number` | X coordinate of the source node |
| `sourceY` | `number` | Y coordinate of the source node |
| `targetX` | `number` | X coordinate of the target node |
| `targetY` | `number` | Y coordinate of the target node |
| `sourcePosition` | `Position` | Position of the source handle |
| `targetPosition` | `Position` | Position of the target handle |
| `data` | `RelationshipEdgeData` | Custom data containing relationship type |
| `markerEnd` | `string` | Optional marker for edge end |

## Color Reference

| Relationship Type | Color Name | Hex Code | Usage |
|------------------|------------|----------|-------|
| `SPOUSE_OF` | Rose | `#E11D48` | Marital relationships |
| `PARENT_OF` | Teal | `#14b8a6` | Parent-child relationships |
| `SIBLING_OF` | Amber | `#F59E0B` | Sibling relationships |
| Default | Gray | `#9CA3AF` | Fallback for undefined types |

## Integration with React Flow

This component is designed to work seamlessly with React Flow's edge system:

1. Register the component in `edgeTypes`
2. Set edge `type` to `'relationshipEdge'`
3. Provide relationship type in edge `data`
4. React Flow handles positioning and rendering

## Example

See `RelationshipEdge.example.tsx` for a complete working example with multiple relationship types.

## Requirements

**Validates: Requirements 3.3** - Interactive family tree visualization with color-coded relationship edges

## Dependencies

- `reactflow` - Core React Flow library
- `react` - React framework

## Notes

- The component uses React Flow's `BaseEdge` for consistent rendering
- Bezier paths are calculated automatically based on node positions
- Edge styling is applied via inline styles for reliable color rendering
- The component is fully typed with TypeScript for type safety
