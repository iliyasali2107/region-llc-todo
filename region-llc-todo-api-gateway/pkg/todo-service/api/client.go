package api

import (
	"fmt"

	"region-llc-todo-api-gateway/pkg/config"

	"region-llc-todo-api-gateway/pkg/todo-service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client pb.TodoServiceClient
}

func InitServiceClient(c *config.Config) pb.TodoServiceClient {
	cc, err := grpc.Dial(c.TodoServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewTodoServiceClient(cc)
}
