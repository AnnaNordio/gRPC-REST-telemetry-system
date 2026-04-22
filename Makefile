PROTO_IMAGE = telemetry-proto-builder

# Comando per generare i file usando Docker
docker-gen:
	# Build dell'immagine di generazione (solo la prima volta o se cambia il Dockerfile.proto)
	docker build -f Dockerfile.proto -t $(PROTO_IMAGE) .
	# Esecuzione del container:
	# --rm: elimina il container dopo l'esecuzione
	# -v: monta la cartella corrente dentro il container così i file generati appaiono sul tuo PC
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

# --- BENCHMARK ---
build-bench:
	docker-compose -f docker-compose.benchmark.yaml build

run-bench:
	docker-compose -f docker-compose.benchmark.yaml up --build --abort-on-container-exit

# --- DASHBOARD ---
build-dashboard:
	docker-compose -f docker-compose.yaml build

run-dashboard:
	docker-compose -f docker-compose.yaml up --build --abort-on-container-exit

# --- CLEANUP ---
down:
	docker-compose -f docker-compose.yaml down -v
	docker-compose -f docker-compose.benchmark.yaml down -v

clean:
	docker-compose down -v --rmi all