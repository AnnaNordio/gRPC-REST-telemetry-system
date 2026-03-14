#!/bin/bash

# Colori per il terminale
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Funzione per pulire tutto alla chiusura
cleanup() {
    echo -e "\n${RED}🛑 Chiusura di tutti i processi in corso...${NC}"
    # Uccide tutti i processi avviati da questo script (il gruppo di processi)
    kill 0
    exit
}

# Associa il segnale di interruzione (CTRL+C) alla funzione cleanup
trap cleanup SIGINT SIGTERM EXIT

echo -e "${BLUE}🧹 Pulizia processi residui...${NC}"
# Usiamo pkill in modo più mirato sui binari se possibile, 
# ma per ora manteniamo la tua logica go run
pkill -f "go run cmd/backend/main.go" || true
pkill -f "go run cmd/gateway/main.go" || true

echo -e "${GREEN}1. 🚀 Avvio BACKEND...${NC}"
go run cmd/backend/main.go > backend.log 2>&1 &
BACKEND_PID=$!
sleep 2

echo -e "${GREEN}2. 🌉 Avvio GATEWAY...${NC}"
go run cmd/gateway/main.go > gateway.log 2>&1 &
GATEWAY_PID=$!
sleep 2

echo -e "${GREEN}3. 📡 Avvio SIMULATORE SENSORE...${NC}"
echo -e "${BLUE}📝 Log: backend.log, gateway.log${NC}"
echo -e "--------------------------------------------------"

# Avviamo il sensore in primo piano
go run cmd/sensor/main.go

# Lo script rimarrà qui finché il sensore non finisce o premi CTRL+C