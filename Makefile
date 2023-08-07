postgres:
	docker run -d --name my-postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=url_redirector -p 5432:5432 postgres:latest

server:
	go run ./cmd/server/main.go

client:
	go run ./cmd/client/main.go

gen:
	protoc -I=./pkg/pb --go_out=./ --go-grpc_out=./ ./pkg/pb/*.proto

mock_storage:
	mockgen -destination=pkg/mocks/mock_storage.go --build_flags=--mod=mod -package=mocks region-llc-todo/pkg/db Storage

test:
	go test -v -cover ./...

mongo:
	docker run --name my-mongodb -d -p 27017:27017 mongo:latest


up:
	docker-compose up

build: 
	docker-compose build