# Go CRUD API

A simple CRUD (Create, Read, Update, Delete) REST API built with Go, using:
- [Gin](https://github.com/gin-gonic/gin) for the web framework
- [GORM](https://gorm.io) as the ORM
- PostgreSQL as the database

## Prerequisites

- Go 1.21 or later
- PostgreSQL installed and running
- Git (optional)

## Setup

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
   - Configure your database connection by setting the following environment variables:
     ```bash
     DB_HOST=localhost      # PostgreSQL host
     DB_USER=postgres      # PostgreSQL user
     DB_PASSWORD=postgres  # PostgreSQL password
     DB_NAME=go_crud      # Database name
     DB_PORT=5432         # PostgreSQL port
     ```

## Running the Application

Start the server:
```bash
go run cmd/api-server/main.go
```

The server will start on `http://localhost:8080`

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
├── go.mod                 # Go module file
└── README.md             # This file
```

## Contributing

Feel free to submit issues and pull requests.

## License

This project is open source and available under the [MIT License](LICENSE). 