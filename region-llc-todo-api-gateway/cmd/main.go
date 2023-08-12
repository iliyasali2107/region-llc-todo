package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"region-llc-todo-api-gateway/pkg/config"
	"region-llc-todo-api-gateway/pkg/todo-service/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// import (
// 	"fmt"
// 	"log"
// 	"net"

// 	"region-llc-todo/pkg/config"
// 	"region-llc-todo/pkg/db"
// 	"region-llc-todo/pkg/pb"
// 	"region-llc-todo/pkg/service"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )

// func main() {
// 	c, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatalln("Failed to config", err)
// 	}

// 	storage := db.Init(c.DBUrl)

// 	lis, err := net.Listen("tcp", c.Port)
// 	if err != nil {
// 		log.Fatalln("Failed to listening")
// 	}

// 	fmt.Println("Url service is on: ", c.Port)

// 	srv := service.NewService(storage)

// 	grpcServer := grpc.NewServer()

// 	pb.RegisterUrlServiceServer(grpcServer, srv)
// 	reflection.Register(grpcServer)

// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalln("Failed to serve: ", err)
// 	}
// }

// gin
// func main() {
// 	cfg, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatalln("failed to configdock", err)
// 	}

// 	r := gin.Default()

// 	todo_svc.RegisterRoutes(r, &cfg)

// 	fmt.Println("api-gateway is running on: ", cfg.Port)

// 	r.Run(cfg.Port)
// }

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	fmt.Println(cfg.TodoServicePort)
	err = pb.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, cfg.TodoServicePort, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Client is on: %s", cfg.Port)

	return http.ListenAndServe(cfg.Port, mux)
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
