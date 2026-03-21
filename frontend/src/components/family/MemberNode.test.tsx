import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import MemberNode from './MemberNode';
import type { MemberNodeData } from './MemberNode';
import type { NodeProps } from 'reactflow';

describe('MemberNode', () => {
  const mockNodeProps: NodeProps<MemberNodeData> = {
    id: '1',
    data: {
      id: '1',
      name: 'Test User',
      gender: 'male',
      avatarUrl: 'https://example.com/avatar.jpg',
      relationBadge: 'Father',
      hasUpcomingEvent: false,
    },
    type: 'memberNode',
    selected: false,
    isConnectable: true,
    xPos: 0,
    yPos: 0,
    dragging: false,
    zIndex: 0,
  };

  it('renders member name correctly', () => {
    render(<MemberNode {...mockNodeProps} />);
    expect(screen.getByText('Test User')).toBeInTheDocument();
  });

  it('renders relation badge when provided', () => {
    render(<MemberNode {...mockNodeProps} />);
    expect(screen.getByText('Father')).toBeInTheDocument();
  });

  it('applies blue border for male gender', () => {
    const { container } = render(<MemberNode {...mockNodeProps} />);
    const node = container.querySelector('.border-blue-500');
    expect(node).toBeInTheDocument();
  });

  it('applies pink border for female gender', () => {
    const femaleProps = {
      ...mockNodeProps,
      data: { ...mockNodeProps.data, gender: 'female' as const },
    };
    const { container } = render(<MemberNode {...femaleProps} />);
    const node = container.querySelector('.border-pink-500');
    expect(node).toBeInTheDocument();
  });

  it('shows upcoming event indicator when hasUpcomingEvent is true', () => {
    const propsWithEvent = {
      ...mockNodeProps,
      data: { ...mockNodeProps.data, hasUpcomingEvent: true },
    };
    const { container } = render(<MemberNode {...propsWithEvent} />);
    const indicator = container.querySelector('.bg-amber-500');
    expect(indicator).toBeInTheDocument();
  });

  it('does not show event indicator when hasUpcomingEvent is false', () => {
    const { container } = render(<MemberNode {...mockNodeProps} />);
    const indicator = container.querySelector('.bg-amber-500');
    expect(indicator).not.toBeInTheDocument();
  });

  it('calls onNodeClick when clicked', () => {
    const mockOnClick = vi.fn();
    const propsWithClick = {
      ...mockNodeProps,
      data: { ...mockNodeProps.data, onNodeClick: mockOnClick },
    };
    const { container } = render(<MemberNode {...propsWithClick} />);
    const node = container.querySelector('.cursor-pointer');
    
    if (node) {
      fireEvent.click(node);
      expect(mockOnClick).toHaveBeenCalledWith('1');
    }
  });

  it('displays avatar image when avatarUrl is provided', () => {
    render(<MemberNode {...mockNodeProps} />);
    const avatar = screen.getByAltText('Test User');
    expect(avatar).toBeInTheDocument();
    expect(avatar).toHaveAttribute('src', 'https://example.com/avatar.jpg');
  });

  it('displays initial letter when no avatarUrl is provided', () => {
    const propsWithoutAvatar = {
      ...mockNodeProps,
      data: { ...mockNodeProps.data, avatarUrl: undefined },
    };
    render(<MemberNode {...propsWithoutAvatar} />);
    expect(screen.getByText('T')).toBeInTheDocument(); // First letter of "Test User"
  });

  it('does not render relation badge when not provided', () => {
    const propsWithoutBadge = {
      ...mockNodeProps,
      data: { ...mockNodeProps.data, relationBadge: undefined },
    };
    render(<MemberNode {...propsWithoutBadge} />);
    expect(screen.queryByText('Father')).not.toBeInTheDocument();
  });
});
