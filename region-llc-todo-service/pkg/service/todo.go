package service

import (
	"context"
	"time"

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
	convertedTime, err := time.Parse("2006-01-02", req.ActiveAt)
	if err != nil || convertedTime.IsZero() {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid active time")
	}

	if req.Title == "" || len(req.Title) > 200 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid title")
	}

	todo := models.Todo{
		Title:    req.Title,
		ActiveAt: convertedTime,
		Status:   db.StatusActive,
	}

	_, err = ts.Storage.InsertTodo(ctx, todo)
	if err != nil {
		if err == db.ErrDuplicate {
			return &emptypb.Empty{}, status.Errorf(codes.AlreadyExists, "already have this todo")
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something unexpected occured")
	}

	return &emptypb.Empty{}, nil
}

func (ts *todoService) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*emptypb.Empty, error) {
	convertedTime, err := time.Parse("2006-01-02", req.ActiveAt)
	if err != nil || convertedTime.IsZero() {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid active time")
	}

	if req.Title == "" || len(req.Title) > 200 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid title")
	}

	todo := models.Todo{
		Id:       req.Id,
		Title:    req.Title,
		ActiveAt: convertedTime,
	}

	_, err = ts.Storage.UpdateTodoById(ctx, todo)
	if err != nil {
		if err == db.ErrNotFound {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "todo is not found")
		}

		if err == db.ErrModify {
			return &emptypb.Empty{}, status.Errorf(codes.AlreadyExists, "already updated")
		}

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something unexpected occured")
	}

	return &emptypb.Empty{}, nil
}

func (ts *todoService) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*emptypb.Empty, error) {
	id := req.Id

	_, err := ts.Storage.DeleteTodoById(ctx, id)
	if err != nil {
		if err == db.ErrNotFound {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "todo is not found")
		}

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something unexpected occured")
	}

	return &emptypb.Empty{}, nil
}

func (ts *todoService) UpdateAsDone(ctx context.Context, req *pb.UpdateAsDoneRequest) (*emptypb.Empty, error) {
	id := req.Id

	_, err := ts.Storage.UpdateAsDone(ctx, id)
	if err != nil {
		if err == db.ErrNotFound {
			return &emptypb.Empty{}, status.Errorf(codes.NotFound, "todo is not found")
		}

		if err == db.ErrModify {
			return &emptypb.Empty{}, status.Errorf(codes.AlreadyExists, "already updated")
		}

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something unexpected occured")
	}

	return &emptypb.Empty{}, nil
}

func (ts *todoService) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	filterStatus := req.Status

	var err error
	var todos []models.Todo

	if filterStatus == db.StatusDone {
		todos, err = ts.Storage.GetTodosByFilterDone(ctx)
	} else {
		todos, err = ts.Storage.GetTodosByFilterActive(ctx)
	}

	if err != nil {
		if err == db.ErrNotFound {
			return &pb.ListTodosResponse{Todos: []*pb.Todo{}}, nil
		}

		return nil, status.Errorf(codes.Internal, "something unexpected occured")
	}

	return &pb.ListTodosResponse{Todos: utils.FromTodosToProtos(todos)}, nil
}
