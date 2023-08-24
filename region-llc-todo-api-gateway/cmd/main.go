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

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	rmux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterTodoServiceHandlerFromEndpoint(ctx, rmux, cfg.TodoServicePort, opts)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	fs := http.FileServer(http.Dir("./pkg/docs"))

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	fmt.Printf("Client is on: %s", cfg.Port)
	err = http.ListenAndServe(cfg.Port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
