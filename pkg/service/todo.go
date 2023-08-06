package service

import (
	"context"

	"region-llc-todo/pkg/db"
	"region-llc-todo/pkg/models"
	"region-llc-todo/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	S db.Storage
	pb.UnimplementedTodoServiceServer
}

type Service interface {
	pb.TodoServiceServer
	CreateTodo(context.Context, *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error)
}

func NewService(s db.Storage) Service {
	return &service{
		S: s,
	}
}

func (s *service) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	convertedTime := req.ActiveAt.AsTime()
	if convertedTime.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid active time")
	}

	todo := models.Todo{
		Title:    req.Title,
		ActiveAt: convertedTime,
	}

	err := s.S.InsertTodo(ctx, todo)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTodoResponse{Id: 200}, nil
}
