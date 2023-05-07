package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var Deployment deployment

type deployment struct{}

// DeploymentResp 定义列表的返回类型
type DeploymentResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// DeployCreate 定义创建 Deployment 使用的结构体
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Cluster       string            `json:"cluster"`
}

// 从 deployment 类型转到 DataCell 类型
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
		logger.Error(fmt.Sprintf("获取Deployment列表失败, %v", err))
		return nil, errors.New(fmt.Sprintf("获取Deployment列表失败, %v", err))
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
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	// 再排序和分页
	data := filtered.Sort().Paginate()

	deployments := d.fromCells(data.GenericDataList)

	return &DeploymentResp{
		Items: deployments,
		Total: total,
	}, nil
}

// GetDeploymentDetail 获取 deployment 详情
func (d *deployment) GetDeploymentDetail(client *kubernetes.Clientset, deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取Deployment详情失败, %v", err))
		return nil, errors.New(fmt.Sprintf("获取Deployment详情失败, %v", err))
	}

	return deployment, nil
}

// UpdateDeployment 更新 deployment
func (d *deployment) UpdateDeployment(client *kubernetes.Clientset, namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化失败, %v", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v", err))
	}

	_, err = client.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新Deployment失败, %v", err))
		return errors.New(fmt.Sprintf("更新Deployment失败, %v", err))
	}
	return nil
}

// DeleteDeployment 删除 deployment
func (d *deployment) DeleteDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	err = client.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("删除Deployment失败, %v", err))
		return errors.New(fmt.Sprintf("删除Deployment失败, %v", err))
	}

	return nil
}

// 修改Deployment副本数
func (d *deployment) ScaleDeployment(client *kubernetes.Clientset, deploymentName, namespace string, scaleNum int) (replica int32, err error) {
	//获取 aotuscalingv1.Scale 类型的对象，能点出当前的副本数
	scale, err := client.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取Deployment副本信息失败, %v", err))
		return 0, errors.New(fmt.Sprintf("获取Deployment副本信息失败, %v", err))
	}
	// 修改副本数
	scale.Spec.Replicas = int32(scaleNum)
	// 副本副本数，传入 scale 对象
	newScale, err := client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新Deployment副本信息失败, %v", err))
		return 0, errors.New(fmt.Sprintf("更新Deployment副本信息失败, %v", err))
	}
	return newScale.Spec.Replicas, nil
}

// 重启 Deployment
func (d *deployment) RestartDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	// 通过 patch 方法实现重启
	// 此功能等同于 kubectl 命令
	// kubectl deployment ${service} -p \
	// '{"spec":{"template":{"spec":{"containers":[{"name":"'"${service}"'","env":[{"name":"RESTART_","value":"'$(date +%s)'"}]}]}}}}'

	// 使用 patchData Map 组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	// 序列化成 json
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化失败, %v", err))
		return errors.New(fmt.Sprintf("序列化失败, %v", err))
	}
	// 调用 patch 方法更新 deployment
	_, err = client.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName,
		"application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("重启Deployment失败, %v", err))
		return errors.New(fmt.Sprintf("重启Deployment失败, %v", err))
	}
	return nil
}

// 创建 Deployment
func (d *deployment) CreateDeployment(client *kubernetes.Clientset, data *DeployCreate) (err error) {
	// 将 data 中的属性组装成 appsv1.Deployment 对象
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}
	// 判断是否打开健康检查功能，若打开，则定义 ReadinessProbe 和 LivenessProbe
	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					// intstr.IntOrString 的作用是端口可以定义为整行，也可以定义为字符串
					// type=0 表示整行，使用 intVal
					// type=1 表示字符串，使用 strVal
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			// 初始化等待时间
			InitialDelaySeconds: 5,
			// 超时时间
			TimeoutSeconds: 15,
			// 执行间隔
			PeriodSeconds: 5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					// intstr.IntOrString 的作用是端口可以定义为整行，也可以定义为字符串
					// type=0 表示整行，使用 intVal
					// type=1 表示字符串，使用 strVal
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			// 初始化等待时间
			InitialDelaySeconds: 15,
			// 超时时间
			TimeoutSeconds: 15,
			// 执行间隔
			PeriodSeconds: 5,
		}
	}
	// 定义容器的 limit 和 request 资源
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits =
		map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests =
		map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}
	//创建deployment
	_, err = client.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("创建Deployment失败, %v\n", err))
		return errors.New(fmt.Sprintf("创建Deployment失败, %v\n", err))
	}
	return nil
}
