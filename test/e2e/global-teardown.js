// global-teardown.js

async function globalTeardown() {
  console.log('🧹 Starting global teardown for E2E tests...');
  
  // Add any cleanup logic here if needed
  // For example, cleaning up test data, stopping services, etc.
  
  console.log('✅ Global teardown completed');
}

module.exports = globalTeardown;
