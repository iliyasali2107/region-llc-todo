package utils

import (
	"math/rand"
	"time"

	"region-llc-todo-service/pkg/models"
	"region-llc-todo-service/pkg/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromTodosToProtos(todos []models.Todo) []*pb.Todo {
	protoTodos := []*pb.Todo{}
	for _, todo := range todos {
		timestampProto := timestamppb.New(todo.ActiveAt)
		protoTodo := &pb.Todo{
			Id:       todo.Id,
			Title:    todo.Title,
			ActiveAt: timestampProto,
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
	hour := rand.Intn(24)
	minute := rand.Intn(60)
	second := rand.Intn(60)
	nanosecond := rand.Intn(1000000000)

	randomTime := time.Date(year, month, day, hour, minute, second, nanosecond, time.UTC)

	return randomTime
}
