import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ReactFlowProvider } from 'reactflow';
import TreeControls from './TreeControls';

// Mock useReactFlow hook
const mockZoomIn = vi.fn();
const mockZoomOut = vi.fn();
const mockFitView = vi.fn();

vi.mock('reactflow', async () => {
  const actual = await vi.importActual('reactflow');
  return {
    ...actual,
    useReactFlow: () => ({
      zoomIn: mockZoomIn,
      zoomOut: mockZoomOut,
      fitView: mockFitView,
    }),
  };
});

describe('TreeControls', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const renderWithProvider = (ui: React.ReactElement) => {
    return render(<ReactFlowProvider>{ui}</ReactFlowProvider>);
  };

  it('renders zoom controls', () => {
    renderWithProvider(<TreeControls />);
    
    expect(screen.getByLabelText('Zoom in')).toBeInTheDocument();
    expect(screen.getByLabelText('Zoom out')).toBeInTheDocument();
    expect(screen.getByLabelText('Fit view')).toBeInTheDocument();
  });

  it('calls zoomIn when zoom in button is clicked', () => {
    renderWithProvider(<TreeControls />);
    
    const zoomInButton = screen.getByLabelText('Zoom in');
    fireEvent.click(zoomInButton);
    
    expect(mockZoomIn).toHaveBeenCalledWith({ duration: 300 });
  });

  it('calls zoomOut when zoom out button is clicked', () => {
    renderWithProvider(<TreeControls />);
    
    const zoomOutButton = screen.getByLabelText('Zoom out');
    fireEvent.click(zoomOutButton);
    
    expect(mockZoomOut).toHaveBeenCalledWith({ duration: 300 });
  });

  it('calls fitView when fit view button is clicked', () => {
    renderWithProvider(<TreeControls />);
    
    const fitViewButton = screen.getByLabelText('Fit view');
    fireEvent.click(fitViewButton);
    
    expect(mockFitView).toHaveBeenCalledWith({ duration: 300, padding: 0.2 });
  });

  it('renders add member button when onAddMember prop is provided', () => {
    const onAddMember = vi.fn();
    renderWithProvider(<TreeControls onAddMember={onAddMember} />);
    
    expect(screen.getByLabelText('Add family member')).toBeInTheDocument();
  });

  it('does not render add member button when onAddMember prop is not provided', () => {
    renderWithProvider(<TreeControls />);
    
    expect(screen.queryByLabelText('Add family member')).not.toBeInTheDocument();
  });

  it('calls onAddMember when add member button is clicked', () => {
    const onAddMember = vi.fn();
    renderWithProvider(<TreeControls onAddMember={onAddMember} />);
    
    const addMemberButton = screen.getByLabelText('Add family member');
    fireEvent.click(addMemberButton);
    
    expect(onAddMember).toHaveBeenCalledTimes(1);
  });

  it('renders add relationship button when onAddRelationship prop is provided', () => {
    const onAddRelationship = vi.fn();
    renderWithProvider(<TreeControls onAddRelationship={onAddRelationship} />);
    
    expect(screen.getByLabelText('Add relationship')).toBeInTheDocument();
  });

  it('does not render add relationship button when onAddRelationship prop is not provided', () => {
    renderWithProvider(<TreeControls />);
    
    expect(screen.queryByLabelText('Add relationship')).not.toBeInTheDocument();
  });

  it('calls onAddRelationship when add relationship button is clicked', () => {
    const onAddRelationship = vi.fn();
    renderWithProvider(<TreeControls onAddRelationship={onAddRelationship} />);
    
    const addRelationshipButton = screen.getByLabelText('Add relationship');
    fireEvent.click(addRelationshipButton);
    
    expect(onAddRelationship).toHaveBeenCalledTimes(1);
  });

  it('applies custom className', () => {
    const { container } = renderWithProvider(
      <TreeControls className="custom-class" />
    );
    
    const toolbar = container.querySelector('[role="toolbar"]');
    expect(toolbar).toHaveClass('custom-class');
  });

  it('has proper accessibility attributes', () => {
    renderWithProvider(<TreeControls />);
    
    const toolbar = screen.getByRole('toolbar');
    expect(toolbar).toHaveAttribute('aria-label', 'Family tree controls');
  });

  it('renders all buttons when both callbacks are provided', () => {
    const onAddMember = vi.fn();
    const onAddRelationship = vi.fn();
    
    renderWithProvider(
      <TreeControls
        onAddMember={onAddMember}
        onAddRelationship={onAddRelationship}
      />
    );
    
    expect(screen.getByLabelText('Zoom in')).toBeInTheDocument();
    expect(screen.getByLabelText('Zoom out')).toBeInTheDocument();
    expect(screen.getByLabelText('Fit view')).toBeInTheDocument();
    expect(screen.getByLabelText('Add family member')).toBeInTheDocument();
    expect(screen.getByLabelText('Add relationship')).toBeInTheDocument();
  });
});
