# Go CRUD API

A simple CRUD (Create, Read, Update, Delete) REST API built with Go, using:
- [Gin](https://github.com/gin-gonic/gin) for the web framework
- [GORM](https://gorm.io) as the ORM
- PostgreSQL as the database

## Prerequisites

- Go 1.24 or later
- PostgreSQL installed and running
- Git (optional)
- Docker and Docker Compose (optional)

## Environment Setup

The application uses environment variables for configuration. Two environment files are provided:

- `.env.local`: For running the application locally
- `.env.docker`: For running the application with Docker

The environment files contain two sets of variables:

1. Application variables (used by the Go application):
```bash
DB_HOST=localhost  # For local, 'db' for Docker
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_crud
DB_PORT=5432
```

2. PostgreSQL variables (used by PostgreSQL container):
```bash
POSTGRES_USER=postgres      # Creates the PostgreSQL user
POSTGRES_PASSWORD=postgres  # Sets the user's password
POSTGRES_DB=go_crud        # Creates the database
```

Note: Both sets of variables should match (e.g., `DB_USER` should match `POSTGRES_USER`) to ensure proper connectivity.

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
   - Create a new PostgreSQL database for the project
   - Adjust `.env.local` if your PostgreSQL configuration differs from the defaults

4. Start the server:
```bash
./scripts/run-local.sh
```

### Option 2: Docker Setup

1. Clone the repository:
```bash
git clone https://github.com/Verano-20/go-crud.git
cd go-crud
```

2. Build and run with Docker Compose:
```bash
docker-compose up --build
```

This will:
- Build the Go application
- Start PostgreSQL database with health checks
- Load environment variables from `.env.docker`
- Create a persistent volume for the database
- Expose the API on port 8080
- Wait for PostgreSQL to be healthy before starting the application

The application includes several reliability features:
- Health checks to ensure PostgreSQL is ready before starting the app
- Automatic restarts if services fail
- Persistent volume for database data

To stop the services:
```bash
docker-compose down
```

To remove the persistent volume as well:
```bash
docker-compose down -v
```

### Connecting to the Database

You can connect to the PostgreSQL database using any PostgreSQL client (e.g., DBeaver, psql) with these credentials:

```
Host: localhost
Port: 5432
Database: go_crud
Username: postgres
Password: postgres
```

Note: The database must be running (either locally or in Docker) before connecting.

## API Endpoints

The API provides the following endpoints:

### Simple Resource

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/simple` | Create a new simple resource |
| GET | `/simple` | Get all simple resources |
| GET | `/simple/:id` | Get a simple resource by ID |
| PUT | `/simple/:id` | Update a simple resource |
| DELETE | `/simple/:id` | Delete a simple resource |

## Project Structure

```
go-crud/
├── cmd/
│   └── api-server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/           # HTTP request handlers
│   │   └── simple.go
│   ├── models/            # Database models
│   │   └── simple.go
│   └── initializers/      # Application initialization code
├── scripts/
│   └── run-local.sh       # Script to run locally
├── .env.local            # Local environment configuration
├── .env.docker          # Docker environment configuration
├── Dockerfile           # Docker build instructions
├── docker-compose.yml   # Docker Compose configuration
├── go.mod               # Go module file
└── README.md           # This file
```

## Contributing

Feel free to submit issues and pull requests.

## License

This project is open source and available under the [MIT License](LICENSE). 