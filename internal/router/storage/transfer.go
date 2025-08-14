package storage

import (
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 通用文件传输处理函数
// 根据传输类型执行相应的文件传输操作
// 传输类型可以是复制、移动、硬链接或软链接
// 如果传输类型未知，则返回错误
func handleFileTransfer(ctx *gin.Context, expectedTransferType schemas.TransferType, transferFunc func(*schemas.FileInfo, *schemas.FileInfo) error) {
	var (
		req  schemas.TransferRequest
		resp schemas.Response[*schemas.FileInfo]
	)

	// 绑定并验证请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	// 验证传输类型
	if req.TransferType != expectedTransferType && req.TransferType != schemas.TransferUnknown {
		resp.Message = "传输类型错误: " + req.TransferType.String()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	// 创建源文件和目标文件信息
	srcFile := schemas.NewBasicFileInfo(req.SrcFile.StorageType, req.SrcFile.Path)
	dstFile := schemas.NewBasicFileInfo(req.DstFile.StorageType, req.DstFile.Path)

	// 执行传输操作
	err := transferFunc(srcFile, dstFile)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Success = true
	resp.Data = dstFile
	ctx.JSON(http.StatusOK, resp)
}

// @Route /storage/copy [post]
// @Summary 复制文件
// @Description 将文件从源位置复制到目标位置
// @Tags 存储,存储文件,文件转移
// @Body {object} schemas.TransferRequest true "传输请求"
// @Accept json
// @Products json

func StorageCopyFile(ctx *gin.Context) {
	handleFileTransfer(ctx, schemas.TransferCopy, storage_controller.Copy)
}

// @Route /storage/move [post]
// @Summary 移动文件
// @Description 将文件从源位置移动到目标位置
// @Tags 存储,存储文件,文件转移
// @Body {object} schemas.TransferRequest true "传输请求"
// @Accept json
// @Products json

func StorageMoveFile(ctx *gin.Context) {
	handleFileTransfer(ctx, schemas.TransferMove, storage_controller.Move)
}

// @Route /storage/link [post]
// @Summary 创建硬链接
// @Description 为文件创建硬链接
// @Tags 存储,存储文件,文件转移
// @Body {object} schemas.TransferRequest true "传输请求"
// @Accept json
// @Products json

func StorageLinkFile(ctx *gin.Context) {
	handleFileTransfer(ctx, schemas.TransferLink, storage_controller.Link)
}

// @Route /storage/softlink [post]
// @Summary 创建软链接
// @Description 为文件创建软链接（符号链接）
// @Tags 存储,存储文件,文件转移
// @Body {object} schemas.TransferRequest true "传输请求"
// @Accept json
// @Products json

func StorageSoftLinkFile(ctx *gin.Context) {
	handleFileTransfer(ctx, schemas.TransferSoftLink, storage_controller.SoftLink)
}

// @Route /storage/transfer [post]
// @Summary 通用文件传输接口
// @Description 根据传输类型执行文件传输操作（复制、移动、硬链接、软链接）
// @Tags 存储,存储文件,文件转移
// @Body {object} schemas.TransferRequest true "传输请求"
// @Accept json
// @Products json

func StorageTransferFile(ctx *gin.Context) {
	var (
		req  schemas.TransferRequest
		resp schemas.Response[*schemas.FileInfo]
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	if req.TransferType == schemas.TransferUnknown {
		resp.Message = "传输类型错误: " + req.TransferType.String()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	srcFile := schemas.NewBasicFileInfo(req.SrcFile.StorageType, req.SrcFile.Path)

	dstFile := schemas.NewBasicFileInfo(req.DstFile.StorageType, req.DstFile.Path)

	err := storage_controller.TransferFile(srcFile, dstFile, req.TransferType)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Success = true
	resp.Data = dstFile
	ctx.JSON(http.StatusOK, resp)
}
