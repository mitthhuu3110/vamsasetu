# TreeCanvas Component

A React Flow-based family tree visualization component with interactive features, mobile gesture support, and custom node/edge rendering.

## Features

- ✅ **React Flow Integration**: Uses React Flow for graph visualization
- ✅ **Custom Node Types**: Renders family members using MemberNode component
- ✅ **Custom Edge Types**: Displays relationships using RelationshipEdge component
- ✅ **Interactive Controls**: Zoom, pan, and fit-view controls
- ✅ **Mobile Gestures**: Pinch-to-zoom and pan gestures on mobile devices
- ✅ **Mini Map**: Navigation overview for large family trees
- ✅ **Background Grid**: Visual grid pattern for better orientation
- ✅ **Loading States**: Displays loading spinner while fetching data
- ✅ **Error Handling**: Shows error messages when data fetch fails
- ✅ **Auto-fit View**: Automatically fits the tree to viewport on load

## Usage

```tsx
import TreeCanvas from './components/family/TreeCanvas';
import { QueryClientProvider } from '@tanstack/react-query';

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="w-full h-screen">
        <TreeCanvas />
      </div>
    </QueryClientProvider>
  );
}
```

## Props

TreeCanvas is a self-contained component with no props. It:
- Fetches family tree data using `useFamilyTree` hook
- Manages its own state for nodes and edges
- Handles loading and error states internally

## Data Structure

The component expects data from the API in the following format:

```typescript
interface FamilyTree {
  nodes: Array<{
    id: string;
    type: 'memberNode';
    position: { x: number; y: number };
    data: {
      id: string;
      name: string;
      gender: 'male' | 'female' | 'other';
      avatarUrl?: string;
      relationBadge?: string;
      hasUpcomingEvent?: boolean;
    };
  }>;
  edges: Array<{
    id: string;
    source: string;
    target: string;
    type: 'relationshipEdge';
    animated: boolean;
    style: {
      stroke: string;
      strokeWidth: string;
    };
  }>;
}
```

## Features in Detail

### 1. Custom Node Types

The component registers `memberNode` as a custom node type, which renders using the `MemberNode` component. Each node displays:
- Member avatar or initial
- Member name
- Gender-based border color
- Relation badge
- Upcoming event indicator

### 2. Custom Edge Types

The component registers `relationshipEdge` as a custom edge type, which renders using the `RelationshipEdge` component. Edges are color-coded by relationship type:
- **Spouse**: Rose color (#E11D48)
- **Parent**: Teal color (#14b8a6)
- **Sibling**: Amber color (#F59E0B)

### 3. Interactive Controls

React Flow Controls provide:
- **Zoom In/Out**: Buttons to zoom in and out
- **Fit View**: Button to fit entire tree in viewport
- **Lock/Unlock**: Toggle interaction lock

### 4. Mobile Gesture Support

The component enables:
- `panOnScroll`: Pan the canvas by scrolling
- `panOnDrag`: Pan by dragging
- `zoomOnPinch`: Pinch-to-zoom on touch devices
- `zoomOnScroll`: Zoom with mouse wheel
- `zoomOnDoubleClick`: Double-click to zoom in

### 5. Mini Map

The mini map provides:
- Overview of entire tree structure
- Color-coded nodes by gender (blue for male, pink for female)
- Clickable navigation to different parts of the tree
- Zoomable and pannable

### 6. Background Pattern

A subtle grid pattern with saffron color (#D4AF37) provides visual orientation and matches the app's theme.

## States

### Loading State
Displays a spinner with "Loading family tree..." message while data is being fetched.

### Error State
Shows an error icon and message when data fetch fails, including the error details.

### Success State
Renders the interactive family tree with all features enabled.

## Styling

The component uses:
- **Tailwind CSS**: For utility classes
- **Custom Colors**: Matches app theme (cream background, saffron accents)
- **React Flow CSS**: Imported from 'reactflow/dist/style.css'

## Dependencies

- `react`: ^19.2.4
- `reactflow`: ^11.11.4
- `@tanstack/react-query`: ^5.91.3
- `framer-motion`: ^12.38.0 (used by MemberNode)

## Integration with Other Components

### MemberNode
Custom node component that renders individual family members. See `MemberNode.README.md` for details.

### RelationshipEdge
Custom edge component that renders relationships between members. See `RelationshipEdge.README.md` for details.

### useFamilyTree Hook
React Query hook that fetches family tree data from the API. Located in `hooks/useFamilyTree.ts`.

## API Integration

The component fetches data from:
```
GET /api/family/tree
```

Expected response:
```json
{
  "data": {
    "nodes": [...],
    "edges": [...]
  }
}
```

## Performance Considerations

1. **React Flow Optimization**: React Flow handles rendering optimization internally
2. **Memoization**: Node click handlers are memoized with `useCallback`
3. **Lazy Updates**: Nodes and edges update only when data changes
4. **Query Caching**: React Query caches tree data for 5 minutes

## Mobile Optimization

The component is fully mobile-optimized:
- Touch gestures for pan and zoom
- Responsive controls
- Mini map for navigation on small screens
- Optimized for portrait and landscape orientations

## Accessibility

- Keyboard navigation supported by React Flow
- Focus management for interactive elements
- ARIA labels on controls
- Color contrast meets WCAG standards

## Future Enhancements

Potential improvements:
- [ ] Node selection and highlighting
- [ ] Search and filter functionality
- [ ] Export tree as image
- [ ] Collaborative editing
- [ ] Undo/redo functionality
- [ ] Custom layouts (hierarchical, radial, etc.)

## Troubleshooting

### Tree not displaying
- Ensure QueryClientProvider wraps the component
- Check API endpoint is accessible
- Verify data format matches expected structure

### Gestures not working on mobile
- Ensure React Flow CSS is imported
- Check touch event handlers are not blocked
- Verify viewport meta tag in HTML

### Performance issues with large trees
- Consider implementing virtualization
- Reduce node complexity
- Optimize edge rendering
- Use React Flow's built-in performance features

## Related Components

- `MemberNode.tsx` - Custom node component
- `RelationshipEdge.tsx` - Custom edge component
- `useFamilyTree.ts` - Data fetching hook
- `familyTreeService.ts` - API service

## Requirements Validation

**Validates: Requirements 3.1, 3.5, 9.2**

- **3.1**: Interactive family tree visualization with zoom and pan
- **3.5**: Mobile-responsive with touch gestures
- **9.2**: Real-time updates via React Query integration
