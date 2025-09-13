# Frontend Challenge - Clean Architecture

This project was refactored following Clean Architecture principles to improve maintainability, testability, and scalability.

## 🏗️ Project Structure

```
frontend-challenge/
├── cmd/
│   └── server/                 # Application entrypoint
│       └── main.go
├── internal/
│   ├── domain/                 # Domain layer (Entities & Business Rules)
│   │   ├── entity/            # Domain entities
│   │   │   ├── user.go
│   │   │   ├── document.go
│   │   │   ├── notification.go
│   │   │   └── errors.go
│   │   └── repository/        # Repository interfaces
│   │       ├── user_repository.go
│   │       ├── document_repository.go
│   │       └── notification_repository.go
│   ├── usecase/               # Use Case layer (Application Logic)
│   │   ├── document_usecase.go
│   │   └── notification_usecase.go
│   ├── infrastructure/        # Infrastructure layer (Implementations)
│   │   └── repository/
│   │       ├── user_repository_impl.go
│   │       ├── document_repository_impl.go
│   │       └── notification_repository_impl.go
│   └── delivery/              # Delivery layer (HTTP, WebSocket, etc.)
│       ├── http/
│       │   └── document_handler.go
│       └── websocket/
│           └── notification_handler.go
├── pkg/                       # Shared packages
│   ├── config/
│   │   └── config.go
│   └── logger/
│       └── logger.go
├── go.mod
├── go.sum
└── README.md
```

## 🎯 Applied Clean Architecture Principles

### 1. Separation of Concerns
- Domain: pure entities and business rules
- Use Cases: application-specific logic
- Infrastructure: implements domain interfaces
- Delivery: external communication (HTTP, WebSocket)

### 2. Dependency Inversion
- Inner layers do not depend on outer layers
- Interfaces are defined in the domain
- Implementations live in infrastructure

### 3. Testability
- Each layer can be tested independently
- Dependencies can be easily mocked
- Business logic is isolated

## 🚀 How to Run

### Requirements
- Go 1.21 or later

### Install and Run
```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go

# Or with custom flags
go run cmd/server/main.go -addr localhost:9090
```

## 📡 Available Endpoints

### 1. Documents API
```
GET http://localhost:8080/documents
```
Returns a list of documents with metadata.

### 2. Real-time Notifications
```
WS ws://localhost:8080/notifications
```
WebSocket connection that emits notifications when orders/documents are created/updated/deleted.

## 🔧 Improvements

### Security
- Security headers
- Input validation in entities
- Safer error handling
- Structured logging

### Architecture
- Clear separation of concerns
- Dependency injection
- Well-defined interfaces
- More maintainable and testable code

### Configuration
- Centralized configuration
- OS signals handling
- Graceful server shutdown
- Configurable timeouts

## 🧪 Testing

The new architecture makes testing easier:

```go
// Example unit test for a use case
func TestDocumentUsecase_GetAllDocuments(t *testing.T) {
    // Mock repository
    mockRepo := &MockDocumentRepository{}
    
    // Create use case with mock
    usecase := NewDocumentUsecase(mockRepo, mockUserRepo)
    
    // Run test
    documents, err := usecase.GetAllDocuments(context.Background())
    
    // Assert results
    assert.NoError(t, err)
    assert.NotNil(t, documents)
}
```

## 🔄 Migration from Original Code

The original `server.go` was refactored keeping functionality while improving architecture:

- ✅ Same endpoint functionality
- ✅ Same simulated data generation
- ✅ Better separation of concerns
- ✅ More testable code
- ✅ Improved error handling
- ✅ More flexible configuration

## 📈 Benefits of the New Architecture

1. Maintainability: easier to maintain and extend
2. Testability: each component can be tested independently
3. Scalability: easy to add features
4. Flexibility: swap implementations (DB, cache, etc.)
5. Security: improved error handling and validations
6. Performance: better resource and timeout management

## 🚧 Recommended Next Steps

1. Implement unit tests across layers
2. Add authentication and authorization
3. Implement a real database (PostgreSQL, MongoDB)
4. Add cache (Redis)
5. Implement metrics and monitoring
6. Add API documentation (Swagger)
7. Add a CI/CD pipeline
