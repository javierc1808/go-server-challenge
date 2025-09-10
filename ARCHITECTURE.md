# Diagrama de Clean Architecture

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

## Flujo de Dependencias

1. **Delivery Layer** → **Use Case Layer** → **Domain Layer**
2. **Infrastructure Layer** → **Domain Layer** (implementa interfaces)
3. **Use Case Layer** → **Domain Layer** (usa interfaces)

## Principios Aplicados

### ✅ Inversión de Dependencias
- Las capas internas no conocen las externas
- Las interfaces están en el dominio
- Las implementaciones están en infraestructura

### ✅ Separación de Responsabilidades
- Cada capa tiene una responsabilidad específica
- El dominio contiene solo lógica de negocio
- Los casos de uso orquestan el flujo de datos

### ✅ Testabilidad
- Cada capa puede ser probada independientemente
- Las dependencias se pueden mockear
- La lógica de negocio está aislada

## Estructura de Archivos

```
internal/
├── domain/                    # 🎯 NÚCLEO DEL NEGOCIO
│   ├── entity/               # Entidades del dominio
│   │   ├── user.go
│   │   ├── document.go
│   │   ├── notification.go
│   │   └── errors.go
│   └── repository/           # Interfaces de repositorios
│       ├── user_repository.go
│       ├── document_repository.go
│       └── notification_repository.go
├── usecase/                  # 🔄 LÓGICA DE APLICACIÓN
│   ├── document_usecase.go
│   └── notification_usecase.go
├── infrastructure/           # 🔧 IMPLEMENTACIONES
│   └── repository/
│       ├── user_repository_impl.go
│       ├── document_repository_impl.go
│       └── notification_repository_impl.go
└── delivery/                 # 🌐 PUNTOS DE ENTRADA
    ├── http/
    │   └── document_handler.go
    └── websocket/
        └── notification_handler.go
```

## Beneficios de esta Arquitectura

1. **Mantenibilidad**: Código organizado y fácil de mantener
2. **Escalabilidad**: Fácil añadir nuevas funcionalidades
3. **Testabilidad**: Cada componente es testeable independientemente
4. **Flexibilidad**: Fácil cambiar implementaciones
5. **Independencia**: El dominio no depende de frameworks externos
6. **Reutilización**: Los casos de uso pueden ser reutilizados
