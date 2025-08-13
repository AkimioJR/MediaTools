package storage

import (
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/schemas"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Route /storage/:storage_type/info [get]
// @Summary 获取文件/目录信息
// @Description 根据路径和存储类型获取文件或目录的详细信息
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件或目录路径"
// @Products json
// @Success 200 {object} schemas.FileInfo
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func StorageGetFileInfo(ctx *gin.Context) {
	var errResp schemas.ErrResponse

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		errResp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	fileInfo, err := storage_controller.GetFile(path, storageType)
	if err != nil {
		errResp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, fileInfo)
}

// @Route /storage/:storage_type/exists [get]
// @Summary 检查文件/目录是否存在
// @Description 根据路径和存储类型检查文件或目录是否存在
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件或目录路径"
// @Products json
// @Success 200 {object} bool
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func StorageCheckExists(ctx *gin.Context) {
	var errResp schemas.ErrResponse

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		errResp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	fileInfo := schemas.NewBasicFileInfo(storageType, path)

	exists, err := storage_controller.Exists(fileInfo)
	if err != nil {
		errResp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, exists)
}

// @Route /storage/:storage_type/list [get]
// @Summary 列出目录内容
// @Description 根据存储类型和路径列出目录下的所有文件和子目录
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path query string true "目录路径"
// @Products json
// @Success 200 {object} []schemas.FileInfo
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func StorageList(ctx *gin.Context) {
	var errResp schemas.ErrResponse

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		errResp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	dirInfo := schemas.NewBasicFileInfo(storageType, path)
	dirInfo.IsDir = true

	files, err := storage_controller.List(dirInfo)
	if err != nil {

		errResp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, files)
}

// @Route /storage/:storage_type/mkdir [post]
// @Summary 创建目录
// @Description 根据存储类型和路径创建一个新目录
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Body {object} schemas.PathRequest true "目录路径"
// @Accept json
// @Products json
// @Success 200 {object} schemas.FileInfo
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func StorageMkdir(ctx *gin.Context) {
	var (
		req     schemas.PathRequest
		errResp schemas.ErrResponse
	)

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errResp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	dirInfo := schemas.NewBasicFileInfo(storageType, req.Path)
	dirInfo.IsDir = true

	err := storage_controller.Mkdir(dirInfo)
	if err != nil {
		errResp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, dirInfo)
}

// @Route /storage/:storage_type/delete [delete]
// @Summary 删除文件或目录
// @Description 根据存储类型和路径删除指定的文件或目录
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Body {object} schemas.PathRequest true "文件或目录路径"
// @Accept json
// @Products json
// @Success 200 {object} schemas.FileInfo
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func StorageDelete(ctx *gin.Context) {
	var (
		req     schemas.PathRequest
		errResp schemas.ErrResponse
	)

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errResp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	fileInfo := schemas.NewBasicFileInfo(storageType, req.Path)

	err := storage_controller.Delete(fileInfo)
	if err != nil {

		errResp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, fileInfo)
}

// @Route /storage/:storage_type/rename [post]
// @Summary 重命名文件或目录
// @Description 根据存储类型和路径重命名指定的文件或目录
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Body {object} schemas.RenameRequest true "重命名请求"
// @Accept json
// @Products json
// @Success 200 {object} schemas.FileInfo
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func StorageRename(ctx *gin.Context) {
	var (
		req     schemas.RenameRequest
		errResp schemas.ErrResponse
	)

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errResp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	fileInfo := schemas.NewBasicFileInfo(storageType, req.Path)

	err := storage_controller.Rename(fileInfo, req.NewName)
	if err != nil {
		errResp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, fileInfo)
}

// @Route /storage/:storage_type/upload [post]
// @Summary 上传文件
// @Description 根据存储类型和路径上传文件
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path formData string true "上传路径"
// @Param file formData file true "上传文件"
// @Accept multipart/form-data
// @Products json
// @Success 200 {object} schemas.FileInfo
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func StorageUploadFile(ctx *gin.Context) {
	var errResp schemas.ErrResponse

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	path := ctx.PostForm("path")
	if path == "" {
		errResp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		errResp.Message = "获取上传文件失败: " + err.Error()
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	src, err := file.Open()
	if err != nil {
		errResp.Message = "打开上传文件失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}
	defer src.Close()

	fileInfo := schemas.NewBasicFileInfo(storageType, path)
	fileInfo.Size = file.Size

	err = storage_controller.CreateFile(fileInfo, src)
	if err != nil {
		errResp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	ctx.JSON(http.StatusOK, fileInfo)
}

// @Route /storage/:storage_type/download [get]
// @Summary 下载文件
// @Description 根据存储类型和路径下载文件
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件路径"
// @Produce application/octet-stream
// @Success 200 {file} file "文件下载成功"
// @Header 200 {string} Content-Disposition "文件下载头，格式：attachment; filename=文件名"
// @Header 200 {string} Content-Type "application/octet-stream"
// @Failure 400 {object} schemas.ErrResponse
// @Failure 500 {object} schemas.ErrResponse
func StorageDownloadFile(ctx *gin.Context) {
	var errResp schemas.ErrResponse

	path := ctx.Query("path")
	storageTypeStr := ctx.Param("storage_type")

	if path == "" {
		errResp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		errResp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	fileInfo := schemas.NewBasicFileInfo(storageType, path)

	reader, err := storage_controller.ReadFile(fileInfo)
	if err != nil {
		errResp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}
	defer reader.Close()

	// 设置下载响应头
	ctx.Header("Content-Disposition", "attachment; filename="+fileInfo.Name)
	ctx.Header("Content-Type", "application/octet-stream")

	// 流式传输文件内容
	_, err = io.Copy(ctx.Writer, reader)
	if err != nil {
		logrus.Error("下载文件时发生错误: ", err)
	}
}
