package dyanmicclient

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func DoDynamicClient(config *rest.Config) {

	dynamicClient, err := dynamic.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	// dynamicClient的唯一关联方法所需的入参
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}

	// 使用dynamicClient的查询列表方法，查询指定namespace下的所有pod，
	// 注意此方法返回的数据结构类型是UnstructuredList
	unstructObj, err := dynamicClient.
		Resource(gvr).
		Namespace("kube-system").
		List(context.TODO(), metav1.ListOptions{Limit: 100})

	if err != nil {
		panic(err.Error())
	}

	// 实例化一个PodList数据结构，用于接收从unstructObj转换后的结果
	podList := &apiv1.PodList{}

	// 转换
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList)

	if err != nil {
		panic(err.Error())
	}

	// 表头
	fmt.Printf("namespace\t status\t\t name\n")

	// 每个pod都打印namespace、status.Phase、name三个字段
	for _, d := range podList.Items {
		fmt.Printf("%v\t %v\t %v\n",
			d.Namespace,
			d.Status.Phase,
			d.Name)
	}
}
