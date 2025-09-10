# Frontend Challenge - Clean Architecture

Este proyecto ha sido refactorizado siguiendo los principios de **Clean Architecture** para mejorar la mantenibilidad, testabilidad y escalabilidad del código.

## 🏗️ Estructura del Proyecto

```
frontend-challenge/
├── cmd/
│   └── server/                 # Punto de entrada de la aplicación
│       └── main.go
├── internal/
│   ├── domain/                 # Capa de Dominio (Entidades y Reglas de Negocio)
│   │   ├── entity/            # Entidades del dominio
│   │   │   ├── user.go
│   │   │   ├── document.go
│   │   │   ├── notification.go
│   │   │   └── errors.go
│   │   └── repository/        # Interfaces de repositorios
│   │       ├── user_repository.go
│   │       ├── document_repository.go
│   │       └── notification_repository.go
│   ├── usecase/               # Capa de Casos de Uso (Lógica de Aplicación)
│   │   ├── document_usecase.go
│   │   └── notification_usecase.go
│   ├── infrastructure/        # Capa de Infraestructura (Implementaciones)
│   │   └── repository/
│   │       ├── user_repository_impl.go
│   │       ├── document_repository_impl.go
│   │       └── notification_repository_impl.go
│   └── delivery/              # Capa de Entrega (HTTP, WebSocket, etc.)
│       ├── http/
│       │   └── document_handler.go
│       └── websocket/
│           └── notification_handler.go
├── pkg/                       # Paquetes compartidos
│   ├── config/
│   │   └── config.go
│   └── logger/
│       └── logger.go
├── go.mod
├── go.sum
└── README.md
```

## 🎯 Principios de Clean Architecture Aplicados

### 1. **Separación de Responsabilidades**
- **Domain**: Contiene las entidades y reglas de negocio puras
- **Use Cases**: Contiene la lógica de aplicación específica
- **Infrastructure**: Implementa las interfaces definidas en el dominio
- **Delivery**: Maneja la comunicación externa (HTTP, WebSocket)

### 2. **Inversión de Dependencias**
- Las capas internas no dependen de las externas
- Las interfaces están definidas en el dominio
- Las implementaciones están en la infraestructura

### 3. **Testabilidad**
- Cada capa puede ser probada independientemente
- Las dependencias se pueden mockear fácilmente
- La lógica de negocio está aislada

## 🚀 Cómo Ejecutar

### Requisitos
- Go 1.21 o superior

### Instalación y Ejecución
```bash
# Instalar dependencias
go mod tidy

# Ejecutar el servidor
go run cmd/server/main.go

# O con flags personalizados
go run cmd/server/main.go -addr localhost:9090
```

## 📡 Endpoints Disponibles

### 1. **API de Documentos**
```
GET http://localhost:8080/documents
```
Retorna una lista de documentos con sus metadatos.

### 2. **Notificaciones en Tiempo Real**
```
WS ws://localhost:8080/notifications
```
Conexión WebSocket que envía notificaciones simuladas en tiempo real.

## 🔧 Mejoras Implementadas

### **Seguridad**
- Headers de seguridad añadidos
- Validación de entrada en entidades
- Manejo seguro de errores
- Logging estructurado

### **Arquitectura**
- Separación clara de responsabilidades
- Inyección de dependencias
- Interfaces bien definidas
- Código más mantenible y testeable

### **Configuración**
- Configuración centralizada
- Manejo de señales del sistema
- Shutdown graceful del servidor
- Timeouts configurables

## 🧪 Testing

La nueva arquitectura facilita el testing:

```go
// Ejemplo de test unitario para un caso de uso
func TestDocumentUsecase_GetAllDocuments(t *testing.T) {
    // Mock del repositorio
    mockRepo := &MockDocumentRepository{}
    
    // Crear caso de uso con mock
    usecase := NewDocumentUsecase(mockRepo, mockUserRepo)
    
    // Ejecutar test
    documents, err := usecase.GetAllDocuments(context.Background())
    
    // Verificar resultados
    assert.NoError(t, err)
    assert.NotNil(t, documents)
}
```

## 🔄 Migración desde Código Original

El código original en `server.go` ha sido refactorizado manteniendo la misma funcionalidad pero con una arquitectura mucho más robusta:

- ✅ Misma funcionalidad de endpoints
- ✅ Misma generación de datos simulados
- ✅ Mejor separación de responsabilidades
- ✅ Código más testeable
- ✅ Mejor manejo de errores
- ✅ Configuración más flexible

## 📈 Beneficios de la Nueva Arquitectura

1. **Mantenibilidad**: Código más fácil de mantener y extender
2. **Testabilidad**: Cada componente puede ser probado independientemente
3. **Escalabilidad**: Fácil añadir nuevas funcionalidades
4. **Flexibilidad**: Fácil cambiar implementaciones (BD, cache, etc.)
5. **Seguridad**: Mejor manejo de errores y validaciones
6. **Performance**: Mejor gestión de recursos y timeouts

## 🚧 Próximos Pasos Recomendados

1. **Implementar tests unitarios** para cada capa
2. **Añadir autenticación y autorización**
3. **Implementar base de datos real** (PostgreSQL, MongoDB)
4. **Añadir cache** (Redis)
5. **Implementar métricas y monitoreo**
6. **Añadir documentación API** (Swagger)
7. **Implementar CI/CD pipeline**
