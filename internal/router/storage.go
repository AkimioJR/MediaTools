package router

import (
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/schemas"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /storage
// @Route /list [get]
// @Summary 获取存储提供者列表
// @Description 返回所有已注册的存储提供者列表
// @Tags storage
// @Products json
// @Success 200 {object} Response[[]schemas.StorageProviderItem]
// @Failure 500 {object} Response
func StorageProviderList(ctx *gin.Context) {
	var resp Response[[]schemas.StorageProviderItem]

	resp.Data = storage_controller.ListStorageProviders()
	resp.Success = true
	logrus.Info(resp)
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage
// @Route /:storage_type/info [get]
// @Summary 获取文件/目录信息
// @Description 根据路径和存储类型获取文件或目录的详细信息
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件或目录路径"
// @Products json
// @Success 200 {object} Response[*schemas.FileInfo]
// @Failure 400 {object} Response[*schemas.FileInfo]
// @Failure 500 {object} Response[*schemas.FileInfo]
func StorageGetFileInfo(ctx *gin.Context) {
	var resp Response[*schemas.FileInfo]

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		resp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	fileInfo, err := storage_controller.GetFile(path, storageType)
	if err != nil {
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = fileInfo
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage
// @Route /:storage_type/exists [get]
// @Summary 检查文件/目录是否存在
// @Description 根据路径和存储类型检查文件或目录是否存在
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件或目录路径"
// @Products json
// @Success 200 {object} Response[bool]
// @Failure 400 {object} Response[bool]
// @Failure 500 {object} Response[bool]
func StorageCheckExists(ctx *gin.Context) {
	var resp Response[bool]

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		resp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	fileInfo := schemas.FileInfo{
		StorageType: storageType,
		Path:        path,
	}

	exists, err := storage_controller.Exists(&fileInfo)
	if err != nil {
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = exists
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage
// @Route /:storage_type/list [get]
// @Summary 列出目录内容
// @Description 根据存储类型和路径列出目录下的所有文件和子目录
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path query string true "目录路径"
// @Products json
// @Success 200 {object} Response[[]schemas.FileInfo]
// @Failure 400 {object} Response[[]schemas.FileInfo]
// @Failure 500 {object} Response[[]schemas.FileInfo]
func StorageList(ctx *gin.Context) {
	var resp Response[[]schemas.FileInfo]

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	path := ctx.Query("path")
	if path == "" {
		resp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	dirInfo := &schemas.FileInfo{
		StorageType: storageType,
		Path:        path,
		IsDir:       true,
	}

	files, err := storage_controller.List(dirInfo)
	if err != nil {

		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = files
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage
// @Route /:storage_type/mkdir [post]
// @Summary 创建目录
// @Description 根据存储类型和路径创建一个新目录
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Body {object} PathRequest true "目录路径"
// @Accept json
// @Products json
// @Success 200 {object} Response[schemas.FileInfo]
// @Failure 400 {object} Response[schemas.FileInfo]
// @Failure 500 {object} Response[schemas.FileInfo]
func StorageMkdir(ctx *gin.Context) {
	var (
		req  PathRequest
		resp Response[schemas.FileInfo]
	)

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	dirInfo := schemas.FileInfo{
		StorageType: storageType,
		Path:        req.Path,
		IsDir:       true,
	}

	err := storage_controller.Mkdir(&dirInfo)
	if err != nil {
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = dirInfo
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage
// @Route /:storage_type/delete [delete]
// @Summary 删除文件或目录
// @Description 根据存储类型和路径删除指定的文件或目录
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Body {object} PathRequest true "文件或目录路径"
// @Accept json
// @Products json
// @Success 200 {object} Response[schemas.FileInfo]
// @Failure 400 {object} Response[schemas.FileInfo]
// @Failure 500 {object} Response[schemas.FileInfo]
func StorageDelete(ctx *gin.Context) {
	var (
		req  PathRequest
		resp Response[schemas.FileInfo]
	)

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	fileInfo := schemas.FileInfo{
		StorageType: storageType,
		Path:        req.Path,
	}

	err := storage_controller.Delete(&fileInfo)
	if err != nil {

		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = fileInfo
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage
// @Route /:storage_type/upload [post]
// @Summary 上传文件
// @Description 根据存储类型和路径上传文件
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path formData string true "上传路径"
// @Param file formData file true "上传文件"
// @Accept multipart/form-data
// @Products json
// @Success 200 {object} Response[schemas.FileInfo]
// @Failure 400 {object} Response[schemas.FileInfo]
// @Failure 500 {object} Response[schemas.FileInfo]
func StorageUploadFile(ctx *gin.Context) {
	var resp Response[schemas.FileInfo]

	storageTypeStr := ctx.Param("storage_type")
	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	path := ctx.PostForm("path")
	if path == "" {
		resp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		resp.Message = "获取上传文件失败: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	src, err := file.Open()
	if err != nil {
		resp.Message = "打开上传文件失败: " + err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer src.Close()

	fileInfo := schemas.FileInfo{
		StorageType: storageType,
		Path:        path,
		Size:        file.Size,
	}

	err = storage_controller.CreateFile(&fileInfo, src)
	if err != nil {
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = fileInfo
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage
// @Route /:storage_type/download [get]
// @Summary 下载文件
// @Description 根据存储类型和路径下载文件
// @Tags storage
// @Param storage_type path string true "存储类型"
// @Param path query string true "文件路径"
// @Failure 400 {object} Response[schemas.FileInfo]
// @Failure 500 {object} Response[schemas.FileInfo]
func StorageDownloadFile(ctx *gin.Context) {
	var resp Response[*schemas.FileInfo]

	path := ctx.Query("path")
	storageTypeStr := ctx.Param("storage_type")

	if path == "" {
		resp.Message = "路径不能为空"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	storageType := schemas.ParseStorageType(storageTypeStr)
	if storageType == schemas.StorageUnknown {
		resp.Message = "无效的存储类型: " + storageTypeStr
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	fileInfo := schemas.FileInfo{
		StorageType: storageType,
		Path:        path,
	}

	reader, err := storage_controller.ReadFile(&fileInfo)
	if err != nil {
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	defer reader.Close()

	// 设置下载响应头
	ctx.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	ctx.Header("Content-Type", "application/octet-stream")

	// 流式传输文件内容
	_, err = io.Copy(ctx.Writer, reader)
	if err != nil {
		logrus.Error("下载文件时发生错误: ", err)
	}
}

// 通用文件传输处理函数
// 根据传输类型执行相应的文件传输操作
// 传输类型可以是复制、移动、硬链接或软链接
// 如果传输类型未知，则返回错误
func handleFileTransfer(ctx *gin.Context, expectedTransferType schemas.TransferType, transferFunc func(*schemas.FileInfo, *schemas.FileInfo) error) {
	var (
		req  TransferRequest
		resp Response[*schemas.FileInfo]
	)

	// 绑定并验证请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 验证传输类型
	if req.TransferType != expectedTransferType && req.TransferType != schemas.TransferUnknown {
		resp.Message = "传输类型错误: " + req.TransferType.String()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// 创建源文件和目标文件信息
	srcFile := schemas.FileInfo{
		StorageType: req.SrcFile.StorageType,
		Path:        req.SrcFile.Path,
	}

	dstFile := schemas.FileInfo{
		StorageType: req.DstFile.StorageType,
		Path:        req.DstFile.Path,
	}

	// 执行传输操作
	err := transferFunc(&srcFile, &dstFile)
	if err != nil {
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	// 返回成功响应
	resp.Data = &dstFile
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}

// @BasePath /storage
// @Route /copy [post]
// @Summary 复制文件
// @Description 将文件从源位置复制到目标位置
// @Tags storage
// @Body {object} TransferRequest true "传输请求"
// @Accept json
// @Products json
// @Success 200 {object} Response[*schemas.FileInfo]
// @Failure 400 {object} Response[*schemas.FileInfo]
// @Failure 500 {object} Response[*schemas.FileInfo]
func StorageCopyFile(ctx *gin.Context) {
	handleFileTransfer(ctx, schemas.TransferCopy, storage_controller.Copy)
}

// @BasePath /storage
// @Route /move [post]
// @Summary 移动文件
// @Description 将文件从源位置移动到目标位置
// @Tags storage
// @Body {object} TransferRequest true "传输请求"
// @Accept json
// @Products json
// @Success 200 {object} Response[*schemas.FileInfo]
// @Failure 400 {object} Response[*schemas.FileInfo]
// @Failure 500 {object} Response[*schemas.FileInfo]
func StorageMoveFile(ctx *gin.Context) {
	handleFileTransfer(ctx, schemas.TransferMove, storage_controller.Move)
}

// @BasePath /storage
// @Route /link [post]
// @Summary 创建硬链接
// @Description 为文件创建硬链接
// @Tags storage
// @Body {object} TransferRequest true "传输请求"
// @Accept json
// @Products json
// @Success 200 {object} Response[*schemas.FileInfo]
// @Failure 400 {object} Response[*schemas.FileInfo]
// @Failure 500 {object} Response[*schemas.FileInfo]
func StorageLinkFile(ctx *gin.Context) {
	handleFileTransfer(ctx, schemas.TransferLink, storage_controller.Link)
}

// @BasePath /storage
// @Route /softlink [post]
// @Summary 创建软链接
// @Description 为文件创建软链接（符号链接）
// @Tags storage
// @Body {object} TransferRequest true "传输请求"
// @Accept json
// @Products json
// @Success 200 {object} Response[*schemas.FileInfo]
// @Failure 400 {object} Response[*schemas.FileInfo]
// @Failure 500 {object} Response[*schemas.FileInfo]
func StorageSoftLinkFile(ctx *gin.Context) {
	handleFileTransfer(ctx, schemas.TransferSoftLink, storage_controller.SoftLink)
}

// @BasePath /storage
// @Route /transfer [post]
// @Summary 通用文件传输接口
// @Description 根据传输类型执行文件传输操作（复制、移动、硬链接、软链接）
// @Tags storage
// @Body {object} TransferRequest true "传输请求"
// @Accept json
// @Products json
// @Success 200 {object} Response[*schemas.FileInfo]
// @Failure 400 {object} Response[*schemas.FileInfo]
// @Failure 500 {object} Response[*schemas.FileInfo]
func StorageTransferFile(ctx *gin.Context) {
	var (
		req  TransferRequest
		resp Response[*schemas.FileInfo]
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Message = "请求参数错误: " + err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if req.TransferType == schemas.TransferUnknown {
		resp.Message = "传输类型错误: " + req.TransferType.String()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	srcFile := schemas.FileInfo{
		StorageType: req.SrcFile.StorageType,
		Path:        req.SrcFile.Path,
	}

	dstFile := schemas.FileInfo{
		StorageType: req.DstFile.StorageType,
		Path:        req.DstFile.Path,
	}

	err := storage_controller.TransferFile(&srcFile, &dstFile, req.TransferType)
	if err != nil {
		resp.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Data = &dstFile
	resp.Success = true
	ctx.JSON(http.StatusOK, resp)
}
