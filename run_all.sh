#!/bin/bash

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
PURPLE='\033[0;35m'
NC='\033[0m'

cleanup() {
    trap - EXIT
    echo -e "\n${RED}🛑 Chiusura di tutti i processi in corso...${NC}"
    pkill -P $$
    exit
}

trap cleanup SIGINT SIGTERM

echo -e "${BLUE}🧹 Pulizia processi residui...${NC}"
pkill -f "iot-dual-stack" || true

# 1. Avvio Infrastruttura
echo -e "${GREEN}1. 🚀 Avvio BACKEND...${NC}"
go run cmd/backend/main.go > backend.log 2>&1 &
sleep 2

echo -e "${GREEN}2. 🌉 Avvio GATEWAY...${NC}"
go run cmd/gateway/main.go > gateway.log 2>&1 &
sleep 2

echo -e "${PURPLE}🧪 AVVIO SCENARI DI TEST...${NC}"
echo -e "--------------------------------------------------"

# --- SCENARIO A: 1 Sensore REST (Baseline) ---
echo -e "${BLUE}📡 Sensore 1: REST (Freq 1s)${NC}"
go run cmd/sensor/main.go -id="SENSOR-REST" -mode="rest" -freq="1s" &

# --- SCENARIO B: 1 Sensore gRPC (Confronto diretto) ---
echo -e "${BLUE}📡 Sensore 2: gRPC (Freq 1s)${NC}"
go run cmd/sensor/main.go -id="SENSOR-GRPC" -mode="grpc" -freq="1s" &

# --- SCENARIO C: Stress Test (Payload Pesante) ---
# Lanciamo un sensore che invia 50KB di dati extra via REST per vedere il lag
echo -e "${RED}🔥 Sensore 3: REST STRESS (Payload 50KB)${NC}"
go run cmd/sensor/main.go -id="STRESS-REST" -mode="rest" -freq="2s" -payload=50000 &

# --- SCENARIO D: gRPC Streaming ---
echo -e "${GREEN}🌊 Sensore 4: gRPC STREAMING${NC}"
go run cmd/sensor/main.go -id="STREAM-SENSOR" -mode="stream" -freq="5s" &

echo -e "--------------------------------------------------"
echo -e "${BLUE}📝 Log attivi in backend.log e gateway.log${NC}"
echo -e "${BLUE}📊 Dashboard disponibile su http://localhost:8080${NC}"
echo -e "${RED}Premi CTRL+C per terminare il test e pulire tutto.${NC}"

# Aspetta che tutti i processi figli terminino (o CTRL+C)
wait