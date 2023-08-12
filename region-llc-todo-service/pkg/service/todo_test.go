package service_test

import (
	"context"
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
	// protoTime := timestamppb.New(utils.RandomDate())

	// okReq := &pb.UpdateTodoRequest{
	// 	Id:       "qwerqwerqwerqwerqwer",
	// 	Title:    utils.RandomString(8),
	// 	ActiveAt: protoTime,
	// }

	// okTodo := models.Todo{
	// 	Title:    okReq.Title,
	// 	ActiveAt: okReq.ActiveAt.AsTime(),
	// 	Status:   db.StatusActive,
	// }

	// testCases := []struct {
	// 	name          string
	// 	req           *pb.UpdateTodoRequest
	// 	buildStubs    func(storage *mocks.MockStorage)
	// 	checkResponse func(t *testing.T, err error)
	// }{{}}
}
