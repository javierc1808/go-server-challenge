#!/usr/bin/env bash
set -euo pipefail

# Config
PORT="${PORT:-8080}"
ADDR=":${PORT}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="${SCRIPT_DIR}/.."

BLUE='\033[1;34m'
GREEN='\033[1;32m'
RED='\033[1;31m'
NC='\033[0m'

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo -e "${RED}ERROR${NC}: required command not found: $1" >&2
    return 1
  fi
}

detect_platform() {
  local uname_s uname_m
  uname_s=$(uname -s)
  uname_m=$(uname -m)

  case "$uname_s" in
    Darwin) OS="mac" ;;
    Linux)
      if grep -qi microsoft /proc/version 2>/dev/null; then
        OS="wsl"
      else
        OS="linux"
      fi
      ;;
    MINGW*|MSYS*) OS="mingw" ;;
    *) OS="unknown" ;;
  esac

  case "$uname_m" in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) ARCH="amd64" ;;
  esac
}

install_cloudflared_if_needed() {
  if command -v cloudflared >/dev/null 2>&1; then
    return 0
  fi

  echo -e "${BLUE}cloudflared not found, attempting to install...${NC}"

  case "$OS" in
    mac)
      if command -v brew >/dev/null 2>&1; then
        brew install cloudflared || true
      else
        echo -e "${RED}Homebrew is not installed. Install it from https://brew.sh or install cloudflared manually.${NC}"
        return 1
      fi
      ;;
    linux|wsl)
      if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update -y || true
        sudo apt-get install -y cloudflared || true
      elif command -v dnf >/dev/null 2>&1; then
        sudo dnf install -y cloudflared || true
      elif command -v pacman >/dev/null 2>&1; then
        sudo pacman -Sy --noconfirm cloudflared || true
      fi

      if ! command -v cloudflared >/dev/null 2>&1; then
        echo -e "${BLUE}Attempting direct binary install...${NC}"
        BIN_URL="https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-${ARCH}"
        TMP_BIN="/tmp/cloudflared"
        curl -fsSL "$BIN_URL" -o "$TMP_BIN"
        chmod +x "$TMP_BIN"
        if [ -w /usr/local/bin ]; then
          mv "$TMP_BIN" /usr/local/bin/cloudflared
        else
          sudo mv "$TMP_BIN" /usr/local/bin/cloudflared
        fi
      fi
      ;;
    *)
      echo -e "${RED}OS not automatically supported. Please install cloudflared manually.${NC}"
      return 1
      ;;
  esac

  if ! command -v cloudflared >/dev/null 2>&1; then
    echo -e "${RED}Failed to install cloudflared automatically.${NC}"
    return 1
  fi
}

check_port_available() {
  echo -e "${BLUE}Checking if port ${PORT} is available...${NC}"
  if lsof -i ":${PORT}" >/dev/null 2>&1; then
    echo -e "${RED}ERROR: Port ${PORT} is already in use.${NC}"
    echo -e "${RED}Please stop the service running on port ${PORT} or use a different port.${NC}"
    echo -e "${BLUE}To check what's using the port:${NC}"
    echo -e "${BLUE}  lsof -i :${PORT}${NC}"
    exit 1
  fi
  echo -e "${GREEN}✅ Port ${PORT} is available${NC}"
}

start_server() {
  echo -e "${BLUE}Starting Go server at ${GREEN}http://localhost:${PORT}${NC}..."
  (
    cd "$REPO_ROOT"
    go run ./cmd/server/main.go --addr="$ADDR"
  ) &
  SERVER_PID=$!
  
  # Wait a moment for server to start
  sleep 2
  
  # Verify server is actually running
  if ! lsof -i ":${PORT}" >/dev/null 2>&1; then
    echo -e "${RED}ERROR: Server failed to start on port ${PORT}${NC}"
    echo -e "${RED}Check the server logs above for errors.${NC}"
    exit 1
  fi
  
  echo -e "${GREEN}✅ Server is running on port ${PORT}${NC}"
}

start_tunnel_and_print_url() {
  echo -e "${BLUE}Opening Cloudflare tunnel to ${GREEN}http://localhost:${PORT}${NC}..."
  LOG_FILE="$(mktemp -t cloudflared.XXXXXX.log)"
  cloudflared tunnel --no-autoupdate --url "http://localhost:${PORT}" 2>&1 | tee "$LOG_FILE" &
  TUNNEL_PID=$!

  # Try to extract the URL for up to 30s
  local URL=""
  for _ in $(seq 1 30); do
    sleep 1
    URL=$(grep -oE 'https://[a-zA-Z0-9.-]+\.trycloudflare\.com' "$LOG_FILE" | head -n1 || true)
    if [[ -n "$URL" ]]; then
      break
    fi
  done

  if [[ -z "$URL" ]]; then
    echo -e "${RED}Could not detect the tunnel URL. Check the log: $LOG_FILE${NC}"
  else
    local WSS_URL
    WSS_URL="${URL/https:/wss:}"
    echo
    echo -e "${GREEN}==============================================================${NC}"
    echo -e "${GREEN}    PUBLIC URL (HTTPS): ${BLUE}$URL${NC}"
    echo -e "${GREEN}    SECURE WS (WSS):     ${BLUE}$WSS_URL${NC}"
    echo -e "${GREEN}==============================================================${NC}"
    echo
  fi
}

cleanup() {
  echo -e "${BLUE}Shutting down processes...${NC}"
  [ -n "${TUNNEL_PID:-}" ] && kill "$TUNNEL_PID" >/dev/null 2>&1 || true
  [ -n "${SERVER_PID:-}" ] && kill "$SERVER_PID" >/dev/null 2>&1 || true
}

main() {
  detect_platform

  require_cmd curl || exit 1
  require_cmd go || { echo -e "${RED}Go is not installed. Please install it and retry.${NC}"; exit 1; }
  require_cmd lsof || { echo -e "${RED}lsof is not installed. Please install it and retry.${NC}"; exit 1; }

  install_cloudflared_if_needed || exit 1

  # Check if port is available before starting
  check_port_available

  trap cleanup INT TERM

  start_server
  start_tunnel_and_print_url

  echo -e "${BLUE}Press Ctrl+C to stop server and tunnel.${NC}"
  wait "$SERVER_PID" "$TUNNEL_PID"
}

main "$@"


