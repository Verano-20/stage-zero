# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /stage-zero

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/server ./cmd/api-server/main.go

# Build the migration binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/migrate ./cmd/migrate/main.go

# Final stage
FROM alpine:latest

WORKDIR /stage-zero

# Install runtime dependencies
RUN apk add --no-cache postgresql-client

# Copy the binaries from builder
COPY --from=builder /bin/server .
COPY --from=builder /bin/migrate .
# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./server"]