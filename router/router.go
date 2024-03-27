package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rancher-setup/app/controller"
	"github.com/rancher-setup/app/middleware"
	"github.com/rancher-setup/internal/server"
	"io"
	"net/http"
)

type AppRouter struct {
	server server.HttpServer
}

func New(server server.HttpServer) *AppRouter {
	server.SetMiddleware(&middleware.Cors{})
	return &AppRouter{
		server,
	}
}

func (*AppRouter) Add(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		resp, err := http.Get("https://wangsiye-test.oss-cn-beijing.aliyuncs.com/web/dist/index.html")
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to fetch file: %v", err)
			return
		}
		defer resp.Body.Close()

		// 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			ctx.String(resp.StatusCode, "Failed to fetch file: %s", resp.Status)
			return
		}

		// 将文件内容写入响应
		_, err = io.Copy(ctx.Writer, resp.Body)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to write file to response: %v", err)
			return
		}
	})

	index := &controller.Index{}
	server.POST("/install", index.InstallRancher)
	server.GET("/config", index.GetRancherConfig)
	server.GET("/rancherState", index.GetRancherState)
}
