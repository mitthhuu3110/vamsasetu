/**
 * Standard API response format used across all backend endpoints
 * Ensures consistent response structure for success and error cases
 */
export interface APIResponse<T> {
  success: boolean;
  data: T | null;
  error: string;
}
