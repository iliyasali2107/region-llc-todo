package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"region-llc-todo-service/pkg/db"
	"region-llc-todo-service/pkg/mocks"
	"region-llc-todo-service/pkg/models"
	"region-llc-todo-service/pkg/pb"
	"region-llc-todo-service/pkg/service"
	"region-llc-todo-service/pkg/utils"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateTodo(t *testing.T) {
	date := utils.RandomDate()
	dateStr := date.Format(utils.DateFormat)

	okReq := &pb.CreateTodoRequest{
		Title:    utils.RandomString(8),
		ActiveAt: dateStr,
	}

	invalidReq := &pb.CreateTodoRequest{
		Title:    utils.RandomString(8),
		ActiveAt: time.Time{}.Format(utils.DateFormat),
	}

	okTodo := models.Todo{
		Title:    okReq.Title,
		ActiveAt: date,
		Status:   db.StatusActive,
	}

	testCases := []struct {
		name          string
		req           *pb.CreateTodoRequest
		buildStubs    func(storage *mocks.MockStorage)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().InsertTodo(gomock.Any(), gomock.Eq(okTodo)).Times(1).Return("qwer", nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Duplicate",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().InsertTodo(gomock.Any(), gomock.Eq(okTodo)).Times(1).Return("", db.ErrDuplicate)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
		{
			name: "InvalidArgument",
			req:  invalidReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().InsertTodo(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "Internal",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().InsertTodo(gomock.Any(), gomock.Any()).Times(1).Return("", fmt.Errorf("internal"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)

			tc.buildStubs(storage)

			serv := service.NewTodoService(storage)

			_, err := serv.CreateTodo(context.Background(), tc.req)
			tc.checkResponse(t, err)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	todo := insertRandomTodo(t, db.StatusActive)

	okReq := &pb.UpdateTodoRequest{
		Id:       todo.Id,
		Title:    todo.Title,
		ActiveAt: todo.ActiveAt.Format(utils.DateFormat),
	}

	date, err := time.Parse("2006-01-02", okReq.ActiveAt)
	require.NoError(t, err)

	okTodo := models.Todo{
		Id:       todo.Id,
		Title:    todo.Title,
		ActiveAt: date,
	}

	notFoundReq := &pb.UpdateTodoRequest{
		Id:       "invalid-id",
		Title:    todo.Title,
		ActiveAt: todo.ActiveAt.Format(utils.DateFormat),
	}

	notFoundTodo := models.Todo{
		Id:       "invalid-id",
		Title:    okReq.Title,
		ActiveAt: date,
	}

	testCases := []struct {
		name          string
		req           *pb.UpdateTodoRequest
		buildStubs    func(storage *mocks.MockStorage)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().UpdateTodoById(gomock.Any(), gomock.Eq(okTodo)).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},

		{
			name: "NotFound",
			req:  notFoundReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().UpdateTodoById(gomock.Any(), gomock.Eq(notFoundTodo)).Times(1).Return(int64(0), db.ErrNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "AlreadyExists",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().UpdateTodoById(gomock.Any(), gomock.Eq(okTodo)).Times(1).Return(int64(0), db.ErrModify)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
		{
			name: "Internal",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().UpdateTodoById(gomock.Any(), gomock.Any()).Times(1).Return(int64(0), fmt.Errorf("internal"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)

			tc.buildStubs(storage)

			serv := service.NewTodoService(storage)

			_, err := serv.UpdateTodo(context.Background(), tc.req)
			tc.checkResponse(t, err)
		})
	}
}

func TestUpdateAsDone(t *testing.T) {
	todo := insertRandomTodo(t, db.StatusActive)

	doneTodo := insertRandomTodo(t, db.StatusDone)

	okReq := &pb.UpdateAsDoneRequest{
		Id: todo.Id,
	}

	notFoundReq := &pb.UpdateAsDoneRequest{
		Id: "invalid-id",
	}

	alreadyExistsReq := &pb.UpdateAsDoneRequest{
		Id: doneTodo.Id,
	}

	testCases := []struct {
		name          string
		req           *pb.UpdateAsDoneRequest
		buildStubs    func(storage *mocks.MockStorage)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().UpdateAsDone(gomock.Any(), gomock.Eq(okReq.Id)).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},

		{
			name: "NotFound",
			req:  notFoundReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().UpdateAsDone(gomock.Any(), gomock.Eq(notFoundReq.Id)).Times(1).Return(int64(0), db.ErrNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "AlreadyExists",
			req:  alreadyExistsReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().UpdateAsDone(gomock.Any(), gomock.Eq(alreadyExistsReq.Id)).Times(1).Return(int64(0), db.ErrModify)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
		{
			name: "Internal",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().UpdateAsDone(gomock.Any(), gomock.Any()).Times(1).Return(int64(0), fmt.Errorf("internal"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)

			tc.buildStubs(storage)

			serv := service.NewTodoService(storage)

			_, err := serv.UpdateAsDone(context.Background(), tc.req)
			tc.checkResponse(t, err)
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	todo := insertRandomTodo(t, db.StatusActive)

	okReq := &pb.DeleteTodoRequest{
		Id: todo.Id,
	}

	notFoundReq := &pb.DeleteTodoRequest{
		Id: "invalid-id",
	}

	testCases := []struct {
		name          string
		req           *pb.DeleteTodoRequest
		buildStubs    func(storage *mocks.MockStorage)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().DeleteTodoById(gomock.Any(), gomock.Eq(okReq.Id)).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},

		{
			name: "NotFound",
			req:  notFoundReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().DeleteTodoById(gomock.Any(), gomock.Eq(notFoundReq.Id)).Times(1).Return(int64(0), db.ErrNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},

		{
			name: "Internal",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().DeleteTodoById(gomock.Any(), gomock.Any()).Times(1).Return(int64(0), fmt.Errorf("internal"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)

			tc.buildStubs(storage)

			serv := service.NewTodoService(storage)

			_, err := serv.DeleteTodo(context.Background(), tc.req)
			tc.checkResponse(t, err)
		})
	}
}

func TestListTodos(t *testing.T) {
	n := 8
	checker := utils.RandomString(n)
	for i := 0; i < n; i++ {
		_, err := TestStorage.InsertTodo(context.Background(), utils.RandomTodo(checker, db.StatusActive))
		require.NoError(t, err)
	}

	okReq := &pb.ListTodosRequest{Status: db.StatusActive}

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		serv := service.NewTodoService(TestStorage)

		todos, err := serv.ListTodos(context.Background(), &pb.ListTodosRequest{Status: db.StatusActive})
		require.NoError(t, err)
		for _, todo := range todos.Todos {
			require.NotNil(t, todo)
			require.Equal(t, db.StatusActive, todo.Status)
		}
	})

	testCases := []struct {
		name          string
		req           *pb.ListTodosRequest
		buildStubs    func(storage *mocks.MockStorage)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "NotFound",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().GetTodosByFilterActive(gomock.Any()).Times(1).Return([]models.Todo{}, db.ErrNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},

		{
			name: "Internal",
			req:  okReq,
			buildStubs: func(storage *mocks.MockStorage) {
				storage.EXPECT().GetTodosByFilterActive(gomock.Any()).Times(1).Return([]models.Todo{}, fmt.Errorf("internal"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				st, _ := status.FromError(err)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)

			tc.buildStubs(storage)

			serv := service.NewTodoService(storage)

			_, err := serv.ListTodos(context.Background(), tc.req)
			tc.checkResponse(t, err)
		})
	}
}

func insertRandomTodo(t *testing.T, status string) models.Todo {
	date, err := time.Parse(utils.DateFormat, utils.RandomDateStr())
	require.NoError(t, err)
	todo := models.Todo{
		Title:    utils.RandomString(8),
		ActiveAt: date,
		Status:   status,
	}

	id, err := TestStorage.InsertTodo(context.Background(), todo)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	todo.Id = id

	return todo
}
