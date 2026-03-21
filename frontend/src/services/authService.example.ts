/**
 * Example usage of the AuthService
 * This file demonstrates how to use the authentication service methods
 */

import authService from './authService';

// Example 1: Register a new user
async function registerExample() {
  const response = await authService.register({
    email: 'user@example.com',
    password: 'SecurePassword123!',
    name: 'John Doe',
  });

  if (response.success && response.data) {
    console.log('Registration successful!');
    console.log('User:', response.data.user);
    console.log('Access Token:', response.data.accessToken);
    // Tokens are automatically stored in localStorage
  } else {
    console.error('Registration failed:', response.error);
  }
}

// Example 2: Login an existing user
async function loginExample() {
  const response = await authService.login({
    email: 'user@example.com',
    password: 'SecurePassword123!',
  });

  if (response.success && response.data) {
    console.log('Login successful!');
    console.log('User:', response.data.user);
    // Tokens are automatically stored in localStorage
  } else {
    console.error('Login failed:', response.error);
  }
}

// Example 3: Refresh access token
async function refreshTokenExample() {
  const refreshToken = authService.getRefreshToken();
  
  if (refreshToken) {
    const response = await authService.refresh(refreshToken);

    if (response.success && response.data) {
      console.log('Token refresh successful!');
      console.log('New Access Token:', response.data.accessToken);
      // New tokens are automatically stored in localStorage
    } else {
      console.error('Token refresh failed:', response.error);
    }
  }
}

// Example 4: Get current user profile
async function getProfileExample() {
  const response = await authService.getProfile();

  if (response.success && response.data) {
    console.log('Profile retrieved successfully!');
    console.log('User:', response.data);
  } else {
    console.error('Failed to get profile:', response.error);
  }
}

// Example 5: Check authentication status
function checkAuthExample() {
  if (authService.isAuthenticated()) {
    console.log('User is authenticated');
    console.log('Access Token:', authService.getAccessToken());
  } else {
    console.log('User is not authenticated');
  }
}

// Example 6: Logout
function logoutExample() {
  authService.logout();
  console.log('User logged out, tokens cleared');
}

// Example 7: Complete authentication flow
async function completeAuthFlowExample() {
  // 1. Register
  const registerResponse = await authService.register({
    email: 'newuser@example.com',
    password: 'Password123!',
    name: 'Jane Smith',
  });

  if (!registerResponse.success) {
    console.error('Registration failed:', registerResponse.error);
    return;
  }

  // 2. Get profile
  const profileResponse = await authService.getProfile();
  console.log('User profile:', profileResponse.data);

  // 3. Logout
  authService.logout();

  // 4. Login again
  const loginResponse = await authService.login({
    email: 'newuser@example.com',
    password: 'Password123!',
  });

  if (loginResponse.success) {
    console.log('Successfully logged back in!');
  }
}

export {
  registerExample,
  loginExample,
  refreshTokenExample,
  getProfileExample,
  checkAuthExample,
  logoutExample,
  completeAuthFlowExample,
};
