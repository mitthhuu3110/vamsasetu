/**
 * TreeCanvas Component - Example Usage
 * 
 * This file demonstrates how to use the TreeCanvas component
 * in different scenarios.
 */

import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import TreeCanvas from './TreeCanvas';

// Create a query client for React Query
const queryClient = new QueryClient();

/**
 * Example 1: Basic Usage
 * Simply wrap TreeCanvas in QueryClientProvider
 */
export const BasicTreeCanvas: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="w-full h-screen">
        <TreeCanvas />
      </div>
    </QueryClientProvider>
  );
};

/**
 * Example 2: With Custom Container
 * TreeCanvas in a custom-sized container
 */
export const CustomSizedTreeCanvas: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="w-full h-[600px] border border-gray-300 rounded-lg overflow-hidden">
        <TreeCanvas />
      </div>
    </QueryClientProvider>
  );
};

/**
 * Example 3: Full Page Layout
 * TreeCanvas as the main content with header
 */
export const FullPageTreeCanvas: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="flex flex-col h-screen">
        {/* Header */}
        <header className="bg-saffron text-white p-4 shadow-md">
          <h1 className="text-2xl font-bold">Family Tree</h1>
        </header>

        {/* Tree Canvas */}
        <main className="flex-1">
          <TreeCanvas />
        </main>
      </div>
    </QueryClientProvider>
  );
};

/**
 * Example 4: Mobile-Optimized Layout
 * TreeCanvas with mobile-friendly controls
 */
export const MobileTreeCanvas: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="flex flex-col h-screen">
        {/* Mobile Header */}
        <header className="bg-saffron text-white p-3 shadow-md">
          <div className="flex items-center justify-between">
            <h1 className="text-lg font-bold">Family Tree</h1>
            <button className="px-3 py-1 bg-white text-saffron rounded-md text-sm">
              Add Member
            </button>
          </div>
        </header>

        {/* Tree Canvas - Full height on mobile */}
        <main className="flex-1 relative">
          <TreeCanvas />
          
          {/* Mobile hint overlay */}
          <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 bg-black/70 text-white px-4 py-2 rounded-full text-sm">
            Pinch to zoom • Drag to pan
          </div>
        </main>
      </div>
    </QueryClientProvider>
  );
};

/**
 * Example 5: With Side Panel
 * TreeCanvas with a side panel for member details
 */
export const TreeCanvasWithSidePanel: React.FC = () => {
  const [selectedMember, setSelectedMember] = React.useState<string | null>(null);

  return (
    <QueryClientProvider client={queryClient}>
      <div className="flex h-screen">
        {/* Tree Canvas */}
        <div className="flex-1">
          <TreeCanvas />
        </div>

        {/* Side Panel */}
        {selectedMember && (
          <aside className="w-80 bg-white border-l border-gray-200 p-6 overflow-y-auto">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-bold text-charcoal">Member Details</h2>
              <button
                onClick={() => setSelectedMember(null)}
                className="text-gray-500 hover:text-gray-700"
              >
                ✕
              </button>
            </div>
            <p className="text-gray-600">Selected: {selectedMember}</p>
            {/* Add more member details here */}
          </aside>
        )}
      </div>
    </QueryClientProvider>
  );
};

export default BasicTreeCanvas;
