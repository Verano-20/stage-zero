#!/bin/bash

# Load environment variables from .env.local
export $(cat .env.local | xargs)

# Run database migrations
go run cmd/migrate/main.go up

# Run the application
go run cmd/api-server/main.go 