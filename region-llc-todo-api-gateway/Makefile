gen:
	protoc -I=./pkg/todo-service/pb/ -I=./ --go_out=./ --go-grpc_out=./ --grpc-gateway_out=./ ./pkg/todo-service/pb/*.proto

client:
	go run ./cmd/main.go


mongo:
	docker run --name test-mongodb -d -p 27017:27017 mongo:latest

swagger:
	docker run -p 80:8080 \
    -e SWAGGER_JSON=./gen/openapiv2/pkg/todo-service/pb/todo.swagger.json \
    swaggerapi/swagger-ui