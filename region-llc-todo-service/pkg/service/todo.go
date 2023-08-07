package service

import (
	"context"
	"region-llc-todo-service/pkg/db"
	"region-llc-todo-service/pkg/models"
	"region-llc-todo-service/pkg/pb"
	"region-llc-todo-service/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type todoService struct {
	Storage db.Storage
	pb.UnimplementedTodoServiceServer
}

type TodoService interface {
	pb.TodoServiceServer
	CreateTodo(context.Context, *pb.CreateTodoRequest) (*emptypb.Empty, error)
	UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*emptypb.Empty, error)
	DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*emptypb.Empty, error)
	UpdateAsDone(ctx context.Context, req *pb.UpdateAsDoneRequest) (*emptypb.Empty, error)
	ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error)
}

func NewTodoService(s db.Storage) TodoService {
	return &todoService{
		Storage: s,
	}
}

func (ts *todoService) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*emptypb.Empty, error) {
	convertedTime := req.ActiveAt.AsTime()
	if convertedTime.IsZero() {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid active time")
	}

	todo := models.Todo{
		Title:    req.Title,
		ActiveAt: convertedTime,
		Status:   db.StatusActive,
	}

	err := ts.Storage.InsertTodo(ctx, todo)
	if err != nil {
		if err == db.ErrDuplicate {
			return &emptypb.Empty{}, status.Errorf(codes.AlreadyExists, "already have this todo")
		}
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (ts *todoService) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*emptypb.Empty, error) {
	convertedTime := req.ActiveAt.AsTime()
	if convertedTime.IsZero() {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid active time")
	}

	todo := models.Todo{
		Id:       req.Id,
		Title:    req.Title,
		ActiveAt: convertedTime,
	}

	err := ts.Storage.UpdateTodoById(ctx, todo)
	if err != nil {
		if err == db.ErrNotFound {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "todo is not found")
		}

		if err == db.ErrModify {
			return &emptypb.Empty{}, status.Errorf(codes.AlreadyExists, "already updated")
		}

		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (ts *todoService) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*emptypb.Empty, error) {
	id := req.Id

	err := ts.Storage.DeleteTodoById(ctx, id)
	if err != nil {
		if err == db.ErrNotFound {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "todo is not found")
		}

		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (ts *todoService) UpdateAsDone(ctx context.Context, req *pb.UpdateAsDoneRequest) (*emptypb.Empty, error) {
	id := req.Id

	err := ts.Storage.UpdateAsDone(ctx, id)
	if err != nil {
		if err == db.ErrNotFound {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "todo is not found")
		}

		if err == db.ErrModify {
			return &emptypb.Empty{}, status.Errorf(codes.AlreadyExists, "already updated")
		}

		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (ts *todoService) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	filterValue := req.Filter

	var err error
	var todos []models.Todo

	if filterValue == db.StatusActive {
		todos, err = ts.Storage.GetTodosByFilterActive(ctx)
	} else {
		todos, err = ts.Storage.GetTodosByFilterDone(ctx)
	}

	if err != nil {
		if err == db.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "todo is not found")
		}

		return nil, err
	}

	return &pb.ListTodosResponse{Todos: utils.FromTodosToProtos(todos)}, nil
}