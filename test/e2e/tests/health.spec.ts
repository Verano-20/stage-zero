import { test, expect } from '@playwright/test';
import { ApiClient } from '../utils/api-client';
import { assertResponse } from '../utils/test-helpers';
import { expectedResponses } from '../fixtures/test-data';

test.describe('Health Check API', () => {
  let apiClient: ApiClient;

  test.beforeEach(async ({ request }) => {
    apiClient = new ApiClient(request);
  });

  test('should return healthy status', async () => {
    const response = await apiClient.healthCheck();
    
    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);
    
    const body = await response.json();
    expect(body).toHaveProperty('status', expectedResponses.health.status);
  });

  test('should respond quickly', async () => {
    const startTime = Date.now();
    const response = await apiClient.healthCheck();
    const endTime = Date.now();
    
    expect(response.ok()).toBeTruthy();
    expect(endTime - startTime).toBeLessThan(1000); // Should respond within 1 second
  });

  test('should have correct content type', async () => {
    const response = await apiClient.healthCheck();
    
    expect(response.ok()).toBeTruthy();
    expect(response.headers()['content-type']).toContain('application/json');
  });
});
