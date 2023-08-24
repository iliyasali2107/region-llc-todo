package utils

import (
	"fmt"
	"math/rand"
	"strconv"
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
	month := strconv.Itoa(rand.Intn(12) + 1)
	day := strconv.Itoa(rand.Intn(28) + 1)

	if len(day) == 1 {
		day = "0" + day
	}

	if len(month) == 1 {
		month = "0" + month
	}

	dateStr := fmt.Sprintf("%d-%s-%s", year, month, day)
	date, _ := time.Parse(DateFormat, dateStr)

	return date
}

func RandomDateStr() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	year := rand.Intn(50) + 1970
	month := strconv.Itoa(rand.Intn(12) + 1)
	day := strconv.Itoa(rand.Intn(28) + 1)

	if len(day) == 1 {
		day = "0" + day
	}

	if len(month) == 1 {
		month = "0" + month
	}

	dateStr := fmt.Sprintf("%d-%s-%s", year, month, day)

	return dateStr
}

func RandomTodo(checker string, status string) models.Todo {
	return models.Todo{
		Title:    checker + RandomString(10),
		ActiveAt: RandomDate(),
		Status:   status,
	}
}
