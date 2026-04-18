# --- Stage 1: Base Builder (Dipendenze comuni) ---
FROM golang:1.24-alpine AS base-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# --- Stage 2: Compilazione specifica Gateway ---
FROM base-builder AS gateway-builder
RUN go build -o /bin/gateway ./gateway/*.go

# --- Stage 3: Compilazione specifica Sensor ---
FROM base-builder AS sensor-builder
RUN go build -o /bin/sensor ./sensor/*.go

# --- Stage 4: Gateway Runtime ---
FROM alpine:latest AS gateway
RUN apk add --no-cache nodejs npm && npm install -g serve
WORKDIR /root/
# Copia solo dal suo builder specifico
COPY --from=gateway-builder /bin/gateway .
COPY dashboard/ ./dashboard/
RUN cd ./dashboard && npm install && npm run build
EXPOSE 50051 8080 3000
CMD ["sh", "-c", "serve -s dashboard/dist -l 3000 & ./gateway"]

# --- Stage 5: Sensor Runtime ---
FROM alpine:latest AS sensor
WORKDIR /root/
# Copia solo dal suo builder specifico
COPY --from=sensor-builder /bin/sensor .
CMD ["./sensor"]