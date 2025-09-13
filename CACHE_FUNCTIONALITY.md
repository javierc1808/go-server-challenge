# 🗄️ Implemented Cache Functionality

This document describes the in-memory cache system implemented for documents in the project.

## 📋 Feature Summary

| Feature | Status | Description |
|---------|--------|-------------|
| **In-memory Cache** | ✅ Implemented | Temporary storage for documents |
| **Configurable TTL** | ✅ Implemented | Default time-to-live: 24 hours |
| **Automatic Cleanup** | ✅ Implemented | Removal of expired items |
| **Statistics** | ✅ Implemented | Cache usage metrics |
| **Session Persistence** | ✅ Implemented | Documents persist during the server session |
| **Loss on Restart** | ✅ Implemented | Documents are lost when the server restarts |

## 🏗️ Cache Architecture

### Cache Interface
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

### In-memory Implementation
```go
type MemoryCache struct {
    documents map[string]*entity.Document
    mutex     sync.RWMutex
    ttl       time.Duration
    expiry    map[string]time.Time
}
```

## 🔄 How It Works

### 1) First Load (Server just started)
```
GET /documents
    ↓
Empty cache
    ↓
Generate simulated documents
    ↓
Store in cache
    ↓
Return documents
```

### 2) Subsequent Loads
```
GET /documents
    ↓
Cache contains data
    ↓
Return cached documents
```

### 3) Document Creation
```
POST /documents
    ↓
Validate document
    ↓
Store in cache
    ↓
Return created document
```

### 4) Server Restart
```
Server stops
    ↓
Cache is lost (in-memory)
    ↓
Server starts
    ↓
Cache is empty again
```

## 📊 Cache Characteristics

### TTL (Time To Live)
- Default: 24 hours
- Automatic cleanup every 5 minutes
- Expired documents are ignored and then removed by the cleaner

### Thread Safety
- Mutex protection for concurrent access
- Read/Write locks to optimize multiple readers
- Operations are consistent and safe

### Real-time Statistics
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

## 🚀 Available Endpoints

### GET /documents
- Description: Fetch all documents
- Behavior:
  - If cache has documents → returns from cache
  - If cache is empty → generates simulated documents and stores them

### POST /documents
- Description: Create a new document
- Body: JSON document payload
- Behavior: Stores the document in the cache

### GET /security/stats
- Description: System stats including cache metrics
- Includes: Document metrics, TTL, expired elements, and more

## 🧪 Functional Tests

### Updated Test Script
```bash
./test_endpoints.sh
```

Includes:
- ✅ GET /documents (first load)
- ✅ POST /documents (create)
- ✅ Cache persistence verification
- ✅ Stats verification

### Manual Tests

#### 1) Verify Session Persistence
```bash
# First request (generates documents)
curl http://localhost:8080/documents

# Second request (should return the same set from cache)
curl http://localhost:8080/documents
```

#### 2) Create a Custom Document
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{
    "id": "my-document-123",
    "title": "My Document",
    "version": "1.0.0",
    "attachments": ["file.pdf"],
    "contributors": [{"id": "user-1", "name": "User 1"}]
  }' \
  http://localhost:8080/documents
```

#### 3) Verify Statistics
```bash
curl http://localhost:8080/security/stats | jq .cache
```

## 📁 Implemented Files

### New Files
- `internal/domain/repository/cache_repository.go` – Cache interface
- `internal/infrastructure/repository/memory_cache.go` – In-memory implementation

### Modified Files
- `internal/infrastructure/repository/document_repository_impl.go` – Cache integration
- `internal/delivery/http/document_handler.go` – POST endpoint to create documents
- `cmd/server/main.go` – Cache initialization
- `scripts/test_endpoints.sh` – Updated tests

## 🔧 Configuration

### Cache TTL
```go
// In main.go
cache := repository.NewMemoryCache(24 * time.Hour) // 24 hours
```

### Automatic Cleanup
```go
// Cleanup runs every 5 minutes
go cache.startCleanup()
```

### Statistics
```go
// Access cache statistics
stats := cache.GetStats()
```

## 📈 Benefits

### Performance
- ⚡ Faster responses for existing documents
- 🔄 Reuse of generated data
- 💾 Less processing for repeated requests

### Functionality
- 📝 Create custom documents
- 🔍 Persistence during server session
- 📊 Real-time cache monitoring

### User Experience
- 🚀 Consistent responses for the same dataset
- ✨ Custom documents remain available during the session
- 🔄 Predictable behavior after server restarts

## ⚠️ Important Considerations

### Current Limitations
- Memory: documents are stored in RAM
- Data Loss: documents are lost on server restart
- Scalability: not designed for distributed environments

### Ideal Use Cases
- Development: quick tests with persistent data during a session
- Demos: show functionality with real-looking data
- Prototypes: validate concepts without a database

### Recommended Next Steps
1. Implement database persistence
2. Add distributed cache (Redis)
3. Implement invalidation strategies
4. Add performance metrics
5. Configure cleanup policies

## 🎯 Use Cases

### Scenario 1: Local Development
```bash
# Start the server
go run cmd/server/main.go

# Create custom documents
curl -X POST ... /documents

# Documents persist during the session
curl /documents  # Returns created + simulated documents

# After restarting the server, custom documents are lost
# and simulated ones are re-generated
```

### Scenario 2: Demo
```bash
# Show basic functionality
curl /documents

# Create specific documents for the demo
curl -X POST ... /documents

# Verify custom documents are available
curl /documents
```

### Scenario 3: Testing
```bash
# Verify initial behavior
curl /documents

# Create test data
curl -X POST ... /documents

# Verify persistence
curl /documents

# Restart the server and verify reset
# (Custom documents are lost)
```

---

**The cache system is fully functional and ready to use!** 🚀🗄️
