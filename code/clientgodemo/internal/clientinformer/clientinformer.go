package clientinformer

import (
	"log"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func DoInformer(config *rest.Config) {

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	//表示每分钟进行一次resync，resync会周期性地执行List操作

	// sharedInformerFactory 是一个工厂类，它可以用来创建各种类型的 informer 对象

	sharedInformers := informers.NewSharedInformerFactory(clientset, time.Second*10)

	// 每种资源都有自己的informer
	informer := sharedInformers.Core().V1().Pods().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("%s: New Pod Added to Store - %s", mObj.GetNamespace(), mObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(v1.Object)
			nObj := newObj.(v1.Object)
			log.Printf("%s Pod Updated to %s", oObj.GetName(), nObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("%s: Pod Deleted from Store: %s", mObj.GetNamespace(), mObj.GetName())
		},
	})

	informer.Run(stopCh)

}

/*
informer 使用方法：
1. 实例化对应资源的informer
2. 添加回调函数：AddFunc、UpdateFunc、DeleteFunc
3. informer.Run()


以pod资源对象为例，底层代码逻辑：
1. 实例化informer工厂函数，加入初始化参数
2. 通过工厂函数，创建pod资源的informer实例
3. 执行run方法


*/
