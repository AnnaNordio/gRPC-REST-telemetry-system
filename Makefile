gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/telemetry.proto

run-gateway:
	go run gateway/main.go

run-sensor:
	go run sensor/main.go