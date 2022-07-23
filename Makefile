.SILENT:

proto:
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. \
	--grpc-gateway_opt generate_unbound_methods=true --openapiv2_out . \
	./internal/proto/news.proto

build:
	go mod tidy && go build -o app.out ./cmd/app

run: build
	./app.out

client:
	go run ./cmd/client/client.go

service:
	docker-compose up

app:
	docker-compose up app

app_build:
	docker-compose up --build app

db:
	docker-compose up db


clear: 
	./scripts/clear_docker.sh