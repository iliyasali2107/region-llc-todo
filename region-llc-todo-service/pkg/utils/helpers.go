package utils

import (
	"fmt"
	"math/rand"
	"time"

	"region-llc-todo-service/pkg/models"
	"region-llc-todo-service/pkg/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const DateFormat = "2006-01-02"

func FromTodosToProtos(todos []models.Todo) []*pb.Todo {
	protoTodos := []*pb.Todo{}
	for _, todo := range todos {
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

func RandomDate() time.Time {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	year := rand.Intn(50) + 1970
	month := time.Month(rand.Intn(12) + 1)
	day := rand.Intn(28) + 1
	dateStr := fmt.Sprintf("%d-%d-%d", year, month, day)

	date, _ := time.Parse(DateFormat, dateStr)

	return date
}
