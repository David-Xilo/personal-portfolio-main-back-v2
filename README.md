# Personal Portfolio Backend

REST API backend for a personal portfolio application built with Go.

## Table of Contents

- [Technologies Used](#technologies-used)
- [API Endpoints](#api-endpoints)
- [Development](#development)
- [Testing](#testing)
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

### **Contact**
```
GET /contact
```
- **Description**: Get contact information
- **Response**: Contact details from database

### **Portfolio Projects**
```
GET /projects
```
- **Description**: Get projects and repositories
- **Response**: Array of tech project groups with repositories

### **Documentation**
```
GET /
GET /swagger/*
```
- **Description**: Swagger API documentation
- **Authentication**: Not required (development only)

## Development

### **Development Environment Setup**

Check https://github.com/David-Xilo/personal-portfolio-orchestration

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
go build -o main ./src/cmd

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

## Related Projects

- **Frontend**: https://github.com/David-Xilo/personal-portfolio-main-front-v2
- **Infrastructure**: https://github.com/David-Xilo/personal-portfolio-orchestration
- **Database**: https://github.com/David-Xilo/personal-portfolio-db-schema

---
