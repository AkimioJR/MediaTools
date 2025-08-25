package router

import (
	configuration "MediaTools/internal/config"
	"MediaTools/internal/router/config"
	"MediaTools/internal/router/history"
	"MediaTools/internal/router/library"
	"MediaTools/internal/router/log"
	"MediaTools/internal/router/recognize"
	"MediaTools/internal/router/scrape"
	"MediaTools/internal/router/storage"
	"MediaTools/internal/router/tmdb"
	"MediaTools/internal/schemas"
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

	ginRouter.GET("/version", func(ctx *gin.Context) {
		var resp schemas.Response[*configuration.VersionInfo]
		resp.RespondSuccessJSON(ctx, &configuration.Version)
	})

	config.RegisterConfigRouter(ginRouter.Group("/config"))         // 配置相关路由
	log.RegisterLogRouter(ginRouter.Group("/log"))                  // 日志相关路由
	tmdb.RegisterTMDBRouter(ginRouter.Group("/tmdb"))               // TMDB 相关接口
	recognize.RegisteRecognizeRouter(ginRouter.Group("/recognize")) // 识别相关接口
	scrape.RegisterScrapeRouter(ginRouter.Group("/scrape"))         // 刮削相关接口
	library.RegisterLibraryRouter(ginRouter.Group("/library"))      // 媒体库相关接口
	storage.RegisterStorageRouter(ginRouter.Group("/storage"))      // 存储相关接口
	history.RegisterHistoryRouter(ginRouter.Group("/history"))      // 历史记录相关接口

	logrus.Info("路由初始化完成")
	return ginRouter
}
