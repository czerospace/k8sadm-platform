package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"

	"kubeadm-platform/config"
	"kubeadm-platform/controller"
	"kubeadm-platform/service"
)

func main() {
	// 初始化gin对象
	r := gin.Default()
	// 初始化 K8s client
	service.K8s.Init()
	// 初始化路由
	controller.Router.InitApiRouter(r)
	// 启动gin server
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", err)
		}
	}()
	//优雅关闭server
	//声明一个系统信号的channel，并监听他，如果没有信号，就一直阻塞，如果有，就继续执行
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	//设置ctx超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//cancel用于释放ctx
	defer cancel()
	//关闭gin
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Gin Server关闭异常: ", err)
	}
	logger.Info("Gin Server退出成功")
}
