# TreeCanvas Integration Guide

This guide explains how to integrate the TreeCanvas component into your application.

## Quick Start

### 1. Basic Integration

```tsx
import TreeCanvas from './components/family/TreeCanvas';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

function FamilyTreePage() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="w-full h-screen">
        <TreeCanvas />
      </div>
    </QueryClientProvider>
  );
}
```

### 2. With React Router

```tsx
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import TreeCanvas from './components/family/TreeCanvas';

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/family-tree" element={
            <div className="w-full h-screen">
              <TreeCanvas />
            </div>
          } />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
```

## Prerequisites

### Required Dependencies
All dependencies are already installed in the project:
- `react`: ^19.2.4
- `reactflow`: ^11.11.4
- `@tanstack/react-query`: ^5.91.3
- `framer-motion`: ^12.38.0

### Required Components
The following components must exist (already created):
- `MemberNode` - Custom node component
- `RelationshipEdge` - Custom edge component
- `useFamilyTree` - Data fetching hook

### Required Services
- `familyTreeService` - API service for fetching tree data
- API endpoint: `GET /api/family/tree`

## Container Requirements

The TreeCanvas component requires a container with defined dimensions:

### ✅ Good Examples

```tsx
// Full screen
<div className="w-full h-screen">
  <TreeCanvas />
</div>

// Fixed height
<div className="w-full h-[600px]">
  <TreeCanvas />
</div>

// Flex container
<div className="flex-1">
  <TreeCanvas />
</div>
```

### ❌ Bad Examples

```tsx
// No height defined - tree won't render
<div className="w-full">
  <TreeCanvas />
</div>

// Inline height without width
<div style={{ height: '600px' }}>
  <TreeCanvas />
</div>
```

## API Integration

### Expected API Response

```json
{
  "data": {
    "nodes": [
      {
        "id": "1",
        "type": "memberNode",
        "position": { "x": 0, "y": 0 },
        "data": {
          "id": "1",
          "name": "John Doe",
          "gender": "male",
          "avatarUrl": "https://example.com/avatar.jpg",
          "relationBadge": "Father",
          "hasUpcomingEvent": false
        }
      }
    ],
    "edges": [
      {
        "id": "e1-2",
        "source": "1",
        "target": "2",
        "type": "relationshipEdge",
        "animated": false,
        "style": {
          "stroke": "#14b8a6",
          "strokeWidth": "2"
        },
        "data": {
          "type": "PARENT_OF"
        }
      }
    ]
  }
}
```

### API Configuration

Ensure your API base URL is configured in `services/api.ts`:

```typescript
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
```

## Styling Integration

### Required CSS

The component requires React Flow CSS to be imported. This is already done in the component:

```tsx
import 'reactflow/dist/style.css';
```

### Tailwind Configuration

Ensure your `tailwind.config.js` includes the custom colors:

```javascript
module.exports = {
  theme: {
    extend: {
      colors: {
        cream: '#FFF8E7',
        saffron: '#D4AF37',
        teal: '#14b8a6',
        charcoal: '#2C2C2C',
        turmeric: '#F59E0B',
      },
    },
  },
};
```

## Mobile Integration

### Viewport Meta Tag

Ensure your `index.html` has the correct viewport meta tag:

```html
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=5.0, user-scalable=yes">
```

### Touch Event Support

The component automatically handles touch events. No additional configuration needed.

## Advanced Integration

### 1. With Navigation

```tsx
import { useNavigate } from 'react-router-dom';

// Modify TreeCanvas to accept onNodeClick prop
function FamilyTreePage() {
  const navigate = useNavigate();

  const handleNodeClick = (nodeId: string) => {
    navigate(`/member/${nodeId}`);
  };

  return (
    <div className="w-full h-screen">
      <TreeCanvas />
    </div>
  );
}
```

### 2. With Side Panel

```tsx
function FamilyTreeWithDetails() {
  const [selectedMember, setSelectedMember] = useState<string | null>(null);

  return (
    <div className="flex h-screen">
      <div className="flex-1">
        <TreeCanvas />
      </div>
      {selectedMember && (
        <aside className="w-80 bg-white border-l p-6">
          <MemberDetails memberId={selectedMember} />
        </aside>
      )}
    </div>
  );
}
```

### 3. With Header and Footer

