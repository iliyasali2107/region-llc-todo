package handlers

import (
	"context"
	"net/http"
	"time"

	"region-llc-todo-api-gateway/pkg/todo-service/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UpdateRequestBody struct {
	Title    string `json:"title" binding:"required,max=200"`
	ActiveAt string `json:"active_at" binding:"required"`
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
func Update(ctx *gin.Context, c pb.TodoServiceClient) {
	var req UpdateRequestBody

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials" + err.Error()})
		return
	}

	parsedTime, err := time.Parse(layout, req.ActiveAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "date validation failed"})
		return
	}

	timestampProto := timestamppb.New(parsedTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "convertation failed"})
		return
	}

	_, err = c.UpdateTodo(context.Background(), &pb.UpdateTodoRequest{
		Id:       id,
		Title:    req.Title,
		ActiveAt: timestampProto,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "todo is not found"})
		} else if st.Code() == codes.AlreadyExists {
			ctx.JSON(http.StatusConflict, gin.H{"error": "already updated"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something unexpected occured"})
		}
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
