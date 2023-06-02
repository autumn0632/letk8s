package config

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewConfig(path string) (*rest.Config, error) {
	// 从本机加载kubeconfig配置文件，因此第一个参数为空字符串
	return clientcmd.BuildConfigFromFlags("", path)
}
