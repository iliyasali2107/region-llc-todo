package main

import (
	"fmt"
	"log"

	"region-llc-todo/cmd"
	"region-llc-todo/pkg/config"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	cmd.RegisterRoutes(r, &cfg)

	fmt.Println("todo service is on: ", cfg.Port)

	r.Run(cfg.ClientPort)

}