# 🗄️ Funcionalidad de Cache Implementada

Este documento describe el sistema de cache implementado para los documentos en el proyecto.

## 📋 Resumen de Funcionalidades

| Funcionalidad | Estado | Descripción |
|---------------|--------|-------------|
| **Cache en Memoria** | ✅ Implementado | Almacenamiento temporal de documentos |
| **TTL Configurable** | ✅ Implementado | Tiempo de vida de 24 horas por defecto |
| **Limpieza Automática** | ✅ Implementado | Eliminación de elementos expirados |
| **Estadísticas** | ✅ Implementado | Métricas de uso del cache |
| **Persistencia de Sesión** | ✅ Implementado | Los documentos persisten durante la ejecución |
| **Pérdida al Reiniciar** | ✅ Implementado | Los documentos se pierden al reiniciar el servidor |

## 🏗️ Arquitectura del Cache

### **Interfaz de Cache**
```go
type CacheRepository interface {
    Set(ctx context.Context, key string, document *entity.Document) error
    Get(ctx context.Context, key string) (*entity.Document, error)
    GetAll(ctx context.Context) ([]*entity.Document, error)
    Delete(ctx context.Context, key string) error
    Clear(ctx context.Context) error
    Exists(ctx context.Context, key string) bool
    Count(ctx context.Context) int
}
```

### **Implementación en Memoria**
```go
type MemoryCache struct {
    documents map[string]*entity.Document
    mutex     sync.RWMutex
    ttl       time.Duration
    expiry    map[string]time.Time
}
```

## 🔄 Flujo de Funcionamiento

### **1. Primera Carga (Servidor Iniciado)**
```
GET /documents
    ↓
Cache vacío
    ↓
Generar documentos simulados
    ↓
Almacenar en cache
    ↓
Devolver documentos
```

### **2. Cargas Subsecuentes**
```
GET /documents
    ↓
Cache con datos
    ↓
Devolver documentos del cache
```

### **3. Creación de Documento**
```
POST /documents
    ↓
Validar documento
    ↓
Almacenar en cache
    ↓
Devolver documento creado
```

### **4. Reinicio del Servidor**
```
Servidor se detiene
    ↓
Cache se pierde (memoria)
    ↓
Servidor se inicia
    ↓
Cache vacío nuevamente
```

## 📊 Características del Cache

### **TTL (Time To Live)**
- **Configuración**: 24 horas por defecto
- **Limpieza**: Automática cada 5 minutos
- **Expiración**: Los documentos se marcan como expirados

### **Thread Safety**
- **Mutex**: Protección contra acceso concurrente
- **Read/Write Locks**: Optimización para lecturas múltiples
- **Operaciones Atómicas**: Garantía de consistencia

### **Estadísticas en Tiempo Real**
```json
{
  "cache": {
    "total_documents": 15,
    "active_documents": 12,
    "expired_documents": 3,
    "ttl_seconds": 86400
  }
}
```

## 🚀 Endpoints Disponibles

### **GET /documents**
- **Descripción**: Obtiene todos los documentos
- **Comportamiento**: 
  - Si hay documentos en cache → devuelve del cache
  - Si cache está vacío → genera documentos simulados y los almacena

### **POST /documents**
- **Descripción**: Crea un nuevo documento
- **Body**: JSON con datos del documento
- **Comportamiento**: Almacena el documento en cache

### **GET /security/stats**
- **Descripción**: Estadísticas del sistema incluyendo cache
- **Incluye**: Métricas de documentos, TTL, elementos expirados

## 🧪 Pruebas de Funcionalidad

### **Script de Prueba Actualizado**
```bash
./test_endpoints.sh
```

**Incluye:**
- ✅ Prueba de GET /documents (carga inicial)
- ✅ Prueba de POST /documents (creación)
- ✅ Verificación de persistencia en cache
- ✅ Prueba de estadísticas

### **Pruebas Manuales**

