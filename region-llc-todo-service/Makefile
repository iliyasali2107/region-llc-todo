test:
	go test -v -cover ./...

mongo:
	docker run --name test-mongodb -d -p 27017:27017 mongo:latest

mock_storage:
	mockgen -destination=pkg/mocks/mock_storage.go --build_flags=--mod=mod -package=mocks region-llc-todo-service/pkg/db Storage

server: 
	go run ./cmd/main.go