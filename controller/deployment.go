package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"

	"kubeadm-platform/service"
)

var Deployment deployment

type deployment struct{}

// GetDeployments 获取 deployment 列表
func (d *deployment) GetDeployments(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Cluster    string `form:"cluster"`
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
	// 调用 service 方法，获取列表
	data, err := service.Deployment.GetDeployments(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Deployment列表成功",
		"data": data,
	})
}

// GetDeploymentDetail 获取 deployment 详情
func (d *deployment) GetDeploymentDetail(ctx *gin.Context) {
	//接收参数,匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
		Cluster        string `form:"cluster"`
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
	data, err := service.Deployment.GetDeploymentDetail(client, params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Deployment详情成功",
		"data": data,
	})
}

// DeleteDeployment 删除 deployment
func (d *deployment) DeleteDeployment(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
		Cluster        string `json:"cluster"`
	})
	// 绑定参数
	// form 格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
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
	// 调用 service 方法，获取列表
	err = service.Deployment.DeleteDeployment(client, params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Deployment成功",
		"data": nil,
	})
}

// UpdateDeployment 更新 deployment
func (d *deployment) UpdateDeployment(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
		Cluster   string `json:"cluster"`
	})
	// 绑定参数
	// form 格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
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
	// 调用 service 方法，获取列表
	err = service.Deployment.UpdateDeployment(client, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新Deployment成功",
		"data": nil,
	})
}

// ScaleDeployment 调整 deployment 副本数
func (d *deployment) ScaleDeployment(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		ScaleNum       int    `json:"scale_num"`
		Namespace      string `json:"namespace"`
		Cluster        string `json:"cluster"`
	})
	// 绑定参数
	// form 格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
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
	// 调用 service 方法，获取列表
	data, err := service.Deployment.ScaleDeployment(client, params.DeploymentName, params.Namespace, params.ScaleNum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "调整Deployment副本数成功",
		"data": data,
	})
}

// RestartDeployment 重启 deployment
func (d *deployment) RestartDeployment(ctx *gin.Context) {
	// 接收参数,匿名结构体，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
		Cluster        string `json:"cluster"`
	})
	// 绑定参数
	// form 格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
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
	// 调用 service 方法，获取列表
	err = service.Deployment.RestartDeployment(client, params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "重启Deployment成功",
		"data": nil,
	})
}

// CreateDeployment 创建 deployment
func (d *deployment) CreateDeployment(ctx *gin.Context) {
	var (
		deployCreate = new(service.DeployCreate)
		err          error
	)
	// 绑定参数
	// form 格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
	if err := ctx.ShouldBindJSON(deployCreate); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v", err),
			"data": nil,
		})
		return
	}
	// 获取 client
	client, err := service.K8s.GetClient(deployCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	// 调用 service 方法，获取列表
	err = service.Deployment.CreateDeployment(client, deployCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建Deployment成功",
		"data": nil,
	})
}
