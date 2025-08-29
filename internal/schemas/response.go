package schemas

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func (r *Response[T]) RespondJSON(ctx *gin.Context, code int) {
	if code == http.StatusOK {
		r.Success = true
		r.Message = "success"
		// logrus.Debugf("响应成功: %+v", r.Data)
	} else {
		logrus.Warning(r.Message)
	}
	ctx.JSON(code, r)
}

func (r *Response[T]) RespondSuccessJSON(ctx *gin.Context, data T) {
	r.Data = data
	r.RespondJSON(ctx, http.StatusOK)
}
