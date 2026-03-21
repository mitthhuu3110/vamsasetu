/**
 * Example usage of the FamilyTreeService
 * This file demonstrates how to use the family tree service methods
 */

import familyTreeService from './familyTreeService';

// Example 1: Get the complete family tree
async function getFamilyTreeExample() {
  const response = await familyTreeService.getTree();

  if (response.success && response.data) {
    console.log('Family tree retrieved successfully!');
    console.log('Total nodes:', response.data.nodes.length);
    console.log('Total edges:', response.data.edges.length);
    console.log('Nodes:', response.data.nodes);
    console.log('Edges:', response.data.edges);
  } else {
    console.error('Failed to get family tree:', response.error);
  }
}

// Example 2: Process family tree data for React Flow
async function processFamilyTreeForReactFlow() {
  const response = await familyTreeService.getTree();

  if (response.success && response.data) {
    const { nodes, edges } = response.data;

    // The data is already in React Flow format
    console.log('React Flow nodes:', nodes);
    console.log('React Flow edges:', edges);

    // You can directly pass this to React Flow
    // <ReactFlow nodes={nodes} edges={edges} />
  } else {
    console.error('Failed to get family tree:', response.error);
  }
}

// Example 3: Check for members with upcoming events
async function checkUpcomingEventsExample() {
  const response = await familyTreeService.getTree();

  if (response.success && response.data) {
    const membersWithEvents = response.data.nodes.filter(
      (node) => node.data.hasUpcomingEvent
    );

    console.log('Members with upcoming events:', membersWithEvents.length);
    membersWithEvents.forEach((node) => {
      console.log(`- ${node.data.name} (${node.data.id})`);
    });
  } else {
    console.error('Failed to get family tree:', response.error);
  }
}

// Example 4: Analyze family tree structure
async function analyzeFamilyTreeExample() {
  const response = await familyTreeService.getTree();

  if (response.success && response.data) {
    const { nodes, edges } = response.data;

    // Count relationship types
    const spouseEdges = edges.filter((edge) => edge.style.stroke === '#E11D48');
    const parentEdges = edges.filter((edge) => edge.style.stroke === '#0D9488');
    const siblingEdges = edges.filter((edge) => edge.style.stroke === '#F59E0B');

    console.log('Family tree statistics:');
    console.log('- Total members:', nodes.length);
    console.log('- Spouse relationships:', spouseEdges.length);
    console.log('- Parent-child relationships:', parentEdges.length);
    console.log('- Sibling relationships:', siblingEdges.length);

    // Count by gender
    const maleMembers = nodes.filter((node) => node.data.gender === 'male');
    const femaleMembers = nodes.filter((node) => node.data.gender === 'female');

    console.log('- Male members:', maleMembers.length);
    console.log('- Female members:', femaleMembers.length);
  } else {
    console.error('Failed to get family tree:', response.error);
  }
}

// Example 5: Complete family tree visualization flow
async function completeFamilyTreeFlowExample() {
  // 1. Fetch the family tree
  const response = await familyTreeService.getTree();

  if (!response.success) {
    console.error('Failed to fetch family tree:', response.error);
    return;
  }

  const { nodes, edges } = response.data!;

  // 2. Log basic information
  console.log('Family tree loaded successfully!');
  console.log(`Found ${nodes.length} members and ${edges.length} relationships`);

  // 3. Find root nodes (members at the top of the tree)
  const rootNodes = nodes.filter((node) => node.position.y === 0);
  console.log('Root members:', rootNodes.map((n) => n.data.name));

  // 4. Find members with upcoming events
  const upcomingEvents = nodes.filter((node) => node.data.hasUpcomingEvent);
  console.log('Members with upcoming events:', upcomingEvents.map((n) => n.data.name));

  // 5. The data is ready to be used with React Flow
  console.log('Ready to render with React Flow!');
}

export {
  getFamilyTreeExample,
  processFamilyTreeForReactFlow,
  checkUpcomingEventsExample,
  analyzeFamilyTreeExample,
  completeFamilyTreeFlowExample,
};
