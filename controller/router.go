package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Router 实例化router对象，可使用该对象点出首字母大写的方法（跨包调用）
var Router router

// 定义router结构体
type router struct{}

// InitApiRouter 初始化路由，创建测试api接口
func (r *router) InitApiRouter(g *gin.Engine) {
	g.GET("/testapi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "testapi success!",
			"data": nil,
		})
	})
	g.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
}