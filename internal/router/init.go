package router

import (
	"MediaTools/internal/router/config"
	"MediaTools/internal/router/library"
	"MediaTools/internal/router/log"
	"MediaTools/internal/router/media"
	"MediaTools/internal/router/scrape"
	"MediaTools/internal/router/storage"
	"MediaTools/internal/router/tmdb"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitRouter() *gin.Engine {
	logrus.Info("开始初始化路由...")
	ginRouter := gin.Default()

	config.RegisterConfigRouter(ginRouter)   // 配置相关路由
	log.RegisterLogRouter(ginRouter)         // 日志相关路由
	tmdb.RegisterTMDBRouter(ginRouter)       // TMDB 相关接口
	media.RegisterMediaRouter(ginRouter)     // 媒体相关接口
	scrape.RegisterScrapeRouter(ginRouter)   // 刮削相关接口
	library.RegisterLibraryRouter(ginRouter) // 媒体库相关接口
	storage.RegisterStorageRouter(ginRouter) // 存储相关接口

	logrus.Info("路由初始化完成")
	return ginRouter
}
