# Development Guidelines

This document provides comprehensive guidelines for developing with Stage Zero, including code standards, architecture patterns, and contribution workflows.

## Table of Contents

- [Code Standards](#code-standards)
- [Architecture Patterns](#architecture-patterns)
- [Adding New Features](#adding-new-features)
- [Database Management](#database-management)
- [Testing Guidelines](#testing-guidelines)
- [Forking & Customization](#forking--customization)
- [Environment Management](#environment-management)
- [Troubleshooting](#troubleshooting)

## Code Standards

### Go Code Style

- **Formatting**: Use `gofmt` and `goimports` for consistent formatting
- **Linting**: Follow `golangci-lint` configuration
- **Naming**: Use descriptive names, avoid abbreviations
- **Comments**: Document all public functions and types
- **Error Handling**: Always handle errors explicitly, use `errors.New()` for custom errors

### Project Structure

Follow the established directory structure:

```
internal/
├── config/          # Configuration management
├── container/       # Dependency injection
├── controller/      # HTTP handlers
├── middleware/      # HTTP middleware
├── model/           # Data models
├── repository/      # Data access layer
├── service/         # Business logic
├── router/          # Route definitions
├── response/        # Response types
├── logger/          # Logging utilities
├── err/             # Error definitions
└── telemetry/       # Observability
```

### Import Organization

```go
import (
    // Standard library
    "context"
    "fmt"
    
    // Third-party packages
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    
    // Internal packages
    "github.com/Verano-20/stage-zero/internal/model"
    "github.com/Verano-20/stage-zero/internal/service"
)
```

## Architecture Patterns

### Dependency Injection

Use the Container pattern for dependency injection:

```go
// In container/container.go
func NewContainerWithDB(db *gorm.DB) *Container {
    userRepository := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepository)
    userController := controller.NewUserController(userService)
    
    return &Container{
        UserRepository: userRepository,
        UserService:    userService,
        UserController: userController,
    }
}
```

### Repository Pattern

Abstract data access with interfaces:

```go
// In repository/user.go
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id uint) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id uint) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}
```

### Service Layer

Implement business logic in services:

```go
// In service/user.go
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*model.User, error)
    GetUser(ctx context.Context, id uint) (*model.User, error)
    UpdateUser(ctx context.Context, id uint, req *UpdateUserRequest) (*model.User, error)
    DeleteUser(ctx context.Context, id uint) error
}

type userService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}
```

### Controller Pattern

Handle HTTP requests and responses:

```go
// In controller/user.go
type UserController struct {
    userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
    return &UserController{userService: userService}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
    var req CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
        return
    }
    
    user, err := c.userService.CreateUser(ctx.Request.Context(), &req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
        return
    }
    
    ctx.JSON(http.StatusCreated, response.SuccessResponse{
        Message: "User created successfully",
        Data:    user,
    })
}
```

## Adding New Features

### 1. Create Model

Define your data model:

```go
// In model/feature.go
type Feature struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null" validate:"required,min=1,max=100"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. Create Repository

Implement data access:

```go
// In repository/feature.go
type FeatureRepository interface {
    Create(ctx context.Context, feature *model.Feature) error
    GetByID(ctx context.Context, id uint) (*model.Feature, error)
    GetAll(ctx context.Context) ([]*model.Feature, error)
    Update(ctx context.Context, feature *model.Feature) error
    Delete(ctx context.Context, id uint) error
}

type featureRepository struct {
    db *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) FeatureRepository {
    return &featureRepository{db: db}
}
```

### 3. Create Service

Implement business logic:

```go
// In service/feature.go
type FeatureService interface {
    CreateFeature(ctx context.Context, req *CreateFeatureRequest) (*model.Feature, error)
    GetFeature(ctx context.Context, id uint) (*model.Feature, error)
    GetAllFeatures(ctx context.Context) ([]*model.Feature, error)
    UpdateFeature(ctx context.Context, id uint, req *UpdateFeatureRequest) (*model.Feature, error)
    DeleteFeature(ctx context.Context, id uint) error
}

type featureService struct {
    featureRepo repository.FeatureRepository
}

func NewFeatureService(featureRepo repository.FeatureRepository) FeatureService {
    return &featureService{featureRepo: featureRepo}
}
```

### 4. Create Controller

Handle HTTP requests:

```go
// In controller/feature.go
type FeatureController struct {
    featureService service.FeatureService
}

func NewFeatureController(featureService service.FeatureService) *FeatureController {
    return &FeatureController{featureService: featureService}
}

func (c *FeatureController) CreateFeature(ctx *gin.Context) {
    // Implementation
}
```

### 5. Update Container

Add to dependency injection:

```go
// In container/container.go
type Container struct {
    // ... existing fields
    
    FeatureRepository repository.FeatureRepository
    FeatureService    service.FeatureService
    FeatureController *controller.FeatureController
}

func NewContainerWithDB(db *gorm.DB) *Container {
    // ... existing code
    
    featureRepository := repository.NewFeatureRepository(db)
    featureService := service.NewFeatureService(featureRepository)
    featureController := controller.NewFeatureController(featureService)
    
    container.FeatureRepository = featureRepository
    container.FeatureService = featureService
    container.FeatureController = featureController
    
    return container
}
```

### 6. Add Routes

Define API endpoints:

```go
// In router/router.go
func setupFeatureRoutes(r *gin.RouterGroup, container *container.Container) {
    features := r.Group("/features")
    {
        features.POST("", container.FeatureController.CreateFeature)
        features.GET("", container.FeatureController.GetAllFeatures)
        features.GET("/:id", container.FeatureController.GetFeature)
        features.PUT("/:id", container.FeatureController.UpdateFeature)
        features.DELETE("/:id", container.FeatureController.DeleteFeature)
    }
}
```

### 7. Create Migration

Add database migration:

```sql
-- In cmd/migrate/migrations/00003_create_features_table.sql
CREATE TABLE features (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_features_name ON features(name);
```

### 8. Add Tests

Create comprehensive tests:

```go
// In test/service/feature_test.go
func TestFeatureService_CreateFeature(t *testing.T) {
    // Mock repository
    mockRepo := &mocks.FeatureRepository{}
    
    // Create service with mock
    service := service.NewFeatureService(mockRepo)
    
    // Test implementation
    // ...
}
```

## Database Management

### Migrations

Use Goose for database migrations:

```bash
# Create new migration
go run cmd/migrate/main.go create add_new_table sql

# Run migrations
go run cmd/migrate/main.go up

# Rollback migration
go run cmd/migrate/main.go down
```

### Migration Best Practices

- **Naming**: Use descriptive names with timestamps
- **Reversibility**: Always provide down migrations
- **Indexes**: Add indexes for frequently queried columns
- **Constraints**: Use appropriate constraints and foreign keys
- **Data**: Avoid data migrations in schema migrations

### Database Connection

Configure connection pooling:

```go
// In database/database.go
func InitDatabase() *gorm.DB {
    config := config.Get()
    
    db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("Failed to connect to database", zap.Error(err))
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get underlying sql.DB", zap.Error(err))
    }
    
    // Configure connection pool
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db
}
```

## Testing Guidelines

### Unit Testing

Test individual components with mocks:

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &mocks.UserRepository{}
    service := service.NewUserService(mockRepo)
    
    req := &service.CreateUserRequest{
        Email:    "test@example.com",
        Password: "SecurePass123!",
    }
    
    expectedUser := &model.User{
        ID:    1,
        Email: req.Email,
    }
    
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
    mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(expectedUser, nil)
    
    // Act
    result, err := service.CreateUser(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser.Email, result.Email)
    mockRepo.AssertExpectations(t)
}
```

### E2E Testing

Test complete workflows:

```typescript
// In test/e2e/tests/feature.spec.ts
test.describe('Feature API', () => {
  test('should create and retrieve feature', async ({ request }) => {
    const apiClient = new ApiClient(request);
    
    // Create feature
    const createResponse = await apiClient.createFeature({
      name: 'Test Feature'
    });
    expect(createResponse.ok()).toBeTruthy();
    
    const feature = await assertResponse<FeatureResponse>(createResponse, 201);
    
    // Retrieve feature
    const getResponse = await apiClient.getFeature(feature.data.id);
    expect(getResponse.ok()).toBeTruthy();
    
    const retrievedFeature = await assertResponse<FeatureResponse>(getResponse, 200);
    expect(retrievedFeature.data.name).toBe('Test Feature');
  });
});
```

### Mock Generation

Generate mocks automatically:

```bash
# Generate mocks
go generate ./...

# Or manually
mockgen -source=internal/repository/user.go -destination=test/mocks/repository/user_mock.go
```

## Forking & Customization

### Fork Setup

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/stage-zero.git
   cd stage-zero
   ```

3. **Update module name** in `go.mod`:
   ```go
   module github.com/YOUR_USERNAME/stage-zero
   ```

4. **Update imports** throughout the codebase:
   ```bash
   find . -name "*.go" -exec sed -i 's|github.com/Verano-20/stage-zero|github.com/YOUR_USERNAME/stage-zero|g' {} \;
   ```

### Customization Points

#### 1. Service Name and Branding

Update configuration:

```go
// In internal/config/config.go
const (
    DefaultServiceName    = "your-service-name"
    DefaultServiceVersion = "1.0.0"
)
```

#### 2. Database Schema

Modify models and create migrations:

```go
// In internal/model/your_model.go
type YourModel struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    // Your fields
}
```

#### 3. API Endpoints

Update routes and controllers:

```go
// In internal/router/router.go
func setupYourRoutes(r *gin.RouterGroup, container *container.Container) {
    yourGroup := r.Group("/your-endpoint")
    {
        yourGroup.GET("", container.YourController.GetAll)
        yourGroup.POST("", container.YourController.Create)
    }
}
```

#### 4. Authentication

Customize JWT claims and validation:

```go
// In internal/service/auth.go
type Claims struct {
    UserID uint   `json:"sub"`
    Email  string `json:"email"`
    Role   string `json:"role"` // Add custom claims
    jwt.RegisteredClaims
}
```

#### 5. Observability

Add custom metrics:

```go
// In internal/telemetry/metrics.go
func NewAppMetrics(meter metric.Meter) (*AppMetrics, error) {
    // Add your custom metrics
    customCounter, err := meter.Int64Counter(
        "your_custom_counter",
        metric.WithDescription("Your custom counter"),
    )
    if err != nil {
        return nil, err
    }
    
    return &AppMetrics{
        // ... existing metrics
        YourCustomCounter: customCounter,
    }, nil
}
```

### Environment Configuration

Create environment-specific configurations:

```env
# .env.development
ENVIRONMENT=development
SERVICE_NAME=your-service-dev
DATABASE_URL=postgres://user:pass@localhost:5432/your_db_dev

# .env.production
ENVIRONMENT=production
SERVICE_NAME=your-service-prod
DATABASE_URL=postgres://user:pass@prod-host:5432/your_db_prod
```

## Environment Management

### Development Environment

Local development with Docker Compose:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Run migrations
docker-compose exec migrate ./migrate up

# Stop services
docker-compose down
```

### Production Environment

Deploy to production:

```bash
# Merge to deployment branch
git checkout deployment
git merge main
git push origin deployment

# Monitor deployment
gh run list --workflow="build-and-deploy.yml"
```

### Environment Variables

Required environment variables:

```env
# Database
POSTGRES_DB=your_database_name
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
DATABASE_URL=postgres://postgres:password@db:5432/your_database_name

# Application
SERVICE_NAME=your-service-name
SERVICE_VERSION=1.0.0
SERVICE_PORT=8080
ENVIRONMENT=development

# JWT
JWT_SECRET=your_jwt_secret_key

# Telemetry
ENABLE_STDOUT=true
ENABLE_OTLP=false
OTLP_ENDPOINT=otel-collector:4317
OTLP_INSECURE=true
METRIC_INTERVAL=30s
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Issues

```bash
# Check database status
docker-compose ps db

# View database logs
docker-compose logs db

# Test connection
docker-compose exec app go run cmd/migrate/main.go status
```

#### 2. Migration Failures

```bash
# Check migration status
go run cmd/migrate/main.go status

# Rollback problematic migration
go run cmd/migrate/main.go down

# Fix migration file and retry
go run cmd/migrate/main.go up
```

#### 3. Test Failures

```bash
# Run tests with verbose output
go test -v ./internal/service

# Run specific test
go test -v -run TestUserService_CreateUser ./internal/service

# Check test coverage
go test -cover ./...
```

#### 4. E2E Test Issues

```bash
# Install Playwright browsers
npx playwright install

# Run tests with debug output
npm run test:e2e:debug

# Check test results
npx playwright show-report
```

#### 5. Docker Issues

```bash
# Rebuild containers
docker-compose build --no-cache

# Clean up containers and volumes
docker-compose down -v
docker system prune -a

# Restart services
docker-compose restart
```

### Debugging

#### Application Logs

```bash
# View application logs
docker-compose logs -f app

# Filter logs by level
docker-compose logs app | grep ERROR
```

#### Database Debugging

```bash
# Connect to database
docker-compose exec db psql -U postgres -d stage_zero

# Check table structure
\dt
\d users
```

#### Performance Monitoring

Access monitoring dashboards:
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Application**: http://localhost:8080/health

### Getting Help

- **Issues**: [GitHub Issues](https://github.com/Verano-20/stage-zero/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Verano-20/stage-zero/discussions)
- **Documentation**: Check the main [README.md](README.md)

---

## Contributing

When contributing to this project:

1. **Fork** the repository
2. **Create** a feature branch
3. **Follow** the code standards outlined in this document
4. **Add tests** for new functionality
5. **Update documentation** as needed
6. **Submit** a pull request

### Pull Request Guidelines

- **Title**: Use clear, descriptive titles
- **Description**: Explain what changes and why
- **Tests**: Ensure all tests pass
- **Documentation**: Update relevant documentation
- **Review**: Address feedback promptly

---

*This guidelines document is maintained alongside the main project. For the most up-to-date information, always refer to the latest version in the repository.*
