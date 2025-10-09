# Stage Zero

A modern, production-ready CRUD (Create, Read, Update, Delete) REST API built with Go, featuring:
- 🚀 **RESTful CRUD API**: Clean, intuitive endpoint design for complete resource management
- 📖 **Auto-generated Documentation**: Swagger UI with interactive testing
- 🔐 **JWT Authentication**: Secure user registration and login
- 🐳 **Docker Ready**: Containerized deployment with Docker Compose
- 📈 **Complete Observability**: Metrics, traces, and logs with Grafana dashboards
- 🔄 **Database Migrations**: Version-controlled schema changes
- 🏗️ **Clean Architecture**: Dependency injection with service/repository pattern
- 📝 **Structured Logging**: JSON logging with Zap for better observability
- 🧪 **Comprehensive Testing**: Unit tests with mocks and E2E tests with Playwright
- 🔧 **Developer Experience**: Hot reload, test automation, and debugging tools


This project serves as a robust foundation for backend applications requiring authentication and CRUD operations.

## Table of Contents

- [TODO](#todo)
- [Packages and Tools](#packages-and-tools)
- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Setup](#setup)
  - [Option 1: Local Setup](#option-1-local-setup)
  - [Option 2: Docker Setup](#option-2-docker-setup)
- [Database Migrations](#database-migrations)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
  - [Swagger Documentation](#swagger-documentation)
  - [Postman Collection](#postman-collection)
- [API Endpoints](#api-endpoints)
  - [System Health](#system-health)
  - [Authentication](#authentication)
  - [Simple Resource](#simple-resource)
- [Response Format](#response-format)
- [API Usage Examples](#api-usage-examples)
  - [Authentication](#authentication-1)
  - [Simple Resource Management](#simple-resource-management)
- [Observability](#observability)
  - [Architecture](#architecture-1)
  - [Quick Start](#quick-start)
  - [Pre-configured Dashboards](#pre-configured-dashboards)
  - [Available Metrics](#available-metrics)
- [Testing](#testing)
  - [Test Architecture](#test-architecture)
  - [Unit Tests](#unit-tests)
  - [End-to-End Tests](#end-to-end-tests)
  - [Test Configuration](#test-configuration)
  - [Continuous Integration](#continuous-integration)
  - [Test Data Management](#test-data-management)
  - [Debugging Tests](#debugging-tests)
- [Architecture](#architecture)
  - [Architectural Layers](#architectural-layers)
  - [Dependency Injection](#dependency-injection)
  - [Layer Responsibilities](#layer-responsibilities)
  - [Request Flow](#request-flow)
  - [Key Design Patterns](#key-design-patterns)
  - [Configuration Management](#configuration-management)
  - [Error Handling](#error-handling)
- [Development](#development)
  - [Generating Swagger Documentation](#generating-swagger-documentation)
  - [Code Structure Guidelines](#code-structure-guidelines)
- [Security Features](#security-features)
  - [Environment Variables](#environment-variables)
- [🚀 Production Deployment](#-production-deployment)
  - [Deployment Architecture](#deployment-architecture)
  - [Complete E2E Deployment Flow](#complete-e2e-deployment-flow)
  - [Deployment Timeline](#deployment-timeline)
  - [Key Features](#key-features)
  - [Service Stack](#service-stack)
  - [Deployment Commands](#deployment-commands)
  - [Environment Variables](#environment-variables-1)
  - [Monitoring & Debugging](#monitoring--debugging)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## TODO
- Finalise README, include fork instructions
- Finalise stack & trace dashboards
- Verify all CI/CD flows

For CI/CD:
- rename main -> develop
- figure out redeployment when droplet exists

## Packages and Tools
- [Gin](https://github.com/gin-gonic/gin) web framework
- [Swagger](https://swagger.io) API documentation
- [GORM](https://gorm.io) ORM with PostgreSQL
- [Goose](https://github.com/pressly/goose) database migrations
- [Zap](https://github.com/uber-go/zap) structured logging
- [stretchr/testify](https://github.com/stretchr/testify) unit testing and mocking
- [Docker](https://docker.com/) containerization
- [OpenTelemetry](https://opentelemetry.io/) comprehensive observability
- [Grafana](https://grafana.com/) dashboards for metrics, traces, and logs
- [Playwright](https://playwright.dev/) end to end test suite
- [Postman](https://www.postman.com/) collection for manual testing
- JWT authentication with bcrypt password hashing

## Prerequisites

- Go 1.24 or later
- PostgreSQL 13+ installed and running
- Node.js 24+ (for E2E testing)
- Git (optional)
- Docker and Docker Compose (optional)

## Environment Setup

The application uses environment variables for configuration.

**For Docker Compose:** Create a `.env.docker` file in the root directory with:
```bash
# Application variables
DB_HOST=db
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_crud
DB_PORT=5432

# JWT Secret (generate using: openssl rand -base64 32)
JWT_SECRET=your-generated-secret-here

# PostgreSQL container variables
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=go_crud
```

**For local development:** Create a `.env.local` file in the root directory with:
```bash
# Application variables
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_crud
DB_PORT=5432

# JWT Secret (generate using: openssl rand -base64 32)
JWT_SECRET=your-generated-secret-here
```

## Setup

### Option 1: Local Setup

1. Clone the repository:
```bash
git clone https://github.com/Verano-20/stage-zero.git
cd stage-zero
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your PostgreSQL database:
   - Create a new PostgreSQL database named `go_crud`
   - Create a `.env.local` file with your database configuration

4. Run database migrations and start the server:
```bash
./scripts/run-local.sh
```

### Option 2: Docker Setup

1. Clone the repository:
```bash
git clone https://github.com/Verano-20/stage-zero.git
cd stage-zero
```

2. Create a `.env.docker` file with the environment variables shown above

3. Build and run with Docker Compose:
```bash
docker-compose up --build
```

This will:
- Build the Go application
- Start PostgreSQL database
- Run database migrations automatically
- Start the API server
- Create a persistent volume for the database
- Expose the API on port 8080

To stop the services:
```bash
docker-compose down
```

To remove the persistent volume as well:
```bash
docker-compose down -v
```

## Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for database migrations. Migrations are stored in `cmd/migrate/migrations/` and can be run using the migrate command.

**Available migration commands:**
```bash
# Run all pending migrations
go run cmd/migrate/main.go up

# Roll back one migration
go run cmd/migrate/main.go down

# Check migration status
go run cmd/migrate/main.go status

# Reset all migrations
go run cmd/migrate/main.go reset
```

Migration files follow this format:
```sql
-- +goose Up
-- SQL to run when migrating up
CREATE TABLE example (...);

-- +goose Down
-- SQL to run when migrating down
DROP TABLE example;
```

## Project Structure

```
stage-zero/
├── cmd/
│   ├── api-server/
│   │   └── main.go          # Application entry point
│   └── migrate/
│       ├── main.go          # Migration command
│       └── migrations/      # SQL migration files
├── internal/
│   ├── config/             # Configuration management
│   │   └── config.go
│   ├── container/          # Dependency injection container
│   │   └── container.go
│   ├── controller/         # HTTP request handlers
│   │   ├── auth.go         # Authentication endpoints
│   │   ├── health.go       # Health check endpoint
│   │   └── simple.go       # Simple resource CRUD
│   ├── database/           # Database initialization
│   │   └── database.go
│   ├── err/               # Custom error types
│   │   └── err.go
│   ├── logger/            # Structured logging
│   │   └── logger.go
│   ├── middleware/        # HTTP middleware
│   │   ├── auth.go        # JWT authentication middleware
│   │   ├── logging.go     # Request logging middleware
│   │   └── metrics.go     # Metrics middleware
│   ├── model/             # Database models and DTOs
│   │   ├── simple.go      # Example resource
│   │   └── user.go        # User model for authentication
│   ├── repository/        # Data access layer
│   │   ├── simple.go      # Example repository
│   │   └── user.go        # User repository
│   ├── response/          # Standardized response types
│   │   └── response.go
│   ├── router/            # Route definitions
│   │   └── router.go
│   ├── service/           # Business logic layer
│   │   ├── auth.go        # Authentication service
│   │   ├── simple.go      # Simple resource service
│   │   └── user.go        # User service
│   └── telemetry/         # OpenTelemetry configuration
│       ├── metrics.go     # Custom metrics
│       └── telemetry.go   # Telemetry setup
├── test/                  # Test suites
│   ├── e2e/              # End-to-end tests (Playwright)
│   │   ├── tests/        # Test specifications
│   │   ├── utils/        # Test utilities
│   │   ├── fixtures/     # Test data
│   │   └── README.md     # E2E testing guide
│   ├── middleware/       # Middleware unit tests
│   ├── mocks/           # Generated mocks
│   ├── service/         # Service unit tests
│   └── testutils/       # Test utilities
├── grafana/              # Grafana dashboards and config
├── docs/                 # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── scripts/
│   ├── run-local.sh     # Script to run locally
│   └── run-e2e-tests.sh # Script to run E2E tests
├── .env.local           # Local environment configuration (create manually)
├── .env.docker         # Docker environment configuration (create manually)
├── Dockerfile          # Docker build instructions
├── docker-compose.yml  # Docker Compose configuration
├── docker-compose.test.yml # Test environment configuration
├── go.mod              # Go module file
├── package.json        # Node.js dependencies for E2E tests
├── playwright.config.ts # Playwright configuration
├── stage_zero.postman_collection.json # Postman collection for API testing
└── README.md          # This file
```

## API Documentation

The API includes comprehensive documentation and testing tools:

### Swagger Documentation
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- The documentation is automatically generated using [swaggo/swag](https://github.com/swaggo/swag) annotations
- Interactive interface for testing all endpoints
- Complete request/response schemas with examples

### Postman Collection
- **Collection File**: `stage_zero.postman_collection.json`
- Import this collection into Postman to test all API endpoints
- The collection includes pre-configured requests for all CRUD operations
- Set the `baseUrl` variable to `http://localhost:8080` in your Postman environment

**To use the Postman collection:**
1. Open Postman
2. Click "Import" and select the `stage_zero.postman_collection.json` file
3. Create a new environment and set `baseUrl` to `http://localhost:8080`
4. Select the environment and start testing the endpoints

## API Endpoints

### System Health
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Check server health status | No |

### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/signup` | Create a new user account | No |
| POST | `/auth/login` | Authenticate user and get JWT token | No |

### Simple Resource
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/simple` | Create a new simple resource | Yes |
| GET | `/simple` | Get all simple resources | Yes |
| GET | `/simple/:id` | Get a simple resource by ID | Yes |
| PUT | `/simple/:id` | Update a simple resource | Yes |
| DELETE | `/simple/:id` | Delete a simple resource | Yes |

## Response Format

All API responses follow a consistent format:

**Success Response:**
```json
{
  "message": "Operation successful",
  "data": {
    "id": 1,
    "name": "Example"
  }
}
```

**Error Response:**
```json
{
  "error": "Error message description"
}
```

## API Usage Examples

### Authentication

**Sign up a new user:**
```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securePassword123"
  }'
```

**Login and get JWT token:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securePassword123"
  }'
```

**Response:**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Simple Resource Management

**Create a simple resource (requires authentication):**
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X POST http://localhost:8080/simple \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "My Simple Resource"}'
```

**Get all resources (requires authentication):**
```bash
curl -X GET http://localhost:8080/simple \
  -H "Authorization: Bearer $TOKEN"
```

**Get a specific resource (requires authentication):**
```bash
curl -X GET http://localhost:8080/simple/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Update a resource (requires authentication):**
```bash
curl -X PUT http://localhost:8080/simple/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Updated Resource Name"}'
```

**Delete a resource (requires authentication):**
```bash
curl -X DELETE http://localhost:8080/simple/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Check server health:**
```bash
curl http://localhost:8080/health
```

## Observability

This application includes enterprise-grade observability with **metrics**, **traces**, and **logs** using OpenTelemetry and Grafana.

### Architecture
```
Go App ────┬─→ OpenTelemetry Collector ├→ Prometheus → Grafana
           │                          └→ Tempo ────┘
           └─→ JSON Logs → Promtail → Loki ──────────┘
```

### Quick Start

1. **Add telemetry configuration** to your `.env.docker`:
```bash
# Telemetry Configuration
ENABLE_STDOUT=false
ENABLE_OTLP=true
OTLP_ENDPOINT=http://otel-collector:4318
OTLP_INSECURE=true
METRIC_INTERVAL=30s
```

2. **Start the complete observability stack**:
```bash
docker-compose up -d
```

3. **Access Grafana dashboards**: http://localhost:3000 (admin/admin)

### Pre-configured Dashboards

- **Application Overview**: HTTP metrics, database performance, auth statistics
- **Distributed Tracing**: Request flows, service maps, error traces  
- **Structured Logs**: Live logs, error filtering, trace correlation

### Available Metrics

The application automatically exports:
- **HTTP**: Request rates, duration, active requests
- **Database**: Connection pools, query performance
- **Authentication**: Login attempts, failures
- **Business**: User counts, entity counts

## Testing

Stage Zero includes comprehensive testing at multiple levels to ensure reliability and maintainability.

### Test Architecture

The project implements a multi-layered testing strategy:
- **Unit Tests**: Test individual components in isolation using mocks
- **End-to-End Tests**: Test complete user workflows via HTTP API

### Unit Tests

Unit tests are written in Go using the standard `testing` package with `testify` for assertions and mocks.

**Running unit tests:**
```bash
# Run all unit tests
go test ./test/...

# Run tests with coverage
go test -cover ./test/...

# Run tests with verbose output
go test -v ./test/...
```

**Test Structure:**
- `test/service/`: Service layer unit tests
- `test/middleware/`: Middleware unit tests  
- `test/mocks/`: Generated mocks for interfaces
- `test/testutils/`: Common test utilities

**Example unit test:**
```go
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := mocks.NewUserRepository(t)
    service := service.NewUserService(mockRepo)
    
    userData := &model.User{Email: "test@example.com"}
    mockRepo.On("Create", userData).Return(nil)
    
    err := service.CreateUser(userData)
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### End-to-End Tests

E2E tests use Playwright to test the complete API functionality through HTTP requests.

**Prerequisites:**
- Node.js 24+
- Docker and Docker Compose

**Setup and run E2E tests:**
```bash
# Install dependencies
npm run test:setup

# Run all E2E tests
npm run test:e2e

# Run tests with browser UI visible
npm run test:e2e:headed

# Run tests in debug mode
npm run test:e2e:debug

# Run specific test suite
npm run test:e2e -- --grep "authentication"

# View test report
npm run test:e2e:report
```

**Test Categories:**
- **Health Check Tests**: API availability and response validation
- **Authentication Tests**: User registration, login, JWT validation
- **CRUD Tests**: Complete resource lifecycle testing
- **Security Tests**: Authentication requirements and error handling

**E2E Test Structure:**
```
test/e2e/
├── tests/                 # Test specifications
│   ├── health.spec.ts     # Health endpoint tests
│   ├── auth.spec.ts       # Authentication flow tests
│   └── simple-crud.spec.ts # CRUD operation tests
├── utils/                 # Test utilities
│   ├── api-client.ts      # API interaction wrapper
│   └── test-helpers.ts    # Common test functions
├── fixtures/              # Test data and constants
└── README.md             # Detailed E2E testing guide
```

**Example E2E test:**
```typescript
test('should create and retrieve simple resource', async ({ request }) => {
  const apiClient = new ApiClient(request);
  
  // Authenticate
  const userData = generateUserData();
  await apiClient.signUp(userData);
  await apiClient.login(userData, true);
  
  // Create resource
  const resourceData = generateSimpleData();
  const createResponse = await apiClient.createSimple(resourceData);
  const body = await assertResponse(createResponse, 201);
  
  // Verify creation
  const getResponse = await apiClient.getSimpleById(body.data.id);
  const retrievedBody = await assertResponse(getResponse, 200);
  expect(retrievedBody.data.name).toBe(resourceData.name);
});
```

### Test Configuration

**Environment Variables for Testing:**
| Variable | Default | Description |
|----------|---------|-------------|
| `BASE_URL` | `http://localhost:8080` | API base URL for E2E tests |
| `CI` | `false` | CI environment flag |

**Docker Test Environment:**
- Uses `docker-compose.test.yml` for isolated test database
- Separate ports to avoid conflicts with development environment
- Automatic cleanup between test runs

### Continuous Integration

The test suite is designed for CI/CD integration:

```yaml
# Example GitHub Actions workflow
- name: Run Unit Tests
  run: go test -cover ./...

- name: Run E2E Tests  
  run: npm run test:e2e

- name: Upload Test Results
  uses: actions/upload-artifact@v3
  if: always()
  with:
    name: test-results
    path: |
      test-results/
      playwright-report/
```

### Test Data Management

- **Fixtures**: Predefined test data in `test/e2e/fixtures/test-data.ts`
- **Generators**: Dynamic test data generation functions
- **Cleanup**: Automatic test data cleanup between runs
- **Isolation**: Each test runs with fresh data to prevent interference

### Debugging Tests

**Unit Tests:**
```bash
# Run specific test with verbose output
go test -v -run TestSpecificFunction ./internal/service

# Run tests with race detection
go test -race ./...
```

**E2E Tests:**
```bash
# Debug specific test
npm run test:e2e:debug -- --grep "failing test name"

# Run with Playwright UI for interactive debugging
npm run test:e2e:ui

# Check API logs during tests
docker-compose -f docker-compose.test.yml logs app-test
```

For detailed E2E testing information, see [test/e2e/README.md](test/e2e/README.md).

## Architecture

Stage Zero follows **Clean Architecture** principles with clear separation of concerns and dependency inversion.

### Architectural Layers

```
┌─────────────────────────────────────────┐
│              HTTP Layer                 │
│  (Controllers, Middleware, Router)      │
├─────────────────────────────────────────┤
│            Business Layer               │
│         (Services, Domain Logic)        │
├─────────────────────────────────────────┤
│             Data Layer                  │
│        (Repositories, Models)           │
├─────────────────────────────────────────┤
│          Infrastructure Layer           │
│    (Database, External Services)        │
└─────────────────────────────────────────┘
```

### Dependency Injection

The application uses a **Container** pattern for dependency injection:

```go
type Container struct {
    // Repositories (Data Layer)
    UserRepository   repository.UserRepository
    SimpleRepository repository.SimpleRepository
    
    // Services (Business Layer)  
    UserService   service.UserService
    AuthService   service.AuthService
    SimpleService service.SimpleService
    
    // Controllers (HTTP Layer)
    AuthController   *controller.AuthController
    SimpleController *controller.SimpleController
}
```

**Benefits:**
- **Testability**: Easy to inject mocks for unit testing
- **Modularity**: Clear component boundaries and responsibilities
- **Maintainability**: Changes in one layer don't affect others
- **Flexibility**: Easy to swap implementations

### Layer Responsibilities

**Controllers (HTTP Layer):**
- Handle HTTP requests and responses
- Input validation and serialization
- Route request to appropriate service
- Return standardized JSON responses

**Services (Business Layer):**
- Implement business logic and rules
- Coordinate between multiple repositories
- Handle complex operations and workflows
- Validate business constraints

**Repositories (Data Layer):**
- Abstract database operations
- Implement data access patterns
- Handle database-specific logic
- Provide clean interface to services

**Models:**
- Define data structures and entities
- Include validation rules and constraints
- Separate DTOs for API contracts

### Request Flow

```
HTTP Request → Middleware → Controller → Service → Repository → Database
                ↓              ↓           ↓          ↓
            Logging,       Validation,  Business   Data Access
            Auth,          Parsing      Logic      Operations
            Metrics
```

### Key Design Patterns

**Repository Pattern:**
```go
type UserRepository interface {
    Create(user *model.User) error
    GetByID(id uint) (*model.User, error)
    GetByEmail(email string) (*model.User, error)
}
```

**Service Pattern:**
```go
type UserService interface {
    CreateUser(userData *model.CreateUserRequest) error
    GetUserByID(id uint) (*model.User, error)
}
```

**Middleware Pattern:**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authentication logic
        c.Next()
    }
}
```

### Configuration Management

Centralized configuration with environment-based overrides:

```go
type Config struct {
    ServiceName    string
    Database       DatabaseConfig
    Telemetry      TelemetryConfig
}
```

### Error Handling

Structured error handling with custom error types:
- **Domain Errors**: Business logic violations
- **Validation Errors**: Input validation failures  
- **Infrastructure Errors**: Database, network issues
- **HTTP Errors**: Standardized API error responses

## Development

### Generating Swagger Documentation

To regenerate the Swagger documentation after making changes to the API:

```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g ./cmd/api-server/main.go
```

### Code Structure Guidelines

- **Config**: Centralized configuration management with environment overrides
- **Container**: Dependency injection container for managing component lifecycle
- **Database**: Database initialization and connection management
- **Router**: Route definitions, middleware registration, and HTTP setup
- **Controllers**: Handle HTTP requests/responses, input validation, and routing
- **Services**: Business logic implementation, workflow coordination
- **Repositories**: Data access abstraction and database operations
- **Models**: Database entities, DTOs, and validation rules
- **Middleware**: Cross-cutting concerns (auth, logging, metrics, CORS)
- **Logger**: Structured logging with Zap for better observability
- **Telemetry**: OpenTelemetry metrics, traces, and monitoring setup
- **Responses**: Standardized API response formats and error handling

## Security Features

- **Password Hashing**: Uses bcrypt for secure password storage
- **JWT Authentication**: Stateless authentication with configurable secret
- **Token Validation**: Comprehensive JWT validation including expiration and signing method
- **Bearer Token Format**: Standard Authorization header format support
- **Input Validation**: Request validation and sanitization
- **Error Handling**: Consistent error responses without sensitive information
- **SQL Injection Prevention**: GORM provides built-in protection

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `SERVICE_NAME` | Service name for telemetry | No | stage-zero-api |
| `SERVICE_VERSION` | Service version | No | 1.0.0 |
| `SERVICE_PORT` | HTTP server port | No | 8080 |
| `ENVIRONMENT` | Environment (develop/production) | No | develop |
| `DB_HOST` | Database host | Yes | localhost |
| `DB_USER` | Database user | Yes | postgres |
| `DB_PASSWORD` | Database password | Yes | postgres |
| `DB_NAME` | Database name | Yes | go_crud |
| `DB_PORT` | Database port | No | 5432 |
| `JWT_SECRET` | JWT signing secret | Yes | - |
| `ENABLE_STDOUT` | Enable stdout telemetry | No | true |
| `ENABLE_OTLP` | Enable OTLP telemetry | No | true |
| `OTLP_ENDPOINT` | OTLP collector endpoint | No | localhost:4317 |
| `OTLP_INSECURE` | Use insecure OTLP connection | No | true |
| `METRIC_INTERVAL` | Metrics collection interval | No | 30s |

## 🚀 Production Deployment

Stage Zero features a **production-grade, zero-downtime deployment system** that automatically deploys to DigitalOcean when you push to the `deployment` branch.

### **Deployment Architecture**

- **Infrastructure**: DigitalOcean Droplet (Ubuntu 25.04, 1GB RAM)
- **Containerization**: Docker Compose with multi-service stack
- **CI/CD**: GitHub Actions with Terraform-native infrastructure management
- **Monitoring**: Grafana, Prometheus, OpenTelemetry, Loki, Tempo
- **Zero Downtime**: Rolling container updates without infrastructure changes

### **Streamlined 3-Step Deployment Process**

#### **Step 1: Build & Push Container**
```bash
# Trigger: Push to deployment branch
1. Change Detection: Monitors Dockerfile, Go source, go.mod/go.sum
2. Docker Build: Multi-platform build (linux/amd64, linux/arm64)
3. Registry Push: ghcr.io/verano-20/stage-zero:deployment
4. Layer Caching: Optimized builds with GitHub Actions cache
```

#### **Step 2: Setup Infrastructure**
```bash
# Infrastructure job handles:
1. Terraform Planning: Detects existing vs new infrastructure
2. Droplet Creation: Creates droplet only if none exists
3. Infrastructure Setup: user-data.sh runs on droplet
   - System updates & Docker installation
   - Configuration file downloads
   - Environment setup
   - NO container operations
4. SSH Readiness: Waits for infrastructure to be accessible
```

#### **Step 3: Deploy Application**
```bash
# Application job handles all Docker operations:
1. Container Deployment: deploy-containers.sh runs via SSH
2. Registry Login: Authenticate with GitHub Container Registry
3. Image Pull: Download latest container images
4. Smart Deployment: Detects initial vs update deployment
   - Initial: Start all services
   - Update: Restart only app & migrate containers
5. Health Verification: Comprehensive service status checks
6. Migration Validation: Ensures database migrations completed successfully
```

#### **Final Health Checks**
```bash
# Comprehensive production health checks:
✅ Application: http://droplet-ip:8080/health
✅ Grafana: http://droplet-ip:3000/api/health
✅ Prometheus: http://droplet-ip:9090/-/healthy
✅ All Services: Docker Compose service status validation
```

### **Deployment Timeline**

```
Push to deployment branch
    ↓
Step 1: Build & Push Container (2-3 min)
    ↓
Step 2: Setup Infrastructure (1-3 min)
    ↓
Step 3: Deploy Application (1-2 min)
    ↓
Final Health Checks (30 sec)
    ↓
Deployment Complete! 🎉
```

**Total Time**: ~4-8 minutes (depending on infrastructure state)

### **Key Features**

- **🔄 Zero Downtime**: Only application containers restart, database persists
- **📊 Data Persistence**: Database and volumes survive deployments
- **🛡️ No Duplicates**: Terraform import prevents duplicate droplets
- **⚡ Fast Updates**: ~1-2 minutes for container updates vs 4+ minutes for full deployment
- **🔙 Rollback Ready**: Easy deployment of previous image tags
- **📈 Full Observability**: Complete monitoring stack with Grafana dashboards
- **🎯 Smart Deployment**: Automatically detects initial vs update deployments
- **🔍 Comprehensive Health Checks**: Validates all services including migrate container
- **🏗️ Separation of Concerns**: Clean infrastructure vs application deployment jobs

### **Service Stack**

| Service | Port | Purpose |
|---------|------|---------|
| **Application** | 8080 | Main REST API |
| **PostgreSQL** | 5432 | Database |
| **Grafana** | 3000 | Dashboards (admin/admin) |
| **Prometheus** | 9090 | Metrics collection |
| **OpenTelemetry** | 4317 | Telemetry collection |
| **Tempo** | - | Distributed tracing |
| **Loki** | - | Log aggregation |
| **Promtail** | - | Log collection |

### **Deployment Commands**

```bash
# Manual deployment (if needed)
git checkout deployment
git merge main
git push origin deployment

# Check deployment status
gh run list --workflow="build-and-deploy.yml"

# View deployment logs
gh run view <run-id> --log
```

### **Environment Variables**

The deployment system uses these GitHub Secrets:
- `DO_TOKEN`: DigitalOcean API token
- `GITHUB_TOKEN`: GitHub Container Registry access
- `JWT_SECRET`: Application JWT signing secret
- `POSTGRES_PASSWORD`: Database password
- `DB_PASSWORD`: Application database password

### **Monitoring & Debugging**

- **Application Health**: `http://droplet-ip:8080/health`
- **Grafana Dashboards**: `http://droplet-ip:3000` (admin/admin)
- **Prometheus Metrics**: `http://droplet-ip:9090`
- **Deployment Logs**: GitHub Actions workflow logs
- **Service Logs**: `docker-compose logs <service>` on droplet

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).

## Support

For questions, issues, or contributions, please:
- Open an issue on GitHub
- Check the Swagger documentation at `/swagger/index.html`
- Review the Postman collection for usage examples 
