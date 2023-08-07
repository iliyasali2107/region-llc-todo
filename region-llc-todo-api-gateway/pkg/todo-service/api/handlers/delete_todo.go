package handlers

import (
	"context"
	"net/http"

	"region-llc-todo-api-gateway/pkg/todo-service/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
func Delete(ctx *gin.Context, c pb.TodoServiceClient) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, err := c.DeleteTodo(context.Background(), &pb.DeleteTodoRequest{
		Id: id,
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

	ctx.JSON(http.StatusNoContent, nil)
}