```tsx
function FamilyTreeLayout() {
  return (
    <div className="flex flex-col h-screen">
      <header className="bg-saffron text-white p-4">
        <h1 className="text-2xl font-bold">Family Tree</h1>
      </header>
      
      <main className="flex-1">
        <TreeCanvas />
      </main>
      
      <footer className="bg-gray-100 p-4 text-center">
        <p className="text-sm text-gray-600">© 2024 VamsaSetu</p>
      </footer>
    </div>
  );
}
```

## State Management

### React Query Configuration

The component uses React Query for data fetching. Configure the query client:

```tsx
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      cacheTime: 10 * 60 * 1000, // 10 minutes
      retry: 3,
      refetchOnWindowFocus: false,
    },
  },
});
```

### Manual Refetch

To manually refetch the tree data:

```tsx
import { useQueryClient } from '@tanstack/react-query';

function RefreshButton() {
  const queryClient = useQueryClient();

  const handleRefresh = () => {
    queryClient.invalidateQueries({ queryKey: ['familyTree'] });
  };

  return (
    <button onClick={handleRefresh}>
      Refresh Tree
    </button>
  );
}
```

## Performance Optimization

### 1. Code Splitting

```tsx
import { lazy, Suspense } from 'react';

const TreeCanvas = lazy(() => import('./components/family/TreeCanvas'));

function FamilyTreePage() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <div className="w-full h-screen">
        <TreeCanvas />
      </div>
    </Suspense>
  );
}
```

### 2. Memoization

The component already uses memoization internally. No additional optimization needed.

## Troubleshooting

### Tree Not Displaying

**Problem**: Component renders but tree is not visible

**Solutions**:
1. Check container has defined height
2. Verify API is returning data
3. Check browser console for errors
4. Ensure React Flow CSS is imported

### Gestures Not Working

**Problem**: Pinch-to-zoom or pan not working on mobile

**Solutions**:
1. Check viewport meta tag
2. Verify touch events are not blocked by parent elements
3. Test on actual device (not just browser emulator)

### Performance Issues

**Problem**: Tree is slow with many nodes

**Solutions**:
1. Implement pagination or lazy loading
2. Reduce node complexity
3. Use React Flow's built-in optimization features
4. Consider server-side layout calculation

### Styling Issues

**Problem**: Colors or layout not matching design

**Solutions**:
1. Verify Tailwind config includes custom colors
2. Check React Flow CSS is imported
3. Ensure no conflicting global styles
4. Use browser DevTools to inspect elements

## Testing

### Unit Tests

```tsx
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import TreeCanvas from './TreeCanvas';

test('renders loading state', () => {
  const queryClient = new QueryClient();
  
  render(
    <QueryClientProvider client={queryClient}>
      <TreeCanvas />
    </QueryClientProvider>
  );
  
  expect(screen.getByText(/loading/i)).toBeInTheDocument();
});
```

### Integration Tests

```tsx
import { renderWithProviders } from './test-utils';
import TreeCanvas from './TreeCanvas';

test('displays family tree after loading', async () => {
  const { findByTestId } = renderWithProviders(<TreeCanvas />);
  
  const reactFlow = await findByTestId('react-flow');
  expect(reactFlow).toBeInTheDocument();
});
```

## Best Practices

1. **Always wrap in QueryClientProvider**: TreeCanvas requires React Query context
2. **Define container dimensions**: Use explicit height for proper rendering
3. **Handle loading states**: Show appropriate feedback during data fetch
4. **Test on mobile devices**: Verify gestures work on actual devices
5. **Monitor performance**: Watch for performance issues with large trees
6. **Cache appropriately**: Configure React Query cache based on your needs
7. **Handle errors gracefully**: Provide clear error messages to users

## Next Steps

After integrating TreeCanvas:

1. **Add navigation**: Implement node click handlers to navigate to member details
2. **Add search**: Implement search functionality to find specific members
3. **Add filters**: Allow filtering by relationship type, generation, etc.
4. **Add export**: Implement export to image or PDF functionality
5. **Add editing**: Allow users to add/edit members directly from the tree
6. **Add animations**: Enhance with smooth transitions on data updates

## Support

For issues or questions:
- Check the README: `TreeCanvas.README.md`
- Review examples: `TreeCanvas.example.tsx`
- Check component tests: `TreeCanvas.test.tsx`
- Review related components: `MemberNode.tsx`, `RelationshipEdge.tsx`
