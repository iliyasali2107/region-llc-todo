package routes

import (
	"context"
	"net/http"
	"time"

	"region-llc-todo/pkg/pb"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
)

type CreateRequestBody struct {
	Title    string `json:"title" binding:"required,max=200"`
	ActiveAt string `json:"active_at" binding:"required"`
}

const layout = "2006-01-02"

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
func Create(ctx *gin.Context, c pb.TodoServiceClient) {
	var req CreateRequestBody

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials" + err.Error()})
		return
	}

	parsedTime, err := time.Parse(layout, req.ActiveAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "couldn't parse date"})
		return
	}

	timestampProto, err := ptypes.TimestampProto(parsedTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "couldn't convert"})
		return
	}

	res, err := c.CreateTodo(context.Background(), &pb.CreateTodoRequest{
		Title:    req.Title,
		ActiveAt: timestampProto,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"success": res})
}
