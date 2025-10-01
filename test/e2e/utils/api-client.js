// api-client.js - Utility functions for API interactions

class ApiClient {
  constructor(request, baseURL = 'http://localhost:8080') {
    this.request = request;
    this.baseURL = baseURL;
    this.authToken = null;
  }

  /**
   * Set authentication token for subsequent requests
   * @param {string} token - JWT token
   */
  setAuthToken(token) {
    this.authToken = token;
  }

  /**
   * Get headers with optional authentication
   * @param {Object} additionalHeaders - Additional headers to include
   * @returns {Object} Headers object
   */
  getHeaders(additionalHeaders = {}) {
    const headers = {
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
   * @returns {Promise<Response>} API response
   */
  async healthCheck() {
    return await this.request.get(`${this.baseURL}/health`);
  }

  /**
   * Sign up a new user
   * @param {Object} userData - User registration data
   * @param {string} userData.email - User email
   * @param {string} userData.password - User password
   * @returns {Promise<Response>} API response
   */
  async signUp(userData) {
    return await this.request.post(`${this.baseURL}/auth/signup`, {
      headers: this.getHeaders(),
      data: userData
    });
  }

  /**
   * Login user and optionally set auth token
   * @param {Object} credentials - Login credentials
   * @param {string} credentials.email - User email
   * @param {string} credentials.password - User password
   * @param {boolean} setToken - Whether to automatically set the auth token
   * @returns {Promise<Response>} API response
   */
  async login(credentials, setToken = true) {
    const response = await this.request.post(`${this.baseURL}/auth/login`, {
      headers: this.getHeaders(),
      data: credentials
    });

    if (setToken && response.ok()) {
      const responseBody = await response.json();
      if (responseBody.data && responseBody.data.token) {
        this.setAuthToken(responseBody.data.token);
      }
    }

    return response;
  }

  /**
   * Create a new simple resource
   * @param {Object} resourceData - Resource data
   * @param {string} resourceData.name - Resource name
   * @returns {Promise<Response>} API response
   */
  async createSimple(resourceData) {
    return await this.request.post(`${this.baseURL}/simple`, {
      headers: this.getHeaders(),
      data: resourceData
    });
  }

  /**
   * Get all simple resources
   * @returns {Promise<Response>} API response
   */
  async getAllSimples() {
    return await this.request.get(`${this.baseURL}/simple`, {
      headers: this.getHeaders()
    });
  }

  /**
   * Get a simple resource by ID
   * @param {number|string} id - Resource ID
   * @returns {Promise<Response>} API response
   */
  async getSimpleById(id) {
    return await this.request.get(`${this.baseURL}/simple/${id}`, {
      headers: this.getHeaders()
    });
  }

  /**
   * Update a simple resource
   * @param {number|string} id - Resource ID
   * @param {Object} updateData - Update data
   * @param {string} updateData.name - Updated resource name
   * @returns {Promise<Response>} API response
   */
  async updateSimple(id, updateData) {
    return await this.request.put(`${this.baseURL}/simple/${id}`, {
      headers: this.getHeaders(),
      data: updateData
    });
  }

  /**
   * Delete a simple resource
   * @param {number|string} id - Resource ID
   * @returns {Promise<Response>} API response
   */
  async deleteSimple(id) {
    return await this.request.delete(`${this.baseURL}/simple/${id}`, {
      headers: this.getHeaders()
    });
  }

  /**
   * Clear authentication token
   */
  clearAuth() {
    this.authToken = null;
  }
}

module.exports = { ApiClient };
