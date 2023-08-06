package cmd

import (
	"fmt"

	"region-llc-todo/pkg/config"
	"region-llc-todo/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client pb.TodoServiceClient
}

func InitServiceClient(c *config.Config) pb.TodoServiceClient {
	cc, err := grpc.Dial(c.ServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewTodoServiceClient(cc)
}
