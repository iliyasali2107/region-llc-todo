package utils

import (
	"region-llc-todo-api-gateway/pkg/models"
	"region-llc-todo-api-gateway/pkg/todo-service/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromTodosToProtos(todos []models.Todo) []*pb.Todo {
	protoTodos := []*pb.Todo{}
	for _, todo := range todos {
		// timestampProto := timestamppb.New(todo.ActiveAt)
		protoTodo := &pb.Todo{
			Id:       todo.Id,
			Title:    todo.Title,
			ActiveAt: todo.ActiveAt.Format("2006-01-02"),
			Status:   todo.Status,
		}

		protoTodos = append(protoTodos, protoTodo)
	}
	return protoTodos
}

func FromProtoToTime(date *timestamppb.Timestamp) string {
	timestampTime := date.AsTime()
	return timestampTime.Format("2006-01-02")
}
