# Go CRUD API

A simple CRUD (Create, Read, Update, Delete) REST API built with Go, using:
- [Gin](https://github.com/gin-gonic/gin) for the web framework
- [GORM](https://gorm.io) as the ORM
- [Goose](https://github.com/pressly/goose) for database migrations
- [Swagger](https://swagger.io) for API documentation
- PostgreSQL as the database

This project is intended to be used as a start point for backend applications.

## Prerequisites

- Go 1.24 or later
- PostgreSQL installed and running
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
│   │   ├── health.go
│   │   └── simple.go
│   ├── initializer/        # Application initialization
│   │   └── database.go
│   ├── model/             # Database models and DTOs
│   │   └── simple.go
│   ├── repository/        # Data access layer
│   │   └── simple.go
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
- The documentation is automatically generated using [swaggo/swag](https://github.com/swaggo/swag) annotations in the controller files

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

### Health Check
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Check server health status |

### Simple Resource
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/simple` | Create a new simple resource |
| GET | `/simple` | Get all simple resources |
| GET | `/simple/:id` | Get a simple resource by ID |
| PUT | `/simple/:id` | Update a simple resource |
| DELETE | `/simple/:id` | Delete a simple resource |

### Example API Usage

You can test the API using either the provided Postman collection or curl commands:

**Create a simple resource:**
```bash
curl -X POST http://localhost:8080/simple \
  -H "Content-Type: application/json" \
  -d '{"name": "Example Name"}'
```

**Get all resources:**
```bash
curl http://localhost:8080/simple
```

**Get a specific resource:**
```bash
curl http://localhost:8080/simple/1
```

**Update a resource:**
```bash
curl -X PUT http://localhost:8080/simple/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name"}'
```

**Delete a resource:**
```bash
curl -X DELETE http://localhost:8080/simple/1
```

## Architecture

The application follows a clean architecture pattern with clear separation of concerns:

- **Config**: Configuration management
- **Initializer**: Application startup and database connection
- **Router**: Route definitions and middleware
- **Controllers**: Handle HTTP requests and responses
- **Repository**: Data access layer with database operations
- **Models**: Data structures and business logic

## Contributing

Feel free to submit issues and pull requests.

## License

This project is open source and available under the [MIT License](LICENSE). 