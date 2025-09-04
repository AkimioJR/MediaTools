package router

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandlerWebRouter(source *embed.FS) gin.HandlerFunc {
	frontendFS, err := fs.Sub(source, "web/dist")
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
