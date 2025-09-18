# Go CRUD API

A modern, production-ready CRUD (Create, Read, Update, Delete) REST API built with Go, featuring:
- [Gin](https://github.com/gin-gonic/gin) web framework
- [GORM](https://gorm.io) ORM with PostgreSQL
- [Goose](https://github.com/pressly/goose) database migrations
- [Swagger](https://swagger.io) API documentation
- JWT authentication with bcrypt password hashing
- Standardized JSON response format
- Docker containerization

This project serves as a robust foundation for backend applications requiring authentication and CRUD operations.

## TODO
- Finalise stack & trace dashboards
- Add unit tests
- Add Playwright test suite (or other)
- Ensure goroutines are used where appropriate
- Probably remove GORM, use sqlc or raw sql instead

## Features

- 🚀 **RESTful CRUD API**: Clean, intuitive endpoint design for complete resource management
- 📖 **Auto-generated Documentation**: Swagger UI with interactive testing
- 🔐 **JWT Authentication**: Secure user registration and login
- 🐳 **Docker Ready**: Containerized deployment with Docker Compose
- 📈 **Complete Observability**: Metrics, traces, and logs with Grafana dashboards
- 🔄 **Database Migrations**: Version-controlled schema changes

## Prerequisites

- Go 1.24 or later
- PostgreSQL 13+ installed and running
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
git clone https://github.com/Verano-20/go-crud.git
cd go-crud
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
git clone https://github.com/Verano-20/go-crud.git
cd go-crud
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
go-crud/
├── cmd/
│   ├── api-server/
│   │   └── main.go          # Application entry point
│   └── migrate/
│       ├── main.go          # Migration command
│       └── migrations/      # SQL migration files
├── internal/
│   ├── config/             # Configuration management
│   │   └── config.go
│   ├── controller/         # HTTP request handlers
│   │   ├── auth.go         # Authentication endpoints
│   │   ├── health.go       # Health check endpoint
│   │   └── simple.go       # Simple resource CRUD
│   ├── initializer/        # Application initialization
│   │   └── database.go
│   ├── middleware/         # HTTP middleware
│   │   └── auth.go         # JWT authentication middleware
│   ├── model/             # Database models and DTOs
│   │   ├── simple.go      # Example resource
│   │   └── user.go        # User model for authentication
│   ├── repository/        # Data access layer
│   │   ├── simple.go      # Example repository
│   │   └── user.go        # User repository
│   ├── response/          # Standardized response types
│   │   └── response.go
│   └── router/            # Route definitions
│       └── router.go
├── docs/                  # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── scripts/
│   └── run-local.sh      # Script to run locally
├── .env.local           # Local environment configuration (create manually)
├── .env.docker         # Docker environment configuration (create manually)
├── Dockerfile          # Docker build instructions
├── docker-compose.yml  # Docker Compose configuration
├── go.mod              # Go module file
├── Go-CRUD.postman_collection.json # Postman collection for API testing
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
- **Collection File**: `Go-CRUD.postman_collection.json`
- Import this collection into Postman to test all API endpoints
- The collection includes pre-configured requests for all CRUD operations
- Set the `baseUrl` variable to `http://localhost:8080` in your Postman environment

**To use the Postman collection:**
1. Open Postman
2. Click "Import" and select the `Go-CRUD.postman_collection.json` file
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

- **Config**: Configuration management
- **Initializer**: Application startup and database connection
- **Router**: Route definitions and middleware
- **Controllers**: Handle HTTP requests and responses
- **Middleware**: HTTP middleware for authentication, logging, etc.
- **Repositories**: Data access layer
- **Models**: Database entities and DTOs
- **Responses**: Standardized API response formats

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
| `DB_HOST` | Database host | Yes | localhost |
| `DB_USER` | Database user | Yes | postgres |
| `DB_PASSWORD` | Database password | Yes | postgres |
| `DB_NAME` | Database name | Yes | go_crud |
| `DB_PORT` | Database port | No | 5432 |
| `JWT_SECRET` | JWT signing secret | Yes | - |

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