#!/bin/bash

set -u

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

BACKEND_PORT="${BACKEND_PORT:-8080}"
FRONTEND_PORT="${FRONTEND_PORT:-3000}"

BACKEND_PID=""
FRONTEND_PID=""

mkdir -p logs
PID_FILE="$PROJECT_ROOT/logs/pids.txt"
> "$PID_FILE"

cleanup() {
    echo -e "\n${YELLOW}Stopping all services...${NC}"
    [ -n "$BACKEND_PID" ] && kill "$BACKEND_PID" 2>/dev/null
    [ -n "$FRONTEND_PID" ] && kill "$FRONTEND_PID" 2>/dev/null
    lsof -ti :$BACKEND_PORT | xargs kill -9 2>/dev/null || true
    lsof -ti :$FRONTEND_PORT | xargs kill -9 2>/dev/null || true
    echo -e "${GREEN}All services stopped.${NC}"
    exit 0
}

trap cleanup SIGINT SIGTERM

echo -e "${BLUE}Starting JobScraper...${NC}\n"

echo -e "${BLUE}[1/2] Starting Backend (air)...${NC}"
air > "$PROJECT_ROOT/logs/api.log" 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID >> "$PID_FILE"
echo -e "${GREEN}✓ Backend started (PID: $BACKEND_PID)${NC}"
echo -e "  Log: logs/api.log"

sleep 2

echo -e "\n${BLUE}[2/2] Starting Frontend...${NC}"
cd "$PROJECT_ROOT/frontend"
npm run dev > "$PROJECT_ROOT/logs/frontend.log" 2>&1 &
FRONTEND_PID=$!
echo $FRONTEND_PID >> "$PID_FILE"
echo -e "${GREEN}✓ Frontend started (PID: $FRONTEND_PID)${NC}"
echo -e "  Log: logs/frontend.log"

cd "$PROJECT_ROOT"

echo -e "\n${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}  JobScraper is running!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

echo -e "${BLUE}URLs:${NC}"
echo -e "  Frontend: ${YELLOW}http://localhost:${FRONTEND_PORT}${NC}"
echo -e "  Backend:  ${YELLOW}http://localhost:${BACKEND_PORT}${NC}"

echo -e "\n${BLUE}Logs:${NC}"
echo -e "  Backend:  ${YELLOW}tail -f logs/api.log${NC}"
echo -e "  Frontend: ${YELLOW}tail -f logs/frontend.log${NC}"

echo -e "\n${RED}Press Ctrl+C to stop${NC}\n"

while true; do
    sleep 2
    if [ -n "$BACKEND_PID" ] && ! kill -0 "$BACKEND_PID" 2>/dev/null; then
        echo -e "${RED}✗ Backend crashed. Stopping...${NC}"
        cleanup
    fi
    if [ -n "$FRONTEND_PID" ] && ! kill -0 "$FRONTEND_PID" 2>/dev/null; then
        echo -e "${RED}✗ Frontend crashed. Stopping...${NC}"
        cleanup
    fi
done
