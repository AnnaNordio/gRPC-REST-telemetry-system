# --- STAGE 1: Builder (Go + Node.js per il build) ---
FROM golang:1.24-alpine AS builder

# Installiamo nodejs e npm solo qui per buildare la dashboard
RUN apk add --no-cache nodejs npm

WORKDIR /app

# 1. Gestione dipendenze Go
COPY go.mod go.sum ./
RUN go mod download

# 2. Copia tutto il codice sorgente
COPY . .

# 3. Build della Dashboard (Frontend)
# Questo genera la cartella /app/dashboard/dist
RUN cd dashboard && npm install && npm run build

# 4. Compilazione dei binari Go
RUN go build -o /bin/gateway ./gateway/*.go
RUN go build -o /bin/sensor ./sensor/*.go


# --- STAGE 2: Gateway Runtime (Immagine Leggera) ---
FROM alpine:latest AS gateway
WORKDIR /root/

# Copiamo il binario dal builder
COPY --from=builder /bin/gateway .

# Copiamo solo i file statici compilati della dashboard
# NOTA: Go dovrà puntare a questa cartella (es: http.Dir("dashboard"))
COPY --from=builder /app/dashboard/dist ./dashboard

# Espone le porte necessarie
EXPOSE 50051 8080

# Avviamo solo Go. Niente Node.js o 'serve' a runtime.
CMD ["./gateway"]


# --- STAGE 3: Sensor Runtime ---
FROM alpine:latest AS sensor
WORKDIR /root/

# Copiamo il binario dal builder
COPY --from=builder /bin/sensor .

# Il sensore non ha bisogno di porte o dashboard
CMD ["./sensor"]