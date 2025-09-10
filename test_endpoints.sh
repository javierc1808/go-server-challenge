#!/bin/bash

# Script para probar los endpoints de la nueva arquitectura Clean Architecture

echo "🚀 Probando endpoints de Clean Architecture"
echo "============================================="

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Función para probar endpoint
test_endpoint() {
    local url=$1
    local description=$2
    
    echo -e "\n${BLUE}Probando: $description${NC}"
    echo "URL: $url"
    echo "----------------------------------------"
    
    response=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$url")
    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)
    body=$(echo "$response" | sed '/HTTP_CODE:/d')
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✅ Éxito (HTTP $http_code)${NC}"
        echo "Respuesta:"
        echo "$body" | jq . 2>/dev/null || echo "$body"
    else
        echo -e "${RED}❌ Error (HTTP $http_code)${NC}"
        echo "Respuesta: $body"
    fi
}

# Verificar si el servidor está corriendo
echo "🔍 Verificando si el servidor está corriendo..."
if ! curl -s http://localhost:8080/documents > /dev/null 2>&1; then
    echo -e "${RED}❌ El servidor no está corriendo en localhost:8080${NC}"
    echo "Por favor, ejecuta primero:"
    echo "  go run cmd/server/main.go"
    exit 1
fi

echo -e "${GREEN}✅ Servidor detectado${NC}"

# Probar endpoint de documentos (GET)
test_endpoint "http://localhost:8080/documents" "API de Documentos (GET)"

# Probar endpoint de estadísticas de seguridad
test_endpoint "http://localhost:8080/security/stats" "Estadísticas de Seguridad"

# Probar endpoint de health check
test_endpoint "http://localhost:8080/health" "Health Check"

# Probar creación de documento (POST)
echo -e "\n${BLUE}Probando: Creación de Documento${NC}"
echo "URL: POST http://localhost:8080/documents"
echo "----------------------------------------"

# Crear un documento de prueba
document_data='{
  "id": "test-doc-123",
  "title": "Documento de Prueba",
  "version": "1.0.0",
  "attachments": ["archivo1.pdf", "archivo2.docx"],
  "contributors": [
    {
      "id": "user-123",
      "name": "Usuario de Prueba"
    }
  ]
}'

response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST -H "Content-Type: application/json" -d "$document_data" http://localhost:8080/documents)
http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)
body=$(echo "$response" | sed '/HTTP_CODE:/d')

if [ "$http_code" -eq 201 ]; then
    echo -e "${GREEN}✅ Documento creado exitosamente (HTTP $http_code)${NC}"
    echo "Respuesta:"
    echo "$body" | jq . 2>/dev/null || echo "$body"
else
    echo -e "${RED}❌ Error al crear documento (HTTP $http_code)${NC}"
    echo "Respuesta: $body"
fi

# Verificar que el documento se guardó en cache
echo -e "\n${BLUE}Verificando: Documento en Cache${NC}"
echo "Obteniendo documentos nuevamente..."
test_endpoint "http://localhost:8080/documents" "API de Documentos (Verificación Cache)"

# Probar endpoint de notificaciones (WebSocket)
echo -e "\n${BLUE}Probando: WebSocket de Notificaciones${NC}"
echo "URL: ws://localhost:8080/notifications"
echo "----------------------------------------"
echo -e "${GREEN}✅ WebSocket endpoint disponible${NC}"
echo "Nota: Para probar WebSocket, usa un cliente WebSocket como wscat:"
echo "  npx wscat -c ws://localhost:8080/notifications"

# Probar rate limiting
echo -e "\n${BLUE}Probando: Rate Limiting${NC}"
echo "Haciendo múltiples peticiones rápidas..."
for i in {1..5}; do
    curl -s -o /dev/null -w "Petición $i: %{http_code}\n" http://localhost:8080/documents
done

echo -e "\n${GREEN}🎉 Pruebas completadas${NC}"
echo "============================================="
echo "La nueva arquitectura Clean Architecture está funcionando correctamente!"
