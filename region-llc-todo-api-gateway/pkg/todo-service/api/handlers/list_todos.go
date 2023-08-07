package handlers

import (
	"context"
	"net/http"
	"time"

	"region-llc-todo-api-gateway/pkg/todo-service/pb"
	"region-llc-todo-api-gateway/pkg/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseTodo struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	ActiveAt string `json:"active_at"`
	Status   string `json:"status"`
}

// @Summary Login existing user
// @Tags auth
// @Description login
// @Accept  json
// @Produce  json
// @Param input body CreateRequestBody true "credentials"
// @Success 200 {object} pb.LoginResponse
// @Failure 404
// @Failure 400
// @Failure 500
// @Router /auth/login [post]
func ListTodos(ctx *gin.Context, c pb.TodoServiceClient) {
	statusFilter := ctx.Query("status")
	if statusFilter != "done" {
		if statusFilter != "" && statusFilter != "active" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}

		statusFilter = "active"
	}

	res, err := c.ListTodos(context.Background(), &pb.ListTodosRequest{
		Filter: statusFilter,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "todo is not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something unexpected occured"})
		}
		return
	}

	responseTodos := []ResponseTodo{}
	for _, pbTodo := range res.Todos {
		weekday := pbTodo.ActiveAt.AsTime().Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			pbTodo.Title = "ВЫХОДНОЙ - " + pbTodo.Title
		}

		todo := ResponseTodo{
			Id:       pbTodo.Id,
			Title:    pbTodo.Title,
			ActiveAt: utils.FromProtoToTime(pbTodo.ActiveAt),
			Status:   pbTodo.Status,
		}
		responseTodos = append(responseTodos, todo)
	}

	ctx.JSON(http.StatusOK, responseTodos)
}
