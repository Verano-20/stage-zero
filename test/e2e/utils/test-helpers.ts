import { APIResponse } from '@playwright/test';
import { UserData, SimpleResourceData, ApiResponse, ApiClient } from './api-client';

/**
 * Generate a unique email for testing
 */
export function generateUniqueEmail(prefix = 'test'): string {
  const timestamp = Date.now();
  const random = Math.random().toString(36).substring(2, 8);
  return `${prefix}-${timestamp}-${random}@example.com`;
}

/**
 * Generate a secure password for testing
 */
export function generatePassword(length = 12): string {
  const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';
  let password = '';
  for (let i = 0; i < length; i++) {
    password += charset.charAt(Math.floor(Math.random() * charset.length));
  }
  return password;
}

/**
 * Generate test user data
 */
export function generateUserData(overrides: Partial<UserData> = {}): UserData {
  return {
    email: generateUniqueEmail(),
    password: generatePassword(),
    ...overrides
  };
}

/**
 * Generate test simple resource data
 */
export function generateSimpleData(overrides: Partial<SimpleResourceData> = {}): SimpleResourceData {
  const timestamp = Date.now();
  const random = Math.random().toString(36).substring(2, 8);
  
  return {
    name: `Test Resource ${timestamp}-${random}`,
    ...overrides
  };
}

/**
 * Wait for a condition to be true with timeout
 */
export async function waitForCondition(
  condition: () => Promise<boolean> | boolean,
  timeout = 5000,
  interval = 100
): Promise<void> {
  const startTime = Date.now();
  
  while (Date.now() - startTime < timeout) {
    if (await condition()) {
      return;
    }
    await new Promise(resolve => setTimeout(resolve, interval));
  }
  
  throw new Error(`Condition not met within ${timeout}ms timeout`);
}

/**
 * Retry an operation with exponential backoff
 */
export async function retryOperation<T>(
  operation: () => Promise<T>,
  maxRetries = 3,
  baseDelay = 1000
): Promise<T> {
  let lastError: Error;
  
  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await operation();
    } catch (error) {
      lastError = error as Error;
      
      if (attempt === maxRetries) {
        break;
      }
      
      const delay = baseDelay * Math.pow(2, attempt);
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }
  
  throw lastError!;
}

/**
 * Assert that a response has the expected status and structure
 */
export async function assertResponse<T = any>(
  response: APIResponse,
  expectedStatus = 200,
  shouldHaveData = true
): Promise<ApiResponse<T>> {
  const body: ApiResponse<T> = await response.json();
  
  if (response.status() !== expectedStatus) {
    throw new Error(`Expected status ${expectedStatus}, got ${response.status()}. Response: ${JSON.stringify(body, null, 2)}`);
  }
  
  if (shouldHaveData && !body.data) {
    throw new Error(`Expected response to have 'data' field. Response: ${JSON.stringify(body, null, 2)}`);
  }
  
  return body;
}

/**
 * Assert that a response is an error with expected status
 */
export async function assertErrorResponse(
  response: APIResponse,
  expectedStatus: number
): Promise<ApiResponse> {
  const body: ApiResponse = await response.json();
  
  if (response.status() !== expectedStatus) {
    throw new Error(`Expected error status ${expectedStatus}, got ${response.status()}. Response: ${JSON.stringify(body, null, 2)}`);
  }
  
  if (!body.error) {
    throw new Error(`Expected response to have 'error' field. Response: ${JSON.stringify(body, null, 2)}`);
  }
  
  return body;
}

/**
 * Clean up test data (placeholder for future implementation)
 */
export async function cleanupTestData(apiClient: ApiClient, resourceIds: (string | number)[] = []): Promise<void> {
  // Implementation depends on your cleanup strategy
  // This could involve deleting created resources, clearing database, etc.
  console.log('Cleaning up test data...', resourceIds);
}
