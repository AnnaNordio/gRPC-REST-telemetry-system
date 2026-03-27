# Stage 1: Compilazione
FROM golang:1.24-alpine AS builder
# Rimuoviamo apk add per ora, così evitiamo blocchi di rete
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/gateway ./gateway/*.go
RUN go build -o /bin/sensor ./sensor/*.go

# Stage 2: Gateway Runtime (Il nome deve essere 'gateway')
FROM alpine:latest AS gateway
WORKDIR /root/
COPY --from=builder /bin/gateway .
EXPOSE 50051 8080
CMD ["./gateway"]

# Stage 3: Sensor Runtime (Il nome deve essere 'sensor')
FROM alpine:latest AS sensor
WORKDIR /root/
COPY --from=builder /bin/sensor .
CMD ["./sensor"]