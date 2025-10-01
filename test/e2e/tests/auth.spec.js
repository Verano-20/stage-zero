// auth.spec.js - Authentication endpoint tests

const { test, expect } = require('@playwright/test');
const { ApiClient } = require('../utils/api-client');
const { 
  generateUserData, 
  assertResponse, 
  assertErrorResponse 
} = require('../utils/test-helpers');
const { 
  testUsers, 
  expectedResponses, 
  apiEndpoints 
} = require('../fixtures/test-data');

test.describe('Authentication API', () => {
  let apiClient;

  test.beforeEach(async ({ request }) => {
    apiClient = new ApiClient(request);
  });

  test.describe('User Signup', () => {
    test('should successfully create a new user', async () => {
      const userData = generateUserData();
      
      const response = await apiClient.signUp(userData);
      
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(201);
      
      const body = await assertResponse(response, 201);
      expect(body.message).toBe(expectedResponses.authSuccess.signup.message);
      expect(body.data).toHaveProperty('id');
      expect(body.data).toHaveProperty('email', userData.email);
      expect(body.data).not.toHaveProperty('password'); // Password should not be returned
    });

    test('should reject duplicate email registration', async () => {
      const userData = generateUserData();
      
      // First registration should succeed
      const firstResponse = await apiClient.signUp(userData);
      expect(firstResponse.ok()).toBeTruthy();
      
      // Second registration with same email should fail
      const secondResponse = await apiClient.signUp(userData);
      expect(secondResponse.ok()).toBeFalsy();
      expect(secondResponse.status()).toBe(409); // Conflict
      
      const body = await assertErrorResponse(secondResponse, 409);
      expect(body.error).toContain('already exists');
    });

    test('should reject invalid email format', async () => {
      const userData = generateUserData({ email: 'invalid-email' });
      
      const response = await apiClient.signUp(userData);
      
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(400);
      
      await assertErrorResponse(response, 400);
    });

    test('should reject weak password', async () => {
      const userData = generateUserData({ password: '123' });
      
      const response = await apiClient.signUp(userData);
      
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(400);
      
      await assertErrorResponse(response, 400);
    });

    test('should reject missing email', async () => {
      const userData = { password: 'ValidPassword123!' };
      
      const response = await apiClient.signUp(userData);
      
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(400);
      
      await assertErrorResponse(response, 400);
    });

    test('should reject missing password', async () => {
      const userData = { email: 'test@example.com' };
      
      const response = await apiClient.signUp(userData);
      
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(400);
      
      await assertErrorResponse(response, 400);
    });
  });

  test.describe('User Login', () => {
    let registeredUser;

    test.beforeEach(async () => {
      // Create a user for login tests
      registeredUser = generateUserData();
      const signupResponse = await apiClient.signUp(registeredUser);
      expect(signupResponse.ok()).toBeTruthy();
    });

    test('should successfully login with valid credentials', async () => {
      const response = await apiClient.login(registeredUser, false);
      
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      
      const body = await assertResponse(response, 200);
      expect(body.message).toBe(expectedResponses.authSuccess.login.message);
      expect(body.data).toHaveProperty('token');
      expect(typeof body.data.token).toBe('string');
      expect(body.data.token.length).toBeGreaterThan(0);
    });

    test('should reject invalid email', async () => {
      const invalidCredentials = {
        email: 'nonexistent@example.com',
        password: registeredUser.password
      };
      
      const response = await apiClient.login(invalidCredentials, false);
      
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(401);
      
      await assertErrorResponse(response, 401);
    });

    test('should reject invalid password', async () => {
      const invalidCredentials = {
        email: registeredUser.email,
        password: 'wrongpassword'
      };
      
      const response = await apiClient.login(invalidCredentials, false);
      
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(401);
      
      await assertErrorResponse(response, 401);
    });

    test('should reject missing email', async () => {
      const invalidCredentials = {
        password: registeredUser.password
      };
      
      const response = await apiClient.login(invalidCredentials, false);
      
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(400);
      
      await assertErrorResponse(response, 400);
    });

    test('should reject missing password', async () => {
      const invalidCredentials = {
        email: registeredUser.email
      };
      
      const response = await apiClient.login(invalidCredentials, false);
      
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(400);
      
      await assertErrorResponse(response, 400);
    });

    test('should return valid JWT token', async () => {
      const response = await apiClient.login(registeredUser, false);
      
      expect(response.ok()).toBeTruthy();
      
      const body = await response.json();
      const token = body.data.token;
      
      // Basic JWT format validation (header.payload.signature)
      const tokenParts = token.split('.');
      expect(tokenParts).toHaveLength(3);
      
      // Each part should be base64 encoded
      tokenParts.forEach(part => {
        expect(part).toMatch(/^[A-Za-z0-9_-]+$/);
      });
    });
  });
});
