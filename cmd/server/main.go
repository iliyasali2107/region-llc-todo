package main

import (
	"fmt"
	"log"
	"net"
	"region-llc-todo/pkg/config"
	"region-llc-todo/pkg/db"
	"region-llc-todo/pkg/pb"
	"region-llc-todo/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to config", err)
	}

	storage := db.Init(cfg)

	lis, err := net.Listen("tcp", cfg.ServicePort)
	if err != nil {
		log.Fatalln("Failed to listening")
	}

	fmt.Println("Url service is on: ", cfg.ServicePort)

	srv := service.NewService(storage)

	grpcServer := grpc.NewServer()

	pb.RegisterTodoServiceServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve: ", err)
	}
}
