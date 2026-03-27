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

run-gateway:
	go run ./gateway/*.go 

run-sensor:
	go run ./sensor/*.go