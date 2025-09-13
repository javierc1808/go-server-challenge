# Clean Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        DELIVERY LAYER                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   HTTP Handler  │  │ WebSocket Handler│  │  CLI Handler    │ │
│  │  (REST API)     │  │  (Real-time)    │  │  (Commands)     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                        USE CASE LAYER                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │ DocumentUsecase │  │NotificationUsecase│ │   UserUsecase   │ │
│  │  (Business      │  │  (Business      │  │  (Business      │ │
│  │   Logic)        │  │   Logic)        │  │   Logic)        │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                         DOMAIN LAYER                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │    Entities     │  │   Repositories  │  │    Errors       │ │
│  │  (User, Doc,    │  │  (Interfaces)   │  │  (Domain        │ │
│  │   Notification) │  │                 │  │   Errors)       │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                ▲
                                │
┌─────────────────────────────────────────────────────────────────┐
│                     INFRASTRUCTURE LAYER                       │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Repository    │  │    Database     │  │   External      │ │
│  │  Implementations│  │   Connections   │  │   Services      │ │
│  │  (PostgreSQL,   │  │                 │  │  (APIs, Cache)  │ │
│  │   MongoDB, etc) │  │                 │  │                 │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Dependency Flow

1. Delivery Layer → Use Case Layer → Domain Layer
2. Infrastructure Layer → Domain Layer (implements interfaces)
3. Use Case Layer → Domain Layer (uses interfaces)

## Applied Principles

### ✅ Dependency Inversion
- Inner layers do not know outer layers
- Interfaces live in the domain layer
- Implementations live in infrastructure

### ✅ Separation of Concerns
- Each layer has a single responsibility
- Domain contains only business logic
- Use cases orchestrate data flow

### ✅ Testability
- Each layer can be tested independently
- Dependencies can be mocked
- Business logic is isolated

## File Structure

```
internal/
├── domain/                    # 🎯 BUSINESS CORE
│   ├── entity/               # Domain entities
│   │   ├── user.go
│   │   ├── document.go
│   │   ├── notification.go
│   │   └── errors.go
│   └── repository/           # Repository interfaces
│       ├── user_repository.go
│       ├── document_repository.go
│       └── notification_repository.go
├── usecase/                  # 🔄 APPLICATION LOGIC
│   ├── document_usecase.go
│   └── notification_usecase.go
├── infrastructure/           # 🔧 IMPLEMENTATIONS
│   └── repository/
│       ├── user_repository_impl.go
│       ├── document_repository_impl.go
│       └── notification_repository_impl.go
└── delivery/                 # 🌐 ENTRY POINTS
    ├── http/
    │   └── document_handler.go
    └── websocket/
        └── notification_handler.go
```

## Benefits of this Architecture

1. Maintainability: organized and easy to evolve
2. Scalability: easy to add new features
3. Testability: each component is independently testable
4. Flexibility: implementations are swappable
5. Independence: domain does not depend on frameworks
6. Reusability: use cases can be reused
