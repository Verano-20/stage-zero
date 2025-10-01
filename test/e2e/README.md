# E2E Tests for Stage Zero API

This directory contains end-to-end tests for the Stage Zero Go REST API using Playwright.

## Overview

The E2E test suite covers:
- **Health Check**: API availability and response time
- **Authentication**: User signup, login, and JWT token validation
- **CRUD Operations**: Complete lifecycle testing for Simple resources
- **Security**: Authentication requirements and error handling

## Quick Start

### Prerequisites

- Node.js 24+ 
- Docker and Docker Compose
- Go 1.24+ (for the API)

### Setup

1. **Install dependencies:**
   ```bash
   npm run test:setup
   ```

2. **Run all tests:**
   ```bash
   npm run test:e2e
   ```

### Test Commands

| Command | Description |
|---------|-------------|
| `npm run test:e2e` | Run all E2E tests (headless) |
| `npm run test:e2e:headed` | Run tests with browser UI visible |
| `npm run test:e2e:debug` | Run tests in debug mode |
| `npm run test:e2e:ui` | Run tests with Playwright UI |
| `npm run test:e2e:report` | Show test report |
| `npm run test:e2e:raw` | Run Playwright directly (no Docker setup) |

### Advanced Usage

**Run specific tests:**
```bash
./scripts/run-e2e-tests.sh --grep "authentication"
```

**Run tests for specific browser:**
```bash
./scripts/run-e2e-tests.sh --project chromium
```

**Debug specific test:**
```bash
./scripts/run-e2e-tests.sh --debug --grep "should successfully create"
```

## Test Structure

```
test/e2e/
├── tests/                     # Test specifications
│   ├── health.spec.ts         # Health check tests
│   ├── auth.spec.ts           # Authentication tests
│   └── simple-crud.spec.ts    # CRUD operation tests
├── utils/                     # Test utilities
│   ├── api-client.ts          # API interaction wrapper
│   └── test-helpers.ts        # Common test functions
├── fixtures/                  # Test data and constants
│   └── test-data.ts           # Test fixtures and expected responses
├── global-setup.ts            # Global test setup
├── global-teardown.ts         # Global test cleanup
└── README.md                  # This file
```

## API Client

The `ApiClient` class provides methods for all API endpoints:

```typescript
const apiClient = new ApiClient(request);

// Authentication
await apiClient.signUp({ email: 'user@example.com', password: 'password' });
await apiClient.login({ email: 'user@example.com', password: 'password' });

// CRUD operations
await apiClient.createSimple({ name: 'Test Resource' });
await apiClient.getAllSimples();
await apiClient.getSimpleById(1);
await apiClient.updateSimple(1, { name: 'Updated Resource' });
await apiClient.deleteSimple(1);
```

## Test Helpers

Common utilities for test setup and assertions:

```typescript
// Generate test data
const userData = generateUserData();
const resourceData = generateSimpleData();

// Response assertions
const body = await assertResponse(response, 200);
const errorBody = await assertErrorResponse(response, 401);

// Utility functions
await waitForCondition(() => condition, 5000);
await retryOperation(() => operation(), 3);
```

## Test Categories

### Health Check Tests (`health.spec.ts`)
- API availability and response time
- Content type verification

### Authentication Tests (`auth.spec.ts`)
- User registration and login
- JWT token validation
- Input validation and error handling

### CRUD Tests (`simple-crud.spec.ts`)
- Create, read, update, delete operations
- Authentication requirements
- Error handling (404, 401, 400)
- Complete CRUD workflow

## Configuration

### Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `BASE_URL` | `http://localhost:8080` | API base URL |
| `CI` | `false` | CI environment flag |

### Docker Configuration
Tests use `docker-compose.test.yml` which provides:
- Isolated test database
- Test-specific environment variables
- Separate ports to avoid conflicts

## Development

### Adding New Tests
1. Create `.spec.ts` files in `test/e2e/tests/`
2. Import utilities from `utils/`
3. Use fixtures from `fixtures/test-data.ts`
4. Follow existing patterns for consistency

### Adding New API Methods
1. Add method to `ApiClient` class
2. Update test helpers if needed
3. Add to test fixtures

## Debugging

### View Test Reports
```bash
npm run test:e2e:report
```

### Debug Failed Tests
```bash
npm run test:e2e:debug --grep "failing test name"
```

### Check API Logs
```bash
docker-compose -f docker-compose.test.yml logs app-test
```

## CI/CD Integration

### GitHub Actions Example
```yaml
- name: Run E2E Tests
  run: npm run test:e2e
  
- name: Upload Test Results
  uses: actions/upload-artifact@v3
  if: always()
  with:
    name: playwright-report
    path: playwright-report/
```

## Troubleshooting

### Common Issues
1. **Port conflicts**: Ensure ports 8080 and 5433 are available
2. **Docker issues**: Check Docker is running and has sufficient resources
3. **Timeout errors**: Increase timeout in `playwright.config.ts`
4. **Database issues**: Check PostgreSQL container logs

### Clean Reset
```bash
# Stop all containers and remove volumes
docker-compose -f docker-compose.test.yml down -v --remove-orphans

# Remove node modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

## Contributing

When adding new tests:
1. Follow existing naming conventions
2. Use the `ApiClient` for API interactions
3. Include both positive and negative test cases
4. Add appropriate test data to `fixtures/test-data.ts`
5. Update this README if adding new test categories