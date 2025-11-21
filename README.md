# Personal Portfolio Backend

REST API backend for a personal portfolio application built with Go.

## Table of Contents

- [Technologies Used](#technologies-used)
- [API Endpoints](#api-endpoints)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Related Projects](#related-projects)

## Technologies Used

### **Technologies**
- **Go 1.24** - Programming language
- **Gin Framework** - HTTP web framework
- **PostgreSQL** - Primary database
- **GORM** - ORM for database operations
- **Swagger/OpenAPI** - API documentation
- **Testify** - Testing framework

## API Endpoints

### **Authentication**
```
POST /auth/token
```
- **Description**: Generate JWT token for frontend authentication (not very useful without proper credentials, but easier to change in the future)
- **Body**: `{"auth_key": "your-auth-key"}`
- **Response**: `{"token": "jwt-token", "expires_in": 1800}`

### **About & Contact**
```
GET /about/contact
```
- **Description**: Get contact information
- **Authentication**: Required (JWT)
- **Response**: Contact details from database

```
GET /about/personal-reviews
```
- **Description**: Get personal reviews carousel data
- **Authentication**: Required (JWT)
- **Response**: Array of personal reviews

### **Portfolio Projects**
```
GET /tech/projects
```
- **Description**: Get technology projects
- **Authentication**: Required (JWT)
- **Response**: Array of tech project groups

```
GET /finance/projects
```
- **Description**: Get finance projects
- **Authentication**: Required (JWT)
- **Response**: Array of finance project groups

```
GET /games/projects
```
- **Description**: Get game projects
- **Authentication**: Required (JWT)
- **Response**: Array of game project groups

```
GET /games/played
```
- **Description**: Get recently played games
- **Authentication**: Required (JWT)
- **Response**: Array of recently played games

### **System**
```
GET /health
```
- **Description**: Health check endpoint
- **Authentication**: Not required
- **Response**: `{"status": "healthy"}`

### **Documentation**
```
GET /
GET /swagger/*
```
- **Description**: Swagger API documentation
- **Authentication**: Not required (development only)

## Development

### **Development Environment Setup**

For development environment please check https://github.com/David-Xilo/personal-portfolio-orchestration

## Testing

### **Essential Go Commands**

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for specific package
go test ./src/internal/config

# Run specific test
go test -run TestLoadConfig ./src/internal/config

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Clean test cache
go clean -testcache

# Build the application
go build -o main ./src/cmd/api

# Format code
go fmt ./...

# Lint code
go vet ./...

# Security scan
gosec ./...

# Check for vulnerabilities
go list -json -m all | nancy sleuth

# Update dependencies
go mod tidy
go mod verify
```

## Deployment

### **Docker Deployment**

```bash
# Build locally
docker build -t ${BACKEND_IMAGE} ${BACKEND_DOCKERFILE}

# Run container
docker run \
        -e ENV=development \
        -e DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" \
        -e FRONTEND_URL=${FRONTEND_URL} \
        -e PORT=${BACKEND_PORT} \
        --network ${NETWORK_NAME} \
        --name ${BACKEND_CONTAINER} \
        -p ${BACKEND_PORT}:${BACKEND_PORT} \
        -d ${BACKEND_IMAGE}
```

## Related Projects

- **Frontend**: https://github.com/David-Xilo/personal-portfolio-main-front
- **Infrastructure**: https://github.com/David-Xilo/personal-portfolio-orchestration
- **Database**: https://github.com/David-Xilo/personal-portfolio-db-schema

---
