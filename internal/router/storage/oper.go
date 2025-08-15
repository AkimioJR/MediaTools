package storage

import (
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/schemas"
	"MediaTools/internal/schemas/storage"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Route /storage/:storage_type/info [get]
// @Summary 获取文件/目录信息
// @Description 根据路径和存储类型获取文件或目录的详细信息
// @Tags 存储,存储文件
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件或目录路径"
// @Products json
func StorageGetFileInfo(ctx *gin.Context) {
	var resp schemas.Response[*storage.StorageFileInfo]

	storageTypeStr := ctx.Param("storage_type")
	storageType := storage.ParseStorageType(storageTypeStr)
	if storageType == storage.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		resp.Message = "路径不能为空"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	filePath, err := storage_controller.GetPath(path, storageType)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	fileInfo, err := storage_controller.GetDetail(filePath)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Data = fileInfo
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/:storage_type/exists [get]
// @Summary 检查文件/目录是否存在
// @Description 根据路径和存储类型检查文件或目录是否存在
// @Tags 存储,存储文件
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件或目录路径"
// @Products json
func StorageCheckExists(ctx *gin.Context) {
	var resp schemas.Response[*bool]

	storageTypeStr := ctx.Param("storage_type")
	storageType := storage.ParseStorageType(storageTypeStr)
	if storageType == storage.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		resp.Message = "路径不能为空"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	exists, err := storage_controller.Exists(storage.NewStoragePath(storageType, path))
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Data = &exists
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/:storage_type/list [get]
// @Summary 列出目录内容
// @Description 根据存储类型和路径列出目录下的所有文件和子目录
// @Tags 存储,存储文件
// @Param storage_type path string true "存储类型"
// @Param path query string true "目录路径"
// @Products json
func StorageList(ctx *gin.Context) {
	var resp schemas.Response[[]storage.StorageFileInfo]

	storageTypeStr := ctx.Param("storage_type")
	storageType := storage.ParseStorageType(storageTypeStr)
	if storageType == storage.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		resp.Message = "路径不能为空"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	files, err := storage_controller.List(storage.NewStoragePath(storageType, path))
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	resp.Data = files
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/:storage_type/mkdir [post]
// @Summary 创建目录
// @Description 根据存储类型和路径创建一个新目录
// @Tags 存储,存储文件
// @Param storage_type path string true "存储类型"
// @Body {object} schemas.PathRequest true "目录路径"
// @Accept json
// @Products json
func StorageMkdir(ctx *gin.Context) {
	var (
		req  schemas.PathRequest
		resp schemas.Response[*string]
	)

	storageTypeStr := ctx.Param("storage_type")
	storageType := storage.ParseStorageType(storageTypeStr)
	if storageType == storage.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	dirPath := storage.NewStoragePath(storageType, req.Path)
	err := storage_controller.Mkdir(dirPath)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	p := dirPath.GetPath()
	resp.Data = &p
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/:storage_type/delete [delete]
// @Summary 删除文件或目录
// @Description 根据存储类型和路径删除指定的文件或目录
// @Tags 存储,存储文件
// @Param storage_type path string true "存储类型"
// @Body {object} schemas.PathRequest true "文件或目录路径"
// @Accept json
// @Products json
func StorageDelete(ctx *gin.Context) {
	var (
		req  schemas.PathRequest
		resp schemas.Response[*string]
	)

	storageTypeStr := ctx.Param("storage_type")
	storageType := storage.ParseStorageType(storageTypeStr)
	if storageType == storage.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	path := storage.NewStoragePath(storageType, req.Path)

	err := storage_controller.Delete(path)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	p := path.GetPath()
	resp.Data = &p
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/:storage_type/rename [post]
// @Summary 重命名文件或目录
// @Description 根据存储类型和路径重命名指定的文件或目录
// @Tags 存储,存储文件
// @Param storage_type path string true "存储类型"
// @Body {object} schemas.RenameRequest true "重命名请求"
// @Accept json
// @Products json
func StorageRename(ctx *gin.Context) {
	var (
		req  schemas.RenameRequest
		resp schemas.Response[*string]
	)

	storageTypeStr := ctx.Param("storage_type")
	storageType := storage.ParseStorageType(storageTypeStr)
	if storageType == storage.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	path := storage.NewStoragePath(storageType, req.Path)

	err := storage_controller.Rename(path, req.NewName)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	p := path.GetPath()
	resp.Data = &p
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/:storage_type/upload [post]
// @Summary 上传文件
// @Description 根据存储类型和路径上传文件
// @Tags 存储,存储文件
// @Param storage_type path string true "存储类型"
// @Param path formData string true "上传路径"
// @Param file formData file true "上传文件"
// @Accept multipart/form-data
// @Products json
func StorageUploadFile(ctx *gin.Context) {
	var resp schemas.Response[*string]

	storageTypeStr := ctx.Param("storage_type")
	storageType := storage.ParseStorageType(storageTypeStr)
	if storageType == storage.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	path := ctx.PostForm("path")
	if path == "" {
		resp.Message = "路径不能为空"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		resp.Message = "获取上传文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	src, err := file.Open()
	if err != nil {
		resp.Message = "打开上传文件失败: " + err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}
	defer src.Close()

	filePath := storage.NewStoragePath(storageType, path)

	err = storage_controller.CreateFile(filePath, src)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}

	p := filePath.GetPath()
	resp.Data = &p
	resp.RespondJSON(ctx, http.StatusOK)
}

// @Route /storage/:storage_type/download [get]
// @Summary 下载文件
// @Description 根据存储类型和路径下载文件
// @Tags 存储,存储文件
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件路径"
// @Produce application/octet-stream
// @Success 200 {file} file "文件下载成功"
// @Header 200 {string} Content-Disposition "文件下载头，格式：attachment; filename=文件名"
// @Header 200 {string} Content-Type "application/octet-stream"
func StorageDownloadFile(ctx *gin.Context) {
	var resp schemas.Response[struct{}]

	path := ctx.Query("path")
	storageTypeStr := ctx.Param("storage_type")

	if path == "" {
		resp.Message = "路径不能为空"
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	storageType := storage.ParseStorageType(storageTypeStr)
	if storageType == storage.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		resp.RespondJSON(ctx, http.StatusBadRequest)
		return
	}

	filePath := storage.NewStoragePath(storageType, path)

	reader, err := storage_controller.ReadFile(filePath)
	if err != nil {
		resp.Message = err.Error()
		resp.RespondJSON(ctx, http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	// 设置下载响应头
	ctx.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath.GetPath()))
	ctx.Header("Content-Type", "application/octet-stream")

	// 流式传输文件内容
	_, err = io.Copy(ctx.Writer, reader)
	if err != nil {
		logrus.Error("下载文件时发生错误: ", err)
	}
}
