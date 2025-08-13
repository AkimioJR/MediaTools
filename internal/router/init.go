package router

import (
	"MediaTools/internal/router/config"
	"MediaTools/internal/router/library"
	"MediaTools/internal/router/log"
	"MediaTools/internal/router/recognize"
	"MediaTools/internal/router/scrape"
	"MediaTools/internal/router/storage"
	"MediaTools/internal/router/tmdb"
	"net/http"

	_ "MediaTools/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func InitRouter(isDev bool) *gin.Engine {
	logrus.Info("开始初始化路由...")
	ginRouter := gin.Default()

	if isDev {
		ginRouter.GET("/docs", func(ctx *gin.Context) {
			ctx.Redirect(http.StatusFound, "/docs/index.html")
		})
		logrus.Info("开发者模式已启用，开启 Swagger API 路由")
		ginRouter.GET("/docs/*any", gs.WrapHandler(swaggerFiles.Handler))
	}

	config.RegisterConfigRouter(ginRouter)      // 配置相关路由
	log.RegisterLogRouter(ginRouter)            // 日志相关路由
	tmdb.RegisterTMDBRouter(ginRouter)          // TMDB 相关接口
	recognize.RegisteRecognizeRouter(ginRouter) // 识别相关接口
	scrape.RegisterScrapeRouter(ginRouter)      // 刮削相关接口
	library.RegisterLibraryRouter(ginRouter)    // 媒体库相关接口
	storage.RegisterStorageRouter(ginRouter)    // 存储相关接口

	logrus.Info("路由初始化完成")
	return ginRouter
}
