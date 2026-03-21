import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import TreeCanvas from './TreeCanvas';
import * as useFamilyTreeHook from '../../hooks/useFamilyTree';

// Mock ReactFlow
vi.mock('reactflow', () => ({
  default: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="react-flow">{children}</div>
  ),
  Background: () => <div data-testid="background" />,
  Controls: () => <div data-testid="controls" />,
  MiniMap: () => <div data-testid="minimap" />,
  useNodesState: () => [[], vi.fn(), vi.fn()],
  useEdgesState: () => [[], vi.fn(), vi.fn()],
}));

// Mock custom components
vi.mock('./MemberNode', () => ({
  default: () => <div data-testid="member-node" />,
}));

vi.mock('./RelationshipEdge', () => ({
  default: () => <div data-testid="relationship-edge" />,
}));

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

describe('TreeCanvas', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state', () => {
    vi.spyOn(useFamilyTreeHook, 'useFamilyTree').mockReturnValue({
      data: undefined,
      isLoading: true,
      error: null,
    } as any);

    render(<TreeCanvas />, { wrapper: createWrapper() });

    expect(screen.getByText(/loading family tree/i)).toBeInTheDocument();
  });

  it('renders error state', () => {
    const errorMessage = 'Failed to fetch tree';
    vi.spyOn(useFamilyTreeHook, 'useFamilyTree').mockReturnValue({
      data: undefined,
      isLoading: false,
      error: new Error(errorMessage),
    } as any);

    render(<TreeCanvas />, { wrapper: createWrapper() });

    expect(screen.getByText(/failed to load family tree/i)).toBeInTheDocument();
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });

  it('renders ReactFlow with data', async () => {
    const mockData = {
      data: {
        nodes: [
          {
            id: '1',
            type: 'memberNode',
            position: { x: 0, y: 0 },
            data: {
              id: '1',
              name: 'John Doe',
              gender: 'male',
              avatarUrl: '',
              relationBadge: 'Father',
              hasUpcomingEvent: false,
            },
          },
        ],
        edges: [
          {
            id: 'e1-2',
            source: '1',
            target: '2',
            type: 'relationshipEdge',
            animated: false,
            style: { stroke: '#14b8a6', strokeWidth: '2' },
          },
        ],
      },
    };

    vi.spyOn(useFamilyTreeHook, 'useFamilyTree').mockReturnValue({
      data: mockData,
      isLoading: false,
      error: null,
    } as any);

    render(<TreeCanvas />, { wrapper: createWrapper() });

    await waitFor(() => {
      expect(screen.getByTestId('react-flow')).toBeInTheDocument();
    });

    expect(screen.getByTestId('background')).toBeInTheDocument();
    expect(screen.getByTestId('controls')).toBeInTheDocument();
    expect(screen.getByTestId('minimap')).toBeInTheDocument();
  });

  it('renders empty tree when no data', () => {
    vi.spyOn(useFamilyTreeHook, 'useFamilyTree').mockReturnValue({
      data: { data: { nodes: [], edges: [] } },
      isLoading: false,
      error: null,
    } as any);

    render(<TreeCanvas />, { wrapper: createWrapper() });

    expect(screen.getByTestId('react-flow')).toBeInTheDocument();
  });
});
