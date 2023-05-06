package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"

	"kubeadm-platform/service"
)

var Pod pod

type pod struct{}

// GetPods 获取pod列表
func (p *pod) GetPods(ctx *gin.Context) {
	// 接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Cluster    string `form:"cluster"`
	})
	// 绑定参数
	// form 格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v", err),
			"data": nil,
		})
		return
	}
	// 获取 client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// 调用 service 方法，获取列表
	data, err := service.Pod.GetPods(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			//"code": 90500, // 业务状态
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": data,
	})
}

// GetPodDetail 获取 pod详 情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})
	// 绑定参数
	// form 格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v", err),
			"data": nil,
		})
		return
	}
	//获取 client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用 service 方法，获取列表
	data, err := service.Pod.GetPodDetail(client, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod详情成功",
		"data": data,
	})
}

// DeletePod 删除 pod
func (p *pod) DeletePod(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})
	// 绑定参数
	// form格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// 调用 service 方法，获取列表
	err = service.Pod.DeletePod(client, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Pod成功",
		"data": nil,
	})
}

// UpdatePod 更新 pod
func (p *pod) UpdatePod(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
		Cluster   string `json:"cluster"`
	})
	// 绑定参数
	// form格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v", err),
			"data": nil,
		})
		return
	}
	// 获取 client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	err = service.Pod.UpdatePod(client, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新Pod成功",
		"data": nil,
	})
}

// GetPodContainer 获取 Pod 容器
func (p *pod) GetPodContainer(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})
	// 绑定参数
	// form格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v", err),
			"data": nil,
		})
		return
	}
	// 获取 client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	data, err := service.Pod.GetPodContainer(client, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod容器成功",
		"data": data,
	})
}

// GetPodLog 获取pod容器日志
func (p *pod) GetPodLog(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		ContainerName string `form:"container_name"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
		Cluster       string `form:"cluster"`
	})
	// 绑定参数
	// form 格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
	if err := ctx.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v", err),
			"data": nil,
		})
		return
	}
	// 获取 client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// 调用 service 方法，获取列表
	data, err := service.Pod.GetPodLog(client, params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod容器日志成功",
		"data": data,
	})
}
