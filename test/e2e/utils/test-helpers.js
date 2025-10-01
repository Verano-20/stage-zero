// test-helpers.js - Common test utilities and helpers

/**
 * Generate a unique email for testing
 * @param {string} prefix - Email prefix (default: 'test')
 * @returns {string} Unique email address
 */
function generateUniqueEmail(prefix = 'test') {
  const timestamp = Date.now();
  const random = Math.random().toString(36).substring(2, 8);
  return `${prefix}-${timestamp}-${random}@example.com`;
}

/**
 * Generate a secure password for testing
 * @param {number} length - Password length (default: 12)
 * @returns {string} Generated password
 */
function generatePassword(length = 12) {
  const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';
  let password = '';
  for (let i = 0; i < length; i++) {
    password += charset.charAt(Math.floor(Math.random() * charset.length));
  }
  return password;
}

/**
 * Generate test user data
 * @param {Object} overrides - Override default values
 * @returns {Object} User data object
 */
function generateUserData(overrides = {}) {
  return {
    email: generateUniqueEmail(),
    password: generatePassword(),
    ...overrides
  };
}

/**
 * Generate test simple resource data
 * @param {Object} overrides - Override default values
 * @returns {Object} Simple resource data object
 */
function generateSimpleData(overrides = {}) {
  const timestamp = Date.now();
  const random = Math.random().toString(36).substring(2, 8);
  
  return {
    name: `Test Simple ${timestamp}-${random}`,
    ...overrides
  };
}

/**
 * Wait for a condition to be true with timeout
 * @param {Function} condition - Function that returns boolean
 * @param {number} timeout - Timeout in milliseconds (default: 5000)
 * @param {number} interval - Check interval in milliseconds (default: 100)
 * @returns {Promise<void>}
 */
async function waitForCondition(condition, timeout = 5000, interval = 100) {
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
 * @param {Function} operation - Async operation to retry
 * @param {number} maxRetries - Maximum number of retries (default: 3)
 * @param {number} baseDelay - Base delay in milliseconds (default: 1000)
 * @returns {Promise<any>} Result of the operation
 */
async function retryOperation(operation, maxRetries = 3, baseDelay = 1000) {
  let lastError;
  
  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await operation();
    } catch (error) {
      lastError = error;
      
      if (attempt === maxRetries) {
        break;
      }
      
      const delay = baseDelay * Math.pow(2, attempt);
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }
  
  throw lastError;
}

/**
 * Assert that a response has the expected status and structure
 * @param {Response} response - API response
 * @param {number} expectedStatus - Expected HTTP status code
 * @param {boolean} shouldHaveData - Whether response should have data field
 * @returns {Promise<Object>} Response body
 */
async function assertResponse(response, expectedStatus = 200, shouldHaveData = true) {
  const body = await response.json();
  
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
 * @param {Response} response - API response
 * @param {number} expectedStatus - Expected HTTP status code
 * @returns {Promise<Object>} Response body
 */
async function assertErrorResponse(response, expectedStatus) {
  const body = await response.json();
  
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
 * @param {ApiClient} apiClient - API client instance
 * @param {Array} resourceIds - Array of resource IDs to clean up
 */
async function cleanupTestData(apiClient, resourceIds = []) {
  // Implementation depends on your cleanup strategy
  // This could involve deleting created resources, clearing database, etc.
  console.log('Cleaning up test data...', resourceIds);
}

module.exports = {
  generateUniqueEmail,
  generatePassword,
  generateUserData,
  generateSimpleData,
  waitForCondition,
  retryOperation,
  assertResponse,
  assertErrorResponse,
  cleanupTestData
};
