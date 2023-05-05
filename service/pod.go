package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/wonderivan/logger"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Pod pod

type pod struct {
}

// PodsResp 定义列表的返回类型
type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

// 从 Pod 类型转到 DataCell 类型
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

// 从 DataCell 类型转到 pod 类型
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

// GetPods 获取 pod 列表
// client 用于选择哪个集群
func (p *pod) GetPods(client *kubernetes.Clientset, fileterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	// context.TODO() 用于声明一个空的 context 上下文，用于 List 方法内设置这个请求的超时（源码），这里的常用用法
	// metav1.ListOptions{} 用于过滤 List 数据，如使用 label , field 等
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取Pod列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取Pod列表失败, %v\n", err))
	}

	// 实例化 dataSelector 对象
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
		dataSelectorQuery: &DataSelectorQuery{
			FilterQuery: &FilterQuery{Name: fileterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	// 再排序和分页
	data := filtered.Sort().Paginate()
	// 将 []DataCell 类型的 pod 列表转为 v1.pod 列表
	pods := p.fromCells(data.GenericDataList)

	return &PodsResp{
		Items: pods,
		Total: total,
	}, nil
}