#### **1. Verificar Persistencia de Sesión**
```bash
# Primera petición (genera documentos)
curl http://localhost:8080/documents

# Segunda petición (debe devolver los mismos documentos)
curl http://localhost:8080/documents
```

#### **2. Crear Documento Personalizado**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{
    "id": "mi-documento-123",
    "title": "Mi Documento",
    "version": "1.0.0",
    "attachments": ["archivo.pdf"],
    "contributors": [{"id": "user-1", "name": "Usuario 1"}]
  }' \
  http://localhost:8080/documents
```

#### **3. Verificar Estadísticas**
```bash
curl http://localhost:8080/security/stats | jq .cache
```

## 📁 Archivos Implementados

### **Nuevos Archivos**
- `internal/domain/repository/cache_repository.go` - Interfaz del cache
- `internal/infrastructure/repository/memory_cache.go` - Implementación en memoria

### **Archivos Modificados**
- `internal/infrastructure/repository/document_repository_impl.go` - Integración con cache
- `internal/delivery/http/document_handler.go` - Endpoint POST para crear documentos
- `cmd/server/main.go` - Inicialización del cache
- `test_endpoints.sh` - Pruebas actualizadas

## 🔧 Configuración

### **TTL del Cache**
```go
// En main.go
cache := repository.NewMemoryCache(24 * time.Hour) // 24 horas
```

### **Limpieza Automática**
```go
// Cada 5 minutos se ejecuta la limpieza
go cache.startCleanup()
```

### **Estadísticas**
```go
// Acceso a estadísticas
stats := cache.GetStats()
```

## 📈 Beneficios Obtenidos

### **Rendimiento**
- ⚡ **Respuestas más rápidas** para documentos existentes
- 🔄 **Reutilización de datos** generados
- 💾 **Menos procesamiento** en peticiones repetidas

### **Funcionalidad**
- 📝 **Creación de documentos** personalizados
- 🔍 **Persistencia durante la sesión** del servidor
- 📊 **Monitoreo en tiempo real** del cache

### **Experiencia de Usuario**
- 🚀 **Respuestas consistentes** para el mismo conjunto de datos
- ✨ **Documentos personalizados** se mantienen disponibles
- 🔄 **Comportamiento predecible** al reiniciar el servidor

## ⚠️ Consideraciones Importantes

### **Limitaciones Actuales**
- **Memoria**: Los documentos se almacenan en RAM
- **Pérdida de Datos**: Al reiniciar el servidor se pierden los documentos
- **Escalabilidad**: No funciona en entornos distribuidos

### **Casos de Uso Ideales**
- **Desarrollo**: Pruebas rápidas con datos persistentes
- **Demostraciones**: Mostrar funcionalidad con datos reales
- **Prototipos**: Validación de conceptos sin base de datos

### **Próximos Pasos Recomendados**
1. **Implementar persistencia** en base de datos
2. **Añadir cache distribuido** (Redis)
3. **Implementar estrategias de invalidación**
4. **Añadir métricas de rendimiento**
5. **Configurar políticas de limpieza**

## 🎯 Casos de Uso

### **Escenario 1: Desarrollo Local**
```bash
# Iniciar servidor
go run cmd/server/main.go

# Crear documentos personalizados
curl -X POST ... /documents

# Los documentos persisten durante la sesión
curl /documents  # Devuelve documentos creados + simulados

# Al reiniciar el servidor, se pierden los personalizados
# pero se regeneran los simulados
```

### **Escenario 2: Demostración**
```bash
# Mostrar funcionalidad básica
curl /documents

# Crear documentos específicos para la demo
curl -X POST ... /documents

# Verificar que los documentos personalizados están disponibles
curl /documents
```

### **Escenario 3: Testing**
```bash
# Verificar comportamiento inicial
curl /documents

# Crear datos de prueba
curl -X POST ... /documents

# Verificar persistencia
curl /documents

# Reiniciar servidor y verificar reset
# (Los documentos personalizados se pierden)
```

---

**¡El sistema de cache está completamente funcional y listo para usar!** 🚀🗄️
