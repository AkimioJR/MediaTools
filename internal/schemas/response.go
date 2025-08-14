package schemas

import (
	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func (r *Response[T]) RespondJSON(ctx *gin.Context, code int) {
	ctx.JSON(code, r)
}
