import { test, expect } from '@playwright/test';
import { ApiClient, UserData, LoginResponse, UserResponse } from '../utils/api-client';
import {
  generateUserData,
  assertResponse,
  assertErrorResponse
} from '../utils/test-helpers';
import { expectedResponses } from '../fixtures/test-data';

test.describe('Authentication API', () => {
  let apiClient: ApiClient;

  test.beforeEach(async ({ request }) => {
    apiClient = new ApiClient(request);
  });

  test.describe('User Signup', () => {
    test('should successfully create a new user', async () => {
      const userData = generateUserData();

      const response = await apiClient.signUp(userData);
      expect(response.ok()).toBeTruthy();

      const body = await assertResponse<UserResponse>(response, 201);
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

      const body = await assertErrorResponse(secondResponse, 409);
      expect(body.error).toContain('already exists');
    });

    test('should reject invalid email format', async () => {
      const userData = generateUserData({ email: 'invalid-email' });

      const response = await apiClient.signUp(userData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 400);
    });

    test('should reject weak password', async () => {
      const userData = generateUserData({ password: '123' });

      const response = await apiClient.signUp(userData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 400);
    });

    test('should reject missing email', async () => {
      const userData = { password: 'ValidPassword123!' } as UserData;

      const response = await apiClient.signUp(userData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 400);
    });

    test('should reject missing password', async () => {
      const userData = { email: 'test@example.com' } as UserData;

      const response = await apiClient.signUp(userData);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 400);
    });
  });

  test.describe('User Login', () => {
    let registeredUser: UserData;

    test.beforeEach(async () => {
      // Create a user for login tests
      registeredUser = generateUserData();
      const signupResponse = await apiClient.signUp(registeredUser);
      expect(signupResponse.ok()).toBeTruthy();
    });

    test('should successfully login with valid credentials', async () => {
      const response = await apiClient.login(registeredUser, false);
      expect(response.ok()).toBeTruthy();

      const body = await assertResponse<LoginResponse>(response, 200);
      expect(body.message).toBe(expectedResponses.authSuccess.login.message);
      expect(body.data).toHaveProperty('token');
      expect(typeof body.data!.token).toBe('string');
      expect(body.data!.token.length).toBeGreaterThan(0);
    });

    test('should reject invalid email', async () => {
      const invalidCredentials: UserData = {
        email: 'nonexistent@example.com',
        password: registeredUser.password
      };

      const response = await apiClient.login(invalidCredentials, false);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 401);
    });

    test('should reject invalid password', async () => {
      const invalidCredentials: UserData = {
        email: registeredUser.email,
        password: 'wrongpassword'
      };

      const response = await apiClient.login(invalidCredentials, false);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 401);
    });

    test('should reject missing email', async () => {
      const invalidCredentials = {
        password: registeredUser.password
      } as UserData;

      const response = await apiClient.login(invalidCredentials, false);
      expect(response.ok()).toBeFalsy();

      await assertErrorResponse(response, 400);
    });

    test('should reject missing password', async () => {
      const invalidCredentials = {
        email: registeredUser.email
      } as UserData;

      const response = await apiClient.login(invalidCredentials, false);
      expect(response.ok()).toBeFalsy();

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
      tokenParts.forEach((part: string) => {
        expect(part).toMatch(/^[A-Za-z0-9_-]+$/);
      });
    });
  });
});
