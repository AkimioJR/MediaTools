package task

import (
	"MediaTools/internal/controller/task_controller"
	"MediaTools/internal/pkg/task"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Router /task/transfer [get]
// @Summary 查询转移任务列表
// @Description 查询转移任务列表
// @Tags 任务管理
// @Produces json
func GetAllTransferTasks(ctx *gin.Context) {
	var resp schemas.Response[[]*task.Task]

	tasks := make([]*task.Task, 0)
	for task := range task_controller.IterTransferTasks {
		resp.Data = append(resp.Data, task)
	}
	resp.RespondSuccessJSON(ctx, tasks)
}

// @Router /task/transfer/{id} [get]
// @Summary 获取转移任务状态
// @Description 获取转移任务状态
// @Tags 任务管理
// @Param id path string true "任务 ID"
// @Produces json
func GetTransferTask(ctx *gin.Context) {
	var resp schemas.Response[*task.Task]
	id := ctx.Param("id")
	if id == "" {
		resp.Message = "任务 ID 不能为空"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	t, err := task_controller.GetTransferTask(id)
	if err != nil {
		resp.Message = "获取任务失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}
	resp.RespondSuccessJSON(ctx, t)
}

// @Router /task/transfer/{id} [delete]
// @Summary 取消转移任务
// @Description 取消转移任务
// @Tags 任务管理
// @Param id path string true "任务 ID"
// @Produces json
func CancelTransferTask(ctx *gin.Context) {
	var resp schemas.Response[*task.Task]
	id := ctx.Param("id")
	if id == "" {
		resp.Message = "任务 ID 不能为空"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	task, err := task_controller.CancelTransferTask(id)
	if err != nil {
		resp.Message = "取消任务失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}
	resp.RespondSuccessJSON(ctx, task)
}
