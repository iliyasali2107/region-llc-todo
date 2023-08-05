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

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to config", err)
	}

	storage := db.Init(cfg)

	lis, err := net.Listen("tcp", cfg.ClientPort)
	if err != nil {
		log.Fatalln("failed to listen", err)
	}

	srv := service.NewService(storage)

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, srv)
	fmt.Println("rumming..........")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("failed to server", err)
	}
}
