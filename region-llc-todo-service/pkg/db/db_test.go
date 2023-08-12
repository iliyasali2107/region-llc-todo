package db_test

import (
	"context"
	"fmt"
	"testing"

	"region-llc-todo-service/pkg/db"
	"region-llc-todo-service/pkg/models"
	"region-llc-todo-service/pkg/utils"

	"github.com/stretchr/testify/require"
)

type Storage interface {
	InsertTodo(ctx context.Context, todo models.Todo) error
	UpdateTodoById(ctx context.Context, todo models.Todo) error
	DeleteTodoById(ctx context.Context, id string) error
	UpdateAsDone(ctx context.Context, id string) error
	GetTodosByFilterDone(ctx context.Context) ([]models.Todo, error)
	GetTodosByFilterActive(ctx context.Context) ([]models.Todo, error)
}

func TestInsertTodo(t *testing.T) {
	duplicateTodo := models.Todo{}

	t.Run("OK", func(t *testing.T) {
		duplicateTodo = insertRandomTodo(t, db.StatusActive)
	})

	t.Run("Duplicate key/field", func(t *testing.T) {
		res, err := TestStorage.InsertTodo(context.Background(), duplicateTodo)
		require.Error(t, err)
		require.Empty(t, res)
		require.Equal(t, "", res)
	})
}

func TestUpdateTodoById(t *testing.T) {
	todo := insertRandomTodo(t, db.StatusActive)

	okArg := models.Todo{
		Id:       todo.Id,
		Title:    utils.RandomString(8),
		ActiveAt: todo.ActiveAt,
		Status:   todo.Status,
	}

	notFoundArg := models.Todo{
		Id:       "invalid id",
		Title:    todo.Title,
		ActiveAt: todo.ActiveAt,
		Status:   todo.Status,
	}

	modifyErrArg := okArg

	testCases := []struct {
		name     string
		arg      models.Todo
		checkRes func(t *testing.T, res int64, err error)
	}{
		{
			name: "OK",
			arg:  okArg,
			checkRes: func(t *testing.T, res int64, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
				checkTodo, err := TestStorage.GetOneTodo(context.Background(), okArg.Id)
				require.Equal(t, checkTodo.Title, okArg.Title)
				require.Equal(t, checkTodo.ActiveAt.Year(), okArg.ActiveAt.Year())
				require.Equal(t, checkTodo.ActiveAt.Month(), okArg.ActiveAt.Month())
				require.Equal(t, checkTodo.ActiveAt.Day(), okArg.ActiveAt.Day())
				require.Equal(t, checkTodo.Title, okArg.Title)
				require.Equal(t, checkTodo.Title, okArg.Title)
				require.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			arg:  notFoundArg,
			checkRes: func(t *testing.T, res int64, err error) {
				require.Error(t, err)
				require.Equal(t, db.ErrNotFound, err)
				require.Equal(t, int64(0), res)
			},
		},
		{
			name: "NotModified",
			arg:  modifyErrArg,
			checkRes: func(t *testing.T, res int64, err error) {
				require.Error(t, err)
				require.Equal(t, db.ErrModify, err)
				require.Equal(t, int64(0), res)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			res, err := TestStorage.UpdateTodoById(context.Background(), tc.arg)
			tc.checkRes(t, res, err)
		})
	}
}

func TestDeleteTodoById(t *testing.T) {
	todo := insertRandomTodo(t, db.StatusActive)
	okArg := todo.Id
	notFoundArg := "not-found-id"

	testCases := []struct {
		name     string
		arg      string
		checkRes func(t *testing.T, res int64, err error)
	}{
		{
			name: "OK",
			arg:  okArg,
			checkRes: func(t *testing.T, res int64, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
				getRes, getErr := TestStorage.GetOneTodo(context.Background(), okArg)
				require.Error(t, getErr)
				require.Empty(t, getRes)
				require.Equal(t, db.ErrNotFound, getErr)
			},
		},
		{
			name: "NotFound",
			arg:  notFoundArg,
			checkRes: func(t *testing.T, res int64, err error) {
				require.Error(t, err)
				require.Equal(t, int64(0), res)
				require.Equal(t, db.ErrNotFound, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			res, err := TestStorage.DeleteTodoById(context.Background(), tc.arg)
			tc.checkRes(t, res, err)
		})
	}
}

func TestUpdateAsDone(t *testing.T) {
	todo := insertRandomTodo(t, db.StatusActive)
	doneTodo := insertRandomTodo(t, db.StatusDone)
	okArg := models.Todo{
		Id:       todo.Id,
		Title:    todo.Title,
		ActiveAt: todo.ActiveAt,
		Status:   db.StatusDone,
	}

	notFoundArg := models.Todo{
		Id:       "invalid id",
		Title:    todo.Title,
		ActiveAt: todo.ActiveAt,
		Status:   todo.Status,
	}

	modifyErrArg := doneTodo

	testCases := []struct {
		name     string
		arg      models.Todo
		checkRes func(t *testing.T, res int64, err error)
	}{
		{
			name: "OK",
			arg:  okArg,
			checkRes: func(t *testing.T, res int64, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
				checkTodo, err := TestStorage.GetOneTodo(context.Background(), okArg.Id)
				require.NoError(t, err)
				require.Equal(t, checkTodo.Status, db.StatusDone)
			},
		},
		{
			name: "NotFound",
			arg:  notFoundArg,
			checkRes: func(t *testing.T, res int64, err error) {
				require.Error(t, err)
				require.Equal(t, db.ErrNotFound, err)
				require.Equal(t, int64(0), res)
			},
		},
		{
			name: "NotModified",
			arg:  modifyErrArg,
			checkRes: func(t *testing.T, res int64, err error) {
				require.Error(t, err)
				require.Equal(t, db.ErrModify, err)
				require.Equal(t, int64(0), res)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			res, err := TestStorage.UpdateAsDone(context.Background(), tc.arg.Id)
			tc.checkRes(t, res, err)
		})
	}
}

func TestGetTodosByFilterDone(t *testing.T) {
	n := 8
	checker := utils.RandomString(n)
	for i := 0; i < n; i++ {
		_, err := TestStorage.InsertTodo(context.Background(), randomTodo(checker, db.StatusDone))
		if err != nil {
			fmt.Print(i, "--")
			fmt.Println(err)
		}
	}

	todos, err := TestStorage.GetTodosByFilterDone(context.Background())
	require.NoError(t, err)

	acc := 0
	for _, todo := range todos {
		if todo.Title[:n] == checker {
			acc++
		}
	}

	require.Equal(t, n, acc)
}

func TestGetTodosByFilterActive(t *testing.T) {
	n := 8
	checker := utils.RandomString(n)
	for i := 0; i < n; i++ {
		TestStorage.InsertTodo(context.Background(), randomTodo(checker, db.StatusActive))
	}

	todos, err := TestStorage.GetTodosByFilterActive(context.Background())

	require.NoError(t, err)

	acc := 0
	for _, todo := range todos {
		if todo.Title[:n] == checker {
			acc++
		}
	}

	require.Equal(t, n, acc)
}

func insertRandomTodo(t *testing.T, status string) models.Todo {
	todo := models.Todo{
		Title:    utils.RandomString(8),
		ActiveAt: utils.RandomDate(),
		Status:   status,
	}

	id, err := TestStorage.InsertTodo(context.Background(), todo)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	todo.Id = id

	return todo
}

func randomTodo(checker string, status string) models.Todo {
	return models.Todo{
		Title:    checker + utils.RandomString(10),
		ActiveAt: utils.RandomDate(),
		Status:   status,
	}
}
