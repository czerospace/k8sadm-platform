package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"kubeadm-platform/config"
)

var K8s k8s

type k8s struct {
	// 提供多集群 client
	ClientMap map[string]*kubernetes.Clientset
	// 提供集群列表功能
	KubeConfMap map[string]string
}

// GetClient 根据集群名获取 client
func (k *k8s) GetClient(cluster string) (*kubernetes.Clientset, error) {
	client, ok := k.ClientMap[cluster]
	if !ok {
		return nil, errors.New(fmt.Sprintf("集群:%s不存在，无法获取client", cluster))
	}
	return client, nil
}

// Init 初始化 client
func (k *k8s) Init() {
	mp := make(map[string]string, 0)
	k.ClientMap = make(map[string]*kubernetes.Clientset, 0)
	// 反序列化
	if err := json.Unmarshal([]byte(config.Kubeconfigs), &mp); err != nil {
		panic(fmt.Sprintf("Kubeconfigs反序列化失败 %v\n", err))
	}
	k.KubeConfMap = mp

	// 初始化 client
	for key, value := range mp {
		conf, err := clientcmd.BuildConfigFromFlags("", value)
		if err != nil {
			panic(fmt.Sprintf("集群%s:创建K8s配置失败 %v", key, err))
		}
		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			panic(fmt.Sprintf("集群%s:创建K8sClient失败 %v", key, err))
		}
		k.ClientMap[key] = clientSet
		logger.Info(fmt.Sprintf("集群%s:创建K8sClient成功", key))
	}
}
