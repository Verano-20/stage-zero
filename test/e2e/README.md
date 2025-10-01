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
e2e/
├── tests/                  # Test specifications
│   ├── health.spec.js      # Health check tests
│   ├── auth.spec.js        # Authentication tests
│   └── simple-crud.spec.js # CRUD operation tests
├── utils/                  # Test utilities
│   ├── api-client.js       # API interaction wrapper
│   └── test-helpers.js     # Common test functions
├── fixtures/               # Test data and constants
│   └── test-data.js        # Test fixtures and expected responses
├── global-setup.js         # Global test setup
└── global-teardown.js      # Global test cleanup
```

## Test Categories

### Health Check Tests (`health.spec.js`)
- API availability
- Response time validation
- Content type verification

### Authentication Tests (`auth.spec.js`)
- User registration (signup)
- User login
- JWT token validation
- Input validation and error handling
- Duplicate registration prevention

### CRUD Tests (`simple-crud.spec.js`)
- Create resources with authentication
- Read operations (get all, get by ID)
- Update operations
- Delete operations
- Complete CRUD workflow
- Authorization requirements
- Error handling (404, 401, 400)

## Configuration

### Environment Variables

The tests use the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `BASE_URL` | `http://localhost:8080` | API base URL |
| `CI` | `false` | CI environment flag |

### Docker Configuration

Tests use `docker-compose.test.yml` which provides:
- Isolated test database (`go_crud_test`)
- Test-specific environment variables
- Separate ports to avoid conflicts
- Automatic cleanup

### Playwright Configuration

Key settings in `playwright.config.js`:
- **Browsers**: Chromium, Firefox, WebKit
- **Retries**: 2 on CI, 0 locally
- **Timeout**: 30s per test, 5s for assertions
- **Artifacts**: Screenshots and videos on failure
- **Reports**: HTML, JSON, and JUnit formats

## API Client

The `ApiClient` class (`utils/api-client.js`) provides:
- Automatic JWT token management
- Standardized request/response handling
- Built-in authentication helpers
- Type-safe API method wrappers

Example usage:
```javascript
const apiClient = new ApiClient(request);

// Login and set token automatically
await apiClient.login(credentials, true);

// Create a resource (uses stored token)
const response = await apiClient.createSimple({ name: 'Test Resource' });
```

## Test Helpers

Common utilities in `test-helpers.js`:
- `generateUserData()` - Create unique test users
- `generateSimpleData()` - Create test resources
- `assertResponse()` - Validate API responses
- `assertErrorResponse()` - Validate error responses
- `waitForCondition()` - Polling with timeout
- `retryOperation()` - Retry with exponential backoff

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

### Manual API Testing
The test environment runs on `http://localhost:8080` and includes:
- Health check: `GET /health`
- Swagger UI: `GET /swagger/index.html`
- All API endpoints as documented

## CI/CD Integration

The test suite is designed for CI/CD with:
- **Exit codes**: Proper exit codes for pipeline integration
- **Artifacts**: Test reports in multiple formats
- **Parallel execution**: Configurable worker count
- **Retry logic**: Automatic retry on CI environments
- **Docker isolation**: No dependency on host environment

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
3. **Timeout errors**: Increase timeout in `playwright.config.js`
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
4. Add appropriate test data to `fixtures/test-data.js`
5. Update this README if adding new test categories
