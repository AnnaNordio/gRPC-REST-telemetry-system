# --- Stage 1: Build Go Binaries ---
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/gateway ./gateway/*.go
RUN go build -o /bin/sensor ./sensor/*.go

# --- Stage 2: Gateway + Dashboard Runtime ---
FROM alpine:latest AS gateway
# Installiamo Node.js e il pacchetto 'serve' per gestire il build sulla 3000
RUN apk add --no-cache nodejs npm
RUN npm install -g serve

WORKDIR /app
COPY --from=builder /bin/gateway .

# Copiamo i sorgenti della dashboard ed eseguiamo il build
# (Nota: facciamo il build qui per assicurarci che l'ambiente Alpine sia coerente)
COPY dashboard/ ./dashboard/
RUN cd ./dashboard && npm install && npm run build

# Espone le porte: 50051 (gRPC), 8080 (Gateway API), 3000 (Dashboard)
EXPOSE 50051 8080 3000

# Script per avviare il server statico sulla 3000 e il gateway
CMD ["sh", "-c", "echo 'Dashboard attiva su localhost:3000' && serve -s dashboard/dist -l 3000 & ./gateway"]

# --- Stage 3: Sensor Runtime ---
FROM alpine:latest AS sensor
WORKDIR /root/
COPY --from=builder /bin/sensor .
CMD ["./sensor"]