#!/bin/bash

# Load environment variables from .env.local
export $(cat .env.local | xargs)

# Run the application
go run cmd/api-server/main.go 