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
	fmt.Println(date)
	fmt.Println(dateStr)

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
				storage.EXPECT().UpdateTodoById(gomock.Any(), gomock.Any()).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
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
