# Variabili
PROTO_DIR=api/proto
GEN_DIR=gen/sensor
GATEWAY_BIN=bin/gateway
SENSOR_BIN=bin/sensor
BACKEND_BIN=bin/backend

.PHONY: all gen build clean run-gateway run-sensor run-backend

# 1. Generazione codice da Proto (Go e gRPC)
gen:
	@echo "Pulizia vecchia generazione..."
	rm -rf $(GEN_DIR)
	mkdir -p $(GEN_DIR)
	@echo "Generazione codice dai file .proto..."
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) --go_opt=module=github.com/anna/iot-dual-stack/gen/sensor \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=module=github.com/anna/iot-dual-stack/gen/sensor \
		$(PROTO_DIR)/*.proto
	@echo "Generazione completata in $(GEN_DIR)!"

# 2. Compilazione dei binari
build: gen
	@echo "Compilazione dei componenti..."
	go build -o $(GATEWAY_BIN) cmd/gateway/main.go
	go build -o $(SENSOR_BIN) cmd/sensor/main.go
	go build -o $(BACKEND_BIN) cmd/backend/main.go

# 3. Comandi per avviare i singoli componenti
run-gateway:
	go run cmd/gateway/main.go

run-sensor:
	go run cmd/sensor/main.go

run-backend:
	go run cmd/backend/main.go

# 4. Pulizia
clean:
	rm -rf gen/
	rm -rf bin/