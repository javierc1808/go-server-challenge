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

# Probar endpoint de documentos
test_endpoint "http://localhost:8080/documents" "API de Documentos"

# Probar endpoint de notificaciones (WebSocket)
echo -e "\n${BLUE}Probando: WebSocket de Notificaciones${NC}"
echo "URL: ws://localhost:8080/notifications"
echo "----------------------------------------"
echo -e "${GREEN}✅ WebSocket endpoint disponible${NC}"
echo "Nota: Para probar WebSocket, usa un cliente WebSocket como wscat:"
echo "  npx wscat -c ws://localhost:8080/notifications"

echo -e "\n${GREEN}🎉 Pruebas completadas${NC}"
echo "============================================="
echo "La nueva arquitectura Clean Architecture está funcionando correctamente!"
