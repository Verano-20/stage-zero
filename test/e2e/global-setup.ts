import { chromium, FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig): Promise<void> {
  console.log('🚀 Starting global setup for E2E tests...');
  
  // Wait for the API to be ready
  const browser = await chromium.launch();
  const page = await browser.newPage();
  
  const maxRetries = 30;
  let retries = 0;
  
  while (retries < maxRetries) {
    try {
      const response = await page.request.get('http://localhost:8080/health');
      if (response.ok()) {
        console.log('✅ API is ready!');
        break;
      }
    } catch (error) {
      console.log(`⏳ Waiting for API... (${retries + 1}/${maxRetries})`);
    }
    
    retries++;
    if (retries === maxRetries) {
      throw new Error('❌ API failed to start within timeout period');
    }
    
    await new Promise(resolve => setTimeout(resolve, 2000));
  }
  
  await browser.close();
  console.log('✅ Global setup completed');
}

export default globalSetup;
