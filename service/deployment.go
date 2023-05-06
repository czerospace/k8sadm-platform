package service

import (
	"context"
	"fmt"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Deployment deployment

type deployment struct{}

// DeploymentResp 定义列表的返回类型
type DeploymentResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// // 从 deployment 类型转到 DataCell 类型
func (d *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])
	}
	return cells
}

// 从 DataCell 类型转到 deployment 类型
func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	pods := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		pods[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return pods
}

// GetDeployments 获取 deployment 列表
func (d *deployment) GetDeployments(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (deploymentResp *DeploymentResp, err error) {
	deploymentList, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取Deployment列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取Deployment列表失败, %v\n", err))
	}
	//实例化dataSelector对象
	selectableData := &dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
		dataSelectorQuery: &DataSelectorQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	//再排序和分页
	data := filtered.Sort().Paginate()

	deployments := d.fromCells(data.GenericDataList)

	return &DeploymentResp{
		Items: deployments,
		Total: total,
	}, nil
}

// 获取deployment详情
func (d *deployment) GetDeploymentDetail(client *kubernetes.Clientset, deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取Deployment详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取Deployment详情失败, %v\n", err))
	}

	return deployment, nil
}

// 更新deployment
func (d *deployment) UpdateDeployment(client *kubernetes.Clientset, namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}

	_, err = client.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新Deployment失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新Deployment失败, %v\n", err))
	}
	return nil
}

// 删除deployment
func (d *deployment) DeleteDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	err = client.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("删除Deployment失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除Deployment失败, %v\n", err))
	}

	return nil
}
