import { test, expect } from '@playwright/test';
import { ApiClient, UserData, SimpleResourceResponse } from '../utils/api-client';
import {
  generateUserData,
  generateSimpleData,
  assertResponse,
  assertErrorResponse
} from '../utils/test-helpers';
import { expectedResponses } from '../fixtures/test-data';

test.describe('Simple Resource CRUD API', () => {
  let apiClient: ApiClient;
  let authenticatedUser: UserData;

  test.beforeAll(async ({ request }) => {
    // Set up authenticated user for all tests
    apiClient = new ApiClient(request);
    authenticatedUser = generateUserData();

    // Register and login user
    const signupResponse = await apiClient.signUp(authenticatedUser);
    expect(signupResponse.ok()).toBeTruthy();

    const loginResponse = await apiClient.login(authenticatedUser, true);
    expect(loginResponse.ok()).toBeTruthy();
  });

  test.beforeEach(async ({ request }) => {
    // Ensure we have a fresh API client with auth token for each test
    apiClient = new ApiClient(request);
    const loginResponse = await apiClient.login(authenticatedUser, true);
    expect(loginResponse.ok()).toBeTruthy();
  });

  test.describe('Create Simple Resource', () => {
    test('should successfully create a simple resource', async () => {
      const resourceData = generateSimpleData();

      const response = await apiClient.createSimple(resourceData);
      expect(response.ok()).toBeTruthy();

      const body = await assertResponse<SimpleResourceResponse>(response, 201);
      expect(body.message).toBe(expectedResponses.simpleSuccess.create.message);
      expect(body.data).toHaveProperty('id');
      expect(body.data).toHaveProperty('name', resourceData.name);
      expect(body.data).toHaveProperty('created_at');
      expect(body.data).toHaveProperty('updated_at');
    });

    test('should reject creation without authentication', async () => {
      const resourceData = generateSimpleData();
      apiClient.clearAuth(); // Remove authentication

      const response = await apiClient.createSimple(resourceData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 401);
    });

    test('should reject creation with invalid data', async () => {
      const invalidData = { not_a_field: 'value' };

      const response = await apiClient.createSimpleWithCustomPayload(invalidData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 400);
    });

    test('should reject creation with missing name', async () => {
      const invalidData = {} as any; // Missing name field

      const response = await apiClient.createSimple(invalidData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 400);
    });
  });

  test.describe('Read Simple Resources', () => {
    let createdResource: SimpleResourceResponse;

    test.beforeEach(async () => {
      // Create a resource for read tests
      const resourceData = generateSimpleData();
      const createResponse = await apiClient.createSimple(resourceData);
      expect(createResponse.ok()).toBeTruthy();

      const createBody = await createResponse.json();
      createdResource = createBody.data;
    });

    test('should get all simple resources', async () => {
      const response = await apiClient.getAllSimples();
      expect(response.ok()).toBeTruthy();

      const body = await assertResponse<SimpleResourceResponse[]>(response, 200);
      expect(body.message).toBe(expectedResponses.simpleSuccess.getAll.message);
      expect(Array.isArray(body.data)).toBeTruthy();
      expect(body.data!.length).toBeGreaterThan(0);

      // Check if our created resource is in the list
      const foundResource = body.data!.find(resource => resource.id === createdResource.id);
      expect(foundResource).toBeDefined();
      expect(foundResource!.name).toBe(createdResource.name);
    });

    test('should get simple resource by ID', async () => {
      const response = await apiClient.getSimpleById(createdResource.id);
      expect(response.ok()).toBeTruthy();

      const body = await assertResponse<SimpleResourceResponse>(response, 200);
      expect(body.message).toBe(expectedResponses.simpleSuccess.getById.message);
      expect(body.data).toHaveProperty('id', createdResource.id);
      expect(body.data).toHaveProperty('name', createdResource.name);
    });

    test('should return 404 for non-existent resource', async () => {
      const nonExistentId = 999999;

      const response = await apiClient.getSimpleById(nonExistentId);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 404);
    });

    test('should reject get all without authentication', async () => {
      apiClient.clearAuth();

      const response = await apiClient.getAllSimples();
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 401);
    });

    test('should reject get by ID without authentication', async () => {
      apiClient.clearAuth();

      const response = await apiClient.getSimpleById(createdResource.id);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 401);
    });
  });

  test.describe('Update Simple Resource', () => {
    let createdResource: SimpleResourceResponse;

    test.beforeEach(async () => {
      // Create a resource for update tests
      const resourceData = generateSimpleData();
      const createResponse = await apiClient.createSimple(resourceData);
      expect(createResponse.ok()).toBeTruthy();

      const createBody = await createResponse.json();
      createdResource = createBody.data;
    });

    test('should successfully update a simple resource', async () => {
      const updateData = { name: 'Updated Resource Name' };

      const response = await apiClient.updateSimple(createdResource.id, updateData);
      expect(response.ok()).toBeTruthy();

      const body = await assertResponse<SimpleResourceResponse>(response, 200);
      expect(body.message).toBe(expectedResponses.simpleSuccess.update.message);
      expect(body.data).toHaveProperty('id', createdResource.id);
      expect(body.data).toHaveProperty('name', updateData.name);
      expect(body.data!.updated_at).not.toBe(createdResource.updated_at);
    });

    test('should return 404 for updating non-existent resource', async () => {
      const nonExistentId = 999999;
      const updateData = { name: 'Updated Name' };

      const response = await apiClient.updateSimple(nonExistentId, updateData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 404);
    });

    test('should reject update without authentication', async () => {
      const updateData = { name: 'Updated Name' };
      apiClient.clearAuth();

      const response = await apiClient.updateSimple(createdResource.id, updateData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 401);
    });

    test('should reject update with invalid data', async () => {
      const invalidData = { not_a_field: 'value' };

      const response = await apiClient.updateSimpleWithCustomPayload(createdResource.id, invalidData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 400);
    });
  });

  test.describe('Delete Simple Resource', () => {
    let createdResource: SimpleResourceResponse;

    test.beforeEach(async () => {
      // Create a resource for delete tests
      const resourceData = generateSimpleData();
      const createResponse = await apiClient.createSimple(resourceData);
      expect(createResponse.ok()).toBeTruthy();

      const createBody = await createResponse.json();
      createdResource = createBody.data;
    });

    test('should successfully delete a simple resource', async () => {
      const response = await apiClient.deleteSimple(createdResource.id);
      expect(response.ok()).toBeTruthy();

      const body = await assertResponse(response, 200, false);
      expect(body.message).toBe(expectedResponses.simpleSuccess.delete.message);

      // Verify resource is actually deleted
      const getResponse = await apiClient.getSimpleById(createdResource.id);
      expect(getResponse.status()).toBe(404);
    });

    test('should return 404 for deleting non-existent resource', async () => {
      const nonExistentId = 999999;

      const response = await apiClient.deleteSimple(nonExistentId);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 404);
    });

    test('should reject delete without authentication', async () => {
      apiClient.clearAuth();

      const response = await apiClient.deleteSimple(createdResource.id);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 401);
    });
  });

  test.describe('Complete CRUD Workflow', () => {
    test('should perform complete CRUD lifecycle', async () => {
      const resourceData = generateSimpleData();

      // 1. Create
      const createResponse = await apiClient.createSimple(resourceData);
      expect(createResponse.ok()).toBeTruthy();
      const createBody = await createResponse.json();
      const resourceId = createBody.data.id;

      // 2. Read (Get by ID)
      const getResponse = await apiClient.getSimpleById(resourceId);
      expect(getResponse.ok()).toBeTruthy();
      const getBody = await getResponse.json();
      expect(getBody.data.name).toBe(resourceData.name);

      // 3. Update
      const updateData = { name: 'Updated in CRUD workflow' };
      const updateResponse = await apiClient.updateSimple(resourceId, updateData);
      expect(updateResponse.ok()).toBeTruthy();
      const updateBody = await updateResponse.json();
      expect(updateBody.data.name).toBe(updateData.name);

      // 4. Delete
      const deleteResponse = await apiClient.deleteSimple(resourceId);
      expect(deleteResponse.ok()).toBeTruthy();

      // 5. Verify deletion
      const verifyResponse = await apiClient.getSimpleById(resourceId);
      expect(verifyResponse.status()).toBe(404);
    });
  });
});
