package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"region-llc-todo-service/pkg/config"
	"region-llc-todo-service/pkg/db"
	"region-llc-todo-service/pkg/pb"
	"region-llc-todo-service/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	// загружаем переменные среды
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to config", err)
	}

	// инициализируем сторадж
	storage := db.Init(cfg)

	// и закрываем при выходе из программы
	defer func() {
		if err := storage.Close(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalln("Failed to listening")
	}

	fmt.Println("Todo service is on: ", cfg.Port)

	srv := service.NewTodoService(storage)

	
	grpcServer := grpc.NewServer()

	pb.RegisterTodoServiceServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve: ", err)
	}
}
