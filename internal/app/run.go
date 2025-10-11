package app

import (
	"MediaTools/internal/info"
	"MediaTools/internal/router"
	"MediaTools/web"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func waitSysSign() <-chan os.Signal {
	sysSignCh := make(chan os.Signal, 1)
	signal.Notify(sysSignCh, syscall.SIGINT, syscall.SIGTERM)
	return sysSignCh
}

func handlerWebRouter() gin.HandlerFunc {
	frontendFS, err := fs.Sub(web.WebDist, "dist")
	if err != nil {
		panic("无法加载前端资源: " + err.Error())
	}
	return func(ctx *gin.Context) {
		var (
			sourcePath  string
			contentType string
		)
		switch {
		case ctx.Request.URL.Path == "" || ctx.Request.URL.Path == "/":
			ctx.Redirect(http.StatusFound, "/dashboard")
			return

		case strings.HasPrefix(ctx.Request.URL.Path, "/assets") || ctx.Request.URL.Path == "/vite.svg":
			sourcePath = strings.TrimPrefix(ctx.Request.URL.Path, "/")
			switch {
			case strings.HasSuffix(sourcePath, ".js"):
				contentType = "application/javascript"

			case strings.HasSuffix(sourcePath, ".css"):
				contentType = "text/css"

			case strings.HasSuffix(sourcePath, ".svg"):
				contentType = "image/svg+xml"
			}

		default:
			sourcePath = "index.html"
			contentType = "text/html"
		}

		data, err := fs.ReadFile(frontendFS, sourcePath)
		if err != nil {
			ctx.String(http.StatusNotFound, "资源不存在: "+err.Error())
			return
		}

		ctx.Data(http.StatusOK, contentType, data)
	}
}

func runServer() {
	ginR := router.InitRouter(info.RuntimeAppStatus.IsDev, handlerWebRouter())
	sysCh := waitSysSign()
	errCh := make(chan error, 1)
	go func() {
		err := ginR.Run(":" + strconv.Itoa(int(info.RuntimeAppStatus.Port)))
		if err != nil {
			errCh <- fmt.Errorf("启动服务器失败: %v", err)
		}
	}()
	select {
	case err := <-errCh:
		if err != nil {
			logrus.Errorf("应用程序运行中发生错误: %v", err)
		}
	case sig := <-sysCh:
		logrus.Infof("收到系统信号: %v, 退出应用程序", sig)
	}
}
