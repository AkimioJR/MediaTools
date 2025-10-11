package router

import (
	"MediaTools/internal/info"
	"MediaTools/internal/router/config"
	"MediaTools/internal/router/history"
	"MediaTools/internal/router/library"
	"MediaTools/internal/router/log"
	"MediaTools/internal/router/recognize"
	"MediaTools/internal/router/runtime"
	"MediaTools/internal/router/scrape"
	"MediaTools/internal/router/storage"
	"MediaTools/internal/router/task"
	"MediaTools/internal/router/tmdb"
	"MediaTools/internal/schemas"
	"net/http"

	_ "MediaTools/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func InitRouter(isDev bool, noRouterHandler gin.HandlerFunc) *gin.Engine {
	logrus.Info("开始初始化路由...")
	ginRouter := gin.Default()
	apiRouter := ginRouter.Group("/api")

	if isDev {
		ginRouter.GET("/docs", func(ctx *gin.Context) {
			ctx.Redirect(http.StatusFound, "/docs/index.html")
		})
		logrus.Info("开发者模式已启用，开启 Swagger API 路由")
		ginRouter.GET("/docs/*any", gs.WrapHandler(swaggerFiles.Handler))
	}

	apiRouter.GET("/version", func(ctx *gin.Context) {
		var resp schemas.Response[*info.VersionInfo]
		resp.RespondSuccessJSON(ctx, &info.Version)
	})
	runtime.RegisterRuntimeRouter(apiRouter.Group("/runtime")) // 程序运行状态相关路由

	config.RegisterConfigRouter(apiRouter.Group("/config"))         // 配置相关路由
	log.RegisterLogRouter(apiRouter.Group("/log"))                  // 日志相关路由
	tmdb.RegisterTMDBRouter(apiRouter.Group("/tmdb"))               // TMDB 相关接口
	recognize.RegisteRecognizeRouter(apiRouter.Group("/recognize")) // 识别相关接口
	scrape.RegisterScrapeRouter(apiRouter.Group("/scrape"))         // 刮削相关接口
	library.RegisterLibraryRouter(apiRouter.Group("/library"))      // 媒体库相关接口
	storage.RegisterStorageRouter(apiRouter.Group("/storage"))      // 存储相关接口
	history.RegisterHistoryRouter(apiRouter.Group("/history"))      // 历史记录相关接口
	task.RegisterTaskRouter(apiRouter.Group("/task"))               // 任务相关接口
	if noRouterHandler != nil {
		ginRouter.NoRoute(noRouterHandler)
	}

	logrus.Info("路由初始化完成")
	return ginRouter
}
