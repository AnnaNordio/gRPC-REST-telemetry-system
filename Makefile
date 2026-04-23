.PHONY: docker-gen gen build-bench run-bench build-dashboard run-dashboard down clean
PROTO_IMAGE = telemetry-proto-builder

# Comando per generare i file usando Docker
docker-gen:
	# Build dell'immagine di generazione (solo la prima volta o se cambia il Dockerfile.proto)
	docker build -f Dockerfile.proto -t $(PROTO_IMAGE) .
	docker run --rm -v $(shell pwd):/app $(PROTO_IMAGE)

gen:
	# 1. Parte Go (Assicurati che la cartella proto esista)
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/telemetry.proto

	# 2. Generazione JS (Corretto il path -I e il puntamento al file)
	mkdir -p dashboard/proto-pkg
	protoc -I=proto telemetry.proto \
		--js_out=import_style=commonjs,binary:./dashboard/proto-pkg \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./dashboard/proto-pkg
	
	# 3. Build del bundle per Vite
	# Usiamo && per assicurarci che se npm install fallisce, il build non parta
	cd dashboard/proto-pkg && npm install && npm run build

COMPOSE_CMD := $(shell docker compose version >/dev/null 2>&1 && echo "docker compose" || echo "docker-compose")

# --- BENCHMARK ---
build-bench:
	$(COMPOSE_CMD)-f docker-compose.benchmark.yaml build

run-bench:
	$(COMPOSE_CMD) -f docker-compose.benchmark.yaml up --build --abort-on-container-exit

# --- DASHBOARD ---
build-dashboard:
	$(COMPOSE_CMD) -f docker-compose.yaml build

run-dashboard:
	$(COMPOSE_CMD) -f docker-compose.yaml up --build --abort-on-container-exit

# --- CLEANUP ---
down:
	$(COMPOSE_CMD) -f docker-compose.yaml down -v
	$(COMPOSE_CMD) -f docker-compose.benchmark.yaml down -v

clean:
	$(COMPOSE_CMD) -f docker-compose.yaml -f docker-compose.benchmark.yaml down -v --rmi all