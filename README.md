# Stage Zero

A production-ready Go REST API showcasing enterprise-grade backend engineering practices. Built with clean architecture, comprehensive observability, and automated deployment pipelines.

## Summary

**Stage Zero** is a modern, scalable backend foundation that demonstrates senior-level software engineering skills. It provides a complete CRUD API with JWT authentication, full observability stack, and automated deployment infrastructure.

## Table of Contents

- [Summary](#summary)
- [TODO](#todo)
- [Key Value Propositions](#key-value-propositions)
- [Features](#features)
  - [🔐 Authentication & Authorization](#-authentication--authorization)
  - [📝 CRUD Operations](#-crud-operations)
  - [📊 Complete Observability Stack](#-complete-observability-stack)
  - [🧪 Comprehensive Testing](#-comprehensive-testing)
  - [🚀 Production Deployment](#-production-deployment)
  - [🗄️ Database Management](#️-database-management)
  - [📖 Documentation](#-documentation)
- [Architecture](#architecture)
  - [Clean Architecture Implementation](#clean-architecture-implementation)
  - [Dependency Injection Container](#dependency-injection-container)
  - [Technology Stack](#technology-stack)
  - [Key Design Patterns](#key-design-patterns)
- [Project Structure](#project-structure)
  - [Layer Responsibilities](#layer-responsibilities)
- [Local Setup](#local-setup)
  - [Prerequisites](#prerequisites)
  - [Quick Start](#quick-start)
  - [First Steps](#first-steps)
- [API docs & Postman](#api-docs--postman)
  - [Swagger Documentation](#swagger-documentation)
  - [Authentication Flow](#authentication-flow)
  - [Postman Collection](#postman-collection)
- [Testing](#testing)
  - [Unit Testing](#unit-testing)
  - [E2E Testing](#e2e-testing)
  - [CI Integration](#ci-integration)
- [Observability](#observability)
  - [Monitoring Stack](#monitoring-stack)
  - [Pre-configured Dashboards](#pre-configured-dashboards)
  - [Custom Metrics](#custom-metrics)
  - [Logging Strategy](#logging-strategy)
- [CI/CD](#cicd)
  - [GitHub Actions Workflows](#github-actions-workflows)
  - [Infrastructure as Code](#infrastructure-as-code)
  - [Deployment Process](#deployment-process)
  - [Environment Management](#environment-management)
- [Security](#security)
  - [Authentication & Authorization](#authentication--authorization-1)
  - [Input Validation](#input-validation)
- [Development Guidelines](#development-guidelines)
- [License](#license)
- [Support](#support)

## TODO
- Finalise README and GUIDELINES
- Finalise stack & trace dashboards

### Key Value Propositions

- **🔒 Production Ready**: Build fully secured enterprise-grade applications from day 1.
- **🚀 Automated Deployments**: Focus on building value, not configuring servers.
- **📊 Complete Observability**: Full visibility at every level of your application, always.
- **🧪 Fully Tested**: Validate every new flow with the simple yet comprehensive test architecture.
- **📖 AI Agent Friendly**: Clear guidelines and project structure to keep AI Agents on the right track.
- **🏗️ Clean Architecture**: Avoid spaghetti code as your project grows with modern architecture patterns.

## Features

### 🔐 Authentication & Authorization
- JWT-based authentication with secure token validation
- Password hashing using bcrypt
- Middleware-based route protection
- User registration and login endpoints

### 📝 CRUD Operations
- RESTful API design with intuitive endpoint structure
- Complete resource management (Create, Read, Update, Delete)
- Input validation and error handling
- Standardized response formats

### 📊 Complete Observability Stack
- **Metrics**: Prometheus integration with custom application metrics
- **Tracing**: OpenTelemetry distributed tracing for every request
- **Logging**: Structured JSON logging with Zap
- **Dashboards**: Pre-configured Grafana dashboards for monitoring

### 🧪 Comprehensive Testing
- **Unit Tests**: Go tests with dependency injection mocks
- **E2E Tests**: Playwright end-to-end test suite
- **CI Integration**: Automated testing on pull requests

### 🚀 Production Deployment
- **Infrastructure as Code**: Terraform for DigitalOcean provisioning
- **Container Orchestration**: Docker Compose for local and production environments
- **CI/CD Pipeline**: GitHub Actions for automated testing and deployment
- **Bootstrap Friendly**: Start with a single VM for minimal infrastructure costs

### 🗄️ Database Management
- **Migrations**: Version-controlled schema changes with Goose
- **Connection Pooling**: Optimized database connections
- **ORM Integration**: GORM for type-safe database operations

### 📖 Documentation
- **Swagger UI**: Auto-generated interactive API documentation
- **Postman Collection**: Ready-to-use API testing collection
- **Code Comments**: Comprehensive inline documentation

## Architecture

### Clean Architecture Implementation

Stage Zero follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│              HTTP Layer                 │
│    (Controllers, Middleware, Router)    │
├─────────────────────────────────────────┤
│            Business Layer               │
│       (Services, Domain Logic)          │
├─────────────────────────────────────────┤
│             Data Layer                  │
│        (Repositories, Models)           │
├─────────────────────────────────────────┤
│          Infrastructure Layer           │
│    (Database, IaC, External Services)   │
└─────────────────────────────────────────┘
```

### Dependency Injection Container

The application uses a **Container** pattern with interfaces for dependency injection, ensuring loose coupling and testability:

```go
type Container struct {
    // Repositories (Data Layer) (Interfaces)
    UserRepository   repository.UserRepository
    SimpleRepository repository.SimpleRepository
    
    // Services (Business Layer) (Interfaces)
    UserService   service.UserService
    AuthService   service.AuthService
    SimpleService service.SimpleService
    
    // Controllers (HTTP Layer)
    AuthController   *controller.AuthController
    SimpleController *controller.SimpleController
}
```

### Technology Stack

- **Language**: Go 1.24
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT with golang-jwt/jwt
- **Observability**: OpenTelemetry, Prometheus, Grafana
- **Testing**: Go testing + Playwright
- **Deployment**: Docker, Terraform, GitHub Actions

## Project Structure

```
stage-zero/
├── cmd/                  # Application entry points
│   ├── api-server/           # Main API server
│   └── migrate/              # Database migration tool
├── internal/             # Private application code
│   ├── config/               # Configuration management
│   ├── container/            # Dependency injection container
│   ├── controller/           # HTTP request handlers
│   ├── database/             # Database initialization
│   ├── err/                  # Custom error types
│   ├── logger/               # Structured logging
│   ├── middleware/           # HTTP middleware (auth, logging, metrics)
│   ├── model/                # Database models and DTOs
│   ├── repository/           # Data access layer
│   ├── response/             # Standardized response types
│   ├── router/               # Route definitions
│   ├── service/              # Business logic layer
│   ├── telemetry/            # OpenTelemetry configuration
|   └── utils/                # Utility functions (Binding error handler)
├── test/                 # Test suites
│   ├── e2e/                  # End-to-end tests (Playwright)
│   ├── middleware/           # Unit tests
│   ├── mocks/                # Test mocks
│   ├── service/              # Unit tests
│   └── testutils/            # Test utility functions
├── terraform/            # Infrastructure as Code
├── scripts/              # Deployment, test, and utility scripts
├── grafana/              # Monitoring dashboards
└── docs/                 # Auto-generated API documentation
```

### Layer Responsibilities

- **Controllers**: Handle HTTP requests, validate input, call services
- **Services**: Implement business logic, orchestrate operations
- **Repositories**: Abstract data access, handle database operations
- **Models**: Define data structures and validation rules
- **Middleware**: Cross-cutting concerns (auth, logging, metrics)

## Local Setup

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Node.js 18+ (for E2E tests)
- Git

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/Verano-20/stage-zero.git
   cd stage-zero
   ```

2. **Create `.env.docker` file:**
   ```env
   SERVICE_NAME=stage-zero-api
   SERVICE_VERSION=1.0.0
   SERVICE_PORT=8080
   ENVIRONMENT=develop
   JWT_SECRET=your-secret-here

   # Application database configuration
   DB_HOST=db
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=stage-zero-db
   DB_PORT=5432

   # PostgreSQL container configuration
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_DB=stage-zero-db

   # Telemetry configuration 
   ENABLE_STDOUT=false
   ENABLE_OTLP=true
   OTLP_ENDPOINT=otel-collector:4317
   OTLP_INSECURE=true
   METRIC_INTERVAL=30s
   ```

2. **Start the application**
   ```bash
   docker-compose up -d
   ```

3. **Verify deployment**
   ```bash
   curl http://localhost:8080/health
   ```

### First Steps

1. **Register a user**
   ```bash
   curl -X POST http://localhost:8080/auth/signup \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"SecurePass123!"}'
   ```

2. **Login and get token**
   ```bash
   curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"SecurePass123!"}'
   ```

3. **Make authenticated requests**
   ```bash
   curl -X GET http://localhost:8080/simples \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## API docs & Postman

### Swagger Documentation

Interactive API documentation is available at:
- **Local**: http://localhost:8080/swagger/index.html

### Authentication Flow

1. **Sign Up**: `POST /auth/signup` with email and password
2. **Login**: `POST /auth/login` to receive JWT token
3. **Authenticate**: Include `Authorization: Bearer <token>` header

### Postman Collection

Import the ready-to-use Postman collection:
- **File**: `Stage-Zero.postman_collection.json`
- **Environment**: Configure base URL.
- **Authentication**: Successful login request sets the 'token' collection variable. This is used as Bearer Token Auth for all requests.

## Testing

All tests are located in the /test/ directory.

### Unit Testing

Run unit tests with mocks:
```bash
# Run all tests
go test ./test/...

# Run with coverage
go test -cover ./test/...

# Run specific package
go test ./test/service
```

### E2E Testing

End-to-end tests using Playwright:
```bash
# Install dependencies
npm install

# Run E2E tests
npm run test:e2e

# Run with UI
npm run test:e2e:ui

# Debug specific test
npm run test:e2e:debug -- --grep "test name"
```

### CI Integration

Tests run automatically on:
- Pull requests to `main` branch
- Pull requests to `deployment` branch
- Manual workflow triggers

## Observability

### Monitoring Stack

Complete observability with modern tools:

- **Grafana**: Visualization and dashboards
- **Prometheus**: Metrics collection and storage
- **Tempo**: Distributed tracing
- **Loki**: Log aggregation
- **OpenTelemetry**: Telemetry collection

### Pre-configured Dashboards

Access Grafana at http://localhost:3000 (admin/admin):

- **Health Dashboard**: Stack health
- **Metrics Dashboard**: Application metrics
- **Logs Dashboard**: Centralized log viewing
- **Tracing Dashboard**: Request tracing and performance

### Custom Metrics

Application-specific metrics:
- Request duration and count
- Authentication success/failure rates
- Database operation metrics
- Custom business metrics

### Logging Strategy

- **Structured Logging**: JSON format with Zap
- **Correlation IDs**: Request tracing across services
- **Log Levels**: Debug, Info, Warn, Error
- **Contextual Logging**: Request-scoped loggers

## CI/CD

### GitHub Actions Workflows

#### Test Pipeline (`run-tests-and-update-docs.yml`)
- Runs on pull requests to `main` and `deployment`
- Executes Go unit tests with race detection
- Runs Playwright E2E tests
- Generates and updates API documentation

#### Deployment Pipeline (`build-and-deploy.yml`)
- Triggers on push to `deployment` branch
- Builds and pushes Docker images to GitHub Container Registry
- Provisions infrastructure with Terraform
- Deploys containers to DigitalOcean

### Infrastructure as Code

Terraform configuration for DigitalOcean:
- **Droplet Management**: Automated server provisioning
- **Environment Variables**: Secure configuration management

### Deployment Process

1. **Build**: Container image creation and registry push
2. **Infrastructure**: Terraform plan and apply - smart handling of existing resources
3. **Deploy**: Container deployment with health checks

### Environment Management

- **Development**: Local Docker Compose setup with .env.docker file
- **Production**: DigitalOcean with automated deployment
- **Secrets**: GitHub Secrets for sensitive deployment configuration

### Deployment Setup

Some configuration is needed to enable the CD pipelines in a fresh project.

1. Create an ssh key pair that will be used to enable the Github Actions runner to access the DigitalOcean droplets.
2. In DigitalOcean, add the public key to your account with the name 'github_actions'.
3. In your Github repository, add the following repository secrets (settings -> secrets and variables -> actions):
  a. DB_PASSWORD (database password)
  b. POSTGRES_PASSWORD (container password)
  c. JWT_SECRET (for application)
  d. DO_TOKEN (DigitalOcean token)
  e. DO_SSH_PRIVATE_KEY (private key from pair generated in step 1)

These secrets will enable the Github workflow to fully automate deployment to a DigitalOcean droplet. The non-sensitive application environment variables are defined directly in main.tf.

## Security

### Authentication & Authorization

- **JWT Tokens**: Secure token-based authentication
- **Password Security**: bcrypt hashing with salt
- **Token Validation**: Comprehensive JWT verification
- **User Verification**: Database-backed user validation
- **Middleware**: Security middleware on all HTTP requests

### Input Validation

- **Request Validation**: Gin binding with validation tags
- **SQL Injection Prevention**: GORM ORM with parameterized queries
- **Comprehensive Responses**: Standard HTTP error response format listing all validation failures

---

## Development Guidelines

For detailed development guidelines, code standards, and contribution instructions, see [GUIDELINES.md](GUIDELINES.md).

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For questions, issues, or contributions:
- **Issues**: [GitHub Issues](https://github.com/Verano-20/stage-zero/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Verano-20/stage-zero/discussions)
- **Documentation**: [Project Wiki](https://github.com/Verano-20/stage-zero/wiki)
