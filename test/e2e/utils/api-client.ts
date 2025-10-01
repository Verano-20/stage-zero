import { APIRequestContext, APIResponse } from '@playwright/test';

export interface UserData {
  email: string;
  password: string;
}

export interface SimpleResourceData {
  name: string;
}

export interface ApiResponse<T = any> {
  message: string;
  data?: T;
  error?: string;
}

export interface UserResponse {
  id: number;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface LoginResponse {
  token: string;
}

export interface SimpleResourceResponse {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export class ApiClient {
  private request: APIRequestContext;
  private baseURL: string;
  private authToken: string | null = null;

  constructor(request: APIRequestContext, baseURL = 'http://localhost:8080') {
    this.request = request;
    this.baseURL = baseURL;
  }

  /**
   * Set authentication token for subsequent requests
   */
  setAuthToken(token: string): void {
    this.authToken = token;
  }

  /**
   * Get headers with optional authentication
   */
  private getHeaders(additionalHeaders: Record<string, string> = {}): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...additionalHeaders
    };

    if (this.authToken) {
      headers['Authorization'] = `Bearer ${this.authToken}`;
    }

    return headers;
  }

  /**
   * Health check endpoint
   */
  async healthCheck(): Promise<APIResponse> {
    return await this.request.get(`${this.baseURL}/health`);
  }

  /**
   * Sign up a new user
   */
  async signUp(userData: UserData): Promise<APIResponse> {
    return await this.request.post(`${this.baseURL}/auth/signup`, {
      headers: this.getHeaders(),
      data: userData
    });
  }

  /**
   * Login user and optionally set auth token
   */
  async login(credentials: UserData, setToken = true): Promise<APIResponse> {
    const response = await this.request.post(`${this.baseURL}/auth/login`, {
      headers: this.getHeaders(),
      data: credentials
    });

    if (setToken && response.ok()) {
      const responseBody: ApiResponse<LoginResponse> = await response.json();
      if (responseBody.data?.token) {
        this.setAuthToken(responseBody.data.token);
      }
    }

    return response;
  }

  /**
   * Create a new simple resource
   */
  async createSimple(resourceData: SimpleResourceData): Promise<APIResponse> {
    return await this.request.post(`${this.baseURL}/simple`, {
      headers: this.getHeaders(),
      data: resourceData
    });
  }

  /**
   * Get all simple resources
   */
  async getAllSimples(): Promise<APIResponse> {
    return await this.request.get(`${this.baseURL}/simple`, {
      headers: this.getHeaders()
    });
  }

  /**
   * Get a simple resource by ID
   */
  async getSimpleById(id: number | string): Promise<APIResponse> {
    return await this.request.get(`${this.baseURL}/simple/${id}`, {
      headers: this.getHeaders()
    });
  }

  /**
   * Update a simple resource
   */
  async updateSimple(id: number | string, updateData: SimpleResourceData): Promise<APIResponse> {
    return await this.request.put(`${this.baseURL}/simple/${id}`, {
      headers: this.getHeaders(),
      data: updateData
    });
  }

  /**
   * Delete a simple resource
   */
  async deleteSimple(id: number | string): Promise<APIResponse> {
    return await this.request.delete(`${this.baseURL}/simple/${id}`, {
      headers: this.getHeaders()
    });
  }

  /**
   * Clear authentication token
   */
  clearAuth(): void {
    this.authToken = null;
  }
}
